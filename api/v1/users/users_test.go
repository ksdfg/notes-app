package users_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"notes-app/api"
	"notes-app/api/v1/users"
	"notes-app/models"
	"notes-app/service"
	"notes-app/utils"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type MockUserService struct{}

func (svc MockUserService) HashPassword(password string) (string, error) {
	return service.UserService{}.HashPassword(password)
}

func (svc MockUserService) ComparePasswords(hashedPassword, password string) error {
	return service.UserService{}.ComparePasswords(hashedPassword, password)
}

func (svc MockUserService) Create(user *models.User, opts *service.DBOpts) error {
	if user.Email == "duplicate@ksdfg.dev" {
		return gorm.ErrDuplicatedKey
	}

	user.ID = 1
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	return nil
}

func (svc MockUserService) GetByEmail(email string, opts *service.DBOpts) (models.User, error) {
	if email == "nosuchuser@ksdfg.dev" {
		return models.User{}, gorm.ErrRecordNotFound
	}

	password, err := svc.HashPassword("securepassword")
	if err != nil {
		return models.User{}, err
	}

	return models.User{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:     "Kshitish Deshpande",
		Email:    "me@ksdfg.dev",
		Password: password,
	}, nil
}

func (svc MockUserService) GetByID(id uint, opts *service.DBOpts) (models.User, error) {
	//TODO implement me
	panic("implement me")
}

type UsersTestSuite struct {
	suite.Suite
	app *fiber.App
}

func (suite *UsersTestSuite) SetupSuite() {
	suite.app = fiber.New(fiber.Config{ErrorHandler: api.ErrorHandler})
	users.RegisterRoutes(suite.app, MockUserService{})
}

func (suite *UsersTestSuite) TestRegister() {
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

func (suite *UsersTestSuite) TestLogin() {
	type testCaseOutput struct {
		status int
		body   utils.ApiResponse
	}

	type testCase struct {
		input  users.LoginRequest
		output testCaseOutput
	}

	testCases := map[string]testCase{
		"successful": {
			input: users.LoginRequest{
				Email:    "me@ksdfg.dev",   // Add a valid email here
				Password: "securepassword", // Add a valid password here
			},
			output: testCaseOutput{
				status: http.StatusOK,
				body: utils.ApiResponse{
					Success: true,
					Message: "User logged in successfully",
				},
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
		})
	}
}

func TestUsersRoutes(t *testing.T) {
	suite.Run(t, new(UsersTestSuite))
}
