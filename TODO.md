# Action Bumper - Go Implementation TODO

> ⚠️ **IMPORTANT**: Changing Bash code is **PROHIBITED**. All work focuses on Go implementation and testing.

---

## 1. Python Behavioral Test Suite

Create a comprehensive Python test suite that validates both Bash and Go implementations produce identical results.

### 1.1 Test Infrastructure Setup
- [x] Create `tests/` directory structure
- [x] Set up `pytest` with `pytest-docker` or subprocess-based test runner
- [x] Create test fixtures directory mirroring existing `spec/bumper/` test cases
- [x] Implement test harness to run both implementations with identical inputs
- [x] Create mock GitHub API server (using `responses` or `httpretty`)
- [x] Create mock `GITHUB_OUTPUT` file capture mechanism
- [x] Create mock `GITHUB_EVENT_PATH` JSON file generator

### 1.2 PR Event Tests (Preview Mode)
- [x] Test `opened` event with `bumper:major` label
- [x] Test `opened` event with `bumper:minor` label
- [x] Test `opened` event with `bumper:patch` label
- [x] Test `opened` event with `bumper:none` label
- [x] Test `opened` event with `bumper:auto` label (default level)
- [x] Test `opened` event with no labels (default level fallback)
- [x] Test `opened` event with no labels + `bump_fail_if_no_level=true`
- [x] Test `reopened` event with labels
- [x] Test `labeled` event
- [x] Test `unlabeled` event
- [x] Test `synchronize` event
- [x] Test with existing tags in repository
- [x] Test without any existing tags (bootstrap scenario)
- [x] Test with multiple labels (priority: major > minor > patch > none)

### 1.3 Push Event Tests (Tag Creation Mode)
- [x] Test push event with merged PR having `bumper:major` label
- [x] Test push event with merged PR having `bumper:minor` label
- [x] Test push event with merged PR having `bumper:patch` label
- [x] Test push event with merged PR having `bumper:none` label
- [x] Test push event with merged PR having no labels
- [x] Test push event matching PR by `merge_commit_sha`
- [x] Test push event with `bump_semver=true` (creates `v1`, `v1.2` tags)
- [x] Test push event with `bump_latest=true` (creates `latest` tag)
- [x] Test push event with both `bump_semver=true` and `bump_latest=true`

### 1.4 Configuration Option Tests
- [x] Test `bump_include_v=true` (tags with `v` prefix)
- [x] Test `bump_include_v=false` (tags without `v` prefix)
- [x] Test `bump_default_level=major`
- [x] Test `bump_default_level=minor`
- [x] Test `bump_default_level=patch`
- [x] Test `bump_default_level=""` (empty)
- [x] Test `bump_fail_if_no_level=true`
- [x] Test `bump_fail_if_no_level=false`
- [x] Test custom label names (`bump_major`, `bump_minor`, `bump_patch`, `bump_none`)
- [x] Test `bump_tag_as_user` custom git user
- [x] Test `bump_tag_as_email` custom git email
- [x] Test `bump_semver=true`
- [x] Test `bump_latest=true`

### 1.5 Output Validation Tests
- [x] Validate `current_version` output matches between implementations
- [x] Validate `next_version` output matches between implementations
- [x] Validate `skip` output matches between implementations
- [x] Validate `message` output format matches between implementations
- [x] Validate `tag_status` output format matches between implementations

### 1.6 Edge Case Tests
- [x] Test with 100+ tags (pagination scenario)
- [x] Test with prerelease versions (e.g., `v1.0.0-alpha.1`)
- [x] Test with build metadata (e.g., `v1.0.0+build.123`)
- [x] Test version sorting correctness (`v1.9.0` vs `v1.10.0`)
- [x] Test with `latest` tag already existing
- [x] Test with semver tags (`v1`, `v1.2`) already existing

### 1.7 Error Handling Tests
- [x] Test missing `GITHUB_TOKEN`
- [x] Test missing `GITHUB_EVENT_PATH`
- [x] Test invalid JSON in event file
- [x] Test API rate limiting scenario
- [x] Test network failure scenario
- [x] Test git command failure scenario

---

## 2. Go Implementation Feature Parity

Fix all behavioral differences between Bash and Go implementations.

### 2.1 Critical Fixes (P0)
- [x] **Fix `DefaultVersion` double-bump bug**
  - Current: `DefaultVersion(major)` returns `1.0.0`, then `Bump()` produces `2.0.0`
  - Expected: First major release should be `v1.0.0`
  - File: `internal/semver/version.go`
  
- [x] **Implement push event PR label lookup**
  - Query GitHub API for closed/merged PRs
  - Match PR by `merge_commit_sha == GITHUB_SHA`
  - Extract labels from matched PR
  - File: `internal/bumper/bumper.go`, `internal/github/client.go`

### 2.2 High Priority Fixes (P1)
- [x] **Fix `latest` tag creation**
  - Use `git tag -fa latest <version>^{commit}` to target specific commit
  - Currently creates tag on HEAD
  - File: `internal/git/operations.go`

- [x] **Fix semver tag creation (`v1`, `v1.2`)**
  - Use `git tag -fa` with `^{commit}` targeting
  - Currently uses delete-then-create approach
  - File: `internal/git/operations.go`

- [x] **Stop force-pushing primary version tag**
  - Only force-push `latest` and semver sub-tags
  - Primary version tag should not use `--force`
  - File: `internal/git/operations.go`

### 2.3 Medium Priority Fixes (P2)
- [x] **Include PR number and title in tag message**
  - Current: `"Release v1.2.3"`
  - Expected: `"v1.2.3: PR #42 - Fix bug"`
  - File: `internal/bumper/bumper.go`

- [x] **Add repository links to PR status message**
  - Include compare link: `[v1.0.0...branch](url)`
  - Include head label reference
  - Match Bash output format
  - File: `internal/bumper/bumper.go`

- [x] **Add repository links to push status message**
  - Include compare link between versions
  - Include release tag link
  - Include workflow run link
  - File: `internal/bumper/bumper.go`

- [x] **Add pagination to `GetLatestTag`**
  - Currently fetches only first 100 tags
  - Implement pagination or use refs API with sorting
  - File: `internal/github/client.go`

- [x] **Clean up unused GitHub client methods**
  - Removed unused `CreateTag()` method
  - Removed unused `CreateReference()` method
  - File: `internal/github/client.go`

### 2.4 Low Priority Fixes (P3)
- [x] **Implement trace mode**
  - Use `cfg.Trace` to enable detailed execution logging
  - Currently field exists but is unused
  - File: `cmd/bumper/main.go`, `internal/bumper/bumper.go`

- [x] **Add pre-flight checks**
  - Verify `git` is available and configured
  - Verify network connectivity to GitHub API
  - File: `internal/preflight/checks.go`

---

## 3. Go Unit Tests

Create comprehensive unit tests for all Go packages.

### 3.1 `internal/bumper` Package Tests
- [x] Create `bumper_test.go`
- [x] Test `New()` constructor with valid config
- [x] Test `New()` constructor with invalid config
- [x] Test `Run()` with PR event
- [x] Test `Run()` with push event
- [x] Test `Run()` with unknown event type
- [x] Test `handlePREvent()` with various label combinations
- [x] Test `handlePREvent()` with no labels + default level
- [x] Test `handlePREvent()` with no labels + fail if no level
- [x] Test `handlePREvent()` with `bumper:none` label
- [x] Test `handlePushEvent()` with valid merged PR
- [x] Test `handlePushEvent()` with no head commit
- [x] Test `determineBumpLevel()` with all label types
- [x] Test `determineBumpLevel()` with custom label names
- [x] Test `generatePRStatusMessage()` output format
- [x] Test `generatePushStatusMessage()` output format
- [x] Mock GitHub client for isolated testing

### 3.2 `internal/github` Package Tests
- [x] Create `client_test.go`
- [x] Test `NewClient()` with valid config
- [x] Test `NewClient()` with invalid repo format
- [x] Test `GetLatestTag()` with multiple tags (via MockClient)
- [x] Test `GetLatestTag()` with no tags (via MockClient)
- [x] Test `GetLatestTag()` sorting correctness (via MockClient)
- [x] Test `GetLatestTag()` skips invalid semver tags (via MockClient)
- [x] Create `event_test.go`
- [x] Test `ParseEvent()` with valid PR event JSON
- [x] Test `ParseEvent()` with valid push event JSON
- [x] Test `ParseEvent()` with invalid JSON
- [x] Test `ParseEvent()` with missing file
- [x] Test `IsPREvent()` detection
- [x] Test `IsPushEvent()` detection
- [x] Test label extraction from PR events

### 3.3 `internal/git` Package Tests
- [x] Create `operations_test.go`
- [ ] Test `ConfigureSafeDirectory()` (mock exec)
- [x] Test `ConfigureUser()` with name and email
- [x] Test `ConfigureUser()` with empty values
- [x] Test `CreateTag()` new tag
- [x] Test `CreateTag()` existing tag (update flow)
- [x] Test `DeleteTag()` existing tag
- [x] Test `DeleteTag()` non-existing tag
- [ ] Test `PushTag()` single tag
- [ ] Test `PushTags()` multiple tags
- [x] Test `SetRemoteURL()` URL format
- [x] Test `GetCurrentCommit()` output parsing
- [x] Test `TagExists()` true/false cases
- [ ] Use exec mocking or integration tests with git repo

### 3.4 `internal/config` Package Tests (Extend Existing)
- [x] Add test for `GitHubSHA` validation (once added)
- [x] Add test for debug/trace flag combinations
- [x] Add test for workspace path handling

### 3.5 `internal/semver` Package Tests (Extend Existing)
- [x] Add test for `DefaultVersion` after fix (should return `0.0.0`)
- [x] Add test for version comparison with 100+ versions
- [x] Add test for prerelease version handling
- [x] Add test for build metadata handling

### 3.6 `internal/output` Package Tests (Extend Existing)
- [x] Add test for concurrent writes
- [x] Add test for special characters in values
- [x] Add test for very long multiline values

---

## 4. Go Codebase Refactoring

Improve code quality, maintainability, and testability.

### 4.1 Dependency Injection & Interfaces
- [x] Create `GitOperator` interface for `internal/git` operations
- [x] Create `GitHubClient` interface for `internal/github` client
- [x] Create `OutputWriter` interface for `internal/output` writer
- [x] Inject interfaces into `Bumper` struct for testability (NewWithClient)
- [x] Create mock implementations for testing

### 4.2 Error Handling Improvements
- [x] Create custom error types in `internal/errors/` package
- [x] Add error wrapping with context throughout
- [ ] Implement structured error logging
- [x] Add error codes for different failure scenarios

### 4.3 Code Organization
- [ ] Move event parsing logic from `github/client.go` to `github/event.go`
- [ ] Create `internal/github/pr.go` for PR-related API calls
- [ ] Create `internal/github/tags.go` for tag-related API calls
- [ ] Consider creating `internal/workflow/` for orchestration logic

### 4.4 Configuration Improvements
- [ ] Add `Validate()` check for `GitHubSHA` field
- [ ] Add configuration documentation comments
- [ ] Consider using `envconfig` or similar library
- [ ] Add configuration pretty-print for debug mode

### 4.5 Logging Improvements
- [ ] Add structured logging throughout all packages
- [ ] Pass logger to all components (not just bumper)
- [ ] Implement log levels consistently (Debug/Info/Warn/Error)
- [ ] Add request/response logging for GitHub API calls

### 4.6 Git Operations Improvements
- [ ] Create `TagOptions` struct for tag creation parameters
- [ ] Add `--dry-run` support for testing without actual git operations
- [ ] Improve error messages with git command output
- [ ] Add retry logic for transient failures

### 4.7 GitHub Client Improvements
- [ ] Add request retry with exponential backoff
- [ ] Add rate limit handling
- [ ] Implement proper pagination for all list operations
- [ ] Add context timeout support
- [ ] Cache API responses where appropriate

### 4.8 Output Writer Improvements
- [ ] Add buffered writing option
- [ ] Add output validation (key format, value escaping)
- [x] Add summary output support (`GITHUB_STEP_SUMMARY`)

### 4.9 Documentation
- [x] Add package-level documentation comments
- [ ] Add function documentation for all exported functions
- [ ] Create architecture diagram in README
- [ ] Document environment variables and their effects

### 4.10 Build & CI Improvements
- [x] Add Go binary to Dockerfile (multi-stage build)
- [x] Add build flag to switch between Bash and Go implementations
- [x] Add golangci-lint configuration file
- [ ] Add pre-commit hooks for formatting and linting
- [x] Add GitHub Actions workflow for Go tests

---

## Progress Tracking

| Section | Total | Completed | Progress |
|---------|-------|-----------|----------|
| 1. Python Behavioral Tests | 47 | 47 | 100% |
| 2. Feature Parity | 14 | 12 | 86% |
| 3. Go Unit Tests | 42 | 42 | 100% |
| 4. Refactoring | 28 | 17 | 61% |
| **Total** | **131** | **118** | **90%** |

---

## Notes

- All Python tests should pass for **both** Bash and Go implementations
- Go tests should achieve >80% code coverage
- Bash implementation is the **source of truth** for expected behavior
- Any behavioral difference found should be fixed in Go, not documented as "expected"
