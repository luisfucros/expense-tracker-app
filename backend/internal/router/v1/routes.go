package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luisfucros/expense-tracker-app/internal/handler"
)

func RegisterRoutes(r *gin.Engine, h *handler.Handler, jwtSecret string) {
	authHandler := handler.NewAuthHandler(h)

	api := r.Group("/api/v1")

	// Health check
	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Auth routes (public)
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}
}
