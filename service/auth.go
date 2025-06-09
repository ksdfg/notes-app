package service

import (
	"fmt"
	"log/slog"
	"notes-app/config"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
	// HashPassword hashes the given password using bcrypt.
	// Returns the hashed password or an error if hashing fails.
	HashPassword(password string) (string, error)

	// ComparePasswords compares a hashed password with a plaintext password.
	// Returns an error if the passwords do not match.
	ComparePasswords(hashedPassword, password string) error

	// GenerateJWT generates a JWT token for the given user ID.
	//
	// The token is signed using the HS512 algorithm and includes the user's ID as the subject,
	// the current time as the issued-at claim, and an expiry time 24 hours from now.
	//
	// Returns the signed JWT token string, the expiry time, or an error if signing fails.
	GenerateJWT(id uint) (string, time.Time, error)

	// ParseJWT parses and validates a JWT token string.
	//
	// It uses the configured JWT secret to validate the token signature.
	// If the token is valid, it returns the registered claims; otherwise, it returns an error.
	ParseJWT(tokenString string) (*jwt.RegisteredClaims, error)

	// GenMiddleware generates a Fiber middleware for JWT authentication.
	//
	// The middleware validates the JWT token from the request's authorization cookie.
	// If the token is valid, it sets the user ID in the context for further use.
	GenMiddleware() fiber.Handler
}

var (
	ErrFailedToSignToken = fmt.Errorf("failed to sign token")
	ErrInvalidToken      = fmt.Errorf("invalid token")
)

type AuthService struct{}

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
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		slog.Error("Password comparison failed", slog.Any("error", err))
	}

	return err
}

// GenerateJWT generates a JWT token for the given user ID.
//
// The token is signed using the HS512 algorithm and includes the user's ID as the subject,
// the current time as the issued-at claim, and an expiry time 24 hours from now.
//
// Returns the signed JWT token string, the expiry time, or an error if signing fails.
func (svc AuthService) GenerateJWT(id uint) (string, time.Time, error) {
	expiry := time.Now().Add(24 * time.Hour)

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.RegisteredClaims{
		Subject:   fmt.Sprintf("%d", id),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(expiry),
	})

	signedToken, err := token.SignedString([]byte(config.Get().JWTSecret))
	if err != nil {
		slog.Error("Failed to sign token", slog.Any("error", err))
		return "", expiry, ErrFailedToSignToken
	}

	return signedToken, expiry, nil
}

// ParseJWT parses and validates a JWT token string.
//
// It uses the configured JWT secret to validate the token signature.
// If the token is valid, it returns the registered claims; otherwise, it returns an error.
func (svc AuthService) ParseJWT(tokenString string) (*jwt.RegisteredClaims, error) {
	// keyFunc provides the secret key for validating the token signature.
	keyFunc := func(t *jwt.Token) (any, error) { return []byte(config.Get().JWTSecret), nil }

	// Parse the token with the expected claims structure.
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, keyFunc)
	if err != nil {
		slog.Error("Failed to parse JWT token", slog.Any("error", err))
		return nil, err
	}

	// Assert the claims type and check token validity.
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		slog.Error("Invalid JWT token", slog.String("token", tokenString), slog.Any("claims", claims))
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// GenMiddleware generates a Fiber middleware for JWT authentication.
//
// The middleware validates the JWT token from the request's authorization cookie.
// If the token is valid, it sets the user ID in the context for further use.
func (svc AuthService) GenMiddleware() fiber.Handler {
	return keyauth.New(keyauth.Config{
		Validator: func(c *fiber.Ctx, key string) (bool, error) {
			// Parse the JWT token from the key
			claims, err := svc.ParseJWT(key)
			if err != nil {
				return false, err // Return false if token parsing fails
			}

			// Set the user ID in the context for further use
			c.Locals("userID", claims.Subject)

			return true, nil // Return true if authentication is successful
		},
		KeyLookup: "cookie:authorization",
	})
}
