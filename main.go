package main

import (
	"fmt"
	"log"
	"log/slog"
	"notes-app/api"
	"notes-app/config"
	"notes-app/database"
	"notes-app/service"
	"notes-app/utils"
)

func init() {
	// Get the config object
	cfg := config.Get()

	// Get the log level from config
	var level slog.Level
	err := level.UnmarshalText([]byte(cfg.LogLevel))
	if err != nil {
		// Panic if the log level is invalid
		log.Fatalln(err)
	}

	// Set the default logger
	utils.SetDefaultLogger(level)
}

func main() {
	// Get the config object
	cfg := config.Get()

	// Connect to the database
	dbService := database.Service{}
	dbService.Connect(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)
	dbService.GetDB()

	// Initialize services
	authService := service.AuthService{}
	userService := service.UserService{Service: service.Service{DBService: dbService}, AuthService: authService}

	// Generate the app
	app := api.GenApp(api.Services{UserService: userService, AuthService: authService})

	// Start the server
	log.Fatalln(app.Listen(fmt.Sprintf(":%d", cfg.Port)))
}
