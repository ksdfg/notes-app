package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"notes-app/api/v1"
	"notes-app/service"
	"time"
)

// GenApp initializes and returns a new fiber.App instance to serve the APIs for the application.
func GenApp(userService service.UserService) *fiber.App {
	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler})

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

	api := app.Group("/api")

	// Register v1 APIs
	v1.RegisterRoutes(api.Group("/v1"), userService)

	return app
}
