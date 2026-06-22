package utils

import "github.com/gofiber/fiber/v3"

// SuccessResponse represents a successful API response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
}

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Success bool                   `json:"success"`
	Error   string                 `json:"error"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// RespondSuccess sends a successful response
func RespondSuccess(c fiber.Ctx, statusCode int, data interface{}, message string) error {
	return c.Status(statusCode).JSON(SuccessResponse{
		Success: true,
		Data:    data,
		Message: message,
	})
}

// RespondError sends an error response
func RespondError(c fiber.Ctx, statusCode int, error string, details map[string]interface{}) error {
	return c.Status(statusCode).JSON(ErrorResponse{
		Success: false,
		Error:   error,
		Details: details,
	})
}

// RespondValidationError sends a validation error response
func RespondValidationError(c fiber.Ctx, errors map[string]string) error {
	details := make(map[string]interface{})
	for key, value := range errors {
		details[key] = value
	}
	return RespondError(c, fiber.StatusBadRequest, "validation failed", details)
}
