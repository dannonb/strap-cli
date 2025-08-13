package cli

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestRootCommand(t *testing.T) {
	// Test that root command is properly configured
	if rootCmd.Use != "strap" {
		t.Errorf("rootCmd.Use = %s, want 'strap'", rootCmd.Use)
	}
	
	if rootCmd.Short == "" {
		t.Error("rootCmd.Short should not be empty")
	}
	
	if rootCmd.Long == "" {
		t.Error("rootCmd.Long should not be empty")
	}
}

func TestRootCommandHelp(t *testing.T) {
	// Test root command long description contains expected content
	expectedSections := []string{
		"Microservice Bootstrapper",
		"SUPPORTED TECHNOLOGIES",
		"FastAPI",
		"React",
		"MongoDB",
		"WHAT YOU GET",
		"COMMON WORKFLOWS",
		"strap create --help",
	}
	
	for _, section := range expectedSections {
		if !strings.Contains(rootCmd.Long, section) {
			t.Errorf("Root command long description does not contain expected section: %s", section)
		}
	}
	
	// Test that examples are present
	if !strings.Contains(rootCmd.Example, "strap create") {
		t.Error("Root command example should contain 'strap create'")
	}
}

func TestVersionCommand(t *testing.T) {
	// Test version command configuration
	if versionCmd.Use != "version" {
		t.Errorf("versionCmd.Use = %s, want 'version'", versionCmd.Use)
	}
	
	if versionCmd.Short == "" {
		t.Error("versionCmd.Short should not be empty")
	}
	
	if versionCmd.Long == "" {
		t.Error("versionCmd.Long should not be empty")
	}
	
	// Test that version command has a run function
	if versionCmd.Run == nil {
		t.Error("versionCmd.Run should not be nil")
	}
}

func TestExamplesCommand(t *testing.T) {
	// Test examples command configuration
	if examplesCmd.Use != "examples" {
		t.Errorf("examplesCmd.Use = %s, want 'examples'", examplesCmd.Use)
	}
	
	if examplesCmd.Short == "" {
		t.Error("examplesCmd.Short should not be empty")
	}
	
	if examplesCmd.Long == "" {
		t.Error("examplesCmd.Long should not be empty")
	}
	
	// Check that examples long description contains expected content
	expectedContent := []string{
		"FULL-STACK APPLICATIONS",
		"API SERVICES",
		"FRONTEND APPLICATIONS",
		"MICROSERVICES ARCHITECTURE",
		"strap create --be=fastapi",
		"strap create --fe=react",
		"--db=postgres",
	}
	
	for _, content := range expectedContent {
		if !strings.Contains(examplesCmd.Long, content) {
			t.Errorf("Examples long description does not contain expected content: %s", content)
		}
	}
}

func TestWorkflowsCommand(t *testing.T) {
	// Test workflows command configuration
	if workflowsCmd.Use != "workflows" {
		t.Errorf("workflowsCmd.Use = %s, want 'workflows'", workflowsCmd.Use)
	}
	
	if workflowsCmd.Short == "" {
		t.Error("workflowsCmd.Short should not be empty")
	}
	
	if workflowsCmd.Long == "" {
		t.Error("workflowsCmd.Long should not be empty")
	}
	
	// Check that workflows long description contains expected content
	expectedContent := []string{
		"QUICK START WORKFLOW",
		"DEVELOPMENT WORKFLOW",
		"PRODUCTION DEPLOYMENT",
		"TESTING WORKFLOW",
		"CUSTOMIZATION WORKFLOW",
		"docker-compose up",
		"npm test",
		"pytest",
	}
	
	for _, content := range expectedContent {
		if !strings.Contains(workflowsCmd.Long, content) {
			t.Errorf("Workflows long description does not contain expected content: %s", content)
		}
	}
}

func TestSubcommands(t *testing.T) {
	// Test that all expected subcommands are registered
	expectedCommands := []string{"create", "version", "examples", "workflows"}
	
	for _, cmdName := range expectedCommands {
		cmd, _, err := rootCmd.Find([]string{cmdName})
		if err != nil {
			t.Errorf("Command %s not found: %v", cmdName, err)
			continue
		}
		
		if cmd.Name() != cmdName {
			t.Errorf("Expected command name %s, got %s", cmdName, cmd.Name())
		}
	}
}

func TestExecute(t *testing.T) {
	// Test that Execute function exists and can be called
	// We can't easily test the actual execution without mocking,
	// but we can test that the function doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Execute() panicked: %v", r)
		}
	}()
	
	// This will fail because no args are provided, but it shouldn't panic
	_ = Execute()
}

func TestCommandConfiguration(t *testing.T) {
	// Test that commands are properly configured
	tests := []struct {
		name    string
		cmd     *cobra.Command
		wantUse string
	}{
		{
			name:    "root command",
			cmd:     rootCmd,
			wantUse: "strap",
		},
		{
			name:    "create command",
			cmd:     createCmd,
			wantUse: "create [flags]",
		},
		{
			name:    "version command",
			cmd:     versionCmd,
			wantUse: "version",
		},
		{
			name:    "examples command",
			cmd:     examplesCmd,
			wantUse: "examples",
		},
		{
			name:    "workflows command",
			cmd:     workflowsCmd,
			wantUse: "workflows",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.cmd.Use != tt.wantUse {
				t.Errorf("%s Use = %s, want %s", tt.name, tt.cmd.Use, tt.wantUse)
			}
			
			if tt.cmd.Short == "" {
				t.Errorf("%s Short description should not be empty", tt.name)
			}
			
			if tt.cmd.Long == "" {
				t.Errorf("%s Long description should not be empty", tt.name)
			}
		})
	}
}

func TestHelpCommand(t *testing.T) {
	// Test that help command is properly configured
	if helpCmd.Use != "help [command]" {
		t.Errorf("helpCmd.Use = %s, want 'help [command]'", helpCmd.Use)
	}
	
	if helpCmd.Short == "" {
		t.Error("helpCmd.Short should not be empty")
	}
	
	if helpCmd.Long == "" {
		t.Error("helpCmd.Long should not be empty")
	}
}

func TestCompletionDisabled(t *testing.T) {
	// Test that completion is disabled as configured
	if !rootCmd.CompletionOptions.DisableDefaultCmd {
		t.Error("Completion should be disabled")
	}
}