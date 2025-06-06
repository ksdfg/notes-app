package service_test

import (
	"log/slog"
	"notes-app/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	password := "password"

	// Hash the password using the service
	hashedPassword, errHash := service.AuthService{}.HashPassword(password)
	assert.NoError(t, errHash)
	assert.NotEmpty(t, hashedPassword)
	slog.Debug("Hashed password", slog.String("hash", hashedPassword))

	// Compare the hashed password
	errCompare := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	assert.NoError(t, errCompare)

	// Compare with an incorrect password
	errIncorrectCompare := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte("wrongpassword"))
	assert.Error(t, errIncorrectCompare)
}

func TestComparePassword(t *testing.T) {
	password := "password"

	// Hash the password
	hashedPassword, errHash := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, errHash)
	assert.NotEmpty(t, hashedPassword)
	slog.Debug("Hashed password", slog.String("hash", string(hashedPassword)))

	// Compare the hashed password using the service
	errCompare := service.AuthService{}.ComparePasswords(string(hashedPassword), password)
	assert.NoError(t, errCompare)

	// Compare with an incorrect password
	errIncorrectCompare := service.AuthService{}.ComparePasswords(string(hashedPassword), "wrongpassword")
	assert.Error(t, errIncorrectCompare)
}
