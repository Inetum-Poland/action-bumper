# Copyright (c) 2024 Inetum Poland.
"""
Mock GitHub API server for testing.
Uses the `responses` library to mock HTTP requests.
"""

import json
import responses
from typing import Optional
from urllib.parse import urljoin


class MockGitHubAPI:
    """
    Mock GitHub API server for testing bumper implementations.
    
    Provides methods to set up mock responses for GitHub API endpoints
    that the bumper uses.
    """
    
    BASE_URL = "https://api.github.com"
    
    def __init__(self):
        """Initialize the mock server."""
        self._tags: list[dict] = []
        self._pulls: list[dict] = []
        self._repo = "test-owner/test-repo"
    
    def set_repo(self, owner: str, repo: str) -> None:
        """Set the repository for API calls."""
        self._repo = f"{owner}/{repo}"
    
    def set_tags(self, tags: list[dict]) -> None:
        """
        Set the tags that will be returned by the API.
        
        Args:
            tags: List of tag objects with 'name' and optionally 'commit.sha'
        """
        self._tags = tags
    
    def set_pulls(self, pulls: list[dict]) -> None:
        """
        Set the pull requests that will be returned by the API.
        
        Args:
            pulls: List of PR objects with labels, merge_commit_sha, etc.
        """
        self._pulls = pulls
    
    def add_tag(self, name: str, sha: str = "abc123") -> None:
        """Add a single tag to the mock."""
        self._tags.append({
            "name": name,
            "commit": {"sha": sha}
        })
    
    def add_pull(
        self,
        number: int,
        title: str,
        labels: list[str],
        merge_commit_sha: str,
    ) -> None:
        """Add a single pull request to the mock."""
        self._pulls.append({
            "number": number,
            "title": title,
            "labels": [{"name": label} for label in labels],
            "merge_commit_sha": merge_commit_sha,
            "state": "closed",
        })
    
    def setup_responses(self) -> None:
        """
        Set up all mock responses.
        Call this within a `responses.activate` context or decorator.
        """
        # Mock tags endpoint
        tags_url = f"{self.BASE_URL}/repos/{self._repo}/tags"
        responses.add(
            responses.GET,
            tags_url,
            json=self._tags,
            status=200,
        )
        
        # Also mock with query parameters
        responses.add(
            responses.GET,
            tags_url,
            json=self._tags,
            status=200,
            match_querystring=False,
        )
        
        # Mock pulls endpoint
        pulls_url = f"{self.BASE_URL}/repos/{self._repo}/pulls"
        responses.add(
            responses.GET,
            pulls_url,
            json=self._pulls,
            status=200,
            match_querystring=False,
        )
    
    @classmethod
    def create_tags_list(cls, tag_names: list[str]) -> list[dict]:
        """
        Create a list of tag objects from tag names.
        
        Args:
            tag_names: List of tag names (e.g., ["v1.0.0", "v0.9.0"])
        
        Returns:
            List of tag objects suitable for API response
        """
        return [
            {
                "name": name,
                "commit": {"sha": f"sha-{name}"}
            }
            for name in tag_names
        ]
    
    @classmethod
    def create_pulls_list(
        cls,
        pulls_data: list[tuple[int, str, list[str], str]],
    ) -> list[dict]:
        """
        Create a list of PR objects.
        
        Args:
            pulls_data: List of tuples (number, title, labels, merge_commit_sha)
        
        Returns:
            List of PR objects suitable for API response
        """
        return [
            {
                "number": num,
                "title": title,
                "labels": [{"name": label} for label in labels],
                "merge_commit_sha": sha,
                "state": "closed",
            }
            for num, title, labels, sha in pulls_data
        ]


def mock_github_api_for_tags(tags: list[str], repo: str = "test-owner/test-repo"):
    """
    Decorator/context manager to mock GitHub API for tag listing.
    
    Usage:
        @mock_github_api_for_tags(["v1.0.0", "v0.9.0"])
        def test_something():
            ...
    """
    def decorator(func):
        @responses.activate
        def wrapper(*args, **kwargs):
            mock = MockGitHubAPI()
            mock.set_repo(*repo.split("/"))
            mock.set_tags(MockGitHubAPI.create_tags_list(tags))
            mock.setup_responses()
            return func(*args, **kwargs)
        return wrapper
    return decorator
