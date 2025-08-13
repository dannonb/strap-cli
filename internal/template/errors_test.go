package template

import (
	"errors"
	"testing"
)

func TestTemplateError(t *testing.T) {
	originalErr := errors.New("original error")
	templateErr := NewTemplateError("test-template", originalErr)

	// Test Error() method
	expectedMsg := "template error in 'test-template': original error"
	if templateErr.Error() != expectedMsg {
		t.Errorf("TemplateError.Error() = %s, expected %s", templateErr.Error(), expectedMsg)
	}

	// Test Unwrap() method
	if templateErr.Unwrap() != originalErr {
		t.Errorf("TemplateError.Unwrap() = %v, expected %v", templateErr.Unwrap(), originalErr)
	}

	// Test fields
	if templateErr.Template != "test-template" {
		t.Errorf("TemplateError.Template = %s, expected test-template", templateErr.Template)
	}

	if templateErr.Cause != originalErr {
		t.Errorf("TemplateError.Cause = %v, expected %v", templateErr.Cause, originalErr)
	}
}

func TestTemplateNotFoundError(t *testing.T) {
	err := NewTemplateNotFoundError("backend", "fastapi")

	// Test Error() method
	expectedMsg := "template not found for service 'backend' with technology 'fastapi'"
	if err.Error() != expectedMsg {
		t.Errorf("TemplateNotFoundError.Error() = %s, expected %s", err.Error(), expectedMsg)
	}

	// Test fields
	if err.Service != "backend" {
		t.Errorf("TemplateNotFoundError.Service = %s, expected backend", err.Service)
	}

	if err.Technology != "fastapi" {
		t.Errorf("TemplateNotFoundError.Technology = %s, expected fastapi", err.Technology)
	}
}

func TestErrorsImplementErrorInterface(t *testing.T) {
	var err error

	// Test TemplateError implements error interface
	templateErr := NewTemplateError("test", errors.New("test"))
	err = templateErr
	if err == nil {
		t.Error("TemplateError does not implement error interface")
	}

	// Test TemplateNotFoundError implements error interface
	notFoundErr := NewTemplateNotFoundError("service", "tech")
	err = notFoundErr
	if err == nil {
		t.Error("TemplateNotFoundError does not implement error interface")
	}
}