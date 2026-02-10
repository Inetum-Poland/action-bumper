// Copyright (c) 2024 Inetum Poland.

package github

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseEvent_PROpened(t *testing.T) {
	tmpDir := t.TempDir()
	eventFile := filepath.Join(tmpDir, "event.json")

	event := map[string]interface{}{
		"action": "opened",
		"number": 42,
		"pull_request": map[string]interface{}{
			"number": 42,
			"title":  "Test PR",
			"labels": []map[string]interface{}{
				{"name": "bumper:patch"},
			},
			"head": map[string]interface{}{
				"sha": "abc123",
			},
		},
	}

	data, err := json.Marshal(event)
	require.NoError(t, err)
	err = os.WriteFile(eventFile, data, 0o644)
	require.NoError(t, err)

	parsed, err := ParseEvent(eventFile)
	require.NoError(t, err)

	assert.Equal(t, "opened", parsed.Action)
	assert.Equal(t, 42, parsed.Number)
	assert.True(t, parsed.IsPREvent())
	assert.False(t, parsed.IsPushEvent())
	assert.NotNil(t, parsed.PullRequest)
	assert.Len(t, parsed.PullRequest.Labels, 1)
	assert.Equal(t, "bumper:patch", parsed.PullRequest.Labels[0].Name)
}

func TestParseEvent_Push(t *testing.T) {
	tmpDir := t.TempDir()
	eventFile := filepath.Join(tmpDir, "event.json")

	event := map[string]interface{}{
		"after": "def456",
		"head_commit": map[string]interface{}{
			"sha":     "def456",
			"message": "Merge pull request #42",
		},
	}

	data, err := json.Marshal(event)
	require.NoError(t, err)
	err = os.WriteFile(eventFile, data, 0o644)
	require.NoError(t, err)

	parsed, err := ParseEvent(eventFile)
	require.NoError(t, err)

	assert.Equal(t, "def456", parsed.After)
	assert.False(t, parsed.IsPREvent())
	assert.True(t, parsed.IsPushEvent())
	assert.NotNil(t, parsed.HeadCommit)
}

func TestParseEvent_MissingFile(t *testing.T) {
	_, err := ParseEvent("/nonexistent/path/event.json")
	assert.Error(t, err)
}

func TestParseEvent_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	eventFile := filepath.Join(tmpDir, "event.json")

	err := os.WriteFile(eventFile, []byte("{ invalid json }"), 0o644)
	require.NoError(t, err)

	_, err = ParseEvent(eventFile)
	assert.Error(t, err)
}

func TestParseEvent_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	eventFile := filepath.Join(tmpDir, "event.json")

	err := os.WriteFile(eventFile, []byte("{}"), 0o644)
	require.NoError(t, err)

	parsed, err := ParseEvent(eventFile)
	require.NoError(t, err)

	assert.Equal(t, "", parsed.Action)
	assert.False(t, parsed.IsPREvent())
	assert.False(t, parsed.IsPushEvent())
}

func TestEvent_IsPREvent(t *testing.T) {
	tests := []struct {
		name   string
		event  Event
		expect bool
	}{
		{
			name: "opened action with PR",
			event: Event{
				Action:      "opened",
				PullRequest: &PullRequest{Number: 1},
			},
			expect: true,
		},
		{
			name: "labeled action with PR",
			event: Event{
				Action:      "labeled",
				PullRequest: &PullRequest{Number: 1},
			},
			expect: true,
		},
		{
			name: "action without PR",
			event: Event{
				Action: "opened",
			},
			expect: false,
		},
		{
			name: "push event",
			event: Event{
				After: "abc123",
			},
			expect: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expect, tt.event.IsPREvent())
		})
	}
}

func TestEvent_IsPushEvent(t *testing.T) {
	tests := []struct {
		name   string
		event  Event
		expect bool
	}{
		{
			name: "push event with after",
			event: Event{
				After: "abc123",
			},
			expect: true,
		},
		{
			name: "PR event",
			event: Event{
				Action:      "opened",
				PullRequest: &PullRequest{Number: 1},
			},
			expect: false,
		},
		{
			name:   "empty event",
			event:  Event{},
			expect: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expect, tt.event.IsPushEvent())
		})
	}
}

func TestParseEvent_LabelExtraction(t *testing.T) {
	tmpDir := t.TempDir()
	eventFile := filepath.Join(tmpDir, "event.json")

	event := map[string]interface{}{
		"action": "opened",
		"number": 1,
		"pull_request": map[string]interface{}{
			"number": 1,
			"title":  "Test",
			"labels": []map[string]interface{}{
				{"name": "bumper:major"},
				{"name": "feature"},
				{"name": "bumper:patch"},
			},
		},
	}

	data, err := json.Marshal(event)
	require.NoError(t, err)
	err = os.WriteFile(eventFile, data, 0o644)
	require.NoError(t, err)

	parsed, err := ParseEvent(eventFile)
	require.NoError(t, err)

	assert.Len(t, parsed.PullRequest.Labels, 3)
	assert.Equal(t, "bumper:major", parsed.PullRequest.Labels[0].Name)
	assert.Equal(t, "feature", parsed.PullRequest.Labels[1].Name)
	assert.Equal(t, "bumper:patch", parsed.PullRequest.Labels[2].Name)
}
