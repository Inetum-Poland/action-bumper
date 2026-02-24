// Copyright (c) 2024 Inetum Poland.

package bumper

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Inetum-Poland/action-bumper/internal/config"
	"github.com/Inetum-Poland/action-bumper/internal/git"
	"github.com/Inetum-Poland/action-bumper/internal/github"
	"github.com/Inetum-Poland/action-bumper/internal/semver"
)

func TestDetermineBumpLevel(t *testing.T) {
	cfg := &config.Config{
		LabelMajor: "bumper:major",
		LabelMinor: "bumper:minor",
		LabelPatch: "bumper:patch",
		LabelNone:  "bumper:none",
	}
	b := &Bumper{cfg: cfg}

	tests := []struct {
		name   string
		labels []github.Label
		want   config.BumpLevel
	}{
		{
			name:   "major label",
			labels: []github.Label{{Name: "bumper:major"}},
			want:   config.BumpLevelMajor,
		},
		{
			name:   "minor label",
			labels: []github.Label{{Name: "bumper:minor"}},
			want:   config.BumpLevelMinor,
		},
		{
			name:   "patch label",
			labels: []github.Label{{Name: "bumper:patch"}},
			want:   config.BumpLevelPatch,
		},
		{
			name:   "none label",
			labels: []github.Label{{Name: "bumper:none"}},
			want:   config.BumpLevelNone,
		},
		{
			name:   "no matching label",
			labels: []github.Label{{Name: "feature"}, {Name: "bug"}},
			want:   config.BumpLevelEmpty,
		},
		{
			name:   "empty labels",
			labels: []github.Label{},
			want:   config.BumpLevelEmpty,
		},
		{
			name:   "multiple labels - major wins over minor",
			labels: []github.Label{{Name: "bumper:minor"}, {Name: "bumper:major"}},
			want:   config.BumpLevelMajor, // Major takes priority (matches Bash behavior)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := b.determineBumpLevel(tt.labels)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDetermineBumpLevelCustomLabels(t *testing.T) {
	cfg := &config.Config{
		LabelMajor: "release:major",
		LabelMinor: "release:minor",
		LabelPatch: "release:patch",
		LabelNone:  "release:skip",
	}
	b := &Bumper{cfg: cfg}

	tests := []struct {
		name   string
		labels []github.Label
		want   config.BumpLevel
	}{
		{
			name:   "custom major label",
			labels: []github.Label{{Name: "release:major"}},
			want:   config.BumpLevelMajor,
		},
		{
			name:   "custom minor label",
			labels: []github.Label{{Name: "release:minor"}},
			want:   config.BumpLevelMinor,
		},
		{
			name:   "custom patch label",
			labels: []github.Label{{Name: "release:patch"}},
			want:   config.BumpLevelPatch,
		},
		{
			name:   "custom none label",
			labels: []github.Label{{Name: "release:skip"}},
			want:   config.BumpLevelNone,
		},
		{
			name:   "default label not recognized with custom config",
			labels: []github.Label{{Name: "bumper:major"}},
			want:   config.BumpLevelEmpty,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := b.determineBumpLevel(tt.labels)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNew_ValidConfig(t *testing.T) {
	cfg := &config.Config{
		GitHubToken:     "test-token",
		GitHubRepo:      "owner/repo",
		GitHubEventName: "pull_request",
		GitHubEventPath: "/path/to/event.json",
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	b, err := New(context.Background(), cfg, logger)
	require.NoError(t, err)
	assert.NotNil(t, b)
	assert.Equal(t, cfg, b.cfg)
}

func TestDetermineBumpLevel_LabelPriority(t *testing.T) {
	// Verify that higher priority labels win (major > minor > patch > none)
	cfg := &config.Config{
		LabelMajor: "bumper:major",
		LabelMinor: "bumper:minor",
		LabelPatch: "bumper:patch",
		LabelNone:  "bumper:none",
	}
	b := &Bumper{cfg: cfg}

	// Order matters - Bash behavior where checks are done in order: none → patch → minor → major
	labels := []github.Label{
		{Name: "other-label"},
		{Name: "bumper:patch"}, // Lower priority
		{Name: "bumper:major"}, // Higher priority - wins
	}

	got := b.determineBumpLevel(labels)
	assert.Equal(t, config.BumpLevelMajor, got)
}

func TestGeneratePRStatusMessage(t *testing.T) {
	cfg := &config.Config{
		GitHubRepo: "owner/repo",
	}
	b := &Bumper{cfg: cfg}

	tests := []struct {
		name            string
		currentVersion  string
		nextVersion     string
		level           config.BumpLevel
		containsVersion bool
		containsLevel   bool
	}{
		{
			name:            "standard PR message",
			currentVersion:  "1.0.0",
			nextVersion:     "1.1.0",
			level:           config.BumpLevelMinor,
			containsVersion: true,
			containsLevel:   true,
		},
		{
			name:            "first release",
			currentVersion:  "0.0.0",
			nextVersion:     "1.0.0",
			level:           config.BumpLevelMajor,
			containsVersion: true,
			containsLevel:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			current, _ := semver.Parse(tt.currentVersion)
			next, _ := semver.Parse(tt.nextVersion)
			msg := b.generatePRStatusMessage(current, next, tt.level)
			if tt.containsVersion {
				assert.Contains(t, msg, tt.nextVersion)
			}
		})
	}
}

func TestDetermineBumpLevel_EmptyLabels(t *testing.T) {
	cfg := &config.Config{
		LabelMajor: "bumper:major",
		LabelMinor: "bumper:minor",
		LabelPatch: "bumper:patch",
		LabelNone:  "bumper:none",
	}
	b := &Bumper{cfg: cfg}

	// nil labels
	got := b.determineBumpLevel(nil)
	assert.Equal(t, config.BumpLevelEmpty, got)

	// empty slice
	got = b.determineBumpLevel([]github.Label{})
	assert.Equal(t, config.BumpLevelEmpty, got)
}

func TestBumperConfig(t *testing.T) {
	// Test that Bumper correctly stores and uses config
	cfg := &config.Config{
		GitHubToken:       "token123",
		GitHubRepo:        "owner/repo",
		GitHubEventName:   "push",
		GitHubEventPath:   "/event.json",
		BumpDefaultLevel:  config.BumpLevelPatch,
		BumpFailIfNoLevel: true,
		BumpIncludeV:      true,
		BumpLatest:        true,
		BumpSemver:        true,
		LabelMajor:        "major",
		LabelMinor:        "minor",
		LabelPatch:        "patch",
		LabelNone:         "none",
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	b, err := New(context.Background(), cfg, logger)
	require.NoError(t, err)

	// Verify all config fields are accessible
	assert.Equal(t, "token123", b.cfg.GitHubToken)
	assert.Equal(t, "owner/repo", b.cfg.GitHubRepo)
	assert.Equal(t, config.BumpLevelPatch, b.cfg.BumpDefaultLevel)
	assert.True(t, b.cfg.BumpFailIfNoLevel)
	assert.True(t, b.cfg.BumpIncludeV)
	assert.True(t, b.cfg.BumpLatest)
	assert.True(t, b.cfg.BumpSemver)
}

func TestRun_InvalidEventPath(t *testing.T) {
	cfg := &config.Config{
		GitHubToken:     "test-token",
		GitHubRepo:      "owner/repo",
		GitHubEventName: "push",
		GitHubEventPath: "/nonexistent/path/event.json",
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	b, err := New(context.Background(), cfg, logger)
	require.NoError(t, err)

	err = b.Run(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse event")
}

func TestRun_UnknownEventType(t *testing.T) {
	// Create a temp event file with neither PR nor push event structure
	tmpDir := t.TempDir()
	eventFile := tmpDir + "/event.json"
	err := os.WriteFile(eventFile, []byte(`{"action":"", "after":"", "pull_request":null}`), 0o644)
	require.NoError(t, err)

	cfg := &config.Config{
		GitHubToken:     "test-token",
		GitHubRepo:      "owner/repo",
		GitHubEventName: "unknown",
		GitHubEventPath: eventFile,
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	b, err := New(context.Background(), cfg, logger)
	require.NoError(t, err)

	err = b.Run(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown event type")
}

func TestNew_InvalidRepoFormat(t *testing.T) {
	cfg := &config.Config{
		GitHubToken:     "test-token",
		GitHubRepo:      "invalid-repo-format",
		GitHubEventName: "push",
		GitHubEventPath: "/event.json",
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	b, err := New(context.Background(), cfg, logger)
	assert.Error(t, err)
	assert.Nil(t, b)
	assert.Contains(t, err.Error(), "invalid GITHUB_REPOSITORY format")
}

func TestGeneratePushStatusMessage(t *testing.T) {
	tests := []struct {
		name         string
		cfg          *config.Config
		version      string
		tags         []string
		prNumber     int
		prTitle      string
		wantContains []string
	}{
		{
			name: "basic push message",
			cfg: &config.Config{
				GitHubRepo:   "owner/repo",
				BumpIncludeV: true,
				BumpSemver:   false,
				BumpLatest:   false,
			},
			version:  "1.2.3",
			tags:     []string{"v1.2.3"},
			prNumber: 42,
			prTitle:  "Add feature",
			wantContains: []string{
				"v1.2.3",
				"owner/repo",
				"Bumped!",
			},
		},
		{
			name: "push message with semver tags",
			cfg: &config.Config{
				GitHubRepo:   "owner/repo",
				BumpIncludeV: true,
				BumpSemver:   true,
				BumpLatest:   false,
			},
			version:  "1.2.3",
			tags:     []string{"v1.2.3", "v1.2", "v1"},
			prNumber: 42,
			prTitle:  "Add feature",
			wantContains: []string{
				"v1.2.3",
				"v1.2",
				"v1",
			},
		},
		{
			name: "push message with latest tag",
			cfg: &config.Config{
				GitHubRepo:   "owner/repo",
				BumpIncludeV: true,
				BumpSemver:   false,
				BumpLatest:   true,
			},
			version:  "1.0.0",
			tags:     []string{"v1.0.0", "latest"},
			prNumber: 1,
			prTitle:  "Initial release",
			wantContains: []string{
				"v1.0.0",
				"latest",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bumper{cfg: tt.cfg}
			v, _ := semver.Parse(tt.version)
			msg := b.generatePushStatusMessage(v, tt.tags, tt.prNumber, tt.prTitle)
			for _, want := range tt.wantContains {
				assert.Contains(t, msg, want)
			}
		})
	}
}

func TestNewWithClient(t *testing.T) {
	cfg := &config.Config{
		GitHubToken:     "test-token",
		GitHubRepo:      "owner/repo",
		GitHubEventName: "push",
	}
	mockClient := github.NewMockClient()
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	b := NewWithClient(cfg, mockClient, logger)

	assert.NotNil(t, b)
	assert.Equal(t, cfg, b.cfg)
	assert.Equal(t, mockClient, b.client)
}

func TestHandlePushEvent_WithValidMergedPR(t *testing.T) {
	// Create temp event file with push event structure
	tmpDir := t.TempDir()
	eventFile := tmpDir + "/event.json"
	eventJSON := `{
		"action": "",
		"after": "abc123def",
		"head_commit": {"id": "abc123def"},
		"pull_request": null
	}`
	err := os.WriteFile(eventFile, []byte(eventJSON), 0o644)
	require.NoError(t, err)

	cfg := &config.Config{
		GitHubToken:     "test-token",
		GitHubRepo:      "owner/repo",
		GitHubEventName: "push",
		GitHubEventPath: eventFile,
		GitHubSHA:       "abc123def",
		LabelMajor:      "bumper:major",
		LabelMinor:      "bumper:minor",
		LabelPatch:      "bumper:patch",
		LabelNone:       "bumper:none",
		BumpIncludeV:    true,
	}

	currentVersion, _ := semver.Parse("1.0.0")
	mockClient := github.NewMockClient()
	mockClient.GetLatestTagFunc = func(_ context.Context) (*semver.Version, error) {
		return currentVersion, nil
	}
	mockClient.GetMergedPRByCommitSHAFunc = func(_ context.Context, _ string) (*github.PullRequest, error) {
		return &github.PullRequest{
			Number: 42,
			Title:  "Feature: Add new capability",
			Labels: []github.Label{{Name: "bumper:minor"}},
		}, nil
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	mockGit := git.NewMockOperator()
	b := NewWithClientAndGit(cfg, mockClient, mockGit, logger)

	// Run with mock git - no real git commands will be executed
	_ = b.Run(context.Background())
	// Verify the mock was called
	assert.Contains(t, mockClient.Calls, "GetMergedPRByCommitSHA")
}

func TestHandlePushEvent_NoHeadCommit(t *testing.T) {
	// Create temp event file without head_commit
	tmpDir := t.TempDir()
	eventFile := tmpDir + "/event.json"
	eventJSON := `{
		"action": "",
		"after": "abc123def",
		"head_commit": null,
		"pull_request": null
	}`
	err := os.WriteFile(eventFile, []byte(eventJSON), 0o644)
	require.NoError(t, err)

	cfg := &config.Config{
		GitHubToken:     "test-token",
		GitHubRepo:      "owner/repo",
		GitHubEventName: "push",
		GitHubEventPath: eventFile,
	}

	mockClient := github.NewMockClient()
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	b := NewWithClient(cfg, mockClient, logger)

	err = b.Run(context.Background())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no head commit in push event")
}

func TestHandlePushEvent_NoBumpLabel_FailIfNoLevel(t *testing.T) {
	// Create temp event file
	tmpDir := t.TempDir()
	eventFile := tmpDir + "/event.json"
	eventJSON := `{
		"action": "",
		"after": "abc123def",
		"head_commit": {"id": "abc123def"},
		"pull_request": null
	}`
	err := os.WriteFile(eventFile, []byte(eventJSON), 0o644)
	require.NoError(t, err)

	cfg := &config.Config{
		GitHubToken:       "test-token",
		GitHubRepo:        "owner/repo",
		GitHubEventName:   "push",
		GitHubEventPath:   eventFile,
		GitHubSHA:         "abc123def",
		LabelMajor:        "bumper:major",
		LabelMinor:        "bumper:minor",
		LabelPatch:        "bumper:patch",
		LabelNone:         "bumper:none",
		BumpFailIfNoLevel: true,
	}

	mockClient := github.NewMockClient()
	mockClient.GetMergedPRByCommitSHAFunc = func(_ context.Context, _ string) (*github.PullRequest, error) {
		return &github.PullRequest{
			Number: 42,
			Title:  "Feature without bump label",
			Labels: []github.Label{{Name: "feature"}}, // No bump label
		}, nil
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	b := NewWithClient(cfg, mockClient, logger)

	err = b.Run(context.Background())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no bump level label found")
}

func TestHandlePushEvent_BumperNoneLabel(t *testing.T) {
	// Create temp event file
	tmpDir := t.TempDir()
	eventFile := tmpDir + "/event.json"
	eventJSON := `{
		"action": "",
		"after": "abc123def",
		"head_commit": {"id": "abc123def"},
		"pull_request": null
	}`
	err := os.WriteFile(eventFile, []byte(eventJSON), 0o644)
	require.NoError(t, err)

	cfg := &config.Config{
		GitHubToken:     "test-token",
		GitHubRepo:      "owner/repo",
		GitHubEventName: "push",
		GitHubEventPath: eventFile,
		GitHubSHA:       "abc123def",
		LabelMajor:      "bumper:major",
		LabelMinor:      "bumper:minor",
		LabelPatch:      "bumper:patch",
		LabelNone:       "bumper:none",
	}

	mockClient := github.NewMockClient()
	mockClient.GetMergedPRByCommitSHAFunc = func(_ context.Context, _ string) (*github.PullRequest, error) {
		return &github.PullRequest{
			Number: 42,
			Title:  "Skip version bump",
			Labels: []github.Label{{Name: "bumper:none"}},
		}, nil
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	b := NewWithClient(cfg, mockClient, logger)

	err = b.Run(context.Background())

	// Should succeed without error (skip scenario)
	assert.NoError(t, err)
}

func TestHandlePREvent_WithMockClient(t *testing.T) {
	// Create temp event file with PR event structure
	tmpDir := t.TempDir()
	eventFile := tmpDir + "/event.json"
	eventJSON := `{
		"action": "opened",
		"pull_request": {
			"number": 123,
			"title": "Test PR",
			"labels": [{"name": "bumper:minor"}]
		}
	}`
	err := os.WriteFile(eventFile, []byte(eventJSON), 0o644)
	require.NoError(t, err)

	cfg := &config.Config{
		GitHubToken:     "test-token",
		GitHubRepo:      "owner/repo",
		GitHubEventName: "pull_request",
		GitHubEventPath: eventFile,
		LabelMajor:      "bumper:major",
		LabelMinor:      "bumper:minor",
		LabelPatch:      "bumper:patch",
		LabelNone:       "bumper:none",
		BumpIncludeV:    true,
	}

	currentVersion, _ := semver.Parse("2.3.4")
	mockClient := github.NewMockClient()
	mockClient.GetLatestTagFunc = func(_ context.Context) (*semver.Version, error) {
		return currentVersion, nil
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	b := NewWithClient(cfg, mockClient, logger)

	err = b.Run(context.Background())

	assert.NoError(t, err)
	// Verify GetLatestTag was called
	assert.Contains(t, mockClient.Calls, "GetLatestTag")
}

func TestHandlePREvent_NoTags(t *testing.T) {
	// Create temp event file with PR event structure
	tmpDir := t.TempDir()
	eventFile := tmpDir + "/event.json"
	eventJSON := `{
		"action": "opened",
		"pull_request": {
			"number": 1,
			"title": "First PR",
			"labels": [{"name": "bumper:major"}]
		}
	}`
	err := os.WriteFile(eventFile, []byte(eventJSON), 0o644)
	require.NoError(t, err)

	cfg := &config.Config{
		GitHubToken:     "test-token",
		GitHubRepo:      "owner/repo",
		GitHubEventName: "pull_request",
		GitHubEventPath: eventFile,
		LabelMajor:      "bumper:major",
		LabelMinor:      "bumper:minor",
		LabelPatch:      "bumper:patch",
		LabelNone:       "bumper:none",
		BumpIncludeV:    true,
	}

	mockClient := github.NewMockClient()
	mockClient.GetLatestTagFunc = func(_ context.Context) (*semver.Version, error) {
		return nil, errors.New("no tags found")
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	b := NewWithClient(cfg, mockClient, logger)

	err = b.Run(context.Background())

	// Should succeed - uses default version when no tags exist
	assert.NoError(t, err)
}

func TestHandlePREvent_SkipOnNoneLabel(t *testing.T) {
	// Create temp event file
	tmpDir := t.TempDir()
	eventFile := tmpDir + "/event.json"
	eventJSON := `{
		"action": "labeled",
		"pull_request": {
			"number": 42,
			"title": "Documentation update",
			"labels": [{"name": "bumper:none"}]
		}
	}`
	err := os.WriteFile(eventFile, []byte(eventJSON), 0o644)
	require.NoError(t, err)

	cfg := &config.Config{
		GitHubToken:     "test-token",
		GitHubRepo:      "owner/repo",
		GitHubEventName: "pull_request",
		GitHubEventPath: eventFile,
		LabelMajor:      "bumper:major",
		LabelMinor:      "bumper:minor",
		LabelPatch:      "bumper:patch",
		LabelNone:       "bumper:none",
	}

	mockClient := github.NewMockClient()
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	b := NewWithClient(cfg, mockClient, logger)

	err = b.Run(context.Background())

	// Should succeed without fetching tags (skip scenario)
	assert.NoError(t, err)
	// GetLatestTag should not be called when bumper:none is set
	_, called := mockClient.Calls["GetLatestTag"]
	assert.False(t, called, "GetLatestTag should not be called when bumper:none label is present")
}
