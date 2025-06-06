package v1

import (
	"notes-app/api/v1/users"
	"notes-app/service"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers v1 routes for the API.
func RegisterRoutes(router fiber.Router, userService service.UserService) {
	// The root path just returns a "Hello, World!" message
	router.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// Register the routes for the users controller
	users.RegisterRoutes(router.Group("/users"), users.Controller{
		UserService: userService,
		AuthService: service.AuthService{},
	})
}
