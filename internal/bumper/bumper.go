// Copyright (c) 2024 Inetum Poland.

package bumper

import (
	"context"
	"fmt"
	"strings"

	"github.com/Inetum-Poland/action-bumper/internal/config"
	"github.com/Inetum-Poland/action-bumper/internal/git"
	"github.com/Inetum-Poland/action-bumper/internal/github"
	"github.com/Inetum-Poland/action-bumper/internal/logger"
	"github.com/Inetum-Poland/action-bumper/internal/output"
	"github.com/Inetum-Poland/action-bumper/internal/semver"
)

// Bumper handles the version bumping logic
type Bumper struct {
	cfg    *config.Config
	client *github.Client
	output *output.Writer
	logger logger.Logger
}

// New creates a new Bumper instance
func New(ctx context.Context, cfg *config.Config, log logger.Logger) (*Bumper, error) {
	client, err := github.NewClient(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create GitHub client: %w", err)
	}

	if log == nil {
		log = logger.NewDefaultLogger()
	}

	return &Bumper{
		cfg:    cfg,
		client: client,
		output: output.NewWriter(),
		logger: log,
	}, nil
}

// Run executes the bumper logic
func (b *Bumper) Run(ctx context.Context) error {
	// Parse GitHub event
	event, err := github.ParseEvent(b.cfg.GitHubEventPath)
	if err != nil {
		return fmt.Errorf("failed to parse event: %w", err)
	}

	// Determine event type and handle accordingly
	if event.IsPREvent() {
		return b.handlePREvent(ctx, event)
	} else if event.IsPushEvent() {
		return b.handlePushEvent(ctx, event)
	}

	return fmt.Errorf("unknown event type")
}

// handlePREvent handles pull request events (preview mode)
func (b *Bumper) handlePREvent(ctx context.Context, event *github.Event) error {
	b.logger.Printf("Handling PR event: action=%s, pr#%d", event.Action, event.Number)

	// Get bump level from labels
	bumpLevel := b.determineBumpLevel(event.PullRequest.Labels)

	// Check if bump level is required
	if bumpLevel == config.BumpLevelEmpty && b.cfg.BumpFailIfNoLevel {
		return fmt.Errorf("no bump level label found and bump_fail_if_no_level is true")
	}

	// Use default if no level found
	if bumpLevel == config.BumpLevelEmpty {
		bumpLevel = b.cfg.BumpDefaultLevel
	}

	// If still empty or none, skip
	if bumpLevel == config.BumpLevelEmpty || bumpLevel == config.BumpLevelNone {
		if err := b.output.SetAll(map[string]string{
			"skip":    "true",
			"message": "No version bump (bumper:none or no label)",
		}); err != nil {
			return fmt.Errorf("failed to set outputs: %w", err)
		}
		b.logger.Println("Skipping bump: no level specified or bumper:none")
		return nil
	}

	// Get current version
	currentVersion, err := b.client.GetLatestTag(ctx)
	if err != nil {
		b.logger.Printf("Warning: failed to get latest tag: %v", err)
		// Use default version
		currentVersion = semver.DefaultVersion(bumpLevel)
	}

	// Calculate next version
	nextVersion := currentVersion.Bump(bumpLevel)

	// Generate status message
	message := b.generatePRStatusMessage(currentVersion, nextVersion, bumpLevel)

	// Set outputs
	outputs := map[string]string{
		"current_version": currentVersion.FullTag(b.cfg.BumpIncludeV),
		"next_version":    nextVersion.FullTag(b.cfg.BumpIncludeV),
		"skip":            "false",
		"message":         fmt.Sprintf("Will bump version %s → %s (%s)", currentVersion.String(), nextVersion.String(), bumpLevel),
	}

	if err := b.output.SetAll(outputs); err != nil {
		return fmt.Errorf("failed to set outputs: %w", err)
	}

	if err := b.output.SetMultiline("tag_status", message); err != nil {
		return fmt.Errorf("failed to set tag_status: %w", err)
	}

	b.logger.Printf("Version bump preview: %s → %s (%s)", currentVersion.String(), nextVersion.String(), bumpLevel)
	return nil
}

// handlePushEvent handles push events (actually create tags)
func (b *Bumper) handlePushEvent(ctx context.Context, event *github.Event) error {
	b.logger.Println("Handling push event")

	// For push events, we need to find the merged PR
	// This is a simplified version - in reality you'd query the GitHub API
	// to find the PR associated with the merge commit

	// For now, we'll check if there's a head_commit
	if event.HeadCommit == nil {
		return fmt.Errorf("no head commit in push event")
	}

	// Try to determine bump level from commit message or default
	bumpLevel := b.cfg.BumpDefaultLevel

	if bumpLevel == config.BumpLevelEmpty {
		b.logger.Println("No default bump level, skipping")
		if err := b.output.Set("skip", "true"); err != nil {
			return fmt.Errorf("failed to set skip output: %w", err)
		}
		return nil
	}

	if bumpLevel == config.BumpLevelNone {
		b.logger.Println("Bump level is 'none', skipping")
		if err := b.output.Set("skip", "true"); err != nil {
			return fmt.Errorf("failed to set skip output: %w", err)
		}
		return nil
	}

	// Get current version
	currentVersion, err := b.client.GetLatestTag(ctx)
	if err != nil {
		b.logger.Printf("Warning: no existing tags, using default")
		currentVersion = semver.DefaultVersion(bumpLevel)
	}

	// Calculate next version
	nextVersion := currentVersion.Bump(bumpLevel)

	// Configure git
	if err := git.ConfigureUser(b.cfg.BumpTagAsUser, b.cfg.BumpTagAsEmail); err != nil {
		return fmt.Errorf("failed to configure git user: %w", err)
	}

	// Set remote URL with authentication
	if err := git.SetRemoteURL(b.cfg.GitHubToken, b.cfg.GitHubRepo); err != nil {
		return fmt.Errorf("failed to set remote URL: %w", err)
	}

	// Create tags
	tagMessage := fmt.Sprintf("Release %s", nextVersion.FullTag(b.cfg.BumpIncludeV))
	tags := []string{nextVersion.FullTag(b.cfg.BumpIncludeV)}

	// Create main tag
	if err := git.CreateTag(tags[0], tagMessage); err != nil {
		return fmt.Errorf("failed to create tag: %w", err)
	}

	// Create semver tags (v1, v1.2) if enabled
	if b.cfg.BumpSemver {
		majorTag := nextVersion.MajorTag(b.cfg.BumpIncludeV)
		minorTag := nextVersion.MinorTag(b.cfg.BumpIncludeV)

		if err := git.CreateTag(majorTag, tagMessage); err != nil {
			b.logger.Printf("Warning: failed to create major tag: %v", err)
		} else {
			tags = append(tags, majorTag)
		}

		if err := git.CreateTag(minorTag, tagMessage); err != nil {
			b.logger.Printf("Warning: failed to create minor tag: %v", err)
		} else {
			tags = append(tags, minorTag)
		}
	}

	// Create latest tag if enabled
	if b.cfg.BumpLatest {
		latestTag := "latest"
		if err := git.CreateTag(latestTag, tagMessage); err != nil {
			b.logger.Printf("Warning: failed to create latest tag: %v", err)
		} else {
			tags = append(tags, latestTag)
		}
	}

	// Push all tags
	if err := git.PushTags(tags); err != nil {
		return fmt.Errorf("failed to push tags: %w", err)
	}

	// Generate status message
	message := b.generatePushStatusMessage(nextVersion, tags)

	// Set outputs
	outputs := map[string]string{
		"current_version": currentVersion.FullTag(b.cfg.BumpIncludeV),
		"next_version":    nextVersion.FullTag(b.cfg.BumpIncludeV),
		"skip":            "false",
		"message":         fmt.Sprintf("Created version %s", nextVersion.FullTag(b.cfg.BumpIncludeV)),
	}

	if err := b.output.SetAll(outputs); err != nil {
		return fmt.Errorf("failed to set outputs: %w", err)
	}

	if err := b.output.SetMultiline("tag_status", message); err != nil {
		return fmt.Errorf("failed to set tag_status: %w", err)
	}

	b.logger.Printf("Successfully created and pushed tags: %v", tags)
	return nil
}

// determineBumpLevel determines the bump level from PR labels
func (b *Bumper) determineBumpLevel(labels []github.Label) config.BumpLevel {
	for _, label := range labels {
		switch label.Name {
		case b.cfg.LabelMajor:
			return config.BumpLevelMajor
		case b.cfg.LabelMinor:
			return config.BumpLevelMinor
		case b.cfg.LabelPatch:
			return config.BumpLevelPatch
		case b.cfg.LabelNone:
			return config.BumpLevelNone
		}
	}
	return config.BumpLevelEmpty
}

// generatePRStatusMessage generates a status message for PR events
func (b *Bumper) generatePRStatusMessage(current, next *semver.Version, level config.BumpLevel) string {
	var sb strings.Builder
	sb.WriteString("## 🏷️ Version Bump Preview\n\n")
	sb.WriteString(fmt.Sprintf("**Current version:** `%s`\n", current.FullTag(b.cfg.BumpIncludeV)))
	sb.WriteString(fmt.Sprintf("**Next version:** `%s`\n", next.FullTag(b.cfg.BumpIncludeV)))
	sb.WriteString(fmt.Sprintf("**Bump level:** `%s`\n\n", level))

	if b.cfg.BumpSemver {
		sb.WriteString("**Additional tags that will be created:**\n")
		sb.WriteString(fmt.Sprintf("- `%s`\n", next.MajorTag(b.cfg.BumpIncludeV)))
		sb.WriteString(fmt.Sprintf("- `%s`\n", next.MinorTag(b.cfg.BumpIncludeV)))
	}

	if b.cfg.BumpLatest {
		sb.WriteString("- `latest`\n")
	}

	return sb.String()
}

// generatePushStatusMessage generates a status message for push events
func (b *Bumper) generatePushStatusMessage(version *semver.Version, tags []string) string {
	var sb strings.Builder
	sb.WriteString("## 🎉 Version Release\n\n")
	sb.WriteString(fmt.Sprintf("**Released version:** `%s`\n\n", version.FullTag(b.cfg.BumpIncludeV)))
	sb.WriteString("**Created tags:**\n")

	for _, tag := range tags {
		sb.WriteString(fmt.Sprintf("- `%s`\n", tag))
	}

	return sb.String()
}
