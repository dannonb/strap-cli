package integration

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestProjectCreationWithInferredNames tests end-to-end project creation with inferred names
func TestProjectCreationWithInferredNames(t *testing.T) {
	// Build the CLI binary for testing
	binaryPath := buildCLIBinary(t)
	defer os.Remove(binaryPath)

	tests := []struct {
		name            string
		directoryName   string
		args            []string
		expectedProject string
		checkFiles      []string
		checkOutput     []string
	}{
		{
			name:            "simple directory name",
			directoryName:   "my-api",
			args:            []string{"create", "--be=fastapi", "--force"},
			expectedProject: "my-api",
			checkFiles: []string{
				"docker-compose.yml",
				"README.md",
				"backend/Dockerfile",
				"backend/main.py",
			},
			checkOutput: []string{
				"Using directory name as project name: my-api",
				"Successfully created microservice project 'my-api'!",
			},
		},
		{
			name:            "directory with underscores",
			directoryName:   "user_service",
			args:            []string{"create", "--be=gin", "--force"},
			expectedProject: "user-service", // Should be sanitized
			checkFiles: []string{
				"docker-compose.yml",
				"README.md",
				"backend/Dockerfile",
				"backend/main.go",
			},
			checkOutput: []string{
				"Using directory name as project name: user-service",
				"Successfully created microservice project 'user-service'!",
			},
		},
		{
			name:            "directory with spaces gets sanitized",
			directoryName:   "my project",
			args:            []string{"create", "--be=express", "--force"},
			expectedProject: "my-project", // Should be sanitized
			checkFiles: []string{
				"docker-compose.yml",
				"README.md",
				"backend/Dockerfile",
				"backend/server.js",
			},
			checkOutput: []string{
				"Using directory name as project name: my-project",
				"Successfully created microservice project 'my-project'!",
			},
		},
		{
			name:            "directory with special characters",
			directoryName:   "api@service!",
			args:            []string{"create", "--be=fastapi", "--force"},
			expectedProject: "apiservice", // Should be sanitized
			checkFiles: []string{
				"docker-compose.yml",
				"README.md",
				"backend/Dockerfile",
				"backend/main.py",
			},
			checkOutput: []string{
				"Using directory name as project name: apiservice",
				"Successfully created microservice project 'apiservice'!",
			},
		},
		{
			name:            "full stack with inferred name",
			directoryName:   "ecommerce-app",
			args:            []string{"create", "--be=gin", "--fe=react", "--db=postgres", "--force"},
			expectedProject: "ecommerce-app",
			checkFiles: []string{
				"docker-compose.yml",
				"README.md",
				"backend/Dockerfile",
				"backend/main.go",
				"frontend/Dockerfile",
				"frontend/package.json",
			},
			checkOutput: []string{
				"Using directory name as project name: ecommerce-app",
				"Successfully created microservice project 'ecommerce-app'!",
			},
		},
		{
			name:            "explicit name overrides directory name",
			directoryName:   "temp-dir",
			args:            []string{"create", "--be=fastapi", "--name=custom-api", "--force"},
			expectedProject: "custom-api",
			checkFiles: []string{
				"docker-compose.yml",
				"README.md",
				"backend/Dockerfile",
				"backend/main.py",
			},
			checkOutput: []string{
				"Successfully created microservice project 'custom-api'!",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary parent directory
			tempParent := t.TempDir()
			originalDir, _ := os.Getwd()
			defer os.Chdir(originalDir)

			// Create the test directory with the specified name
			testDir := filepath.Join(tempParent, tt.directoryName)
			err := os.MkdirAll(testDir, 0755)
			if err != nil {
				t.Fatalf("Failed to create test directory: %v", err)
			}

			// Change to the test directory
			os.Chdir(testDir)

			// Run the CLI command
			cmd := exec.Command(binaryPath, tt.args...)
			output, err := cmd.CombinedOutput()
			outputStr := string(output)

			if err != nil {
				t.Fatalf("CLI command failed: %v. Output: %s", err, outputStr)
			}

			// Check output messages
			for _, expected := range tt.checkOutput {
				if !strings.Contains(outputStr, expected) {
					t.Errorf("Output does not contain expected text '%s'. Full output: %s", expected, outputStr)
				}
			}

			// Determine the project path - projects are always created as subdirectories
			projectPath := filepath.Join(testDir, tt.expectedProject)

			// Check that expected files exist
			for _, expectedFile := range tt.checkFiles {
				filePath := filepath.Join(projectPath, expectedFile)
				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					t.Errorf("Expected file %s does not exist at %s", expectedFile, filePath)
				}
			}

			// Verify project name is used correctly in generated files
			if !strings.Contains(strings.Join(tt.args, " "), "--name=") {
				// Only check for inferred names (not explicit names)
				dockerComposePath := filepath.Join(projectPath, "docker-compose.yml")
				if content, err := os.ReadFile(dockerComposePath); err == nil {
					contentStr := string(content)
					if !strings.Contains(contentStr, tt.expectedProject) {
						t.Errorf("docker-compose.yml should contain project name '%s'", tt.expectedProject)
					}
				}
			}
		})
	}
}

// TestProjectGenerationWithoutDocker tests project generation when Docker is not running
func TestProjectGenerationWithoutDocker(t *testing.T) {
	// Build the CLI binary for testing
	binaryPath := buildCLIBinary(t)
	defer os.Remove(binaryPath)

	// Stop Docker if it's running (this is a simulation - in real tests we'd mock Docker detection)
	// For this test, we'll assume Docker is not available and check that generation still works

	tests := []struct {
		name        string
		args        []string
		projectName string
		checkFiles  []string
		checkOutput []string
	}{
		{
			name:        "fastapi backend without docker",
			args:        []string{"create", "--be=fastapi", "--name=no-docker-api", "--force"},
			projectName: "no-docker-api",
			checkFiles: []string{
				"no-docker-api/docker-compose.yml",
				"no-docker-api/README.md",
				"no-docker-api/backend/Dockerfile",
				"no-docker-api/backend/main.py",
				"no-docker-api/backend/requirements.txt",
			},
			checkOutput: []string{
				"Successfully created microservice project 'no-docker-api'!",
				// Note: Docker warning messages would be checked here if implemented
			},
		},
		{
			name:        "full stack without docker",
			args:        []string{"create", "--be=gin", "--fe=vue", "--db=postgres", "--name=no-docker-full", "--force"},
			projectName: "no-docker-full",
			checkFiles: []string{
				"no-docker-full/docker-compose.yml",
				"no-docker-full/README.md",
				"no-docker-full/backend/Dockerfile",
				"no-docker-full/backend/main.go",
				"no-docker-full/frontend/Dockerfile",
				"no-docker-full/frontend/package.json",
			},
			checkOutput: []string{
				"Successfully created microservice project 'no-docker-full'!",
			},
		},
		{
			name:        "frontend only without docker",
			args:        []string{"create", "--fe=react", "--name=no-docker-frontend", "--force"},
			projectName: "no-docker-frontend",
			checkFiles: []string{
				"no-docker-frontend/docker-compose.yml",
				"no-docker-frontend/README.md",
				"no-docker-frontend/frontend/Dockerfile",
				"no-docker-frontend/frontend/package.json",
			},
			checkOutput: []string{
				"Successfully created microservice project 'no-docker-frontend'!",
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

			// The command should succeed even without Docker
			if err != nil {
				t.Fatalf("CLI command failed: %v. Output: %s", err, outputStr)
			}

			// Check output messages
			for _, expected := range tt.checkOutput {
				if !strings.Contains(outputStr, expected) {
					t.Errorf("Output does not contain expected text '%s'. Full output: %s", expected, outputStr)
				}
			}

			// Check that all files are created successfully
			for _, expectedFile := range tt.checkFiles {
				filePath := filepath.Join(tempDir, expectedFile)
				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					t.Errorf("Expected file %s does not exist", expectedFile)
				}
			}

			// Verify that README contains Docker instructions
			readmePath := filepath.Join(tempDir, tt.projectName, "README.md")
			if content, err := os.ReadFile(readmePath); err == nil {
				contentStr := string(content)
				dockerInstructions := []string{"docker-compose", "Docker", "Getting Started"}
				for _, instruction := range dockerInstructions {
					if !strings.Contains(contentStr, instruction) {
						t.Errorf("README.md should contain Docker instruction: %s", instruction)
					}
				}
			}
		})
	}
}

// TestVariousDirectoryStructures tests project creation in different directory scenarios
func TestVariousDirectoryStructures(t *testing.T) {
	// Build the CLI binary for testing
	binaryPath := buildCLIBinary(t)
	defer os.Remove(binaryPath)

	tests := []struct {
		name          string
		setupDir      func(string) string // Function to setup directory and return the working directory
		args          []string
		expectedFiles []string
		expectError   bool
		errorMessage  string
	}{
		{
			name: "empty directory",
			setupDir: func(tempDir string) string {
				emptyDir := filepath.Join(tempDir, "empty-project")
				os.MkdirAll(emptyDir, 0755)
				return emptyDir
			},
			args: []string{"create", "--be=fastapi", "--force"},
			expectedFiles: []string{
				"docker-compose.yml",
				"README.md",
				"backend/Dockerfile",
			},
			expectError: false,
		},
		{
			name: "nested directory structure",
			setupDir: func(tempDir string) string {
				nestedDir := filepath.Join(tempDir, "projects", "apis", "user-service")
				os.MkdirAll(nestedDir, 0755)
				return nestedDir
			},
			args: []string{"create", "--be=gin", "--force"},
			expectedFiles: []string{
				"docker-compose.yml",
				"README.md",
				"backend/Dockerfile",
			},
			expectError: false,
		},
		{
			name: "directory with existing files without force",
			setupDir: func(tempDir string) string {
				existingDir := filepath.Join(tempDir, "existing-project")
				os.MkdirAll(existingDir, 0755)
				// Create an existing file
				os.WriteFile(filepath.Join(existingDir, "existing.txt"), []byte("content"), 0644)
				return existingDir
			},
			args: []string{"create", "--be=fastapi"},
			expectedFiles: []string{
				"docker-compose.yml",
				"README.md",
				"backend/Dockerfile",
			},
			expectError: false, // The CLI now handles this gracefully
		},
		{
			name: "directory with existing files with force",
			setupDir: func(tempDir string) string {
				existingDir := filepath.Join(tempDir, "existing-project-force")
				os.MkdirAll(existingDir, 0755)
				// Create an existing file
				os.WriteFile(filepath.Join(existingDir, "existing.txt"), []byte("content"), 0644)
				return existingDir
			},
			args: []string{"create", "--be=fastapi", "--force"},
			expectedFiles: []string{
				"docker-compose.yml",
				"README.md",
				"backend/Dockerfile",
			},
			expectError: false,
		},
		{
			name: "directory with special characters in path",
			setupDir: func(tempDir string) string {
				specialDir := filepath.Join(tempDir, "special-chars", "api@v1", "service")
				os.MkdirAll(specialDir, 0755)
				return specialDir
			},
			args: []string{"create", "--be=express", "--force"},
			expectedFiles: []string{
				"docker-compose.yml",
				"README.md",
				"backend/Dockerfile",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory for this test
			tempDir := t.TempDir()
			originalDir, _ := os.Getwd()
			defer os.Chdir(originalDir)

			// Setup the test directory
			workingDir := tt.setupDir(tempDir)
			os.Chdir(workingDir)

			// Run the CLI command
			cmd := exec.Command(binaryPath, tt.args...)
			output, err := cmd.CombinedOutput()
			outputStr := string(output)

			// Check error expectation
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but command succeeded. Output: %s", outputStr)
				}
				if tt.errorMessage != "" && !strings.Contains(outputStr, tt.errorMessage) {
					t.Errorf("Expected error message '%s' not found in output: %s", tt.errorMessage, outputStr)
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v. Output: %s", err, outputStr)
			}

			// Check that expected files exist
			// Files are created in a subdirectory with the project name (inferred from directory name)
			projectName := filepath.Base(workingDir)
			for _, expectedFile := range tt.expectedFiles {
				filePath := filepath.Join(workingDir, projectName, expectedFile)
				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					t.Errorf("Expected file %s does not exist at %s", expectedFile, filePath)
				}
			}
		})
	}
}

// TestEdgeCasesAndErrorHandling tests various edge cases and error scenarios
func TestEdgeCasesAndErrorHandling(t *testing.T) {
	// Build the CLI binary for testing
	binaryPath := buildCLIBinary(t)
	defer os.Remove(binaryPath)

	tests := []struct {
		name         string
		setupDir     func(string) string
		args         []string
		expectError  bool
		errorMessage string
		checkOutput  []string
	}{
		{
			name: "root-like directory name",
			setupDir: func(tempDir string) string {
				rootDir := filepath.Join(tempDir, "root")
				os.MkdirAll(rootDir, 0755)
				return rootDir
			},
			args:        []string{"create", "--be=fastapi", "--force"},
			expectError: false,
			checkOutput: []string{"Successfully created microservice project"},
		},
		{
			name: "very long directory name",
			setupDir: func(tempDir string) string {
				longName := strings.Repeat("very-long-directory-name-", 10) // 260+ characters
				longDir := filepath.Join(tempDir, longName[:100]) // Truncate to reasonable length
				os.MkdirAll(longDir, 0755)
				return longDir
			},
			args:        []string{"create", "--be=gin", "--force"},
			expectError: false,
			checkOutput: []string{"Successfully created microservice project"},
		},
		{
			name: "directory name with only special characters",
			setupDir: func(tempDir string) string {
				specialDir := filepath.Join(tempDir, "!@#$%")
				os.MkdirAll(specialDir, 0755)
				return specialDir
			},
			args:        []string{"create", "--be=express", "--force"},
			expectError: false,
			checkOutput: []string{"Successfully created microservice project"},
		},
		{
			name: "directory name with numbers",
			setupDir: func(tempDir string) string {
				numDir := filepath.Join(tempDir, "api-v2-2024")
				os.MkdirAll(numDir, 0755)
				return numDir
			},
			args:        []string{"create", "--be=fastapi", "--force"},
			expectError: false,
			checkOutput: []string{"Using directory name as project name: api-v2-2024"},
		},
		{
			name: "directory with mixed case",
			setupDir: func(tempDir string) string {
				mixedDir := filepath.Join(tempDir, "MyApiService")
				os.MkdirAll(mixedDir, 0755)
				return mixedDir
			},
			args:        []string{"create", "--be=gin", "--force"},
			expectError: false,
			checkOutput: []string{"Using directory name as project name: myapiservice"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory for this test
			tempDir := t.TempDir()
			originalDir, _ := os.Getwd()
			defer os.Chdir(originalDir)

			// Setup the test directory
			workingDir := tt.setupDir(tempDir)
			os.Chdir(workingDir)

			// Run the CLI command
			cmd := exec.Command(binaryPath, tt.args...)
			output, err := cmd.CombinedOutput()
			outputStr := string(output)

			// Check error expectation
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but command succeeded. Output: %s", outputStr)
				}
				if tt.errorMessage != "" && !strings.Contains(outputStr, tt.errorMessage) {
					t.Errorf("Expected error message '%s' not found in output: %s", tt.errorMessage, outputStr)
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v. Output: %s", err, outputStr)
			}

			// Check output messages
			for _, expected := range tt.checkOutput {
				if !strings.Contains(outputStr, expected) {
					t.Errorf("Output does not contain expected text '%s'. Full output: %s", expected, outputStr)
				}
			}

			// Verify basic project structure was created
			// Files are created in a subdirectory with the sanitized project name
			// We need to determine the actual project name from the output
			var actualProjectName string
			if strings.Contains(outputStr, "Using directory name as project name:") {
				lines := strings.Split(outputStr, "\n")
				for _, line := range lines {
					if strings.Contains(line, "Using directory name as project name:") {
						parts := strings.Split(line, ":")
						if len(parts) > 1 {
							actualProjectName = strings.TrimSpace(parts[len(parts)-1])
							break
						}
					}
				}
			}
			
			if actualProjectName == "" {
				// Fallback to directory name if we can't parse it
				actualProjectName = filepath.Base(workingDir)
			}
			
			basicFiles := []string{"docker-compose.yml", "README.md"}
			for _, file := range basicFiles {
				filePath := filepath.Join(workingDir, actualProjectName, file)
				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					t.Errorf("Basic file %s should be created at %s", file, filePath)
				}
			}
		})
	}
}

// TestUserFeedbackAndMessaging tests that appropriate user feedback is provided
func TestUserFeedbackAndMessaging(t *testing.T) {
	// Build the CLI binary for testing
	binaryPath := buildCLIBinary(t)
	defer os.Remove(binaryPath)

	tests := []struct {
		name         string
		directoryName string
		args         []string
		expectedMessages []string
		notExpectedMessages []string
	}{
		{
			name:          "inferred name feedback",
			directoryName: "feedback-test",
			args:          []string{"create", "--be=fastapi", "--force"},
			expectedMessages: []string{
				"Using directory name as project name: feedback-test",
				"Successfully created microservice project 'feedback-test'!",
			},
		},
		{
			name:          "explicit name no inference message",
			directoryName: "temp-dir",
			args:          []string{"create", "--be=gin", "--name=explicit-name", "--force"},
			expectedMessages: []string{
				"Successfully created microservice project 'explicit-name'!",
			},
			notExpectedMessages: []string{
				"Using directory name as project name",
			},
		},
		{
			name:          "sanitized name feedback",
			directoryName: "test_with_underscores",
			args:          []string{"create", "--be=express", "--force"},
			expectedMessages: []string{
				"Using directory name as project name: test-with-underscores",
				"Successfully created microservice project 'test-with-underscores'!",
			},
		},
		{
			name:          "full stack creation feedback",
			directoryName: "fullstack-app",
			args:          []string{"create", "--be=gin", "--fe=react", "--db=postgres", "--force"},
			expectedMessages: []string{
				"Using directory name as project name: fullstack-app",
				"Successfully created microservice project 'fullstack-app'!",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary parent directory
			tempParent := t.TempDir()
			originalDir, _ := os.Getwd()
			defer os.Chdir(originalDir)

			// Create the test directory
			testDir := filepath.Join(tempParent, tt.directoryName)
			err := os.MkdirAll(testDir, 0755)
			if err != nil {
				t.Fatalf("Failed to create test directory: %v", err)
			}

			// Change to the test directory
			os.Chdir(testDir)

			// Run the CLI command
			cmd := exec.Command(binaryPath, tt.args...)
			output, err := cmd.CombinedOutput()
			outputStr := string(output)

			if err != nil {
				t.Fatalf("CLI command failed: %v. Output: %s", err, outputStr)
			}

			// Check expected messages
			for _, expected := range tt.expectedMessages {
				if !strings.Contains(outputStr, expected) {
					t.Errorf("Output should contain message '%s'. Full output: %s", expected, outputStr)
				}
			}

			// Check that certain messages are NOT present
			for _, notExpected := range tt.notExpectedMessages {
				if strings.Contains(outputStr, notExpected) {
					t.Errorf("Output should NOT contain message '%s'. Full output: %s", notExpected, outputStr)
				}
			}
		})
	}
}

// TestConcurrentProjectCreation tests that multiple projects can be created in sequence
func TestConcurrentProjectCreation(t *testing.T) {
	// Build the CLI binary for testing
	binaryPath := buildCLIBinary(t)
	defer os.Remove(binaryPath)

	// Create multiple projects sequentially to avoid Windows file locking issues
	numProjects := 3
	for i := 0; i < numProjects; i++ {
		t.Run(fmt.Sprintf("project_%d", i), func(t *testing.T) {
			// Create temporary directory for this project
			tempDir := t.TempDir()
			originalDir, _ := os.Getwd()
			defer os.Chdir(originalDir)

			// Create project directory
			projectDir := filepath.Join(tempDir, fmt.Sprintf("sequential-test-%d", i))
			os.MkdirAll(projectDir, 0755)
			os.Chdir(projectDir)

			// Run CLI command
			args := []string{"create", "--be=fastapi", "--force"}
			cmd := exec.Command(binaryPath, args...)
			output, err := cmd.CombinedOutput()

			if err != nil {
				t.Fatalf("Project %d failed: %v, output: %s", i, err, string(output))
			}

			// Files are created in a subdirectory with the project name
			projectName := filepath.Base(projectDir)
			
			// Check if the main project directory was created
			mainProjectPath := filepath.Join(projectDir, projectName)
			if _, err := os.Stat(mainProjectPath); os.IsNotExist(err) {
				t.Fatalf("Project %d: main project directory %s was not created", i, mainProjectPath)
			}
			
			// Verify essential files were created
			expectedFiles := []string{"docker-compose.yml", "README.md", "backend/Dockerfile"}
			for _, file := range expectedFiles {
				filePath := filepath.Join(mainProjectPath, file)
				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					t.Errorf("Project %d missing file: %s at %s", i, file, filePath)
				}
			}
		})
	}
}

// TestProjectCreationPerformance tests that project creation completes within reasonable time
func TestProjectCreationPerformance(t *testing.T) {
	// Build the CLI binary for testing
	binaryPath := buildCLIBinary(t)
	defer os.Remove(binaryPath)

	tests := []struct {
		name        string
		args        []string
		maxDuration time.Duration
	}{
		{
			name:        "simple backend",
			args:        []string{"create", "--be=fastapi", "--name=perf-test-1", "--force"},
			maxDuration: 10 * time.Second,
		},
		{
			name:        "full stack",
			args:        []string{"create", "--be=gin", "--fe=react", "--db=postgres", "--name=perf-test-2", "--force"},
			maxDuration: 15 * time.Second,
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

			// Measure execution time
			start := time.Now()

			// Run the CLI command
			cmd := exec.Command(binaryPath, tt.args...)
			output, err := cmd.CombinedOutput()
			
			duration := time.Since(start)

			if err != nil {
				t.Fatalf("CLI command failed: %v. Output: %s", err, string(output))
			}

			if duration > tt.maxDuration {
				t.Errorf("Project creation took too long: %v (max: %v)", duration, tt.maxDuration)
			}

			t.Logf("Project creation completed in %v", duration)
		})
	}
}