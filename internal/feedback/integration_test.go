package feedback

import (
	"bytes"
	"strings"
	"testing"

	"microservice-bootstrapper/internal/interfaces"
)

// TestMessengerIntegration tests the complete messenger workflow
func TestMessengerIntegration(t *testing.T) {
	var buf bytes.Buffer
	messenger := NewMessengerWithOutput(&buf)

	// Test complete workflow
	config := interfaces.CLIConfig{
		Backend:     "fastapi",
		Frontend:    "react",
		Database:    "postgres",
		ProjectName: "test-project",
		Force:       false,
	}

	// Test project name inference
	messenger.ShowProjectNameInference("test-project", "test project!")
	
	// Test Docker warning
	messenger.ShowDockerWarning()
	
	// Test generation success
	messenger.ShowGenerationSuccess("test-project", ".")
	
	// Test next steps
	messenger.ShowNextSteps(config)

	output := buf.String()

	// Verify all expected content is present
	expectedContent := []string{
		"📁 Using directory name as project name: test-project",
		"(sanitized from directory: test project!)",
		"⚠️  Docker Warning:",
		"Docker is not running or not available",
		"✅ Successfully created microservice project 'test-project'!",
		"📋 Next steps:",
		"cd test-project",
		"Backend (fastapi): http://localhost:8000",
		"Frontend (react): http://localhost:3000",
		"Database (postgres): localhost:5432",
		"🚀 Happy coding!",
	}

	for _, expected := range expectedContent {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected output to contain %q, but it was missing from:\n%s", expected, output)
		}
	}
}

// TestMessengerWithMinimalConfig tests messenger with minimal configuration
func TestMessengerWithMinimalConfig(t *testing.T) {
	var buf bytes.Buffer
	messenger := NewMessengerWithOutput(&buf)

	config := interfaces.CLIConfig{
		Backend:     "gin",
		ProjectName: "", // No project name (current directory)
	}

	messenger.ShowNextSteps(config)

	output := buf.String()

	expectedContent := []string{
		"📋 Next steps:",
		"Review the generated files in your current directory",
		"Backend (gin): http://localhost:8080",
		"🚀 Happy coding!",
	}

	for _, expected := range expectedContent {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected output to contain %q, but it was missing from:\n%s", expected, output)
		}
	}

	// Should not contain project-specific navigation
	if strings.Contains(output, "cd ") {
		t.Errorf("Output should not contain 'cd' command when ProjectName is empty, got:\n%s", output)
	}
}