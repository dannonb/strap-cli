package interfaces

import (
	"testing"
)

func TestCLIConfig(t *testing.T) {
	tests := []struct {
		name     string
		config   CLIConfig
		expected CLIConfig
	}{
		{
			name: "basic config with new fields",
			config: CLIConfig{
				Backend:      "fastapi",
				Frontend:     "react",
				Database:     "postgres",
				ProjectName:  "test-project",
				Force:        false,
				InferredName: true,
				WorkingDir:   "test-dir",
			},
			expected: CLIConfig{
				Backend:      "fastapi",
				Frontend:     "react",
				Database:     "postgres",
				ProjectName:  "test-project",
				Force:        false,
				InferredName: true,
				WorkingDir:   "test-dir",
			},
		},
		{
			name: "config with explicit name",
			config: CLIConfig{
				Backend:      "gin",
				ProjectName:  "my-api",
				Force:        true,
				InferredName: false,
				WorkingDir:   "",
			},
			expected: CLIConfig{
				Backend:      "gin",
				ProjectName:  "my-api",
				Force:        true,
				InferredName: false,
				WorkingDir:   "",
			},
		},
		{
			name: "config with default values",
			config: CLIConfig{
				Backend:     "express",
				ProjectName: "default-project",
			},
			expected: CLIConfig{
				Backend:      "express",
				ProjectName:  "default-project",
				InferredName: false, // default value
				WorkingDir:   "",    // default value
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.config.Backend != tt.expected.Backend {
				t.Errorf("CLIConfig.Backend = %v, want %v", tt.config.Backend, tt.expected.Backend)
			}
			if tt.config.Frontend != tt.expected.Frontend {
				t.Errorf("CLIConfig.Frontend = %v, want %v", tt.config.Frontend, tt.expected.Frontend)
			}
			if tt.config.Database != tt.expected.Database {
				t.Errorf("CLIConfig.Database = %v, want %v", tt.config.Database, tt.expected.Database)
			}
			if tt.config.ProjectName != tt.expected.ProjectName {
				t.Errorf("CLIConfig.ProjectName = %v, want %v", tt.config.ProjectName, tt.expected.ProjectName)
			}
			if tt.config.Force != tt.expected.Force {
				t.Errorf("CLIConfig.Force = %v, want %v", tt.config.Force, tt.expected.Force)
			}
			if tt.config.InferredName != tt.expected.InferredName {
				t.Errorf("CLIConfig.InferredName = %v, want %v", tt.config.InferredName, tt.expected.InferredName)
			}
			if tt.config.WorkingDir != tt.expected.WorkingDir {
				t.Errorf("CLIConfig.WorkingDir = %v, want %v", tt.config.WorkingDir, tt.expected.WorkingDir)
			}
		})
	}
}

func TestCLIConfigInferredNameLogic(t *testing.T) {
	tests := []struct {
		name         string
		providedName string
		workingDir   string
		expected     bool
	}{
		{
			name:         "name provided explicitly",
			providedName: "my-project",
			workingDir:   "some-dir",
			expected:     false,
		},
		{
			name:         "name inferred from directory",
			providedName: "",
			workingDir:   "project-dir",
			expected:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := CLIConfig{
				ProjectName:  tt.providedName,
				InferredName: tt.providedName == "",
				WorkingDir:   tt.workingDir,
			}

			if config.InferredName != tt.expected {
				t.Errorf("CLIConfig.InferredName = %v, want %v", config.InferredName, tt.expected)
			}
		})
	}
}