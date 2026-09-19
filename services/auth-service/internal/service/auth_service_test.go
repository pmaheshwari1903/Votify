package service

import (
	"context"
	"testing"

	"github.com/votify/auth-service/internal/dto"
	"github.com/votify/auth-service/internal/repository"
)

func TestAuthService_RegisterAndLogin(t *testing.T) {
	repo := repository.NewMemoryUserRepository()
	svc := NewAuthService(repo, "test-jwt-secret-key-min-32-bytes")
	ctx := context.Background()

	// 1. Successful Registration
	regReq := dto.RegisterRequest{
		Email:    "testvoter@example.com",
		Name:     "Test Voter",
		Password: "password123",
	}
	regRes, err := svc.Register(ctx, regReq)
	if err != nil {
		t.Fatalf("Registration failed: %v", err)
	}
	if regRes.Token == "" {
		t.Errorf("Expected JWT token in AuthResponse, got empty string")
	}
	if regRes.User.Email != "testvoter@example.com" {
		t.Errorf("Expected email testvoter@example.com, got %s", regRes.User.Email)
	}

	// 2. Duplicate Registration Rejection
	_, err = svc.Register(ctx, regReq)
	if err == nil {
		t.Errorf("Expected error on duplicate registration, got nil")
	}

	// 3. Successful Login
	loginReq := dto.LoginRequest{
		Email:    "testvoter@example.com",
		Password: "password123",
	}
	loginRes, err := svc.Login(ctx, loginReq)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	if loginRes.Token == "" {
		t.Errorf("Expected JWT token on login, got empty string")
	}

	// 4. Invalid Password Rejection
	badLoginReq := dto.LoginRequest{
		Email:    "testvoter@example.com",
		Password: "wrongpassword",
	}
	_, err = svc.Login(ctx, badLoginReq)
	if err == nil {
		t.Errorf("Expected error on invalid password, got nil")
	}
}
