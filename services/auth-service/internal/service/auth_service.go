package service

import (
	"context"
	"strings"
	"time"

	"github.com/votify/auth-service/internal/dto"
	"github.com/votify/auth-service/internal/model"
	"github.com/votify/auth-service/internal/repository"
	"github.com/votify/auth-service/internal/security"
	"github.com/votify/pkg/errors"
)

// AuthService defines the business logic interface for authentication.
type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error)
	GetUserByID(ctx context.Context, id string) (*dto.UserResponse, error)
}

// DefaultAuthService is the implementation of AuthService.
type DefaultAuthService struct {
	repo      repository.UserRepository
	jwtSecret string
}

func NewAuthService(repo repository.UserRepository, jwtSecret string) *DefaultAuthService {
	return &DefaultAuthService{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

func (s *DefaultAuthService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	email := strings.TrimSpace(strings.ToLower(req.Email))
	name := strings.TrimSpace(req.Name)

	existing, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}
	if existing != nil {
		return nil, errors.NewConflictError("An account with this email address already exists", "EMAIL_ALREADY_EXISTS")
	}

	hash, err := security.HashPassword(req.Password)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	user := &model.User{
		Name:         name,
		Email:        email,
		PasswordHash: hash,
		Role:         "user",
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, errors.NewInternalError(err)
	}

	token, err := security.GenerateJWT(user.ID, user.Email, user.Name, user.Role, s.jwtSecret, 24*time.Hour)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	return &dto.AuthResponse{
		Token: token,
		User: dto.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

func (s *DefaultAuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	email := strings.TrimSpace(strings.ToLower(req.Email))

	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}
	if user == nil {
		return nil, errors.NewUnauthorizedError("Invalid email or password")
	}

	if !security.VerifyPassword(req.Password, user.PasswordHash) {
		return nil, errors.NewUnauthorizedError("Invalid email or password")
	}

	token, err := security.GenerateJWT(user.ID, user.Email, user.Name, user.Role, s.jwtSecret, 24*time.Hour)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	return &dto.AuthResponse{
		Token: token,
		User: dto.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

func (s *DefaultAuthService) GetUserByID(ctx context.Context, id string) (*dto.UserResponse, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}
	if user == nil {
		return nil, errors.NewNotFoundError("User", id)
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}, nil
}
