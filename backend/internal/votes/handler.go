package votes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/votify/backend/internal/config"
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

	clientIP := c.ClientIP()
	anonVoterID, _ := c.Cookie("votify_anonymous_voter_id")

	// If unauthenticated public user and no voter cookie present, issue a new HttpOnly anonymous voter ID cookie
	if userID == "" && anonVoterID == "" {
		anonVoterID = "anon_voter_" + uuid.New().String()
		cfg := config.Load()
		isSecure := cfg.Env == "production"

		if isSecure {
			c.SetSameSite(http.SameSiteNoneMode)
		} else {
			c.SetSameSite(http.SameSiteLaxMode)
		}

		c.SetCookie("votify_anonymous_voter_id", anonVoterID, 31536000, "/", "", isSecure, true)
	}

	var req CastVoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid vote payload", map[string]string{"body": err.Error()})
		return
	}

	res, err := h.service.CastVote(c.Request.Context(), userID, anonVoterID, clientIP, req)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	response.Created(c, res)
}
