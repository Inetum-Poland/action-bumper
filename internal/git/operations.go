// Copyright (c) 2024 Inetum Poland.

// Package git provides low-level git command wrappers for the action-bumper.
// It executes git commands via exec.Command and handles output parsing.
//
// Tag Operations:
//   - CreateTag: Create new annotated tags (fails if tag exists)
//   - CreateOrUpdateTag: Create or update tags using -fa flag (for semver/latest)
//   - DeleteTag: Remove local tags
//   - PushTag/PushTagForce: Push tags to remote with or without force
//   - PushTags: Push multiple tags (first without force, rest with force)
//
// Configuration:
//   - ConfigureSafeDirectory: Set up safe.directory for container environments
//   - ConfigureUser: Set git user.name and user.email
//   - SetRemoteURL: Configure remote with authentication token
//
// Queries:
//   - GetCurrentCommit: Get HEAD commit SHA
//   - TagExists: Check if a tag exists locally
package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// ConfigureSafeDirectory configures git safe directory
func ConfigureSafeDirectory(dir string) error {
	cmd := exec.Command("git", "config", "--global", "--add", "safe.directory", dir)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to configure safe directory: %w\n%s", err, output)
	}
	return nil
}

// ConfigureUser configures git user name and email
func ConfigureUser(name, email string) error {
	if name != "" {
		cmd := exec.Command("git", "config", "user.name", name)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to configure user name: %w\n%s", err, output)
		}
	}

	if email != "" {
		cmd := exec.Command("git", "config", "user.email", email)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to configure user email: %w\n%s", err, output)
		}
	}

	return nil
}

// CreateTag creates an annotated git tag (does not force)
func CreateTag(tag, message string) error {
	cmd := exec.Command("git", "tag", "-a", tag, "-m", message)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to create tag: %w\n%s", err, output)
	}
	return nil
}

// CreateOrUpdateTag creates or updates an annotated git tag using -fa flag
// This should be used for semver tags (v1, v1.2) and latest tag
// The refSpec can be a tag name with ^{commit} to target specific commit
func CreateOrUpdateTag(tag, message, refSpec string) error {
	// Use -fa to force create/update annotated tag
	// refSpec should be like "v1.2.3^{commit}" to target the commit of another tag
	var cmd *exec.Cmd
	if refSpec != "" {
		cmd = exec.Command("git", "tag", "-fa", tag, refSpec, "-m", message)
	} else {
		cmd = exec.Command("git", "tag", "-fa", tag, "-m", message)
	}
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to create/update tag: %w\n%s", err, output)
	}
	return nil
}

// DeleteTag deletes a local git tag
func DeleteTag(tag string) error {
	cmd := exec.Command("git", "tag", "-d", tag)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to delete tag: %w\n%s", err, output)
	}
	return nil
}

// PushTag pushes a tag to remote (without force)
func PushTag(tag string) error {
	cmd := exec.Command("git", "push", "origin", tag)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to push tag: %w\n%s", err, output)
	}
	return nil
}

// PushTagForce pushes a tag to remote with --force flag
// Use for semver tags (v1, v1.2) and latest tag that may already exist
func PushTagForce(tag string) error {
	cmd := exec.Command("git", "push", "--force", "origin", tag)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to push tag: %w\n%s", err, output)
	}
	return nil
}

// PushTags pushes multiple tags to remote (first tag without force, rest with force)
func PushTags(tags []string) error {
	for i, tag := range tags {
		var err error
		if i == 0 {
			// Primary version tag - no force
			err = PushTag(tag)
		} else {
			// Semver/latest tags - force push
			err = PushTagForce(tag)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// SetRemoteURL sets the remote URL with authentication token
func SetRemoteURL(token, repo string) error {
	// Format: https://x-access-token:TOKEN@github.com/owner/repo
	url := fmt.Sprintf("https://x-access-token:%s@github.com/%s", token, repo)

	cmd := exec.Command("git", "remote", "set-url", "origin", url)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to set remote URL: %w\n%s", err, output)
	}
	return nil
}

// GetCurrentCommit returns the current commit SHA
func GetCurrentCommit() (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current commit: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// TagExists checks if a tag exists locally
func TagExists(tag string) bool {
	cmd := exec.Command("git", "rev-parse", tag)
	err := cmd.Run()
	return err == nil
}
