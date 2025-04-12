package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"time"
)

// GenApp initializes and returns a new fiber.App instance to serve the APIs for the application.
func GenApp() *fiber.App {
	app := fiber.New()

	// Recover middleware recovers from panics anywhere in the app
	app.Use(recover.New())

	// RequestID middleware generates a unique ID for each request
	app.Use(requestid.New())

	// Logger middleware logs HTTP requests
	app.Use(logger.New(logger.Config{
		Format:     "${locals:requestid} | ${time} | ${status} - ${method} ${path}\n",
		TimeFormat: time.RFC3339,
	}))

	// Define a GET route for the root path
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	return app
}
