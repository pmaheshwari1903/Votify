package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/votify/backend/internal/config"
	"github.com/votify/backend/internal/errors"
	"github.com/votify/backend/internal/middleware"
)

type UserRepo interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id string) (*User, error)
}

type AuthService interface {
	Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error)
	Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)
	GetUserByID(ctx context.Context, id string) (*UserResponse, error)
	HandleOAuthCallback(ctx context.Context, code, codeVerifier string) (*AuthResponse, error)
}

type DefaultAuthService struct {
	repo      UserRepo
	jwtSecret string
}

func NewAuthService(repo UserRepo, jwtSecret string) *DefaultAuthService {
	return &DefaultAuthService{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

func (s *DefaultAuthService) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	email := strings.TrimSpace(strings.ToLower(req.Email))
	name := strings.TrimSpace(req.Name)

	existing, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}
	if existing != nil {
		return nil, errors.NewConflictError("An account with this email address already exists", "EMAIL_ALREADY_EXISTS")
	}

	hash, err := HashPassword(req.Password)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	user := &User{
		Name:         name,
		Email:        email,
		PasswordHash: hash,
		Role:         "user",
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, errors.NewInternalError(err)
	}

	token, err := middleware.GenerateToken(user.ID, user.Email, user.Name, user.Role, s.jwtSecret, 24*time.Hour)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	return &AuthResponse{
		Token: token,
		User: UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

func (s *DefaultAuthService) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	email := strings.TrimSpace(strings.ToLower(req.Email))

	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}
	if user == nil {
		return nil, errors.NewUnauthorizedError("Invalid email or password")
	}

	if !VerifyPassword(req.Password, user.PasswordHash) {
		return nil, errors.NewUnauthorizedError("Invalid email or password")
	}

	token, err := middleware.GenerateToken(user.ID, user.Email, user.Name, user.Role, s.jwtSecret, 24*time.Hour)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	return &AuthResponse{
		Token: token,
		User: UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

func (s *DefaultAuthService) GetUserByID(ctx context.Context, id string) (*UserResponse, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}
	if user == nil {
		return nil, errors.NewNotFoundError("User", id)
	}

	return &UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *DefaultAuthService) HandleOAuthCallback(ctx context.Context, code, codeVerifier string) (*AuthResponse, error) {
	cfg := config.Load()

	// 1. OIDC Discovery
	meta, err := FetchOIDCDiscovery(ctx, cfg.OAuthIssuer)
	if err != nil {
		// Fallback to standard provider paths if discovery endpoint fails
		meta = &OIDCDiscoveryMetadata{
			Issuer:                cfg.OAuthIssuer,
			AuthorizationEndpoint: strings.TrimRight(cfg.OAuthIssuer, "/") + "/authorize",
			TokenEndpoint:         strings.TrimRight(cfg.OAuthIssuer, "/") + "/token",
			UserinfoEndpoint:      strings.TrimRight(cfg.OAuthIssuer, "/") + "/userinfo",
		}
	}

	// 2. Exchange Code for Access Token
	tokens, err := ExchangeCodeForTokens(
		ctx,
		meta.TokenEndpoint,
		cfg.OAuthClientID,
		cfg.OAuthClientSecret,
		code,
		cfg.OAuthRedirectURI,
		codeVerifier,
	)
	if err != nil {
		return nil, errors.NewUnauthorizedError("Failed to exchange OAuth code: " + err.Error())
	}

	// 3. Fetch UserInfo from OIDC provider
	userInfo, err := FetchUserInfo(ctx, meta.UserinfoEndpoint, tokens.AccessToken)
	if err != nil {
		return nil, errors.NewUnauthorizedError("Failed to fetch OIDC userinfo: " + err.Error())
	}

	email := strings.TrimSpace(strings.ToLower(userInfo.Email))
	if email == "" {
		if userInfo.Sub != "" {
			email = userInfo.Sub + "@oidcauth.user"
		} else {
			return nil, errors.NewUnauthorizedError("No verified email or subject claim found in OIDC identity")
		}
	}

	name := strings.TrimSpace(userInfo.Name)
	if name == "" {
		name = strings.TrimSpace(userInfo.GivenName + " " + userInfo.FamilyName)
	}
	if name == "" {
		parts := strings.Split(email, "@")
		name = parts[0]
	}

	// 4. Find or Create User in MongoDB
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	if user == nil {
		user = &User{
			Name:         name,
			Email:        email,
			PasswordHash: "oauth_authenticated",
			Role:         "user",
		}
		if err := s.repo.Create(ctx, user); err != nil {
			return nil, errors.NewInternalError(err)
		}
	}

	// 5. Generate Votify JWT
	jwtToken, err := middleware.GenerateToken(user.ID, user.Email, user.Name, user.Role, s.jwtSecret, 24*time.Hour)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	return &AuthResponse{
		Token: jwtToken,
		User: UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

func HashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	saltHex := hex.EncodeToString(salt)
	hash := sha256.Sum256([]byte(password + saltHex))
	hashHex := hex.EncodeToString(hash[:])

	return fmt.Sprintf("%s:%s", saltHex, hashHex), nil
}

func VerifyPassword(password, storedHash string) bool {
	parts := strings.Split(storedHash, ":")
	if len(parts) != 2 {
		return false
	}

	saltHex := parts[0]
	expectedHash := parts[1]

	hash := sha256.Sum256([]byte(password + saltHex))
	computedHash := hex.EncodeToString(hash[:])

	return computedHash == expectedHash
}
