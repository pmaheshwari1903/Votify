package service

import (
	"context"
	"testing"

	voteDto "github.com/votify/vote-service/internal/dto"
	"github.com/votify/vote-service/internal/model"
	"github.com/votify/vote-service/internal/repository"
)

func TestCastVote_RejectsMissingPollID(t *testing.T) {
	repo := repository.NewMemoryVoteRepository()

	service := NewVoteService(
		repo,
		nil,
		"http://localhost:8082",
	)

	_, err := service.CastVote(
		context.Background(),
		"user-1",
		"127.0.0.1",
		voteDto.CastVoteRequest{
			PollID:   "",
			OptionID: "option-1",
		},
	)

	if err == nil {
		t.Fatal("expected error for missing poll ID")
	}
}

func TestCastVote_RejectsMissingOptionID(t *testing.T) {
	repo := repository.NewMemoryVoteRepository()

	service := NewVoteService(
		repo,
		nil,
		"http://localhost:8082",
	)

	_, err := service.CastVote(
		context.Background(),
		"user-1",
		"127.0.0.1",
		voteDto.CastVoteRequest{
			PollID:   "poll-1",
			OptionID: "",
		},
	)

	if err == nil {
		t.Fatal("expected error for missing option ID")
	}
}

func TestVoteRepositoryDuplicateProtection(t *testing.T) {
	repo := repository.NewMemoryVoteRepository()

	vote := &model.Vote{
		PollID:        "poll-1",
		OptionID:      "option-1",
		VoterIdentity: "user-1",
		VoterIP:       "127.0.0.1",
		UserID:        "user-1",
	}

	if err := repo.Create(context.Background(), vote); err != nil {
		t.Fatalf("first vote should succeed: %v", err)
	}

	alreadyVoted, err := repo.HasVoted(
		context.Background(),
		"poll-1",
		"user-1",
	)

	if err != nil {
		t.Fatalf("failed checking duplicate vote: %v", err)
	}

	if !alreadyVoted {
		t.Fatal("expected voter to be detected as already voted")
	}
}