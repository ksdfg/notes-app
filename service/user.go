package service

import (
	"log/slog"
	"notes-app/models"
)

type IUserService interface {
	// Create creates a new user record in the database.
	// The user's password is hashed before saving.
	// Accepts optional DBOpts to specify a DB instance.
	Create(user *models.User, opts *DBOpts) error

	// GetByEmail retrieves a user by their email from the database.
	// Accepts optional DBOpts to specify a DB instance.
	// Returns the user or an error if the user is not found.
	GetByEmail(email string, opts *DBOpts) (models.User, error)

	// GetByID retrieves a user by their ID from the database.
	// Accepts optional DBOpts to specify a DB instance.
	// Returns the user or an error if the user is not found.
	GetByID(id uint, opts *DBOpts) (models.User, error)
}

type UserService struct {
	Service
}

// Create creates a new user record in the database.
// The user's password is hashed before saving.
//
// Accepts optional DBOpts to specify a DB instance.
func (svc UserService) Create(user *models.User, opts *DBOpts) error {
	db := svc.getDB(opts)

	var err error
	user.Password, err = AuthService{}.HashPassword(user.Password)
	if err != nil {
		return err
	}

	result := db.Create(user)
	if result.Error != nil {
		slog.Error("Failed to create user", slog.Any("error", result.Error))
	}

	return result.Error
}

// GetByID retrieves a user by their ID from the database.
// Accepts optional DBOpts to specify a DB instance.
//
// Returns the user or an error if the user is not found.
func (svc UserService) GetByID(id uint, opts *DBOpts) (models.User, error) {
	db := svc.getDB(opts)

	var user models.User
	result := db.Where("id = ?", id).First(&user)
	if result.Error != nil {
		slog.Error("Failed to fetch user", slog.Any("error", result.Error))
	}

	return user, result.Error
}

// GetByEmail retrieves a user by their email from the database.
// Accepts optional DBOpts to specify a DB instance.
//
// Returns the user or an error if the user is not found.
func (svc UserService) GetByEmail(email string, opts *DBOpts) (models.User, error) {
	db := svc.getDB(opts)

	var user models.User
	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		slog.Error("Failed to fetch user", slog.Any("error", result.Error))
	}

	return user, result.Error
}
