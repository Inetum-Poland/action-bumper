# Copyright (c) 2024 Inetum Poland.
"""
Test runner for executing Bash and Go bumper implementations.
"""

import subprocess
import os
from dataclasses import dataclass
from pathlib import Path
from typing import Optional
from enum import Enum


class Implementation(Enum):
    """Bumper implementation type."""
    BASH = "bash"
    GO = "go"


@dataclass
class BumperResult:
    """Result of running the bumper."""
    exit_code: int
    stdout: str
    stderr: str
    outputs: dict[str, str]  # Parsed GITHUB_OUTPUT contents
    
    @property
    def success(self) -> bool:
        """Return True if the bumper exited successfully."""
        return self.exit_code == 0


class BumperRunner:
    """
    Runner for executing bumper implementations.
    
    Supports running both Bash and Go implementations with the same
    environment and comparing their outputs.
    """
    
    def __init__(
        self,
        project_root: Path,
        implementation: Implementation = Implementation.BASH,
    ):
        """
        Initialize the runner.
        
        Args:
            project_root: Path to the project root directory
            implementation: Which implementation to run
        """
        self.project_root = project_root
        self.implementation = implementation
        
        if implementation == Implementation.BASH:
            self.executable = project_root / "bumper.sh"
        else:
            self.executable = project_root / "bin" / "bumper"
    
    def run(
        self,
        env: dict[str, str],
        timeout: int = 30,
    ) -> BumperResult:
        """
        Run the bumper with the given environment.
        
        Args:
            env: Environment variables to set
            timeout: Timeout in seconds
        
        Returns:
            BumperResult with exit code, output, and parsed outputs
        """
        # Ensure the executable exists
        if not self.executable.exists():
            raise FileNotFoundError(
                f"Bumper executable not found: {self.executable}"
            )
        
        # Build command based on implementation
        if self.implementation == Implementation.BASH:
            cmd = ["bash", str(self.executable)]
        else:
            cmd = [str(self.executable)]
        
        # Ensure HOME is set (git requires it)
        run_env = env.copy()
        if "HOME" not in run_env:
            run_env["HOME"] = os.environ.get("HOME", "/tmp")
        
        # Ensure PATH is set (for git and other tools)
        if "PATH" not in run_env:
            run_env["PATH"] = os.environ.get("PATH", "/usr/bin:/bin")
        
        # Run the bumper
        result = subprocess.run(
            cmd,
            env=run_env,
            capture_output=True,
            text=True,
            timeout=timeout,
            cwd=run_env.get("GITHUB_WORKSPACE", self.project_root),
        )
        
        # Parse GITHUB_OUTPUT file
        outputs = {}
        github_output_path = run_env.get("GITHUB_OUTPUT")
        if github_output_path and Path(github_output_path).exists():
            from .output_parser import parse_github_output
            outputs = parse_github_output(Path(github_output_path))
        
        return BumperResult(
            exit_code=result.returncode,
            stdout=result.stdout,
            stderr=result.stderr,
            outputs=outputs,
        )
    
    @classmethod
    def run_both(
        cls,
        project_root: Path,
        env: dict[str, str],
        timeout: int = 30,
    ) -> tuple[BumperResult, BumperResult]:
        """
        Run both implementations and return their results.
        
        Args:
            project_root: Path to the project root directory
            env: Environment variables to set
            timeout: Timeout in seconds
        
        Returns:
            Tuple of (bash_result, go_result)
        """
        # Create separate output files for each implementation
        import tempfile
        
        bash_env = env.copy()
        go_env = env.copy()
        
        with tempfile.NamedTemporaryFile(mode='w', delete=False, suffix='_bash.txt') as f:
            bash_output_file = f.name
        with tempfile.NamedTemporaryFile(mode='w', delete=False, suffix='_go.txt') as f:
            go_output_file = f.name
        
        bash_env["GITHUB_OUTPUT"] = bash_output_file
        go_env["GITHUB_OUTPUT"] = go_output_file
        
        try:
            bash_runner = cls(project_root, Implementation.BASH)
            go_runner = cls(project_root, Implementation.GO)
            
            bash_result = bash_runner.run(bash_env, timeout)
            go_result = go_runner.run(go_env, timeout)
            
            return bash_result, go_result
        finally:
            # Cleanup temp files
            Path(bash_output_file).unlink(missing_ok=True)
            Path(go_output_file).unlink(missing_ok=True)
