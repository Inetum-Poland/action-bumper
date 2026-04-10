// Copyright (c) 2024-2026 Inetum Poland.

package preflight

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewChecker(t *testing.T) {
	c := NewChecker()
	assert.NotNil(t, c)
	assert.Equal(t, 10*time.Second, c.timeout)
}

func TestWithTimeout(t *testing.T) {
	c := NewChecker().WithTimeout(5 * time.Second)
	assert.Equal(t, 5*time.Second, c.timeout)
}

func TestCheckGitAvailable(t *testing.T) {
	c := NewChecker()
	ctx := context.Background()

	result := c.CheckGitAvailable(ctx)

	assert.Equal(t, "git-available", result.Name)
	// Git should be available in the test environment
	assert.True(t, result.Passed, "git should be available")
	assert.Contains(t, result.Message, "git version")
	assert.NoError(t, result.Error)
}

func TestCheckGitAvailable_Canceled(t *testing.T) {
	c := NewChecker()
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	result := c.CheckGitAvailable(ctx)

	// Context canceled, should fail
	assert.Equal(t, "git-available", result.Name)
	assert.False(t, result.Passed)
	assert.Error(t, result.Error)
}

func TestCheckGitHubReachable(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping network test in short mode")
	}

	c := NewChecker().WithTimeout(30 * time.Second)
	ctx := context.Background()

	result := c.CheckGitHubReachable(ctx)

	assert.Equal(t, "github-reachable", result.Name)
	// In a normal environment with internet, this should pass
	// but we don't want to fail CI if there are network issues
	if result.Passed {
		assert.Equal(t, "GitHub API is reachable", result.Message)
		assert.NoError(t, result.Error)
	}
}

func TestCheckGitHubReachable_Timeout(t *testing.T) {
	c := NewChecker().WithTimeout(1 * time.Nanosecond) // Very short timeout
	ctx := context.Background()

	result := c.CheckGitHubReachable(ctx)

	assert.Equal(t, "github-reachable", result.Name)
	// Should fail due to timeout
	assert.False(t, result.Passed)
	assert.Error(t, result.Error)
}

func TestCheckAll(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping network test in short mode")
	}

	c := NewChecker()
	ctx := context.Background()

	results := c.CheckAll(ctx)

	assert.Len(t, results, 2)
	assert.Equal(t, "git-available", results[0].Name)
	assert.Equal(t, "github-reachable", results[1].Name)
}

func TestCheckRequired_AllPass(t *testing.T) {
	results := []Result{
		{Name: "check1", Passed: true, Message: "ok"},
		{Name: "check2", Passed: true, Message: "ok"},
	}

	err := CheckRequired(results)
	assert.NoError(t, err)
}

func TestCheckRequired_OneFails(t *testing.T) {
	results := []Result{
		{Name: "check1", Passed: true, Message: "ok"},
		{Name: "check2", Passed: false, Message: "failed"},
	}

	err := CheckRequired(results)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "pre-flight checks failed")
	assert.Contains(t, err.Error(), "check2: failed")
}

func TestCheckRequired_MultipleFail(t *testing.T) {
	results := []Result{
		{Name: "check1", Passed: false, Message: "error1"},
		{Name: "check2", Passed: false, Message: "error2"},
	}

	err := CheckRequired(results)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "check1: error1")
	assert.Contains(t, err.Error(), "check2: error2")
}

func TestCheckRequired_Empty(t *testing.T) {
	err := CheckRequired([]Result{})
	assert.NoError(t, err)
}
