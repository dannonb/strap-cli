package config

import (
	"os"
	"strings"
	"testing"
)

func TestDirectoryService_GetCurrentDirectoryName(t *testing.T) {
	service := NewDirectoryService()

	tests := []struct {
		name        string
		setup       func() (cleanup func(), err error)
		expectError bool
		errorContains string
	}{
		{
			name: "normal directory",
			setup: func() (func(), error) {
				// Create a temporary directory and change to it
				tempDir, err := os.MkdirTemp("", "test-dir")
				if err != nil {
					return nil, err
				}
				
				originalDir, err := os.Getwd()
				if err != nil {
					os.RemoveAll(tempDir)
					return nil, err
				}
				
				err = os.Chdir(tempDir)
				if err != nil {
					os.RemoveAll(tempDir)
					return nil, err
				}
				
				return func() {
					os.Chdir(originalDir)
					os.RemoveAll(tempDir)
				}, nil
			},
			expectError: false,
		},
		{
			name: "directory with special characters",
			setup: func() (func(), error) {
				// Create a directory with special characters
				tempDir, err := os.MkdirTemp("", "test-dir-with-spaces and symbols!")
				if err != nil {
					return nil, err
				}
				
				originalDir, err := os.Getwd()
				if err != nil {
					os.RemoveAll(tempDir)
					return nil, err
				}
				
				err = os.Chdir(tempDir)
				if err != nil {
					os.RemoveAll(tempDir)
					return nil, err
				}
				
				return func() {
					os.Chdir(originalDir)
					os.RemoveAll(tempDir)
				}, nil
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup, err := tt.setup()
			if err != nil {
				t.Fatalf("Setup failed: %v", err)
			}
			defer cleanup()

			result, err := service.GetCurrentDirectoryName()
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain '%s', got: %v", tt.errorContains, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result == "" {
					t.Errorf("Expected non-empty result")
				}
			}
		})
	}
}

func TestDirectoryService_SanitizeProjectName(t *testing.T) {
	service := NewDirectoryService()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normal name",
			input:    "my-project",
			expected: "my-project",
		},
		{
			name:     "name with spaces",
			input:    "my project",
			expected: "my-project",
		},
		{
			name:     "name with underscores",
			input:    "my_project",
			expected: "my-project",
		},
		{
			name:     "name with special characters",
			input:    "my@project#with$symbols",
			expected: "myprojectwithsymbols",
		},
		{
			name:     "empty name",
			input:    "",
			expected: "microservice-project",
		},
		{
			name:     "only special characters",
			input:    "@#$%",
			expected: "microservice-project",
		},
		{
			name:     "reserved name",
			input:    "node",
			expected: "node-project",
		},
		{
			name:     "name starting with number",
			input:    "123project",
			expected: "project-123project",
		},
		{
			name:     "very long name",
			input:    strings.Repeat("a", 60),
			expected: strings.Repeat("a", 50),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.SanitizeProjectName(tt.input)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestDirectoryService_ValidateDirectoryForProject(t *testing.T) {
	service := NewDirectoryService()

	tests := []struct {
		name          string
		setup         func() (path string, cleanup func(), err error)
		expectError   bool
		errorContains string
	}{
		{
			name: "valid directory",
			setup: func() (string, func(), error) {
				tempDir, err := os.MkdirTemp("", "valid-dir")
				if err != nil {
					return "", nil, err
				}
				return tempDir, func() { os.RemoveAll(tempDir) }, nil
			},
			expectError: false,
		},
		{
			name: "non-existent directory",
			setup: func() (string, func(), error) {
				return "/non/existent/path", func() {}, nil
			},
			expectError:   true,
			errorContains: "does not exist",
		},
		{
			name: "file instead of directory",
			setup: func() (string, func(), error) {
				tempFile, err := os.CreateTemp("", "test-file")
				if err != nil {
					return "", nil, err
				}
				tempFile.Close()
				return tempFile.Name(), func() { os.Remove(tempFile.Name()) }, nil
			},
			expectError:   true,
			errorContains: "not a directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, cleanup, err := tt.setup()
			if err != nil {
				t.Fatalf("Setup failed: %v", err)
			}
			defer cleanup()

			err = service.ValidateDirectoryForProject(path)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain '%s', got: %v", tt.errorContains, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestDirectoryService_ErrorHandling(t *testing.T) {
	service := &directoryService{}

	t.Run("handleCurrentDirectoryAccessError", func(t *testing.T) {
		tests := []struct {
			name          string
			inputError    error
			expectContains []string
		}{
			{
				name:       "permission denied error",
				inputError: os.ErrPermission,
				expectContains: []string{"permission denied", "fallback", "microservice-project"},
			},
			{
				name:       "generic error",
				inputError: os.ErrInvalid,
				expectContains: []string{"unable to determine", "fallback", "microservice-project"},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := service.handleCurrentDirectoryAccessError(tt.inputError)
				
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				
				if result != "microservice-project" {
					t.Errorf("Expected fallback name 'microservice-project', got '%s'", result)
				}
				
				errMsg := err.Error()
				for _, expected := range tt.expectContains {
					if !strings.Contains(errMsg, expected) {
						t.Errorf("Expected error message to contain '%s', got: %s", expected, errMsg)
					}
				}
			})
		}
	})

	t.Run("handleRootDirectoryCase", func(t *testing.T) {
		result, err := service.handleRootDirectoryCase("/")
		
		if err == nil {
			t.Errorf("Expected error but got none")
		}
		
		if result != "microservice-project" {
			t.Errorf("Expected fallback name 'microservice-project', got '%s'", result)
		}
		
		errMsg := err.Error()
		expectedPhrases := []string{"cannot use root directory", "fallback", "mkdir myproject"}
		for _, phrase := range expectedPhrases {
			if !strings.Contains(errMsg, phrase) {
				t.Errorf("Expected error message to contain '%s', got: %s", phrase, errMsg)
			}
		}
	})
}