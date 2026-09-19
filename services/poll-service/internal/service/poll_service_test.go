package service

import (
	"context"
	"testing"

	"github.com/votify/poll-service/internal/dto"
	"github.com/votify/poll-service/internal/repository"
)

func TestPollService_FullLifecycle(t *testing.T) {
	repo := repository.NewMemoryPollRepository()
	svc := NewPollService(repo)
	ctx := context.Background()

	ownerID := "usr_owner_1"
	otherUser := "usr_other_2"

	// 1. Create Poll
	createReq := dto.CreatePollRequest{
		Question:    "What design pattern should we use?",
		Description: "Voting for architecture style",
		Options:     []string{"Monolith", "Microservices", "Serverless"},
	}

	poll, err := svc.CreatePoll(ctx, ownerID, createReq)
	if err != nil {
		t.Fatalf("CreatePoll failed: %v", err)
	}
	if poll.ID == "" || len(poll.Options) != 3 {
		t.Fatalf("Expected poll with 3 options, got %v", poll)
	}

	// 2. Retrieve Public Poll
	pubPoll, err := svc.GetPublicPoll(ctx, poll.ID)
	if err != nil {
		t.Fatalf("GetPublicPoll failed: %v", err)
	}
	if pubPoll.Question != createReq.Question {
		t.Errorf("Expected question '%s', got '%s'", createReq.Question, pubPoll.Question)
	}

	// 3. Unauthorized Owner Management Check
	_, err = svc.ClosePoll(ctx, poll.ID, otherUser)
	if err == nil {
		t.Errorf("Expected 403 Forbidden error for non-owner, got nil")
	}

	// 4. Authorized Close Poll
	closedPoll, err := svc.ClosePoll(ctx, poll.ID, ownerID)
	if err != nil {
		t.Fatalf("ClosePoll failed: %v", err)
	}
	if closedPoll.Status != "closed" {
		t.Errorf("Expected status 'closed', got '%s'", closedPoll.Status)
	}

	// 5. Open Poll again
	openPoll, err := svc.OpenPoll(ctx, poll.ID, ownerID)
	if err != nil {
		t.Fatalf("OpenPoll failed: %v", err)
	}
	if openPoll.Status != "open" {
		t.Errorf("Expected status 'open', got '%s'", openPoll.Status)
	}
}
