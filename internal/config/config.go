// Copyright (c) 2024 Inetum Poland.

package config

import (
	"fmt"
	"os"
	"strconv"
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
}

// LoadFromEnv loads configuration from environment variables
func LoadFromEnv() (*Config, error) {
	cfg := &Config{
		// GitHub environment variables
		GitHubToken:     os.Getenv("INPUT_GITHUB_TOKEN"),
		GitHubEventPath: os.Getenv("GITHUB_EVENT_PATH"),
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
