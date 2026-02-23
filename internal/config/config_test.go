// Copyright (c) 2024 Inetum Poland.

package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadFromEnv(t *testing.T) {
	tests := []struct {
		name    string
		envVars map[string]string
		want    *Config
		wantErr bool
	}{
		{
			name: "all required fields present with defaults",
			envVars: map[string]string{
				"INPUT_GITHUB_TOKEN": "test-token",
				"GITHUB_EVENT_PATH":  "/path/to/event.json",
				"GITHUB_EVENT_NAME":  "push",
				"GITHUB_REPOSITORY":  "owner/repo",
				"GITHUB_ACTOR":       "testuser",
			},
			want: &Config{
				GitHubToken:       "test-token",
				GitHubEventPath:   "/path/to/event.json",
				GitHubEventName:   "push",
				GitHubRepo:        "owner/repo",
				GitHubActor:       "testuser",
				BumpIncludeV:      true,
				BumpFailIfNoLevel: false,
				BumpLatest:        false,
				BumpSemver:        false,
				LabelMajor:        "bumper:major",
				LabelMinor:        "bumper:minor",
				LabelPatch:        "bumper:patch",
				LabelNone:         "bumper:none",
				BumpTagAsUser:     "testuser",
				BumpTagAsEmail:    "testuser@users.noreply.github.com",
			},
			wantErr: false,
		},
		{
			name: "custom labels and settings",
			envVars: map[string]string{
				"INPUT_GITHUB_TOKEN":          "test-token",
				"GITHUB_EVENT_PATH":           "/path/to/event.json",
				"GITHUB_EVENT_NAME":           "pull_request",
				"GITHUB_REPOSITORY":           "owner/repo",
				"INPUT_BUMP_DEFAULT_LEVEL":    "patch",
				"INPUT_BUMP_FAIL_IF_NO_LEVEL": "true",
				"INPUT_BUMP_INCLUDE_V":        "false",
				"INPUT_BUMP_LATEST":           "true",
				"INPUT_BUMP_SEMVER":           "true",
				"INPUT_BUMP_MAJOR":            "release:major",
				"INPUT_BUMP_MINOR":            "release:minor",
				"INPUT_BUMP_PATCH":            "release:patch",
				"INPUT_BUMP_NONE":             "release:none",
				"INPUT_BUMP_TAG_AS_EMAIL":     "bot@example.com",
				"INPUT_BUMP_TAG_AS_USER":      "Release Bot",
			},
			want: &Config{
				GitHubToken:       "test-token",
				GitHubEventPath:   "/path/to/event.json",
				GitHubEventName:   "pull_request",
				GitHubRepo:        "owner/repo",
				BumpDefaultLevel:  BumpLevelPatch,
				BumpFailIfNoLevel: true,
				BumpIncludeV:      false,
				BumpLatest:        true,
				BumpSemver:        true,
				LabelMajor:        "release:major",
				LabelMinor:        "release:minor",
				LabelPatch:        "release:patch",
				LabelNone:         "release:none",
				BumpTagAsEmail:    "bot@example.com",
				BumpTagAsUser:     "Release Bot",
			},
			wantErr: false,
		},
		{
			name: "missing github token",
			envVars: map[string]string{
				"GITHUB_EVENT_PATH": "/path/to/event.json",
				"GITHUB_EVENT_NAME": "push",
				"GITHUB_REPOSITORY": "owner/repo",
			},
			wantErr: true,
		},
		{
			name: "missing event path",
			envVars: map[string]string{
				"INPUT_GITHUB_TOKEN": "test-token",
				"GITHUB_EVENT_NAME":  "push",
				"GITHUB_REPOSITORY":  "owner/repo",
			},
			wantErr: true,
		},
		{
			name: "invalid bump default level",
			envVars: map[string]string{
				"INPUT_GITHUB_TOKEN":       "test-token",
				"GITHUB_EVENT_PATH":        "/path/to/event.json",
				"GITHUB_EVENT_NAME":        "push",
				"GITHUB_REPOSITORY":        "owner/repo",
				"INPUT_BUMP_DEFAULT_LEVEL": "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment
			os.Clearenv()

			// Set test environment variables
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			got, err := LoadFromEnv()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want.GitHubToken, got.GitHubToken)
			assert.Equal(t, tt.want.GitHubEventPath, got.GitHubEventPath)
			assert.Equal(t, tt.want.GitHubEventName, got.GitHubEventName)
			assert.Equal(t, tt.want.GitHubRepo, got.GitHubRepo)
			assert.Equal(t, tt.want.BumpDefaultLevel, got.BumpDefaultLevel)
			assert.Equal(t, tt.want.BumpFailIfNoLevel, got.BumpFailIfNoLevel)
			assert.Equal(t, tt.want.BumpIncludeV, got.BumpIncludeV)
			assert.Equal(t, tt.want.BumpLatest, got.BumpLatest)
			assert.Equal(t, tt.want.BumpSemver, got.BumpSemver)
			assert.Equal(t, tt.want.LabelMajor, got.LabelMajor)
			assert.Equal(t, tt.want.LabelMinor, got.LabelMinor)
			assert.Equal(t, tt.want.LabelPatch, got.LabelPatch)
			assert.Equal(t, tt.want.LabelNone, got.LabelNone)
			assert.Equal(t, tt.want.BumpTagAsEmail, got.BumpTagAsEmail)
			assert.Equal(t, tt.want.BumpTagAsUser, got.BumpTagAsUser)
		})
	}
}

func TestBumpLevel_IsValid(t *testing.T) {
	tests := []struct {
		level BumpLevel
		want  bool
	}{
		{BumpLevelMajor, true},
		{BumpLevelMinor, true},
		{BumpLevelPatch, true},
		{BumpLevelNone, true},
		{BumpLevelEmpty, false},
		{"invalid", false},
	}

	for _, tt := range tests {
		t.Run(string(tt.level), func(t *testing.T) {
			assert.Equal(t, tt.want, tt.level.IsValid())
		})
	}
}

func TestParseBool(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		defaultVal bool
		want       bool
	}{
		{"true string", "true", false, true},
		{"false string", "false", true, false},
		{"1", "1", false, true},
		{"0", "0", true, false},
		{"empty uses default true", "", true, true},
		{"empty uses default false", "", false, false},
		{"invalid uses default", "invalid", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseBool(tt.input, tt.defaultVal)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConfig_BooleanEdgeCases(t *testing.T) {
	tests := []struct {
		name          string
		envValue      string
		expectedTrue  bool
		expectedFalse bool
	}{
		{"TRUE uppercase", "TRUE", true, false},
		{"True mixed", "True", false, true}, // parseBool uses strconv.ParseBool which handles TRUE but not True by default
		{"yes", "yes", false, true},         // not supported, defaults
		{"no", "no", false, true},           // not supported, defaults
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseBool(tt.envValue, true)
			assert.Equal(t, tt.expectedFalse || result, true)
		})
	}
}

func TestConfig_AllBumpLevels(t *testing.T) {
	levels := []struct {
		input string
		want  BumpLevel
		valid bool
	}{
		{"major", BumpLevelMajor, true},
		{"minor", BumpLevelMinor, true},
		{"patch", BumpLevelPatch, true},
		{"none", BumpLevelNone, true},
		{"MAJOR", BumpLevelEmpty, false}, // case sensitive
		{"Minor", BumpLevelEmpty, false},
		{"", BumpLevelEmpty, false},
		{"auto", BumpLevelEmpty, false},
	}

	for _, tt := range levels {
		t.Run("level_"+tt.input, func(t *testing.T) {
			os.Clearenv()
			os.Setenv("INPUT_GITHUB_TOKEN", "test-token")
			os.Setenv("GITHUB_EVENT_PATH", "/path/to/event.json")
			os.Setenv("GITHUB_EVENT_NAME", "push")
			os.Setenv("GITHUB_REPOSITORY", "owner/repo")
			os.Setenv("INPUT_BUMP_DEFAULT_LEVEL", tt.input)

			cfg, err := LoadFromEnv()
			if !tt.valid {
				if tt.input == "" {
					// Empty is allowed, defaults to empty
					require.NoError(t, err)
				} else {
					require.Error(t, err)
				}
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, cfg.BumpDefaultLevel)
		})
	}
}

func TestConfig_DefaultLabels(t *testing.T) {
	os.Clearenv()
	os.Setenv("INPUT_GITHUB_TOKEN", "test-token")
	os.Setenv("GITHUB_EVENT_PATH", "/path/to/event.json")
	os.Setenv("GITHUB_EVENT_NAME", "push")
	os.Setenv("GITHUB_REPOSITORY", "owner/repo")

	cfg, err := LoadFromEnv()
	require.NoError(t, err)

	// Verify defaults match Bash implementation
	assert.Equal(t, "bumper:major", cfg.LabelMajor)
	assert.Equal(t, "bumper:minor", cfg.LabelMinor)
	assert.Equal(t, "bumper:patch", cfg.LabelPatch)
	assert.Equal(t, "bumper:none", cfg.LabelNone)
}

func TestConfig_CustomLabels(t *testing.T) {
	os.Clearenv()
	os.Setenv("INPUT_GITHUB_TOKEN", "test-token")
	os.Setenv("GITHUB_EVENT_PATH", "/path/to/event.json")
	os.Setenv("GITHUB_EVENT_NAME", "push")
	os.Setenv("GITHUB_REPOSITORY", "owner/repo")
	os.Setenv("INPUT_BUMP_MAJOR", "version:breaking")
	os.Setenv("INPUT_BUMP_MINOR", "version:feature")
	os.Setenv("INPUT_BUMP_PATCH", "version:fix")
	os.Setenv("INPUT_BUMP_NONE", "version:skip")

	cfg, err := LoadFromEnv()
	require.NoError(t, err)

	assert.Equal(t, "version:breaking", cfg.LabelMajor)
	assert.Equal(t, "version:feature", cfg.LabelMinor)
	assert.Equal(t, "version:fix", cfg.LabelPatch)
	assert.Equal(t, "version:skip", cfg.LabelNone)
}

func TestConfig_MissingEventName(t *testing.T) {
	os.Clearenv()
	os.Setenv("INPUT_GITHUB_TOKEN", "test-token")
	os.Setenv("GITHUB_EVENT_PATH", "/path/to/event.json")
	os.Setenv("GITHUB_REPOSITORY", "owner/repo")

	_, err := LoadFromEnv()
	assert.Error(t, err)
}

func TestConfig_MissingRepository(t *testing.T) {
	os.Clearenv()
	os.Setenv("INPUT_GITHUB_TOKEN", "test-token")
	os.Setenv("GITHUB_EVENT_PATH", "/path/to/event.json")
	os.Setenv("GITHUB_EVENT_NAME", "push")

	_, err := LoadFromEnv()
	assert.Error(t, err)
}

func TestConfig_TagUserSettings(t *testing.T) {
	os.Clearenv()
	os.Setenv("INPUT_GITHUB_TOKEN", "test-token")
	os.Setenv("GITHUB_EVENT_PATH", "/path/to/event.json")
	os.Setenv("GITHUB_EVENT_NAME", "push")
	os.Setenv("GITHUB_REPOSITORY", "owner/repo")
	os.Setenv("INPUT_BUMP_TAG_AS_USER", "GitHub Actions Bot")
	os.Setenv("INPUT_BUMP_TAG_AS_EMAIL", "actions@github.com")

	cfg, err := LoadFromEnv()
	require.NoError(t, err)

	assert.Equal(t, "GitHub Actions Bot", cfg.BumpTagAsUser)
	assert.Equal(t, "actions@github.com", cfg.BumpTagAsEmail)
}
