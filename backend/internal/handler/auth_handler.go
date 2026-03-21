package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/luisfucros/expense-tracker-app/internal/domain/model"
	appvalidator "github.com/luisfucros/expense-tracker-app/pkg/validator"
)

// AuthHandler handles authentication-related HTTP routes.
type AuthHandler struct {
	*Handler
}

// NewAuthHandler creates an AuthHandler.
func NewAuthHandler(h *Handler) *AuthHandler {
	return &AuthHandler{Handler: h}
}

// Register handles POST /api/v1/auth/register
//
// @Summary     Register a new user
// @Description Creates a new user account and returns a JWT token
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body model.RegisterInput true "Registration payload"
// @Success     201 {object} successResponse{data=model.AuthResponse}
// @Failure     400 {object} errorResponse
// @Failure     409 {object} errorResponse "Email already registered"
// @Router      /auth/register [post]
func (ah *AuthHandler) Register(c *gin.Context) {
	var input model.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"code":    "BAD_REQUEST",
			"message": "invalid request body",
		}})
		return
	}

	if err := appvalidator.Validate(input); err != nil {
		ah.Fail(c, err)
		return
	}

	resp, err := ah.AuthService.Register(c.Request.Context(), input)
	if err != nil {
		ah.Fail(c, err)
		return
	}

	ah.Respond(c, http.StatusCreated, resp)
}

// Login handles POST /api/v1/auth/login
//
// @Summary     Login
// @Description Authenticates a user and returns a JWT token
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body model.LoginInput true "Login credentials"
// @Success     200 {object} successResponse{data=model.AuthResponse}
// @Failure     400 {object} errorResponse
// @Failure     401 {object} errorResponse "Invalid credentials"
// @Router      /auth/login [post]
func (ah *AuthHandler) Login(c *gin.Context) {
	var input model.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"code":    "BAD_REQUEST",
			"message": "invalid request body",
		}})
		return
	}

	if err := appvalidator.Validate(input); err != nil {
		ah.Fail(c, err)
		return
	}

	resp, err := ah.AuthService.Login(c.Request.Context(), input)
	if err != nil {
		ah.Fail(c, err)
		return
	}

	ah.Respond(c, http.StatusOK, resp)
}
