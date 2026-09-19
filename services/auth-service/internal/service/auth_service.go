package service

// AuthService defines the interface for authentication business logic.
// This interface allows the handler to depend on an abstraction,
// making it possible to swap implementations (e.g., for testing,
// or when integrating with Myshri.com authentication later).
type AuthService interface {
	// Register creates a new user account.
	// Business rules enforced here:
	// - Email uniqueness
	// - Password hashing
	// - Default role assignment
	// Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error)

	// Login authenticates a user and returns a JWT token.
	// Business rules enforced here:
	// - Credential verification
	// - Token generation
	// Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error)

	// GetUserByID retrieves a user by their ID.
	// GetUserByID(ctx context.Context, id string) (*dto.UserResponse, error)
}

// NOTE: The concrete implementation of AuthService will be built
// in the next development phase. It will depend on:
// - UserRepository (for MongoDB access)
// - JWT utility (for token generation/validation)
// - Password hasher (bcrypt)
//
// The interface is defined now to establish the architectural boundary.
