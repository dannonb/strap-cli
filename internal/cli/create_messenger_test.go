package cli

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"microservice-bootstrapper/internal/feedback"
)

// TestResolveProjectNameWithDetails tests the enhanced project name resolution
func TestResolveProjectNameWithDetails(t *testing.T) {
	tests := []struct {
		name         string
		providedName string
		expectError  bool
	}{
		{
			name:         "explicit name provided",
			providedName: "my-project",
			expectError:  false,
		},
		{
			name:         "empty name - should infer from directory",
			providedName: "",
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolvedName, originalDir, err := resolveProjectNameWithDetails(tt.providedName)
			
			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			
			if tt.providedName != "" {
				// When name is provided, should return that name
				if resolvedName != tt.providedName {
					t.Errorf("Expected resolved name %q, got %q", tt.providedName, resolvedName)
				}
				// Original directory should be empty when name is provided
				if originalDir != "" {
					t.Errorf("Expected empty original directory when name provided, got %q", originalDir)
				}
			} else {
				// When name is not provided, should have resolved name and original directory
				if resolvedName == "" {
					t.Errorf("Expected non-empty resolved name when inferring from directory")
				}
				if originalDir == "" {
					t.Errorf("Expected non-empty original directory when inferring from directory")
				}
			}
		})
	}
}

// TestMessengerIntegrationInCLI tests that the messenger is properly integrated into the CLI
func TestMessengerIntegrationInCLI(t *testing.T) {
	// This test verifies that the messenger interface is properly used
	// We can't easily test the full CLI execution without mocking the generator,
	// but we can test the messenger creation and basic functionality
	
	messenger := feedback.NewMessenger()
	if messenger == nil {
		t.Fatal("Failed to create messenger")
	}
	
	// Test that messenger methods can be called without panicking
	var buf bytes.Buffer
	testMessenger := feedback.NewMessengerWithOutput(&buf)
	
	// These should not panic
	testMessenger.ShowProjectNameInference("test-project", "test-dir")
	testMessenger.ShowDockerWarning()
	testMessenger.ShowGenerationSuccess("test-project", ".")
	
	output := buf.String()
	if output == "" {
		t.Error("Expected messenger to produce output")
	}
	
	// Verify key messages are present
	expectedMessages := []string{
		"📁 Using directory name as project name",
		"⚠️  Docker Warning:",
		"✅ Successfully created microservice project",
	}
	
	for _, msg := range expectedMessages {
		if !strings.Contains(output, msg) {
			t.Errorf("Expected output to contain %q, got: %s", msg, output)
		}
	}
}

// TestMessengerErrorHandling tests that messenger integrates well with error handling
func TestMessengerErrorHandling(t *testing.T) {
	// This test verifies that the messenger doesn't interfere with error handling
	var buf bytes.Buffer
	messenger := feedback.NewMessengerWithOutput(&buf)
	
	// Test that messenger methods work even when called in error scenarios
	messenger.ShowDockerWarning()
	
	output := buf.String()
	if !strings.Contains(output, "⚠️  Docker Warning:") {
		t.Error("Expected Docker warning to be shown")
	}
}

// TestCLIFlagValidation tests that CLI flags are properly validated
func TestCLIFlagValidation(t *testing.T) {
	// Save original stdout
	originalStdout := os.Stdout
	defer func() { os.Stdout = originalStdout }()

	// Create a pipe to capture output
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Test validation with invalid backend
	backend = "invalid-backend"
	frontend = ""
	database = ""
	projectName = "test-project"
	force = false

	err := validateFlags(nil, []string{})
	
	// Close writer and restore stdout
	w.Close()
	os.Stdout = originalStdout

	// Read captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)

	if err == nil {
		t.Error("Expected validation error for invalid backend")
	}

	if !strings.Contains(err.Error(), "Invalid backend") {
		t.Errorf("Expected validation error to mention invalid backend, got: %v", err)
	}

	// Reset flags for other tests
	backend = ""
	frontend = ""
	database = ""
	projectName = ""
	force = false
}