package middlewares

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// RequestResponseLogger logs the request and response bodies
func RequestResponseLogger(logger *logrus.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// Process request
			err := next(c)
			if err != nil {
				c.Error(err)
			}

			// Calculate execution time
			duration := time.Since(start)

			// Get HTTP method and path
			req := c.Request()
			method := req.Method
			path := req.URL.Path

			// Log request details
			logger.WithFields(logrus.Fields{
				"method":        method,
				"path":          path,
				"executionTime": duration.Milliseconds(),
			}).Info("HTTP Request")

			return nil
		}
	}
}
