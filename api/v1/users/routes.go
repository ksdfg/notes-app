package users

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(router fiber.Router, controller Controller) {
	router.Post("/", controller.Register)
	router.Post("/login", controller.Login)
}
