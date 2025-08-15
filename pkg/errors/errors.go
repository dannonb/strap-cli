package errors

import (
	"fmt"
	"strings"
)

// ValidationError represents an input validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error for field '%s': %s", e.Field, e.Message)
}

// NewValidationError creates a new ValidationError
func NewValidationError(field, message string) ValidationError {
	return ValidationError{
		Field:   field,
		Message: message,
	}
}

// FileSystemError represents a file system operation error
type FileSystemError struct {
	Operation string
	Path      string
	Cause     error
}

func (e FileSystemError) Error() string {
	return fmt.Sprintf("filesystem error during %s at path '%s': %v", e.Operation, e.Path, e.Cause)
}

// Unwrap returns the underlying error
func (e FileSystemError) Unwrap() error {
	return e.Cause
}

// NewFileSystemError creates a new FileSystemError
func NewFileSystemError(operation, path string, cause error) FileSystemError {
	return FileSystemError{
		Operation: operation,
		Path:      path,
		Cause:     cause,
	}
}

// TemplateError represents a template processing error
type TemplateError struct {
	Template string
	Cause    error
}

func (e TemplateError) Error() string {
	return fmt.Sprintf("template error for '%s': %v", e.Template, e.Cause)
}

// Unwrap returns the underlying error
func (e TemplateError) Unwrap() error {
	return e.Cause
}

// NewTemplateError creates a new TemplateError
func NewTemplateError(template string, cause error) TemplateError {
	return TemplateError{
		Template: template,
		Cause:    cause,
	}
}

// PrerequisiteError represents a system prerequisite error
type PrerequisiteError struct {
	Component string
	Message   string
	Cause     error
}

func (e PrerequisiteError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("prerequisite error for %s: %s (cause: %v)", e.Component, e.Message, e.Cause)
	}
	return fmt.Sprintf("prerequisite error for %s: %s", e.Component, e.Message)
}

// Unwrap returns the underlying error
func (e PrerequisiteError) Unwrap() error {
	return e.Cause
}

// NewPrerequisiteError creates a new PrerequisiteError
func NewPrerequisiteError(component, message string, cause error) PrerequisiteError {
	return PrerequisiteError{
		Component: component,
		Message:   message,
		Cause:     cause,
	}
}

// DirectoryConflictError represents a directory conflict error
type DirectoryConflictError struct {
	Path    string
	Message string
	Files   []string
}

func (e DirectoryConflictError) Error() string {
	if len(e.Files) > 0 {
		fileList := strings.Join(e.Files, ", ")
		if len(e.Files) > 5 {
			fileList = strings.Join(e.Files[:5], ", ") + "..."
		}
		return fmt.Sprintf("directory conflict at '%s': %s (contains: %s)", e.Path, e.Message, fileList)
	}
	return fmt.Sprintf("directory conflict at '%s': %s", e.Path, e.Message)
}

// NewDirectoryConflictError creates a new DirectoryConflictError
func NewDirectoryConflictError(path, message string, files []string) DirectoryConflictError {
	return DirectoryConflictError{
		Path:    path,
		Message: message,
		Files:   files,
	}
}

// CleanupError represents an error during cleanup operations
type CleanupError struct {
	Path    string
	Message string
	Cause   error
}

func (e CleanupError) Error() string {
	return fmt.Sprintf("cleanup error at '%s': %s (cause: %v)", e.Path, e.Message, e.Cause)
}

// Unwrap returns the underlying error
func (e CleanupError) Unwrap() error {
	return e.Cause
}

// NewCleanupError creates a new CleanupError
func NewCleanupError(path, message string, cause error) CleanupError {
	return CleanupError{
		Path:    path,
		Message: message,
		Cause:   cause,
	}
}

// DirectoryInferenceError represents an error during directory name inference
type DirectoryInferenceError struct {
	Path    string
	Reason  string
	Cause   error
}

func (e DirectoryInferenceError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("directory name inference failed for '%s': %s (cause: %v)", e.Path, e.Reason, e.Cause)
	}
	return fmt.Sprintf("directory name inference failed for '%s': %s", e.Path, e.Reason)
}

// Unwrap returns the underlying error
func (e DirectoryInferenceError) Unwrap() error {
	return e.Cause
}

// NewDirectoryInferenceError creates a new DirectoryInferenceError
func NewDirectoryInferenceError(path, reason string, cause error) DirectoryInferenceError {
	return DirectoryInferenceError{
		Path:   path,
		Reason: reason,
		Cause:  cause,
	}
}

// DockerError represents a Docker-related error with user-friendly guidance
type DockerError struct {
	Component string
	Issue     string
	Guidance  string
	Cause     error
}

func (e DockerError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("Docker %s issue: %s. %s (cause: %v)", e.Component, e.Issue, e.Guidance, e.Cause)
	}
	return fmt.Sprintf("Docker %s issue: %s. %s", e.Component, e.Issue, e.Guidance)
}

// Unwrap returns the underlying error
func (e DockerError) Unwrap() error {
	return e.Cause
}

// NewDockerError creates a new DockerError
func NewDockerError(component, issue, guidance string, cause error) DockerError {
	return DockerError{
		Component: component,
		Issue:     issue,
		Guidance:  guidance,
		Cause:     cause,
	}
}

// ProjectNameError represents an error with project name resolution or validation
type ProjectNameError struct {
	ProvidedName string
	InferredName string
	Issue        string
	Suggestion   string
}

func (e ProjectNameError) Error() string {
	if e.ProvidedName != "" {
		return fmt.Sprintf("project name error for '%s': %s. %s", e.ProvidedName, e.Issue, e.Suggestion)
	}
	if e.InferredName != "" {
		return fmt.Sprintf("project name error for inferred name '%s': %s. %s", e.InferredName, e.Issue, e.Suggestion)
	}
	return fmt.Sprintf("project name error: %s. %s", e.Issue, e.Suggestion)
}

// NewProjectNameError creates a new ProjectNameError
func NewProjectNameError(providedName, inferredName, issue, suggestion string) ProjectNameError {
	return ProjectNameError{
		ProvidedName: providedName,
		InferredName: inferredName,
		Issue:        issue,
		Suggestion:   suggestion,
	}
}