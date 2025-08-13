package template

import (
	"fmt"
	"strings"

	"microservice-bootstrapper/internal/interfaces"
)

// ComposeGenerator handles Docker Compose file generation
type ComposeGenerator struct {
	engine *Engine
}

// NewComposeGenerator creates a new Docker Compose generator
func NewComposeGenerator(engine *Engine) *ComposeGenerator {
	return &ComposeGenerator{
		engine: engine,
	}
}

// GenerateDockerCompose generates a complete docker-compose.yml file
func (c *ComposeGenerator) GenerateDockerCompose(config interfaces.CLIConfig) ([]byte, error) {
	templateData, err := c.buildTemplateData(config)
	if err != nil {
		return nil, fmt.Errorf("failed to build template data: %w", err)
	}

	return c.engine.ProcessTemplate("docker-compose.yml.tmpl", templateData)
}

// buildTemplateData constructs the template data for Docker Compose generation
func (c *ComposeGenerator) buildTemplateData(config interfaces.CLIConfig) (*interfaces.TemplateData, error) {
	data := &interfaces.TemplateData{
		ProjectName: config.ProjectName,
		Ports:       c.buildPortConfig(config),
		Environment: make(map[string]string),
	}

	// Build backend service data
	if config.Backend != "" {
		backend, err := c.buildBackendData(config.Backend)
		if err != nil {
			return nil, fmt.Errorf("failed to build backend data: %w", err)
		}
		data.Backend = backend
	}

	// Build frontend service data
	if config.Frontend != "" {
		frontend, err := c.buildFrontendData(config.Frontend)
		if err != nil {
			return nil, fmt.Errorf("failed to build frontend data: %w", err)
		}
		data.Frontend = frontend
	}

	// Build database service data
	if config.Database != "" {
		database, err := c.buildDatabaseData(config.Database, config.ProjectName)
		if err != nil {
			return nil, fmt.Errorf("failed to build database data: %w", err)
		}
		data.Database = database
	}

	return data, nil
}

// buildPortConfig creates port configuration for all services
func (c *ComposeGenerator) buildPortConfig(config interfaces.CLIConfig) interfaces.PortConfig {
	ports := interfaces.PortConfig{}

	if config.Backend != "" {
		ports.Backend = c.engine.getDefaultBackendPort(config.Backend)
	}

	if config.Frontend != "" {
		ports.Frontend = c.engine.getDefaultFrontendPort(config.Frontend)
	}

	if config.Database != "" {
		ports.Database = c.engine.getDefaultDatabasePort(config.Database)
	}

	return ports
}

// buildBackendData creates backend service template data
func (c *ComposeGenerator) buildBackendData(technology string) (*interfaces.ServiceTemplateData, error) {
	port := c.engine.getDefaultBackendPort(technology)
	env := c.getBackendEnvironment(technology, port)

	return &interfaces.ServiceTemplateData{
		Type:        "backend",
		Technology:  technology,
		Port:        port,
		Environment: env,
	}, nil
}

// buildFrontendData creates frontend service template data
func (c *ComposeGenerator) buildFrontendData(technology string) (*interfaces.ServiceTemplateData, error) {
	port := c.engine.getDefaultFrontendPort(technology)
	env := c.getFrontendEnvironment(technology, port)

	return &interfaces.ServiceTemplateData{
		Type:        "frontend",
		Technology:  technology,
		Port:        port,
		Environment: env,
	}, nil
}

// buildDatabaseData creates database service template data
func (c *ComposeGenerator) buildDatabaseData(dbType, projectName string) (*interfaces.DatabaseConfig, error) {
	port := c.engine.getDefaultDatabasePort(dbType)
	env := c.getDatabaseEnvironment(dbType, projectName)
	volume := c.getDatabaseVolume(dbType, projectName)

	return &interfaces.DatabaseConfig{
		Type:        dbType,
		Port:        port,
		Volume:      volume,
		Environment: env,
	}, nil
}

// getBackendEnvironment returns environment variables for backend services
func (c *ComposeGenerator) getBackendEnvironment(technology string, port int) map[string]string {
	env := make(map[string]string)

	// Common environment variables
	env["NODE_ENV"] = "development"
	env["PORT"] = fmt.Sprintf("%d", port)

	// Technology-specific environment variables
	switch technology {
	case "fastapi":
		env["PYTHONPATH"] = "/app"
		env["UVICORN_HOST"] = "0.0.0.0"
		env["UVICORN_PORT"] = fmt.Sprintf("%d", port)
	case "express":
		env["NODE_ENV"] = "development"
	case "gin":
		env["GIN_MODE"] = "debug"
		env["GO_ENV"] = "development"
	}

	return env
}

// getFrontendEnvironment returns environment variables for frontend services
func (c *ComposeGenerator) getFrontendEnvironment(technology string, port int) map[string]string {
	env := make(map[string]string)

	// Common environment variables
	env["NODE_ENV"] = "development"
	env["PORT"] = fmt.Sprintf("%d", port)

	// Technology-specific environment variables
	switch technology {
	case "react":
		env["REACT_APP_ENV"] = "development"
		env["GENERATE_SOURCEMAP"] = "true"
	case "vue":
		env["VUE_APP_ENV"] = "development"
	case "angular":
		env["NG_ENV"] = "development"
	}

	return env
}

// getDatabaseEnvironment returns environment variables for database services
func (c *ComposeGenerator) getDatabaseEnvironment(dbType, projectName string) map[string]string {
	env := make(map[string]string)
	dbName := c.sanitizeProjectName(projectName)

	switch dbType {
	case "postgres":
		env["POSTGRES_DB"] = dbName
		env["POSTGRES_USER"] = "postgres"
		env["POSTGRES_PASSWORD"] = "postgres"
	case "mongo":
		env["MONGO_INITDB_ROOT_USERNAME"] = "admin"
		env["MONGO_INITDB_ROOT_PASSWORD"] = "admin123"
		env["MONGO_INITDB_DATABASE"] = dbName
	case "mysql":
		env["MYSQL_ROOT_PASSWORD"] = "rootpassword"
		env["MYSQL_DATABASE"] = dbName
		env["MYSQL_USER"] = "user"
		env["MYSQL_PASSWORD"] = "password"
	case "redis":
		// Redis doesn't need authentication in development
	}

	return env
}

// getDatabaseVolume returns the volume name for database persistence
func (c *ComposeGenerator) getDatabaseVolume(dbType, projectName string) string {
	sanitized := c.sanitizeProjectName(projectName)
	return fmt.Sprintf("%s-%s-data", sanitized, dbType)
}

// sanitizeProjectName removes invalid characters from project name for use in Docker resources
func (c *ComposeGenerator) sanitizeProjectName(name string) string {
	// Replace invalid characters with hyphens and convert to lowercase
	sanitized := strings.ToLower(name)
	sanitized = strings.ReplaceAll(sanitized, "_", "-")
	sanitized = strings.ReplaceAll(sanitized, " ", "-")
	
	// Remove any characters that aren't alphanumeric or hyphens
	var result strings.Builder
	for _, r := range sanitized {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}
	
	return result.String()
}

// ValidateConfiguration validates the Docker Compose configuration
func (c *ComposeGenerator) ValidateConfiguration(config interfaces.CLIConfig) error {
	// Validate that at least one service is specified
	if config.Backend == "" && config.Frontend == "" && config.Database == "" {
		return fmt.Errorf("at least one service (backend, frontend, or database) must be specified")
	}

	// Validate backend technology
	if config.Backend != "" {
		if !c.isValidBackendTechnology(config.Backend) {
			return fmt.Errorf("unsupported backend technology: %s", config.Backend)
		}
	}

	// Validate frontend technology
	if config.Frontend != "" {
		if !c.isValidFrontendTechnology(config.Frontend) {
			return fmt.Errorf("unsupported frontend technology: %s", config.Frontend)
		}
	}

	// Validate database technology
	if config.Database != "" {
		if !c.isValidDatabaseTechnology(config.Database) {
			return fmt.Errorf("unsupported database technology: %s", config.Database)
		}
	}

	// Validate project name
	if config.ProjectName == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	return nil
}

// isValidBackendTechnology checks if the backend technology is supported
func (c *ComposeGenerator) isValidBackendTechnology(tech string) bool {
	validTechs := []string{"fastapi", "express", "gin"}
	for _, valid := range validTechs {
		if tech == valid {
			return true
		}
	}
	return false
}

// isValidFrontendTechnology checks if the frontend technology is supported
func (c *ComposeGenerator) isValidFrontendTechnology(tech string) bool {
	validTechs := []string{"react", "vue", "angular"}
	for _, valid := range validTechs {
		if tech == valid {
			return true
		}
	}
	return false
}

// isValidDatabaseTechnology checks if the database technology is supported
func (c *ComposeGenerator) isValidDatabaseTechnology(tech string) bool {
	validTechs := []string{"postgres", "mongo", "mysql", "redis"}
	for _, valid := range validTechs {
		if tech == valid {
			return true
		}
	}
	return false
}