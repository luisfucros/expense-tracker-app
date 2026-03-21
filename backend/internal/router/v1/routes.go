package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luisfucros/expense-tracker-app/internal/handler"
	"github.com/luisfucros/expense-tracker-app/internal/middleware"
)

// RegisterRoutes mounts all v1 API routes onto the provided engine.
func RegisterRoutes(r *gin.Engine, h *handler.Handler, jwtSecret string) {
	authHandler := handler.NewAuthHandler(h)
	expenseHandler := handler.NewExpenseHandler(h)

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

	// Expense routes (require authentication)
	expenses := api.Group("/expenses")
	expenses.Use(middleware.Auth(jwtSecret))
	{
		expenses.GET("", expenseHandler.List)
		expenses.POST("", expenseHandler.Create)
		expenses.GET("/:id", expenseHandler.GetByID)
		expenses.PUT("/:id", expenseHandler.Update)
		expenses.DELETE("/:id", expenseHandler.Delete)
	}
}
