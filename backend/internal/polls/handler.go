package polls

import (
	"github.com/gin-gonic/gin"
	"github.com/votify/backend/internal/errors"
	"github.com/votify/backend/internal/response"
)

type PollHandler struct {
	service PollService
}

func NewPollHandler(service PollService) *PollHandler {
	return &PollHandler{service: service}
}

func getUserID(c *gin.Context) string {
	if uid := c.GetString("userID"); uid != "" {
		return uid
	}
	return c.GetHeader("X-User-ID")
}

func (h *PollHandler) Create(c *gin.Context) {
	ownerID := getUserID(c)
	if ownerID == "" {
		response.Unauthorized(c, "Authentication required to create a poll")
		return
	}

	var req CreatePollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request payload", map[string]string{"body": err.Error()})
		return
	}

	poll, err := h.service.CreatePoll(c.Request.Context(), ownerID, req)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	response.Created(c, poll)
}

func (h *PollHandler) List(c *gin.Context) {
	ownerID := getUserID(c)
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

func (h *PollHandler) GetByID(c *gin.Context) {
	pollID := c.Param("id")
	ownerID := getUserID(c)

	poll, err := h.service.GetPollByID(c.Request.Context(), pollID, ownerID)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	response.OK(c, poll)
}

func (h *PollHandler) GetPublic(c *gin.Context) {
	pollID := c.Param("id")

	pubPoll, err := h.service.GetPublicPoll(c.Request.Context(), pollID)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	response.OK(c, pubPoll)
}

func (h *PollHandler) Update(c *gin.Context) {
	pollID := c.Param("id")
	ownerID := getUserID(c)
	if ownerID == "" {
		response.Unauthorized(c, "Authentication required")
		return
	}

	var req UpdatePollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request payload", map[string]string{"body": err.Error()})
		return
	}

	poll, err := h.service.UpdatePoll(c.Request.Context(), pollID, ownerID, req)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	response.OK(c, poll)
}

func (h *PollHandler) Delete(c *gin.Context) {
	pollID := c.Param("id")
	ownerID := getUserID(c)
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

func (h *PollHandler) Open(c *gin.Context) {
	pollID := c.Param("id")
	ownerID := getUserID(c)
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

func (h *PollHandler) Close(c *gin.Context) {
	pollID := c.Param("id")
	ownerID := getUserID(c)
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

func (h *PollHandler) Results(c *gin.Context) {
	pollID := c.Param("id")

	results, err := h.service.GetPollResults(c.Request.Context(), pollID)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	response.OK(c, results)
}
