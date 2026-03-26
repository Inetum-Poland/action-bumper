// Copyright (c) 2024-2026 Inetum Poland.

// Package semver provides semantic versioning utilities for parsing, comparing,
// and bumping version numbers according to the Semantic Versioning 2.0.0 specification.
//
// Semantic versions have the format: MAJOR.MINOR.PATCH[-PRERELEASE][+BUILD]
//   - MAJOR: Incremented for incompatible API changes
//   - MINOR: Incremented for backward-compatible functionality additions
//   - PATCH: Incremented for backward-compatible bug fixes
//
// Example usage:
//
//	v, err := semver.Parse("v1.2.3")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	next := v.Bump("minor")  // Returns 1.3.0
//	tag := next.Tag(true)    // Returns "v1.3.0"
package semver

import (
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"

	"github.com/Inetum-Poland/action-bumper/internal/config"
)

// Version wraps semver.Version with additional functionality
type Version struct {
	v *semver.Version
}

// Parse parses a version string (with or without 'v' prefix)
func Parse(s string) (*Version, error) {
	// Remove 'v' prefix if present
	s = strings.TrimPrefix(s, "v")

	v, err := semver.NewVersion(s)
	if err != nil {
		return nil, fmt.Errorf("failed to parse version %q: %w", s, err)
	}

	return &Version{v: v}, nil
}

// MustParse parses a version string and panics on error
func MustParse(s string) *Version {
	v, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return v
}

// DefaultVersion returns the starting version (0.0.0) to be bumped for first release.
// The caller should call Bump() on this to get the actual first version.
func DefaultVersion() *Version {
	// Always return 0.0.0 - the Bump() call will produce:
	// - major: 1.0.0
	// - minor: 0.1.0
	// - patch: 0.0.1
	return MustParse("0.0.0")
}

// Bump returns a new version with the specified level bumped
func (v *Version) Bump(level config.BumpLevel) *Version {
	current := v.v

	var next semver.Version
	switch level {
	case config.BumpLevelMajor:
		next = current.IncMajor()
	case config.BumpLevelMinor:
		next = current.IncMinor()
	case config.BumpLevelPatch:
		next = current.IncPatch()
	case config.BumpLevelNone, config.BumpLevelEmpty:
		// No bump, return current version
		return &Version{v: current}
	default:
		// Unknown level, return current version
		return &Version{v: current}
	}

	return &Version{v: &next}
}

// String returns the version string without 'v' prefix
func (v *Version) String() string {
	return v.v.String()
}

// StringWithV returns the version string with 'v' prefix
func (v *Version) StringWithV() string {
	return "v" + v.v.String()
}

// Major returns the major version number
func (v *Version) Major() uint64 {
	return v.v.Major()
}

// Minor returns the minor version number
func (v *Version) Minor() uint64 {
	return v.v.Minor()
}

// Patch returns the patch version number
func (v *Version) Patch() uint64 {
	return v.v.Patch()
}

// MajorTag returns the major version tag (e.g., "v1")
func (v *Version) MajorTag(includeV bool) string {
	if includeV {
		return fmt.Sprintf("v%d", v.Major())
	}
	return fmt.Sprintf("%d", v.Major())
}

// MinorTag returns the major.minor version tag (e.g., "v1.2")
func (v *Version) MinorTag(includeV bool) string {
	if includeV {
		return fmt.Sprintf("v%d.%d", v.Major(), v.Minor())
	}
	return fmt.Sprintf("%d.%d", v.Major(), v.Minor())
}

// FullTag returns the full version tag with optional 'v' prefix
func (v *Version) FullTag(includeV bool) string {
	if includeV {
		return v.StringWithV()
	}
	return v.String()
}

// Compare compares two versions
// Returns -1 if v < other, 0 if v == other, 1 if v > other
func (v *Version) Compare(other *Version) int {
	return v.v.Compare(other.v)
}

// Equal checks if two versions are equal
func (v *Version) Equal(other *Version) bool {
	return v.v.Equal(other.v)
}

// LessThan checks if v < other
func (v *Version) LessThan(other *Version) bool {
	return v.v.LessThan(other.v)
}

// GreaterThan checks if v > other
func (v *Version) GreaterThan(other *Version) bool {
	return v.v.GreaterThan(other.v)
}
