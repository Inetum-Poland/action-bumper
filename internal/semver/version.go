// Copyright (c) 2024 Inetum Poland.

package semver

import (
	"fmt"
	"strings"

	"github.com/Inetum-Poland/action-bumper/internal/config"
	"github.com/Masterminds/semver/v3"
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

// DefaultVersion returns the default starting version for a given bump level
func DefaultVersion(level config.BumpLevel) *Version {
	switch level {
	case config.BumpLevelMajor:
		return MustParse("1.0.0")
	case config.BumpLevelMinor:
		return MustParse("0.1.0")
	case config.BumpLevelPatch:
		return MustParse("0.0.1")
	default:
		return MustParse("0.0.0")
	}
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
	case config.BumpLevelNone:
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
