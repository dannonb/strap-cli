package validation

import (
	"os"
	"path/filepath"
	"testing"

	"microservice-bootstrapper/internal/interfaces"
)

func TestValidator_ValidateConfig(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name    string
		config  interfaces.CLIConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid backend only",
			config: interfaces.CLIConfig{
				Backend: "fastapi",
			},
			wantErr: false,
		},
		{
			name: "valid frontend only",
			config: interfaces.CLIConfig{
				Frontend: "react",
			},
			wantErr: false,
		},
		{
			name: "valid database only",
			config: interfaces.CLIConfig{
				Database: "postgres",
			},
			wantErr: false,
		},
		{
			name: "valid full stack",
			config: interfaces.CLIConfig{
				Backend:  "gin",
				Frontend: "vue",
				Database: "mysql",
			},
			wantErr: false,
		},
		{
			name: "invalid backend",
			config: interfaces.CLIConfig{
				Backend: "invalid",
			},
			wantErr: true,
			errMsg:  "unsupported backend",
		},
		{
			name: "invalid frontend",
			config: interfaces.CLIConfig{
				Frontend: "invalid",
			},
			wantErr: true,
			errMsg:  "unsupported frontend",
		},
		{
			name: "invalid database",
			config: interfaces.CLIConfig{
				Database: "invalid",
			},
			wantErr: true,
			errMsg:  "unsupported database",
		},
		{
			name:    "no services specified",
			config:  interfaces.CLIConfig{},
			wantErr: true,
			errMsg:  "at least one service must be specified",
		},
		{
			name: "invalid project name - too short",
			config: interfaces.CLIConfig{
				Backend:     "fastapi",
				ProjectName: "a",
			},
			wantErr: true,
			errMsg:  "at least 2 characters",
		},
		{
			name: "invalid project name - special characters",
			config: interfaces.CLIConfig{
				Backend:     "fastapi",
				ProjectName: "my@project",
			},
			wantErr: true,
			errMsg:  "invalid project name",
		},
		{
			name: "reserved project name",
			config: interfaces.CLIConfig{
				Backend:     "fastapi",
				ProjectName: "docker",
			},
			wantErr: true,
			errMsg:  "reserved and cannot be used",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateConfig(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" {
				if err == nil || !containsString(err.Error(), tt.errMsg) {
					t.Errorf("ValidateConfig() error = %v, want error containing %v", err, tt.errMsg)
				}
			}
		})
	}
}

func TestValidator_ValidateDirectoryConflicts(t *testing.T) {
	validator := NewValidator()

	// Create a temporary directory for testing
	tempDir := t.TempDir()

	tests := []struct {
		name     string
		setup    func() string
		force    bool
		wantErr  bool
		errMsg   string
	}{
		{
			name: "non-existent directory",
			setup: func() string {
				return filepath.Join(tempDir, "non-existent")
			},
			force:   false,
			wantErr: false,
		},
		{
			name: "empty directory",
			setup: func() string {
				dir := filepath.Join(tempDir, "empty")
				os.MkdirAll(dir, 0755)
				return dir
			},
			force:   false,
			wantErr: false,
		},
		{
			name: "non-empty directory without force",
			setup: func() string {
				dir := filepath.Join(tempDir, "non-empty")
				os.MkdirAll(dir, 0755)
				// Create a file in the directory
				file := filepath.Join(dir, "test.txt")
				os.WriteFile(file, []byte("test"), 0644)
				return dir
			},
			force:   false,
			wantErr: true,
			errMsg:  "not empty",
		},
		{
			name: "non-empty directory with force",
			setup: func() string {
				dir := filepath.Join(tempDir, "non-empty-force")
				os.MkdirAll(dir, 0755)
				// Create a file in the directory
				file := filepath.Join(dir, "test.txt")
				os.WriteFile(file, []byte("test"), 0644)
				return dir
			},
			force:   true,
			wantErr: false,
		},
		{
			name: "file instead of directory",
			setup: func() string {
				file := filepath.Join(tempDir, "file.txt")
				os.WriteFile(file, []byte("test"), 0644)
				return file
			},
			force:   false,
			wantErr: true,
			errMsg:  "not a directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup()
			err := validator.ValidateDirectoryConflicts(path, tt.force)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDirectoryConflicts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if tt.wantErr && tt.errMsg != "" {
				if err == nil || !containsString(err.Error(), tt.errMsg) {
					t.Errorf("ValidateDirectoryConflicts() error = %v, want error containing %v", err, tt.errMsg)
				}
			}
		})
	}
}

func TestValidator_validateProjectName(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name        string
		projectName string
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "empty name (valid - uses default)",
			projectName: "",
			wantErr:     false,
		},
		{
			name:        "valid name with hyphens",
			projectName: "my-project",
			wantErr:     false,
		},
		{
			name:        "valid name with underscores",
			projectName: "my_project",
			wantErr:     false,
		},
		{
			name:        "valid alphanumeric name",
			projectName: "myproject123",
			wantErr:     false,
		},
		{
			name:        "invalid name with special characters",
			projectName: "my@project",
			wantErr:     true,
			errMsg:      "invalid project name",
		},
		{
			name:        "name too short",
			projectName: "a",
			wantErr:     true,
			errMsg:      "at least 2 characters",
		},
		{
			name:        "name too long",
			projectName: "this-is-a-very-long-project-name-that-exceeds-fifty-characters-limit",
			wantErr:     true,
			errMsg:      "50 characters or less",
		},
		{
			name:        "reserved name",
			projectName: "docker",
			wantErr:     true,
			errMsg:      "reserved",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.validateProjectName(tt.projectName)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateProjectName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" {
				if err == nil || !containsString(err.Error(), tt.errMsg) {
					t.Errorf("validateProjectName() error = %v, want error containing %v", err, tt.errMsg)
				}
			}
		})
	}
}

func TestValidator_validatePortConflicts(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name    string
		config  interfaces.CLIConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "no port conflicts - different services",
			config: interfaces.CLIConfig{
				Backend:  "fastapi", // port 8000
				Frontend: "react",   // port 3000
				Database: "postgres", // port 5432
			},
			wantErr: false,
		},
		{
			name: "potential port conflict - express backend and react frontend",
			config: interfaces.CLIConfig{
				Backend:  "express", // port 3000
				Frontend: "react",   // port 3000
			},
			wantErr: true,
			errMsg:  "port conflict",
		},
		{
			name: "single service - no conflict",
			config: interfaces.CLIConfig{
				Backend: "gin", // port 8080
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.validatePortConflicts(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("validatePortConflicts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" {
				if err == nil || !containsString(err.Error(), tt.errMsg) {
					t.Errorf("validatePortConflicts() error = %v, want error containing %v", err, tt.errMsg)
				}
			}
		})
	}
}

// Helper function to check if a string contains a substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || 
		containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}