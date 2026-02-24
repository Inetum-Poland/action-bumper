// Copyright (c) 2024-2026 Inetum Poland.

package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRepo(t *testing.T) string {
	t.Helper()

	tmpDir := t.TempDir()

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "git init failed: %s", output)

	// Configure git user
	cmd = exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = tmpDir
	_, err = cmd.CombinedOutput()
	require.NoError(t, err)

	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = tmpDir
	_, err = cmd.CombinedOutput()
	require.NoError(t, err)

	// Create initial commit
	readmePath := filepath.Join(tmpDir, "README.md")
	err = os.WriteFile(readmePath, []byte("# Test\n"), 0o644)
	require.NoError(t, err)

	cmd = exec.Command("git", "add", ".")
	cmd.Dir = tmpDir
	_, err = cmd.CombinedOutput()
	require.NoError(t, err)

	cmd = exec.Command("git", "commit", "-m", "Initial commit")
	cmd.Dir = tmpDir
	_, err = cmd.CombinedOutput()
	require.NoError(t, err)

	return tmpDir
}

func TestConfigureUser(t *testing.T) {
	repoDir := setupTestRepo(t)

	// Change to repo directory for the test
	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(repoDir)

	t.Run("configure both name and email", func(t *testing.T) {
		err := ConfigureUser("New Name", "new@example.com")
		assert.NoError(t, err)

		// Verify name was set
		cmd := exec.Command("git", "config", "user.name")
		output, err := cmd.Output()
		require.NoError(t, err)
		assert.Contains(t, string(output), "New Name")

		// Verify email was set
		cmd = exec.Command("git", "config", "user.email")
		output, err = cmd.Output()
		require.NoError(t, err)
		assert.Contains(t, string(output), "new@example.com")
	})

	t.Run("configure only name", func(t *testing.T) {
		err := ConfigureUser("Another Name", "")
		assert.NoError(t, err)
	})

	t.Run("configure only email", func(t *testing.T) {
		err := ConfigureUser("", "another@example.com")
		assert.NoError(t, err)
	})

	t.Run("configure neither", func(t *testing.T) {
		err := ConfigureUser("", "")
		assert.NoError(t, err)
	})
}

func TestCreateTag(t *testing.T) {
	repoDir := setupTestRepo(t)

	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(repoDir)

	t.Run("create new tag", func(t *testing.T) {
		err := CreateTag("v1.0.0", "Release v1.0.0")
		assert.NoError(t, err)

		// Verify tag exists
		assert.True(t, TagExists("v1.0.0"))
	})

	t.Run("create tag that already exists fails", func(t *testing.T) {
		// First create a tag
		err := CreateTag("v2.0.0", "Release v2.0.0")
		require.NoError(t, err)

		// Try to create same tag again - should fail
		err = CreateTag("v2.0.0", "Release v2.0.0 again")
		assert.Error(t, err)
	})
}

func TestCreateOrUpdateTag(t *testing.T) {
	repoDir := setupTestRepo(t)

	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(repoDir)

	t.Run("create new tag with force", func(t *testing.T) {
		err := CreateOrUpdateTag("v3.0.0", "Release v3.0.0", "")
		assert.NoError(t, err)
		assert.True(t, TagExists("v3.0.0"))
	})

	t.Run("update existing tag", func(t *testing.T) {
		// Create initial tag
		err := CreateOrUpdateTag("v4.0.0", "Initial", "")
		require.NoError(t, err)

		// Update the tag
		err = CreateOrUpdateTag("v4.0.0", "Updated", "")
		assert.NoError(t, err)
	})
}

func TestDeleteTag(t *testing.T) {
	repoDir := setupTestRepo(t)

	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(repoDir)

	t.Run("delete existing tag", func(t *testing.T) {
		// Create tag first
		err := CreateTag("v5.0.0", "Release v5.0.0")
		require.NoError(t, err)
		require.True(t, TagExists("v5.0.0"))

		// Delete it
		err = DeleteTag("v5.0.0")
		assert.NoError(t, err)
		assert.False(t, TagExists("v5.0.0"))
	})

	t.Run("delete non-existing tag fails", func(t *testing.T) {
		err := DeleteTag("v999.0.0")
		assert.Error(t, err)
	})
}

func TestTagExists(t *testing.T) {
	repoDir := setupTestRepo(t)

	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(repoDir)

	t.Run("existing tag", func(t *testing.T) {
		err := CreateTag("v6.0.0", "Release")
		require.NoError(t, err)

		assert.True(t, TagExists("v6.0.0"))
	})

	t.Run("non-existing tag", func(t *testing.T) {
		assert.False(t, TagExists("v888.0.0"))
	})
}

func TestGetCurrentCommit(t *testing.T) {
	repoDir := setupTestRepo(t)

	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(repoDir)

	sha, err := GetCurrentCommit()
	assert.NoError(t, err)
	assert.Len(t, sha, 40) // Git SHA is 40 characters
}
