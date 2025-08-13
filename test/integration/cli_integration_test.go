package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLIIntegration(t *testing.T) {
	// Build the CLI binary for testing
	binaryPath := buildCLIBinary(t)
	defer os.Remove(binaryPath)

	tests := []struct {
		name        string
		args        []string
		expectError bool
		checkOutput []string
		checkFiles  []string
	}{
		{
			name:        "help command",
			args:        []string{"--help"},
			expectError: false,
			checkOutput: []string{"Microservice Bootstrapper", "Available Commands", "create"},
		},
		{
			name:        "version command",
			args:        []string{"version"},
			expectError: false,
			checkOutput: []string{"Microservice Bootstrapper", "Build Information", "Supported Technologies"},
		},
		{
			name:        "examples command",
			args:        []string{"examples"},
			expectError: false,
			checkOutput: []string{"FULL-STACK APPLICATIONS", "API SERVICES", "FRONTEND APPLICATIONS"},
		},
		{
			name:        "workflows command",
			args:        []string{"workflows"},
			expectError: false,
			checkOutput: []string{"QUICK START WORKFLOW", "DEVELOPMENT WORKFLOW", "docker-compose up"},
		},
		{
			name:        "create help",
			args:        []string{"create", "--help"},
			expectError: false,
			checkOutput: []string{"Create a new microservice project", "--be", "--fe", "--db"},
		},
		{
			name:        "invalid backend",
			args:        []string{"create", "--be=invalid"},
			expectError: true,
			checkOutput: []string{"Invalid backend", "Supported backends"},
		},
		{
			name:        "invalid frontend",
			args:        []string{"create", "--fe=invalid"},
			expectError: true,
			checkOutput: []string{"Invalid frontend", "Supported frontends"},
		},
		{
			name:        "invalid database",
			args:        []string{"create", "--db=invalid"},
			expectError: true,
			checkOutput: []string{"Invalid database", "Supported databases"},
		},
		{
			name:        "no services specified",
			args:        []string{"create"},
			expectError: true,
			checkOutput: []string{"At least one service must be specified"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory for this test
			tempDir := t.TempDir()
			originalDir, _ := os.Getwd()
			defer os.Chdir(originalDir)

			// Change to temp directory
			os.Chdir(tempDir)

			// Run the CLI command
			cmd := exec.Command(binaryPath, tt.args...)
			output, err := cmd.CombinedOutput()
			outputStr := string(output)

			// Check error expectation
			if tt.expectError && err == nil {
				t.Errorf("Expected error but command succeeded. Output: %s", outputStr)
			} else if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v. Output: %s", err, outputStr)
			}

			// Check output content
			for _, expected := range tt.checkOutput {
				if !strings.Contains(outputStr, expected) {
					t.Errorf("Output does not contain expected text '%s'. Full output: %s", expected, outputStr)
				}
			}

			// Check created files
			for _, file := range tt.checkFiles {
				filePath := filepath.Join(tempDir, file)
				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					t.Errorf("Expected file %s was not created", file)
				}
			}
		})
	}
}

func TestCLIProjectGeneration(t *testing.T) {
	// Build the CLI binary for testing
	binaryPath := buildCLIBinary(t)
	defer os.Remove(binaryPath)

	tests := []struct {
		name        string
		args        []string
		projectName string
		checkFiles  []string
		checkDirs   []string
	}{
		{
			name:        "fastapi backend only",
			args:        []string{"create", "--be=fastapi", "--name=fastapi-test", "--force"},
			projectName: "fastapi-test",
			checkFiles: []string{
				"fastapi-test/docker-compose.yml",
				"fastapi-test/README.md",
				"fastapi-test/backend/Dockerfile",
				"fastapi-test/backend/main.py",
			},
			checkDirs: []string{
				"fastapi-test/backend",
				"fastapi-test/backend/app",
			},
		},
		{
			name:        "react frontend only",
			args:        []string{"create", "--fe=react", "--name=react-test", "--force"},
			projectName: "react-test",
			checkFiles: []string{
				"react-test/docker-compose.yml",
				"react-test/README.md",
				"react-test/frontend/Dockerfile",
				"react-test/frontend/package.json",
			},
			checkDirs: []string{
				"react-test/frontend",
				"react-test/frontend/src",
			},
		},
		{
			name:        "full stack application",
			args:        []string{"create", "--be=gin", "--fe=vue", "--db=postgres", "--name=fullstack-test", "--force"},
			projectName: "fullstack-test",
			checkFiles: []string{
				"fullstack-test/docker-compose.yml",
				"fullstack-test/README.md",
				"fullstack-test/backend/Dockerfile",
				"fullstack-test/backend/main.go",
				"fullstack-test/frontend/Dockerfile",
				"fullstack-test/frontend/package.json",
			},
			checkDirs: []string{
				"fullstack-test/backend",
				"fullstack-test/frontend",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory for this test
			tempDir := t.TempDir()
			originalDir, _ := os.Getwd()
			defer os.Chdir(originalDir)

			// Change to temp directory
			os.Chdir(tempDir)

			// Run the CLI command
			cmd := exec.Command(binaryPath, tt.args...)
			output, err := cmd.CombinedOutput()
			outputStr := string(output)

			if err != nil {
				t.Fatalf("CLI command failed: %v. Output: %s", err, outputStr)
			}

			// Check that success message is displayed
			if !strings.Contains(outputStr, "Successfully created") {
				t.Errorf("Success message not found in output: %s", outputStr)
			}

			// Check created files
			for _, file := range tt.checkFiles {
				filePath := filepath.Join(tempDir, file)
				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					t.Errorf("Expected file %s was not created", file)
				}
			}

			// Check created directories
			for _, dir := range tt.checkDirs {
				dirPath := filepath.Join(tempDir, dir)
				if stat, err := os.Stat(dirPath); os.IsNotExist(err) || !stat.IsDir() {
					t.Errorf("Expected directory %s was not created or is not a directory", dir)
				}
			}
		})
	}
}

func TestCLIErrorHandling(t *testing.T) {
	// Build the CLI binary for testing
	binaryPath := buildCLIBinary(t)
	defer os.Remove(binaryPath)

	tests := []struct {
		name         string
		args         []string
		expectError  bool
		errorMessage string
	}{
		{
			name:         "project name with spaces",
			args:         []string{"create", "--be=fastapi", "--name=test project", "--force"},
			expectError:  true,
			errorMessage: "cannot contain spaces",
		},
		{
			name:         "project name starting with dash",
			args:         []string{"create", "--be=fastapi", "--name=-test", "--force"},
			expectError:  true,
			errorMessage: "cannot start with - or _",
		},
		{
			name:         "multiple invalid flags",
			args:         []string{"create", "--be=invalid", "--fe=invalid", "--force"},
			expectError:  true,
			errorMessage: "Invalid backend",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory for this test
			tempDir := t.TempDir()
			originalDir, _ := os.Getwd()
			defer os.Chdir(originalDir)

			// Change to temp directory
			os.Chdir(tempDir)

			// Run the CLI command
			cmd := exec.Command(binaryPath, tt.args...)
			output, err := cmd.CombinedOutput()
			outputStr := string(output)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but command succeeded. Output: %s", outputStr)
				}
				if !strings.Contains(outputStr, tt.errorMessage) {
					t.Errorf("Expected error message '%s' not found in output: %s", tt.errorMessage, outputStr)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v. Output: %s", err, outputStr)
				}
			}
		})
	}
}

func TestCLIDirectoryConflicts(t *testing.T) {
	// Build the CLI binary for testing
	binaryPath := buildCLIBinary(t)
	defer os.Remove(binaryPath)

	// Create temporary directory for this test
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	// Change to temp directory
	os.Chdir(tempDir)

	// Create a directory with some files
	projectDir := filepath.Join(tempDir, "existing-project")
	os.MkdirAll(projectDir, 0755)
	os.WriteFile(filepath.Join(projectDir, "existing-file.txt"), []byte("content"), 0644)

	t.Run("directory not empty without force", func(t *testing.T) {
		// Try to create project in non-empty directory without --force
		cmd := exec.Command(binaryPath, "create", "--be=fastapi", "--name=existing-project")
		output, err := cmd.CombinedOutput()
		outputStr := string(output)

		if err == nil {
			t.Errorf("Expected error for non-empty directory but command succeeded. Output: %s", outputStr)
		}

		if !strings.Contains(outputStr, "not empty") {
			t.Errorf("Expected 'not empty' error message not found in output: %s", outputStr)
		}
	})

	t.Run("directory not empty with force", func(t *testing.T) {
		// Try to create project in non-empty directory with --force
		cmd := exec.Command(binaryPath, "create", "--be=fastapi", "--name=existing-project", "--force")
		output, err := cmd.CombinedOutput()
		outputStr := string(output)

		if err != nil {
			t.Errorf("Unexpected error with --force flag: %v. Output: %s", err, outputStr)
		}

		if !strings.Contains(outputStr, "Successfully created") {
			t.Errorf("Success message not found with --force flag. Output: %s", outputStr)
		}
	})
}

// buildCLIBinary builds the CLI binary for testing and returns the path
func buildCLIBinary(t *testing.T) string {
	// Get the project root directory
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	// Navigate to project root (assuming we're in test/integration)
	projectRoot := filepath.Join(wd, "..", "..")
	
	// Create temporary binary path
	binaryPath := filepath.Join(t.TempDir(), "strap-test")
	if os.Getenv("GOOS") == "windows" || strings.Contains(os.Getenv("OS"), "Windows") {
		binaryPath += ".exe"
	}

	// Build the binary
	cmd := exec.Command("go", "build", "-o", binaryPath, "./cmd")
	cmd.Dir = projectRoot
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build CLI binary: %v. Output: %s", err, string(output))
	}

	return binaryPath
}

func TestCLIBinaryExists(t *testing.T) {
	// Test that we can build the CLI binary
	binaryPath := buildCLIBinary(t)
	defer os.Remove(binaryPath)

	// Check that binary exists and is executable
	stat, err := os.Stat(binaryPath)
	if err != nil {
		t.Fatalf("Built binary does not exist: %v", err)
	}

	if stat.IsDir() {
		t.Fatal("Built binary is a directory, not a file")
	}

	// Try to run the binary with --help
	cmd := exec.Command(binaryPath, "--help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run built binary: %v. Output: %s", err, string(output))
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "Microservice Bootstrapper") {
		t.Errorf("Binary help output does not contain expected text. Output: %s", outputStr)
	}
}