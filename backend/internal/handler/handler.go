package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/luisfucros/expense-tracker-app/internal/domain/apierror"
	"github.com/luisfucros/expense-tracker-app/internal/service"
)

// successResponse is the envelope for successful JSON responses.
// @Description Wrapper for successful API responses.
type successResponse struct {
	Data any `json:"data"`
}

// errorResponse is the envelope for error JSON responses.
// @Description Wrapper for API error responses.
type errorResponse struct {
	Error apierror.APIError `json:"error"`
}

// Handler holds shared dependencies for all HTTP handlers.
type Handler struct {
	AuthService    service.AuthService
	// ExpenseService service.ExpenseService
	Log            *zap.Logger
}

// NewHandler creates a Handler with the provided services and logger.
func NewHandler(authSvc service.AuthService, log *zap.Logger) *Handler {
	return &Handler{
		AuthService:    authSvc,
		// ExpenseService: expenseSvc,
		Log:            log,
	}
}

// Respond writes a JSON success response with the given status code and data payload.
func (h *Handler) Respond(c *gin.Context, status int, data any) {
	c.JSON(status, gin.H{"data": data})
}

// Fail maps an error to the appropriate HTTP status code and JSON error response.
func (h *Handler) Fail(c *gin.Context, err error) {
	var apiErr *apierror.APIError
	if errors.As(err, &apiErr) {
		status := codeToStatus(apiErr.Code)
		c.AbortWithStatusJSON(status, gin.H{"error": apiErr})
		return
	}

	h.Log.Error("unexpected error", zap.Error(err))
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"error": apierror.Internal("internal server error"),
	})
}

func codeToStatus(code string) int {
	switch code {
	case apierror.CodeNotFound:
		return http.StatusNotFound
	case apierror.CodeUnauthorized:
		return http.StatusUnauthorized
	case apierror.CodeForbidden:
		return http.StatusForbidden
	case apierror.CodeConflict:
		return http.StatusConflict
	case apierror.CodeBadRequest, apierror.CodeValidation, apierror.CodeInvalidCategory:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
