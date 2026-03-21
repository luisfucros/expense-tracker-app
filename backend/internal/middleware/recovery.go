package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/luisfucros/expense-tracker-app/internal/domain/apierror"
)

// Recovery returns a Gin middleware that recovers from panics and logs them.
func Recovery(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Error("panic recovered", zap.Any("error", err))
				c.AbortWithStatusJSON(500, gin.H{"error": apierror.Internal("internal server error")})
			}
		}()
		c.Next()
	}
}
