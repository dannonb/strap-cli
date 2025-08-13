package template

import (
	"bytes"
	"strings"
	"testing"
	"text/template"

	"microservice-bootstrapper/internal/interfaces"
)

func TestGetProjectTemplate(t *testing.T) {
	engine := NewEngine()

	template, err := engine.GetProjectTemplate()
	if err != nil {
		t.Fatalf("Failed to get project template: %v", err)
	}

	// Check that all expected files are present
	expectedFiles := []string{
		"README.md",
		".env.example",
		".gitignore",
	}

	for _, expectedFile := range expectedFiles {
		if _, exists := template.Files[expectedFile]; !exists {
			t.Errorf("Expected file %s not found in project template", expectedFile)
		}
	}

	// Verify that files have content
	for filename, content := range template.Files {
		if strings.TrimSpace(content) == "" {
			t.Errorf("File %s has empty content", filename)
		}
	}
}

func TestProjectDocumentationGeneration(t *testing.T) {
	engine := NewEngine()

	testCases := []struct {
		name   string
		config interfaces.CLIConfig
	}{
		{
			name: "full_stack_project",
			config: interfaces.CLIConfig{
				ProjectName: "test-project",
				Backend:     "fastapi",
				Frontend:    "react",
				Database:    "postgres",
			},
		},
		{
			name: "backend_only_project",
			config: interfaces.CLIConfig{
				ProjectName: "api-service",
				Backend:     "gin",
				Database:    "mongo",
			},
		},
		{
			name: "frontend_only_project",
			config: interfaces.CLIConfig{
				ProjectName: "web-app",
				Frontend:    "vue",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Get project template
			template, err := engine.GetProjectTemplate()
			if err != nil {
				t.Fatalf("Failed to get project template: %v", err)
			}

			// Create template data
			templateData := createTestTemplateData(tc.config)

			// Test README.md generation
			readmeContent, exists := template.Files["README.md"]
			if !exists {
				t.Fatal("README.md template not found")
			}

			processedReadme, err := processTestTemplate(readmeContent, templateData)
			if err != nil {
				t.Fatalf("Failed to process README template: %v", err)
			}

			// Verify README contains project name
			if !strings.Contains(processedReadme, tc.config.ProjectName) {
				t.Errorf("README does not contain project name: %s", tc.config.ProjectName)
			}

			// Verify backend information is included if specified
			if tc.config.Backend != "" {
				if !strings.Contains(processedReadme, tc.config.Backend) {
					t.Errorf("README does not contain backend technology: %s", tc.config.Backend)
				}
			}

			// Verify frontend information is included if specified
			if tc.config.Frontend != "" {
				if !strings.Contains(processedReadme, tc.config.Frontend) {
					t.Errorf("README does not contain frontend technology: %s", tc.config.Frontend)
				}
			}

			// Verify database information is included if specified
			if tc.config.Database != "" {
				if !strings.Contains(processedReadme, tc.config.Database) {
					t.Errorf("README does not contain database technology: %s", tc.config.Database)
				}
			}

			// Test .env.example generation
			envContent, exists := template.Files[".env.example"]
			if !exists {
				t.Fatal(".env.example template not found")
			}

			processedEnv, err := processTestTemplate(envContent, templateData)
			if err != nil {
				t.Fatalf("Failed to process .env.example template: %v", err)
			}

			// Verify .env.example contains project name
			if !strings.Contains(processedEnv, tc.config.ProjectName) {
				t.Errorf(".env.example does not contain project name: %s", tc.config.ProjectName)
			}

			// Test .gitignore generation
			gitignoreContent, exists := template.Files[".gitignore"]
			if !exists {
				t.Fatal(".gitignore template not found")
			}

			processedGitignore, err := processTestTemplate(gitignoreContent, templateData)
			if err != nil {
				t.Fatalf("Failed to process .gitignore template: %v", err)
			}

			// Verify .gitignore contains project name
			if !strings.Contains(processedGitignore, tc.config.ProjectName) {
				t.Errorf(".gitignore does not contain project name: %s", tc.config.ProjectName)
			}

			// Verify technology-specific ignores are included
			if tc.config.Backend == "fastapi" && !strings.Contains(processedGitignore, "__pycache__") {
				t.Error(".gitignore does not contain Python-specific ignores for FastAPI backend")
			}

			if tc.config.Backend == "express" && !strings.Contains(processedGitignore, "node_modules") {
				t.Error(".gitignore does not contain Node.js-specific ignores for Express backend")
			}

			if tc.config.Backend == "gin" && !strings.Contains(processedGitignore, "*.exe") {
				t.Error(".gitignore does not contain Go-specific ignores for Gin backend")
			}
		})
	}
}

func TestProjectDocumentationTemplateContent(t *testing.T) {
	engine := NewEngine()

	template, err := engine.GetProjectTemplate()
	if err != nil {
		t.Fatalf("Failed to get project template: %v", err)
	}

	// Test README.md template structure
	readmeContent := template.Files["README.md"]
	
	// Check for essential sections
	essentialSections := []string{
		"Architecture Overview",
		"Getting Started",
		"Prerequisites",
		"Quick Start",
		"Service URLs",
		"Development",
		"Project Structure",
		"Environment Configuration",
		"Troubleshooting",
	}

	for _, section := range essentialSections {
		if !strings.Contains(readmeContent, section) {
			t.Errorf("README template missing essential section: %s", section)
		}
	}

	// Test .env.example template structure
	envContent := template.Files[".env.example"]
	
	// Check for essential environment sections
	envSections := []string{
		"SERVICE PORTS",
		"DATABASE CONFIGURATION", 
		"SECURITY SETTINGS",
		"LOGGING AND MONITORING",
		"CORS SETTINGS",
		"DEVELOPMENT SETTINGS",
		"PRODUCTION OVERRIDES",
	}

	for _, section := range envSections {
		if !strings.Contains(envContent, section) {
			t.Errorf(".env.example template missing essential section: %s", section)
		}
	}

	// Test .gitignore template structure
	gitignoreContent := template.Files[".gitignore"]
	
	// Check for essential ignore categories
	ignoreCategories := []string{
		"ENVIRONMENT AND CONFIGURATION FILES",
		"LOGS AND TEMPORARY FILES",
		"DOCKER AND CONTAINERIZATION",
		"OPERATING SYSTEM FILES",
		"IDE AND EDITOR FILES",
		"BACKEND SPECIFIC FILES",
		"FRONTEND SPECIFIC FILES",
		"DATABASE FILES",
		"TESTING AND COVERAGE",
		"DEPLOYMENT AND CI/CD",
		"SECURITY AND SECRETS",
	}

	for _, category := range ignoreCategories {
		if !strings.Contains(gitignoreContent, category) {
			t.Errorf(".gitignore template missing essential category: %s", category)
		}
	}
}

// Helper functions for testing

func createTestTemplateData(config interfaces.CLIConfig) interfaces.TemplateData {
	var backend *interfaces.ServiceTemplateData
	var frontend *interfaces.ServiceTemplateData
	var database *interfaces.DatabaseConfig

	if config.Backend != "" {
		backend = &interfaces.ServiceTemplateData{
			Type:        "backend",
			Technology:  config.Backend,
			Port:        getTestDefaultPort("backend", config.Backend),
			Environment: make(map[string]string),
		}
	}

	if config.Frontend != "" {
		frontend = &interfaces.ServiceTemplateData{
			Type:        "frontend",
			Technology:  config.Frontend,
			Port:        getTestDefaultPort("frontend", config.Frontend),
			Environment: make(map[string]string),
		}
	}

	if config.Database != "" {
		database = &interfaces.DatabaseConfig{
			Type:        config.Database,
			Port:        getTestDefaultPort("database", config.Database),
			Volume:      config.ProjectName + "_" + config.Database + "_data",
			Environment: getTestDatabaseEnvironment(config.Database),
		}
	}

	return interfaces.TemplateData{
		ProjectName: config.ProjectName,
		Backend:     backend,
		Frontend:    frontend,
		Database:    database,
		Ports: interfaces.PortConfig{
			Backend:  getTestDefaultPort("backend", config.Backend),
			Frontend: getTestDefaultPort("frontend", config.Frontend),
			Database: getTestDefaultPort("database", config.Database),
		},
		Environment: make(map[string]string),
	}
}

func getTestDefaultPort(serviceType, technology string) int {
	switch serviceType {
	case "backend":
		switch technology {
		case "fastapi":
			return 8000
		case "express":
			return 3000
		case "gin":
			return 8080
		default:
			return 8000
		}
	case "frontend":
		switch technology {
		case "react", "vue":
			return 3000
		case "angular":
			return 4200
		default:
			return 3000
		}
	case "database":
		switch technology {
		case "mongo":
			return 27017
		case "postgres":
			return 5432
		case "mysql":
			return 3306
		case "redis":
			return 6379
		default:
			return 5432
		}
	}
	return 8000
}

func getTestDatabaseEnvironment(dbType string) map[string]string {
	env := make(map[string]string)

	switch dbType {
	case "postgres":
		env["POSTGRES_DB"] = "microservice_db"
		env["POSTGRES_USER"] = "postgres"
		env["POSTGRES_PASSWORD"] = "password"
	case "mysql":
		env["MYSQL_DATABASE"] = "microservice_db"
		env["MYSQL_USER"] = "mysql"
		env["MYSQL_PASSWORD"] = "password"
		env["MYSQL_ROOT_PASSWORD"] = "rootpassword"
	case "mongo":
		env["MONGO_INITDB_DATABASE"] = "microservice_db"
		env["MONGO_INITDB_ROOT_USERNAME"] = "mongo"
		env["MONGO_INITDB_ROOT_PASSWORD"] = "password"
	case "redis":
		env["REDIS_PASSWORD"] = "password"
	}

	return env
}

func processTestTemplate(templateContent string, data interfaces.TemplateData) (string, error) {
	tmpl, err := template.New("test").Parse(templateContent)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}