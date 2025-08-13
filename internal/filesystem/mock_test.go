package filesystem

import (
	"fmt"
	"testing"
)

func TestMockManager_CreateDirectory(t *testing.T) {
	mock := NewMockManager()

	// Test successful directory creation
	err := mock.CreateDirectory("test/dir")
	if err != nil {
		t.Errorf("CreateDirectory() error = %v", err)
	}

	// Verify directory exists
	if !mock.FileExists("test/dir") {
		t.Error("Directory should exist after creation")
	}

	// Verify parent directories were created
	if !mock.FileExists("test") {
		t.Error("Parent directory should exist after creation")
	}

	// Verify call was tracked
	if len(mock.CreateDirectoryCalls) != 1 || mock.CreateDirectoryCalls[0] != "test/dir" {
		t.Errorf("CreateDirectoryCalls = %v, want [test/dir]", mock.CreateDirectoryCalls)
	}

	// Test error behavior
	mock.CreateDirectoryError = fmt.Errorf("mock error")
	err = mock.CreateDirectory("error/dir")
	if err == nil {
		t.Error("Expected error from CreateDirectory")
	}
}

func TestMockManager_WriteFile(t *testing.T) {
	mock := NewMockManager()

	content := []byte("test content")
	err := mock.WriteFile("test/file.txt", content)
	if err != nil {
		t.Errorf("WriteFile() error = %v", err)
	}

	// Verify file exists
	if !mock.FileExists("test/file.txt") {
		t.Error("File should exist after writing")
	}

	// Verify content
	storedContent, exists := mock.GetFileContent("test/file.txt")
	if !exists {
		t.Error("File content should be stored")
	}
	if string(storedContent) != string(content) {
		t.Errorf("File content = %s, want %s", string(storedContent), string(content))
	}

	// Verify call was tracked
	if len(mock.WriteFileCalls) != 1 {
		t.Errorf("WriteFileCalls length = %d, want 1", len(mock.WriteFileCalls))
	}
	if mock.WriteFileCalls[0].Path != "test/file.txt" {
		t.Errorf("WriteFileCalls[0].Path = %s, want test/file.txt", mock.WriteFileCalls[0].Path)
	}

	// Test error behavior
	mock.WriteFileError = fmt.Errorf("mock error")
	err = mock.WriteFile("error/file.txt", content)
	if err == nil {
		t.Error("Expected error from WriteFile")
	}
}

func TestMockManager_FileExists(t *testing.T) {
	mock := NewMockManager()

	// Test non-existing file
	if mock.FileExists("nonexistent.txt") {
		t.Error("Non-existing file should not exist")
	}

	// Create a file and test
	mock.SetFileExists("test.txt", []byte("content"))
	if !mock.FileExists("test.txt") {
		t.Error("Existing file should exist")
	}

	// Create a directory and test
	mock.SetDirectoryExists("testdir", true)
	if !mock.FileExists("testdir") {
		t.Error("Existing directory should exist")
	}

	// Test empty path
	if mock.FileExists("") {
		t.Error("Empty path should not exist")
	}
}

func TestMockManager_IsDirectoryEmpty(t *testing.T) {
	mock := NewMockManager()

	// Test non-existing directory
	if mock.IsDirectoryEmpty("nonexistent") {
		t.Error("Non-existing directory should not be considered empty")
	}

	// Create empty directory
	mock.SetDirectoryExists("empty", true)
	if !mock.IsDirectoryEmpty("empty") {
		t.Error("Empty directory should be empty")
	}

	// Add file to directory
	mock.SetFileExists("empty/file.txt", []byte("content"))
	if mock.IsDirectoryEmpty("empty") {
		t.Error("Directory with file should not be empty")
	}

	// Test with subdirectory
	mock.Reset()
	mock.SetDirectoryExists("parent", true)
	mock.SetDirectoryExists("parent/child", true)
	if mock.IsDirectoryEmpty("parent") {
		t.Error("Directory with subdirectory should not be empty")
	}

	// Test empty path
	if mock.IsDirectoryEmpty("") {
		t.Error("Empty path should not be considered empty")
	}
}

func TestMockManager_CleanupOnFailure(t *testing.T) {
	mock := NewMockManager()

	// Set up some files and directories
	mock.SetFileExists("project/file1.txt", []byte("content1"))
	mock.SetFileExists("project/subdir/file2.txt", []byte("content2"))
	mock.SetDirectoryExists("project", true)
	mock.SetDirectoryExists("project/subdir", true)
	mock.SetFileExists("outside/file3.txt", []byte("content3"))

	// Cleanup project directory
	err := mock.CleanupOnFailure("project")
	if err != nil {
		t.Errorf("CleanupOnFailure() error = %v", err)
	}

	// Verify files within project were cleaned up
	if mock.FileExists("project/file1.txt") {
		t.Error("File within project should be cleaned up")
	}
	if mock.FileExists("project/subdir/file2.txt") {
		t.Error("File within project subdirectory should be cleaned up")
	}
	if mock.FileExists("project") {
		t.Error("Project directory should be cleaned up")
	}

	// Verify files outside project were not cleaned up
	if !mock.FileExists("outside/file3.txt") {
		t.Error("File outside project should not be cleaned up")
	}

	// Verify call was tracked
	if len(mock.CleanupCalls) != 1 || mock.CleanupCalls[0] != "project" {
		t.Errorf("CleanupCalls = %v, want [project]", mock.CleanupCalls)
	}

	// Test error behavior
	mock.CleanupError = fmt.Errorf("mock error")
	err = mock.CleanupOnFailure("test")
	if err == nil {
		t.Error("Expected error from CleanupOnFailure")
	}
}

func TestMockManager_Reset(t *testing.T) {
	mock := NewMockManager()

	// Set up some state
	mock.SetFileExists("file.txt", []byte("content"))
	mock.SetDirectoryExists("dir", true)
	mock.CreateDirectoryError = fmt.Errorf("error")
	mock.WriteFileError = fmt.Errorf("error")
	mock.CleanupError = fmt.Errorf("error")

	// Make some calls to track
	mock.CreateDirectory("test")
	mock.WriteFile("test.txt", []byte("content"))
	mock.CleanupOnFailure("test")

	// Reset
	mock.Reset()

	// Verify everything was cleared
	if mock.FileExists("file.txt") {
		t.Error("File should not exist after reset")
	}
	if mock.FileExists("dir") {
		t.Error("Directory should not exist after reset")
	}
	if mock.CreateDirectoryError != nil {
		t.Error("CreateDirectoryError should be nil after reset")
	}
	if mock.WriteFileError != nil {
		t.Error("WriteFileError should be nil after reset")
	}
	if mock.CleanupError != nil {
		t.Error("CleanupError should be nil after reset")
	}
	if len(mock.CreateDirectoryCalls) != 0 {
		t.Error("CreateDirectoryCalls should be empty after reset")
	}
	if len(mock.WriteFileCalls) != 0 {
		t.Error("WriteFileCalls should be empty after reset")
	}
	if len(mock.CleanupCalls) != 0 {
		t.Error("CleanupCalls should be empty after reset")
	}
}