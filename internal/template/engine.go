package template

import (
	"bytes"
	"embed"
	"fmt"
	"strings"
	"text/template"

	"microservice-bootstrapper/internal/interfaces"
)

//go:embed templates
var templateFS embed.FS

// Engine implements the TemplateEngine interface
type Engine struct {
	templates       map[string]*template.Template
	composeGenerator *ComposeGenerator
}

// NewEngine creates a new template engine instance
func NewEngine() *Engine {
	engine := &Engine{
		templates: make(map[string]*template.Template),
	}
	engine.composeGenerator = NewComposeGenerator(engine)
	return engine
}

// ProcessTemplate processes a template with the given data
func (e *Engine) ProcessTemplate(templateName string, data interface{}) ([]byte, error) {
	tmpl, err := e.getOrLoadTemplate(templateName)
	if err != nil {
		return nil, NewTemplateError(templateName, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, NewTemplateError(templateName, err)
	}

	return buf.Bytes(), nil
}

// GetTemplate retrieves a template configuration for a specific service and technology
func (e *Engine) GetTemplate(service, technology string) (interfaces.Template, error) {
	templateKey := fmt.Sprintf("%s/%s", service, technology)
	
	// Load template files for this service/technology combination
	files, err := e.loadTemplateFiles(templateKey)
	if err != nil {
		return interfaces.Template{}, NewTemplateNotFoundError(service, technology)
	}

	// Get directories that should be created
	directories := e.getTemplateDirectories(service, technology)

	// Get default variables for this template
	variables := e.getTemplateVariables(service, technology)

	return interfaces.Template{
		Files:       files,
		Directories: directories,
		Variables:   variables,
	}, nil
}

// GetProjectTemplate retrieves project-level documentation templates
func (e *Engine) GetProjectTemplate() (interfaces.Template, error) {
	templateKey := "project"
	
	// Load template files for project documentation
	files, err := e.loadTemplateFiles(templateKey)
	if err != nil {
		return interfaces.Template{}, NewTemplateNotFoundError("project", "documentation")
	}

	return interfaces.Template{
		Files:       files,
		Directories: []string{}, // No additional directories needed for project templates
		Variables:   make(map[string]interface{}),
	}, nil
}

// getOrLoadTemplate loads a template from the embedded filesystem or returns cached version
func (e *Engine) getOrLoadTemplate(templateName string) (*template.Template, error) {
	if tmpl, exists := e.templates[templateName]; exists {
		return tmpl, nil
	}

	templatePath := fmt.Sprintf("templates/%s", templateName)
	content, err := templateFS.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("template file not found: %s", templatePath)
	}

	tmpl, err := template.New(templateName).Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	e.templates[templateName] = tmpl
	return tmpl, nil
}

// loadTemplateFiles loads all template files for a given service/technology
func (e *Engine) loadTemplateFiles(templateKey string) (map[string]string, error) {
	files := make(map[string]string)
	templateDir := fmt.Sprintf("templates/%s", templateKey)

	entries, err := templateFS.ReadDir(templateDir)
	if err != nil {
		return nil, fmt.Errorf("template directory not found: %s", templateDir)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filePath := fmt.Sprintf("%s/%s", templateDir, entry.Name())
		content, err := templateFS.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read template file %s: %w", filePath, err)
		}

		// Remove .tmpl extension from filename if present
		filename := entry.Name()
		if strings.HasSuffix(filename, ".tmpl") {
			filename = strings.TrimSuffix(filename, ".tmpl")
		}

		// Handle special cases for project template files that need dot prefixes
		if templateKey == "project" {
			if filename == "env.example" {
				filename = ".env.example"
			} else if filename == "gitignore" {
				filename = ".gitignore"
			}
		}

		files[filename] = string(content)
	}

	return files, nil
}

// getTemplateDirectories returns the directories that should be created for a template
func (e *Engine) getTemplateDirectories(service, technology string) []string {
	switch service {
	case "backend":
		return e.getBackendDirectories(technology)
	case "frontend":
		return e.getFrontendDirectories(technology)
	case "database":
		return e.getDatabaseDirectories(technology)
	default:
		return []string{}
	}
}

// getBackendDirectories returns directories for backend services
func (e *Engine) getBackendDirectories(technology string) []string {
	switch technology {
	case "fastapi":
		return []string{"app", "app/api", "app/core", "app/models"}
	case "express":
		return []string{"src", "src/routes", "src/middleware", "src/models"}
	case "gin":
		return []string{"cmd", "internal", "pkg", "api"}
	default:
		return []string{"src"}
	}
}

// getFrontendDirectories returns directories for frontend services
func (e *Engine) getFrontendDirectories(technology string) []string {
	switch technology {
	case "react", "vue":
		return []string{"src", "src/components", "src/pages", "public"}
	case "angular":
		return []string{"src", "src/app", "src/assets", "src/environments"}
	default:
		return []string{"src", "public"}
	}
}

// getDatabaseDirectories returns directories for database services
func (e *Engine) getDatabaseDirectories(technology string) []string {
	switch technology {
	case "mongo":
		return []string{"mongo-init"}
	case "postgres":
		return []string{"postgres-init"}
	case "mysql":
		return []string{"mysql-init"}
	case "redis":
		return []string{"redis-config"}
	default:
		return []string{}
	}
}

// getTemplateVariables returns default variables for a template
func (e *Engine) getTemplateVariables(service, technology string) map[string]interface{} {
	variables := make(map[string]interface{})
	
	// Common variables
	variables["service"] = service
	variables["technology"] = technology
	
	// Service-specific variables
	switch service {
	case "backend":
		variables["port"] = e.getDefaultBackendPort(technology)
	case "frontend":
		variables["port"] = e.getDefaultFrontendPort(technology)
	case "database":
		variables["port"] = e.getDefaultDatabasePort(technology)
	}
	
	return variables
}

// getDefaultBackendPort returns the default port for backend technologies
func (e *Engine) getDefaultBackendPort(technology string) int {
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
}

// getDefaultFrontendPort returns the default port for frontend technologies
func (e *Engine) getDefaultFrontendPort(technology string) int {
	switch technology {
	case "react", "vue":
		return 3000
	case "angular":
		return 4200
	default:
		return 3000
	}
}

// getDefaultDatabasePort returns the default port for database technologies
func (e *Engine) getDefaultDatabasePort(technology string) int {
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

// GenerateDockerCompose generates a complete docker-compose.yml file
func (e *Engine) GenerateDockerCompose(config interfaces.CLIConfig) ([]byte, error) {
	return e.composeGenerator.GenerateDockerCompose(config)
}

// ValidateComposeConfiguration validates the Docker Compose configuration
func (e *Engine) ValidateComposeConfiguration(config interfaces.CLIConfig) error {
	return e.composeGenerator.ValidateConfiguration(config)
}