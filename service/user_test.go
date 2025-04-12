package service_test

import (
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"notes-app/config"
	"notes-app/database"
	"notes-app/models"
	"notes-app/service"
	"notes-app/utils"
	"testing"
)

type UserServiceTestSuite struct {
	suite.Suite
	dbService   database.Service
	userService service.UserService
}

func (suite *UserServiceTestSuite) SetupSuite() {
	// Setup logger at debug level for easy visibility during tests
	utils.SetDefaultLogger(slog.LevelDebug)

	cfg := config.Get()

	// Connect to the database
	suite.dbService = database.Service{}
	suite.dbService.Connect(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.TestDBName, cfg.DBSSLMode)

	// Create the user service instance to use for testing
	suite.userService = service.UserService{Service: service.Service{DBService: suite.dbService}}

	slog.Debug("Setup suite")
}

func (suite *UserServiceTestSuite) SetupTest() {
	// Clear all tables before each test
	suite.dbService.ClearAllTables()
	slog.Debug("Setup test")
}

func (suite *UserServiceTestSuite) TestHashPassword() {
	password := "password"

	// Hash the password using the service
	hashedPassword, errHash := suite.userService.HashPassword(password)
	suite.NoError(errHash)
	suite.NotEmpty(hashedPassword)
	slog.Debug("Hashed password", slog.String("hash", hashedPassword))

	// Compare the hashed password
	errCompare := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	suite.NoError(errCompare)
}

func (suite *UserServiceTestSuite) TestComparePassword() {
	password := "password"

	// Hash the password
	hashedPassword, errHash := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	suite.NoError(errHash)
	suite.NotEmpty(hashedPassword)
	slog.Debug("Hashed password", slog.String("hash", string(hashedPassword)))

	// Compare the hashed password using the service
	errCompare := suite.userService.ComparePasswords(string(hashedPassword), password)
	suite.NoError(errCompare)
}

func (suite *UserServiceTestSuite) TestCreate() {
	svc := suite.userService

	password := "password"
	user := models.User{
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Password: password,
	}

	// Create the user using the service
	errCreate := svc.Create(&user, nil)
	suite.NoError(errCreate)
	slog.Debug("Created user", slog.Any("user", user))

	// Get the user from the database
	db := suite.dbService.GetDB()
	var userFromDB models.User
	errSearch := db.Where("id = ?", user.ID).First(&userFromDB).Error
	suite.NoError(errSearch)
	slog.Debug("Got user from DB", slog.Any("user", userFromDB))

	// Assert that the details of the created and retrieved user are the same
	suite.Equal(user.ID, userFromDB.ID)
	suite.Equal(user.Name, userFromDB.Name)
	suite.Equal(user.Email, userFromDB.Email)
	suite.Equal(user.Password, userFromDB.Password)

	// Compare the hashed password
	errCompare := bcrypt.CompareHashAndPassword([]byte(userFromDB.Password), []byte(password))
	suite.NoError(errCompare)
}

func (suite *UserServiceTestSuite) TestGetByID() {
	user := models.User{
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Password: "password",
	}

	// Create the user
	errCreate := suite.dbService.GetDB().Create(&user).Error
	suite.NoError(errCreate)
	slog.Debug("Created user", slog.Any("user", user))

	// Get the user from the database using the service
	svc := suite.userService
	userFromDB, errSearch := svc.GetByID(user.ID, nil)
	suite.NoError(errSearch)
	slog.Debug("Got user from DB", slog.Any("user", userFromDB))

	// Assert that the details of the created and retrieved user are the same
	suite.Equal(user.ID, userFromDB.ID)
	suite.Equal(user.Name, userFromDB.Name)
	suite.Equal(user.Email, userFromDB.Email)
	suite.Equal(user.Password, userFromDB.Password)
}

func (suite *UserServiceTestSuite) TestGetByEmail() {
	user := models.User{
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Password: "password",
	}

	// Create the user
	errCreate := suite.dbService.GetDB().Create(&user).Error
	suite.NoError(errCreate)
	slog.Debug("Created user", slog.Any("user", user))

	// Get the user from the database using the service
	svc := suite.userService
	userFromDB, errSearch := svc.GetByEmail("john.doe@example.com", nil)
	suite.NoError(errSearch)
	slog.Debug("Got user from DB", slog.Any("user", userFromDB))

	// Assert that the details of the created and retrieved user are the same
	suite.Equal(user.ID, userFromDB.ID)
	suite.Equal(user.Name, userFromDB.Name)
	suite.Equal(user.Email, userFromDB.Email)
	suite.Equal(user.Password, userFromDB.Password)
}

func TestUserService(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}
