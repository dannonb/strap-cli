package validation

import (
	"errors"
	"strings"
	"testing"
)

func TestValidator_DockerErrorHandling(t *testing.T) {
	validator := NewValidator()

	t.Run("HandleDockerNotInstalled", func(t *testing.T) {
		tests := []struct {
			name          string
			inputError    error
			expectedGuidance string
		}{
			{
				name:          "Executable not found",
				inputError:    errors.New("executable file not found in $PATH"),
				expectedGuidance: "Docker Installation Guide:",
			},
			{
				name:          "Permission denied on command",
				inputError:    errors.New("permission denied"),
				expectedGuidance: "Docker Permission Solutions:",
			},
			{
				name:          "Access denied",
				inputError:    errors.New("access denied"),
				expectedGuidance: "Docker Access Solutions:",
			},
			{
				name:          "Generic error",
				inputError:    errors.New("some other error"),
				expectedGuidance: "General Docker Setup:",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := validator.handleDockerNotInstalled(tt.inputError)

				if err == nil {
					t.Error("Expected error, got nil")
				}

				if !strings.Contains(err.Error(), tt.expectedGuidance) {
					t.Errorf("Expected error to contain '%s', got: %s", tt.expectedGuidance, err.Error())
				}

				// Verify comprehensive guidance is provided
				if !strings.Contains(err.Error(), "🐳") {
					t.Error("Expected Docker emoji in error message")
				}

				// Only check for documentation link in generic error case
				if tt.name == "Generic error" {
					if !strings.Contains(err.Error(), "https://docs.docker.com") {
						t.Error("Expected Docker documentation link")
					}
				}
			})
		}
	})

	t.Run("HandleDockerNotRunning", func(t *testing.T) {
		tests := []struct {
			name          string
			inputError    error
			expectedGuidance string
		}{
			{
				name:          "Cannot connect to daemon",
				inputError:    errors.New("Cannot connect to the Docker daemon at unix:///var/run/docker.sock"),
				expectedGuidance: "Docker Startup Solutions:",
			},
			{
				name:          "Permission denied to daemon",
				inputError:    errors.New("permission denied while trying to connect"),
				expectedGuidance: "Docker Permission Solutions:",
			},
			{
				name:          "Timeout error",
				inputError:    errors.New("timeout waiting for response"),
				expectedGuidance: "Docker Timeout Solutions:",
			},
			{
				name:          "Context deadline exceeded",
				inputError:    errors.New("context deadline exceeded"),
				expectedGuidance: "Docker Timeout Solutions:",
			},
			{
				name:          "Socket connection error",
				inputError:    errors.New("dial unix /var/run/docker.sock: connect: no such file or directory"),
				expectedGuidance: "Docker Socket Solutions:",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := validator.handleDockerNotRunning(tt.inputError)

				if err == nil {
					t.Error("Expected error, got nil")
				}

				if !strings.Contains(err.Error(), tt.expectedGuidance) {
					t.Errorf("Expected error to contain '%s', got: %s", tt.expectedGuidance, err.Error())
				}

				// Verify platform-specific guidance
				if strings.Contains(tt.inputError.Error(), "daemon") {
					if !strings.Contains(err.Error(), "Docker Desktop") {
						t.Error("Expected Docker Desktop guidance for daemon errors")
					}
				}
			})
		}
	})

	t.Run("HandleDockerComposeNotFound", func(t *testing.T) {
		tests := []struct {
			name          string
			modernError   error
			legacyError   error
			expectedGuidance string
		}{
			{
				name:          "Daemon not running",
				modernError:   errors.New("Cannot connect to the Docker daemon"),
				legacyError:   nil,
				expectedGuidance: "Docker Daemon Solutions:",
			},
			{
				name:          "Permission denied",
				modernError:   errors.New("permission denied"),
				legacyError:   nil,
				expectedGuidance: "Docker Compose Permission Solutions:",
			},
			{
				name:          "Unknown command",
				modernError:   errors.New("unknown command: compose"),
				legacyError:   errors.New("executable file not found"),
				expectedGuidance: "Docker Compose Installation:",
			},
			{
				name:          "Version compatibility",
				modernError:   errors.New("version compatibility issue"),
				legacyError:   nil,
				expectedGuidance: "Docker Compose Version Solutions:",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := validator.handleDockerComposeNotFound(tt.modernError, tt.legacyError)

				if err == nil {
					t.Error("Expected error, got nil")
				}

				if !strings.Contains(err.Error(), tt.expectedGuidance) {
					t.Errorf("Expected error to contain '%s', got: %s", tt.expectedGuidance, err.Error())
				}

				// Verify comprehensive guidance
				if !strings.Contains(err.Error(), "🐳") {
					t.Error("Expected Docker emoji in error message")
				}
			})
		}
	})

	t.Run("CheckDockerPrerequisites_Integration", func(t *testing.T) {
		// This test verifies the integration of Docker error handling
		// Note: This will fail in environments without Docker, which is expected
		
		err := validator.CheckExecutionPrerequisites()
		
		// If Docker is not available, verify error contains helpful guidance
		if err != nil {
			errorMsg := err.Error()
			
			// Should contain helpful guidance
			hasGuidance := strings.Contains(errorMsg, "🐳") ||
						  strings.Contains(errorMsg, "Docker") ||
						  strings.Contains(errorMsg, "https://docs.docker.com")
			
			if !hasGuidance {
				t.Error("Docker error should contain helpful guidance")
			}
		}
	})
}

func TestValidator_EnhancedErrorMessages(t *testing.T) {
	validator := NewValidator()

	t.Run("ComprehensiveDockerGuidance", func(t *testing.T) {
		// Test that all Docker error scenarios provide comprehensive guidance
		testErrors := []error{
			errors.New("executable file not found in $PATH"),
			errors.New("Cannot connect to the Docker daemon"),
			errors.New("permission denied"),
			errors.New("timeout waiting for response"),
			errors.New("unknown command: compose"),
		}

		for _, testErr := range testErrors {
			var resultErr error
			
			// Test different error handling paths
			if strings.Contains(testErr.Error(), "executable file not found") {
				resultErr = validator.handleDockerNotInstalled(testErr)
			} else if strings.Contains(testErr.Error(), "Cannot connect") {
				resultErr = validator.handleDockerNotRunning(testErr)
			} else if strings.Contains(testErr.Error(), "unknown command") {
				resultErr = validator.handleDockerComposeNotFound(testErr, nil)
			} else {
				resultErr = validator.handleDockerNotRunning(testErr)
			}

			if resultErr == nil {
				t.Errorf("Expected error for input: %v", testErr)
				continue
			}

			errorMsg := resultErr.Error()

			// Verify comprehensive guidance elements
			checks := []struct {
				element string
				description string
				required bool
			}{
				{"🐳", "Docker emoji for visual identification", true},
				{"Solutions:", "Solution section header", false},
				{"•", "Bullet points for readability", true},
				{"docker", "Docker command references", true},
				{"Original error:", "Original error preservation", false},
			}

			for _, check := range checks {
				if check.required && !strings.Contains(errorMsg, check.element) {
					t.Errorf("Error message missing %s: %s", check.description, errorMsg)
				}
			}
		}
	})

	t.Run("PlatformSpecificGuidance", func(t *testing.T) {
		// Test that errors provide platform-specific guidance where appropriate
		permissionErr := errors.New("permission denied")
		
		err := validator.handleDockerNotInstalled(permissionErr)
		errorMsg := err.Error()

		// Should contain both Linux and Windows/Mac guidance
		if !strings.Contains(errorMsg, "Linux:") {
			t.Error("Expected Linux-specific guidance")
		}
		
		if !strings.Contains(errorMsg, "Windows/Mac:") {
			t.Error("Expected Windows/Mac-specific guidance")
		}

		if !strings.Contains(errorMsg, "usermod -aG docker") {
			t.Error("Expected Linux docker group guidance")
		}
	})

	t.Run("RecoveryStepsProgression", func(t *testing.T) {
		// Test that error messages provide logical progression of recovery steps
		daemonErr := errors.New("Cannot connect to the Docker daemon")
		
		err := validator.handleDockerNotRunning(daemonErr)
		errorMsg := err.Error()

		// Should provide steps in logical order
		expectedSteps := []string{
			"Start Docker Desktop",
			"Wait",
			"Verify",
		}

		lastIndex := -1
		for _, step := range expectedSteps {
			index := strings.Index(errorMsg, step)
			if index == -1 {
				t.Errorf("Expected step '%s' not found in error message", step)
				continue
			}
			if index <= lastIndex {
				t.Errorf("Steps not in logical order: '%s' should come after previous step", step)
			}
			lastIndex = index
		}
	})
}