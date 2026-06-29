package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Logger(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		status := c.Writer.Status()
		fields := []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("path", c.FullPath()),
			zap.Int("status", c.Writer.Status()),
			zap.Float64("duration_ms", float64(time.Since(start).Milliseconds())),
		}

		switch {
		case status >= 500:
			log.Error("request", fields...)

		case status >= 400:
			log.Warn("request", fields...)

		default:
			log.Info("request", fields...)
		}
	}
}
