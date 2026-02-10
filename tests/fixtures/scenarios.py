# Copyright (c) 2024 Inetum Poland.
"""
Test data converted from existing spec/bumper test fixtures.
"""

from pathlib import Path
import json

# Path to original spec fixtures
SPEC_FIXTURES_PATH = Path(__file__).parent.parent.parent / "spec" / "bumper"


def load_spec_fixture(fixture_name: str) -> dict:
    """
    Load a test fixture from the spec/bumper directory.
    
    Args:
        fixture_name: Name of the fixture directory
    
    Returns:
        Dictionary with 'data', 'tags', and optionally 'pull_request' keys
    """
    fixture_path = SPEC_FIXTURES_PATH / fixture_name
    result = {}
    
    data_file = fixture_path / "data.json"
    if data_file.exists():
        result["data"] = json.loads(data_file.read_text())
    
    tags_file = fixture_path / "tags.json"
    if tags_file.exists():
        result["tags"] = json.loads(tags_file.read_text())
    
    pr_file = fixture_path / "pull_request.json"
    if pr_file.exists():
        result["pull_request"] = json.loads(pr_file.read_text())
    
    assert_file = fixture_path / "assert.txt"
    if assert_file.exists():
        result["assert"] = assert_file.read_text().strip()
    
    return result


# Pre-defined test scenarios based on existing spec fixtures
TEST_SCENARIOS = {
    # PR Events
    "opened_event_bumper_patch": {
        "description": "PR opened with bumper:patch label",
        "event_type": "pull_request",
        "action": "opened",
        "labels": ["bumper:patch"],
        "existing_tags": ["v0.9.1", "v0.9.0", "v0.8.4"],
        "expected_current": "v0.9.1",
        "expected_next": "v0.9.2",
    },
    "opened_event_bumper_minor": {
        "description": "PR opened with bumper:minor label",
        "event_type": "pull_request",
        "action": "opened",
        "labels": ["bumper:minor"],
        "existing_tags": ["v0.9.1", "v0.9.0", "v0.8.4"],
        "expected_current": "v0.9.1",
        "expected_next": "v0.10.0",
    },
    "opened_event_bumper_major": {
        "description": "PR opened with bumper:major label",
        "event_type": "pull_request",
        "action": "opened",
        "labels": ["bumper:major"],
        "existing_tags": ["v0.9.1", "v0.9.0", "v0.8.4"],
        "expected_current": "v0.9.1",
        "expected_next": "v1.0.0",
    },
    "opened_event_bumper_none": {
        "description": "PR opened with bumper:none label",
        "event_type": "pull_request",
        "action": "opened",
        "labels": ["bumper:none"],
        "existing_tags": ["v0.9.1", "v0.9.0", "v0.8.4"],
        "expected_skip": True,
    },
    "opened_event_without_tags": {
        "description": "PR opened with patch label but no existing tags",
        "event_type": "pull_request",
        "action": "opened",
        "labels": ["bumper:patch"],
        "existing_tags": [],
        "expected_next": "v0.0.1",
    },
    "opened_event_without_tags_without_labels": {
        "description": "PR opened with no labels and no existing tags",
        "event_type": "pull_request",
        "action": "opened",
        "labels": [],
        "existing_tags": [],
        "default_level": "",
        "expected_skip": True,
    },
    "opened_event_without_tags_without_labels_allow": {
        "description": "PR opened with no labels but default level set",
        "event_type": "pull_request",
        "action": "opened",
        "labels": [],
        "existing_tags": [],
        "default_level": "patch",
        "expected_next": "v0.0.1",
    },
    "opened_event_bumper_major_semver": {
        "description": "PR opened with major label and semver enabled",
        "event_type": "pull_request",
        "action": "opened",
        "labels": ["bumper:major"],
        "existing_tags": ["v0.9.1"],
        "bump_semver": True,
        "expected_next": "v1.0.0",
    },
    "opened_event_bumper_major_latest": {
        "description": "PR opened with major label and latest enabled",
        "event_type": "pull_request",
        "action": "opened",
        "labels": ["bumper:major"],
        "existing_tags": ["v0.9.1"],
        "bump_latest": True,
        "expected_next": "v1.0.0",
    },
    # Push Events
    "push_event_with_labels": {
        "description": "Push event with merged PR having bumper:patch label",
        "event_type": "push",
        "merge_commit_sha": "54fa23aef40b58c8f22350c830f7a89dad0121bc",
        "pr_labels": ["feature", "bumper:patch", "doc"],
        "pr_title": "feat(gha): align the gh actions before publish",
        "existing_tags": ["v0.9.1", "v0.9.0", "v0.8.4"],
        "expected_current": "v0.9.1",
        "expected_next": "v0.9.2",
    },
    "push_event_with_labels_semver": {
        "description": "Push event with semver enabled",
        "event_type": "push",
        "merge_commit_sha": "54fa23aef40b58c8f22350c830f7a89dad0121bc",
        "pr_labels": ["bumper:patch"],
        "bump_semver": True,
        "existing_tags": ["v0.9.1"],
        "expected_next": "v0.9.2",
    },
    "push_event_with_labels_without_v": {
        "description": "Push event without v prefix",
        "event_type": "push",
        "merge_commit_sha": "54fa23aef40b58c8f22350c830f7a89dad0121bc",
        "pr_labels": ["bumper:patch"],
        "bump_include_v": False,
        "existing_tags": ["0.9.1"],
        "expected_next": "0.9.2",
    },
    "push_event_without_labels": {
        "description": "Push event with merged PR having no bumper labels",
        "event_type": "push",
        "merge_commit_sha": "54fa23aef40b58c8f22350c830f7a89dad0121bc",
        "pr_labels": ["feature"],
        "existing_tags": ["v0.9.1"],
        "expected_skip": True,
    },
    "synchronize_event": {
        "description": "PR synchronize event",
        "event_type": "pull_request",
        "action": "synchronize",
        "labels": ["bumper:patch"],
        "existing_tags": ["v0.9.1"],
        "expected_next": "v0.9.2",
    },
}


def get_test_scenario(name: str) -> dict:
    """Get a pre-defined test scenario by name."""
    return TEST_SCENARIOS.get(name, {})


def list_test_scenarios() -> list[str]:
    """List all available test scenario names."""
    return list(TEST_SCENARIOS.keys())
