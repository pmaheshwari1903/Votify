package errors

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/votify/backend/internal/response"
)

type ErrorType string

const (
	ErrorTypeValidation   ErrorType = "VALIDATION_ERROR"
	ErrorTypeNotFound     ErrorType = "NOT_FOUND"
	ErrorTypeUnauthorized ErrorType = "UNAUTHORIZED"
	ErrorTypeForbidden    ErrorType = "FORBIDDEN"
	ErrorTypeConflict     ErrorType = "CONFLICT"
	ErrorTypeBadRequest   ErrorType = "BAD_REQUEST"
	ErrorTypeInternal     ErrorType = "INTERNAL_ERROR"
)

type AppError struct {
	Type    ErrorType         `json:"type"`
	Message string            `json:"message"`
	Code    string            `json:"code"`
	Details map[string]string `json:"details,omitempty"`
	Err     error             `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Type, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewBadRequestError(message string, code string) *AppError {
	if code == "" {
		code = string(ErrorTypeBadRequest)
	}
	return &AppError{
		Type:    ErrorTypeBadRequest,
		Message: message,
		Code:    code,
	}
}

func NewValidationError(message string, details map[string]string) *AppError {
	return &AppError{
		Type:    ErrorTypeValidation,
		Message: message,
		Code:    string(ErrorTypeValidation),
		Details: details,
	}
}

func NewNotFoundError(resource string, id string) *AppError {
	msg := fmt.Sprintf("%s not found", resource)
	if id != "" {
		msg = fmt.Sprintf("%s with id '%s' not found", resource, id)
	}
	return &AppError{
		Type:    ErrorTypeNotFound,
		Message: msg,
		Code:    string(ErrorTypeNotFound),
	}
}

func NewUnauthorizedError(message string) *AppError {
	if message == "" {
		message = "Authentication required"
	}
	return &AppError{
		Type:    ErrorTypeUnauthorized,
		Message: message,
		Code:    string(ErrorTypeUnauthorized),
	}
}

func NewForbiddenError(message string) *AppError {
	if message == "" {
		message = "Access denied"
	}
	return &AppError{
		Type:    ErrorTypeForbidden,
		Message: message,
		Code:    string(ErrorTypeForbidden),
	}
}

func NewConflictError(message string, code string) *AppError {
	if code == "" {
		code = string(ErrorTypeConflict)
	}
	return &AppError{
		Type:    ErrorTypeConflict,
		Message: message,
		Code:    code,
	}
}

func NewInternalError(err error) *AppError {
	return &AppError{
		Type:    ErrorTypeInternal,
		Message: "An unexpected error occurred",
		Code:    string(ErrorTypeInternal),
		Err:     err,
	}
}

func RespondWithError(c *gin.Context, err error) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		switch appErr.Type {
		case ErrorTypeBadRequest:
			c.JSON(http.StatusBadRequest, response.APIResponse{
				Success: false,
				Error: &response.APIError{
					Code:    appErr.Code,
					Message: appErr.Message,
					Details: appErr.Details,
				},
			})
		case ErrorTypeValidation:
			response.ValidationError(c, appErr.Details)
		case ErrorTypeNotFound:
			response.NotFound(c, appErr.Message)
		case ErrorTypeUnauthorized:
			response.Unauthorized(c, appErr.Message)
		case ErrorTypeForbidden:
			response.Forbidden(c, appErr.Message)
		case ErrorTypeConflict:
			response.Conflict(c, appErr.Message)
		default:
			response.InternalError(c)
		}
		return
	}

	response.InternalError(c)
}
