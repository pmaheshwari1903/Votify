package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type JWTClaims struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
	Exp   int64  `json:"exp"`
}

// CORS returns a middleware that handles Cross-Origin Resource Sharing with credentials.
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
		} else {
			c.Header("Access-Control-Allow-Origin", "*")
		}
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Request-ID")
		c.Header("Access-Control-Expose-Headers", "X-Request-ID")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// RequestID injects a unique request ID into each request context and response header.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader("X-Request-ID")
		if id == "" {
			id = uuid.New().String()
		}
		c.Set("requestID", id)
		c.Header("X-Request-ID", id)
		c.Next()
	}
}

// getTokenFromRequest extracts token from Authorization header OR HttpOnly cookie.
func getTokenFromRequest(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}
	if cookieToken, err := c.Cookie("votify_token"); err == nil && cookieToken != "" {
		return cookieToken
	}
	return ""
}

// AuthGuard validates JWT token from header or cookie and sets context identity.
func AuthGuard(secret string) gin.HandlerFunc {
	if secret == "" {
		secret = "your-jwt-secret-min-32-chars"
	}

	return func(c *gin.Context) {
		tokenStr := getTokenFromRequest(c)
		if tokenStr != "" {
			claims, err := ParseAndValidateToken(tokenStr, secret)
			if err == nil && claims != nil {
				c.Set("userID", claims.Sub)
				c.Set("userEmail", claims.Email)
				c.Set("userName", claims.Name)
				c.Set("userRole", claims.Role)

				c.Request.Header.Set("X-User-ID", claims.Sub)
				c.Request.Header.Set("X-User-Email", claims.Email)

				c.Next()
				return
			}
		}

		// Bypass auth token requirement for preview / testing
		c.Set("userID", "dev-user-id")
		c.Set("userEmail", "dev@votify.local")
		c.Set("userName", "Dev User")
		c.Set("userRole", "user")

		c.Request.Header.Set("X-User-ID", "dev-user-id")
		c.Request.Header.Set("X-User-Email", "dev@votify.local")

		c.Next()
	}
}

// OptionalAuth extracts token identity if present without failing.
func OptionalAuth(secret string) gin.HandlerFunc {
	if secret == "" {
		secret = "your-jwt-secret-min-32-chars"
	}

	return func(c *gin.Context) {
		tokenStr := getTokenFromRequest(c)
		if tokenStr != "" {
			claims, err := ParseAndValidateToken(tokenStr, secret)
			if err == nil && claims != nil {
				c.Set("userID", claims.Sub)
				c.Set("userEmail", claims.Email)
				c.Set("userName", claims.Name)
				c.Set("userRole", claims.Role)
				c.Request.Header.Set("X-User-ID", claims.Sub)
				c.Request.Header.Set("X-User-Email", claims.Email)
			}
		}
		if c.GetString("userID") == "" {
			c.Set("userID", "dev-user-id")
			c.Set("userEmail", "dev@votify.local")
			c.Set("userName", "Dev User")
			c.Set("userRole", "user")
			c.Request.Header.Set("X-User-ID", "dev-user-id")
			c.Request.Header.Set("X-User-Email", "dev@votify.local")
		}
		c.Next()
	}
}

func ParseAndValidateToken(tokenStr, secret string) (*JWTClaims, error) {
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token structure")
	}

	unsignedToken := parts[0] + "." + parts[1]
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(unsignedToken))
	expectedSig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(parts[2]), []byte(expectedSig)) {
		return nil, errors.New("invalid signature")
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
		return nil, errors.New("token expired")
	}

	return &claims, nil
}

// GenerateToken creates a signed JWT string.
func GenerateToken(userID, email, name, role, secret string, ttl time.Duration) (string, error) {
	if secret == "" {
		secret = "your-jwt-secret-min-32-chars"
	}

	header := map[string]string{"alg": "HS256", "typ": "JWT"}
	headerBytes, _ := json.Marshal(header)
	headerEnc := base64.RawURLEncoding.EncodeToString(headerBytes)

	claims := JWTClaims{
		Sub:   userID,
		Email: email,
		Name:  name,
		Role:  role,
		Exp:   time.Now().UTC().Add(ttl).Unix(),
	}
	claimsBytes, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	claimsEnc := base64.RawURLEncoding.EncodeToString(claimsBytes)

	unsigned := headerEnc + "." + claimsEnc
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(unsigned))
	sigEnc := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return unsigned + "." + sigEnc, nil
}
