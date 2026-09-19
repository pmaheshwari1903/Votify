package repository

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/votify/auth-service/internal/model"
)

// UserRepository defines data access methods for users.
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByID(ctx context.Context, id string) (*model.User, error)
}

// MemoryUserRepository is an in-memory thread-safe implementation of UserRepository.
// Provides instant persistence when running without a live MongoDB cluster.
type MemoryUserRepository struct {
	mu     sync.RWMutex
	users  map[string]*model.User
	byEmail map[string]*model.User
	counter int64
}

func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		users:   make(map[string]*model.User),
		byEmail: make(map[string]*model.User),
	}
}

func (r *MemoryUserRepository) Create(ctx context.Context, user *model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	emailKey := strings.ToLower(user.Email)
	if _, exists := r.byEmail[emailKey]; exists {
		return fmt.Errorf("user with email '%s' already exists", user.Email)
	}

	r.counter++
	user.ID = fmt.Sprintf("usr_%d_%d", time.Now().UnixNano(), r.counter)
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = time.Now().UTC()

	// Store copy
	copy := *user
	r.users[user.ID] = &copy
	r.byEmail[emailKey] = &copy

	return nil
}

func (r *MemoryUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	emailKey := strings.ToLower(email)
	user, exists := r.byEmail[emailKey]
	if !exists {
		return nil, nil
	}
	copy := *user
	return &copy, nil
}

func (r *MemoryUserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, nil
	}
	copy := *user
	return &copy, nil
}
