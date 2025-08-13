package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"microservice-bootstrapper/internal/generator"
	"microservice-bootstrapper/internal/interfaces"
)

func TestEndToEndProjectGeneration(t *testing.T) {
	tests := []struct {
		name           string
		config         interfaces.CLIConfig
		expectedFiles  []string
		expectedDirs   []string
		checkContent   map[string][]string // file -> content to check
	}{
		{
			name: "full stack fastapi react postgres",
			config: interfaces.CLIConfig{
				Backend:     "fastapi",
				Frontend:    "react",
				Database:    "postgres",
				ProjectName: "fullstack-app",
				Force:       true,
			},
			expectedFiles: []string{
				"docker-compose.yml",
				"README.md",
				".env.example",
				".gitignore",
				"backend/Dockerfile",
				"backend/main.py",
				"backend/requirements.txt",
				"frontend/Dockerfile",
				"frontend/package.json",
				"frontend/App.js",
			},
			expectedDirs: []string{
				"backend",
				"frontend",
				"backend/app",
				"backend/app/api",
				"backend/app/core",
				"backend/app/models",
				"frontend/src",
				"frontend/src/components",
				"frontend/src/pages",
				"frontend/public",
			},
			checkContent: map[string][]string{
				"docker-compose.yml": {"version: '3.8'", "backend:", "frontend:", "postgres:", "fullstack-app"},
				"README.md":          {"fullstack-app", "FastAPI", "React", "Getting Started"},
				"backend/main.py":    {"FastAPI", "app = FastAPI", "@app.get"},
				"frontend/package.json": {"fullstack-app-frontend", "react", "scripts"},
			},
		},
		{
			name: "backend only gin mongo",
			config: interfaces.CLIConfig{
				Backend:     "gin",
				Database:    "mongo",
				ProjectName: "api-service",
				Force:       true,
			},
			expectedFiles: []string{
				"docker-compose.yml",
				"README.md",
				".env.example",
				".gitignore",
				"backend/Dockerfile",
				"backend/main.go",
				"backend/go.mod",
			},
			expectedDirs: []string{
				"backend",
				"backend/cmd",
				"backend/internal",
				"backend/pkg",
				"backend/api",
			},
			checkContent: map[string][]string{
				"docker-compose.yml": {"version: '3.8'", "backend:", "mongo:", "api-service"},
				"README.md":          {"api-service", "Gin", "Getting Started"},
				"backend/main.go":    {"gin", "gin.Default", "GET"},
				"backend/go.mod":     {"api-service-backend", "go 1.21"},
			},
		},
		{
			name: "frontend only vue",
			config: interfaces.CLIConfig{
				Frontend:    "vue",
				ProjectName: "vue-app",
				Force:       true,
			},
			expectedFiles: []string{
				"docker-compose.yml",
				"README.md",
				".env.example",
				".gitignore",
				"frontend/Dockerfile",
				"frontend/package.json",
				"frontend/App.vue",
			},
			expectedDirs: []string{
				"frontend",
				"frontend/src",
				"frontend/src/components",
				"frontend/src/pages",
				"frontend/public",
			},
			checkContent: map[string][]string{
				"docker-compose.yml": {"version: '3.8'", "frontend:", "vue-app"},
				"README.md":          {"vue-app", "Vue.js", "Getting Started"},
				"frontend/package.json": {"vue-app-frontend", "vue", "scripts"},
			},
		},
		{
			name: "express angular mysql",
			config: interfaces.CLIConfig{
				Backend:     "express",
				Frontend:    "angular",
				Database:    "mysql",
				ProjectName: "enterprise-app",
				Force:       true,
			},
			expectedFiles: []string{
				"docker-compose.yml",
				"README.md",
				".env.example",
				".gitignore",
				"backend/Dockerfile",
				"backend/server.js",
				"backend/package.json",
				"frontend/Dockerfile",
				"frontend/package.json",
				"frontend/angular.json",
			},
			expectedDirs: []string{
				"backend",
				"frontend",
				"backend/src",
				"backend/src/routes",
				"backend/src/middleware",
				"backend/src/models",
				"frontend/src",
				"frontend/src/app",
				"frontend/src/assets",
				"frontend/src/environments",
			},
			checkContent: map[string][]string{
				"docker-compose.yml": {"version: '3.8'", "backend:", "frontend:", "mysql:", "enterprise-app"},
				"README.md":          {"enterprise-app", "Express.js", "Angular", "Getting Started"},
				"backend/server.js":  {"express", "app = express", "app.get"},
				"frontend/package.json": {"enterprise-app-frontend", "@angular/core", "scripts"},
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

			// Create generator and generate project
			gen := generator.NewGenerator()
			err := gen.Generate(tt.config)
			if err != nil {
				t.Fatalf("Failed to generate project: %v", err)
			}

			projectPath := filepath.Join(tempDir, tt.config.ProjectName)

			// Check that expected files exist
			for _, expectedFile := range tt.expectedFiles {
				filePath := filepath.Join(projectPath, expectedFile)
				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					t.Errorf("Expected file %s does not exist", expectedFile)
				}
			}

			// Check that expected directories exist
			for _, expectedDir := range tt.expectedDirs {
				dirPath := filepath.Join(projectPath, expectedDir)
				if stat, err := os.Stat(dirPath); os.IsNotExist(err) || !stat.IsDir() {
					t.Errorf("Expected directory %s does not exist or is not a directory", expectedDir)
				}
			}

			// Check file content
			for file, expectedContent := range tt.checkContent {
				filePath := filepath.Join(projectPath, file)
				content, err := os.ReadFile(filePath)
				if err != nil {
					t.Errorf("Failed to read file %s: %v", file, err)
					continue
				}

				contentStr := string(content)
				for _, expected := range expectedContent {
					if !strings.Contains(contentStr, expected) {
						t.Errorf("File %s does not contain expected content: %s", file, expected)
					}
				}
			}
		})
	}
}

func TestDockerComposeValidation(t *testing.T) {
	tests := []struct {
		name   string
		config interfaces.CLIConfig
	}{
		{
			name: "fastapi react postgres",
			config: interfaces.CLIConfig{
				Backend:     "fastapi",
				Frontend:    "react",
				Database:    "postgres",
				ProjectName: "test-app",
				Force:       true,
			},
		},
		{
			name: "gin vue mysql",
			config: interfaces.CLIConfig{
				Backend:     "gin",
				Frontend:    "vue",
				Database:    "mysql",
				ProjectName: "test-app",
				Force:       true,
			},
		},
		{
			name: "express angular mongo",
			config: interfaces.CLIConfig{
				Backend:     "express",
				Frontend:    "angular",
				Database:    "mongo",
				ProjectName: "test-app",
				Force:       true,
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

			// Create generator and generate project
			gen := generator.NewGenerator()
			err := gen.Generate(tt.config)
			if err != nil {
				t.Fatalf("Failed to generate project: %v", err)
			}

			projectPath := filepath.Join(tempDir, tt.config.ProjectName)
			composePath := filepath.Join(projectPath, "docker-compose.yml")

			// Read and validate docker-compose.yml
			content, err := os.ReadFile(composePath)
			if err != nil {
				t.Fatalf("Failed to read docker-compose.yml: %v", err)
			}

			composeStr := string(content)

			// Basic validation checks
			requiredSections := []string{
				"version: '3.8'",
				"services:",
				"networks:",
			}

			for _, section := range requiredSections {
				if !strings.Contains(composeStr, section) {
					t.Errorf("docker-compose.yml missing required section: %s", section)
				}
			}

			// Service-specific checks
			if tt.config.Backend != "" {
				if !strings.Contains(composeStr, "backend:") {
					t.Error("docker-compose.yml should contain backend service")
				}
			}

			if tt.config.Frontend != "" {
				if !strings.Contains(composeStr, "frontend:") {
					t.Error("docker-compose.yml should contain frontend service")
				}
			}

			if tt.config.Database != "" {
				dbServiceName := tt.config.Database
				if tt.config.Database == "mongo" {
					dbServiceName = "mongo"
				} else if tt.config.Database == "postgres" {
					dbServiceName = "postgres"
				}
				
				if !strings.Contains(composeStr, dbServiceName+":") {
					t.Errorf("docker-compose.yml should contain %s service", dbServiceName)
				}
			}

			// Network configuration check
			networkName := tt.config.ProjectName + "-network"
			if !strings.Contains(composeStr, networkName) {
				t.Errorf("docker-compose.yml should contain network: %s", networkName)
			}
		})
	}
}

func TestProjectStructureConsistency(t *testing.T) {
	// Test that generated projects have consistent structure across different technology combinations
	configs := []interfaces.CLIConfig{
		{Backend: "fastapi", ProjectName: "fastapi-test", Force: true},
		{Backend: "express", ProjectName: "express-test", Force: true},
		{Backend: "gin", ProjectName: "gin-test", Force: true},
		{Frontend: "react", ProjectName: "react-test", Force: true},
		{Frontend: "vue", ProjectName: "vue-test", Force: true},
		{Frontend: "angular", ProjectName: "angular-test", Force: true},
		{Database: "postgres", ProjectName: "postgres-test", Force: true},
		{Database: "mongo", ProjectName: "mongo-test", Force: true},
		{Database: "mysql", ProjectName: "mysql-test", Force: true},
		{Database: "redis", ProjectName: "redis-test", Force: true},
	}

	for _, config := range configs {
		t.Run(config.ProjectName, func(t *testing.T) {
			// Create temporary directory for this test
			tempDir := t.TempDir()
			originalDir, _ := os.Getwd()
			defer os.Chdir(originalDir)

			// Change to temp directory
			os.Chdir(tempDir)

			// Create generator and generate project
			gen := generator.NewGenerator()
			err := gen.Generate(config)
			if err != nil {
				t.Fatalf("Failed to generate project: %v", err)
			}

			projectPath := filepath.Join(tempDir, config.ProjectName)

			// All projects should have these base files
			baseFiles := []string{
				"docker-compose.yml",
				"README.md",
				".env.example",
				".gitignore",
			}

			for _, file := range baseFiles {
				filePath := filepath.Join(projectPath, file)
				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					t.Errorf("Base file %s should exist in all projects", file)
				}
			}

			// Check service-specific structure
			if config.Backend != "" {
				backendDir := filepath.Join(projectPath, "backend")
				if stat, err := os.Stat(backendDir); os.IsNotExist(err) || !stat.IsDir() {
					t.Error("Backend projects should have backend directory")
				}

				dockerfilePath := filepath.Join(backendDir, "Dockerfile")
				if _, err := os.Stat(dockerfilePath); os.IsNotExist(err) {
					t.Error("Backend projects should have Dockerfile")
				}
			}

			if config.Frontend != "" {
				frontendDir := filepath.Join(projectPath, "frontend")
				if stat, err := os.Stat(frontendDir); os.IsNotExist(err) || !stat.IsDir() {
					t.Error("Frontend projects should have frontend directory")
				}

				dockerfilePath := filepath.Join(frontendDir, "Dockerfile")
				if _, err := os.Stat(dockerfilePath); os.IsNotExist(err) {
					t.Error("Frontend projects should have Dockerfile")
				}
			}
		})
	}
}

func TestGeneratedProjectsCanBuild(t *testing.T) {
	// This test verifies that generated projects have valid configuration files
	// Note: We don't actually build Docker images in tests, but we validate the configuration
	
	configs := []interfaces.CLIConfig{
		{
			Backend:     "fastapi",
			Frontend:    "react",
			Database:    "postgres",
			ProjectName: "build-test-1",
			Force:       true,
		},
		{
			Backend:     "gin",
			Database:    "redis",
			ProjectName: "build-test-2",
			Force:       true,
		},
		{
			Frontend:    "vue",
			ProjectName: "build-test-3",
			Force:       true,
		},
	}

	for _, config := range configs {
		t.Run(config.ProjectName, func(t *testing.T) {
			// Create temporary directory for this test
			tempDir := t.TempDir()
			originalDir, _ := os.Getwd()
			defer os.Chdir(originalDir)

			// Change to temp directory
			os.Chdir(tempDir)

			// Create generator and generate project
			gen := generator.NewGenerator()
			err := gen.Generate(config)
			if err != nil {
				t.Fatalf("Failed to generate project: %v", err)
			}

			projectPath := filepath.Join(tempDir, config.ProjectName)

			// Validate Dockerfiles have proper structure
			if config.Backend != "" {
				dockerfilePath := filepath.Join(projectPath, "backend", "Dockerfile")
				content, err := os.ReadFile(dockerfilePath)
				if err != nil {
					t.Fatalf("Failed to read backend Dockerfile: %v", err)
				}

				dockerfileStr := string(content)
				requiredInstructions := []string{"FROM", "WORKDIR", "COPY", "RUN", "EXPOSE", "CMD"}
				for _, instruction := range requiredInstructions {
					if !strings.Contains(dockerfileStr, instruction) {
						t.Errorf("Backend Dockerfile missing instruction: %s", instruction)
					}
				}
			}

			if config.Frontend != "" {
				dockerfilePath := filepath.Join(projectPath, "frontend", "Dockerfile")
				content, err := os.ReadFile(dockerfilePath)
				if err != nil {
					t.Fatalf("Failed to read frontend Dockerfile: %v", err)
				}

				dockerfileStr := string(content)
				requiredInstructions := []string{"FROM", "WORKDIR", "COPY", "RUN", "EXPOSE", "CMD"}
				for _, instruction := range requiredInstructions {
					if !strings.Contains(dockerfileStr, instruction) {
						t.Errorf("Frontend Dockerfile missing instruction: %s", instruction)
					}
				}
			}

			// Validate package.json files (for Node.js projects)
			if config.Backend == "express" {
				packagePath := filepath.Join(projectPath, "backend", "package.json")
				content, err := os.ReadFile(packagePath)
				if err != nil {
					t.Fatalf("Failed to read backend package.json: %v", err)
				}

				packageStr := string(content)
				requiredFields := []string{"name", "version", "scripts", "dependencies"}
				for _, field := range requiredFields {
					if !strings.Contains(packageStr, `"`+field+`"`) {
						t.Errorf("Backend package.json missing field: %s", field)
					}
				}
			}

			if config.Frontend != "" {
				packagePath := filepath.Join(projectPath, "frontend", "package.json")
				content, err := os.ReadFile(packagePath)
				if err != nil {
					t.Fatalf("Failed to read frontend package.json: %v", err)
				}

				packageStr := string(content)
				requiredFields := []string{"name", "version", "scripts", "dependencies"}
				for _, field := range requiredFields {
					if !strings.Contains(packageStr, `"`+field+`"`) {
						t.Errorf("Frontend package.json missing field: %s", field)
					}
				}
			}

			// Validate go.mod files (for Go projects)
			if config.Backend == "gin" {
				goModPath := filepath.Join(projectPath, "backend", "go.mod")
				content, err := os.ReadFile(goModPath)
				if err != nil {
					t.Fatalf("Failed to read backend go.mod: %v", err)
				}

				goModStr := string(content)
				if !strings.Contains(goModStr, "module") {
					t.Error("Backend go.mod should contain module declaration")
				}
				if !strings.Contains(goModStr, "go 1.21") {
					t.Error("Backend go.mod should specify Go version")
				}
			}

			// Validate requirements.txt files (for Python projects)
			if config.Backend == "fastapi" {
				reqPath := filepath.Join(projectPath, "backend", "requirements.txt")
				content, err := os.ReadFile(reqPath)
				if err != nil {
					t.Fatalf("Failed to read backend requirements.txt: %v", err)
				}

				reqStr := string(content)
				if !strings.Contains(reqStr, "fastapi") {
					t.Error("Backend requirements.txt should contain fastapi")
				}
				if !strings.Contains(reqStr, "uvicorn") {
					t.Error("Backend requirements.txt should contain uvicorn")
				}
			}
		})
	}
}

func TestErrorHandlingInGeneration(t *testing.T) {
	// Test error handling scenarios
	tests := []struct {
		name        string
		config      interfaces.CLIConfig
		expectError bool
		errorMsg    string
	}{
		{
			name: "invalid backend technology",
			config: interfaces.CLIConfig{
				Backend:     "invalid",
				ProjectName: "error-test",
				Force:       true,
			},
			expectError: true,
			errorMsg:    "unsupported backend",
		},
		{
			name: "invalid frontend technology",
			config: interfaces.CLIConfig{
				Frontend:    "invalid",
				ProjectName: "error-test",
				Force:       true,
			},
			expectError: true,
			errorMsg:    "unsupported frontend",
		},
		{
			name: "invalid database technology",
			config: interfaces.CLIConfig{
				Database:    "invalid",
				ProjectName: "error-test",
				Force:       true,
			},
			expectError: true,
			errorMsg:    "unsupported database",
		},
		{
			name: "no services specified",
			config: interfaces.CLIConfig{
				ProjectName: "error-test",
				Force:       true,
			},
			expectError: true,
			errorMsg:    "at least one service must be specified",
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

			// Create generator and attempt to generate project
			gen := generator.NewGenerator()
			err := gen.Generate(tt.config)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error containing '%s', got: %v", tt.errorMsg, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}