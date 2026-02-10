// Copyright (c) 2024 Inetum Poland.

package output

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriter_Set(t *testing.T) {
	// Create temporary output file
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.txt")

	tests := []struct {
		name  string
		key   string
		value string
	}{
		{"simple value", "version", "1.2.3"},
		{"with special chars", "message", "test: value"},
		{"empty value", "empty", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment
			os.Setenv("GITHUB_OUTPUT", outputFile)
			defer os.Unsetenv("GITHUB_OUTPUT")

			w := NewWriter()
			err := w.Set(tt.key, tt.value)
			require.NoError(t, err)

			// Read file and verify
			content, err := os.ReadFile(outputFile)
			require.NoError(t, err)

			expected := tt.key + "=" + tt.value + "\n"
			assert.Contains(t, string(content), expected)
		})
	}
}

func TestWriter_SetMultiline(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.txt")

	os.Setenv("GITHUB_OUTPUT", outputFile)
	defer os.Unsetenv("GITHUB_OUTPUT")

	w := NewWriter()
	multilineValue := "line1\nline2\nline3"

	err := w.SetMultiline("message", multilineValue)
	require.NoError(t, err)

	content, err := os.ReadFile(outputFile)
	require.NoError(t, err)

	// Verify format: key<<EOF\nvalue\nEOF\n
	lines := strings.Split(string(content), "\n")
	assert.Equal(t, "message<<EOF", lines[0])
	assert.Contains(t, string(content), multilineValue)
	assert.Contains(t, string(content), "EOF")
}

func TestWriter_SetAll(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.txt")

	os.Setenv("GITHUB_OUTPUT", outputFile)
	defer os.Unsetenv("GITHUB_OUTPUT")

	w := NewWriter()
	outputs := map[string]string{
		"version": "1.2.3",
		"skip":    "false",
		"message": "test message",
	}

	err := w.SetAll(outputs)
	require.NoError(t, err)

	content, err := os.ReadFile(outputFile)
	require.NoError(t, err)

	for key, value := range outputs {
		expected := key + "=" + value
		assert.Contains(t, string(content), expected)
	}
}

func TestWriter_NoGitHubOutput(t *testing.T) {
	// When GITHUB_OUTPUT is not set, should not error
	os.Unsetenv("GITHUB_OUTPUT")

	w := NewWriter()
	err := w.Set("test", "value")
	assert.NoError(t, err)

	err = w.SetMultiline("test", "multiline\nvalue")
	assert.NoError(t, err)
}

func TestWriter_SpecialCharacters(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.txt")

	os.Setenv("GITHUB_OUTPUT", outputFile)
	defer os.Unsetenv("GITHUB_OUTPUT")

	tests := []struct {
		name  string
		key   string
		value string
	}{
		{"unicode", "message", "🏷️ Release v1.0.0"},
		{"quotes", "desc", `He said "hello"`},
		{"backslash", "path", `C:\path\to\file`},
		{"equals sign", "equation", "a=b+c"},
		{"newline escaped", "single", "line1\\nline2"},
		{"html-like", "html", "<b>bold</b>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear file
			os.WriteFile(outputFile, []byte{}, 0o644)

			w := NewWriter()
			err := w.Set(tt.key, tt.value)
			require.NoError(t, err)

			content, err := os.ReadFile(outputFile)
			require.NoError(t, err)

			expected := tt.key + "=" + tt.value + "\n"
			assert.Equal(t, expected, string(content))
		})
	}
}

func TestWriter_LongMultilineValue(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.txt")

	os.Setenv("GITHUB_OUTPUT", outputFile)
	defer os.Unsetenv("GITHUB_OUTPUT")

	// Create a very long multiline value (1000 lines)
	var builder strings.Builder
	for i := 0; i < 1000; i++ {
		builder.WriteString(fmt.Sprintf("Line %d: This is test content for line number %d\n", i, i))
	}
	longValue := builder.String()

	w := NewWriter()
	err := w.SetMultiline("long_output", longValue)
	require.NoError(t, err)

	content, err := os.ReadFile(outputFile)
	require.NoError(t, err)

	// Verify the multiline format
	assert.True(t, strings.HasPrefix(string(content), "long_output<<EOF\n"))
	assert.True(t, strings.HasSuffix(string(content), "\nEOF\n"))
	assert.Contains(t, string(content), "Line 999:")
}

func TestWriter_ConcurrentWrites(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.txt")

	os.Setenv("GITHUB_OUTPUT", outputFile)
	defer os.Unsetenv("GITHUB_OUTPUT")

	// Write multiple outputs concurrently
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(idx int) {
			w := NewWriter()
			key := fmt.Sprintf("key_%d", idx)
			value := fmt.Sprintf("value_%d", idx)
			_ = w.Set(key, value)
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	content, err := os.ReadFile(outputFile)
	require.NoError(t, err)

	// Verify all keys are present (order may vary)
	for i := 0; i < 10; i++ {
		assert.Contains(t, string(content), fmt.Sprintf("key_%d=value_%d", i, i))
	}
}

func TestWriter_WriteSummary(t *testing.T) {
	tmpDir := t.TempDir()
	summaryFile := filepath.Join(tmpDir, "summary.md")

	os.Setenv("GITHUB_STEP_SUMMARY", summaryFile)
	defer os.Unsetenv("GITHUB_STEP_SUMMARY")

	w := NewWriter()

	markdown := `## Version Bump Summary

| Field | Value |
|-------|-------|
| Current | v1.0.0 |
| Next | v1.1.0 |
| Level | minor |
`

	err := w.WriteSummary(markdown)
	require.NoError(t, err)

	content, err := os.ReadFile(summaryFile)
	require.NoError(t, err)

	assert.Contains(t, string(content), "Version Bump Summary")
	assert.Contains(t, string(content), "v1.1.0")
}

func TestWriter_WriteSummary_NoFile(t *testing.T) {
	// When GITHUB_STEP_SUMMARY is not set, should not error
	os.Unsetenv("GITHUB_STEP_SUMMARY")

	w := NewWriter()
	err := w.WriteSummary("# Test Summary")
	assert.NoError(t, err)
}
