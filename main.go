package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"dndgoldtracker/ui"

	tea "charm.land/bubbletea/v2"
)

func main() {
	if err := run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func run() error {
	fileName := "logFile.log"

	// open log file
	logFile, err := os.OpenFile(fileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return fmt.Errorf("opening log file: %w", err)
	}
	defer func() {
		if err := logFile.Close(); err != nil {
			fmt.Println("Error:", fmt.Errorf("closing log file: %w", err))
		}
	}()

	// set logging level from ENV (if available)
	level := parseLogLevel(os.Getenv("DNDGOLD_LOGGING"))
	handler := slog.NewTextHandler(logFile, &slog.HandlerOptions{Level: level})
	slog.SetDefault(slog.New(handler))

	p := tea.NewProgram(ui.NewModel())

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("application runtime: %w", err)
	}

	return nil
}

// parseLogLevel converts a string into a log level
func parseLogLevel(s string) slog.Level {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
