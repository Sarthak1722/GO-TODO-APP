package handlers

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v3"
)

var errMissingAuthUserID = errors.New("missing authenticated user id")

// authenticatedUser returns the Clerk user ID from Fiber locals and a context
// enriched with the same identity for downstream layers.
func authenticatedUser(c fiber.Ctx) (string, context.Context, error) {
	userID, ok := c.Locals("auth_user_id").(string)
	if !ok || userID == "" {
		return "", nil, errMissingAuthUserID
	}

	return userID, c.Context(), nil
}
