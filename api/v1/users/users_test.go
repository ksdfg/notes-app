package users_test

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"notes-app/api"
	"notes-app/api/v1/users"
	"notes-app/models"
	"notes-app/service"
	"notes-app/utils"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type mockUserService struct{}

func (svc mockUserService) Create(user *models.User, opts *service.DBOpts) error {
	if user.Email == "duplicate@ksdfg.dev" {
		return gorm.ErrDuplicatedKey
	}

	user.ID = 1
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	return nil
}

func (svc mockUserService) GetByEmail(email string, opts *service.DBOpts) (models.User, error) {
	if email == "nosuchuser@ksdfg.dev" {
		return models.User{}, gorm.ErrRecordNotFound
	}

	return models.User{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:     "Kshitish Deshpande",
		Email:    "me@ksdfg.dev",
		Password: "hashedpassword",
	}, nil
}

func (svc mockUserService) GetByID(id uint, opts *service.DBOpts) (models.User, error) {
	//TODO implement me
	panic("implement me")
}

type mockAuthService struct{}

func (svc mockAuthService) HashPassword(password string) (string, error) {
	return "hashedpassword", nil
}

func (svc mockAuthService) ComparePasswords(hashedPassword string, password string) error {
	if hashedPassword == "hashedpassword" && password == "securepassword" {
		return nil
	}

	return bcrypt.ErrMismatchedHashAndPassword
}

func (svc mockAuthService) GenerateJWT(id uint) (string, time.Time, error) {
	return "jwt-token", time.Now().Add(24 * time.Hour), nil
}

func (svc mockAuthService) ParseJWT(tokenString string) (*jwt.RegisteredClaims, error) {
	panic("not implemented") // TODO: Implement
}

func (svc mockAuthService) GenMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Mock middleware that sets a user ID in the context
		c.Locals("userID", "1")
		return c.Next()
	}
}

type usersTestSuite struct {
	suite.Suite
	app *fiber.App
}

func (suite *usersTestSuite) SetupSuite() {
	utils.SetDefaultLogger(slog.LevelDebug)

	suite.app = fiber.New(fiber.Config{ErrorHandler: api.ErrorHandler})
	users.RegisterRoutes(suite.app, users.Controller{UserService: mockUserService{}, AuthService: mockAuthService{}})
}

func (suite *usersTestSuite) TestRegister() {
	type testCaseOutput struct {
		status int
		body   users.RegisterResponse
	}

	type testCase struct {
		input  users.RegisterRequest
		output testCaseOutput
	}

	testCases := map[string]testCase{
		"succesful": {
			input: users.RegisterRequest{
				Name:     "Kshitish Deshpande",
				Email:    "me@ksdfg.dev",
				Password: "securepassword",
			},
			output: testCaseOutput{
				status: http.StatusCreated,
				body: users.RegisterResponse{
					ApiResponse: utils.ApiResponse{
						Success: true,
						Message: "User created successfully",
					},
					User: models.User{
						Name:  "Kshitish Deshpande",
						Email: "me@ksdfg.dev",
					},
				},
			},
		},
		"duplicate email": {
			input: users.RegisterRequest{
				Name:     "Kshitish Deshpande",
				Email:    "duplicate@ksdfg.dev",
				Password: "securepassword",
			},
			output: testCaseOutput{
				status: http.StatusConflict,
				body: users.RegisterResponse{
					ApiResponse: utils.ApiResponse{
						Success: false,
						Message: "User already exists",
					},
				},
			},
		},
	}

	for name, tc := range testCases {
		suite.Run(name, func() {
			// Marshal test user to json []byte body
			requestBody, err := json.Marshal(tc.input)
			if err != nil {
				suite.T().Error(err)
				return
			}

			// Generate POST request to /api/v1/users endpoint with body generated above
			request, err := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
			if err != nil {
				suite.T().Error(err)
				return
			}

			// Add content type header so that the app can parse the body
			request.Header.Add("Content-Type", "application/json")

			// Send the request
			response, err := suite.app.Test(request)
			if err != nil {
				suite.T().Error(err)
				return
			}
			defer response.Body.Close()

			// Assert that the response status code is as expected
			suite.Equal(tc.output.status, response.StatusCode)

			// Read the response body
			body, err := io.ReadAll(response.Body)
			if err != nil {
				suite.T().Error(err)
				return
			}

			// Unmarshal the response body into the expected response type
			var responseBody users.RegisterResponse
			if err = json.Unmarshal(body, &responseBody); err != nil {
				suite.T().Error(err)
				return
			}

			// Assert that the response body matches the expected output
			suite.Equal(tc.output.body.Success, responseBody.Success)
			suite.Equal(tc.output.body.Message, responseBody.Message)
			suite.Equal(tc.output.body.User.Name, responseBody.User.Name)
			suite.Equal(tc.output.body.User.Email, responseBody.User.Email)
			suite.Empty(responseBody.User.Password) // Password should not be returned in the response
		})
	}
}

func (suite *usersTestSuite) TestLogin() {
	type testCaseOutput struct {
		status int
		body   utils.ApiResponse
		userId string
	}

	type testCase struct {
		input  users.LoginRequest
		output testCaseOutput
	}

	testCases := map[string]testCase{
		"successful": {
			input: users.LoginRequest{
				Email:    "me@ksdfg.dev",
				Password: "securepassword",
			},
			output: testCaseOutput{
				status: http.StatusOK,
				body: utils.ApiResponse{
					Success: true,
					Message: "User logged in successfully",
				},
				userId: "1",
			},
		},
		"user not found": {
			input: users.LoginRequest{
				Email:    "nosuchuser@ksdfg.dev",
				Password: "securepassword",
			},
			output: testCaseOutput{
				status: http.StatusNotFound,
				body: utils.ApiResponse{
					Success: false,
					Message: "User not found",
				},
			},
		},
		"incorrect password": {
			input: users.LoginRequest{
				Email:    "me@ksdfg.dev",
				Password: "wrongpassword",
			},
			output: testCaseOutput{
				status: http.StatusUnauthorized,
				body: utils.ApiResponse{
					Success: false,
					Message: "Incorrect password",
				},
			},
		},
	}

	for name, tc := range testCases {
		suite.Run(name, func() {
			// Marshal test user to json []byte body
			requestBody, err := json.Marshal(tc.input)
			if err != nil {
				suite.T().Error(err)
				return
			}

			// Generate POST request to /api/v1/users endpoint with body generated above
			request, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(requestBody))
			if err != nil {
				suite.T().Error(err)
				return
			}

			// Add content type header so that the app can parse the body
			request.Header.Add("Content-Type", "application/json")

			// Send the request
			response, err := suite.app.Test(request)
			if err != nil {
				suite.T().Error(err)
				return
			}
			defer response.Body.Close()

			// Assert that the response status code is as expected
			suite.Equal(tc.output.status, response.StatusCode)

			// Read the response body
			body, err := io.ReadAll(response.Body)
			if err != nil {
				suite.T().Error(err)
				return
			}

			// Unmarshal the response body into the expected response type
			var responseBody utils.ApiResponse
			if err = json.Unmarshal(body, &responseBody); err != nil {
				suite.T().Error(err)
				return
			}

			// Assert that the response body matches the expected output
			suite.Equal(tc.output.body.Success, responseBody.Success)
			suite.Equal(tc.output.body.Message, responseBody.Message)

			// Search for the authorization cookie in the response
			foundAuthCookie := false
			for _, cookie := range response.Cookies() {
				if cookie.Name == "authorization" {
					foundAuthCookie = true
					break
				}
			}

			if tc.output.userId != "" && !foundAuthCookie {
				// If a user ID is specified in the expected output, then the auth cookie must be present in the response.
				suite.T().Error("Authorization cookie not found in response")
				return
			} else if tc.output.userId == "" && foundAuthCookie {
				// If a user ID is not specified in the expected output, then the auth cookie must not be present in the response.
				suite.T().Error("Authorization cookie found in response, but it should not be present")
				return
			}
		})
	}
}

func TestUsersRoutes(t *testing.T) {
	suite.Run(t, new(usersTestSuite))
}
