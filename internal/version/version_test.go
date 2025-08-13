package version

import (
	"runtime"
	"strings"
	"testing"
)

func TestInfo(t *testing.T) {
	info := Info()
	
	// Check that info contains expected elements
	expectedElements := []string{
		"Microservice Bootstrapper",
		Version,
		"Build Information",
		"Version:",
		"Built:",
		"Commit:",
		"Go Version:",
		"Platform:",
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
	}
	
	for _, element := range expectedElements {
		if !strings.Contains(info, element) {
			t.Errorf("Info() does not contain expected element: %s", element)
		}
	}
}

func TestShort(t *testing.T) {
	short := Short()
	expected := "v" + Version
	
	if short != expected {
		t.Errorf("Short() = %s, want %s", short, expected)
	}
}

func TestDetailed(t *testing.T) {
	detailed := Detailed()
	
	// Check that detailed contains expected elements
	expectedElements := []string{
		"Microservice Bootstrapper",
		Version,
		"System Information",
		"Build Details",
		"Features",
		"Multi-technology project generation",
		"Docker and Docker Compose integration",
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
		runtime.Compiler,
		BuildDate,
		GitCommit,
	}
	
	for _, element := range expectedElements {
		if !strings.Contains(detailed, element) {
			t.Errorf("Detailed() does not contain expected element: %s", element)
		}
	}
}

func TestVersionVariables(t *testing.T) {
	// Test that version variables are set to expected defaults
	if Version == "" {
		t.Error("Version should not be empty")
	}
	
	// BuildDate and GitCommit can be "unknown" in development
	if BuildDate == "" {
		t.Error("BuildDate should not be empty (can be 'unknown')")
	}
	
	if GitCommit == "" {
		t.Error("GitCommit should not be empty (can be 'unknown')")
	}
}

func TestVersionFormat(t *testing.T) {
	// Test that version follows semantic versioning pattern
	versionParts := strings.Split(Version, ".")
	if len(versionParts) != 3 {
		t.Errorf("Version should follow semantic versioning (x.y.z), got %s", Version)
	}
	
	// Each part should be numeric (basic check)
	for i, part := range versionParts {
		if part == "" {
			t.Errorf("Version part %d should not be empty", i)
		}
	}
}

func TestInfoFormatting(t *testing.T) {
	info := Info()
	
	// Check that info is properly formatted with emojis and sections
	if !strings.Contains(info, "🚀") {
		t.Error("Info() should contain rocket emoji")
	}
	
	if !strings.Contains(info, "📋") {
		t.Error("Info() should contain clipboard emoji")
	}
	
	// Check that it's multi-line
	lines := strings.Split(info, "\n")
	if len(lines) < 5 {
		t.Error("Info() should be multi-line with at least 5 lines")
	}
}

func TestDetailedFormatting(t *testing.T) {
	detailed := Detailed()
	
	// Check that detailed contains all expected emojis
	expectedEmojis := []string{"🚀", "📋", "🔧", "📦", "🌟"}
	for _, emoji := range expectedEmojis {
		if !strings.Contains(detailed, emoji) {
			t.Errorf("Detailed() should contain emoji: %s", emoji)
		}
	}
	
	// Check that it's multi-line with substantial content
	lines := strings.Split(detailed, "\n")
	if len(lines) < 15 {
		t.Error("Detailed() should be multi-line with at least 15 lines")
	}
}

func TestRuntimeInformation(t *testing.T) {
	info := Info()
	
	// Verify that runtime information is correctly included
	if !strings.Contains(info, runtime.Version()) {
		t.Error("Info() should contain Go version from runtime")
	}
	
	if !strings.Contains(info, runtime.GOOS) {
		t.Error("Info() should contain OS from runtime")
	}
	
	if !strings.Contains(info, runtime.GOARCH) {
		t.Error("Info() should contain architecture from runtime")
	}
}

func TestBuildVariablesInOutput(t *testing.T) {
	info := Info()
	detailed := Detailed()
	
	// Both should contain build variables
	outputs := []string{info, detailed}
	variables := []string{BuildDate, GitCommit}
	
	for _, output := range outputs {
		for _, variable := range variables {
			if !strings.Contains(output, variable) {
				t.Errorf("Output should contain build variable: %s", variable)
			}
		}
	}
}