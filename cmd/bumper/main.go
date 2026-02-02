// Copyright (c) 2024 Inetum Poland.

package main

import (
	"context"
	"log"
	"os"

	"github.com/Inetum-Poland/action-bumper/internal/bumper"
	"github.com/Inetum-Poland/action-bumper/internal/config"
	"github.com/Inetum-Poland/action-bumper/internal/git"
	"github.com/Inetum-Poland/action-bumper/internal/logger"
)

func main() {
	// Load configuration from environment
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Configure logging
	var appLogger logger.Logger
	if cfg.Debug {
		appLogger = logger.NewStandardLogger(os.Stdout, "[action-bumper] ", log.LstdFlags|log.Lshortfile)
		appLogger.Println("Debug mode enabled")
	} else {
		appLogger = logger.NewDefaultLogger()
	}

	// Change to workspace directory if specified
	if cfg.Workspace != "" {
		if err := os.Chdir(cfg.Workspace); err != nil {
			log.Fatalf("Failed to change to workspace directory: %v", err)
		}
		appLogger.Printf("Changed to workspace: %s", cfg.Workspace)

		// Configure git safe directory
		if err := git.ConfigureSafeDirectory(cfg.Workspace); err != nil {
			appLogger.Printf("Warning: failed to configure safe directory: %v", err)
		}
	}

	// Create bumper instance
	ctx := context.Background()
	b, err := bumper.New(ctx, cfg, appLogger)
	if err != nil {
		log.Fatalf("Failed to create bumper: %v", err)
	}

	// Run the bumper
	if err := b.Run(ctx); err != nil {
		log.Fatalf("Bumper failed: %v", err)
	}

	appLogger.Println("Bumper completed successfully")
}
