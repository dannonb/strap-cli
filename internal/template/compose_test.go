package template

import (
	"strings"
	"testing"

	"microservice-bootstrapper/internal/interfaces"
)

func TestComposeGenerator_GenerateDockerCompose(t *testing.T) {
	engine := NewEngine()
	generator := NewComposeGenerator(engine)

	tests := []struct {
		name     string
		config   interfaces.CLIConfig
		wantErr  bool
		contains []string
	}{
		{
			name: "backend only",
			config: interfaces.CLIConfig{
				Backend:     "fastapi",
				ProjectName: "test-project",
			},
			wantErr: false,
			contains: []string{
				"backend:",
				"test-project-backend",
				"8000:8000",
				"UVICORN_HOST: 0.0.0.0",
				"test-project-network",
			},
		},
		{
			name: "frontend only",
			config: interfaces.CLIConfig{
				Frontend:    "react",
				ProjectName: "test-project",
			},
			wantErr: false,
			contains: []string{
				"frontend:",
				"test-project-frontend",
				"3000:3000",
				"REACT_APP_ENV: development",
				"test-project-network",
			},
		},
		{
			name: "database only",
			config: interfaces.CLIConfig{
				Database:    "postgres",
				ProjectName: "test-project",
			},
			wantErr: false,
			contains: []string{
				"postgres:",
				"test-project-postgres",
				"5432:5432",
				"POSTGRES_DB: test-project",
				"test-project-postgres-data",
			},
		},
		{
			name: "full stack",
			config: interfaces.CLIConfig{
				Backend:     "express",
				Frontend:    "vue",
				Database:    "mongo",
				ProjectName: "full-stack-app",
			},
			wantErr: false,
			contains: []string{
				"backend:",
				"frontend:",
				"mongo:",
				"full-stack-app-backend",
				"full-stack-app-frontend",
				"full-stack-app-mongo",
				"depends_on:",
				"VUE_APP_API_URL: http://localhost:3000",
				"MONGODB_URL:",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := generator.GenerateDockerCompose(tt.config)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateDockerCompose() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				resultStr := string(result)
				for _, expected := range tt.contains {
					if !strings.Contains(resultStr, expected) {
						t.Errorf("GenerateDockerCompose() result does not contain expected string: %s", expected)
						t.Logf("Generated content:\n%s", resultStr)
					}
				}
			}
		})
	}
}

func TestComposeGenerator_ValidateConfiguration(t *testing.T) {
	engine := NewEngine()
	generator := NewComposeGenerator(engine)

	tests := []struct {
		name    string
		config  interfaces.CLIConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid backend only",
			config: interfaces.CLIConfig{
				Backend:     "fastapi",
				ProjectName: "test-project",
			},
			wantErr: false,
		},
		{
			name: "valid full stack",
			config: interfaces.CLIConfig{
				Backend:     "gin",
				Frontend:    "angular",
				Database:    "redis",
				ProjectName: "test-project",
			},
			wantErr: false,
		},
		{
			name: "no services specified",
			config: interfaces.CLIConfig{
				ProjectName: "test-project",
			},
			wantErr: true,
			errMsg:  "at least one service",
		},
		{
			name: "invalid backend technology",
			config: interfaces.CLIConfig{
				Backend:     "django",
				ProjectName: "test-project",
			},
			wantErr: true,
			errMsg:  "unsupported backend technology",
		},
		{
			name: "invalid frontend technology",
			config: interfaces.CLIConfig{
				Frontend:    "svelte",
				ProjectName: "test-project",
			},
			wantErr: true,
			errMsg:  "unsupported frontend technology",
		},
		{
			name: "invalid database technology",
			config: interfaces.CLIConfig{
				Database:    "cassandra",
				ProjectName: "test-project",
			},
			wantErr: true,
			errMsg:  "unsupported database technology",
		},
		{
			name: "empty project name",
			config: interfaces.CLIConfig{
				Backend: "fastapi",
			},
			wantErr: true,
			errMsg:  "project name cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := generator.ValidateConfiguration(tt.config)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateConfiguration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("ValidateConfiguration() error = %v, expected to contain %v", err, tt.errMsg)
			}
		})
	}
}

func TestComposeGenerator_buildPortConfig(t *testing.T) {
	engine := NewEngine()
	generator := NewComposeGenerator(engine)

	tests := []struct {
		name     string
		config   interfaces.CLIConfig
		expected interfaces.PortConfig
	}{
		{
			name: "fastapi backend",
			config: interfaces.CLIConfig{
				Backend: "fastapi",
			},
			expected: interfaces.PortConfig{
				Backend: 8000,
			},
		},
		{
			name: "express backend with react frontend",
			config: interfaces.CLIConfig{
				Backend:  "express",
				Frontend: "react",
			},
			expected: interfaces.PortConfig{
				Backend:  3000,
				Frontend: 3000,
			},
		},
		{
			name: "full stack with postgres",
			config: interfaces.CLIConfig{
				Backend:  "gin",
				Frontend: "angular",
				Database: "postgres",
			},
			expected: interfaces.PortConfig{
				Backend:  8080,
				Frontend: 4200,
				Database: 5432,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generator.buildPortConfig(tt.config)
			
			if result.Backend != tt.expected.Backend {
				t.Errorf("buildPortConfig() Backend port = %v, expected %v", result.Backend, tt.expected.Backend)
			}
			if result.Frontend != tt.expected.Frontend {
				t.Errorf("buildPortConfig() Frontend port = %v, expected %v", result.Frontend, tt.expected.Frontend)
			}
			if result.Database != tt.expected.Database {
				t.Errorf("buildPortConfig() Database port = %v, expected %v", result.Database, tt.expected.Database)
			}
		})
	}
}

func TestComposeGenerator_sanitizeProjectName(t *testing.T) {
	engine := NewEngine()
	generator := NewComposeGenerator(engine)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple name",
			input:    "myproject",
			expected: "myproject",
		},
		{
			name:     "name with spaces",
			input:    "my project",
			expected: "my-project",
		},
		{
			name:     "name with underscores",
			input:    "my_project_name",
			expected: "my-project-name",
		},
		{
			name:     "mixed case with special chars",
			input:    "My_Project Name!",
			expected: "my-project-name",
		},
		{
			name:     "numbers and letters",
			input:    "project123",
			expected: "project123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generator.sanitizeProjectName(tt.input)
			if result != tt.expected {
				t.Errorf("sanitizeProjectName() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestComposeGenerator_getDatabaseEnvironment(t *testing.T) {
	engine := NewEngine()
	generator := NewComposeGenerator(engine)

	tests := []struct {
		name        string
		dbType      string
		projectName string
		expected    map[string]string
	}{
		{
			name:        "postgres",
			dbType:      "postgres",
			projectName: "test-project",
			expected: map[string]string{
				"POSTGRES_DB":       "test-project",
				"POSTGRES_USER":     "postgres",
				"POSTGRES_PASSWORD": "postgres",
			},
		},
		{
			name:        "mongo",
			dbType:      "mongo",
			projectName: "test-project",
			expected: map[string]string{
				"MONGO_INITDB_ROOT_USERNAME": "admin",
				"MONGO_INITDB_ROOT_PASSWORD": "admin123",
				"MONGO_INITDB_DATABASE":      "test-project",
			},
		},
		{
			name:        "mysql",
			dbType:      "mysql",
			projectName: "test-project",
			expected: map[string]string{
				"MYSQL_ROOT_PASSWORD": "rootpassword",
				"MYSQL_DATABASE":      "test-project",
				"MYSQL_USER":          "user",
				"MYSQL_PASSWORD":      "password",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generator.getDatabaseEnvironment(tt.dbType, tt.projectName)
			
			for key, expectedValue := range tt.expected {
				if result[key] != expectedValue {
					t.Errorf("getDatabaseEnvironment() %s = %v, expected %v", key, result[key], expectedValue)
				}
			}
		})
	}
}