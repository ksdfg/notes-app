package service

import (
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"notes-app/models"
)

type IUserService interface {
	// HashPassword hashes the given password using bcrypt.
	// Returns the hashed password or an error if hashing fails.
	HashPassword(password string) (string, error)

	// ComparePasswords compares a hashed password with a plaintext password.
	// Returns an error if the passwords do not match.
	ComparePasswords(hashedPassword, password string) error

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

// HashPassword hashes the given password using bcrypt.
//
// Returns the hashed password or an error if hashing fails.
func (svc UserService) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Failed to hash password", slog.Any("error", err))
		return "", err
	}

	return string(hashedPassword), nil
}

// ComparePasswords compares a hashed password with a plaintext password.
//
// Returns an error if the passwords do not match.
func (svc UserService) ComparePasswords(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// Create creates a new user record in the database.
// The user's password is hashed before saving.
//
// Accepts optional DBOpts to specify a DB instance.
func (svc UserService) Create(user *models.User, opts *DBOpts) error {
	db := svc.getDB(opts)

	var err error
	user.Password, err = svc.HashPassword(user.Password)
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
		slog.Error("Failed to create user", slog.Any("error", result.Error))
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
		slog.Error("Failed to create user", slog.Any("error", result.Error))
	}

	return user, result.Error
}
