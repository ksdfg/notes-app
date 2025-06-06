package service_test

import (
	"fmt"
	"log/slog"
	"notes-app/config"
	"notes-app/service"
	"strconv"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

func TestGenerateJWT(t *testing.T) {
	userID := uint(1)

	token, expiry, err := service.AuthService{}.GenerateJWT(userID)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	slog.Debug("Generated JWT token", slog.String("token", token), slog.Time("expiry", expiry))

	parsedToken, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) { return []byte(config.Get().JWTSecret), nil })
	if err != nil {
		t.Error("Failed to parse JWT token", err)
		return
	}

	claims, ok := parsedToken.Claims.(*jwt.RegisteredClaims)
	if !ok {
		t.Error("Invalid JWT token")
		return
	}

	assert.Equal(t, strconv.FormatUint(uint64(userID), 10), claims.Subject)
	assert.Equal(t, expiry.Unix(), claims.ExpiresAt.Unix())
}

func TestParseJWT(t *testing.T) {
	userID := uint(1)
	
	expiry := time.Now().Add(24 * time.Hour)

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.RegisteredClaims{
		Subject:   fmt.Sprintf("%d", userID),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(expiry),
	})

	signedToken, err := token.SignedString([]byte(config.Get().JWTSecret))
	if err != nil {
		t.Error("Failed to sign token", err)
		return
	}

	claims, err := service.AuthService{}.ParseJWT(signedToken)
	if err != nil {
		t.Error("Failed to parse JWT token", err)
		return
	}

	assert.Equal(t, strconv.FormatUint(uint64(userID), 10), claims.Subject)
}
