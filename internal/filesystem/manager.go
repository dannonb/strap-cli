package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"microservice-bootstrapper/internal/interfaces"
	"microservice-bootstrapper/pkg/errors"
)

// Manager implements the FileSystemManager interface
type Manager struct {
	createdPaths []string // Track created paths for cleanup
}

// NewManager creates a new filesystem manager instance
func NewManager() interfaces.FileSystemManager {
	return &Manager{
		createdPaths: make([]string, 0),
	}
}

// CreateDirectory creates a directory and all necessary parent directories
func (m *Manager) CreateDirectory(path string) error {
	if path == "" {
		return errors.NewFileSystemError("create_directory", path, fmt.Errorf("path cannot be empty"))
	}

	// Check if directory already exists
	if info, err := os.Stat(path); err == nil {
		if !info.IsDir() {
			return errors.NewFileSystemError("create_directory", path, 
				fmt.Errorf("path exists but is not a directory"))
		}
		return nil // Directory already exists
	}

	// Create directory with proper permissions
	err := os.MkdirAll(path, 0755)
	if err != nil {
		// Provide more specific error messages
		if os.IsPermission(err) {
			return errors.NewFileSystemError("create_directory", path, 
				fmt.Errorf("permission denied - check write permissions for parent directory"))
		}
		if strings.Contains(err.Error(), "no space left") {
			return errors.NewFileSystemError("create_directory", path, 
				fmt.Errorf("insufficient disk space"))
		}
		return errors.NewFileSystemError("create_directory", path, err)
	}

	// Track created path for potential cleanup
	m.createdPaths = append(m.createdPaths, path)
	
	return nil
}

// WriteFile writes content to a file, creating parent directories if needed
func (m *Manager) WriteFile(path string, content []byte) error {
	if path == "" {
		return errors.NewFileSystemError("write_file", path, fmt.Errorf("path cannot be empty"))
	}

	if content == nil {
		return errors.NewFileSystemError("write_file", path, fmt.Errorf("content cannot be nil"))
	}

	// Create parent directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := m.CreateDirectory(dir); err != nil {
		return errors.NewFileSystemError("write_file", path, 
			fmt.Errorf("failed to create parent directory: %w", err))
	}

	// Check if file already exists and handle appropriately
	if info, err := os.Stat(path); err == nil {
		if info.IsDir() {
			return errors.NewFileSystemError("write_file", path, 
				fmt.Errorf("path exists but is a directory"))
		}
		// File exists, will be overwritten (this is expected behavior)
	}

	// Write file with proper permissions
	err := os.WriteFile(path, content, 0644)
	if err != nil {
		// Provide more specific error messages
		if os.IsPermission(err) {
			return errors.NewFileSystemError("write_file", path, 
				fmt.Errorf("permission denied - check write permissions"))
		}
		if strings.Contains(err.Error(), "no space left") {
			return errors.NewFileSystemError("write_file", path, 
				fmt.Errorf("insufficient disk space"))
		}
		if strings.Contains(err.Error(), "file name too long") {
			return errors.NewFileSystemError("write_file", path, 
				fmt.Errorf("file name too long"))
		}
		return errors.NewFileSystemError("write_file", path, err)
	}

	// Track created file for potential cleanup
	m.createdPaths = append(m.createdPaths, path)
	
	return nil
}

// FileExists checks if a file or directory exists
func (m *Manager) FileExists(path string) bool {
	if path == "" {
		return false
	}

	_, err := os.Stat(path)
	return err == nil
}

// IsDirectoryEmpty checks if a directory is empty
func (m *Manager) IsDirectoryEmpty(path string) bool {
	if path == "" {
		return false
	}

	// Check if path exists and is a directory
	info, err := os.Stat(path)
	if err != nil {
		return false // Path doesn't exist or can't be accessed
	}

	if !info.IsDir() {
		return false // Path exists but is not a directory
	}

	// Read directory contents
	entries, err := os.ReadDir(path)
	if err != nil {
		return false // Can't read directory
	}

	return len(entries) == 0
}

// CleanupOnFailure removes all files and directories created during this session
func (m *Manager) CleanupOnFailure(basePath string) error {
	if basePath == "" {
		return errors.NewCleanupError("", "basePath cannot be empty", nil)
	}

	var cleanupErrors []error
	var cleanedPaths []string

	fmt.Printf("🧹 Cleaning up partially created files...\n")

	// Remove all tracked paths in reverse order (files before directories)
	for i := len(m.createdPaths) - 1; i >= 0; i-- {
		path := m.createdPaths[i]
		
		// Only clean up paths that are within the base path
		if !filepath.HasPrefix(path, basePath) {
			continue
		}

		// Check if path still exists
		if !m.FileExists(path) {
			continue
		}

		// Remove file or directory
		err := os.RemoveAll(path)
		if err != nil {
			cleanupErrors = append(cleanupErrors, 
				errors.NewCleanupError(path, "failed to remove", err))
		} else {
			cleanedPaths = append(cleanedPaths, path)
		}
	}

	// Clear the tracked paths after cleanup attempt
	m.createdPaths = make([]string, 0)

	// Report cleanup results
	if len(cleanedPaths) > 0 {
		fmt.Printf("✅ Cleaned up %d files/directories\n", len(cleanedPaths))
	}

	// Return combined errors if any occurred
	if len(cleanupErrors) > 0 {
		var errorMessages []string
		for _, err := range cleanupErrors {
			errorMessages = append(errorMessages, err.Error())
		}
		return fmt.Errorf("cleanup partially failed:\n  - %s", 
			strings.Join(errorMessages, "\n  - "))
	}

	return nil
}

// GetCreatedPaths returns the list of paths created during this session
func (m *Manager) GetCreatedPaths() []string {
	return append([]string(nil), m.createdPaths...) // Return a copy
}

// ClearCreatedPaths clears the list of tracked created paths
func (m *Manager) ClearCreatedPaths() {
	m.createdPaths = make([]string, 0)
}