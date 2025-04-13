package users_test

import (
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"io"
	"net/http"
	"notes-app/api"
	"notes-app/api/v1/users"
	"notes-app/models"
	"notes-app/service"
	"notes-app/utils"
	"testing"
	"time"
)

type MockUserService struct{}

func (svc MockUserService) HashPassword(password string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (svc MockUserService) ComparePasswords(hashedPassword, password string) error {
	//TODO implement me
	panic("implement me")
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
	//TODO implement me
	panic("implement me")
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

func (suite *UsersTestSuite) TestRegister_Successful() {
	// Marshal test animal to json []bye body
	requestBody, err := json.Marshal(users.RegisterRequest{
		Name:     "Kshitish Deshpande",
		Email:    "me@ksdfg.dev",
		Password: "securepassword",
	})
	if err != nil {
		suite.T().Error(err)
		return
	}

	// Generate POST request to /api/v1/animal endpoint with body generated above
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

	// Assert that the response status code is 201
	suite.Equal(http.StatusCreated, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		suite.T().Error(err)
		return
	}

	var responseBody users.RegisterResponse
	if err = json.Unmarshal(body, &responseBody); err != nil {
		suite.T().Error(err)
		return
	}

	suite.Equal(true, responseBody.Success)
	suite.Equal("User created successfully", responseBody.Message)
	suite.Equal("Kshitish Deshpande", responseBody.User.Name)
	suite.Equal("me@ksdfg.dev", responseBody.User.Email)
	suite.Empty(responseBody.User.Password)
}

func (suite *UsersTestSuite) TestRegister_DuplicateEmail() {
	// Marshal test animal to json []bye body
	requestBody, err := json.Marshal(users.RegisterRequest{
		Name:     "Kshitish Deshpande",
		Email:    "duplicate@ksdfg.dev",
		Password: "securepassword",
	})
	if err != nil {
		suite.T().Error(err)
		return
	}

	// Generate POST request to /api/v1/animal endpoint with body generated above
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

	// Assert that the response status code is 201
	suite.Equal(http.StatusConflict, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		suite.T().Error(err)
		return
	}

	var responseBody utils.ApiResponse
	if err = json.Unmarshal(body, &responseBody); err != nil {
		suite.T().Error(err)
		return
	}

	suite.Equal(false, responseBody.Success)
	suite.Equal("User already exists", responseBody.Message)
}

func TestUsersRoutes(t *testing.T) {
	suite.Run(t, new(UsersTestSuite))
}
