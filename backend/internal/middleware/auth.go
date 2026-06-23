package middleware

import (
	"strings"

	"github.com/Sarthak1722/todo_app/internal/logger"
	"github.com/Sarthak1722/todo_app/internal/utils"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/gofiber/fiber/v3"
)

// ClerkAuth is the middleware that protects your routes using the official SDK
func ClerkAuth() fiber.Handler {
	return func(c fiber.Ctx) error {
		requestID, _ := c.Locals("request_id").(string)
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			logger.Log.Warn().Str("requestID", requestID).Msg("missing authorization token")
			return utils.RespondError(c, fiber.StatusUnauthorized, "missing authorization token", nil)
		}

		// Extract the actual token from the "Bearer <token>" string
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Log.Warn().Str("requestID", requestID).Msg("invalid authorization header format")
			return utils.RespondError(c, fiber.StatusUnauthorized, "invalid authorization token format", nil)
		}
		tokenStr := parts[1]

		// Use the official Clerk SDK to verify the token cryptographically
		// c.Context() passes the standard request context for timeouts/cancellations
		claims, err := jwt.Verify(c.Context(), &jwt.VerifyParams{
			Token: tokenStr,
		})

		if err != nil {
			logger.Log.Warn().
				Str("requestID", requestID).
				Err(err).
				Msg("unauthorized request: token verification failed")

			return utils.RespondError(c, fiber.StatusUnauthorized, "invalid or expired authorization token", nil)
		}

		// Inject the authenticated user's ID into the Fiber context
		// claims.Subject is the unique Clerk User ID (e.g., "user_2xyz...")
		c.Locals("auth_user_id", claims.Subject)

		return c.Next()
	}
}
