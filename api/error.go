package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"notes-app/utils"
)

// ErrorHandler formats any errors as a proper JSON response
var ErrorHandler = func(c *fiber.Ctx, err error) error {
	// Status code defaults to 500
	code := fiber.StatusInternalServerError

	// Retrieve the custom status code if it's a *fiber.Error
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	// Set Content-Type: text/plain; charset=utf-8
	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

	// Return status code with error message
	return c.Status(code).JSON(utils.ApiResponse{
		Success: false,
		Message: err.Error(),
	})
}
