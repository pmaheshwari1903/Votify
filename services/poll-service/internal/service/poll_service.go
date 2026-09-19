package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/votify/pkg/errors"
	"github.com/votify/poll-service/internal/dto"
	"github.com/votify/poll-service/internal/model"
	"github.com/votify/poll-service/internal/repository"
)

type PollService interface {
	CreatePoll(ctx context.Context, ownerID string, req dto.CreatePollRequest) (*model.Poll, error)
	ListOwnerPolls(ctx context.Context, ownerID string) ([]*model.Poll, error)
	GetPollByID(ctx context.Context, id string, ownerID string) (*model.Poll, error)
	GetPublicPoll(ctx context.Context, id string) (*dto.PublicPollResponse, error)
	UpdatePoll(ctx context.Context, id string, ownerID string, req dto.UpdatePollRequest) (*model.Poll, error)
	DeletePoll(ctx context.Context, id string, ownerID string) error
	OpenPoll(ctx context.Context, id string, ownerID string) (*model.Poll, error)
	ClosePoll(ctx context.Context, id string, ownerID string) (*model.Poll, error)
	GetPollResults(ctx context.Context, id string) (*dto.PollResultsResponse, error)
	IncrementVote(ctx context.Context, pollID string, optionID string) error
}

type DefaultPollService struct {
	repo repository.PollRepository
}

func NewPollService(repo repository.PollRepository) *DefaultPollService {
	return &DefaultPollService{
		repo: repo,
	}
}

func (s *DefaultPollService) CreatePoll(ctx context.Context, ownerID string, req dto.CreatePollRequest) (*model.Poll, error) {
	if ownerID == "" {
		return nil, errors.NewUnauthorizedError("Authentication required to create a poll")
	}

	question := strings.TrimSpace(req.Question)
	if len(question) < 5 || len(question) > 300 {
		return nil, errors.NewBadRequestError("Poll question must be between 5 and 300 characters", "INVALID_QUESTION_LENGTH")
	}

	if len(req.Options) < 2 || len(req.Options) > 10 {
		return nil, errors.NewBadRequestError("Poll must contain between 2 and 10 options", "INVALID_OPTION_COUNT")
	}

	seenText := make(map[string]bool)
	var options []model.PollOption

	for i, optText := range req.Options {
		trimmed := strings.TrimSpace(optText)
		if trimmed == "" {
			return nil, errors.NewBadRequestError(fmt.Sprintf("Option %d cannot be empty", i+1), "EMPTY_OPTION")
		}
		if len(trimmed) > 150 {
			return nil, errors.NewBadRequestError(fmt.Sprintf("Option %d exceeds maximum length of 150 characters", i+1), "OPTION_TOO_LONG")
		}
		lower := strings.ToLower(trimmed)
		if seenText[lower] {
			return nil, errors.NewConflictError(fmt.Sprintf("Duplicate option text '%s' is not allowed", trimmed), "DUPLICATE_OPTION")
		}
		seenText[lower] = true

		options = append(options, model.PollOption{
			ID:        fmt.Sprintf("opt_%d", i+1),
			Text:      trimmed,
			VoteCount: 0,
		})
	}

	poll := &model.Poll{
		OwnerID:     ownerID,
		Question:    question,
		Description: strings.TrimSpace(req.Description),
		Options:     options,
		Status:      "open", // Default status is open for voting
		TotalVotes:  0,
	}

	if err := s.repo.Create(ctx, poll); err != nil {
		return nil, errors.NewInternalError(err)
	}

	return poll, nil
}

func (s *DefaultPollService) ListOwnerPolls(ctx context.Context, ownerID string) ([]*model.Poll, error) {
	if ownerID == "" {
		return nil, errors.NewUnauthorizedError("Authentication required")
	}
	polls, err := s.repo.FindByOwnerID(ctx, ownerID)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}
	if polls == nil {
		polls = []*model.Poll{}
	}
	return polls, nil
}

func (s *DefaultPollService) GetPollByID(ctx context.Context, id string, ownerID string) (*model.Poll, error) {
	poll, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}
	if poll == nil {
		return nil, errors.NewNotFoundError("Poll", id)
	}

	if ownerID != "" && poll.OwnerID != ownerID {
		return nil, errors.NewForbiddenError("You do not have permission to view or manage this poll")
	}

	return poll, nil
}

func (s *DefaultPollService) GetPublicPoll(ctx context.Context, id string) (*dto.PublicPollResponse, error) {
	poll, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}
	if poll == nil {
		return nil, errors.NewNotFoundError("Poll", id)
	}

	if poll.Status == "draft" {
		return nil, errors.NewForbiddenError("This poll is currently in draft mode and not public")
	}

	// Sanitize response — exclude owner ID and private metrics
	return &dto.PublicPollResponse{
		ID:          poll.ID,
		Question:    poll.Question,
		Description: poll.Description,
		Options:     poll.Options,
		Status:      poll.Status,
		CreatedAt:   poll.CreatedAt,
	}, nil
}

func (s *DefaultPollService) UpdatePoll(ctx context.Context, id string, ownerID string, req dto.UpdatePollRequest) (*model.Poll, error) {
	poll, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}
	if poll == nil {
		return nil, errors.NewNotFoundError("Poll", id)
	}

	if poll.OwnerID != ownerID {
		return nil, errors.NewForbiddenError("You are not authorized to update this poll")
	}

	if req.Question != "" {
		poll.Question = strings.TrimSpace(req.Question)
	}
	if req.Description != "" {
		poll.Description = strings.TrimSpace(req.Description)
	}

	if err := s.repo.Update(ctx, poll); err != nil {
		return nil, errors.NewInternalError(err)
	}

	return poll, nil
}

func (s *DefaultPollService) DeletePoll(ctx context.Context, id string, ownerID string) error {
	poll, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return errors.NewInternalError(err)
	}
	if poll == nil {
		return errors.NewNotFoundError("Poll", id)
	}

	if poll.OwnerID != ownerID {
		return errors.NewForbiddenError("You are not authorized to delete this poll")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return errors.NewInternalError(err)
	}

	return nil
}

func (s *DefaultPollService) OpenPoll(ctx context.Context, id string, ownerID string) (*model.Poll, error) {
	poll, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}
	if poll == nil {
		return nil, errors.NewNotFoundError("Poll", id)
	}

	if poll.OwnerID != ownerID {
		return nil, errors.NewForbiddenError("You are not authorized to modify this poll")
	}

	if poll.Status == "open" {
		return poll, nil
	}

	poll.Status = "open"
	if err := s.repo.Update(ctx, poll); err != nil {
		return nil, errors.NewInternalError(err)
	}

	return poll, nil
}

func (s *DefaultPollService) ClosePoll(ctx context.Context, id string, ownerID string) (*model.Poll, error) {
	poll, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}
	if poll == nil {
		return nil, errors.NewNotFoundError("Poll", id)
	}

	if poll.OwnerID != ownerID {
		return nil, errors.NewForbiddenError("You are not authorized to modify this poll")
	}

	if poll.Status == "closed" {
		return poll, nil
	}

	poll.Status = "closed"
	if err := s.repo.Update(ctx, poll); err != nil {
		return nil, errors.NewInternalError(err)
	}

	return poll, nil
}

func (s *DefaultPollService) GetPollResults(ctx context.Context, id string) (*dto.PollResultsResponse, error) {
	poll, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}
	if poll == nil {
		return nil, errors.NewNotFoundError("Poll", id)
	}

	return &dto.PollResultsResponse{
		PollID:     poll.ID,
		Question:   poll.Question,
		Status:     poll.Status,
		TotalVotes: poll.TotalVotes,
		Options:    poll.Options,
	}, nil
}

func (s *DefaultPollService) IncrementVote(ctx context.Context, pollID string, optionID string) error {
	return s.repo.IncrementVote(ctx, pollID, optionID)
}
