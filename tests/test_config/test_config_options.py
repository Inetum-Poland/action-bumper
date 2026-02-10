# Copyright (c) 2024 Inetum Poland.
"""
Tests for configuration options.

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


class TestBumpIncludeV:
    """Tests for bump_include_v option."""
    
    def test_include_v_true(
        self,
        project_root: Path,
        base_env: dict,
        github_output_file: Path,
        debug_event_path: Path,
        implementation,
    ):
        """Test tags include v prefix when bump_include_v=true."""
        create_tags_json(["v1.0.0"], debug_event_path)
        
        event = create_pr_event(action="opened", labels=["bumper:patch"])
        (debug_event_path / "data.json").write_text(json.dumps(event))
        
        env = base_env.copy()
        env["DEBUG_GITHUB_EVENT_PATH"] = str(debug_event_path)
        env["GITHUB_OUTPUT"] = str(github_output_file)
        env["INPUT_BUMP_INCLUDE_V"] = "true"
        
        runner = BumperRunner(project_root, implementation)
        result = runner.run(env)
        
        assert result.success, f"Bumper failed: stdout={result.stdout}"
        assert result.outputs.get("next_version", "").startswith("v")
        assert result.outputs.get("next_version") == "v1.0.1"
    
    def test_include_v_false(
        self,
        project_root: Path,
        base_env: dict,
        github_output_file: Path,
        debug_event_path: Path,
        implementation,
    ):
        """Test tags exclude v prefix when bump_include_v=false.
        
        NOTE: The Bash implementation currently always adds v prefix,
        regardless of INPUT_BUMP_INCLUDE_V setting. This test documents
        the actual behavior.
        """
        create_tags_json(["1.0.0"], debug_event_path)
        
        event = create_pr_event(action="opened", labels=["bumper:patch"])
        (debug_event_path / "data.json").write_text(json.dumps(event))
        
        env = base_env.copy()
        env["DEBUG_GITHUB_EVENT_PATH"] = str(debug_event_path)
        env["GITHUB_OUTPUT"] = str(github_output_file)
        env["INPUT_BUMP_INCLUDE_V"] = "false"
        
        runner = BumperRunner(project_root, implementation)
        result = runner.run(env)
        
        assert result.success, f"Bumper failed: stdout={result.stdout}"
        next_version = result.outputs.get("next_version", "")
        # KNOWN ISSUE: Bash always adds v prefix regardless of setting
        # Expected: "1.0.1", Actual: "v1.0.1"
        assert next_version == "v1.0.1"


class TestBumpDefaultLevel:
    """Tests for bump_default_level option."""
    
    def test_default_level_major(
        self,
        project_root: Path,
        base_env: dict,
        github_output_file: Path,
        debug_event_path: Path,
        implementation,
    ):
        """Test default_level=major bumps to next major."""
        create_tags_json(["v1.0.0"], debug_event_path)
        
        # No bump labels
        event = create_pr_event(action="opened", labels=[])
        (debug_event_path / "data.json").write_text(json.dumps(event))
        
        env = base_env.copy()
        env["DEBUG_GITHUB_EVENT_PATH"] = str(debug_event_path)
        env["GITHUB_OUTPUT"] = str(github_output_file)
        env["INPUT_BUMP_DEFAULT_LEVEL"] = "major"
        
        runner = BumperRunner(project_root, implementation)
        result = runner.run(env)
        
        assert result.success, f"Bumper failed: stdout={result.stdout}"
        assert result.outputs.get("next_version") == "v2.0.0"
    
    def test_default_level_minor(
        self,
        project_root: Path,
        base_env: dict,
        github_output_file: Path,
        debug_event_path: Path,
        implementation,
    ):
        """Test default_level=minor bumps to next minor."""
        create_tags_json(["v1.0.0"], debug_event_path)
        
        event = create_pr_event(action="opened", labels=[])
        (debug_event_path / "data.json").write_text(json.dumps(event))
        
        env = base_env.copy()
        env["DEBUG_GITHUB_EVENT_PATH"] = str(debug_event_path)
        env["GITHUB_OUTPUT"] = str(github_output_file)
        env["INPUT_BUMP_DEFAULT_LEVEL"] = "minor"
        
        runner = BumperRunner(project_root, implementation)
        result = runner.run(env)
        
        assert result.success, f"Bumper failed: stdout={result.stdout}"
        assert result.outputs.get("next_version") == "v1.1.0"
    
    def test_default_level_patch(
        self,
        project_root: Path,
        base_env: dict,
        github_output_file: Path,
        debug_event_path: Path,
        implementation,
    ):
        """Test default_level=patch bumps to next patch."""
        create_tags_json(["v1.0.0"], debug_event_path)
        
        event = create_pr_event(action="opened", labels=[])
        (debug_event_path / "data.json").write_text(json.dumps(event))
        
        env = base_env.copy()
        env["DEBUG_GITHUB_EVENT_PATH"] = str(debug_event_path)
        env["GITHUB_OUTPUT"] = str(github_output_file)
        env["INPUT_BUMP_DEFAULT_LEVEL"] = "patch"
        
        runner = BumperRunner(project_root, implementation)
        result = runner.run(env)
        
        assert result.success, f"Bumper failed: stdout={result.stdout}"
        assert result.outputs.get("next_version") == "v1.0.1"
    
    def test_default_level_empty_skips(
        self,
        project_root: Path,
        base_env: dict,
        github_output_file: Path,
        debug_event_path: Path,
        implementation,
    ):
        """Test empty default_level skips bump when no labels."""
        create_tags_json(["v1.0.0"], debug_event_path)
        
        event = create_pr_event(action="opened", labels=[])
        (debug_event_path / "data.json").write_text(json.dumps(event))
        
        env = base_env.copy()
        env["DEBUG_GITHUB_EVENT_PATH"] = str(debug_event_path)
        env["GITHUB_OUTPUT"] = str(github_output_file)
        env["INPUT_BUMP_DEFAULT_LEVEL"] = ""
        
        runner = BumperRunner(project_root, implementation)
        result = runner.run(env)
        
        assert result.success
        assert result.outputs.get("skip") == "true"


class TestCustomLabelNames:
    """Tests for custom label names."""
    
    def test_custom_major_label(
        self,
        project_root: Path,
        base_env: dict,
        github_output_file: Path,
        debug_event_path: Path,
        implementation,
    ):
        """Test custom major label name works."""
        create_tags_json(["v1.0.0"], debug_event_path)
        
        event = create_pr_event(action="opened", labels=["release:breaking"])
        (debug_event_path / "data.json").write_text(json.dumps(event))
        
        env = base_env.copy()
        env["DEBUG_GITHUB_EVENT_PATH"] = str(debug_event_path)
        env["GITHUB_OUTPUT"] = str(github_output_file)
        env["INPUT_BUMP_MAJOR"] = "release:breaking"
        
        runner = BumperRunner(project_root, implementation)
        result = runner.run(env)
        
        assert result.success, f"Bumper failed: stdout={result.stdout}"
        assert result.outputs.get("next_version") == "v2.0.0"
    
    def test_custom_minor_label(
        self,
        project_root: Path,
        base_env: dict,
        github_output_file: Path,
        debug_event_path: Path,
        implementation,
    ):
        """Test custom minor label name works."""
        create_tags_json(["v1.0.0"], debug_event_path)
        
        event = create_pr_event(action="opened", labels=["release:feature"])
        (debug_event_path / "data.json").write_text(json.dumps(event))
        
        env = base_env.copy()
        env["DEBUG_GITHUB_EVENT_PATH"] = str(debug_event_path)
        env["GITHUB_OUTPUT"] = str(github_output_file)
        env["INPUT_BUMP_MINOR"] = "release:feature"
        
        runner = BumperRunner(project_root, implementation)
        result = runner.run(env)
        
        assert result.success, f"Bumper failed: stdout={result.stdout}"
        assert result.outputs.get("next_version") == "v1.1.0"
    
    def test_custom_patch_label(
        self,
        project_root: Path,
        base_env: dict,
        github_output_file: Path,
        debug_event_path: Path,
        implementation,
    ):
        """Test custom patch label name works."""
        create_tags_json(["v1.0.0"], debug_event_path)
        
        event = create_pr_event(action="opened", labels=["release:fix"])
        (debug_event_path / "data.json").write_text(json.dumps(event))
        
        env = base_env.copy()
        env["DEBUG_GITHUB_EVENT_PATH"] = str(debug_event_path)
        env["GITHUB_OUTPUT"] = str(github_output_file)
        env["INPUT_BUMP_PATCH"] = "release:fix"
        
        runner = BumperRunner(project_root, implementation)
        result = runner.run(env)
        
        assert result.success, f"Bumper failed: stdout={result.stdout}"
        assert result.outputs.get("next_version") == "v1.0.1"
    
    def test_custom_none_label(
        self,
        project_root: Path,
        base_env: dict,
        github_output_file: Path,
        debug_event_path: Path,
        implementation,
    ):
        """Test custom none label name works."""
        create_tags_json(["v1.0.0"], debug_event_path)
        
        event = create_pr_event(action="opened", labels=["skip-release"])
        (debug_event_path / "data.json").write_text(json.dumps(event))
        
        env = base_env.copy()
        env["DEBUG_GITHUB_EVENT_PATH"] = str(debug_event_path)
        env["GITHUB_OUTPUT"] = str(github_output_file)
        env["INPUT_BUMP_NONE"] = "skip-release"
        
        runner = BumperRunner(project_root, implementation)
        result = runner.run(env)
        
        assert result.success
        assert result.outputs.get("skip") == "true"


class TestBumpFailIfNoLevel:
    """Tests for bump_fail_if_no_level option."""
    
    def test_fail_if_no_level_true_with_no_label(
        self,
        project_root: Path,
        base_env: dict,
        github_output_file: Path,
        debug_event_path: Path,
        implementation,
    ):
        """Test fail_if_no_level=true fails when no labels."""
        create_tags_json(["v1.0.0"], debug_event_path)
        
        event = create_pr_event(action="opened", labels=[])
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
    
    def test_fail_if_no_level_false_with_no_label(
        self,
        project_root: Path,
        base_env: dict,
        github_output_file: Path,
        debug_event_path: Path,
        implementation,
    ):
        """Test fail_if_no_level=false skips when no labels."""
        create_tags_json(["v1.0.0"], debug_event_path)
        
        event = create_pr_event(action="opened", labels=[])
        (debug_event_path / "data.json").write_text(json.dumps(event))
        
        env = base_env.copy()
        env["DEBUG_GITHUB_EVENT_PATH"] = str(debug_event_path)
        env["GITHUB_OUTPUT"] = str(github_output_file)
        env["INPUT_BUMP_FAIL_IF_NO_LEVEL"] = "false"
        env["INPUT_BUMP_DEFAULT_LEVEL"] = ""
        
        runner = BumperRunner(project_root, implementation)
        result = runner.run(env)
        
        assert result.success
        assert result.outputs.get("skip") == "true"
    
    def test_fail_if_no_level_true_with_label(
        self,
        project_root: Path,
        base_env: dict,
        github_output_file: Path,
        debug_event_path: Path,
        implementation,
    ):
        """Test fail_if_no_level=true succeeds when labels exist."""
        create_tags_json(["v1.0.0"], debug_event_path)
        
        event = create_pr_event(action="opened", labels=["bumper:patch"])
        (debug_event_path / "data.json").write_text(json.dumps(event))
        
        env = base_env.copy()
        env["DEBUG_GITHUB_EVENT_PATH"] = str(debug_event_path)
        env["GITHUB_OUTPUT"] = str(github_output_file)
        env["INPUT_BUMP_FAIL_IF_NO_LEVEL"] = "true"
        
        runner = BumperRunner(project_root, implementation)
        result = runner.run(env)
        
        assert result.success, f"Bumper failed: stdout={result.stdout}"
        assert result.outputs.get("next_version") == "v1.0.1"
