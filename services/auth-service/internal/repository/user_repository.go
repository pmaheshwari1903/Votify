package repository

// UserRepository defines the interface for user data access.
// This abstracts MongoDB operations so the service layer
// never directly interacts with the database driver.
type UserRepository interface {
	// Create persists a new user document.
	// Create(ctx context.Context, user *model.User) error

	// FindByEmail looks up a user by email address.
	// FindByEmail(ctx context.Context, email string) (*model.User, error)

	// FindByID looks up a user by their unique ID.
	// FindByID(ctx context.Context, id string) (*model.User, error)

	// ExistsByEmail checks if a user with the given email already exists.
	// ExistsByEmail(ctx context.Context, email string) (bool, error)
}

// NOTE: The concrete MongoDB implementation of UserRepository
// will be built in the next development phase. It will use the
// official MongoDB Go driver (go.mongodb.org/mongo-driver).
//
// The interface is defined now to establish the data access boundary.
// No business logic should exist in repository implementations —
// only CRUD operations and query construction.
