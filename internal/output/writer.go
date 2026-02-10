// Copyright (c) 2024 Inetum Poland.

// Package output provides functionality for writing GitHub Actions outputs.
// It handles both single-value and multiline outputs using the GITHUB_OUTPUT
// file format introduced in GitHub Actions Node.js 16+.
//
// Output Format:
//   - Single value: key=value
//   - Multiline: key<<EOF\nvalue\nEOF
//
// The package falls back to the legacy ::set-output command format when
// GITHUB_OUTPUT is not set (useful for local testing).
package output

import (
	"fmt"
	"os"
)

// Writer handles GitHub Actions outputs
type Writer struct {
	outputFile string
}

// NewWriter creates a new output writer
func NewWriter() *Writer {
	return &Writer{
		outputFile: os.Getenv("GITHUB_OUTPUT"),
	}
}

// Set writes a key-value pair to GitHub Actions output
func (w *Writer) Set(key, value string) error {
	if w.outputFile == "" {
		// For local testing or when GITHUB_OUTPUT is not set
		fmt.Printf("::set-output name=%s::%s\n", key, value)
		return nil
	}

	f, err := os.OpenFile(w.outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("failed to open output file: %w", err)
	}
	defer f.Close()

	// GitHub Actions output format: key=value
	_, err = fmt.Fprintf(f, "%s=%s\n", key, value)
	if err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	return nil
}

// SetMultiline writes a multi-line value to GitHub Actions output
func (w *Writer) SetMultiline(key, value string) error {
	if w.outputFile == "" {
		// For local testing
		fmt.Printf("::set-output name=%s::%s\n", key, value)
		return nil
	}

	f, err := os.OpenFile(w.outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("failed to open output file: %w", err)
	}
	defer f.Close()

	// GitHub Actions multiline format:
	// key<<EOF
	// value
	// EOF
	_, err = fmt.Fprintf(f, "%s<<EOF\n%s\nEOF\n", key, value)
	if err != nil {
		return fmt.Errorf("failed to write multiline output: %w", err)
	}

	return nil
}

// SetAll sets multiple outputs at once
func (w *Writer) SetAll(outputs map[string]string) error {
	for key, value := range outputs {
		if err := w.Set(key, value); err != nil {
			return err
		}
	}
	return nil
}

// WriteSummary writes markdown content to GITHUB_STEP_SUMMARY
// This appears as a summary in the GitHub Actions UI
func (w *Writer) WriteSummary(markdown string) error {
	summaryFile := os.Getenv("GITHUB_STEP_SUMMARY")
	if summaryFile == "" {
		// For local testing, print to stdout
		fmt.Println("--- STEP SUMMARY ---")
		fmt.Println(markdown)
		fmt.Println("--- END SUMMARY ---")
		return nil
	}

	f, err := os.OpenFile(summaryFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("failed to open summary file: %w", err)
	}
	defer f.Close()

	_, err = fmt.Fprintln(f, markdown)
	if err != nil {
		return fmt.Errorf("failed to write summary: %w", err)
	}

	return nil
}
