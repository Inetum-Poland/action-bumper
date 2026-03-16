// Copyright (c) 2024-2026 Inetum Poland.

package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/Inetum-Poland/action-bumper/internal/bumper"
	"github.com/Inetum-Poland/action-bumper/internal/config"
	"github.com/Inetum-Poland/action-bumper/internal/git"
	"github.com/Inetum-Poland/action-bumper/internal/preflight"
)

func main() {
	// Load configuration from environment
	cfg, err := config.LoadFromEnv()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	// Configure logging
	var appLogger *slog.Logger
	switch {
	case cfg.Trace:
		// Trace mode: most verbose, includes everything
		handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug - 4, // Even more verbose than debug
		})
		appLogger = slog.New(handler)
		appLogger.Info("Trace mode enabled")
	case cfg.Debug:
		handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		})
		appLogger = slog.New(handler)
		appLogger.Info("Debug mode enabled")
	default:
		appLogger = slog.Default()
	}

	// Run pre-flight checks
	ctx := context.Background()
	checker := preflight.NewChecker()
	results := checker.CheckAll(ctx)
	for _, r := range results {
		if r.Passed {
			appLogger.Debug("Pre-flight check passed", "check", r.Name, "message", r.Message)
		} else {
			appLogger.Warn("Pre-flight check failed", "check", r.Name, "message", r.Message, "error", r.Error)
		}
	}
	// Only git check is mandatory; GitHub reachability is informational
	gitPassed := false
	for _, r := range results {
		if r.Name == "git-available" {
			gitPassed = r.Passed
			break
		}
	}
	if !gitPassed {
		appLogger.Error("Required pre-flight check failed", "check", "git-available")
		os.Exit(1)
	}

	// Change to workspace directory if specified
	if cfg.Workspace != "" {
		if chErr := os.Chdir(cfg.Workspace); chErr != nil {
			appLogger.Error("Failed to change to workspace directory", "error", chErr)
			os.Exit(1)
		}
		appLogger.Info("Changed to workspace", "path", cfg.Workspace)

		// Configure git safe directory
		if safeErr := git.ConfigureSafeDirectory(cfg.Workspace); safeErr != nil {
			appLogger.Warn("Failed to configure safe directory", "error", safeErr)
		}
	}

	// Create bumper instance
	b, err := bumper.New(ctx, cfg, appLogger)
	if err != nil {
		appLogger.Error("Failed to create bumper", "error", err)
		os.Exit(1)
	}

	// Run the bumper
	if err := b.Run(ctx); err != nil {
		appLogger.Error("Bumper failed", "error", err)
		os.Exit(1)
	}

	appLogger.Info("Bumper completed successfully")
}
