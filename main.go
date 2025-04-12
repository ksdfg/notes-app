package main

import (
	"fmt"
	"log"
	"log/slog"
	"notes-app/api"
	"notes-app/config"
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

	// Generate the app
	app := api.GenApp()

	// Start the server
	log.Fatalln(app.Listen(fmt.Sprintf(":%d", cfg.Port)))
}
