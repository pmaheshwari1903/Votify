package votes_test

import (
	"context"
	"testing"

	"github.com/votify/backend/internal/polls"
	"github.com/votify/backend/internal/realtime"
	"github.com/votify/backend/internal/repository"
	"github.com/votify/backend/internal/votes"
)

func TestCastVoteAndDuplicateProtection(t *testing.T) {
	pollRepo := repository.NewMemoryPollRepository()
	pollSvc := polls.NewPollService(pollRepo)

	voteRepo := repository.NewMemoryVoteRepository()
	rtSvc := realtime.NewRealtimeService("", "", 0, pollSvc)

	voteSvc := votes.NewVoteService(voteRepo, pollSvc, rtSvc)
	ctx := context.Background()

	// Create Poll
	poll, err := pollSvc.CreatePoll(ctx, "usr_owner", polls.CreatePollRequest{
		Question: "Which backend framework is best?",
		Options:  []string{"Gin", "Fiber", "Echo"},
	})
	if err != nil {
		t.Fatalf("CreatePoll failed: %v", err)
	}

	optID := poll.Options[0].ID

	// Cast Vote
	res, err := voteSvc.CastVote(ctx, "voter_1", "127.0.0.1", votes.CastVoteRequest{
		PollID:   poll.ID,
		OptionID: optID,
	})
	if err != nil {
		t.Fatalf("CastVote failed: %v", err)
	}
	if res.PollID != poll.ID || res.OptionID != optID {
		t.Fatalf("CastVote returned unexpected response: %+v", res)
	}

	// Verify vote count increment
	pResults, err := pollSvc.GetPollResults(ctx, poll.ID)
	if err != nil || pResults.TotalVotes != 1 {
		t.Fatalf("Expected total votes 1, got %d", pResults.TotalVotes)
	}

	// Duplicate Vote from same user -> should fail
	_, err = voteSvc.CastVote(ctx, "voter_1", "127.0.0.1", votes.CastVoteRequest{
		PollID:   poll.ID,
		OptionID: optID,
	})
	if err == nil {
		t.Fatalf("Expected duplicate vote to be rejected")
	}
}
