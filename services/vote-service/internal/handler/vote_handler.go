package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/votify/pkg/errors"
	"github.com/votify/pkg/response"
	val "github.com/votify/pkg/validator"
	"github.com/votify/vote-service/internal/dto"
	"github.com/votify/vote-service/internal/service"
)

type VoteHandler struct {
	service service.VoteService
}

func NewVoteHandler(svc service.VoteService) *VoteHandler {
	return &VoteHandler{
		service: svc,
	}
}

// CastVote handles POST /votes
func (h *VoteHandler) CastVote(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	clientIP := c.ClientIP()

	var req dto.CastVoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			response.ValidationError(c, val.FormatErrors(ve))
			return
		}
		response.BadRequest(c, "Invalid request payload", nil)
		return
	}

	voteResp, err := h.service.CastVote(c.Request.Context(), userID, clientIP, req)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	response.Created(c, voteResp)
}
