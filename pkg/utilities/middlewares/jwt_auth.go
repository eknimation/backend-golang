package middlewares

import (
	"net/http"
	"strings"
	"time"

	"backend-service/config"
	"backend-service/pkg/utilities/jwt"
	"backend-service/pkg/utilities/responses"

	"github.com/labstack/echo/v4"
)

// JWTAuth middleware for protecting routes with JWT authentication
func JWTAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, responses.Error(http.StatusUnauthorized, "Authorization header required"))
			}

			var tokenString string
			// Check if it starts with "Bearer "
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			} else {
				// If no "Bearer " prefix, assume the entire header value is the token
				// This allows Swagger UI to work without requiring users to add "Bearer " manually
				tokenString = authHeader
			}

			if tokenString == "" {
				return c.JSON(http.StatusUnauthorized, responses.Error(http.StatusUnauthorized, "Token is required"))
			}

			// Get JWT configuration
			appConfig := config.GetAppConfig()
			if appConfig.JWTSecret == "" {
				return c.JSON(http.StatusInternalServerError, responses.Error(http.StatusInternalServerError, "JWT secret not configured"))
			}

			// Create JWT manager and verify token
			jwtManager := jwt.NewJWTManager(appConfig.JWTSecret, 24*time.Hour)
			claims, err := jwtManager.VerifyToken(tokenString)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, responses.Error(http.StatusUnauthorized, "Invalid token"))
			}

			// Store user information in context for use in handlers
			c.Set("user_id", claims.UserID)
			c.Set("user_email", claims.Email)

			return next(c)
		}
	}
}
