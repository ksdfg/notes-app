package utils

import (
	"log/slog"
	"os"
)

// SetDefaultLogger sets the default logger.
//
// This sets the default logger to write logs to the console with the given
// level. The logger will include the source of the log message.
func SetDefaultLogger(level slog.Level) {
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
