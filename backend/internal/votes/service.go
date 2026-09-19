package votes

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/votify/backend/internal/errors"
	"github.com/votify/backend/internal/polls"
	"github.com/votify/backend/internal/realtime"
)

type VoteRepo interface {
	Create(ctx context.Context, vote *Vote) error
	HasVoted(ctx context.Context, pollID, voterIdentity string) (bool, error)
}

type VoteService interface {
	CastVote(ctx context.Context, userID, anonVoterID, clientIP string, req CastVoteRequest) (*VoteResponse, error)
}

type DefaultVoteService struct {
	voteRepo        VoteRepo
	pollService     polls.PollService
	realtimeService *realtime.RealtimeService
}

func NewVoteService(
	voteRepo VoteRepo,
	pollService polls.PollService,
	realtimeService *realtime.RealtimeService,
) *DefaultVoteService {
	return &DefaultVoteService{
		voteRepo:        voteRepo,
		pollService:     pollService,
		realtimeService: realtimeService,
	}
}

func (s *DefaultVoteService) CastVote(
	ctx context.Context,
	userID, anonVoterID, clientIP string,
	req CastVoteRequest,
) (*VoteResponse, error) {
	pollID := strings.TrimSpace(req.PollID)
	optionID := strings.TrimSpace(req.OptionID)

	if pollID == "" || optionID == "" {
		return nil, errors.NewBadRequestError("pollId and optionId are required", "MISSING_REQUIRED_FIELDS")
	}

	pubPoll, err := s.pollService.GetPublicPoll(ctx, pollID)
	if err != nil {
		return nil, err
	}

	if pubPoll.Status != "open" {
		return nil, errors.NewConflictError("This poll is closed for voting", "POLL_CLOSED")
	}

	optionValid := false
	for _, opt := range pubPoll.Options {
		if opt.ID == optionID {
			optionValid = true
			break
		}
	}

	if !optionValid {
		return nil, errors.NewBadRequestError(fmt.Sprintf("Option '%s' does not exist in poll '%s'", optionID, pollID), "INVALID_OPTION_ID")
	}

	voterIdentity := userID
	if voterIdentity == "" {
		if anonVoterID != "" {
			voterIdentity = anonVoterID
		} else {
			h := sha256.Sum256([]byte(clientIP))
			voterIdentity = "anon_ip_" + hex.EncodeToString(h[:12])
		}
	}

	alreadyVoted, err := s.voteRepo.HasVoted(ctx, pollID, voterIdentity)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}
	if alreadyVoted {
		return nil, errors.NewConflictError("You have already voted on this poll", "DUPLICATE_VOTE")
	}

	vote := &Vote{
		PollID:        pollID,
		OptionID:      optionID,
		VoterIdentity: voterIdentity,
		VoterIP:       clientIP,
		UserID:        userID,
		CreatedAt:     time.Now().UTC(),
	}

	if err := s.voteRepo.Create(ctx, vote); err != nil {
		if strings.Contains(err.Error(), "duplicate vote") {
			return nil, errors.NewConflictError("You have already voted on this poll", "DUPLICATE_VOTE")
		}
		return nil, errors.NewInternalError(err)
	}

	_ = s.pollService.IncrementVote(ctx, pollID, optionID)

	// Broadcast live update asynchronously so HTTP vote submission response returns immediately
	go s.realtimeService.PublishUpdate(context.Background(), pollID)

	return &VoteResponse{
		ID:        vote.ID,
		PollID:    vote.PollID,
		OptionID:  vote.OptionID,
		CreatedAt: vote.CreatedAt,
	}, nil
}
