package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// CustomError represents a structured error that can be returned to the client
type CustomError struct {
	Code    int    `json:"status"`
	Message string `json:"msg"`
	Err     error  `json:"-"` // Don't expose internal error to client
}

func (e *CustomError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// NewBadRequestError creates a new bad request error
func NewBadRequestError(message string, err error) *CustomError {
	return &CustomError{
		Code:    http.StatusBadRequest,
		Message: message,
		Err:     err,
	}
}

// NewInternalServerError creates a new internal server error
func NewInternalServerError(message string, err error) *CustomError {
	return &CustomError{
		Code:    http.StatusInternalServerError,
		Message: message,
		Err:     err,
	}
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(message string, err error) *CustomError {
	return &CustomError{
		Code:    http.StatusNotFound,
		Message: message,
		Err:     err,
	}
}

// ErrorHandler is a Gin middleware that handles panics and errors
func ErrorHandler() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		var err error
		
		switch v := recovered.(type) {
		case *CustomError:
			logrus.WithError(v.Err).Errorf("Custom error: %s", v.Message)
			c.JSON(v.Code, gin.H{
				"status": v.Code,
				"msg":    v.Message,
			})
			c.Abort()
			return
		case error:
			err = v
		case string:
			err = fmt.Errorf("%s", v)
		default:
			err = fmt.Errorf("unknown error: %v", v)
		}

		logrus.Errorf("Panic recovered: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"msg":    "An error occurred internally",
		})
		c.Abort()
	})
}

// AbortWithError is a helper function to abort with a custom error
func AbortWithError(c *gin.Context, err *CustomError) {
	logrus.WithError(err.Err).Errorf("Request error: %s", err.Message)
	c.JSON(err.Code, gin.H{
		"status": err.Code,
		"msg":    err.Message,
	})
	c.Abort()
}