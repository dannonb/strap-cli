package template

import (
	"fmt"
	"strings"
	"testing"

	"microservice-bootstrapper/internal/interfaces"
)

// TestDockerComposeIntegration tests the complete Docker Compose generation workflow
func TestDockerComposeIntegration(t *testing.T) {
	engine := NewEngine()

	tests := []struct {
		name           string
		config         interfaces.CLIConfig
		expectedServices []string
		expectedNetworks []string
		expectedVolumes  []string
		expectedPorts    []string
		expectedEnvVars  []string
	}{
		{
			name: "full stack application",
			config: interfaces.CLIConfig{
				Backend:     "fastapi",
				Frontend:    "react",
				Database:    "postgres",
				ProjectName: "my-awesome-app",
			},
			expectedServices: []string{
				"backend:",
				"frontend:",
				"postgres:",
			},
			expectedNetworks: []string{
				"my-awesome-app-network:",
				"driver: bridge",
			},
			expectedVolumes: []string{
				"my-awesome-app-postgres-data:",
				"driver: local",
			},
			expectedPorts: []string{
				"8000:8000",  // FastAPI backend
				"3000:3000",  // React frontend
				"5432:5432",  // PostgreSQL
			},
			expectedEnvVars: []string{
				"UVICORN_HOST: 0.0.0.0",
				"REACT_APP_ENV: development",
				"POSTGRES_DB: my-awesome-app",
				"DATABASE_URL: postgresql://",
			},
		},
		{
			name: "backend with mongodb",
			config: interfaces.CLIConfig{
				Backend:     "gin",
				Database:    "mongo",
				ProjectName: "api-service",
			},
			expectedServices: []string{
				"backend:",
				"mongo:",
			},
			expectedNetworks: []string{
				"api-service-network:",
			},
			expectedVolumes: []string{
				"api-service-mongo-data:",
			},
			expectedPorts: []string{
				"8080:8080",  // Gin backend
				"27017:27017", // MongoDB
			},
			expectedEnvVars: []string{
				"GIN_MODE: debug",
				"MONGO_INITDB_ROOT_USERNAME: admin",
				"MONGODB_URL: mongodb://",
			},
		},
		{
			name: "frontend only",
			config: interfaces.CLIConfig{
				Frontend:    "vue",
				ProjectName: "frontend-app",
			},
			expectedServices: []string{
				"frontend:",
			},
			expectedNetworks: []string{
				"frontend-app-network:",
			},
			expectedVolumes: []string{},
			expectedPorts: []string{
				"3000:3000", // Vue frontend
			},
			expectedEnvVars: []string{
				"VUE_APP_ENV: development",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Generate Docker Compose file
			result, err := engine.GenerateDockerCompose(tt.config)
			if err != nil {
				t.Fatalf("GenerateDockerCompose() failed: %v", err)
			}

			resultStr := string(result)

			// Validate services
			for _, service := range tt.expectedServices {
				if !strings.Contains(resultStr, service) {
					t.Errorf("Generated Docker Compose missing expected service: %s", service)
				}
			}

			// Validate networks
			for _, network := range tt.expectedNetworks {
				if !strings.Contains(resultStr, network) {
					t.Errorf("Generated Docker Compose missing expected network config: %s", network)
				}
			}

			// Validate volumes
			for _, volume := range tt.expectedVolumes {
				if !strings.Contains(resultStr, volume) {
					t.Errorf("Generated Docker Compose missing expected volume config: %s", volume)
				}
			}

			// Validate ports
			for _, port := range tt.expectedPorts {
				if !strings.Contains(resultStr, port) {
					t.Errorf("Generated Docker Compose missing expected port mapping: %s", port)
				}
			}

			// Validate environment variables
			for _, envVar := range tt.expectedEnvVars {
				if !strings.Contains(resultStr, envVar) {
					t.Errorf("Generated Docker Compose missing expected environment variable: %s", envVar)
				}
			}

			// Validate that the generated YAML is well-formed
			if !strings.Contains(resultStr, "version: '3.8'") {
				t.Error("Generated Docker Compose missing version declaration")
			}

			if !strings.Contains(resultStr, "services:") {
				t.Error("Generated Docker Compose missing services section")
			}

			// Log the generated content for debugging if needed
			t.Logf("Generated Docker Compose for %s:\n%s", tt.name, resultStr)
		})
	}
}

// TestDockerComposeValidation tests the validation logic
func TestDockerComposeValidation(t *testing.T) {
	engine := NewEngine()

	validConfigs := []interfaces.CLIConfig{
		{Backend: "fastapi", ProjectName: "test"},
		{Frontend: "react", ProjectName: "test"},
		{Database: "postgres", ProjectName: "test"},
		{Backend: "gin", Frontend: "vue", Database: "mongo", ProjectName: "full-stack"},
	}

	for i, config := range validConfigs {
		t.Run(fmt.Sprintf("valid_config_%d", i), func(t *testing.T) {
			err := engine.ValidateComposeConfiguration(config)
			if err != nil {
				t.Errorf("ValidateComposeConfiguration() failed for valid config: %v", err)
			}
		})
	}

	invalidConfigs := []struct {
		config interfaces.CLIConfig
		errMsg string
	}{
		{
			config: interfaces.CLIConfig{ProjectName: "test"},
			errMsg: "at least one service",
		},
		{
			config: interfaces.CLIConfig{Backend: "django", ProjectName: "test"},
			errMsg: "unsupported backend technology",
		},
		{
			config: interfaces.CLIConfig{Frontend: "svelte", ProjectName: "test"},
			errMsg: "unsupported frontend technology",
		},
		{
			config: interfaces.CLIConfig{Database: "cassandra", ProjectName: "test"},
			errMsg: "unsupported database technology",
		},
		{
			config: interfaces.CLIConfig{Backend: "fastapi"},
			errMsg: "project name cannot be empty",
		},
	}

	for i, tt := range invalidConfigs {
		t.Run(fmt.Sprintf("invalid_config_%d", i), func(t *testing.T) {
			err := engine.ValidateComposeConfiguration(tt.config)
			if err == nil {
				t.Error("ValidateComposeConfiguration() should have failed for invalid config")
			}
			if !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("ValidateComposeConfiguration() error = %v, expected to contain %v", err, tt.errMsg)
			}
		})
	}
}

// TestDockerComposeServiceDependencies tests that service dependencies are correctly configured
func TestDockerComposeServiceDependencies(t *testing.T) {
	engine := NewEngine()

	config := interfaces.CLIConfig{
		Backend:     "express",
		Frontend:    "react",
		Database:    "postgres",
		ProjectName: "dependency-test",
	}

	result, err := engine.GenerateDockerCompose(config)
	if err != nil {
		t.Fatalf("GenerateDockerCompose() failed: %v", err)
	}

	resultStr := string(result)

	// Backend should depend on database
	if !strings.Contains(resultStr, "depends_on:") {
		t.Error("Generated Docker Compose missing service dependencies")
	}

	// Frontend should depend on backend
	frontendSection := extractServiceSection(resultStr, "frontend:")
	if !strings.Contains(frontendSection, "depends_on:") {
		t.Error("Frontend service should depend on backend")
	}

	// Backend should have database connection URL
	backendSection := extractServiceSection(resultStr, "backend:")
	if !strings.Contains(backendSection, "DATABASE_URL:") {
		t.Error("Backend service should have database connection URL")
	}

	// Frontend should have API URL
	if !strings.Contains(frontendSection, "REACT_APP_API_URL:") {
		t.Error("Frontend service should have backend API URL")
	}
}

// extractServiceSection extracts a service section from the Docker Compose YAML
func extractServiceSection(content, serviceName string) string {
	lines := strings.Split(content, "\n")
	var serviceLines []string
	inService := false
	
	for _, line := range lines {
		if strings.Contains(line, serviceName) && strings.HasSuffix(strings.TrimSpace(line), ":") {
			inService = true
			serviceLines = append(serviceLines, line)
			continue
		}
		
		if inService {
			// Check if we've reached the next service or section
			if strings.HasPrefix(line, "  ") && !strings.HasPrefix(line, "    ") && strings.HasSuffix(strings.TrimSpace(line), ":") {
				// This is another top-level service
				break
			}
			if strings.HasPrefix(line, "networks:") || strings.HasPrefix(line, "volumes:") {
				// We've reached the networks or volumes section
				break
			}
			serviceLines = append(serviceLines, line)
		}
	}
	
	return strings.Join(serviceLines, "\n")
}