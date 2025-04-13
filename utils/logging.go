package utils

import (
	"github.com/lmittmann/tint"
	"log/slog"
	"os"
)

// SetDefaultLogger sets the default logger.
//
// This sets the default logger to write colourised logs to the console with the given
// level. The logger will include the source of the log message.
func SetDefaultLogger(level slog.Level) {
	// Create the logger
	logger := slog.New(
		tint.NewHandler(os.Stdout, &tint.Options{AddSource: true, Level: level}),
	)

	// Set the default logger
	slog.SetDefault(logger)
}
