package main

import (
	"log/slog"
	"os"
)

func init() {
	logger := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug, ReplaceAttr: nil}),
	)
	slog.SetDefault(logger)
}

func main() {
	slog.Info("Hello world!")
}
