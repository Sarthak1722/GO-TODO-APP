package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func GetRequestID() fiber.Handler {
	return func(c fiber.Ctx) error {
		

		// generate UUID
		requestID := uuid.New().String()

		// store in request context
		c.Locals("request_id", requestID)
		c.Set("X-Request-ID", requestID)

		return c.Next()
	}
}
