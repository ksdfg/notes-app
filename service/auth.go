package service

import (
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
	// HashPassword hashes the given password using bcrypt.
	// Returns the hashed password or an error if hashing fails.
	HashPassword(password string) (string, error)

	// ComparePasswords compares a hashed password with a plaintext password.
	// Returns an error if the passwords do not match.
	ComparePasswords(hashedPassword, password string) error
}

type AuthService struct {}

// HashPassword hashes the given password using bcrypt.
//
// Returns the hashed password or an error if hashing fails.
func (svc AuthService) HashPassword(password string) (string, error) {
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
func (svc AuthService) ComparePasswords(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}