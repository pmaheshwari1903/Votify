package votes

import (
	"github.com/gin-gonic/gin"
	"github.com/votify/backend/internal/errors"
	"github.com/votify/backend/internal/response"
)

type VoteHandler struct {
	service VoteService
}

func NewVoteHandler(service VoteService) *VoteHandler {
	return &VoteHandler{service: service}
}

func (h *VoteHandler) CastVote(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		userID = c.GetHeader("X-User-ID")
	}

	clientIP := c.ClientIP()

	var req CastVoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid vote payload", map[string]string{"body": err.Error()})
		return
	}

	res, err := h.service.CastVote(c.Request.Context(), userID, clientIP, req)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	response.Created(c, res)
}
