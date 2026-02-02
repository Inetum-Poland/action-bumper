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
			},
			want: &Config{
				GitHubToken:       "test-token",
				GitHubEventPath:   "/path/to/event.json",
				GitHubEventName:   "push",
				GitHubRepo:        "owner/repo",
				BumpIncludeV:      true,
				BumpFailIfNoLevel: false,
				BumpLatest:        false,
				BumpSemver:        false,
				LabelMajor:        "bumper:major",
				LabelMinor:        "bumper:minor",
				LabelPatch:        "bumper:patch",
				LabelNone:         "bumper:none",
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
