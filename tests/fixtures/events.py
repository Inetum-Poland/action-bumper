# Copyright (c) 2024 Inetum Poland.
"""
Sample GitHub event JSON generators for testing.
"""

import json
from typing import Optional


def create_pr_event(
    action: str,
    pr_number: int = 42,
    pr_title: str = "Test PR",
    labels: Optional[list[str]] = None,
    head_sha: str = "abc123",
    head_label: str = "feature-branch",
    merge_commit_sha: Optional[str] = None,
) -> dict:
    """
    Create a pull_request event payload.
    
    Args:
        action: PR action (opened, labeled, unlabeled, synchronize, reopened)
        pr_number: Pull request number
        pr_title: Pull request title
        labels: List of label names
        head_sha: SHA of the head commit
        head_label: Label of the head branch
        merge_commit_sha: Merge commit SHA (if merged)
    
    Returns:
        Dictionary representing the event payload
    """
    labels = labels or []
    
    return {
        "action": action,
        "number": pr_number,
        "pull_request": {
            "number": pr_number,
            "title": pr_title,
            "labels": [{"name": label} for label in labels],
            "head": {
                "sha": head_sha,
                "label": head_label,
            },
            "merge_commit_sha": merge_commit_sha,
        },
        "repository": {
            "name": "test-repo",
            "full_name": "test-owner/test-repo",
            "default_branch": "main",
        },
    }


def create_push_event(
    head_commit_sha: str = "def456",
    head_commit_message: str = "Merge pull request #42",
    after_sha: str = "def456",
) -> dict:
    """
    Create a push event payload.
    
    Args:
        head_commit_sha: SHA of the head commit
        head_commit_message: Message of the head commit
        after_sha: The SHA after the push
    
    Returns:
        Dictionary representing the event payload
    """
    return {
        "after": after_sha,
        "head_commit": {
            "sha": head_commit_sha,
            "message": head_commit_message,
        },
        "commits": [
            {
                "sha": head_commit_sha,
                "message": head_commit_message,
            }
        ],
        "repository": {
            "name": "test-repo",
            "full_name": "test-owner/test-repo",
            "default_branch": "main",
        },
    }


def create_tags_response(tags: list[str]) -> list[dict]:
    """
    Create a GitHub API tags response.
    
    Args:
        tags: List of tag names (e.g., ["v1.0.0", "v0.9.0"])
    
    Returns:
        List of tag objects as returned by GitHub API
    """
    return [
        {
            "name": tag,
            "commit": {
                "sha": f"sha-for-{tag}",
                "url": f"https://api.github.com/repos/test-owner/test-repo/commits/sha-for-{tag}",
            },
        }
        for tag in tags
    ]


def create_pulls_response(
    pulls: list[dict],
) -> list[dict]:
    """
    Create a GitHub API pulls response for push event PR lookup.
    
    Args:
        pulls: List of pull request data dicts with keys:
               - number: PR number
               - title: PR title
               - labels: list of label names
               - merge_commit_sha: merge commit SHA
    
    Returns:
        List of PR objects as returned by GitHub API
    """
    result = []
    for pr in pulls:
        result.append({
            "number": pr.get("number", 1),
            "title": pr.get("title", "Test PR"),
            "labels": [{"name": label} for label in pr.get("labels", [])],
            "merge_commit_sha": pr.get("merge_commit_sha", "abc123"),
            "state": "closed",
            "merged_at": "2024-01-01T00:00:00Z",
        })
    return result
