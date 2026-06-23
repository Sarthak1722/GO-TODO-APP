package middleware

import (
	"time"

	"github.com/Sarthak1722/todo_app/internal/logger"
	"github.com/gofiber/fiber/v3"
)

func RequestLogger() fiber.Handler {

	return func(c fiber.Ctx) error {
		logger.Log.Info().Str("route", c.Path()).Msg("Received a " + c.Method() + " request")

		// STEP 1 → start timer
		start := time.Now()

		// STEP 2 → execute next handler
		err := c.Next()

		// STEP 3 → calculate duration
		latency := time.Since(start)

		requestID := c.Locals("request_id").(string)

		// STEP 4 → log everything
		logger.Log.Info().
			Str("requestID", requestID).
			Str("method", c.Method()).
			Str("path", c.Path()).
			Int("status", c.Response().StatusCode()).
			Str("ip", c.IP()).
			Dur("latency", latency).
			Msg("request completed")

		return err
	}
}
