package v1

import (
	"notes-app/api/v1/users"
	"notes-app/service"

	"github.com/gofiber/fiber/v2"
)

type Services struct {
	UserService service.IUserService
	AuthService service.IAuthService
}

// RegisterRoutes registers v1 routes for the API.
func RegisterRoutes(router fiber.Router, services Services) {
	// The root path just returns a "Hello, World!" message
	router.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// Register the routes for the users controller
	users.RegisterRoutes(router.Group("/users"), users.Controller{
		UserService: services.UserService,
		AuthService: services.AuthService,
	})
}
