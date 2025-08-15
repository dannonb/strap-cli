package config

import (
	"testing"
)

func TestDirectoryInfo(t *testing.T) {
	tests := []struct {
		name     string
		dirInfo  DirectoryInfo
		expected DirectoryInfo
	}{
		{
			name: "basic directory info",
			dirInfo: DirectoryInfo{
				Path:         "/test/path",
				Name:         "test-project",
				IsEmpty:      true,
				HasConflicts: false,
				Permissions:  0755,
			},
			expected: DirectoryInfo{
				Path:         "/test/path",
				Name:         "test-project",
				IsEmpty:      true,
				HasConflicts: false,
				Permissions:  0755,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.dirInfo.Path != tt.expected.Path {
				t.Errorf("DirectoryInfo.Path = %v, want %v", tt.dirInfo.Path, tt.expected.Path)
			}
			if tt.dirInfo.Name != tt.expected.Name {
				t.Errorf("DirectoryInfo.Name = %v, want %v", tt.dirInfo.Name, tt.expected.Name)
			}
			if tt.dirInfo.IsEmpty != tt.expected.IsEmpty {
				t.Errorf("DirectoryInfo.IsEmpty = %v, want %v", tt.dirInfo.IsEmpty, tt.expected.IsEmpty)
			}
			if tt.dirInfo.HasConflicts != tt.expected.HasConflicts {
				t.Errorf("DirectoryInfo.HasConflicts = %v, want %v", tt.dirInfo.HasConflicts, tt.expected.HasConflicts)
			}
			if tt.dirInfo.Permissions != tt.expected.Permissions {
				t.Errorf("DirectoryInfo.Permissions = %v, want %v", tt.dirInfo.Permissions, tt.expected.Permissions)
			}
		})
	}
}

func TestProjectNameInfo(t *testing.T) {
	tests := []struct {
		name     string
		nameInfo ProjectNameInfo
		expected ProjectNameInfo
	}{
		{
			name: "inferred project name",
			nameInfo: ProjectNameInfo{
				Original:  "my-project-dir",
				Sanitized: "my-project-dir",
				Source:    "inferred",
			},
			expected: ProjectNameInfo{
				Original:  "my-project-dir",
				Sanitized: "my-project-dir",
				Source:    "inferred",
			},
		},
		{
			name: "provided project name",
			nameInfo: ProjectNameInfo{
				Original:  "MyProject",
				Sanitized: "myproject",
				Source:    "provided",
			},
			expected: ProjectNameInfo{
				Original:  "MyProject",
				Sanitized: "myproject",
				Source:    "provided",
			},
		},
		{
			name: "fallback project name",
			nameInfo: ProjectNameInfo{
				Original:  "",
				Sanitized: "microservice-project",
				Source:    "fallback",
			},
			expected: ProjectNameInfo{
				Original:  "",
				Sanitized: "microservice-project",
				Source:    "fallback",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.nameInfo.Original != tt.expected.Original {
				t.Errorf("ProjectNameInfo.Original = %v, want %v", tt.nameInfo.Original, tt.expected.Original)
			}
			if tt.nameInfo.Sanitized != tt.expected.Sanitized {
				t.Errorf("ProjectNameInfo.Sanitized = %v, want %v", tt.nameInfo.Sanitized, tt.expected.Sanitized)
			}
			if tt.nameInfo.Source != tt.expected.Source {
				t.Errorf("ProjectNameInfo.Source = %v, want %v", tt.nameInfo.Source, tt.expected.Source)
			}
		})
	}
}