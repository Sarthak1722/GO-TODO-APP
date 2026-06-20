package middleware

import (
	"github.com/Sarthak1722/todo_app/internal/logger"
	"github.com/Sarthak1722/todo_app/internal/utils"
	"github.com/gofiber/fiber/v3"
)

// RecoverPanic recovers from panics and returns a proper error response
func RecoverPanic() fiber.Handler {
	return func(c fiber.Ctx) error {
		defer func() {
			if err := recover(); err != nil {
				requestID := c.Locals("request_id").(string)
				logger.Log.Error().
					Str("requestID", requestID).
					Interface("panic", err).
					Msg("panic recovered")

				utils.RespondError(
					c,
					fiber.StatusInternalServerError,
					"internal server error",
					nil,
				)
			}
		}()
		return c.Next()
	}
}
