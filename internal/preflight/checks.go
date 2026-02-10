// Copyright (c) 2024 Inetum Poland.

// Package preflight provides pre-flight checks for the bumper application.
// These checks verify that the required external dependencies and environment
// are properly configured before the main application logic runs.
package preflight

import (
	"context"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

// Result represents the result of a pre-flight check.
type Result struct {
	Name    string
	Passed  bool
	Message string
	Error   error
}

// Checker performs pre-flight checks.
type Checker struct {
	timeout time.Duration
}

// NewChecker creates a new pre-flight checker.
func NewChecker() *Checker {
	return &Checker{
		timeout: 10 * time.Second,
	}
}

// WithTimeout sets the timeout for network checks.
func (c *Checker) WithTimeout(d time.Duration) *Checker {
	c.timeout = d
	return c
}

// CheckAll runs all pre-flight checks and returns the results.
func (c *Checker) CheckAll(ctx context.Context) []Result {
	return []Result{
		c.CheckGitAvailable(ctx),
		c.CheckGitHubReachable(ctx),
	}
}

// CheckGitAvailable verifies that git is installed and accessible.
func (c *Checker) CheckGitAvailable(ctx context.Context) Result {
	result := Result{Name: "git-available"}

	cmd := exec.CommandContext(ctx, "git", "--version")
	output, err := cmd.Output()
	if err != nil {
		result.Passed = false
		result.Message = "git is not available"
		result.Error = fmt.Errorf("git not found: %w", err)
		return result
	}

	version := strings.TrimSpace(string(output))
	result.Passed = true
	result.Message = version
	return result
}

// CheckGitHubReachable verifies that the GitHub API is reachable.
func (c *Checker) CheckGitHubReachable(ctx context.Context) Result {
	result := Result{Name: "github-reachable"}

	client := &http.Client{
		Timeout: c.timeout,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/zen", http.NoBody)
	if err != nil {
		result.Passed = false
		result.Message = "failed to create request"
		result.Error = err
		return result
	}

	resp, err := client.Do(req)
	if err != nil {
		result.Passed = false
		result.Message = "GitHub API is unreachable"
		result.Error = fmt.Errorf("failed to reach GitHub API: %w", err)
		return result
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		result.Passed = false
		result.Message = fmt.Sprintf("GitHub API returned status %d", resp.StatusCode)
		result.Error = fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		return result
	}

	result.Passed = true
	result.Message = "GitHub API is reachable"
	return result
}

// CheckRequired returns an error if any required checks failed.
func CheckRequired(results []Result) error {
	var failures []string
	for _, r := range results {
		if !r.Passed {
			failures = append(failures, fmt.Sprintf("%s: %s", r.Name, r.Message))
		}
	}
	if len(failures) > 0 {
		return fmt.Errorf("pre-flight checks failed: %s", strings.Join(failures, "; "))
	}
	return nil
}
