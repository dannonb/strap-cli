package filesystem

import (
	"fmt"
	"path/filepath"
	"strings"

	"microservice-bootstrapper/internal/interfaces"
)

// MockManager is a mock implementation of FileSystemManager for testing
type MockManager struct {
	files       map[string][]byte // path -> content
	directories map[string]bool   // path -> exists
	
	// Control behavior
	CreateDirectoryError error
	WriteFileError       error
	CleanupError         error
	
	// Track calls
	CreateDirectoryCalls []string
	WriteFileCalls       []WriteFileCall
	CleanupCalls         []string
}

type WriteFileCall struct {
	Path    string
	Content []byte
}

// NewMockManager creates a new mock filesystem manager
func NewMockManager() *MockManager {
	return &MockManager{
		files:       make(map[string][]byte),
		directories: make(map[string]bool),
	}
}

// Verify interface compliance
var _ interfaces.FileSystemManager = (*MockManager)(nil)

// CreateDirectory mocks directory creation
func (m *MockManager) CreateDirectory(path string) error {
	m.CreateDirectoryCalls = append(m.CreateDirectoryCalls, path)
	
	if m.CreateDirectoryError != nil {
		return m.CreateDirectoryError
	}
	
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}
	
	// Mark directory as existing
	m.directories[path] = true
	
	// Also mark all parent directories as existing
	parent := filepath.Dir(path)
	for parent != "." && parent != "/" && parent != path {
		m.directories[parent] = true
		parent = filepath.Dir(parent)
	}
	
	return nil
}

// WriteFile mocks file writing
func (m *MockManager) WriteFile(path string, content []byte) error {
	m.WriteFileCalls = append(m.WriteFileCalls, WriteFileCall{
		Path:    path,
		Content: content,
	})
	
	if m.WriteFileError != nil {
		return m.WriteFileError
	}
	
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}
	
	// Create parent directory
	dir := filepath.Dir(path)
	if err := m.CreateDirectory(dir); err != nil {
		return err
	}
	
	// Store file content
	m.files[path] = content
	
	return nil
}

// FileExists mocks file existence check
func (m *MockManager) FileExists(path string) bool {
	if path == "" {
		return false
	}
	
	// Check if it's a file
	if _, exists := m.files[path]; exists {
		return true
	}
	
	// Check if it's a directory
	if exists, ok := m.directories[path]; ok && exists {
		return true
	}
	
	return false
}

// IsDirectoryEmpty mocks directory emptiness check
func (m *MockManager) IsDirectoryEmpty(path string) bool {
	if path == "" {
		return false
	}
	
	// Check if directory exists
	if exists, ok := m.directories[path]; !ok || !exists {
		return false
	}
	
	// Check if any files or subdirectories exist within this path
	for filePath := range m.files {
		if strings.HasPrefix(filePath, path+"/") || strings.HasPrefix(filePath, path+"\\") {
			return false
		}
	}
	
	for dirPath := range m.directories {
		if dirPath != path && (strings.HasPrefix(dirPath, path+"/") || strings.HasPrefix(dirPath, path+"\\")) {
			return false
		}
	}
	
	return true
}

// CleanupOnFailure mocks cleanup functionality
func (m *MockManager) CleanupOnFailure(basePath string) error {
	m.CleanupCalls = append(m.CleanupCalls, basePath)
	
	if m.CleanupError != nil {
		return m.CleanupError
	}
	
	if basePath == "" {
		return fmt.Errorf("basePath cannot be empty")
	}
	
	// Remove all files and directories within basePath
	for filePath := range m.files {
		if strings.HasPrefix(filePath, basePath) {
			delete(m.files, filePath)
		}
	}
	
	for dirPath := range m.directories {
		if strings.HasPrefix(dirPath, basePath) {
			delete(m.directories, dirPath)
		}
	}
	
	return nil
}

// Helper methods for testing

// GetFileContent returns the content of a mock file
func (m *MockManager) GetFileContent(path string) ([]byte, bool) {
	content, exists := m.files[path]
	return content, exists
}

// SetFileExists sets whether a file exists in the mock
func (m *MockManager) SetFileExists(path string, content []byte) {
	m.files[path] = content
}

// SetDirectoryExists sets whether a directory exists in the mock
func (m *MockManager) SetDirectoryExists(path string, exists bool) {
	m.directories[path] = exists
}

// Reset clears all mock state
func (m *MockManager) Reset() {
	m.files = make(map[string][]byte)
	m.directories = make(map[string]bool)
	m.CreateDirectoryError = nil
	m.WriteFileError = nil
	m.CleanupError = nil
	m.CreateDirectoryCalls = nil
	m.WriteFileCalls = nil
	m.CleanupCalls = nil
}