package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	"github.com/luisfucros/expense-tracker-app/internal/config"
	v1 "github.com/luisfucros/expense-tracker-app/internal/router/v1"
	"github.com/luisfucros/expense-tracker-app/internal/handler"
	_ "github.com/luisfucros/expense-tracker-app/docs"
)

// New creates and configures a gin.Engine with all middleware and routes.
func New(cfg *config.Config, h *handler.Handler, log *zap.Logger) *gin.Engine {
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// Apply global middleware
	r.Use(middleware.Recovery(log))
	r.Use(middleware.Logger(log))

	// CORS configuration
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000", "http://frontend:3000"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	r.Use(cors.New(corsConfig))

	// Swagger UI — available at /swagger/index.html (disabled in production)
	if cfg.Env != "production" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	// Mount v1 routes
	v1.RegisterRoutes(r, h, cfg.JWTSecret)

	return r
}
