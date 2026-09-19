package repository

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/votify/backend/internal/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository interface {
	Create(ctx context.Context, user *auth.User) error
	Update(ctx context.Context, user *auth.User) error
	FindByEmail(ctx context.Context, email string) (*auth.User, error)
	FindByID(ctx context.Context, id string) (*auth.User, error)
	FindByOIDCSubject(ctx context.Context, issuer, subject string) (*auth.User, error)
}

// MongoUserRepository implements UserRepository using MongoDB.
type MongoUserRepository struct {
	coll *mongo.Collection
}

func NewMongoUserRepository(db *mongo.Database) *MongoUserRepository {
	coll := db.Collection("users")

	// Ensure unique index on email
	_, _ = coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	})

	// Index on oidc_issuer + oidc_subject
	_, _ = coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "oidc_issuer", Value: 1},
			{Key: "oidc_subject", Value: 1},
		},
	})

	return &MongoUserRepository{coll: coll}
}

func (r *MongoUserRepository) Create(ctx context.Context, user *auth.User) error {
	emailKey := strings.ToLower(user.Email)
	user.Email = emailKey
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = time.Now().UTC()

	if user.ID == "" {
		user.ID = fmt.Sprintf("usr_%d", time.Now().UnixNano())
	}

	_, err := r.coll.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) || strings.Contains(err.Error(), "duplicate key") {
			return fmt.Errorf("user with email '%s' already exists", user.Email)
		}
		return err
	}
	return nil
}

func (r *MongoUserRepository) Update(ctx context.Context, user *auth.User) error {
	user.UpdatedAt = time.Now().UTC()
	res, err := r.coll.ReplaceOne(ctx, bson.M{"_id": user.ID}, user)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("user '%s' not found", user.ID)
	}
	return nil
}

func (r *MongoUserRepository) FindByEmail(ctx context.Context, email string) (*auth.User, error) {
	emailKey := strings.ToLower(email)
	var user auth.User
	err := r.coll.FindOne(ctx, bson.M{"email": emailKey}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *MongoUserRepository) FindByID(ctx context.Context, id string) (*auth.User, error) {
	var user auth.User
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *MongoUserRepository) FindByOIDCSubject(ctx context.Context, issuer, subject string) (*auth.User, error) {
	if issuer == "" || subject == "" {
		return nil, nil
	}
	var user auth.User
	err := r.coll.FindOne(ctx, bson.M{"oidc_issuer": issuer, "oidc_subject": subject}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// MemoryUserRepository is an in-memory thread-safe fallback implementation.
type MemoryUserRepository struct {
	mu         sync.RWMutex
	users      map[string]*auth.User
	byEmail    map[string]*auth.User
	byOIDCSub  map[string]*auth.User // key: "issuer:subject"
	counter    int64
}

func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		users:     make(map[string]*auth.User),
		byEmail:   make(map[string]*auth.User),
		byOIDCSub: make(map[string]*auth.User),
	}
}

func (r *MemoryUserRepository) Create(ctx context.Context, user *auth.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	emailKey := strings.ToLower(user.Email)
	if _, exists := r.byEmail[emailKey]; exists {
		return fmt.Errorf("user with email '%s' already exists", user.Email)
	}

	r.counter++
	if user.ID == "" {
		user.ID = fmt.Sprintf("usr_%d_%d", time.Now().UnixNano(), r.counter)
	}
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = time.Now().UTC()

	cp := *user
	r.users[user.ID] = &cp
	r.byEmail[emailKey] = &cp
	if user.OIDCIssuer != "" && user.OIDCSubject != "" {
		r.byOIDCSub[user.OIDCIssuer+":"+user.OIDCSubject] = &cp
	}

	return nil
}

func (r *MemoryUserRepository) Update(ctx context.Context, user *auth.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; !exists {
		return fmt.Errorf("user '%s' not found", user.ID)
	}

	user.UpdatedAt = time.Now().UTC()
	cp := *user
	r.users[user.ID] = &cp
	r.byEmail[strings.ToLower(user.Email)] = &cp
	if user.OIDCIssuer != "" && user.OIDCSubject != "" {
		r.byOIDCSub[user.OIDCIssuer+":"+user.OIDCSubject] = &cp
	}
	return nil
}

func (r *MemoryUserRepository) FindByEmail(ctx context.Context, email string) (*auth.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	emailKey := strings.ToLower(email)
	user, exists := r.byEmail[emailKey]
	if !exists {
		return nil, nil
	}
	cp := *user
	return &cp, nil
}

func (r *MemoryUserRepository) FindByID(ctx context.Context, id string) (*auth.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, nil
	}
	cp := *user
	return &cp, nil
}

func (r *MemoryUserRepository) FindByOIDCSubject(ctx context.Context, issuer, subject string) (*auth.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	key := issuer + ":" + subject
	user, exists := r.byOIDCSub[key]
	if !exists {
		return nil, nil
	}
	cp := *user
	return &cp, nil
}
