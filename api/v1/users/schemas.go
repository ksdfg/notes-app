package users

import (
	"notes-app/models"
	"notes-app/utils"
)

// RegisterRequest is a struct that represents the request for a user registration API.
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterResponse is a struct that represents the response for a user registration API.
type RegisterResponse struct {
	utils.ApiResponse
	User models.User `json:"user"`
}

// LoginRequest is a struct that represents the request for a user login API.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
