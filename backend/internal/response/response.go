package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse is the standard envelope for all API responses.
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// APIError represents a structured error in the API response.
type APIError struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

// Meta holds optional metadata for paginated or enriched responses.
type Meta struct {
	Page       int   `json:"page,omitempty"`
	PerPage    int   `json:"perPage,omitempty"`
	Total      int64 `json:"total,omitempty"`
	TotalPages int   `json:"totalPages,omitempty"`
}

// OK sends a 200 response with data.
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    data,
	})
}

// Created sends a 201 response with data.
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Data:    data,
	})
}

// Paginated sends a 200 response with data and pagination metadata.
func Paginated(c *gin.Context, data interface{}, meta Meta) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    data,
		Meta:    &meta,
	})
}

// BadRequest sends a 400 error response.
func BadRequest(c *gin.Context, message string, details map[string]string) {
	c.JSON(http.StatusBadRequest, APIResponse{
		Success: false,
		Error: &APIError{
			Code:    "BAD_REQUEST",
			Message: message,
			Details: details,
		},
	})
}

// Unauthorized sends a 401 error response.
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, APIResponse{
		Success: false,
		Error: &APIError{
			Code:    "UNAUTHORIZED",
			Message: message,
		},
	})
}

// Forbidden sends a 403 error response.
func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, APIResponse{
		Success: false,
		Error: &APIError{
			Code:    "FORBIDDEN",
			Message: message,
		},
	})
}

// NotFound sends a 404 error response.
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, APIResponse{
		Success: false,
		Error: &APIError{
			Code:    "NOT_FOUND",
			Message: message,
		},
	})
}

// Conflict sends a 409 error response.
func Conflict(c *gin.Context, message string) {
	c.JSON(http.StatusConflict, APIResponse{
		Success: false,
		Error: &APIError{
			Code:    "CONFLICT",
			Message: message,
		},
	})
}

// ValidationError sends a 422 error response with field-level details.
func ValidationError(c *gin.Context, details map[string]string) {
	c.JSON(http.StatusUnprocessableEntity, APIResponse{
		Success: false,
		Error: &APIError{
			Code:    "VALIDATION_ERROR",
			Message: "One or more fields failed validation.",
			Details: details,
		},
	})
}

// InternalError sends a 500 error response.
func InternalError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, APIResponse{
		Success: false,
		Error: &APIError{
			Code:    "INTERNAL_ERROR",
			Message: "An unexpected error occurred. Please try again later.",
		},
	})
}
