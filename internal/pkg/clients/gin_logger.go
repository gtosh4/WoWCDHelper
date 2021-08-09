package clients

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Ginzap(logger *zap.Logger, timeFormat string, utc bool, level zapcore.Level) gin.HandlerFunc {
	var log func(msg string, fields ...zap.Field)

	switch level {
	case zap.DebugLevel:
		log = logger.Debug

	case zap.InfoLevel:
		log = logger.Info

	case zap.WarnLevel:
		log = logger.Warn

	case zap.DPanicLevel:
		log = logger.DPanic

	case zap.PanicLevel:
		log = logger.Panic

	case zap.FatalLevel:
		log = logger.Fatal
	}

	return func(c *gin.Context) {
		start := time.Now()
		// some evil middlewares modify this values
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		if utc {
			start = start.UTC()
			end = end.UTC()
		}

		if len(c.Errors) > 0 {
			// Append error field if this is an erroneous request.
			for _, e := range c.Errors.Errors() {
				logger.Error(e)
			}
		} else {
			log(path,
				zap.Int("status", c.Writer.Status()),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", c.ClientIP()),
				zap.String("user-agent", c.Request.UserAgent()),
				zap.String("time", end.Format(timeFormat)),
				zap.Duration("latency", latency),
				zap.String("handler", c.HandlerName()),
			)
		}
	}
}
