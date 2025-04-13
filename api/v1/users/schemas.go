package users

import (
	"notes-app/models"
	"notes-app/utils"
)

// RegisterResponse is a struct that represents the response for a user registration API.
// It contains the API success, message and the created user details.
type RegisterResponse struct {
	utils.ApiResponse
	User models.User `json:"user"`
}
