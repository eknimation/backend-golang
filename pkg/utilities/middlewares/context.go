package middlewares

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
)

// RequestContext adds a context with timeout to each request
func RequestContext(timeout time.Duration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Create a context with timeout for the request
			ctx, cancel := context.WithTimeout(c.Request().Context(), timeout)
			defer cancel()

			// Set the context in the request
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}

// RequestContextWithCancel adds a cancellable context to each request
func RequestContextWithCancel() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Create a cancellable context for the request
			ctx, cancel := context.WithCancel(c.Request().Context())
			defer cancel()

			// Set the context in the request
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}

// RequestContextWithDeadline adds a context with deadline to each request
func RequestContextWithDeadline(deadline time.Time) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Create a context with deadline for the request
			ctx, cancel := context.WithDeadline(c.Request().Context(), deadline)
			defer cancel()

			// Set the context in the request
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}
