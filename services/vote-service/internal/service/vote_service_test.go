package service

import (
	"context"
	"testing"

	pollDto "github.com/votify/poll-service/internal/dto"
	pollRepo "github.com/votify/poll-service/internal/repository"
	pollRoutes "github.com/votify/poll-service/internal/routes"
	pollSvc "github.com/votify/poll-service/internal/service"

	voteDto "github.com/votify/vote-service/internal/dto"
	voteRepo "github.com/votify/vote-service/internal/repository"
)

func TestVoteService_VotingAndDuplicateProtection(t *testing.T) {
	// Setup shared poll repository
	sharedPollRepo := pollRoutes.GetSharedPollRepository()
	pSvc := pollSvc.NewPollService(sharedPollRepo)
	ctx := context.Background()

	// Create test poll
	poll, err := pSvc.CreatePoll(ctx, "usr_100", pollDto.CreatePollRequest{
		Question: "Which language do you prefer?",
		Options:  []string{"Go", "TypeScript", "Rust"},
	})
	if err != nil {
		t.Fatalf("Failed to create test poll: %v", err)
	}

	vRepo := voteRepo.NewMemoryVoteRepository()
	vSvc := NewVoteService(vRepo, "http://localhost:8082")

	// 1. Cast Valid Vote
	optID := poll.Options[0].ID
	res, err := vSvc.CastVote(ctx, "voter_1", "127.0.0.1", voteDto.CastVoteRequest{
		PollID:   poll.ID,
		OptionID: optID,
	})
	if err != nil {
		t.Fatalf("CastVote failed: %v", err)
	}
	if res.PollID != poll.ID || res.OptionID != optID {
		t.Errorf("Unexpected vote response: %v", res)
	}

	// 2. Duplicate Vote Rejection for same user
	_, err = vSvc.CastVote(ctx, "voter_1", "127.0.0.1", voteDto.CastVoteRequest{
		PollID:   poll.ID,
		OptionID: optID,
	})
	if err == nil {
		t.Errorf("Expected error on duplicate vote, got nil")
	}

	// 3. Invalid Option Vote Rejection
	_, err = vSvc.CastVote(ctx, "voter_2", "127.0.0.1", voteDto.CastVoteRequest{
		PollID:   poll.ID,
		OptionID: "invalid_option_id",
	})
	if err == nil {
		t.Errorf("Expected error on invalid option vote, got nil")
	}

	// 4. Vote on Closed Poll Rejection
	_, err = pSvc.ClosePoll(ctx, poll.ID, "usr_100")
	if err != nil {
		t.Fatalf("Failed to close poll: %v", err)
	}

	_, err = vSvc.CastVote(ctx, "voter_3", "127.0.0.1", voteDto.CastVoteRequest{
		PollID:   poll.ID,
		OptionID: optID,
	})
	if err == nil {
		t.Errorf("Expected error when voting on closed poll, got nil")
	}

	// 5. Verify Results Vote Count
	results, err := pSvc.GetPollResults(ctx, poll.ID)
	if err != nil {
		t.Fatalf("GetPollResults failed: %v", err)
	}
	if results.TotalVotes != 1 {
		t.Errorf("Expected 1 total vote, got %d", results.TotalVotes)
	}
	if results.Options[0].VoteCount != 1 {
		t.Errorf("Expected 1 vote for option 0, got %d", results.Options[0].VoteCount)
	}
}
