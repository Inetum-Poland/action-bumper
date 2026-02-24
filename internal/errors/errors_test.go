// Copyright (c) 2024-2026 Inetum Poland.

package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBumperError_Error(t *testing.T) {
	t.Run("without cause", func(t *testing.T) {
		err := New(CodeConfig, "missing token")
		assert.Equal(t, "[CONFIG] missing token", err.Error())
	})

	t.Run("with cause", func(t *testing.T) {
		cause := errors.New("underlying error")
		err := Wrap(cause, CodeGitHub, "API call failed")
		assert.Equal(t, "[GITHUB] API call failed: underlying error", err.Error())
	})
}

func TestBumperError_Unwrap(t *testing.T) {
	cause := errors.New("underlying error")
	err := Wrap(cause, CodeGit, "git operation failed")

	assert.Equal(t, cause, err.Unwrap())
	assert.True(t, errors.Is(err, cause))
}

func TestBumperError_Is(t *testing.T) {
	err1 := New(CodeConfig, "error 1")
	err2 := New(CodeConfig, "error 2")
	err3 := New(CodeGitHub, "error 3")

	assert.True(t, errors.Is(err1, err2))  // Same code
	assert.False(t, errors.Is(err1, err3)) // Different code
}

func TestErrorCreators(t *testing.T) {
	tests := []struct {
		name     string
		creator  func() *BumperError
		expected string
	}{
		{
			name:     "ConfigError",
			creator:  func() *BumperError { return ConfigError("test") },
			expected: "[CONFIG] test",
		},
		{
			name:     "ConfigErrorf",
			creator:  func() *BumperError { return ConfigErrorf("value: %d", 42) },
			expected: "[CONFIG] value: 42",
		},
		{
			name:     "GitHubError",
			creator:  func() *BumperError { return GitHubError("test") },
			expected: "[GITHUB] test",
		},
		{
			name:     "GitError",
			creator:  func() *BumperError { return GitError("test") },
			expected: "[GIT] test",
		},
		{
			name:     "SemverError",
			creator:  func() *BumperError { return SemverError("test") },
			expected: "[SEMVER] test",
		},
		{
			name:     "EventError",
			creator:  func() *BumperError { return EventError("test") },
			expected: "[EVENT] test",
		},
		{
			name:     "BumpError",
			creator:  func() *BumperError { return BumpError("test") },
			expected: "[BUMP] test",
		},
		{
			name:     "OutputError",
			creator:  func() *BumperError { return OutputError("test") },
			expected: "[OUTPUT] test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.creator()
			assert.Equal(t, tt.expected, err.Error())
		})
	}
}

func TestWrapFunctions(t *testing.T) {
	cause := errors.New("underlying")

	tests := []struct {
		name     string
		wrapper  func(error, string) *BumperError
		checker  func(error) bool
		expected string
	}{
		{
			name:     "ConfigWrap",
			wrapper:  ConfigWrap,
			checker:  IsConfigError,
			expected: "[CONFIG] message: underlying",
		},
		{
			name:     "GitHubWrap",
			wrapper:  GitHubWrap,
			checker:  IsGitHubError,
			expected: "[GITHUB] message: underlying",
		},
		{
			name:     "GitWrap",
			wrapper:  GitWrap,
			checker:  IsGitError,
			expected: "[GIT] message: underlying",
		},
		{
			name:     "SemverWrap",
			wrapper:  SemverWrap,
			checker:  IsSemverError,
			expected: "[SEMVER] message: underlying",
		},
		{
			name:     "EventWrap",
			wrapper:  EventWrap,
			checker:  IsEventError,
			expected: "[EVENT] message: underlying",
		},
		{
			name:     "BumpWrap",
			wrapper:  BumpWrap,
			checker:  IsBumpError,
			expected: "[BUMP] message: underlying",
		},
		{
			name:     "OutputWrap",
			wrapper:  OutputWrap,
			checker:  func(e error) bool { var b *BumperError; return errors.As(e, &b) && b.Code == CodeOutput },
			expected: "[OUTPUT] message: underlying",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.wrapper(cause, "message")
			assert.Equal(t, tt.expected, err.Error())
			assert.True(t, tt.checker(err))
			assert.True(t, errors.Is(err, cause))
		})
	}
}

func TestNoLevelError(t *testing.T) {
	err := NewNoLevelError()

	assert.True(t, IsNoLevelError(err))
	assert.True(t, IsBumpError(err))
	assert.Contains(t, err.Error(), "no bump level specified")
}

func TestSkipBumpError(t *testing.T) {
	err := NewSkipBumpError("bumper:none label")

	assert.True(t, IsSkipBumpError(err))
	assert.True(t, IsBumpError(err))
	assert.Equal(t, "bumper:none label", err.Reason)
}

func TestIsCheckers_NilAndOther(t *testing.T) {
	stdErr := errors.New("standard error")

	assert.False(t, IsConfigError(stdErr))
	assert.False(t, IsGitHubError(stdErr))
	assert.False(t, IsGitError(stdErr))
	assert.False(t, IsSemverError(stdErr))
	assert.False(t, IsEventError(stdErr))
	assert.False(t, IsBumpError(stdErr))
	assert.False(t, IsNoLevelError(stdErr))
	assert.False(t, IsSkipBumpError(stdErr))
}
