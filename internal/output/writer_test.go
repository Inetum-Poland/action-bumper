// Copyright (c) 2024 Inetum Poland.

package output

import (
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
