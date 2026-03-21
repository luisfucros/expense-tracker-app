package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Logger returns a Gin middleware that logs each request using zap.
func Logger(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		fields := []zap.Field{
			zap.Int("status", statusCode),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.Duration("latency", latency),
			zap.String("user-agent", c.Request.UserAgent()),
		}

		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				log.Error(e, fields...)
			}
		} else {
			if statusCode >= 500 {
				log.Error("request", fields...)
			} else if statusCode >= 400 {
				log.Warn("request", fields...)
			} else {
				log.Info("request", fields...)
			}
		}
	}
}
