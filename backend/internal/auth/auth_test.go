package auth_test

import (
	"context"
	"testing"

	"github.com/votify/backend/internal/auth"
	"github.com/votify/backend/internal/repository"
)

func TestAuthRegisterAndLogin(t *testing.T) {
	repo := repository.NewMemoryUserRepository()
	svc := auth.NewAuthService(repo, "test-secret-key-min-32-chars-long")

	ctx := context.Background()

	// Register
	regReq := auth.RegisterRequest{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "password123",
	}

	regRes, err := svc.Register(ctx, regReq)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}
	if regRes.Token == "" {
		t.Fatalf("Expected non-empty token")
	}
	if regRes.User.Email != "test@example.com" {
		t.Fatalf("Expected email test@example.com, got %s", regRes.User.Email)
	}

	// Duplicate Register
	_, err = svc.Register(ctx, regReq)
	if err == nil {
		t.Fatalf("Expected duplicate email registration to fail")
	}

	// Login
	loginRes, err := svc.Login(ctx, auth.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	if loginRes.Token == "" {
		t.Fatalf("Expected non-empty login token")
	}

	// Get User By ID
	userRes, err := svc.GetUserByID(ctx, regRes.User.ID)
	if err != nil {
		t.Fatalf("GetUserByID failed: %v", err)
	}
	if userRes.Name != "Test User" {
		t.Fatalf("Expected name 'Test User', got '%s'", userRes.Name)
	}
}
