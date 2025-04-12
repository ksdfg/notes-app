package main

import (
	"fmt"
	"log"
	"log/slog"
	"notes-app/api"
	"notes-app/config"
	"os"
)

func init() {
	// Get the config object
	cfg := config.Get()

	// Set the log level
	var level slog.Level
	err := level.UnmarshalText([]byte(cfg.LogLevel))
	if err != nil {
		// Panic if the log level is invalid
		log.Fatalln(err)
	}

	// Create the logger
	logger := slog.New(
		// Use the text handler
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			// Show the source of the log message
			AddSource: true,
			// Set the log level
			Level: level,
			// Don't replace any attributes
			ReplaceAttr: nil,
		}),
	)

	// Set the default logger
	slog.SetDefault(logger)
}

func main() {
	// Get the config object
	cfg := config.Get()

	// Generate the app
	app := api.GenApp()

	// Start the server
	log.Fatalln(app.Listen(fmt.Sprintf(":%d", cfg.Port)))
}
