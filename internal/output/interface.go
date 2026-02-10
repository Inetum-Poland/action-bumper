// Copyright (c) 2024 Inetum Poland.

package output

// WriterInterface defines the interface for output operations.
// This allows for mocking in tests.
type WriterInterface interface {
	// Set writes a key-value pair to GitHub Actions output
	Set(key, value string) error

	// SetMultiline writes a multi-line value to GitHub Actions output
	SetMultiline(key, value string) error

	// SetAll sets multiple outputs at once
	SetAll(outputs map[string]string) error
}

// Ensure Writer implements WriterInterface
var _ WriterInterface = (*Writer)(nil)

// MockWriter implements WriterInterface for testing.
type MockWriter struct {
	SetFunc          func(key, value string) error
	SetMultilineFunc func(key, value string) error
	SetAllFunc       func(outputs map[string]string) error

	// Captured outputs for assertions
	Outputs map[string]string

	// Track calls for assertions
	Calls map[string][]interface{}
}

// NewMockWriter creates a new MockWriter.
func NewMockWriter() *MockWriter {
	m := &MockWriter{
		Outputs: make(map[string]string),
		Calls:   make(map[string][]interface{}),
	}

	// Default implementations that capture outputs
	m.SetFunc = func(key, value string) error {
		m.Outputs[key] = value
		return nil
	}
	m.SetMultilineFunc = func(key, value string) error {
		m.Outputs[key] = value
		return nil
	}
	m.SetAllFunc = func(outputs map[string]string) error {
		for k, v := range outputs {
			m.Outputs[k] = v
		}
		return nil
	}

	return m
}

func (m *MockWriter) recordCall(method string, args ...interface{}) {
	m.Calls[method] = append(m.Calls[method], args)
}

// Set implements WriterInterface.
func (m *MockWriter) Set(key, value string) error {
	m.recordCall("Set", key, value)
	return m.SetFunc(key, value)
}

// SetMultiline implements WriterInterface.
func (m *MockWriter) SetMultiline(key, value string) error {
	m.recordCall("SetMultiline", key, value)
	return m.SetMultilineFunc(key, value)
}

// SetAll implements WriterInterface.
func (m *MockWriter) SetAll(outputs map[string]string) error {
	m.recordCall("SetAll", outputs)
	return m.SetAllFunc(outputs)
}

// GetOutput returns a captured output value.
func (m *MockWriter) GetOutput(key string) string {
	return m.Outputs[key]
}

// HasOutput checks if an output was set.
func (m *MockWriter) HasOutput(key string) bool {
	_, ok := m.Outputs[key]
	return ok
}

// Reset clears all captured outputs and calls.
func (m *MockWriter) Reset() {
	m.Outputs = make(map[string]string)
	m.Calls = make(map[string][]interface{})
}
