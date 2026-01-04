package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrBadRequest     = NewAPIError(http.StatusBadRequest, "Bad request")
	ErrUnauthorized   = NewAPIError(http.StatusUnauthorized, "Unauthorized")
	ErrForbidden      = NewAPIError(http.StatusForbidden, "Forbidden")
	ErrNotFound       = NewAPIError(http.StatusNotFound, "Resource not found")
	ErrAlreadyExists  = NewAPIError(http.StatusConflict, "Already exists")
	ErrInternalServer = NewAPIError(http.StatusInternalServerError, "Internal server error")
)

// APIError defines a structured API error
type APIError struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
	Code       int    `json:"code,omitempty"`
}

// NewAPIError creates a new APIError
func NewAPIError(status int, msg string, code ...int) *APIError {
	err := &APIError{
		StatusCode: status,
		Message:    msg,
	}
	if len(code) > 0 {
		err.Code = code[0]
	}
	return err
}

// SendError sends the error response and aborts the request
func SendError(c *gin.Context, err *APIError) {
	c.JSON(err.StatusCode, gin.H{
		"error": err.Message,
		"code":  err.Code,
	})
	c.Abort()
}

// SendSuccess sends a success response with optional data
func SendSuccess(c *gin.Context, message string, data interface{}, status ...int) {
	resStatus := http.StatusOK
	if len(status) > 0 {
		resStatus = status[0]
	}

	c.JSON(resStatus, gin.H{
		"message": message,
		"data":    data,
	})
}
