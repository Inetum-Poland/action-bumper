// Copyright (c) 2024 Inetum Poland.

// Package config provides configuration management for the action-bumper GitHub Action.
// It handles loading configuration from environment variables, including GitHub Actions
// inputs (INPUT_*) and built-in GitHub environment variables (GITHUB_*).
//
// Environment Variables:
//
// Required:
//   - INPUT_GITHUB_TOKEN or GITHUB_TOKEN: GitHub API authentication token
//   - GITHUB_EVENT_PATH: Path to the JSON file containing event payload
//   - GITHUB_EVENT_NAME: Name of the event that triggered the action
//   - GITHUB_REPOSITORY: Repository in "owner/repo" format
//
// Optional (with defaults):
//   - INPUT_BUMP_DEFAULT_LEVEL: Default bump level when no label is present (default: "")
//   - INPUT_BUMP_FAIL_IF_NO_LEVEL: Exit with error if no bump level (default: false)
//   - INPUT_BUMP_INCLUDE_V: Include 'v' prefix in tags (default: true)
//   - INPUT_BUMP_SEMVER: Create semver tags like v1, v1.2 (default: false)
//   - INPUT_BUMP_LATEST: Create/update 'latest' tag (default: false)
//   - INPUT_BUMP_MAJOR/MINOR/PATCH/NONE: Custom label names
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// BumpLevel represents the type of version bump
type BumpLevel string

const (
	BumpLevelMajor BumpLevel = "major"
	BumpLevelMinor BumpLevel = "minor"
	BumpLevelPatch BumpLevel = "patch"
	BumpLevelNone  BumpLevel = "none"
	BumpLevelEmpty BumpLevel = ""
)

// Config holds all configuration for the action
type Config struct {
	// GitHub configuration
	GitHubToken     string
	GitHubEventPath string
	GitHubEventName string
	GitHubRepo      string // owner/repo format
	GitHubSHA       string
	Workspace       string

	// Bump configuration
	BumpDefaultLevel  BumpLevel
	BumpFailIfNoLevel bool
	BumpIncludeV      bool
	BumpLatest        bool
	BumpSemver        bool
	BumpTagAsEmail    string
	BumpTagAsUser     string

	// Label names
	LabelMajor string
	LabelMinor string
	LabelPatch string
	LabelNone  string

	// Debug flags
	Debug bool
	Trace bool

	// Debug event path for testing (reads tags from local files instead of API)
	DebugEventPath string
}

// LoadFromEnv loads configuration from environment variables
func LoadFromEnv() (*Config, error) {
	// Check for debug event path first and load .input.env if present
	debugEventPath := os.Getenv("DEBUG_GITHUB_EVENT_PATH")
	if debugEventPath != "" {
		if err := loadDebugEnvFile(debugEventPath); err != nil {
			// Non-fatal: just log if the file doesn't exist
			fmt.Fprintf(os.Stderr, "Warning: could not load debug env file: %v\n", err)
		}
	}

	// Determine event path - use debug path's data.json if set
	eventPath := os.Getenv("GITHUB_EVENT_PATH")
	if debugEventPath != "" {
		eventPath = filepath.Join(debugEventPath, "data.json")
	}

	cfg := &Config{
		// GitHub environment variables
		GitHubToken:     os.Getenv("INPUT_GITHUB_TOKEN"),
		GitHubEventPath: eventPath,
		GitHubEventName: os.Getenv("GITHUB_EVENT_NAME"),
		GitHubRepo:      os.Getenv("GITHUB_REPOSITORY"),
		GitHubSHA:       os.Getenv("GITHUB_SHA"),
		Workspace:       os.Getenv("GITHUB_WORKSPACE"),

		// Bump configuration with defaults
		BumpDefaultLevel:  BumpLevel(os.Getenv("INPUT_BUMP_DEFAULT_LEVEL")),
		BumpFailIfNoLevel: parseBool(os.Getenv("INPUT_BUMP_FAIL_IF_NO_LEVEL"), false),
		BumpIncludeV:      parseBool(os.Getenv("INPUT_BUMP_INCLUDE_V"), true),
		BumpLatest:        parseBool(os.Getenv("INPUT_BUMP_LATEST"), false),
		BumpSemver:        parseBool(os.Getenv("INPUT_BUMP_SEMVER"), false),
		BumpTagAsEmail:    os.Getenv("INPUT_BUMP_TAG_AS_EMAIL"),
		BumpTagAsUser:     os.Getenv("INPUT_BUMP_TAG_AS_USER"),

		// Label names with defaults
		LabelMajor: getEnvOrDefault("INPUT_BUMP_MAJOR", "bumper:major"),
		LabelMinor: getEnvOrDefault("INPUT_BUMP_MINOR", "bumper:minor"),
		LabelPatch: getEnvOrDefault("INPUT_BUMP_PATCH", "bumper:patch"),
		LabelNone:  getEnvOrDefault("INPUT_BUMP_NONE", "bumper:none"),

		// Debug flags
		Debug: parseBool(os.Getenv("INETUM_POLAND_ACTION_BUMPER_DEBUG"), false),
		Trace: parseBool(os.Getenv("INETUM_POLAND_ACTION_BUMPER_TRACE"), false),

		// Debug event path for testing
		DebugEventPath: os.Getenv("DEBUG_GITHUB_EVENT_PATH"),
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks that required fields are present
func (c *Config) Validate() error {
	if c.GitHubToken == "" {
		return fmt.Errorf("GITHUB_TOKEN is required")
	}
	if c.GitHubEventPath == "" {
		return fmt.Errorf("GITHUB_EVENT_PATH is required")
	}
	if c.GitHubEventName == "" {
		return fmt.Errorf("GITHUB_EVENT_NAME is required")
	}
	if c.GitHubRepo == "" {
		return fmt.Errorf("GITHUB_REPOSITORY is required")
	}

	// Validate bump default level if provided
	if c.BumpDefaultLevel != BumpLevelEmpty {
		if !c.BumpDefaultLevel.IsValid() {
			return fmt.Errorf("invalid bump_default_level: %s (must be major, minor, patch, or none)", c.BumpDefaultLevel)
		}
	}

	return nil
}

// IsValid checks if the bump level is valid
func (b BumpLevel) IsValid() bool {
	switch b {
	case BumpLevelMajor, BumpLevelMinor, BumpLevelPatch, BumpLevelNone:
		return true
	case BumpLevelEmpty:
		return false
	default:
		return false
	}
}

// parseBool parses a string to bool with a default value
func parseBool(s string, defaultVal bool) bool {
	if s == "" {
		return defaultVal
	}
	val, err := strconv.ParseBool(s)
	if err != nil {
		return defaultVal
	}
	return val
}

// getEnvOrDefault gets environment variable or returns default
func getEnvOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// loadDebugEnvFile loads environment variables from a .input.env file
// in the debug event path directory. This mimics the Bash behavior of
// sourcing the .input.env file for testing purposes.
func loadDebugEnvFile(debugEventPath string) error {
	envFile := filepath.Join(debugEventPath, ".input.env")
	data, err := os.ReadFile(envFile)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse KEY=value or KEY="value"
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove surrounding quotes if present
		if len(value) >= 2 && ((value[0] == '"' && value[len(value)-1] == '"') ||
			(value[0] == '\'' && value[len(value)-1] == '\'')) {
			value = value[1 : len(value)-1]
		}

		// Set the environment variable
		os.Setenv(key, value)
	}

	return nil
}
