// Copyright (c) 2024-2026 Inetum Poland.

// Package bumper implements the core version bumping logic for the action-bumper
// GitHub Action. It orchestrates the interaction between GitHub API, git operations,
// and version management based on PR labels.
//
// The bumper supports two modes:
//   - PR events (opened, labeled, synchronize): Preview mode that calculates
//     the next version based on labels without creating tags
//   - Push events (merged PRs): Create mode that actually creates and pushes
//     version tags based on the merged PR's labels
//
// Label-based version bumping:
//   - bumper:major - Increments the major version (1.0.0 -> 2.0.0)
//   - bumper:minor - Increments the minor version (1.0.0 -> 1.1.0)
//   - bumper:patch - Increments the patch version (1.0.0 -> 1.0.1)
//   - bumper:none  - Skips version bumping for this PR
package bumper

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Inetum-Poland/action-bumper/internal/config"
	"github.com/Inetum-Poland/action-bumper/internal/git"
	"github.com/Inetum-Poland/action-bumper/internal/github"
	"github.com/Inetum-Poland/action-bumper/internal/output"
	"github.com/Inetum-Poland/action-bumper/internal/semver"
)

// Bumper handles the version bumping logic for GitHub Actions.
// It coordinates between the GitHub API client, git operations,
// and output writer to implement the version bumping workflow.
type Bumper struct {
	cfg    *config.Config
	client github.ClientInterface
	git    git.Operator
	output *output.Writer
	logger *slog.Logger
}

// New creates a new Bumper instance with the given configuration.
// It initializes a GitHub API client for repository operations.
// Returns an error if the GitHub client cannot be created.
func New(ctx context.Context, cfg *config.Config, log *slog.Logger) (*Bumper, error) {
	client, err := github.NewClient(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create GitHub client: %w", err)
	}

	return &Bumper{
		cfg:    cfg,
		client: client,
		git:    git.NewOperator(),
		output: output.NewWriter(),
		logger: log,
	}, nil
}

// NewWithClient creates a new Bumper instance with a custom GitHub client.
// This is useful for testing with mock clients.
func NewWithClient(cfg *config.Config, client github.ClientInterface, log *slog.Logger) *Bumper {
	return &Bumper{
		cfg:    cfg,
		client: client,
		git:    git.NewOperator(),
		output: output.NewWriter(),
		logger: log,
	}
}

// NewWithClientAndGit creates a new Bumper instance with custom GitHub and Git clients.
// This is useful for testing with mock clients.
func NewWithClientAndGit(cfg *config.Config, client github.ClientInterface, gitOp git.Operator, log *slog.Logger) *Bumper {
	return &Bumper{
		cfg:    cfg,
		client: client,
		git:    gitOp,
		output: output.NewWriter(),
		logger: log,
	}
}

// trace logs at trace level (more verbose than debug)
// This is a custom level below Debug (-4)
func (b *Bumper) trace(ctx context.Context, msg string, args ...any) {
	if b.cfg.Trace {
		// Log at level -4 (below Debug which is -4)
		b.logger.Log(ctx, slog.LevelDebug-4, msg, args...)
	}
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
	b.logger.Info("Handling PR event", "action", event.Action, "pr", event.Number)
	b.trace(ctx, "PR event details", "labels", len(event.PullRequest.Labels), "title", event.PullRequest.Title)

	// Get bump level from labels
	bumpLevel := b.determineBumpLevel(event.PullRequest.Labels)
	b.trace(ctx, "Determined bump level from labels", "level", bumpLevel)

	if err := b.requireExplicitBumpLabel(bumpLevel, nil); err != nil {
		return err
	}

	rawBumpLevel := bumpLevel
	bumpLevel = b.applyDefaultBumpLevel(bumpLevel)
	if rawBumpLevel == config.BumpLevelEmpty && bumpLevel != config.BumpLevelEmpty {
		b.trace(ctx, "Using default level", "level", bumpLevel)
	}

	// If still empty or none, skip
	if bumpLevel == config.BumpLevelEmpty || bumpLevel == config.BumpLevelNone {
		if err := b.output.SetAll(map[string]string{
			"skip":    "true",
			"message": "No version bump (bumper:none or no label)",
		}); err != nil {
			return fmt.Errorf("failed to set outputs: %w", err)
		}
		b.logger.Info("Skipping bump: no level specified or bumper:none")
		return nil
	}

	// Get current version
	currentVersion, err := b.client.GetLatestTag(ctx)
	if err != nil {
		b.logger.Warn("Failed to get latest tag", "error", err)
		// Use default version
		currentVersion = semver.DefaultVersion(bumpLevel)
	} else if currentVersion == nil {
		// No tags exist yet - use default version
		b.logger.Info("No existing tags found, using default version")
		currentVersion = semver.DefaultVersion(bumpLevel)
	}

	// Calculate next version
	nextVersion := currentVersion.Bump(bumpLevel)

	// Generate status message
	message := b.generatePRStatusMessage(currentVersion, nextVersion, bumpLevel)

	// Set outputs
	// Note: Bash always outputs next_version with "v" prefix regardless of BumpIncludeV setting
	// BumpIncludeV only affects actual tag creation in push events
	outputs := map[string]string{
		"current_version": currentVersion.FullTag(true), // Always include v for PR events (matches Bash)
		"next_version":    nextVersion.FullTag(true),    // Always include v for PR events (matches Bash)
		"skip":            "false",
		"message":         fmt.Sprintf("Will bump version %s → %s (%s)", currentVersion.FullTag(true), nextVersion.FullTag(true), bumpLevel),
	}

	if err := b.output.SetAll(outputs); err != nil {
		return fmt.Errorf("failed to set outputs: %w", err)
	}

	if err := b.output.SetMultiline("tag_status", message); err != nil {
		return fmt.Errorf("failed to set tag_status: %w", err)
	}

	b.logger.Info("Version bump preview", "current", currentVersion.String(), "next", nextVersion.String(), "level", bumpLevel)
	return nil
}

// handlePushEvent handles push events (actually create tags)
func (b *Bumper) handlePushEvent(ctx context.Context, event *github.Event) error {
	b.logger.Info("Handling push event")
	b.trace(ctx, "Push event details", "after", event.After)

	// For push events, we need to find the merged PR to get labels
	if event.HeadCommit == nil {
		return fmt.Errorf("no head commit in push event")
	}

	// Get the merge commit SHA to find the associated PR
	mergeCommitSHA := b.cfg.GitHubSHA
	if event.After != "" {
		mergeCommitSHA = event.After
	}
	b.trace(ctx, "Looking for merged PR", "sha", mergeCommitSHA)

	bumpLevel, prNumber, prTitle, err := b.resolvePushBumpLevel(ctx, mergeCommitSHA)
	if err != nil {
		return err
	}

	skipped, err := b.handlePushSkip(bumpLevel)
	if err != nil {
		return err
	}
	if skipped {
		return nil
	}

	currentVersion, nextVersion, tags, tagMessage, err := b.createAndPushTags(ctx, bumpLevel, prNumber, prTitle)
	if err != nil {
		return err
	}

	// Generate status message
	message := b.generatePushStatusMessage(nextVersion, tags, prNumber, prTitle)

	// Set outputs
	outputs := map[string]string{
		"current_version": currentVersion.FullTag(b.cfg.BumpIncludeV),
		"next_version":    nextVersion.FullTag(b.cfg.BumpIncludeV),
		"skip":            "false",
		"message":         tagMessage,
	}

	if err := b.output.SetAll(outputs); err != nil {
		return fmt.Errorf("failed to set outputs: %w", err)
	}

	if err := b.output.SetMultiline("tag_status", message); err != nil {
		return fmt.Errorf("failed to set tag_status: %w", err)
	}

	b.logger.Info("Successfully created and pushed tags", "tags", tags)
	return nil
}

func (b *Bumper) resolvePushBumpLevel(ctx context.Context, mergeCommitSHA string) (bumpLevel config.BumpLevel, prNumber int, prTitle string, err error) {
	mergedPR, err := b.client.GetMergedPRByCommitSHA(ctx, mergeCommitSHA)
	if err != nil {
		if guardErr := b.requireExplicitBumpLabel(config.BumpLevelEmpty, err); guardErr != nil {
			return config.BumpLevelEmpty, 0, "", guardErr
		}
		b.logger.Warn("Failed to find merged PR, using default level", "error", err, "sha", mergeCommitSHA)
		return b.cfg.BumpDefaultLevel, 0, "", nil
	}

	b.logger.Info("Found merged PR", "number", mergedPR.Number, "title", mergedPR.Title)
	b.trace(ctx, "PR labels", "count", len(mergedPR.Labels))

	bumpLevel = b.determineBumpLevel(mergedPR.Labels)
	b.trace(ctx, "Determined bump level from PR", "level", bumpLevel)
	if guardErr := b.requireExplicitBumpLabel(bumpLevel, nil); guardErr != nil {
		return config.BumpLevelEmpty, 0, "", guardErr
	}

	rawBumpLevel := bumpLevel
	bumpLevel = b.applyDefaultBumpLevel(bumpLevel)
	if rawBumpLevel == config.BumpLevelEmpty && bumpLevel != config.BumpLevelEmpty {
		b.trace(ctx, "Using default level", "level", bumpLevel)
	}

	return bumpLevel, mergedPR.Number, mergedPR.Title, nil
}

func (b *Bumper) requireExplicitBumpLabel(bumpLevel config.BumpLevel, cause error) error {
	if !b.cfg.BumpFailIfNoLevel || bumpLevel != config.BumpLevelEmpty {
		return nil
	}

	fmt.Println("::error ::Job failed as no bump label is found.")
	if cause != nil {
		return fmt.Errorf("no bump level label found and bump_fail_if_no_level is true: %w", cause)
	}

	return fmt.Errorf("no bump level label found and bump_fail_if_no_level is true")
}

func (b *Bumper) applyDefaultBumpLevel(bumpLevel config.BumpLevel) config.BumpLevel {
	if bumpLevel != config.BumpLevelEmpty {
		return bumpLevel
	}
	return b.cfg.BumpDefaultLevel
}

func (b *Bumper) handlePushSkip(bumpLevel config.BumpLevel) (bool, error) {
	switch bumpLevel {
	case config.BumpLevelEmpty:
		b.logger.Info("No bump level, skipping")
	case config.BumpLevelNone:
		b.logger.Info("Bump level is 'none', skipping")
	case config.BumpLevelMajor, config.BumpLevelMinor, config.BumpLevelPatch:
		return false, nil
	default:
		return false, nil
	}

	if err := b.output.Set("skip", "true"); err != nil {
		return false, fmt.Errorf("failed to set skip output: %w", err)
	}
	return true, nil
}

func (b *Bumper) createAndPushTags(ctx context.Context, bumpLevel config.BumpLevel, prNumber int, prTitle string) (currentVersion, nextVersion *semver.Version, tags []string, tagMessage string, err error) {
	currentVersion, err = b.client.GetLatestTag(ctx)
	if err != nil {
		b.logger.Warn("No existing tags, using default")
		currentVersion = semver.DefaultVersion(bumpLevel)
	} else if currentVersion == nil {
		b.logger.Info("No existing tags found, using default version")
		currentVersion = semver.DefaultVersion(bumpLevel)
	}

	nextVersion = currentVersion.Bump(bumpLevel)

	if err = b.git.ConfigureUser(b.cfg.BumpTagAsUser, b.cfg.BumpTagAsEmail); err != nil {
		return nil, nil, nil, "", fmt.Errorf("failed to configure git user: %w", err)
	}

	if err = b.git.SetRemoteURL(b.cfg.GitHubToken, b.cfg.GitHubRepo); err != nil {
		return nil, nil, nil, "", fmt.Errorf("failed to set remote URL: %w", err)
	}

	tagMessage = b.pushTagMessage(nextVersion, prNumber, prTitle)
	tags, err = b.createVersionTags(nextVersion, tagMessage)
	if err != nil {
		return nil, nil, nil, "", err
	}

	if err := b.git.PushTags(tags); err != nil {
		return nil, nil, nil, "", fmt.Errorf("failed to push tags: %w", err)
	}

	return currentVersion, nextVersion, tags, tagMessage, nil
}

func (b *Bumper) pushTagMessage(version *semver.Version, prNumber int, prTitle string) string {
	if prNumber > 0 {
		return fmt.Sprintf("%s: PR #%d - %s", version.FullTag(b.cfg.BumpIncludeV), prNumber, prTitle)
	}
	return fmt.Sprintf("Release %s", version.FullTag(b.cfg.BumpIncludeV))
}

func (b *Bumper) createVersionTags(version *semver.Version, tagMessage string) ([]string, error) {
	mainTag := version.FullTag(b.cfg.BumpIncludeV)
	tags := []string{mainTag}

	if err := b.git.CreateTag(mainTag, tagMessage); err != nil {
		return nil, fmt.Errorf("failed to create tag: %w", err)
	}

	refSpec := fmt.Sprintf("%s^{commit}", mainTag)
	if b.cfg.BumpSemver {
		majorTag := version.MajorTag(b.cfg.BumpIncludeV)
		minorTag := version.MinorTag(b.cfg.BumpIncludeV)

		if err := b.git.CreateOrUpdateTag(majorTag, tagMessage, refSpec); err != nil {
			b.logger.Warn("Failed to create major tag", "error", err)
		} else {
			tags = append(tags, majorTag)
		}

		if err := b.git.CreateOrUpdateTag(minorTag, tagMessage, refSpec); err != nil {
			b.logger.Warn("Failed to create minor tag", "error", err)
		} else {
			tags = append(tags, minorTag)
		}
	}

	if b.cfg.BumpLatest {
		latestTag := "latest"
		if err := b.git.CreateOrUpdateTag(latestTag, tagMessage, refSpec); err != nil {
			b.logger.Warn("Failed to create latest tag", "error", err)
		} else {
			tags = append(tags, latestTag)
		}
	}

	return tags, nil
}

// determineBumpLevel determines the bump level from PR labels.
// Matches Bash behavior: processes in order none → patch → minor → major,
// with later matches taking priority.
func (b *Bumper) determineBumpLevel(labels []github.Label) config.BumpLevel {
	level := config.BumpLevelEmpty

	// Check labels in priority order (lowest to highest)
	// Each match overwrites the previous, so highest priority wins
	for _, label := range labels {
		if label.Name == b.cfg.LabelNone {
			level = config.BumpLevelNone
		}
	}
	for _, label := range labels {
		if label.Name == b.cfg.LabelPatch {
			level = config.BumpLevelPatch
		}
	}
	for _, label := range labels {
		if label.Name == b.cfg.LabelMinor {
			level = config.BumpLevelMinor
		}
	}
	for _, label := range labels {
		if label.Name == b.cfg.LabelMajor {
			level = config.BumpLevelMajor
		}
	}

	return level
}

// generatePRStatusMessage generates a status message for PR events (matching Bash format)
func (b *Bumper) generatePRStatusMessage(current, next *semver.Version, _ config.BumpLevel) string {
	var sb strings.Builder

	// Build additional info for semver/latest tags
	additionalInfo := ""
	if b.cfg.BumpSemver {
		additionalInfo += fmt.Sprintf(" / %s / %s", next.MinorTag(b.cfg.BumpIncludeV), next.MajorTag(b.cfg.BumpIncludeV))
	}
	if b.cfg.BumpLatest {
		additionalInfo += " / latest"
	}

	// Format matching Bash: 🏷️ [[bumper]](url) @ action<br>**Next version**: version<br>**Changes**: compare_link
	sb.WriteString("🏷️ [[bumper]](https://github.com/inetum-poland/action-bumper) @ preview<br>")
	sb.WriteString(fmt.Sprintf("**Next version**: %s%s<br>", next.FullTag(b.cfg.BumpIncludeV), additionalInfo))

	if current != nil {
		sb.WriteString(fmt.Sprintf("**Changes**: [%s...HEAD](https://%s/%s/compare/%s...HEAD)",
			current.FullTag(b.cfg.BumpIncludeV),
			"github.com",
			b.cfg.GitHubRepo,
			current.FullTag(b.cfg.BumpIncludeV),
		))
	}

	return sb.String()
}

// generatePushStatusMessage generates a status message for push events (matching Bash format)
func (b *Bumper) generatePushStatusMessage(version *semver.Version, _ []string, _ int, _ string) string {
	var sb strings.Builder

	// Build additional info for semver/latest tags
	additionalInfo := ""
	if b.cfg.BumpSemver {
		additionalInfo += fmt.Sprintf(" / %s / %s", version.MinorTag(b.cfg.BumpIncludeV), version.MajorTag(b.cfg.BumpIncludeV))
	}
	if b.cfg.BumpLatest {
		additionalInfo += " / latest"
	}

	// Format matching Bash: 🚀 [[bumper]](url) [Bumped!](run_url)<br>**New version**: [version](release_url)<br>**Changes**: compare_link
	sb.WriteString(fmt.Sprintf("🚀 [[bumper]](https://github.com/inetum-poland/action-bumper) [Bumped!](https://github.com/%s/actions)<br>",
		b.cfg.GitHubRepo,
	))
	sb.WriteString(fmt.Sprintf("**New version**: [%s%s](https://github.com/%s/releases/tag/%s)<br>",
		version.FullTag(b.cfg.BumpIncludeV),
		additionalInfo,
		b.cfg.GitHubRepo,
		version.FullTag(b.cfg.BumpIncludeV),
	))

	return sb.String()
}
