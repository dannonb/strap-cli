package generator

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"microservice-bootstrapper/internal/filesystem"
	"microservice-bootstrapper/internal/interfaces"
	templatepkg "microservice-bootstrapper/internal/template"
	"microservice-bootstrapper/internal/validation"
	"microservice-bootstrapper/pkg/errors"
)

// Generator implements the ProjectGenerator interface
type Generator struct {
	templateEngine interfaces.TemplateEngine
	fileManager    interfaces.FileSystemManager
	validator      *validation.Validator
}

// NewGenerator creates a new project generator instance
func NewGenerator() interfaces.ProjectGenerator {
	return &Generator{
		templateEngine: templatepkg.NewEngine(),
		fileManager:    filesystem.NewManager(),
		validator:      validation.NewValidator(),
	}
}

// Generate orchestrates the complete project generation process
func (g *Generator) Generate(config interfaces.CLIConfig) error {
	// Step 1: Validate configuration and prerequisites
	fmt.Println("🔍 Validating configuration...")
	if err := g.ValidateConfig(config); err != nil {
		return fmt.Errorf("configuration validation failed:\n%w", err)
	}

	fmt.Println("🔧 Checking prerequisites...")
	if err := g.CheckPrerequisites(); err != nil {
		return fmt.Errorf("prerequisite check failed:\n%w", err)
	}

	// Step 2: Determine project name and base path
	projectName := g.getProjectName(config)
	basePath := filepath.Join(".", projectName)

	// Step 3: Check directory and handle conflicts
	fmt.Println("📂 Checking directory conflicts...")
	if err := g.handleDirectoryConflicts(basePath, config.Force); err != nil {
		return err
	}

	fmt.Printf("🚀 Creating microservice project '%s'...\n", projectName)

	// Step 4: Create project structure
	if err := g.createProjectStructure(basePath, config); err != nil {
		g.handleGenerationFailure(basePath, "project structure creation", err)
		return fmt.Errorf("failed to create project structure: %w", err)
	}

	// Step 5: Generate services
	if err := g.generateServices(basePath, config); err != nil {
		g.handleGenerationFailure(basePath, "service generation", err)
		return fmt.Errorf("failed to generate services: %w", err)
	}

	// Step 6: Generate Docker Compose file
	if err := g.generateDockerCompose(basePath, config); err != nil {
		g.handleGenerationFailure(basePath, "Docker Compose generation", err)
		return fmt.Errorf("failed to generate Docker Compose: %w", err)
	}

	// Step 7: Generate project documentation
	if err := g.generateProjectDocumentation(basePath, config); err != nil {
		g.handleGenerationFailure(basePath, "documentation generation", err)
		return fmt.Errorf("failed to generate documentation: %w", err)
	}

	// Success message
	g.printSuccessMessage(projectName, config)

	return nil
}

// ValidateConfig validates the CLI configuration using the enhanced validator
func (g *Generator) ValidateConfig(config interfaces.CLIConfig) error {
	// Use the enhanced validator for comprehensive validation
	if err := g.validator.ValidateConfig(config); err != nil {
		return err
	}

	// Validate Docker Compose configuration
	if err := g.templateEngine.ValidateComposeConfiguration(config); err != nil {
		return errors.NewValidationError("docker-compose", 
			fmt.Sprintf("Docker Compose validation failed: %v. "+
				"Suggestion: Check that your service combination is supported", err))
	}

	return nil
}

// CheckPrerequisites verifies that required tools are available using the enhanced validator
func (g *Generator) CheckPrerequisites() error {
	return g.validator.CheckPrerequisites()
}

// getProjectName determines the project name from config or current directory
func (g *Generator) getProjectName(config interfaces.CLIConfig) string {
	if config.ProjectName != "" {
		return config.ProjectName
	}

	// Use current directory name as default
	cwd, err := os.Getwd()
	if err != nil {
		return "microservice-project"
	}

	return filepath.Base(cwd)
}

// handleDirectoryConflicts checks for existing files and handles conflicts using enhanced validation
func (g *Generator) handleDirectoryConflicts(basePath string, force bool) error {
	return g.validator.ValidateDirectoryConflicts(basePath, force)
}

// createProjectStructure creates the basic project directory structure
func (g *Generator) createProjectStructure(basePath string, config interfaces.CLIConfig) error {
	fmt.Println("📁 Creating project structure...")

	// Create base project directory
	if err := g.fileManager.CreateDirectory(basePath); err != nil {
		return err
	}

	// Create service directories based on configuration
	if config.Backend != "" {
		backendPath := filepath.Join(basePath, "backend")
		if err := g.fileManager.CreateDirectory(backendPath); err != nil {
			return err
		}
	}

	if config.Frontend != "" {
		frontendPath := filepath.Join(basePath, "frontend")
		if err := g.fileManager.CreateDirectory(frontendPath); err != nil {
			return err
		}
	}

	// Create docs directory
	docsPath := filepath.Join(basePath, "docs")
	if err := g.fileManager.CreateDirectory(docsPath); err != nil {
		return err
	}

	return nil
}

// generateServices generates all specified services
func (g *Generator) generateServices(basePath string, config interfaces.CLIConfig) error {
	// Generate backend service
	if config.Backend != "" {
		fmt.Printf("🔧 Generating %s backend service...\n", config.Backend)
		if err := g.generateService(basePath, "backend", config.Backend, config); err != nil {
			return fmt.Errorf("failed to generate backend service: %w", err)
		}
	}

	// Generate frontend service
	if config.Frontend != "" {
		fmt.Printf("🎨 Generating %s frontend service...\n", config.Frontend)
		if err := g.generateService(basePath, "frontend", config.Frontend, config); err != nil {
			return fmt.Errorf("failed to generate frontend service: %w", err)
		}
	}

	return nil
}

// generateService generates a specific service (backend or frontend)
func (g *Generator) generateService(basePath, serviceType, technology string, config interfaces.CLIConfig) error {
	// Get template for this service/technology combination
	tmpl, err := g.templateEngine.GetTemplate(serviceType, technology)
	if err != nil {
		return fmt.Errorf("failed to get template for %s/%s: %w", serviceType, technology, err)
	}

	servicePath := filepath.Join(basePath, serviceType)

	// Create service directories
	for _, dir := range tmpl.Directories {
		dirPath := filepath.Join(servicePath, dir)
		if err := g.fileManager.CreateDirectory(dirPath); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dirPath, err)
		}
	}

	// Create template data
	templateData := g.createTemplateData(config, serviceType, technology)

	// Generate files from templates
	for filename, templateContent := range tmpl.Files {
		filePath := filepath.Join(servicePath, filename)

		// Process template content
		processedContent, err := g.processTemplateContent(templateContent, templateData)
		if err != nil {
			return fmt.Errorf("failed to process template for %s: %w", filename, err)
		}

		// Write file
		if err := g.fileManager.WriteFile(filePath, processedContent); err != nil {
			return fmt.Errorf("failed to write file %s: %w", filePath, err)
		}
	}

	return nil
}

// generateDockerCompose generates the Docker Compose file
func (g *Generator) generateDockerCompose(basePath string, config interfaces.CLIConfig) error {
	fmt.Println("🐳 Generating Docker Compose configuration...")

	composeContent, err := g.templateEngine.GenerateDockerCompose(config)
	if err != nil {
		return fmt.Errorf("failed to generate Docker Compose: %w", err)
	}

	composePath := filepath.Join(basePath, "docker-compose.yml")
	if err := g.fileManager.WriteFile(composePath, composeContent); err != nil {
		return fmt.Errorf("failed to write Docker Compose file: %w", err)
	}

	return nil
}

// generateProjectDocumentation generates project documentation files
func (g *Generator) generateProjectDocumentation(basePath string, config interfaces.CLIConfig) error {
	fmt.Println("📚 Generating project documentation...")

	// Get project template
	projectTemplate, err := g.templateEngine.GetProjectTemplate()
	if err != nil {
		return fmt.Errorf("failed to get project template: %w", err)
	}

	// Create template data
	templateData := g.createTemplateData(config, "", "")

	// Generate files from templates
	for filename, templateContent := range projectTemplate.Files {
		filePath := filepath.Join(basePath, filename)

		// Process template content
		processedContent, err := g.processTemplateContent(templateContent, templateData)
		if err != nil {
			return fmt.Errorf("failed to process template for %s: %w", filename, err)
		}

		// Write file
		if err := g.fileManager.WriteFile(filePath, processedContent); err != nil {
			return fmt.Errorf("failed to write file %s: %w", filePath, err)
		}
	}

	return nil
}

// createTemplateData creates template data for processing
func (g *Generator) createTemplateData(config interfaces.CLIConfig, serviceType, technology string) interfaces.TemplateData {
	projectName := g.getProjectName(config)

	// Create service configurations
	var services []interfaces.ServiceConfig
	var backend, frontend *interfaces.ServiceTemplateData
	var database *interfaces.DatabaseConfig

	// Backend service
	if config.Backend != "" {
		backendService := interfaces.ServiceConfig{
			Type:        "backend",
			Technology:  config.Backend,
			Port:        g.getDefaultPort("backend", config.Backend),
			Environment: make(map[string]string),
		}
		services = append(services, backendService)

		backend = &interfaces.ServiceTemplateData{
			Type:        "backend",
			Technology:  config.Backend,
			Port:        backendService.Port,
			Environment: backendService.Environment,
		}
	}

	// Frontend service
	if config.Frontend != "" {
		frontendService := interfaces.ServiceConfig{
			Type:        "frontend",
			Technology:  config.Frontend,
			Port:        g.getDefaultPort("frontend", config.Frontend),
			Environment: make(map[string]string),
		}
		services = append(services, frontendService)

		frontend = &interfaces.ServiceTemplateData{
			Type:        "frontend",
			Technology:  config.Frontend,
			Port:        frontendService.Port,
			Environment: frontendService.Environment,
		}
	}

	// Database service
	if config.Database != "" {
		database = &interfaces.DatabaseConfig{
			Type:        config.Database,
			Port:        g.getDefaultPort("database", config.Database),
			Volume:      fmt.Sprintf("%s_%s_data", projectName, config.Database),
			Environment: g.getDatabaseEnvironment(config.Database),
		}
	}

	// Port configuration
	ports := interfaces.PortConfig{
		Backend:  g.getDefaultPort("backend", config.Backend),
		Frontend: g.getDefaultPort("frontend", config.Frontend),
		Database: g.getDefaultPort("database", config.Database),
	}

	return interfaces.TemplateData{
		ProjectName: projectName,
		Services:    services,
		Backend:     backend,
		Frontend:    frontend,
		Database:    database,
		Ports:       ports,
		Environment: make(map[string]string),
	}
}

// Helper functions

// checkCommand checks if a command is available
func (g *Generator) checkCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	return cmd.Run()
}

// contains checks if a slice contains a specific string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// getDefaultPort returns the default port for a service type and technology
func (g *Generator) getDefaultPort(serviceType, technology string) int {
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

// getDatabaseEnvironment returns default environment variables for database
func (g *Generator) getDatabaseEnvironment(dbType string) map[string]string {
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

// processTemplateContent processes template content with data
func (g *Generator) processTemplateContent(templateContent string, data interfaces.TemplateData) ([]byte, error) {
	tmpl, err := template.New("content").Parse(templateContent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.Bytes(), nil
}



// handleGenerationFailure handles cleanup and error reporting when generation fails
func (g *Generator) handleGenerationFailure(basePath, operation string, err error) {
	fmt.Printf("❌ Failed during %s: %v\n", operation, err)
	
	// Attempt cleanup
	if cleanupErr := g.fileManager.CleanupOnFailure(basePath); cleanupErr != nil {
		fmt.Printf("⚠️  Warning: Failed to cleanup after error: %v\n", cleanupErr)
		fmt.Printf("💡 You may need to manually remove the directory '%s'\n", basePath)
	}
}

// printSuccessMessage prints the success message with next steps
func (g *Generator) printSuccessMessage(projectName string, config interfaces.CLIConfig) {
	fmt.Printf("\n✅ Successfully created microservice project '%s'!\n\n", projectName)
	
	fmt.Println("📋 Next steps:")
	fmt.Printf("   1. cd %s\n", projectName)
	fmt.Println("   2. cp .env.example .env")
	fmt.Println("   3. docker-compose up -d")
	fmt.Println("   4. docker-compose ps")
	
	fmt.Println("\n🔗 Service URLs:")
	if config.Backend != "" {
		port := g.getDefaultPort("backend", config.Backend)
		fmt.Printf("   Backend (%s): http://localhost:%d\n", config.Backend, port)
	}
	if config.Frontend != "" {
		port := g.getDefaultPort("frontend", config.Frontend)
		fmt.Printf("   Frontend (%s): http://localhost:%d\n", config.Frontend, port)
	}
	if config.Database != "" {
		port := g.getDefaultPort("database", config.Database)
		fmt.Printf("   Database (%s): localhost:%d\n", config.Database, port)
	}
	
	fmt.Println("\n🚀 Happy coding!")
}