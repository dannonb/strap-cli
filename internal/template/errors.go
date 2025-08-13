package template

import "fmt"

// TemplateError represents an error that occurred during template processing
type TemplateError struct {
	Template string
	Cause    error
}

// Error implements the error interface
func (e TemplateError) Error() string {
	return fmt.Sprintf("template error in '%s': %v", e.Template, e.Cause)
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

// TemplateNotFoundError represents an error when a template is not found
type TemplateNotFoundError struct {
	Service    string
	Technology string
}

// Error implements the error interface
func (e TemplateNotFoundError) Error() string {
	return fmt.Sprintf("template not found for service '%s' with technology '%s'", e.Service, e.Technology)
}

// NewTemplateNotFoundError creates a new TemplateNotFoundError
func NewTemplateNotFoundError(service, technology string) TemplateNotFoundError {
	return TemplateNotFoundError{
		Service:    service,
		Technology: technology,
	}
}