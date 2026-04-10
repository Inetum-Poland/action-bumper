// Copyright (c) 2024-2026 Inetum Poland.

package github

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Inetum-Poland/action-bumper/internal/config"
	"github.com/Inetum-Poland/action-bumper/internal/semver"
)

func TestNewClient_ValidConfig(t *testing.T) {
	cfg := &config.Config{
		GitHubToken: "test-token",
		GitHubRepo:  "owner/repo",
	}

	client, err := NewClient(context.Background(), cfg)

	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, "owner", client.owner)
	assert.Equal(t, "repo", client.repo)
}

func TestNewClient_InvalidRepoFormat(t *testing.T) {
	tests := []struct {
		name string
		repo string
	}{
		{"no slash", "ownerrepo"},
		{"empty", ""},
		{"just slash", "/"},
		{"multiple slashes", "owner/repo/extra"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				GitHubToken: "test-token",
				GitHubRepo:  tt.repo,
			}

			client, err := NewClient(context.Background(), cfg)

			assert.Error(t, err)
			assert.Nil(t, client)
			assert.Contains(t, err.Error(), "invalid GITHUB_REPOSITORY format")
		})
	}
}

func TestNewClient_OwnerRepoParsing(t *testing.T) {
	tests := []struct {
		name          string
		repo          string
		expectedOwner string
		expectedRepo  string
	}{
		{
			name:          "simple",
			repo:          "myorg/myrepo",
			expectedOwner: "myorg",
			expectedRepo:  "myrepo",
		},
		{
			name:          "with hyphen",
			repo:          "my-org/my-repo",
			expectedOwner: "my-org",
			expectedRepo:  "my-repo",
		},
		{
			name:          "with underscore",
			repo:          "my_org/my_repo",
			expectedOwner: "my_org",
			expectedRepo:  "my_repo",
		},
		{
			name:          "with numbers",
			repo:          "org123/repo456",
			expectedOwner: "org123",
			expectedRepo:  "repo456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				GitHubToken: "test-token",
				GitHubRepo:  tt.repo,
			}

			client, err := NewClient(context.Background(), cfg)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedOwner, client.owner)
			assert.Equal(t, tt.expectedRepo, client.repo)
		})
	}
}

func TestMockClient_GetLatestTag(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func(*MockClient)
		wantVersion string
		wantError   bool
	}{
		{
			name: "returns version",
			setupMock: func(m *MockClient) {
				m.GetLatestTagFunc = func(_ context.Context) (*semver.Version, error) {
					return semver.MustParse("1.2.3"), nil
				}
			},
			wantVersion: "1.2.3",
			wantError:   false,
		},
		{
			name: "returns nil when no tags",
			setupMock: func(m *MockClient) {
				m.GetLatestTagFunc = func(_ context.Context) (*semver.Version, error) {
					return nil, nil
				}
			},
			wantVersion: "",
			wantError:   false,
		},
		{
			name: "returns error",
			setupMock: func(m *MockClient) {
				m.GetLatestTagFunc = func(_ context.Context) (*semver.Version, error) {
					return nil, assert.AnError
				}
			},
			wantVersion: "",
			wantError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := NewMockClient()
			tt.setupMock(mock)

			v, err := mock.GetLatestTag(context.Background())

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			if tt.wantVersion != "" {
				assert.NotNil(t, v)
				assert.Equal(t, tt.wantVersion, v.String())
			}
			// Verify call was recorded
			assert.Contains(t, mock.Calls, "GetLatestTag")
		})
	}
}

func TestMockClient_GetMergedPRByCommitSHA(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(*MockClient)
		sha       string
		wantPR    *PullRequest
		wantError bool
	}{
		{
			name: "returns PR",
			setupMock: func(m *MockClient) {
				m.GetMergedPRByCommitSHAFunc = func(_ context.Context, _ string) (*PullRequest, error) {
					return &PullRequest{
						Number: 42,
						Title:  "Test PR",
						Labels: []Label{{Name: "bumper:minor"}},
					}, nil
				}
			},
			sha: "abc123",
			wantPR: &PullRequest{
				Number: 42,
				Title:  "Test PR",
				Labels: []Label{{Name: "bumper:minor"}},
			},
			wantError: false,
		},
		{
			name: "returns nil when no PR found",
			setupMock: func(m *MockClient) {
				m.GetMergedPRByCommitSHAFunc = func(_ context.Context, _ string) (*PullRequest, error) {
					return nil, nil
				}
			},
			sha:       "unknown",
			wantPR:    nil,
			wantError: false,
		},
		{
			name: "returns error",
			setupMock: func(m *MockClient) {
				m.GetMergedPRByCommitSHAFunc = func(_ context.Context, _ string) (*PullRequest, error) {
					return nil, assert.AnError
				}
			},
			sha:       "abc123",
			wantPR:    nil,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := NewMockClient()
			tt.setupMock(mock)

			pr, err := mock.GetMergedPRByCommitSHA(context.Background(), tt.sha)

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			if tt.wantPR != nil {
				assert.NotNil(t, pr)
				assert.Equal(t, tt.wantPR.Number, pr.Number)
				assert.Equal(t, tt.wantPR.Title, pr.Title)
			} else if !tt.wantError {
				assert.Nil(t, pr)
			}
			// Verify call was recorded with SHA
			assert.Contains(t, mock.Calls, "GetMergedPRByCommitSHA")
		})
	}
}

func TestMockClient_WithLatestTag(t *testing.T) {
	v := semver.MustParse("3.2.1")
	mock := NewMockClient().WithLatestTag(v)

	result, err := mock.GetLatestTag(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, "3.2.1", result.String())
}

func TestMockClient_WithMergedPR(t *testing.T) {
	pr := &PullRequest{
		Number: 100,
		Title:  "Big Feature",
		Labels: []Label{{Name: "bumper:major"}},
	}
	mock := NewMockClient().WithMergedPR(pr)

	result, err := mock.GetMergedPRByCommitSHA(context.Background(), "any-sha")

	assert.NoError(t, err)
	assert.Equal(t, 100, result.Number)
	assert.Equal(t, "Big Feature", result.Title)
}

func TestMockClient_CallRecording(t *testing.T) {
	mock := NewMockClient()

	// Make multiple calls
	_, _ = mock.GetLatestTag(context.Background())
	_, _ = mock.GetLatestTag(context.Background())
	_, _ = mock.GetMergedPRByCommitSHA(context.Background(), "sha1")
	_, _ = mock.GetMergedPRByCommitSHA(context.Background(), "sha2")

	// Verify all calls were recorded
	assert.Len(t, mock.Calls["GetLatestTag"], 2)
	assert.Len(t, mock.Calls["GetMergedPRByCommitSHA"], 2)
}
