package middleware

import (
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"example.com/rest-core-arch/pkg/util/log"
)

// ZapLogger middleware logs HTTP requests using Zap.
func ZapLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// Extract or Generate Trace ID
			reqID := c.Request().Header.Get(echo.HeaderXRequestID)
			if reqID == "" {
				reqID = uuid.New().String()
				c.Request().Header.Set(echo.HeaderXRequestID, reqID)
			}
			c.Response().Header().Set(echo.HeaderXRequestID, reqID)

			// Process request
			err := next(c)

			// Calculate latency
			latency := time.Since(start)

			// Log fields
			fields := []zap.Field{
				zap.String("req_id", reqID),
				zap.String("method", c.Request().Method),
				zap.String("path", c.Path()),
				zap.Int("status", c.Response().Status),
				zap.Duration("latency", latency),
				zap.String("ip", c.RealIP()),
				zap.String("user_agent", c.Request().UserAgent()),
			}

			if err != nil {
				fields = append(fields, zap.Error(err))
				// If the error is handled by Echo's HTTPErrorHandler, it might not return here
				// but Echo usually returns the error from next(c)
				c.Error(err) // Ensure response is sent if next(c) returned raw error
			}

			log.L().Info("http_request", fields...)

			return nil
		}
	}
}
