// Copyright (c) 2024-2026 Inetum Poland.

package bumper

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Inetum-Poland/action-bumper/internal/config"
	"github.com/Inetum-Poland/action-bumper/internal/github"
	"github.com/Inetum-Poland/action-bumper/internal/semver"
)

// TestGeneratePRStatusMessage_WithSemverTags tests semver tag info in PR messages
func TestGeneratePRStatusMessage_WithSemverTags(t *testing.T) {
	tests := []struct {
		name         string
		cfg          *config.Config
		wantContains []string
		wantMissing  []string
	}{
		{
			name: "with semver enabled",
			cfg: &config.Config{
				GitHubRepo:   "owner/repo",
				BumpIncludeV: true,
				BumpSemver:   true,
				BumpLatest:   false,
			},
			wantContains: []string{"v1.2.3", "v1.2", "v1"},
			wantMissing:  []string{"latest"},
		},
		{
			name: "with latest enabled",
			cfg: &config.Config{
				GitHubRepo:   "owner/repo",
				BumpIncludeV: true,
				BumpSemver:   false,
				BumpLatest:   true,
			},
			wantContains: []string{"v1.2.3", "latest"},
		},
		{
			name: "without v prefix",
			cfg: &config.Config{
				GitHubRepo:   "owner/repo",
				BumpIncludeV: false,
				BumpSemver:   true,
				BumpLatest:   false,
			},
			wantContains: []string{"1.2.3", "1.2", "1"},
			wantMissing:  []string{"v1.2.3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bumper{cfg: tt.cfg}
			current, _ := semver.Parse("1.2.2")
			next, _ := semver.Parse("1.2.3")

			msg := b.generatePRStatusMessage(current, next, config.BumpLevelPatch)

			for _, want := range tt.wantContains {
				assert.Contains(t, msg, want, "Message should contain %q", want)
			}
			for _, missing := range tt.wantMissing {
				assert.NotContains(t, msg, missing, "Message should not contain %q", missing)
			}
		})
	}
}

// TestDetermineBumpLevel_AllPriorityCombinations tests all label priority combinations
func TestDetermineBumpLevel_AllPriorityCombinations(t *testing.T) {
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
			name: "patch + minor = minor wins",
			labels: []github.Label{
				{Name: "bumper:patch"},
				{Name: "bumper:minor"},
			},
			want: config.BumpLevelMinor,
		},
		{
			name: "none + patch = patch wins",
			labels: []github.Label{
				{Name: "bumper:none"},
				{Name: "bumper:patch"},
			},
			want: config.BumpLevelPatch,
		},
		{
			name: "all labels present = major wins",
			labels: []github.Label{
				{Name: "bumper:none"},
				{Name: "bumper:patch"},
				{Name: "bumper:minor"},
				{Name: "bumper:major"},
			},
			want: config.BumpLevelMajor,
		},
		{
			name: "with other labels mixed in = major wins",
			labels: []github.Label{
				{Name: "feature"},
				{Name: "bumper:patch"},
				{Name: "enhancement"},
				{Name: "bumper:major"},
				{Name: "bug"},
			},
			want: config.BumpLevelMajor,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := b.determineBumpLevel(tt.labels)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestHandlePREvent_WithDefaultLevel tests default level fallback
func TestHandlePREvent_WithDefaultLevel(t *testing.T) {
	tmpDir := t.TempDir()
	eventFile := tmpDir + "/event.json"
	eventJSON := `{
		"action": "opened",
		"pull_request": {
			"number": 1,
			"title": "No bump label",
			"labels": [{"name": "feature"}]
		}
	}`
	err := os.WriteFile(eventFile, []byte(eventJSON), 0o644)
	require.NoError(t, err)

	cfg := &config.Config{
		GitHubToken:       "test-token",
		GitHubRepo:        "owner/repo",
		GitHubEventName:   "pull_request",
		GitHubEventPath:   eventFile,
		LabelMajor:        "bumper:major",
		LabelMinor:        "bumper:minor",
		LabelPatch:        "bumper:patch",
		LabelNone:         "bumper:none",
		BumpDefaultLevel:  config.BumpLevelPatch,
		BumpFailIfNoLevel: false,
		BumpIncludeV:      true,
	}

	currentVersion, _ := semver.Parse("1.0.0")
	mockClient := github.NewMockClient()
	mockClient.GetLatestTagFunc = func(_ context.Context) (*semver.Version, error) {
		return currentVersion, nil
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	b := NewWithClient(cfg, mockClient, logger)

	err = b.Run(context.Background())

	assert.NoError(t, err)
	assert.Contains(t, mockClient.Calls, "GetLatestTag")
}

// TestHandlePREvent_FailIfNoLevelWithDefaultLevel tests that enforcement requires an explicit label.
func TestHandlePREvent_FailIfNoLevelWithDefaultLevel(t *testing.T) {
	tmpDir := t.TempDir()
	eventFile := tmpDir + "/event.json"
	eventJSON := `{
		"action": "opened",
		"pull_request": {
			"number": 1,
			"title": "No bump label",
			"labels": []
		}
	}`
	err := os.WriteFile(eventFile, []byte(eventJSON), 0o644)
	require.NoError(t, err)

	cfg := &config.Config{
		GitHubToken:       "test-token",
		GitHubRepo:        "owner/repo",
		GitHubEventName:   "pull_request",
		GitHubEventPath:   eventFile,
		LabelMajor:        "bumper:major",
		LabelMinor:        "bumper:minor",
		LabelPatch:        "bumper:patch",
		LabelNone:         "bumper:none",
		BumpDefaultLevel:  config.BumpLevelMinor,
		BumpFailIfNoLevel: true,
		BumpIncludeV:      true,
	}

	currentVersion, _ := semver.Parse("1.0.0")
	mockClient := github.NewMockClient()
	mockClient.GetLatestTagFunc = func(_ context.Context) (*semver.Version, error) {
		return currentVersion, nil
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	b := NewWithClient(cfg, mockClient, logger)

	err = b.Run(context.Background())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no bump level label found")
}

// TestVersionTagGeneration tests version tag string generation
func TestVersionTagGeneration(t *testing.T) {
	tests := []struct {
		name      string
		version   string
		includeV  bool
		wantFull  string
		wantMajor string
		wantMinor string
	}{
		{
			name:      "v1.2.3 with v prefix",
			version:   "1.2.3",
			includeV:  true,
			wantFull:  "v1.2.3",
			wantMajor: "v1",
			wantMinor: "v1.2",
		},
		{
			name:      "v1.2.3 without v prefix",
			version:   "1.2.3",
			includeV:  false,
			wantFull:  "1.2.3",
			wantMajor: "1",
			wantMinor: "1.2",
		},
		{
			name:      "v10.20.30 with v prefix",
			version:   "10.20.30",
			includeV:  true,
			wantFull:  "v10.20.30",
			wantMajor: "v10",
			wantMinor: "v10.20",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := semver.Parse(tt.version)
			require.NoError(t, err)

			assert.Equal(t, tt.wantFull, v.FullTag(tt.includeV))
			assert.Equal(t, tt.wantMajor, v.MajorTag(tt.includeV))
			assert.Equal(t, tt.wantMinor, v.MinorTag(tt.includeV))
		})
	}
}
