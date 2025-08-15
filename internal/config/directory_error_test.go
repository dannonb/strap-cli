package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDirectoryService_EnhancedErrorHandling(t *testing.T) {
	service := NewDirectoryService()

	t.Run("HandleCurrentDirectoryAccessError", func(t *testing.T) {
		tests := []struct {
			name          string
			inputError    error
			expectedName  string
			expectedError string
		}{
			{
				name:          "Permission denied error",
				inputError:    os.ErrPermission,
				expectedName:  "microservice-project",
				expectedError: "permission restrictions",
			},
			{
				name:          "Directory not found error",
				inputError:    &testError{msg: "no such file or directory"},
				expectedName:  "microservice-project",
				expectedError: "no longer exists",
			},
			{
				name:          "Generic error",
				inputError:    os.ErrInvalid,
				expectedName:  "microservice-project",
				expectedError: "unable to determine",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ds := &directoryService{}
				name, err := ds.handleCurrentDirectoryAccessError(tt.inputError)

				if name != tt.expectedName {
					t.Errorf("Expected name %s, got %s", tt.expectedName, name)
				}

				if err == nil {
					t.Error("Expected error, got nil")
				}

				if !strings.Contains(err.Error(), tt.expectedError) {
					t.Errorf("Expected error to contain '%s', got: %s", tt.expectedError, err.Error())
				}

				// Verify recovery guidance is included
				if !strings.Contains(err.Error(), "🔧 Recovery Steps:") {
					t.Error("Expected recovery guidance in error message")
				}
			})
		}
	})

	t.Run("HandleRootDirectoryCase", func(t *testing.T) {
		ds := &directoryService{}
		name, err := ds.handleRootDirectoryCase("/")

		if name != "microservice-project" {
			t.Errorf("Expected fallback name, got %s", name)
		}

		if err == nil {
			t.Error("Expected error, got nil")
		}

		if !strings.Contains(err.Error(), "Root Directory Solutions:") {
			t.Error("Expected root directory solutions in error message")
		}

		if !strings.Contains(err.Error(), "mkdir myproject") {
			t.Error("Expected mkdir guidance in error message")
		}
	})

	t.Run("ValidateWritePermissions_ErrorHandling", func(t *testing.T) {
		// Create a temporary directory for testing
		tempDir, err := os.MkdirTemp("", "strap-test")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		// Test with a non-existent directory
		nonExistentDir := filepath.Join(tempDir, "nonexistent")
		ds := &directoryService{}
		err = ds.validateWritePermissions(nonExistentDir)

		if err == nil {
			t.Error("Expected error for non-existent directory")
		}

		// Verify error contains helpful guidance
		if !strings.Contains(err.Error(), "🔧") {
			t.Error("Expected recovery guidance in error message")
		}
	})

	t.Run("SanitizeProjectName_EdgeCases", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			expected string
		}{
			{
				name:     "Empty string",
				input:    "",
				expected: "microservice-project",
			},
			{
				name:     "Only special characters",
				input:    "!@#$%^&*()",
				expected: "microservice-project",
			},
			{
				name:     "Reserved name",
				input:    "con",
				expected: "con-project",
			},
			{
				name:     "Starts with number",
				input:    "123project",
				expected: "project-123project",
			},
			{
				name:     "Very long name",
				input:    strings.Repeat("a", 100),
				expected: strings.Repeat("a", 50),
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := service.SanitizeProjectName(tt.input)
				if result != tt.expected {
					t.Errorf("Expected %s, got %s", tt.expected, result)
				}
			})
		}
	})

	t.Run("ValidateDirectoryForProject_ErrorScenarios", func(t *testing.T) {
		ds := &directoryService{}

		// Test non-existent directory
		err := ds.ValidateDirectoryForProject("/nonexistent/path")
		if err == nil {
			t.Error("Expected error for non-existent directory")
		}
		if !strings.Contains(err.Error(), "Directory Creation Solutions:") {
			t.Error("Expected directory creation guidance")
		}

		// Test with a file instead of directory
		tempFile, err := os.CreateTemp("", "test-file")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tempFile.Name())
		tempFile.Close()

		err = ds.ValidateDirectoryForProject(tempFile.Name())
		if err == nil {
			t.Error("Expected error when path is a file")
		}
		if !strings.Contains(err.Error(), "not a directory") {
			t.Error("Expected 'not a directory' error message")
		}
	})
}

func TestDirectoryService_FallbackStrategies(t *testing.T) {
	service := NewDirectoryService()

	t.Run("HandleEmptySanitizedName", func(t *testing.T) {
		ds := &directoryService{}
		result := ds.handleEmptySanitizedName("!@#$%")

		if result != "microservice-project" {
			t.Errorf("Expected fallback name, got %s", result)
		}
	})

	t.Run("HandleEmptyDirectoryName", func(t *testing.T) {
		ds := &directoryService{}
		name, err := ds.handleEmptyDirectoryName("/some/path")

		if name != "microservice-project" {
			t.Errorf("Expected fallback name, got %s", name)
		}

		if err == nil {
			t.Error("Expected error, got nil")
		}

		if !strings.Contains(err.Error(), "Empty Directory Name Solutions:") {
			t.Error("Expected empty directory name solutions")
		}
	})

	t.Run("ReservedNamesHandling", func(t *testing.T) {
		reservedNames := []string{"con", "prn", "aux", "nul", "node", "src", "test"}
		
		for _, reserved := range reservedNames {
			result := service.SanitizeProjectName(reserved)
			if result == reserved {
				t.Errorf("Reserved name %s was not modified", reserved)
			}
			if !strings.Contains(result, "project") {
				t.Errorf("Expected reserved name %s to include 'project', got %s", reserved, result)
			}
		}
	})
}

func TestDirectoryService_ComprehensiveErrorGuidance(t *testing.T) {
	ds := &directoryService{}

	t.Run("NetworkErrorHandling", func(t *testing.T) {
		networkErr := &os.PathError{
			Op:   "stat",
			Path: "/network/path",
			Err:  os.ErrInvalid,
		}
		// Simulate network error by modifying error string
		networkErr.Err = &testError{msg: "network error occurred"}

		name, err := ds.handleCurrentDirectoryAccessError(networkErr)

		if name != "microservice-project" {
			t.Errorf("Expected fallback name, got %s", name)
		}

		if !strings.Contains(err.Error(), "network connection") {
			t.Error("Expected network-specific guidance")
		}
	})

	t.Run("PathLengthErrorHandling", func(t *testing.T) {
		longPath := "/" + strings.Repeat("very-long-directory-name/", 20)
		err := ds.validateDirectoryStructure(longPath)

		if err == nil {
			t.Error("Expected error for very long path")
		}

		if !strings.Contains(err.Error(), "path is very long") {
			t.Error("Expected path length warning")
		}
	})

	t.Run("InvalidCharactersHandling", func(t *testing.T) {
		invalidPath := "/path/with<invalid>chars"
		err := ds.validateDirectoryStructure(invalidPath)

		if err == nil {
			t.Error("Expected error for invalid characters")
		}

		if !strings.Contains(err.Error(), "invalid characters") {
			t.Error("Expected invalid characters warning")
		}
	})
}

// testError is a helper for testing specific error conditions
type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}