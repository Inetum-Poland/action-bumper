// Copyright (c) 2024 Inetum Poland.

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

// CreateTag creates an annotated git tag
func CreateTag(tag, message string) error {
	cmd := exec.Command("git", "tag", "-a", tag, "-m", message)
	if output, err := cmd.CombinedOutput(); err != nil {
		// Check if tag already exists
		if strings.Contains(string(output), "already exists") {
			// Tag exists, update it
			if err := DeleteTag(tag); err != nil {
				return err
			}
			// Try creating again
			cmd = exec.Command("git", "tag", "-a", tag, "-m", message)
			if output, err := cmd.CombinedOutput(); err != nil {
				return fmt.Errorf("failed to create tag after delete: %w\n%s", err, output)
			}
			return nil
		}
		return fmt.Errorf("failed to create tag: %w\n%s", err, output)
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

// PushTag pushes a tag to remote
func PushTag(tag string) error {
	cmd := exec.Command("git", "push", "origin", tag, "--force")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to push tag: %w\n%s", err, output)
	}
	return nil
}

// PushTags pushes multiple tags to remote
func PushTags(tags []string) error {
	for _, tag := range tags {
		if err := PushTag(tag); err != nil {
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
