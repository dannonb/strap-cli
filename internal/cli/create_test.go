package cli

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestValidateFlags(t *testing.T) {
	tests := []struct {
		name     string
		backend  string
		frontend string
		database string
		project  string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "valid backend only",
			backend:  "fastapi",
			frontend: "",
			database: "",
			project:  "test-project",
			wantErr:  false,
		},
		{
			name:     "valid frontend only",
			backend:  "",
			frontend: "react",
			database: "",
			project:  "test-project",
			wantErr:  false,
		},
		{
			name:     "valid database only",
			backend:  "",
			frontend: "",
			database: "postgres",
			project:  "test-project",
			wantErr:  false,
		},
		{
			name:     "valid full stack",
			backend:  "gin",
			frontend: "vue",
			database: "mysql",
			project:  "test-project",
			wantErr:  false,
		},
		{
			name:     "invalid backend",
			backend:  "invalid",
			frontend: "",
			database: "",
			project:  "test-project",
			wantErr:  true,
			errMsg:   "Invalid backend",
		},
		{
			name:     "invalid frontend",
			backend:  "",
			frontend: "invalid",
			database: "",
			project:  "test-project",
			wantErr:  true,
			errMsg:   "Invalid frontend",
		},
		{
			name:     "invalid database",
			backend:  "",
			frontend: "",
			database: "invalid",
			project:  "test-project",
			wantErr:  true,
			errMsg:   "Invalid database",
		},
		{
			name:     "no services specified",
			backend:  "",
			frontend: "",
			database: "",
			project:  "test-project",
			wantErr:  true,
			errMsg:   "At least one service must be specified",
		},
		{
			name:     "project name with spaces",
			backend:  "fastapi",
			frontend: "",
			database: "",
			project:  "test project",
			wantErr:  true,
			errMsg:   "cannot contain spaces",
		},
		{
			name:     "project name starting with dash",
			backend:  "fastapi",
			frontend: "",
			database: "",
			project:  "-test-project",
			wantErr:  true,
			errMsg:   "cannot start with - or _",
		},
		{
			name:     "project name starting with underscore",
			backend:  "fastapi",
			frontend: "",
			database: "",
			project:  "_test-project",
			wantErr:  true,
			errMsg:   "cannot start with - or _",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set global variables
			backend = tt.backend
			frontend = tt.frontend
			database = tt.database
			projectName = tt.project

			// Create a mock command
			cmd := &cobra.Command{}
			
			err := validateFlags(cmd, []string{})
			
			if (err != nil) != tt.wantErr {
				t.Errorf("validateFlags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if tt.wantErr && tt.errMsg != "" {
				if err == nil || !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("validateFlags() error = %v, want error containing %v", err, tt.errMsg)
				}
			}
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name  string
		slice []string
		item  string
		want  bool
	}{
		{
			name:  "item exists",
			slice: []string{"fastapi", "express", "gin"},
			item:  "fastapi",
			want:  true,
		},
		{
			name:  "item does not exist",
			slice: []string{"fastapi", "express", "gin"},
			item:  "django",
			want:  false,
		},
		{
			name:  "empty slice",
			slice: []string{},
			item:  "fastapi",
			want:  false,
		},
		{
			name:  "empty item",
			slice: []string{"fastapi", "express", "gin"},
			item:  "",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := contains(tt.slice, tt.item)
			if got != tt.want {
				t.Errorf("contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandleGenerationError(t *testing.T) {
	tests := []struct {
		name     string
		inputErr error
		contains []string
	}{
		{
			name:     "validation error",
			inputErr: &mockError{msg: "validation failed: invalid backend"},
			contains: []string{"validation failed", "strap create --help", "strap examples"},
		},
		{
			name:     "prerequisite error",
			inputErr: &mockError{msg: "prerequisite check failed: docker not found"},
			contains: []string{"prerequisite check failed", "Docker", "Docker Compose"},
		},
		{
			name:     "directory not empty error",
			inputErr: &mockError{msg: "directory not empty"},
			contains: []string{"not empty", "--force", "mkdir myproject"},
		},
		{
			name:     "permission error",
			inputErr: &mockError{msg: "permission denied"},
			contains: []string{"permission denied", "write permissions", "administrator"},
		},
		{
			name:     "generic error",
			inputErr: &mockError{msg: "some other error"},
			contains: []string{"Project generation failed", "strap create --help", "docker --version"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handleGenerationError(tt.inputErr)
			
			if err == nil {
				t.Fatal("handleGenerationError() returned nil")
			}
			
			errStr := err.Error()
			for _, expected := range tt.contains {
				if !strings.Contains(errStr, expected) {
					t.Errorf("handleGenerationError() error does not contain expected string '%s'", expected)
				}
			}
		})
	}
}

func TestCreateCommandFlags(t *testing.T) {
	// Reset global variables
	backend = ""
	frontend = ""
	database = ""
	projectName = ""
	force = false

	// Test that flags are properly defined
	if createCmd.Flags().Lookup("be") == nil {
		t.Error("--be flag not defined")
	}
	if createCmd.Flags().Lookup("fe") == nil {
		t.Error("--fe flag not defined")
	}
	if createCmd.Flags().Lookup("db") == nil {
		t.Error("--db flag not defined")
	}
	if createCmd.Flags().Lookup("name") == nil {
		t.Error("--name flag not defined")
	}
	if createCmd.Flags().Lookup("force") == nil {
		t.Error("--force flag not defined")
	}
}

func TestCreateCommandHelp(t *testing.T) {
	// Test that create command help contains expected content
	helpText := createCmd.Long
	
	// Check that help contains expected sections
	expectedSections := []string{
		"Create a new microservice project",
		"backend",
		"frontend", 
		"database",
		"Docker",
		"boilerplate",
	}
	
	for _, section := range expectedSections {
		if !strings.Contains(helpText, section) {
			t.Errorf("Help text does not contain expected section: %s", section)
		}
	}
	
	// Test flag descriptions
	beFlag := createCmd.Flags().Lookup("be")
	if beFlag == nil {
		t.Error("--be flag not found")
	} else if !strings.Contains(beFlag.Usage, "fastapi") {
		t.Error("--be flag usage should mention fastapi")
	}
	
	feFlag := createCmd.Flags().Lookup("fe")
	if feFlag == nil {
		t.Error("--fe flag not found")
	} else if !strings.Contains(feFlag.Usage, "react") {
		t.Error("--fe flag usage should mention react")
	}
	
	dbFlag := createCmd.Flags().Lookup("db")
	if dbFlag == nil {
		t.Error("--db flag not found")
	} else if !strings.Contains(dbFlag.Usage, "postgres") {
		t.Error("--db flag usage should mention postgres")
	}
}

func TestSupportedTechnologies(t *testing.T) {
	// Test that supported technology slices contain expected values
	expectedBackends := []string{"fastapi", "express", "gin"}
	expectedFrontends := []string{"react", "vue", "angular"}
	expectedDatabases := []string{"mongo", "postgres", "mysql", "redis"}
	
	if len(supportedBackends) != len(expectedBackends) {
		t.Errorf("supportedBackends length = %d, want %d", len(supportedBackends), len(expectedBackends))
	}
	
	for _, backend := range expectedBackends {
		if !contains(supportedBackends, backend) {
			t.Errorf("supportedBackends does not contain %s", backend)
		}
	}
	
	if len(supportedFrontends) != len(expectedFrontends) {
		t.Errorf("supportedFrontends length = %d, want %d", len(supportedFrontends), len(expectedFrontends))
	}
	
	for _, frontend := range expectedFrontends {
		if !contains(supportedFrontends, frontend) {
			t.Errorf("supportedFrontends does not contain %s", frontend)
		}
	}
	
	if len(supportedDatabases) != len(expectedDatabases) {
		t.Errorf("supportedDatabases length = %d, want %d", len(supportedDatabases), len(expectedDatabases))
	}
	
	for _, database := range expectedDatabases {
		if !contains(supportedDatabases, database) {
			t.Errorf("supportedDatabases does not contain %s", database)
		}
	}
}

// mockError is a simple error implementation for testing
type mockError struct {
	msg string
}

func (e *mockError) Error() string {
	return e.msg
}