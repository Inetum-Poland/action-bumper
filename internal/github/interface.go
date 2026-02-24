// Copyright (c) 2024-2026 Inetum Poland.

package github

import (
	"context"

	"github.com/Inetum-Poland/action-bumper/internal/semver"
)

// ClientInterface defines the interface for GitHub API operations.
// This allows for mocking in tests.
type ClientInterface interface {
	// GetLatestTag fetches the latest semantic version tag from the repository
	GetLatestTag(ctx context.Context) (*semver.Version, error)

	// GetMergedPRByCommitSHA finds a merged PR by its merge commit SHA
	GetMergedPRByCommitSHA(ctx context.Context, sha string) (*PullRequest, error)
}

// Ensure Client implements ClientInterface
var _ ClientInterface = (*Client)(nil)

// MockClient implements ClientInterface for testing.
type MockClient struct {
	GetLatestTagFunc           func(ctx context.Context) (*semver.Version, error)
	GetMergedPRByCommitSHAFunc func(ctx context.Context, sha string) (*PullRequest, error)

	// Track calls for assertions
	Calls map[string][]interface{}
}

// NewMockClient creates a new MockClient with default no-op implementations.
func NewMockClient() *MockClient {
	return &MockClient{
		Calls: make(map[string][]interface{}),
		GetLatestTagFunc: func(_ context.Context) (*semver.Version, error) {
			return nil, nil
		},
		GetMergedPRByCommitSHAFunc: func(_ context.Context, _ string) (*PullRequest, error) {
			return nil, nil
		},
	}
}

func (m *MockClient) recordCall(method string, args ...interface{}) {
	m.Calls[method] = append(m.Calls[method], args)
}

// GetLatestTag implements ClientInterface.
func (m *MockClient) GetLatestTag(ctx context.Context) (*semver.Version, error) {
	m.recordCall("GetLatestTag")
	return m.GetLatestTagFunc(ctx)
}

// GetMergedPRByCommitSHA implements ClientInterface.
func (m *MockClient) GetMergedPRByCommitSHA(ctx context.Context, sha string) (*PullRequest, error) {
	m.recordCall("GetMergedPRByCommitSHA", sha)
	return m.GetMergedPRByCommitSHAFunc(ctx, sha)
}

// WithLatestTag sets the mock to return a specific version from GetLatestTag.
func (m *MockClient) WithLatestTag(v *semver.Version) *MockClient {
	m.GetLatestTagFunc = func(_ context.Context) (*semver.Version, error) {
		return v, nil
	}
	return m
}

// WithMergedPR sets the mock to return a specific PR from GetMergedPRByCommitSHA.
func (m *MockClient) WithMergedPR(pr *PullRequest) *MockClient {
	m.GetMergedPRByCommitSHAFunc = func(_ context.Context, _ string) (*PullRequest, error) {
		return pr, nil
	}
	return m
}
