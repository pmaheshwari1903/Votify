package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// JWTClaims represents the payload embedded in Votify tokens.
type JWTClaims struct {
	Sub   string `json:"sub"`   // User ID
	Email string `json:"email"` // User Email
	Name  string `json:"name"`  // User Name
	Role  string `json:"role"`  // User Role ("user" or "admin")
	Exp   int64  `json:"exp"`   // Expiration Epoch
	Iat   int64  `json:"iat"`   // Issued At Epoch
}

// GenerateJWT creates a signed HMAC-SHA256 JWT token for an authenticated user.
func GenerateJWT(userID, email, name, role, secret string, duration time.Duration) (string, error) {
	if secret == "" {
		secret = "votify-default-secret-key-min-32-bytes"
	}

	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}

	now := time.Now().UTC()
	claims := JWTClaims{
		Sub:   userID,
		Email: email,
		Name:  name,
		Role:  role,
		Iat:   now.Unix(),
		Exp:   now.Add(duration).Unix(),
	}

	headerJSON, _ := json.Marshal(header)
	claimsJSON, _ := json.Marshal(claims)

	headerB64 := base64.RawURLEncoding.EncodeToString(headerJSON)
	claimsB64 := base64.RawURLEncoding.EncodeToString(claimsJSON)

	unsignedToken := fmt.Sprintf("%s.%s", headerB64, claimsB64)

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(unsignedToken))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return fmt.Sprintf("%s.%s", unsignedToken, signature), nil
}

// ValidateJWT verifies signature and expiration of a JWT token.
func ValidateJWT(tokenStr, secret string) (*JWTClaims, error) {
	if secret == "" {
		secret = "votify-default-secret-key-min-32-bytes"
	}

	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}

	unsignedToken := fmt.Sprintf("%s.%s", parts[0], parts[1])

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(unsignedToken))
	expectedSig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(parts[2]), []byte(expectedSig)) {
		return nil, errors.New("invalid token signature")
	}

	claimsJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errors.New("failed to decode claims")
	}

	var claims JWTClaims
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		return nil, errors.New("invalid claims format")
	}

	if time.Now().UTC().Unix() > claims.Exp {
		return nil, errors.New("token has expired")
	}

	return &claims, nil
}
