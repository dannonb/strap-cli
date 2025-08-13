package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"microservice-bootstrapper/internal/interfaces"
)

func TestGenerateProjectDocumentation(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	
	// Change to temp directory
	os.Chdir(tempDir)

	generator := NewGenerator()

	testConfig := interfaces.CLIConfig{
		ProjectName: "test-microservice",
		Backend:     "fastapi",
		Frontend:    "react",
		Database:    "postgres",
		Force:       true,
	}

	// Generate the project
	err := generator.Generate(testConfig)
	if err != nil {
		t.Fatalf("Failed to generate project: %v", err)
	}

	// Verify that project documentation files were created
	projectPath := filepath.Join(tempDir, "test-microservice")
	
	expectedFiles := []string{
		"README.md",
		".env.example",
		".gitignore",
		"docker-compose.yml",
	}

	for _, expectedFile := range expectedFiles {
		filePath := filepath.Join(projectPath, expectedFile)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("Expected file %s was not created", expectedFile)
		}
	}

	// Verify README.md contains expected content
	readmePath := filepath.Join(projectPath, "README.md")
	readmeContent, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("Failed to read README.md: %v", err)
	}

	readmeStr := string(readmeContent)
	expectedContent := []string{
		"test-microservice",
		"fastapi",
		"react", 
		"postgres",
		"Getting Started",
		"Prerequisites",
		"Docker",
		"Docker Compose",
	}

	for _, expected := range expectedContent {
		if !strings.Contains(readmeStr, expected) {
			t.Errorf("README.md does not contain expected content: %s", expected)
		}
	}

	// Verify .env.example contains expected content
	envPath := filepath.Join(projectPath, ".env.example")
	envContent, err := os.ReadFile(envPath)
	if err != nil {
		t.Fatalf("Failed to read .env.example: %v", err)
	}

	envStr := string(envContent)
	expectedEnvContent := []string{
		"test-microservice",
		"BACKEND_PORT=8000",
		"FRONTEND_PORT=3000",
		"DATABASE_PORT=5432",
		"POSTGRES_DB=microservice_db",
	}

	for _, expected := range expectedEnvContent {
		if !strings.Contains(envStr, expected) {
			t.Errorf(".env.example does not contain expected content: %s", expected)
		}
	}

	// Verify .gitignore contains expected content
	gitignorePath := filepath.Join(projectPath, ".gitignore")
	gitignoreContent, err := os.ReadFile(gitignorePath)
	if err != nil {
		t.Fatalf("Failed to read .gitignore: %v", err)
	}

	gitignoreStr := string(gitignoreContent)
	expectedGitignoreContent := []string{
		"test-microservice",
		".env",
		"__pycache__", // Python-specific for FastAPI
		"node_modules", // Node.js-specific for React
		"*.log",
		"docker-compose.override.yml",
	}

	for _, expected := range expectedGitignoreContent {
		if !strings.Contains(gitignoreStr, expected) {
			t.Errorf(".gitignore does not contain expected content: %s", expected)
		}
	}
}

func TestGenerateProjectDocumentationBackendOnly(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	
	// Change to temp directory
	os.Chdir(tempDir)

	generator := NewGenerator()

	testConfig := interfaces.CLIConfig{
		ProjectName: "api-service",
		Backend:     "gin",
		Database:    "mongo",
		Force:       true,
	}

	// Generate the project
	err := generator.Generate(testConfig)
	if err != nil {
		t.Fatalf("Failed to generate project: %v", err)
	}

	// Verify that project documentation files were created
	projectPath := filepath.Join(tempDir, "api-service")
	
	// Verify README.md contains expected content for backend-only project
	readmePath := filepath.Join(projectPath, "README.md")
	readmeContent, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("Failed to read README.md: %v", err)
	}

	readmeStr := string(readmeContent)
	
	// Should contain backend info
	if !strings.Contains(readmeStr, "gin") {
		t.Error("README.md should contain backend technology 'gin'")
	}
	
	// Should contain database info
	if !strings.Contains(readmeStr, "mongo") {
		t.Error("README.md should contain database technology 'mongo'")
	}
	
	// Should NOT contain frontend info
	if strings.Contains(readmeStr, "Frontend Service") {
		t.Error("README.md should not contain frontend information for backend-only project")
	}
}

