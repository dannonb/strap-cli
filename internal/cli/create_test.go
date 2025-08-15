package cli

import (
	"testing"

	"microservice-bootstrapper/internal/interfaces"
)

func TestResolveProjectNameWithDetailsIntegration(t *testing.T) {
	tests := []struct {
		name         string
		providedName string
		expectError  bool
	}{
		{
			name:         "explicit name provided",
			providedName: "my-project",
			expectError:  false,
		},
		{
			name:         "empty name - should infer from directory",
			providedName: "",
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolvedName, originalDir, err := resolveProjectNameWithDetails(tt.providedName)
			
			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			
			if tt.providedName != "" {
				// When name is provided, it should be used as-is
				if resolvedName != tt.providedName {
					t.Errorf("Expected resolved name %s, got %s", tt.providedName, resolvedName)
				}
				if originalDir != "" {
					t.Errorf("Expected empty original directory when name is provided, got %s", originalDir)
				}
			} else {
				// When name is empty, it should be inferred from directory
				if resolvedName == "" {
					t.Errorf("Expected non-empty resolved name when inferring from directory")
				}
				if originalDir == "" {
					t.Errorf("Expected non-empty original directory when inferring name")
				}
			}
		})
	}
}

func TestCLIConfigCreation(t *testing.T) {
	tests := []struct {
		name         string
		providedName string
		backend      string
		frontend     string
		database     string
		force        bool
	}{
		{
			name:         "explicit name with backend",
			providedName: "my-api",
			backend:      "fastapi",
			force:        false,
		},
		{
			name:         "inferred name with full stack",
			providedName: "",
			backend:      "express",
			frontend:     "react",
			database:     "postgres",
			force:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the CLI config creation logic
			resolvedName, originalDir, err := resolveProjectNameWithDetails(tt.providedName)
			if err != nil {
				t.Fatalf("Failed to resolve project name: %v", err)
			}

			config := interfaces.CLIConfig{
				Backend:      tt.backend,
				Frontend:     tt.frontend,
				Database:     tt.database,
				ProjectName:  resolvedName,
				Force:        tt.force,
				InferredName: tt.providedName == "",
				WorkingDir:   originalDir,
			}

			// Verify the config fields
			if config.Backend != tt.backend {
				t.Errorf("Expected backend %s, got %s", tt.backend, config.Backend)
			}
			if config.Frontend != tt.frontend {
				t.Errorf("Expected frontend %s, got %s", tt.frontend, config.Frontend)
			}
			if config.Database != tt.database {
				t.Errorf("Expected database %s, got %s", tt.database, config.Database)
			}
			if config.Force != tt.force {
				t.Errorf("Expected force %v, got %v", tt.force, config.Force)
			}

			// Verify inference logic
			expectedInferred := tt.providedName == ""
			if config.InferredName != expectedInferred {
				t.Errorf("Expected InferredName %v, got %v", expectedInferred, config.InferredName)
			}

			// Verify working directory logic
			if tt.providedName == "" {
				// When name is inferred, WorkingDir should be set
				if config.WorkingDir == "" {
					t.Errorf("Expected non-empty WorkingDir when name is inferred")
				}
			} else {
				// When name is provided, WorkingDir should be empty
				if config.WorkingDir != "" {
					t.Errorf("Expected empty WorkingDir when name is provided, got %s", config.WorkingDir)
				}
			}

			// Verify project name is set
			if config.ProjectName == "" {
				t.Errorf("Expected non-empty ProjectName")
			}
		})
	}
}