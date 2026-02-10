# Copyright (c) 2024 Inetum Poland.
"""
Tests for PR 'opened' event with various labels.

These tests use DEBUG_GITHUB_EVENT_PATH to mock GitHub API responses
for both Bash and Go implementations.

Run with: pytest --impl=bash (default) or pytest --impl=go
"""

import json
import pytest
from pathlib import Path

from tests.fixtures.events import create_pr_event
from tests.helpers.runner import BumperRunner, Implementation
from tests.conftest import create_tags_json


class TestOpenedEventWithLabels:
    """Tests for PR opened events with bump labels."""
    
    def test_opened_event_bumper_patch(
        self,
        project_root: Path,
        base_env: dict,
        github_output_file: Path,
        debug_event_path: Path,
        implementation,
    ):
        """Test PR opened with bumper:patch label bumps patch version."""
        create_tags_json(["v0.9.1", "v0.9.0"], debug_event_path)
        
        event = create_pr_event(
            action="opened",
            pr_number=42,
            pr_title="Fix bug",
            labels=["bumper:patch"],
        )
        (debug_event_path / "data.json").write_text(json.dumps(event))
        
        env = base_env.copy()
        env["DEBUG_GITHUB_EVENT_PATH"] = str(debug_event_path)
        env["GITHUB_OUTPUT"] = str(github_output_file)
        
        runner = BumperRunner(project_root, implementation)
        result = runner.run(env)
        
        assert result.success, f"Bumper failed: stdout={result.stdout}, stderr={result.stderr}"
        assert result.outputs.get("next_version") == "v0.9.2"
        assert result.outputs.get("skip") != "true"
    
    def test_opened_event_bumper_minor(
        self,
        project_root: Path,
        base_env: dict,
        github_output_file: Path,
        debug_event_path: Path,
        implementation,
    ):
        """Test PR opened with bumper:minor label bumps minor version."""
        create_tags_json(["v0.9.1"], debug_event_path)
        
        event = create_pr_event(
            action="opened",
            labels=["bumper:minor"],
        )
        (debug_event_path / "data.json").write_text(json.dumps(event))
        
        env = base_env.copy()
        env["DEBUG_GITHUB_EVENT_PATH"] = str(debug_event_path)
        env["GITHUB_OUTPUT"] = str(github_output_file)
        
        runner = BumperRunner(project_root, implementation)
        result = runner.run(env)
        
        assert result.success, f"Bumper failed: stdout={result.stdout}"
        assert result.outputs.get("next_version") == "v0.10.0"
    
    def test_opened_event_bumper_major(
        self,
        project_root: Path,
        base_env: dict,
        github_output_file: Path,
        debug_event_path: Path,
        implementation,
    ):
        """Test PR opened with bumper:major label bumps major version."""
        create_tags_json(["v0.9.1"], debug_event_path)
        
        event = create_pr_event(
            action="opened",
            labels=["bumper:major"],
        )
        (debug_event_path / "data.json").write_text(json.dumps(event))
        
        env = base_env.copy()
        env["DEBUG_GITHUB_EVENT_PATH"] = str(debug_event_path)
        env["GITHUB_OUTPUT"] = str(github_output_file)
        
        runner = BumperRunner(project_root, implementation)
        result = runner.run(env)
        
        assert result.success, f"Bumper failed: stdout={result.stdout}"
        assert result.outputs.get("next_version") == "v1.0.0"
    
    def test_opened_event_bumper_none(
        self,
        project_root: Path,
        base_env: dict,
        github_output_file: Path,
        debug_event_path: Path,
        implementation,
    ):
        """Test PR opened with bumper:none label skips bump."""
        create_tags_json(["v0.9.1"], debug_event_path)
        
        event = create_pr_event(
            action="opened",
            labels=["bumper:none"],
        )
        (debug_event_path / "data.json").write_text(json.dumps(event))
        
        env = base_env.copy()
        env["DEBUG_GITHUB_EVENT_PATH"] = str(debug_event_path)
        env["GITHUB_OUTPUT"] = str(github_output_file)
        
        runner = BumperRunner(project_root, implementation)
        result = runner.run(env)
        
        assert result.success
        assert result.outputs.get("skip") == "true"


class TestOpenedEventWithoutTags:
    """Tests for PR opened events when no tags exist (bootstrap)."""
    
    def test_opened_event_patch_no_tags(
        self,
        project_root: Path,
        base_env: dict,
        github_output_file: Path,
        debug_event_path: Path,
        implementation,
    ):
        """Test first patch release creates v0.0.1."""
        # Empty tags list
        create_tags_json([], debug_event_path)
        
        event = create_pr_event(
            action="opened",
            labels=["bumper:patch"],
        )
        (debug_event_path / "data.json").write_text(json.dumps(event))
        
        env = base_env.copy()
        env["DEBUG_GITHUB_EVENT_PATH"] = str(debug_event_path)
        env["GITHUB_OUTPUT"] = str(github_output_file)
        
        runner = BumperRunner(project_root, implementation)
        result = runner.run(env)
        
        assert result.success, f"Bumper failed: stdout={result.stdout}"
        assert result.outputs.get("next_version") == "v0.0.1"
    
    def test_opened_event_minor_no_tags(
        self,
        project_root: Path,
        base_env: dict,
        github_output_file: Path,
        debug_event_path: Path,
        implementation,
    ):
        """Test first minor release creates v0.1.0."""
        create_tags_json([], debug_event_path)
        
        event = create_pr_event(
            action="opened",
            labels=["bumper:minor"],
        )
        (debug_event_path / "data.json").write_text(json.dumps(event))
        
        env = base_env.copy()
        env["DEBUG_GITHUB_EVENT_PATH"] = str(debug_event_path)
        env["GITHUB_OUTPUT"] = str(github_output_file)
        
        runner = BumperRunner(project_root, implementation)
        result = runner.run(env)
        
        assert result.success, f"Bumper failed: stdout={result.stdout}"
        assert result.outputs.get("next_version") == "v0.1.0"
    
    def test_opened_event_major_no_tags(
        self,
        project_root: Path,
        base_env: dict,
        github_output_file: Path,
        debug_event_path: Path,
        implementation,
    ):
        """Test first major release creates v1.0.0."""
        create_tags_json([], debug_event_path)
        
        event = create_pr_event(
            action="opened",
            labels=["bumper:major"],
        )
        (debug_event_path / "data.json").write_text(json.dumps(event))
        
        env = base_env.copy()
        env["DEBUG_GITHUB_EVENT_PATH"] = str(debug_event_path)
        env["GITHUB_OUTPUT"] = str(github_output_file)
        
        runner = BumperRunner(project_root, implementation)
        result = runner.run(env)
        
        assert result.success, f"Bumper failed: stdout={result.stdout}"
        assert result.outputs.get("next_version") == "v1.0.0"


class TestOpenedEventWithoutLabels:
    """Tests for PR opened events when no bump labels are present."""
    
    def test_opened_event_no_labels_no_default(
        self,
        project_root: Path,
        base_env: dict,
        github_output_file: Path,
        debug_event_path: Path,
        implementation,
    ):
        """Test PR without labels and no default skips bump."""
        create_tags_json(["v1.0.0"], debug_event_path)
        
        event = create_pr_event(
            action="opened",
            labels=[],
        )
        (debug_event_path / "data.json").write_text(json.dumps(event))
        
        env = base_env.copy()
        env["DEBUG_GITHUB_EVENT_PATH"] = str(debug_event_path)
        env["GITHUB_OUTPUT"] = str(github_output_file)
        env["INPUT_BUMP_DEFAULT_LEVEL"] = ""
        
        runner = BumperRunner(project_root, implementation)
        result = runner.run(env)
        
        assert result.success
        assert result.outputs.get("skip") == "true"
    
    def test_opened_event_no_labels_with_default_patch(
        self,
        project_root: Path,
        base_env: dict,
        github_output_file: Path,
        debug_event_path: Path,
        implementation,
    ):
        """Test PR without labels uses default_level=patch."""
        create_tags_json(["v1.0.0"], debug_event_path)
        
        event = create_pr_event(
            action="opened",
            labels=[],
        )
        (debug_event_path / "data.json").write_text(json.dumps(event))
        
        env = base_env.copy()
        env["DEBUG_GITHUB_EVENT_PATH"] = str(debug_event_path)
        env["GITHUB_OUTPUT"] = str(github_output_file)
        env["INPUT_BUMP_DEFAULT_LEVEL"] = "patch"
        
        runner = BumperRunner(project_root, implementation)
        result = runner.run(env)
        
        assert result.success, f"Bumper failed: stdout={result.stdout}"
        assert result.outputs.get("next_version") == "v1.0.1"
    
    def test_opened_event_no_labels_fail_if_no_level(
        self,
        project_root: Path,
        base_env: dict,
        github_output_file: Path,
        debug_event_path: Path,
        implementation,
    ):
        """Test fail_if_no_level causes failure when no labels."""
        create_tags_json(["v1.0.0"], debug_event_path)
        
        event = create_pr_event(
            action="opened",
            labels=[],
        )
        (debug_event_path / "data.json").write_text(json.dumps(event))
        
        env = base_env.copy()
        env["DEBUG_GITHUB_EVENT_PATH"] = str(debug_event_path)
        env["GITHUB_OUTPUT"] = str(github_output_file)
        env["INPUT_BUMP_FAIL_IF_NO_LEVEL"] = "true"
        env["INPUT_BUMP_DEFAULT_LEVEL"] = ""
        
        runner = BumperRunner(project_root, implementation)
        result = runner.run(env)
        
        assert not result.success
        assert "no bump label" in result.stdout.lower()


class TestOtherPREvents:
    """Tests for other PR event types (labeled, unlabeled, synchronize, reopened)."""
    
    def test_labeled_event(
        self,
        project_root: Path,
        base_env: dict,
        github_output_file: Path,
        debug_event_path: Path,
        implementation,
    ):
        """Test 'labeled' event works like 'opened'."""
        create_tags_json(["v1.0.0"], debug_event_path)
        
        event = create_pr_event(
            action="labeled",
            labels=["bumper:patch"],
        )
        (debug_event_path / "data.json").write_text(json.dumps(event))
        
        env = base_env.copy()
        env["DEBUG_GITHUB_EVENT_PATH"] = str(debug_event_path)
        env["GITHUB_OUTPUT"] = str(github_output_file)
        
        runner = BumperRunner(project_root, implementation)
        result = runner.run(env)
        
        assert result.success, f"Bumper failed: stdout={result.stdout}"
        assert result.outputs.get("next_version") == "v1.0.1"
    
    def test_unlabeled_event(
        self,
        project_root: Path,
        base_env: dict,
        github_output_file: Path,
        debug_event_path: Path,
        implementation,
    ):
        """Test 'unlabeled' event recalculates bump."""
        create_tags_json(["v1.0.0"], debug_event_path)
        
        event = create_pr_event(
            action="unlabeled",
            labels=["bumper:minor"],
        )
        (debug_event_path / "data.json").write_text(json.dumps(event))
        
        env = base_env.copy()
        env["DEBUG_GITHUB_EVENT_PATH"] = str(debug_event_path)
        env["GITHUB_OUTPUT"] = str(github_output_file)
        
        runner = BumperRunner(project_root, implementation)
        result = runner.run(env)
        
        assert result.success, f"Bumper failed: stdout={result.stdout}"
        assert result.outputs.get("next_version") == "v1.1.0"
    
    def test_synchronize_event(
        self,
        project_root: Path,
        base_env: dict,
        github_output_file: Path,
        debug_event_path: Path,
        implementation,
    ):
        """Test 'synchronize' event recalculates bump."""
        create_tags_json(["v1.0.0"], debug_event_path)
        
        event = create_pr_event(
            action="synchronize",
            labels=["bumper:patch"],
        )
        (debug_event_path / "data.json").write_text(json.dumps(event))
        
        env = base_env.copy()
        env["DEBUG_GITHUB_EVENT_PATH"] = str(debug_event_path)
        env["GITHUB_OUTPUT"] = str(github_output_file)
        
        runner = BumperRunner(project_root, implementation)
        result = runner.run(env)
        
        assert result.success, f"Bumper failed: stdout={result.stdout}"
        assert result.outputs.get("next_version") == "v1.0.1"
    
    def test_reopened_event(
        self,
        project_root: Path,
        base_env: dict,
        github_output_file: Path,
        debug_event_path: Path,
        implementation,
    ):
        """Test 'reopened' event recalculates bump."""
        create_tags_json(["v1.0.0"], debug_event_path)
        
        event = create_pr_event(
            action="reopened",
            labels=["bumper:major"],
        )
        (debug_event_path / "data.json").write_text(json.dumps(event))
        
        env = base_env.copy()
        env["DEBUG_GITHUB_EVENT_PATH"] = str(debug_event_path)
        env["GITHUB_OUTPUT"] = str(github_output_file)
        
        runner = BumperRunner(project_root, implementation)
        result = runner.run(env)
        
        assert result.success, f"Bumper failed: stdout={result.stdout}"
        assert result.outputs.get("next_version") == "v2.0.0"


class TestLabelPriority:
    """Tests for label priority when multiple bump labels are present."""
    
    def test_major_takes_priority_over_minor(
        self,
        project_root: Path,
        base_env: dict,
        github_output_file: Path,
        debug_event_path: Path,
        implementation,
    ):
        """Major label should take priority over minor."""
        create_tags_json(["v1.0.0"], debug_event_path)
        
        event = create_pr_event(
            action="opened",
            labels=["bumper:minor", "bumper:major"],
        )
        (debug_event_path / "data.json").write_text(json.dumps(event))
        
        env = base_env.copy()
        env["DEBUG_GITHUB_EVENT_PATH"] = str(debug_event_path)
        env["GITHUB_OUTPUT"] = str(github_output_file)
        
        runner = BumperRunner(project_root, implementation)
        result = runner.run(env)
        
        assert result.success, f"Bumper failed: stdout={result.stdout}"
        assert result.outputs.get("next_version") == "v2.0.0"
    
    def test_minor_takes_priority_over_patch(
        self,
        project_root: Path,
        base_env: dict,
        github_output_file: Path,
        debug_event_path: Path,
        implementation,
    ):
        """Minor label should take priority over patch."""
        create_tags_json(["v1.0.0"], debug_event_path)
        
        event = create_pr_event(
            action="opened",
            labels=["bumper:patch", "bumper:minor"],
        )
        (debug_event_path / "data.json").write_text(json.dumps(event))
        
        env = base_env.copy()
        env["DEBUG_GITHUB_EVENT_PATH"] = str(debug_event_path)
        env["GITHUB_OUTPUT"] = str(github_output_file)
        
        runner = BumperRunner(project_root, implementation)
        result = runner.run(env)
        
        assert result.success, f"Bumper failed: stdout={result.stdout}"
        assert result.outputs.get("next_version") == "v1.1.0"
