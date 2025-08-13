package template

import (
	"strings"
	"testing"

	"microservice-bootstrapper/internal/interfaces"
)

func TestNewEngine(t *testing.T) {
	engine := NewEngine()
	if engine == nil {
		t.Fatal("NewEngine() returned nil")
	}
	if engine.templates == nil {
		t.Fatal("NewEngine() did not initialize templates map")
	}
}

func TestProcessTemplate(t *testing.T) {
	engine := NewEngine()

	tests := []struct {
		name         string
		templateName string
		data         interface{}
		expectError  bool
		contains     string
	}{
		{
			name:         "valid fastapi dockerfile",
			templateName: "backend/fastapi/Dockerfile.tmpl",
			data: interfaces.TemplateData{
				ProjectName: "test-project",
				Ports:       interfaces.PortConfig{Backend: 8000},
			},
			expectError: false,
			contains:    "FROM python:3.11-slim",
		},
		{
			name:         "valid express package.json",
			templateName: "backend/express/package.json.tmpl",
			data: interfaces.TemplateData{
				ProjectName: "my-service",
			},
			expectError: false,
			contains:    "my-service-backend",
		},
		{
			name:         "nonexistent template",
			templateName: "nonexistent/template.tmpl",
			data:         interfaces.TemplateData{},
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := engine.ProcessTemplate(tt.templateName, tt.data)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("ProcessTemplate() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("ProcessTemplate() unexpected error: %v", err)
				return
			}

			if tt.contains != "" && !strings.Contains(string(result), tt.contains) {
				t.Errorf("ProcessTemplate() result does not contain expected string '%s'", tt.contains)
			}
		})
	}
}

func TestGetTemplate(t *testing.T) {
	engine := NewEngine()

	tests := []struct {
		name        string
		service     string
		technology  string
		expectError bool
		checkFiles  []string
		checkDirs   []string
	}{
		{
			name:        "fastapi backend template",
			service:     "backend",
			technology:  "fastapi",
			expectError: false,
			checkFiles:  []string{"Dockerfile", "requirements.txt", "main.py"},
			checkDirs:   []string{"app", "app/api", "app/core", "app/models"},
		},
		{
			name:        "express backend template",
			service:     "backend",
			technology:  "express",
			expectError: false,
			checkFiles:  []string{"Dockerfile", "package.json", "server.js"},
			checkDirs:   []string{"src", "src/routes", "src/middleware", "src/models"},
		},
		{
			name:        "gin backend template",
			service:     "backend",
			technology:  "gin",
			expectError: false,
			checkFiles:  []string{"Dockerfile", "go.mod", "main.go"},
			checkDirs:   []string{"cmd", "internal", "pkg", "api"},
		},
		{
			name:        "react frontend template",
			service:     "frontend",
			technology:  "react",
			expectError: false,
			checkFiles:  []string{"Dockerfile", "package.json", "App.js"},
			checkDirs:   []string{"src", "src/components", "src/pages", "public"},
		},
		{
			name:        "unsupported service",
			service:     "unsupported",
			technology:  "unknown",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template, err := engine.GetTemplate(tt.service, tt.technology)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("GetTemplate() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("GetTemplate() unexpected error: %v", err)
				return
			}

			// Check that expected files are present
			for _, file := range tt.checkFiles {
				if _, exists := template.Files[file]; !exists {
					t.Errorf("GetTemplate() missing expected file: %s", file)
				}
			}

			// Check that expected directories are present
			for _, dir := range tt.checkDirs {
				found := false
				for _, templateDir := range template.Directories {
					if templateDir == dir {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("GetTemplate() missing expected directory: %s", dir)
				}
			}

			// Check that variables are set
			if template.Variables == nil {
				t.Errorf("GetTemplate() variables should not be nil")
			}

			if template.Variables["service"] != tt.service {
				t.Errorf("GetTemplate() service variable incorrect, got %v, want %s", template.Variables["service"], tt.service)
			}

			if template.Variables["technology"] != tt.technology {
				t.Errorf("GetTemplate() technology variable incorrect, got %v, want %s", template.Variables["technology"], tt.technology)
			}
		})
	}
}

func TestGetTemplateDirectories(t *testing.T) {
	engine := NewEngine()

	tests := []struct {
		name       string
		service    string
		technology string
		expected   []string
	}{
		{
			name:       "fastapi backend",
			service:    "backend",
			technology: "fastapi",
			expected:   []string{"app", "app/api", "app/core", "app/models"},
		},
		{
			name:       "express backend",
			service:    "backend",
			technology: "express",
			expected:   []string{"src", "src/routes", "src/middleware", "src/models"},
		},
		{
			name:       "gin backend",
			service:    "backend",
			technology: "gin",
			expected:   []string{"cmd", "internal", "pkg", "api"},
		},
		{
			name:       "react frontend",
			service:    "frontend",
			technology: "react",
			expected:   []string{"src", "src/components", "src/pages", "public"},
		},
		{
			name:       "vue frontend",
			service:    "frontend",
			technology: "vue",
			expected:   []string{"src", "src/components", "src/pages", "public"},
		},
		{
			name:       "angular frontend",
			service:    "frontend",
			technology: "angular",
			expected:   []string{"src", "src/app", "src/assets", "src/environments"},
		},
		{
			name:       "mongo database",
			service:    "database",
			technology: "mongo",
			expected:   []string{"mongo-init"},
		},
		{
			name:       "postgres database",
			service:    "database",
			technology: "postgres",
			expected:   []string{"postgres-init"},
		},
		{
			name:       "mysql database",
			service:    "database",
			technology: "mysql",
			expected:   []string{"mysql-init"},
		},
		{
			name:       "redis database",
			service:    "database",
			technology: "redis",
			expected:   []string{"redis-config"},
		},
		{
			name:       "unknown service",
			service:    "unknown",
			technology: "unknown",
			expected:   []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := engine.getTemplateDirectories(tt.service, tt.technology)
			
			if len(result) != len(tt.expected) {
				t.Errorf("getTemplateDirectories() returned %d directories, expected %d", len(result), len(tt.expected))
				return
			}

			for i, dir := range tt.expected {
				if result[i] != dir {
					t.Errorf("getTemplateDirectories() directory %d = %s, expected %s", i, result[i], dir)
				}
			}
		})
	}
}

func TestGetDefaultPorts(t *testing.T) {
	engine := NewEngine()

	backendTests := []struct {
		technology string
		expected   int
	}{
		{"fastapi", 8000},
		{"express", 3000},
		{"gin", 8080},
		{"unknown", 8000},
	}

	for _, tt := range backendTests {
		t.Run("backend_"+tt.technology, func(t *testing.T) {
			result := engine.getDefaultBackendPort(tt.technology)
			if result != tt.expected {
				t.Errorf("getDefaultBackendPort(%s) = %d, expected %d", tt.technology, result, tt.expected)
			}
		})
	}

	frontendTests := []struct {
		technology string
		expected   int
	}{
		{"react", 3000},
		{"vue", 3000},
		{"angular", 4200},
		{"unknown", 3000},
	}

	for _, tt := range frontendTests {
		t.Run("frontend_"+tt.technology, func(t *testing.T) {
			result := engine.getDefaultFrontendPort(tt.technology)
			if result != tt.expected {
				t.Errorf("getDefaultFrontendPort(%s) = %d, expected %d", tt.technology, result, tt.expected)
			}
		})
	}

	databaseTests := []struct {
		technology string
		expected   int
	}{
		{"mongo", 27017},
		{"postgres", 5432},
		{"mysql", 3306},
		{"redis", 6379},
		{"unknown", 5432},
	}

	for _, tt := range databaseTests {
		t.Run("database_"+tt.technology, func(t *testing.T) {
			result := engine.getDefaultDatabasePort(tt.technology)
			if result != tt.expected {
				t.Errorf("getDefaultDatabasePort(%s) = %d, expected %d", tt.technology, result, tt.expected)
			}
		})
	}
}

func TestTemplateVariables(t *testing.T) {
	engine := NewEngine()

	tests := []struct {
		name       string
		service    string
		technology string
		checkPort  int
	}{
		{
			name:       "fastapi backend variables",
			service:    "backend",
			technology: "fastapi",
			checkPort:  8000,
		},
		{
			name:       "react frontend variables",
			service:    "frontend",
			technology: "react",
			checkPort:  3000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			variables := engine.getTemplateVariables(tt.service, tt.technology)
			
			if variables["service"] != tt.service {
				t.Errorf("getTemplateVariables() service = %v, expected %s", variables["service"], tt.service)
			}

			if variables["technology"] != tt.technology {
				t.Errorf("getTemplateVariables() technology = %v, expected %s", variables["technology"], tt.technology)
			}

			if port, ok := variables["port"].(int); ok {
				if port != tt.checkPort {
					t.Errorf("getTemplateVariables() port = %d, expected %d", port, tt.checkPort)
				}
			} else {
				t.Errorf("getTemplateVariables() port not found or not an int")
			}
		})
	}
}