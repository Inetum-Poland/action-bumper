# Copyright (c) 2024 Inetum Poland.
"""
Git helper functions for test setup.
"""

import subprocess
from pathlib import Path


def create_git_tag(
    repo_path: Path,
    tag_name: str,
    message: str | None = None,
) -> None:
    """
    Create a git tag in the repository.
    
    Args:
        repo_path: Path to the git repository
        tag_name: Name of the tag to create
        message: Optional message for annotated tag
    """
    if message:
        cmd = ["git", "tag", "-a", tag_name, "-m", message]
    else:
        cmd = ["git", "tag", tag_name]
    
    subprocess.run(
        cmd,
        cwd=repo_path,
        capture_output=True,
        check=True,
    )


def list_git_tags(repo_path: Path) -> list[str]:
    """
    List all tags in the repository.
    
    Args:
        repo_path: Path to the git repository
    
    Returns:
        List of tag names
    """
    result = subprocess.run(
        ["git", "tag", "-l"],
        cwd=repo_path,
        capture_output=True,
        text=True,
        check=True,
    )
    
    return [tag.strip() for tag in result.stdout.split('\n') if tag.strip()]


def get_current_commit_sha(repo_path: Path) -> str:
    """
    Get the current commit SHA.
    
    Args:
        repo_path: Path to the git repository
    
    Returns:
        Current commit SHA
    """
    result = subprocess.run(
        ["git", "rev-parse", "HEAD"],
        cwd=repo_path,
        capture_output=True,
        text=True,
        check=True,
    )
    
    return result.stdout.strip()


def create_commit(
    repo_path: Path,
    message: str = "Test commit",
    file_name: str = "test.txt",
    file_content: str = "test content",
) -> str:
    """
    Create a new commit in the repository.
    
    Args:
        repo_path: Path to the git repository
        message: Commit message
        file_name: Name of file to create/modify
        file_content: Content to write to file
    
    Returns:
        SHA of the new commit
    """
    file_path = repo_path / file_name
    file_path.write_text(file_content)
    
    subprocess.run(
        ["git", "add", file_name],
        cwd=repo_path,
        capture_output=True,
        check=True,
    )
    
    subprocess.run(
        ["git", "commit", "-m", message],
        cwd=repo_path,
        capture_output=True,
        check=True,
    )
    
    return get_current_commit_sha(repo_path)
