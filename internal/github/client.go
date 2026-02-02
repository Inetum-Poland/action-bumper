// Copyright (c) 2024 Inetum Poland.

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
	client *github.Client
	owner  string
	repo   string
}

// NewClient creates a new GitHub client
func NewClient(ctx context.Context, cfg *config.Config) (*Client, error) {
	// Parse owner/repo from GITHUB_REPOSITORY
	parts := strings.Split(cfg.GitHubRepo, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid GITHUB_REPOSITORY format: %s (expected owner/repo)", cfg.GitHubRepo)
	}

	// Create OAuth2 token source
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.GitHubToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	return &Client{
		client: github.NewClient(tc),
		owner:  parts[0],
		repo:   parts[1],
	}, nil
}

// GetLatestTag fetches the latest semantic version tag from the repository
func (c *Client) GetLatestTag(ctx context.Context) (*semver.Version, error) {
	tags, _, err := c.client.Repositories.ListTags(ctx, c.owner, c.repo, &github.ListOptions{
		PerPage: defaultPerPage,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}

	// Find the latest valid semver tag
	var latestVersion *semver.Version
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

	return latestVersion, nil
}

// CreateTag creates an annotated tag in the repository
func (c *Client) CreateTag(ctx context.Context, tagName, sha, message string) error {
	// Create tag object
	tag := &github.Tag{
		Tag:     github.String(tagName),
		Message: github.String(message),
		Object: &github.GitObject{
			Type: github.String("commit"),
			SHA:  github.String(sha),
		},
	}

	_, _, err := c.client.Git.CreateTag(ctx, c.owner, c.repo, tag)
	if err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}

	return nil
}

// CreateReference creates or updates a reference (tag)
func (c *Client) CreateReference(ctx context.Context, ref, sha string) error {
	refName := refTagsPrefix + ref

	reference := &github.Reference{
		Ref: github.String(refName),
		Object: &github.GitObject{
			SHA: github.String(sha),
		},
	}

	_, _, err := c.client.Git.CreateRef(ctx, c.owner, c.repo, reference)
	if err != nil {
		// If ref exists, try to update it
		if strings.Contains(err.Error(), "Reference already exists") {
			_, _, err = c.client.Git.UpdateRef(ctx, c.owner, c.repo, reference, false)
			if err != nil {
				return fmt.Errorf("failed to update reference: %w", err)
			}
			return nil
		}
		return fmt.Errorf("failed to create reference: %w", err)
	}

	return nil
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
