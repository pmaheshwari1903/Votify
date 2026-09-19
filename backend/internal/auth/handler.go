package auth

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/votify/backend/internal/config"
	"github.com/votify/backend/internal/errors"
	"github.com/votify/backend/internal/response"
)

type AuthHandler struct {
	service AuthService
}

func NewAuthHandler(service AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid registration payload", map[string]string{"body": err.Error()})
		return
	}

	res, err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	response.Created(c, res)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid login payload", map[string]string{"body": err.Error()})
		return
	}

	res, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	response.OK(c, res)
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		userID = c.GetHeader("X-User-ID")
	}

	if userID == "" {
		response.Unauthorized(c, "User context missing")
		return
	}

	res, err := h.service.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	response.OK(c, res)
}

func (h *AuthHandler) OAuthLogin(c *gin.Context) {
	cfg := config.Load()
	ctx := c.Request.Context()
	isSecure := cfg.Env == "production"

	meta, err := FetchOIDCDiscovery(ctx, cfg.OAuthIssuer)
	authEndpoint := ""
	if err == nil && meta != nil {
		authEndpoint = meta.AuthorizationEndpoint
	}
	if authEndpoint == "" {
		authEndpoint = strings.TrimRight(cfg.OAuthIssuer, "/") + "/authorize"
	}

	stateRandom := GenerateSecureRandomState()
	codeVerifier := GeneratePKCEVerifier()
	codeChallenge := ComputePKCEChallengeS256(codeVerifier)

	// Combine random state & PKCE verifier into state payload so verifier survives cross-site redirects
	fullState := fmt.Sprintf("%s.%s", stateRandom, codeVerifier)

	if isSecure {
		c.SetSameSite(http.SameSiteNoneMode)
	} else {
		c.SetSameSite(http.SameSiteLaxMode)
	}

	c.SetCookie("oauth_state", stateRandom, 600, "/", "", isSecure, true)
	c.SetCookie("oauth_code_verifier", codeVerifier, 600, "/", "", isSecure, true)
	if c.Query("popup") == "true" {
		c.SetCookie("oauth_popup", "true", 600, "/", "", isSecure, false)
	}

	u, err := url.Parse(authEndpoint)
	if err != nil {
		response.BadRequest(c, "Invalid authorization endpoint", nil)
		return
	}

	q := u.Query()
	q.Set("response_type", "code")
	q.Set("client_id", cfg.OAuthClientID)
	q.Set("redirect_uri", cfg.OAuthRedirectURI)
	q.Set("scope", "openid profile email")
	q.Set("state", fullState)
	q.Set("code_challenge", codeChallenge)
	q.Set("code_challenge_method", "S256")
	u.RawQuery = q.Encode()

	c.Redirect(http.StatusFound, u.String())
}

func (h *AuthHandler) OAuthCallback(c *gin.Context) {
	cfg := config.Load()
	isSecure := cfg.Env == "production"

	stateParam := c.Query("state")
	code := c.Query("code")

	var stateRandom, codeVerifier string
	if parts := strings.Split(stateParam, "."); len(parts) == 2 {
		stateRandom = parts[0]
		codeVerifier = parts[1]
	} else {
		stateRandom = stateParam
	}

	if codeVerifier == "" {
		codeVerifier, _ = c.Cookie("oauth_code_verifier")
	}

	cookieState, _ := c.Cookie("oauth_state")
	popupCookie, _ := c.Cookie("oauth_popup")
	isPopup := popupCookie == "true" || c.Query("popup") == "true"

	if stateParam == "" || (cookieState != "" && stateRandom != cookieState) {
		response.BadRequest(c, "Invalid state parameter", nil)
		return
	}

	if code == "" {
		response.BadRequest(c, "Missing authorization code", nil)
		return
	}

	authRes, err := h.service.HandleOAuthCallback(c.Request.Context(), code, codeVerifier)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	if isSecure {
		c.SetSameSite(http.SameSiteNoneMode)
	} else {
		c.SetSameSite(http.SameSiteLaxMode)
	}

	// Clear OAuth flow cookies
	c.SetCookie("oauth_state", "", -1, "/", "", isSecure, true)
	c.SetCookie("oauth_code_verifier", "", -1, "/", "", isSecure, true)
	c.SetCookie("oauth_popup", "", -1, "/", "", isSecure, false)

	// Set JWT as HttpOnly cookie
	c.SetCookie("votify_token", authRes.Token, 86400, "/", "", isSecure, true)

	if isPopup {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, `<!DOCTYPE html>
<html>
<head><title>Authentication Successful</title></head>
<body style="font-family: system-ui, sans-serif; display: flex; align-items: center; justify-content: center; height: 100vh; margin: 0; background: #faf9f6; color: #1a1a18;">
  <div style="text-align: center; padding: 24px; background: white; border-radius: 8px; box-shadow: 0 4px 12px rgba(0,0,0,0.08);">
    <h2 style="margin-bottom: 8px;">✓ Signed in successfully</h2>
    <p style="color: #6b6460; font-size: 14px;">Closing login window...</p>
  </div>
  <script>
    if (window.opener) {
      window.opener.postMessage({ type: 'OAUTH_SUCCESS' }, '*');
      window.close();
    } else {
      window.location.href = '`+strings.TrimRight(cfg.FrontendURL, "/")+`';
    }
  </script>
</body>
</html>`)
		return
	}

	c.Redirect(http.StatusFound, strings.TrimRight(cfg.FrontendURL, "/"))
}

func (h *AuthHandler) Logout(c *gin.Context) {
	cfg := config.Load()
	isSecure := cfg.Env == "production"

	// Clear the JWT cookie
	c.SetCookie("votify_token", "", -1, "/", "", isSecure, true)

	response.OK(c, gin.H{"message": "Logged out successfully"})
}
