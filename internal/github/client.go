// Copyright (c) 2024 Inetum Poland.

// Package github provides GitHub API client functionality for the action-bumper.
// It handles authentication, tag management, pull request queries, and event parsing.
//
// The package uses the go-github library for API communication and supports:
//   - Fetching repository tags with semver sorting
//   - Creating annotated tags via the Git API
//   - Querying merged pull requests by commit SHA
//   - Parsing GitHub Action event payloads (push, pull_request)
//
// Authentication is performed using the GITHUB_TOKEN provided as an OAuth2 token.
package github

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Inetum-Poland/action-bumper/internal/config"
	"github.com/Inetum-Poland/action-bumper/internal/semver"
	"github.com/google/go-github/v57/github"
	"golang.org/x/oauth2"
)

const (
	defaultPerPage = 100
	refTagsPrefix  = "refs/tags/"
)

// Client wraps GitHub API client
type Client struct {
	client         *github.Client
	owner          string
	repo           string
	debugEventPath string // For testing: path to directory with tags.json, pull_request.json
}

// NewClient creates a new GitHub client
func NewClient(ctx context.Context, cfg *config.Config) (*Client, error) {
	// Parse owner/repo from GITHUB_REPOSITORY
	parts := strings.Split(cfg.GitHubRepo, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("invalid GITHUB_REPOSITORY format: %s (expected owner/repo)", cfg.GitHubRepo)
	}

	// Create OAuth2 token source
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.GitHubToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	return &Client{
		client:         github.NewClient(tc),
		owner:          parts[0],
		repo:           parts[1],
		debugEventPath: cfg.DebugEventPath,
	}, nil
}

// GetLatestTag fetches the latest semantic version tag from the repository.
// If debugEventPath is set, it reads from tags.json instead of the API.
// It finds the highest semver version among all tags.
func (c *Client) GetLatestTag(ctx context.Context) (*semver.Version, error) {
	// Debug mode: read from local tags.json file
	if c.debugEventPath != "" {
		return c.getLatestTagFromFile()
	}

	// Normal mode: fetch from GitHub API
	return c.getLatestTagFromAPI(ctx)
}

// getLatestTagFromFile reads tags from a local JSON file (for testing).
// Matches Bash behavior: returns the first non-"latest" tag from the array.
func (c *Client) getLatestTagFromFile() (*semver.Version, error) {
	tagsFile := c.debugEventPath + "/tags.json"
	data, err := os.ReadFile(tagsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read tags file: %w", err)
	}

	// Parse the tags array
	var tags []struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(data, &tags); err != nil {
		return nil, fmt.Errorf("failed to parse tags JSON: %w", err)
	}

	// Return the first non-"latest" tag (matching Bash behavior: .[0].name)
	for _, tag := range tags {
		if tag.Name == "latest" {
			continue
		}
		return semver.Parse(tag.Name)
	}

	return nil, nil
}

// getLatestTagFromAPI fetches tags from the GitHub API
func (c *Client) getLatestTagFromAPI(ctx context.Context) (*semver.Version, error) {
	var latestVersion *semver.Version

	opts := &github.ListOptions{
		PerPage: defaultPerPage,
	}

	for {
		tags, resp, err := c.client.Repositories.ListTags(ctx, c.owner, c.repo, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to list tags: %w", err)
		}

		// Find the latest valid semver tag in this page
		for _, tag := range tags {
			v, err := semver.Parse(tag.GetName())
			if err != nil {
				// Skip invalid semver tags
				continue
			}

			if latestVersion == nil || v.GreaterThan(latestVersion) {
				latestVersion = v
			}
		}

		// Check if there are more pages
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return latestVersion, nil
}

// ParseEvent parses the GitHub event from GITHUB_EVENT_PATH
func ParseEvent(eventPath string) (*Event, error) {
	data, err := os.ReadFile(eventPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read event file: %w", err)
	}

	var event Event
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, fmt.Errorf("failed to parse event JSON: %w", err)
	}

	return &event, nil
}

// Event represents a GitHub webhook event
type Event struct {
	Action      string       `json:"action"`
	Number      int          `json:"number"`
	PullRequest *PullRequest `json:"pull_request"`
	Repository  *Repository  `json:"repository"`
	Commits     []Commit     `json:"commits"`
	HeadCommit  *Commit      `json:"head_commit"`
	After       string       `json:"after"`
}

// PullRequest represents a pull request in the event
type PullRequest struct {
	Number         int     `json:"number"`
	Title          string  `json:"title"`
	Labels         []Label `json:"labels"`
	MergeCommitSHA string  `json:"merge_commit_sha"`
	Head           *Head   `json:"head"`
}

// Label represents a label on a pull request
type Label struct {
	Name string `json:"name"`
}

// Head represents the head branch of a PR
type Head struct {
	SHA string `json:"sha"`
}

// Repository represents repository information
type Repository struct {
	Name          string `json:"name"`
	FullName      string `json:"full_name"`
	DefaultBranch string `json:"default_branch"`
}

// Commit represents a commit in the event
type Commit struct {
	SHA     string `json:"sha"`
	Message string `json:"message"`
}

// IsPREvent checks if the event is a PR event
func (e *Event) IsPREvent() bool {
	return e.Action != "" && e.PullRequest != nil
}

// IsPushEvent checks if the event is a push event
func (e *Event) IsPushEvent() bool {
	return e.Action == "" && e.After != ""
}

// GetMergedPRByCommitSHA finds a merged PR by its merge commit SHA.
// If debugEventPath is set, it reads from pull_request.json instead of the API.
func (c *Client) GetMergedPRByCommitSHA(ctx context.Context, sha string) (*PullRequest, error) {
	// Debug mode: read from local pull_request.json file
	if c.debugEventPath != "" {
		return c.getMergedPRFromFile(sha)
	}

	// Normal mode: fetch from GitHub API
	return c.getMergedPRFromAPI(ctx, sha)
}

// getMergedPRFromFile reads PRs from a local JSON file (for testing)
func (c *Client) getMergedPRFromFile(sha string) (*PullRequest, error) {
	prFile := c.debugEventPath + "/pull_request.json"
	data, err := os.ReadFile(prFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read pull_request file: %w", err)
	}

	// Parse the PRs array
	var prs []PullRequest
	if err := json.Unmarshal(data, &prs); err != nil {
		return nil, fmt.Errorf("failed to parse pull_request JSON: %w", err)
	}

	// Find the PR with matching merge commit SHA
	for _, pr := range prs {
		if pr.MergeCommitSHA == sha {
			return &pr, nil
		}
	}

	return nil, fmt.Errorf("no merged PR found for commit SHA: %s", sha)
}

// getMergedPRFromAPI fetches the merged PR from the GitHub API
func (c *Client) getMergedPRFromAPI(ctx context.Context, sha string) (*PullRequest, error) {
	// List closed PRs sorted by updated time (most recent first)
	pulls, _, err := c.client.PullRequests.List(ctx, c.owner, c.repo, &github.PullRequestListOptions{
		State:     "closed",
		Sort:      "updated",
		Direction: "desc",
		ListOptions: github.ListOptions{
			PerPage: defaultPerPage,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list pull requests: %w", err)
	}

	// Find the PR with matching merge commit SHA
	for _, pr := range pulls {
		if pr.GetMergeCommitSHA() == sha {
			// Convert to our PullRequest type
			labels := make([]Label, 0, len(pr.Labels))
			for _, label := range pr.Labels {
				labels = append(labels, Label{Name: label.GetName()})
			}

			return &PullRequest{
				Number:         pr.GetNumber(),
				Title:          pr.GetTitle(),
				Labels:         labels,
				MergeCommitSHA: pr.GetMergeCommitSHA(),
			}, nil
		}
	}

	return nil, fmt.Errorf("no merged PR found for commit SHA: %s", sha)
}
