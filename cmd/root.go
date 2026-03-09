package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"dndgoldtracker/internal/ui"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"
)

// GetRootCmd returns the root command, initializing it if necessary.
func GetRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "dndgoldtracker",
		Version: "0.2.0",
		Short:   "A basic tracker for gold and experience",
		RunE:    run,
	}

	// Add persistent flags for configuration
	cmd.Flags().StringP("party-file", "f", "party.json", "file used to store party data")

	return cmd
}

// Execute runs the root command using fang for enhanced styling and error handling.
// This is called by main.main() and is the application entry point.
func Execute(ctx context.Context) error {
	return fang.Execute(ctx, GetRootCmd(), fang.WithoutVersion())
}

func run(cmd *cobra.Command, _ []string) error {
	fileName := "logFile.log"

	// open log file
	logFile, err := os.OpenFile(fileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return fmt.Errorf("opening log file: %w", err)
	}
	defer func() {
		if logErr := logFile.Close(); logErr != nil {
			fmt.Println("Error:", fmt.Errorf("closing log file: %w", logErr))
		}
	}()

	// set logging level from ENV (if available)
	level := parseLogLevel(os.Getenv("DNDGOLD_LOGGING"))
	handler := slog.NewTextHandler(logFile, &slog.HandlerOptions{Level: level})
	slog.SetDefault(slog.New(handler))

	f, err := cmd.Flags().GetString("party-file")
	if err != nil {
		return err
	}

	pf, err := filepath.Abs(f)
	if err != nil {
		return err
	}

	p := tea.NewProgram(ui.NewModel(pf))

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
