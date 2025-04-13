package users

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"log/slog"
	"notes-app/models"
	"notes-app/service"
	"notes-app/utils"
)

// Controller defines the handlers for the v1/users API.
type Controller struct {
	userService service.IUserService
}

// Register creates a new user in the database.
//
// The request body should contain the user data.
//
// Returns a 201 Created response with the created user in the response body.
func (c Controller) Register(ctx *fiber.Ctx) error {
	// Parse the request body into a User struct
	user := new(models.User)

	if err := ctx.BodyParser(user); err != nil {
		slog.Error("Failed to parse request body", slog.Any("error", err))
		// Return a 400 Bad Request response if the request body is invalid
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Register the user in the database
	if err := c.userService.Create(user, nil); err != nil {
		// Return a 409 Conflict response if the user already exists
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return fiber.NewError(fiber.StatusConflict, "User already exists")
		}
		// Return the error if anything else goes wrong
		return err
	}

	// Return a 201 Created response with the created user in the response body
	return ctx.Status(fiber.StatusCreated).JSON(RegisterResponse{
		ApiResponse: utils.ApiResponse{
			Success: true,
			Message: "User created successfully",
		},
		User: *user,
	})
}
