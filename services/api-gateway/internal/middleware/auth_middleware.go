package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type JWTClaims struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
	Exp   int64  `json:"exp"`
}

// AuthGuard validates Authorization: Bearer <token> and sets X-User-* headers on requests.
func AuthGuard(secret string) gin.HandlerFunc {
	if secret == "" {
		secret = "votify-default-secret-key-min-32-bytes"
	}

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "Authentication token required",
				},
			})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := parseAndValidateToken(tokenStr, secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "Invalid or expired authentication token: " + err.Error(),
				},
			})
			c.Abort()
			return
		}

		// Inject user identity headers for downstream microservices
		c.Request.Header.Set("X-User-ID", claims.Sub)
		c.Request.Header.Set("X-User-Email", claims.Email)
		c.Request.Header.Set("X-User-Role", claims.Role)

		c.Next()
	}
}

// OptionalAuth extracts identity headers if present, but does not abort if missing.
func OptionalAuth(secret string) gin.HandlerFunc {
	if secret == "" {
		secret = "votify-default-secret-key-min-32-bytes"
	}

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := parseAndValidateToken(tokenStr, secret)
			if err == nil && claims != nil {
				c.Request.Header.Set("X-User-ID", claims.Sub)
				c.Request.Header.Set("X-User-Email", claims.Email)
				c.Request.Header.Set("X-User-Role", claims.Role)
			}
		}
		c.Next()
	}
}

func parseAndValidateToken(tokenStr, secret string) (*JWTClaims, error) {
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return nil, http.ErrAbortHandler
	}

	unsignedToken := parts[0] + "." + parts[1]
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(unsignedToken))
	expectedSig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(parts[2]), []byte(expectedSig)) {
		return nil, http.ErrAbortHandler
	}

	claimsBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	var claims JWTClaims
	if err := json.Unmarshal(claimsBytes, &claims); err != nil {
		return nil, err
	}

	if time.Now().UTC().Unix() > claims.Exp {
		return nil, http.ErrAbortHandler
	}

	return &claims, nil
}
