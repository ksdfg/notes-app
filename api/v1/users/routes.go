package users

import (
	"github.com/gofiber/fiber/v2"
	"notes-app/service"
)

func RegisterRoutes(router fiber.Router, service service.UserService) {
	controller := Controller{userService: service}

	router.Post("/", controller.Register)
}
