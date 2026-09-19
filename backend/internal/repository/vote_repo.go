package repository

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/votify/backend/internal/votes"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type VoteRepository interface {
	Create(ctx context.Context, vote *votes.Vote) error
	HasVoted(ctx context.Context, pollID, voterIdentity string) (bool, error)
}

// MongoVoteRepository implements VoteRepository using MongoDB.
type MongoVoteRepository struct {
	coll *mongo.Collection
}

func NewMongoVoteRepository(db *mongo.Database) *MongoVoteRepository {
	coll := db.Collection("votes")

	// Ensure compound unique index on (poll_id, voter_identity) for concurrency-safe duplicate protection
	_, _ = coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "poll_id", Value: 1},
			{Key: "voter_identity", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})

	return &MongoVoteRepository{coll: coll}
}

func (r *MongoVoteRepository) HasVoted(ctx context.Context, pollID, voterIdentity string) (bool, error) {
	filter := bson.M{
		"poll_id":        pollID,
		"voter_identity": voterIdentity,
	}
	count, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *MongoVoteRepository) Create(ctx context.Context, vote *votes.Vote) error {
	vote.CreatedAt = time.Now().UTC()
	if vote.ID == "" {
		vote.ID = fmt.Sprintf("vt_%d", time.Now().UnixNano())
	}

	_, err := r.coll.InsertOne(ctx, vote)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) || strings.Contains(err.Error(), "duplicate key") {
			return fmt.Errorf("duplicate vote detected for poll '%s' and voter '%s'", vote.PollID, vote.VoterIdentity)
		}
		return err
	}
	return nil
}

// MemoryVoteRepository is an in-memory thread-safe fallback using compound unique key `pollID:voterIdentity`.
type MemoryVoteRepository struct {
	mu          sync.RWMutex
	votes       map[string]*votes.Vote
	votedKeyMap map[string]bool // key: "pollID:voterIdentity"
	counter     int64
}

func NewMemoryVoteRepository() *MemoryVoteRepository {
	return &MemoryVoteRepository{
		votes:       make(map[string]*votes.Vote),
		votedKeyMap: make(map[string]bool),
	}
}

func (r *MemoryVoteRepository) HasVoted(ctx context.Context, pollID, voterIdentity string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	key := fmt.Sprintf("%s:%s", pollID, voterIdentity)
	return r.votedKeyMap[key], nil
}

func (r *MemoryVoteRepository) Create(ctx context.Context, vote *votes.Vote) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := fmt.Sprintf("%s:%s", vote.PollID, vote.VoterIdentity)
	if r.votedKeyMap[key] {
		return fmt.Errorf("duplicate vote detected for poll '%s' and voter '%s'", vote.PollID, vote.VoterIdentity)
	}

	r.counter++
	if vote.ID == "" {
		vote.ID = fmt.Sprintf("vt_%d_%d", time.Now().UnixNano(), r.counter)
	}
	vote.CreatedAt = time.Now().UTC()

	r.votedKeyMap[key] = true
	cp := *vote
	r.votes[vote.ID] = &cp

	return nil
}
