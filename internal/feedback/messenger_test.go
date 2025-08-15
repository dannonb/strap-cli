package feedback

import (
	"bytes"
	"strings"
	"testing"

	"microservice-bootstrapper/internal/interfaces"
)

func TestNewMessenger(t *testing.T) {
	messenger := NewMessenger()
	if messenger == nil {
		t.Fatal("NewMessenger() returned nil")
	}
}

func TestNewMessengerWithOutput(t *testing.T) {
	var buf bytes.Buffer
	messenger := NewMessengerWithOutput(&buf)
	if messenger == nil {
		t.Fatal("NewMessengerWithOutput() returned nil")
	}
}

func TestShowProjectNameInference(t *testing.T) {
	tests := []struct {
		name      string
		inferred  string
		directory string
		expected  []string
	}{
		{
			name:      "same name",
			inferred:  "my-project",
			directory: "my-project",
			expected:  []string{"📁 Using directory name as project name: my-project"},
		},
		{
			name:      "sanitized name",
			inferred:  "my-project",
			directory: "my project!",
			expected: []string{
				"📁 Using directory name as project name: my-project",
				"(sanitized from directory: my project!)",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			messenger := NewMessengerWithOutput(&buf)
			
			messenger.ShowProjectNameInference(tt.inferred, tt.directory)
			
			output := buf.String()
			for _, expected := range tt.expected {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain %q, got: %s", expected, output)
				}
			}
		})
	}
}

func TestShowDockerWarning(t *testing.T) {
	var buf bytes.Buffer
	messenger := NewMessengerWithOutput(&buf)
	
	messenger.ShowDockerWarning()
	
	output := buf.String()
	expectedPhrases := []string{
		"⚠️  Docker Warning:",
		"Docker is not running or not available",
		"Project files will be generated successfully",
		"https://docs.docker.com/get-docker/",
		"docker --version",
	}
	
	for _, phrase := range expectedPhrases {
		if !strings.Contains(output, phrase) {
			t.Errorf("Expected output to contain %q, got: %s", phrase, output)
		}
	}
}

func TestShowGenerationSuccess(t *testing.T) {
	tests := []struct {
		name        string
		projectName string
		path        string
		expected    []string
	}{
		{
			name:        "current directory",
			projectName: "my-project",
			path:        ".",
			expected:    []string{"✅ Successfully created microservice project 'my-project'!"},
		},
		{
			name:        "specific path",
			projectName: "my-project",
			path:        "/path/to/project",
			expected: []string{
				"✅ Successfully created microservice project 'my-project'!",
				"📂 Location: /path/to/project",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			messenger := NewMessengerWithOutput(&buf)
			
			messenger.ShowGenerationSuccess(tt.projectName, tt.path)
			
			output := buf.String()
			for _, expected := range tt.expected {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain %q, got: %s", expected, output)
				}
			}
		})
	}
}

func TestShowNextSteps(t *testing.T) {
	tests := []struct {
		name     string
		config   interfaces.CLIConfig
		expected []string
	}{
		{
			name: "full stack with project name",
			config: interfaces.CLIConfig{
				Backend:     "fastapi",
				Frontend:    "react",
				Database:    "postgres",
				ProjectName: "my-project",
			},
			expected: []string{
				"📋 Next steps:",
				"cd my-project",
				"cp .env.example .env",
				"docker-compose up -d",
				"🔗 Service URLs:",
				"Backend (fastapi): http://localhost:8000",
				"Frontend (react): http://localhost:3000",
				"Database (postgres): localhost:5432",
				"🚀 Happy coding!",
			},
		},
		{
			name: "backend only without project name",
			config: interfaces.CLIConfig{
				Backend:     "gin",
				ProjectName: "",
			},
			expected: []string{
				"📋 Next steps:",
				"Review the generated files in your current directory",
				"cp .env.example .env",
				"Backend (gin): http://localhost:8080",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			messenger := NewMessengerWithOutput(&buf)
			
			messenger.ShowNextSteps(tt.config)
			
			output := buf.String()
			for _, expected := range tt.expected {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain %q, got: %s", expected, output)
				}
			}
		})
	}
}

func TestGetDefaultPort(t *testing.T) {
	tests := []struct {
		serviceType string
		technology  string
		expected    int
	}{
		{"backend", "fastapi", 8000},
		{"backend", "express", 3000},
		{"backend", "gin", 8080},
		{"backend", "unknown", 8000},
		{"frontend", "react", 3000},
		{"frontend", "vue", 3000},
		{"frontend", "angular", 4200},
		{"frontend", "unknown", 3000},
		{"database", "postgres", 5432},
		{"database", "mysql", 3306},
		{"database", "mongo", 27017},
		{"database", "redis", 6379},
		{"database", "unknown", 5432},
		{"unknown", "unknown", 8000},
	}

	for _, tt := range tests {
		t.Run(tt.serviceType+"_"+tt.technology, func(t *testing.T) {
			result := getDefaultPort(tt.serviceType, tt.technology)
			if result != tt.expected {
				t.Errorf("getDefaultPort(%q, %q) = %d, expected %d", 
					tt.serviceType, tt.technology, result, tt.expected)
			}
		})
	}
}