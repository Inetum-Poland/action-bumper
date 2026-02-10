// Copyright (c) 2024 Inetum Poland.

// Package git provides interfaces and implementations for git operations.
package git

// Operator defines the interface for git operations.
// This allows for mocking in tests.
type Operator interface {
	// ConfigureSafeDirectory configures git safe directory
	ConfigureSafeDirectory(dir string) error

	// ConfigureUser configures git user name and email
	ConfigureUser(name, email string) error

	// CreateTag creates an annotated git tag (does not force)
	CreateTag(tag, message string) error

	// CreateOrUpdateTag creates or updates an annotated git tag using -fa flag
	CreateOrUpdateTag(tag, message, refSpec string) error

	// DeleteTag deletes a local git tag
	DeleteTag(tag string) error

	// PushTag pushes a tag to remote (without force)
	PushTag(tag string) error

	// PushTagForce pushes a tag to remote with --force flag
	PushTagForce(tag string) error

	// PushTags pushes multiple tags to remote
	PushTags(tags []string) error

	// SetRemoteURL sets the remote URL with authentication token
	SetRemoteURL(token, repo string) error

	// GetCurrentCommit returns the current commit SHA
	GetCurrentCommit() (string, error)

	// TagExists checks if a tag exists locally
	TagExists(tag string) bool
}

// DefaultOperator implements Operator using actual git commands.
type DefaultOperator struct{}

// NewOperator creates a new DefaultOperator.
func NewOperator() *DefaultOperator {
	return &DefaultOperator{}
}

// ConfigureSafeDirectory implements Operator.
func (o *DefaultOperator) ConfigureSafeDirectory(dir string) error {
	return ConfigureSafeDirectory(dir)
}

// ConfigureUser implements Operator.
func (o *DefaultOperator) ConfigureUser(name, email string) error {
	return ConfigureUser(name, email)
}

// CreateTag implements Operator.
func (o *DefaultOperator) CreateTag(tag, message string) error {
	return CreateTag(tag, message)
}

// CreateOrUpdateTag implements Operator.
func (o *DefaultOperator) CreateOrUpdateTag(tag, message, refSpec string) error {
	return CreateOrUpdateTag(tag, message, refSpec)
}

// DeleteTag implements Operator.
func (o *DefaultOperator) DeleteTag(tag string) error {
	return DeleteTag(tag)
}

// PushTag implements Operator.
func (o *DefaultOperator) PushTag(tag string) error {
	return PushTag(tag)
}

// PushTagForce implements Operator.
func (o *DefaultOperator) PushTagForce(tag string) error {
	return PushTagForce(tag)
}

// PushTags implements Operator.
func (o *DefaultOperator) PushTags(tags []string) error {
	return PushTags(tags)
}

// SetRemoteURL implements Operator.
func (o *DefaultOperator) SetRemoteURL(token, repo string) error {
	return SetRemoteURL(token, repo)
}

// GetCurrentCommit implements Operator.
func (o *DefaultOperator) GetCurrentCommit() (string, error) {
	return GetCurrentCommit()
}

// TagExists implements Operator.
func (o *DefaultOperator) TagExists(tag string) bool {
	return TagExists(tag)
}

// MockOperator implements Operator for testing.
type MockOperator struct {
	ConfigureSafeDirectoryFunc func(dir string) error
	ConfigureUserFunc          func(name, email string) error
	CreateTagFunc              func(tag, message string) error
	CreateOrUpdateTagFunc      func(tag, message, refSpec string) error
	DeleteTagFunc              func(tag string) error
	PushTagFunc                func(tag string) error
	PushTagForceFunc           func(tag string) error
	PushTagsFunc               func(tags []string) error
	SetRemoteURLFunc           func(token, repo string) error
	GetCurrentCommitFunc       func() (string, error)
	TagExistsFunc              func(tag string) bool

	// Track calls for assertions
	Calls map[string][]interface{}
}

// NewMockOperator creates a new MockOperator with default no-op implementations.
func NewMockOperator() *MockOperator {
	return &MockOperator{
		Calls:                      make(map[string][]interface{}),
		ConfigureSafeDirectoryFunc: func(dir string) error { return nil },
		ConfigureUserFunc:          func(name, email string) error { return nil },
		CreateTagFunc:              func(tag, message string) error { return nil },
		CreateOrUpdateTagFunc:      func(tag, message, refSpec string) error { return nil },
		DeleteTagFunc:              func(tag string) error { return nil },
		PushTagFunc:                func(tag string) error { return nil },
		PushTagForceFunc:           func(tag string) error { return nil },
		PushTagsFunc:               func(tags []string) error { return nil },
		SetRemoteURLFunc:           func(token, repo string) error { return nil },
		GetCurrentCommitFunc:       func() (string, error) { return "abc123", nil },
		TagExistsFunc:              func(tag string) bool { return false },
	}
}

func (m *MockOperator) recordCall(method string, args ...interface{}) {
	m.Calls[method] = append(m.Calls[method], args)
}

// ConfigureSafeDirectory implements Operator.
func (m *MockOperator) ConfigureSafeDirectory(dir string) error {
	m.recordCall("ConfigureSafeDirectory", dir)
	return m.ConfigureSafeDirectoryFunc(dir)
}

// ConfigureUser implements Operator.
func (m *MockOperator) ConfigureUser(name, email string) error {
	m.recordCall("ConfigureUser", name, email)
	return m.ConfigureUserFunc(name, email)
}

// CreateTag implements Operator.
func (m *MockOperator) CreateTag(tag, message string) error {
	m.recordCall("CreateTag", tag, message)
	return m.CreateTagFunc(tag, message)
}

// CreateOrUpdateTag implements Operator.
func (m *MockOperator) CreateOrUpdateTag(tag, message, refSpec string) error {
	m.recordCall("CreateOrUpdateTag", tag, message, refSpec)
	return m.CreateOrUpdateTagFunc(tag, message, refSpec)
}

// DeleteTag implements Operator.
func (m *MockOperator) DeleteTag(tag string) error {
	m.recordCall("DeleteTag", tag)
	return m.DeleteTagFunc(tag)
}

// PushTag implements Operator.
func (m *MockOperator) PushTag(tag string) error {
	m.recordCall("PushTag", tag)
	return m.PushTagFunc(tag)
}

// PushTagForce implements Operator.
func (m *MockOperator) PushTagForce(tag string) error {
	m.recordCall("PushTagForce", tag)
	return m.PushTagForceFunc(tag)
}

// PushTags implements Operator.
func (m *MockOperator) PushTags(tags []string) error {
	m.recordCall("PushTags", tags)
	return m.PushTagsFunc(tags)
}

// SetRemoteURL implements Operator.
func (m *MockOperator) SetRemoteURL(token, repo string) error {
	m.recordCall("SetRemoteURL", token, repo)
	return m.SetRemoteURLFunc(token, repo)
}

// GetCurrentCommit implements Operator.
func (m *MockOperator) GetCurrentCommit() (string, error) {
	m.recordCall("GetCurrentCommit")
	return m.GetCurrentCommitFunc()
}

// TagExists implements Operator.
func (m *MockOperator) TagExists(tag string) bool {
	m.recordCall("TagExists", tag)
	return m.TagExistsFunc(tag)
}
