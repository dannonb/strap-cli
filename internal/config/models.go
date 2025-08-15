package config

import "os"

// DirectoryInfo contains information about a directory used for project creation
type DirectoryInfo struct {
	Path         string      // Full path to the directory
	Name         string      // Directory name
	IsEmpty      bool        // Whether the directory is empty
	HasConflicts bool        // Whether there are potential file conflicts
	Permissions  os.FileMode // Directory permissions
}

// ProjectNameInfo contains information about project name resolution
type ProjectNameInfo struct {
	Original   string // Original name (directory name or provided name)
	Sanitized  string // Sanitized version suitable for project naming
	Source     string // Source of the name: "provided", "inferred", "fallback"
}