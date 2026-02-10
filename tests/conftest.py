# Copyright (c) 2024 Inetum Poland.
"""
Pytest configuration and shared fixtures for action-bumper tests.
"""

import pytest
import os
import json
import tempfile
import shutil
from pathlib import Path

from tests.helpers.runner import Implementation

# Root directory of the project
PROJECT_ROOT = Path(__file__).parent.parent.absolute()

# Paths to implementations
BASH_BUMPER = PROJECT_ROOT / "bumper.sh"
GO_BUMPER = PROJECT_ROOT / "bin" / "bumper"


def pytest_addoption(parser):
    """Add command line option for selecting implementation."""
    parser.addoption(
        "--impl",
        action="store",
        default=os.environ.get("BUMPER_IMPL", "bash"),
        choices=["bash", "go"],
        help="Implementation to test: bash or go (default: bash, or BUMPER_IMPL env var)",
    )


@pytest.fixture
def implementation(request) -> Implementation:
    """Return the implementation to test based on --impl option."""
    impl = request.config.getoption("--impl")
    return Implementation.GO if impl == "go" else Implementation.BASH


@pytest.fixture
def project_root() -> Path:
    """Return the project root directory."""
    return PROJECT_ROOT


@pytest.fixture
def debug_event_path(tmp_path: Path) -> Path:
    """
    Create a debug event path directory for Bash bumper testing.
    
    The Bash bumper supports DEBUG_GITHUB_EVENT_PATH which reads
    from local files instead of GitHub API:
    - data.json: The GitHub event payload
    - tags.json: The repository tags (simulates API response)
    - pull_request.json: PR data for push events
    """
    debug_path = tmp_path / "debug_event"
    debug_path.mkdir()
    return debug_path


def create_tags_json(tags: list[str], debug_path: Path) -> None:
    """
    Create a tags.json file and .input.env for DEBUG_GITHUB_EVENT_PATH.
    
    Args:
        tags: List of tag names (e.g., ["v1.0.0", "v0.9.0"])
        debug_path: Path to the debug event directory
    """
    tags_data = [
        {
            "name": tag,
            "commit": {"sha": f"sha-for-{tag}"}
        }
        for tag in tags
    ]
    (debug_path / "tags.json").write_text(json.dumps(tags_data, indent=2))
    
    # Create empty .input.env file (required by debug.sh)
    # The actual env vars are passed via the test environment
    (debug_path / ".input.env").write_text("# Empty - env vars passed via test\n")


def setup_test_environment(
    implementation: "Implementation",
    tags: list[str],
    event: dict,
    base_env: dict,
    debug_event_path: Path,
    temp_workspace: Path,
    github_output_file: Path,
) -> dict:
    """
    Set up test environment for both Bash and Go implementations.
    
    For Bash: Uses DEBUG_GITHUB_EVENT_PATH with tags.json and data.json
    For Go: Uses GITHUB_EVENT_PATH and creates git tags in temp_workspace
    
    Args:
        implementation: Which implementation is being tested
        tags: List of tag names to create
        event: The GitHub event payload dict
        base_env: Base environment variables
        debug_event_path: Path for Bash debug mode
        temp_workspace: Git workspace for Go
        github_output_file: Path to output file
    
    Returns:
        Environment dict configured for the implementation
    """
    from tests.helpers.runner import Implementation
    from tests.helpers.git_helpers import create_git_tag
    
    env = base_env.copy()
    env["GITHUB_OUTPUT"] = str(github_output_file)
    
    if implementation == Implementation.BASH:
        # Bash uses DEBUG_GITHUB_EVENT_PATH mechanism
        create_tags_json(tags, debug_event_path)
        (debug_event_path / "data.json").write_text(json.dumps(event))
        env["DEBUG_GITHUB_EVENT_PATH"] = str(debug_event_path)
    else:
        # Go uses GITHUB_EVENT_PATH and local git tags
        event_file = debug_event_path / "event.json"
        event_file.write_text(json.dumps(event))
        env["GITHUB_EVENT_PATH"] = str(event_file)
        env["GITHUB_WORKSPACE"] = str(temp_workspace)
        
        # Create git tags in the workspace
        for tag in tags:
            create_git_tag(temp_workspace, tag, f"Release {tag}")
    
    return env


@pytest.fixture
def bash_bumper_path() -> Path:
    """Return path to the Bash bumper script."""
    return BASH_BUMPER


@pytest.fixture
def go_bumper_path() -> Path:
    """Return path to the Go bumper binary."""
    return GO_BUMPER


@pytest.fixture
def temp_workspace(tmp_path: Path):
    """
    Create a temporary workspace with a git repository initialized.
    Yields the path and cleans up after the test.
    """
    workspace = tmp_path / "workspace"
    workspace.mkdir()
    
    # Initialize git repo
    import subprocess
    subprocess.run(
        ["git", "init"],
        cwd=workspace,
        capture_output=True,
        check=True
    )
    subprocess.run(
        ["git", "config", "user.email", "test@example.com"],
        cwd=workspace,
        capture_output=True,
        check=True
    )
    subprocess.run(
        ["git", "config", "user.name", "Test User"],
        cwd=workspace,
        capture_output=True,
        check=True
    )
    
    # Create initial commit
    readme = workspace / "README.md"
    readme.write_text("# Test Repository\n")
    subprocess.run(
        ["git", "add", "."],
        cwd=workspace,
        capture_output=True,
        check=True
    )
    subprocess.run(
        ["git", "commit", "-m", "Initial commit"],
        cwd=workspace,
        capture_output=True,
        check=True
    )
    
    yield workspace


@pytest.fixture
def github_output_file(tmp_path: Path) -> Path:
    """Create a temporary file for GITHUB_OUTPUT."""
    output_file = tmp_path / "github_output.txt"
    output_file.touch()
    return output_file


@pytest.fixture
def github_event_file(tmp_path: Path) -> Path:
    """Create a temporary file for GITHUB_EVENT_PATH."""
    event_file = tmp_path / "event.json"
    return event_file


@pytest.fixture
def base_env(github_output_file: Path, github_event_file: Path, temp_workspace: Path) -> dict:
    """
    Return base environment variables required for running bumper.
    """
    return {
        "GITHUB_OUTPUT": str(github_output_file),
        "GITHUB_EVENT_PATH": str(github_event_file),
        "GITHUB_EVENT_NAME": "pull_request",
        "GITHUB_REPOSITORY": "test-owner/test-repo",
        "GITHUB_WORKSPACE": str(temp_workspace),
        "GITHUB_SHA": "abc123def456",
        "GITHUB_ACTOR": "test-user",
        "GITHUB_API_URL": "https://api.github.com",
        "GITHUB_SERVER_URL": "https://github.com",
        "GITHUB_RUN_ID": "12345",
        "INPUT_GITHUB_TOKEN": "test-token-12345",
        "INPUT_BUMP_MAJOR": "bumper:major",
        "INPUT_BUMP_MINOR": "bumper:minor",
        "INPUT_BUMP_PATCH": "bumper:patch",
        "INPUT_BUMP_NONE": "bumper:none",
        "INPUT_BUMP_INCLUDE_V": "true",
        "INPUT_BUMP_SEMVER": "false",
        "INPUT_BUMP_LATEST": "false",
        "INPUT_BUMP_FAIL_IF_NO_LEVEL": "false",
        "INPUT_BUMP_DEFAULT_LEVEL": "",
    }
