package polls_test

import (
	"context"
	"testing"

	"github.com/votify/backend/internal/polls"
	"github.com/votify/backend/internal/repository"
)

func TestPollLifecycle(t *testing.T) {
	repo := repository.NewMemoryPollRepository()
	svc := polls.NewPollService(repo)

	ctx := context.Background()
	ownerID := "usr_123"

	// Create Poll
	createReq := polls.CreatePollRequest{
		Question:    "Favorite Programming Language?",
		Description: "Choose your primary language",
		Options:     []string{"Go", "TypeScript", "Python"},
	}

	poll, err := svc.CreatePoll(ctx, ownerID, createReq)
	if err != nil {
		t.Fatalf("CreatePoll failed: %v", err)
	}
	if len(poll.Options) != 3 {
		t.Fatalf("Expected 3 options, got %d", len(poll.Options))
	}
	if poll.Status != "open" {
		t.Fatalf("Expected poll status open, got %s", poll.Status)
	}

	// List Polls
	list, err := svc.ListOwnerPolls(ctx, ownerID)
	if err != nil || len(list) != 1 {
		t.Fatalf("ListOwnerPolls failed, count=%d, err=%v", len(list), err)
	}

	// Close & Open Poll
	poll, err = svc.ClosePoll(ctx, poll.ID, ownerID)
	if err != nil || poll.Status != "closed" {
		t.Fatalf("ClosePoll failed: %v", err)
	}

	poll, err = svc.OpenPoll(ctx, poll.ID, ownerID)
	if err != nil || poll.Status != "open" {
		t.Fatalf("OpenPoll failed: %v", err)
	}
}
