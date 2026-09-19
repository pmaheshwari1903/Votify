package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/votify/pkg/errors"
	"github.com/votify/pkg/response"
	val "github.com/votify/pkg/validator"
	"github.com/votify/poll-service/internal/dto"
	"github.com/votify/poll-service/internal/service"
)

type PollHandler struct {
	service service.PollService
}

func NewPollHandler(svc service.PollService) *PollHandler {
	return &PollHandler{
		service: svc,
	}
}

// Create handles POST /polls
func (h *PollHandler) Create(c *gin.Context) {
	ownerID := c.GetHeader("X-User-ID")
	if ownerID == "" {
		response.Unauthorized(c, "Authentication required to create a poll")
		return
	}

	var req dto.CreatePollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			response.ValidationError(c, val.FormatErrors(ve))
			return
		}
		response.BadRequest(c, "Invalid request payload", nil)
		return
	}

	poll, err := h.service.CreatePoll(c.Request.Context(), ownerID, req)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	response.Created(c, poll)
}

// List handles GET /polls
func (h *PollHandler) List(c *gin.Context) {
	ownerID := c.GetHeader("X-User-ID")
	if ownerID == "" {
		response.Unauthorized(c, "Authentication required to list your polls")
		return
	}

	polls, err := h.service.ListOwnerPolls(c.Request.Context(), ownerID)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	response.OK(c, polls)
}

// GetByID handles GET /polls/:id
func (h *PollHandler) GetByID(c *gin.Context) {
	pollID := c.Param("id")
	ownerID := c.GetHeader("X-User-ID")

	poll, err := h.service.GetPollByID(c.Request.Context(), pollID, ownerID)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	response.OK(c, poll)
}

// GetPublic handles GET /polls/public/:id
func (h *PollHandler) GetPublic(c *gin.Context) {
	pollID := c.Param("id")

	pubPoll, err := h.service.GetPublicPoll(c.Request.Context(), pollID)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	response.OK(c, pubPoll)
}

// Update handles PATCH /polls/:id
func (h *PollHandler) Update(c *gin.Context) {
	pollID := c.Param("id")
	ownerID := c.GetHeader("X-User-ID")
	if ownerID == "" {
		response.Unauthorized(c, "Authentication required")
		return
	}

	var req dto.UpdatePollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			response.ValidationError(c, val.FormatErrors(ve))
			return
		}
		response.BadRequest(c, "Invalid request payload", nil)
		return
	}

	poll, err := h.service.UpdatePoll(c.Request.Context(), pollID, ownerID, req)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	response.OK(c, poll)
}

// Delete handles DELETE /polls/:id
func (h *PollHandler) Delete(c *gin.Context) {
	pollID := c.Param("id")
	ownerID := c.GetHeader("X-User-ID")
	if ownerID == "" {
		response.Unauthorized(c, "Authentication required")
		return
	}

	if err := h.service.DeletePoll(c.Request.Context(), pollID, ownerID); err != nil {
		errors.RespondWithError(c, err)
		return
	}

	response.OK(c, gin.H{"message": "Poll deleted successfully"})
}

// Open handles POST /polls/:id/open
func (h *PollHandler) Open(c *gin.Context) {
	pollID := c.Param("id")
	ownerID := c.GetHeader("X-User-ID")
	if ownerID == "" {
		response.Unauthorized(c, "Authentication required")
		return
	}

	poll, err := h.service.OpenPoll(c.Request.Context(), pollID, ownerID)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	response.OK(c, poll)
}

// Close handles POST /polls/:id/close
func (h *PollHandler) Close(c *gin.Context) {
	pollID := c.Param("id")
	ownerID := c.GetHeader("X-User-ID")
	if ownerID == "" {
		response.Unauthorized(c, "Authentication required")
		return
	}

	poll, err := h.service.ClosePoll(c.Request.Context(), pollID, ownerID)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	response.OK(c, poll)
}

// Results handles GET /polls/:id/results
func (h *PollHandler) Results(c *gin.Context) {
	pollID := c.Param("id")

	results, err := h.service.GetPollResults(c.Request.Context(), pollID)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	response.OK(c, results)
}

// InternalIncrementVote handles internal POST /internal/polls/:id/options/:optionId/vote
func (h *PollHandler) InternalIncrementVote(c *gin.Context) {
	pollID := c.Param("id")
	optionID := c.Param("optionId")

	if err := h.service.IncrementVote(c.Request.Context(), pollID, optionID); err != nil {
		errors.RespondWithError(c, err)
		return
	}

	response.OK(c, gin.H{"status": "incremented"})
}
