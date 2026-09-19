package repository

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/votify/backend/internal/polls"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PollRepository interface {
	Create(ctx context.Context, poll *polls.Poll) error
	FindByID(ctx context.Context, id string) (*polls.Poll, error)
	FindByOwnerID(ctx context.Context, ownerID string) ([]*polls.Poll, error)
	Update(ctx context.Context, poll *polls.Poll) error
	Delete(ctx context.Context, id string) error
	IncrementVote(ctx context.Context, pollID string, optionID string) error
}

// MongoPollRepository implements PollRepository using MongoDB.
type MongoPollRepository struct {
	coll *mongo.Collection
}

func NewMongoPollRepository(db *mongo.Database) *MongoPollRepository {
	coll := db.Collection("polls")

	// Create index on owner_id for query speed
	_, _ = coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{{Key: "owner_id", Value: 1}},
	})

	return &MongoPollRepository{coll: coll}
}

func (r *MongoPollRepository) Create(ctx context.Context, poll *polls.Poll) error {
	poll.CreatedAt = time.Now().UTC()
	poll.UpdatedAt = time.Now().UTC()
	if poll.ID == "" {
		poll.ID = fmt.Sprintf("poll_%d", time.Now().UnixNano())
	}

	_, err := r.coll.InsertOne(ctx, poll)
	return err
}

func (r *MongoPollRepository) FindByID(ctx context.Context, id string) (*polls.Poll, error) {
	var poll polls.Poll
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&poll)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &poll, nil
}

func (r *MongoPollRepository) FindByOwnerID(ctx context.Context, ownerID string) ([]*polls.Poll, error) {
	cursor, err := r.coll.Find(ctx, bson.M{"owner_id": ownerID}, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var list []*polls.Poll
	for cursor.Next(ctx) {
		var p polls.Poll
		if err := cursor.Decode(&p); err == nil {
			cp := p
			list = append(list, &cp)
		}
	}
	return list, nil
}

func (r *MongoPollRepository) Update(ctx context.Context, poll *polls.Poll) error {
	poll.UpdatedAt = time.Now().UTC()
	res, err := r.coll.ReplaceOne(ctx, bson.M{"_id": poll.ID}, poll)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("poll '%s' not found", poll.ID)
	}
	return nil
}

func (r *MongoPollRepository) Delete(ctx context.Context, id string) error {
	res, err := r.coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("poll '%s' not found", id)
	}
	return nil
}

func (r *MongoPollRepository) IncrementVote(ctx context.Context, pollID string, optionID string) error {
	filter := bson.M{
		"_id":        pollID,
		"options.id": optionID,
	}
	update := bson.M{
		"$inc": bson.M{
			"total_votes":       1,
			"options.$.vote_count": 1,
		},
		"$set": bson.M{
			"updated_at": time.Now().UTC(),
		},
	}

	res, err := r.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("poll '%s' or option '%s' not found", pollID, optionID)
	}
	return nil
}

// MemoryPollRepository is an in-memory thread-safe fallback.
type MemoryPollRepository struct {
	mu      sync.RWMutex
	polls   map[string]*polls.Poll
	counter int64
}

func NewMemoryPollRepository() *MemoryPollRepository {
	return &MemoryPollRepository{
		polls: make(map[string]*polls.Poll),
	}
}

func (r *MemoryPollRepository) Create(ctx context.Context, poll *polls.Poll) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.counter++
	if poll.ID == "" {
		poll.ID = fmt.Sprintf("poll_%d_%d", time.Now().UnixNano(), r.counter)
	}
	poll.CreatedAt = time.Now().UTC()
	poll.UpdatedAt = time.Now().UTC()

	cp := *poll
	r.polls[poll.ID] = &cp
	return nil
}

func (r *MemoryPollRepository) FindByID(ctx context.Context, id string) (*polls.Poll, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, exists := r.polls[id]
	if !exists {
		return nil, nil
	}
	cp := *p
	return &cp, nil
}

func (r *MemoryPollRepository) FindByOwnerID(ctx context.Context, ownerID string) ([]*polls.Poll, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*polls.Poll
	for _, p := range r.polls {
		if p.OwnerID == ownerID {
			cp := *p
			result = append(result, &cp)
		}
	}
	return result, nil
}

func (r *MemoryPollRepository) Update(ctx context.Context, poll *polls.Poll) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.polls[poll.ID]; !exists {
		return fmt.Errorf("poll '%s' not found", poll.ID)
	}

	poll.UpdatedAt = time.Now().UTC()
	cp := *poll
	r.polls[poll.ID] = &cp
	return nil
}

func (r *MemoryPollRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.polls[id]; !exists {
		return fmt.Errorf("poll '%s' not found", id)
	}

	delete(r.polls, id)
	return nil
}

func (r *MemoryPollRepository) IncrementVote(ctx context.Context, pollID string, optionID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	p, exists := r.polls[pollID]
	if !exists {
		return fmt.Errorf("poll '%s' not found", pollID)
	}

	optionFound := false
	for i := range p.Options {
		if p.Options[i].ID == optionID {
			p.Options[i].VoteCount++
			p.TotalVotes++
			optionFound = true
			break
		}
	}

	if !optionFound {
		return fmt.Errorf("option '%s' not found in poll '%s'", optionID, pollID)
	}

	p.UpdatedAt = time.Now().UTC()
	return nil
}
