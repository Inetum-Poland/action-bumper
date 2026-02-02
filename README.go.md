# Action Bumper - Go Implementation

Go implementation of action-bumper, migrated from Bash for better maintainability and testability.

## Quick Start

### Build
```bash
make build          # or: go build -o bin/bumper ./cmd/bumper
```

### Test
```bash
make test           # or: go test ./...
make test-coverage  # Generate coverage report
```

## Project Structure

```
cmd/bumper/          # Main entry point
internal/
├── config/          # Configuration from env vars
├── semver/          # Semantic versioning logic
├── github/          # GitHub API client + interfaces
├── git/             # Git operations
├── bumper/          # Core business logic
├── output/          # GitHub Actions output + interface
└── logger/          # Logging abstraction
```

## Key Features

- **Zero external dependencies** - Everything in single binary
- **Testable** - Interfaces for all major components
- **Type-safe** - Strong typing with Go
- **Fast** - Compiled binary, no subprocess calls
- **Context-aware** - Proper context propagation

## Dependencies

- `github.com/google/go-github/v57` - GitHub API client
- `github.com/Masterminds/semver/v3` - Semantic versioning
- `golang.org/x/oauth2` - OAuth2 authentication
- `github.com/stretchr/testify` - Testing framework

## Development

```bash
make build      # Build binary
make test       # Run tests
make fmt        # Format code
make vet        # Run go vet
make tidy       # Tidy dependencies
make clean      # Clean artifacts
```

## Testing with Mocks

```go
import "github.com/Inetum-Poland/action-bumper/internal/logger"

// Use NoopLogger for tests
mockLogger := logger.NewNoopLogger()
```

## Architecture Highlights

- **Clean Architecture** - Separation of concerns
- **Dependency Injection** - Testable components via interfaces
- **Context Propagation** - Context passed through call chain
- **Error Wrapping** - All errors wrapped with context

## Status

| Package | Description | Test Coverage |
|---------|-------------|---------------|
| config | Configuration management | 92.9% |
| semver | Version bumping logic | 100% |
| output | GitHub Actions output | 81.5% |
| github | GitHub API client | - |
| git | Git operations | - |
| bumper | Core business logic | - |

## Migration from Bash

- **Before**: ~400 LOC Bash + external tools (jq, semver-tool, curl)
- **After**: ~1,600 LOC Go, self-contained binary
- **Benefits**: Faster, type-safe, testable, maintainable

## License

Copyright (c) 2024 Inetum Poland.
