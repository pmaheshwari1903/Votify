package auth

import (
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

	state := GenerateSecureRandomState()
	codeVerifier := GeneratePKCEVerifier()
	codeChallenge := ComputePKCEChallengeS256(codeVerifier)

	c.SetCookie("oauth_state", state, 600, "/", "", isSecure, true)
	c.SetCookie("oauth_code_verifier", codeVerifier, 600, "/", "", isSecure, true)

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
	q.Set("state", state)
	q.Set("code_challenge", codeChallenge)
	q.Set("code_challenge_method", "S256")
	u.RawQuery = q.Encode()

	c.Redirect(http.StatusFound, u.String())
}

func (h *AuthHandler) OAuthCallback(c *gin.Context) {
	cfg := config.Load()
	isSecure := cfg.Env == "production"

	state := c.Query("state")
	code := c.Query("code")

	cookieState, _ := c.Cookie("oauth_state")
	codeVerifier, _ := c.Cookie("oauth_code_verifier")

	if state == "" || (cookieState != "" && state != cookieState) {
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

	// Clear OAuth flow cookies
	c.SetCookie("oauth_state", "", -1, "/", "", isSecure, true)
	c.SetCookie("oauth_code_verifier", "", -1, "/", "", isSecure, true)

	// Set JWT as HttpOnly cookie instead of passing in URL
	c.SetCookie("votify_token", authRes.Token, 86400, "/", "", isSecure, true)

	c.Redirect(http.StatusFound, strings.TrimRight(cfg.FrontendURL, "/"))
}

func (h *AuthHandler) Logout(c *gin.Context) {
	cfg := config.Load()
	isSecure := cfg.Env == "production"

	// Clear the JWT cookie
	c.SetCookie("votify_token", "", -1, "/", "", isSecure, true)

	response.OK(c, gin.H{"message": "Logged out successfully"})
}
