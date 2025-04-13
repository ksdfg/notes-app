package users

import (
	"notes-app/models"
	"notes-app/utils"
)

// RegisterRequest is a struct that represents the request for a user registration API.
// It contains the user's name, email and password.
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterResponse is a struct that represents the response for a user registration API.
// It contains the API success, message and the created user details.
type RegisterResponse struct {
	utils.ApiResponse
	User models.User `json:"user"`
}
