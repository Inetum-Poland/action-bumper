// Copyright (c) 2024 Inetum Poland.

// Package errors provides custom error types for the bumper application.
package errors

import (
	"errors"
	"fmt"
)

// Error codes for different failure scenarios
const (
	CodeConfig  = "CONFIG"
	CodeGitHub  = "GITHUB"
	CodeGit     = "GIT"
	CodeSemver  = "SEMVER"
	CodeEvent   = "EVENT"
	CodeBump    = "BUMP"
	CodeOutput  = "OUTPUT"
	CodeUnknown = "UNKNOWN"
)

// BumperError is the base error type for all bumper errors
type BumperError struct {
	Code    string
	Message string
	Cause   error
}

// Error implements the error interface
func (e *BumperError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap returns the underlying cause
func (e *BumperError) Unwrap() error {
	return e.Cause
}

// Is checks if the error matches the target
func (e *BumperError) Is(target error) bool {
	t, ok := target.(*BumperError)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

// New creates a new BumperError
func New(code, message string) *BumperError {
	return &BumperError{
		Code:    code,
		Message: message,
	}
}

// Wrap wraps an existing error with additional context
func Wrap(err error, code, message string) *BumperError {
	return &BumperError{
		Code:    code,
		Message: message,
		Cause:   err,
	}
}

// Common error types

// ConfigError creates a configuration error
func ConfigError(message string) *BumperError {
	return New(CodeConfig, message)
}

// ConfigErrorf creates a formatted configuration error
func ConfigErrorf(format string, args ...interface{}) *BumperError {
	return New(CodeConfig, fmt.Sprintf(format, args...))
}

// ConfigWrap wraps an error as a configuration error
func ConfigWrap(err error, message string) *BumperError {
	return Wrap(err, CodeConfig, message)
}

// GitHubError creates a GitHub API error
func GitHubError(message string) *BumperError {
	return New(CodeGitHub, message)
}

// GitHubErrorf creates a formatted GitHub API error
func GitHubErrorf(format string, args ...interface{}) *BumperError {
	return New(CodeGitHub, fmt.Sprintf(format, args...))
}

// GitHubWrap wraps an error as a GitHub API error
func GitHubWrap(err error, message string) *BumperError {
	return Wrap(err, CodeGitHub, message)
}

// GitError creates a git operation error
func GitError(message string) *BumperError {
	return New(CodeGit, message)
}

// GitErrorf creates a formatted git operation error
func GitErrorf(format string, args ...interface{}) *BumperError {
	return New(CodeGit, fmt.Sprintf(format, args...))
}

// GitWrap wraps an error as a git operation error
func GitWrap(err error, message string) *BumperError {
	return Wrap(err, CodeGit, message)
}

// SemverError creates a semver parsing error
func SemverError(message string) *BumperError {
	return New(CodeSemver, message)
}

// SemverErrorf creates a formatted semver parsing error
func SemverErrorf(format string, args ...interface{}) *BumperError {
	return New(CodeSemver, fmt.Sprintf(format, args...))
}

// SemverWrap wraps an error as a semver parsing error
func SemverWrap(err error, message string) *BumperError {
	return Wrap(err, CodeSemver, message)
}

// EventError creates an event parsing error
func EventError(message string) *BumperError {
	return New(CodeEvent, message)
}

// EventErrorf creates a formatted event parsing error
func EventErrorf(format string, args ...interface{}) *BumperError {
	return New(CodeEvent, fmt.Sprintf(format, args...))
}

// EventWrap wraps an error as an event parsing error
func EventWrap(err error, message string) *BumperError {
	return Wrap(err, CodeEvent, message)
}

// BumpError creates a bump operation error
func BumpError(message string) *BumperError {
	return New(CodeBump, message)
}

// BumpErrorf creates a formatted bump operation error
func BumpErrorf(format string, args ...interface{}) *BumperError {
	return New(CodeBump, fmt.Sprintf(format, args...))
}

// BumpWrap wraps an error as a bump operation error
func BumpWrap(err error, message string) *BumperError {
	return Wrap(err, CodeBump, message)
}

// OutputError creates an output writing error
func OutputError(message string) *BumperError {
	return New(CodeOutput, message)
}

// OutputWrap wraps an error as an output writing error
func OutputWrap(err error, message string) *BumperError {
	return Wrap(err, CodeOutput, message)
}

// IsConfigError checks if the error is a configuration error
func IsConfigError(err error) bool {
	var berr *BumperError
	if errors.As(err, &berr) {
		return berr.Code == CodeConfig
	}
	return false
}

// IsGitHubError checks if the error is a GitHub API error
func IsGitHubError(err error) bool {
	var berr *BumperError
	if errors.As(err, &berr) {
		return berr.Code == CodeGitHub
	}
	return false
}

// IsGitError checks if the error is a git operation error
func IsGitError(err error) bool {
	var berr *BumperError
	if errors.As(err, &berr) {
		return berr.Code == CodeGit
	}
	return false
}

// IsSemverError checks if the error is a semver parsing error
func IsSemverError(err error) bool {
	var berr *BumperError
	if errors.As(err, &berr) {
		return berr.Code == CodeSemver
	}
	return false
}

// IsEventError checks if the error is an event parsing error
func IsEventError(err error) bool {
	var berr *BumperError
	if errors.As(err, &berr) {
		return berr.Code == CodeEvent
	}
	return false
}

// IsBumpError checks if the error is a bump operation error
func IsBumpError(err error) bool {
	var berr *BumperError
	if errors.As(err, &berr) {
		return berr.Code == CodeBump
	}
	return false
}

// NoLevelError is returned when no bump level is specified and fail_if_no_level is true
type NoLevelError struct {
	*BumperError
}

// NewNoLevelError creates a new NoLevelError
func NewNoLevelError() *NoLevelError {
	return &NoLevelError{
		BumperError: New(CodeBump, "no bump level specified and fail_if_no_level is enabled"),
	}
}

// Error implements the error interface
func (e *NoLevelError) Error() string {
	return e.BumperError.Error()
}

// Unwrap returns the underlying BumperError
func (e *NoLevelError) Unwrap() error {
	return e.BumperError
}

// IsNoLevelError checks if the error is a NoLevelError
func IsNoLevelError(err error) bool {
	var nerr *NoLevelError
	return errors.As(err, &nerr)
}

// SkipBumpError is returned when bumper:none label is set
type SkipBumpError struct {
	*BumperError
	Reason string
}

// NewSkipBumpError creates a new SkipBumpError
func NewSkipBumpError(reason string) *SkipBumpError {
	return &SkipBumpError{
		BumperError: New(CodeBump, "bump skipped"),
		Reason:      reason,
	}
}

// Error implements the error interface
func (e *SkipBumpError) Error() string {
	return e.BumperError.Error()
}

// Unwrap returns the underlying BumperError
func (e *SkipBumpError) Unwrap() error {
	return e.BumperError
}

// IsSkipBumpError checks if the error is a SkipBumpError
func IsSkipBumpError(err error) bool {
	var serr *SkipBumpError
	return errors.As(err, &serr)
}
