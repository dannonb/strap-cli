package errors

import (
	"errors"
	"strings"
	"testing"
)

func TestValidationError(t *testing.T) {
	// Test ValidationError creation and methods
	field := "backend"
	message := "unsupported technology"
	
	err := NewValidationError(field, message)
	
	if err.Field != field {
		t.Errorf("ValidationError.Field = %s, want %s", err.Field, field)
	}
	
	if err.Message != message {
		t.Errorf("ValidationError.Message = %s, want %s", err.Message, message)
	}
	
	expectedError := "validation error for field 'backend': unsupported technology"
	if err.Error() != expectedError {
		t.Errorf("ValidationError.Error() = %s, want %s", err.Error(), expectedError)
	}
}

func TestFileSystemError(t *testing.T) {
	// Test FileSystemError creation and methods
	operation := "create directory"
	path := "/tmp/test"
	cause := errors.New("permission denied")
	
	err := NewFileSystemError(operation, path, cause)
	
	if err.Operation != operation {
		t.Errorf("FileSystemError.Operation = %s, want %s", err.Operation, operation)
	}
	
	if err.Path != path {
		t.Errorf("FileSystemError.Path = %s, want %s", err.Path, path)
	}
	
	if err.Cause != cause {
		t.Errorf("FileSystemError.Cause = %v, want %v", err.Cause, cause)
	}
	
	expectedError := "filesystem error during create directory at path '/tmp/test': permission denied"
	if err.Error() != expectedError {
		t.Errorf("FileSystemError.Error() = %s, want %s", err.Error(), expectedError)
	}
	
	// Test Unwrap
	if err.Unwrap() != cause {
		t.Errorf("FileSystemError.Unwrap() = %v, want %v", err.Unwrap(), cause)
	}
}

func TestTemplateError(t *testing.T) {
	// Test TemplateError creation and methods
	template := "backend/fastapi/Dockerfile.tmpl"
	cause := errors.New("template not found")
	
	err := NewTemplateError(template, cause)
	
	if err.Template != template {
		t.Errorf("TemplateError.Template = %s, want %s", err.Template, template)
	}
	
	if err.Cause != cause {
		t.Errorf("TemplateError.Cause = %v, want %v", err.Cause, cause)
	}
	
	expectedError := "template error for 'backend/fastapi/Dockerfile.tmpl': template not found"
	if err.Error() != expectedError {
		t.Errorf("TemplateError.Error() = %s, want %s", err.Error(), expectedError)
	}
	
	// Test Unwrap
	if err.Unwrap() != cause {
		t.Errorf("TemplateError.Unwrap() = %v, want %v", err.Unwrap(), cause)
	}
}

func TestPrerequisiteError(t *testing.T) {
	// Test PrerequisiteError with cause
	component := "Docker"
	message := "Docker is not installed or not running"
	cause := errors.New("command not found")
	
	err := NewPrerequisiteError(component, message, cause)
	
	if err.Component != component {
		t.Errorf("PrerequisiteError.Component = %s, want %s", err.Component, component)
	}
	
	if err.Message != message {
		t.Errorf("PrerequisiteError.Message = %s, want %s", err.Message, message)
	}
	
	if err.Cause != cause {
		t.Errorf("PrerequisiteError.Cause = %v, want %v", err.Cause, cause)
	}
	
	expectedError := "prerequisite error for Docker: Docker is not installed or not running (cause: command not found)"
	if err.Error() != expectedError {
		t.Errorf("PrerequisiteError.Error() = %s, want %s", err.Error(), expectedError)
	}
	
	// Test Unwrap
	if err.Unwrap() != cause {
		t.Errorf("PrerequisiteError.Unwrap() = %v, want %v", err.Unwrap(), cause)
	}
}

func TestPrerequisiteErrorWithoutCause(t *testing.T) {
	// Test PrerequisiteError without cause
	component := "Docker Compose"
	message := "Docker Compose version is too old"
	
	err := NewPrerequisiteError(component, message, nil)
	
	if err.Component != component {
		t.Errorf("PrerequisiteError.Component = %s, want %s", err.Component, component)
	}
	
	if err.Message != message {
		t.Errorf("PrerequisiteError.Message = %s, want %s", err.Message, message)
	}
	
	if err.Cause != nil {
		t.Errorf("PrerequisiteError.Cause = %v, want nil", err.Cause)
	}
	
	expectedError := "prerequisite error for Docker Compose: Docker Compose version is too old"
	if err.Error() != expectedError {
		t.Errorf("PrerequisiteError.Error() = %s, want %s", err.Error(), expectedError)
	}
	
	// Test Unwrap returns nil
	if err.Unwrap() != nil {
		t.Errorf("PrerequisiteError.Unwrap() = %v, want nil", err.Unwrap())
	}
}

func TestDirectoryConflictError(t *testing.T) {
	// Test DirectoryConflictError with files
	path := "/tmp/project"
	message := "directory is not empty"
	files := []string{"file1.txt", "file2.txt", "directory1"}
	
	err := NewDirectoryConflictError(path, message, files)
	
	if err.Path != path {
		t.Errorf("DirectoryConflictError.Path = %s, want %s", err.Path, path)
	}
	
	if err.Message != message {
		t.Errorf("DirectoryConflictError.Message = %s, want %s", err.Message, message)
	}
	
	if len(err.Files) != len(files) {
		t.Errorf("DirectoryConflictError.Files length = %d, want %d", len(err.Files), len(files))
	}
	
	for i, file := range files {
		if err.Files[i] != file {
			t.Errorf("DirectoryConflictError.Files[%d] = %s, want %s", i, err.Files[i], file)
		}
	}
	
	errorStr := err.Error()
	expectedContains := []string{path, message, "file1.txt", "file2.txt", "directory1"}
	for _, expected := range expectedContains {
		if !strings.Contains(errorStr, expected) {
			t.Errorf("DirectoryConflictError.Error() should contain %s", expected)
		}
	}
}

func TestDirectoryConflictErrorWithManyFiles(t *testing.T) {
	// Test DirectoryConflictError with many files (should truncate)
	path := "/tmp/project"
	message := "directory is not empty"
	files := []string{"file1.txt", "file2.txt", "file3.txt", "file4.txt", "file5.txt", "file6.txt", "file7.txt"}
	
	err := NewDirectoryConflictError(path, message, files)
	
	errorStr := err.Error()
	
	// Should contain first 5 files and "..."
	expectedContains := []string{"file1.txt", "file2.txt", "file3.txt", "file4.txt", "file5.txt", "..."}
	for _, expected := range expectedContains {
		if !strings.Contains(errorStr, expected) {
			t.Errorf("DirectoryConflictError.Error() should contain %s", expected)
		}
	}
	
	// Should NOT contain files beyond the first 5
	if strings.Contains(errorStr, "file6.txt") {
		t.Error("DirectoryConflictError.Error() should not contain file6.txt (should be truncated)")
	}
}

func TestDirectoryConflictErrorWithoutFiles(t *testing.T) {
	// Test DirectoryConflictError without files
	path := "/tmp/project"
	message := "directory is not empty"
	
	err := NewDirectoryConflictError(path, message, nil)
	
	expectedError := "directory conflict at '/tmp/project': directory is not empty"
	if err.Error() != expectedError {
		t.Errorf("DirectoryConflictError.Error() = %s, want %s", err.Error(), expectedError)
	}
}

func TestCleanupError(t *testing.T) {
	// Test CleanupError creation and methods
	path := "/tmp/cleanup"
	message := "failed to remove directory"
	cause := errors.New("directory not empty")
	
	err := NewCleanupError(path, message, cause)
	
	if err.Path != path {
		t.Errorf("CleanupError.Path = %s, want %s", err.Path, path)
	}
	
	if err.Message != message {
		t.Errorf("CleanupError.Message = %s, want %s", err.Message, message)
	}
	
	if err.Cause != cause {
		t.Errorf("CleanupError.Cause = %v, want %v", err.Cause, cause)
	}
	
	expectedError := "cleanup error at '/tmp/cleanup': failed to remove directory (cause: directory not empty)"
	if err.Error() != expectedError {
		t.Errorf("CleanupError.Error() = %s, want %s", err.Error(), expectedError)
	}
	
	// Test Unwrap
	if err.Unwrap() != cause {
		t.Errorf("CleanupError.Unwrap() = %v, want %v", err.Unwrap(), cause)
	}
}

func TestErrorsImplementErrorInterface(t *testing.T) {
	// Test that all custom errors implement the error interface
	var err error
	
	err = NewValidationError("field", "message")
	if err.Error() == "" {
		t.Error("ValidationError should implement error interface")
	}
	
	err = NewFileSystemError("op", "path", errors.New("cause"))
	if err.Error() == "" {
		t.Error("FileSystemError should implement error interface")
	}
	
	err = NewTemplateError("template", errors.New("cause"))
	if err.Error() == "" {
		t.Error("TemplateError should implement error interface")
	}
	
	err = NewPrerequisiteError("component", "message", errors.New("cause"))
	if err.Error() == "" {
		t.Error("PrerequisiteError should implement error interface")
	}
	
	err = NewDirectoryConflictError("path", "message", []string{"file"})
	if err.Error() == "" {
		t.Error("DirectoryConflictError should implement error interface")
	}
	
	err = NewCleanupError("path", "message", errors.New("cause"))
	if err.Error() == "" {
		t.Error("CleanupError should implement error interface")
	}
}

func TestErrorUnwrapping(t *testing.T) {
	// Test that errors that should support unwrapping do so correctly
	cause := errors.New("root cause")
	
	// Test FileSystemError unwrapping
	fsErr := NewFileSystemError("op", "path", cause)
	if !errors.Is(fsErr, cause) {
		t.Error("FileSystemError should support error unwrapping")
	}
	
	// Test TemplateError unwrapping
	tmplErr := NewTemplateError("template", cause)
	if !errors.Is(tmplErr, cause) {
		t.Error("TemplateError should support error unwrapping")
	}
	
	// Test PrerequisiteError unwrapping
	prereqErr := NewPrerequisiteError("component", "message", cause)
	if !errors.Is(prereqErr, cause) {
		t.Error("PrerequisiteError should support error unwrapping")
	}
	
	// Test CleanupError unwrapping
	cleanupErr := NewCleanupError("path", "message", cause)
	if !errors.Is(cleanupErr, cause) {
		t.Error("CleanupError should support error unwrapping")
	}
}