package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"microservice-bootstrapper/internal/generator"
	"microservice-bootstrapper/internal/interfaces"
)

var createCmd = &cobra.Command{
	Use:   "create [flags]",
	Short: "Create a new microservice project",
	Long: `Create a new microservice project with the specified backend, frontend, 
and database technologies. The project will include Docker configuration,
boilerplate code, and proper directory structure.

🔧 AVAILABLE OPTIONS:
  --be     Backend technology (fastapi, express, gin)
  --fe     Frontend technology (react, vue, angular)  
  --db     Database technology (mongo, postgres, mysql, redis)
  --name   Project name (defaults to current directory name)
  --force  Force creation even if directory is not empty

📁 GENERATED STRUCTURE:
  project-name/
  ├── docker-compose.yml    # Orchestration for all services
  ├── README.md            # Setup and usage instructions
  ├── .env.example         # Environment variables template
  ├── .gitignore           # Git ignore patterns
  ├── backend/             # Backend service (if specified)
  │   ├── Dockerfile
  │   ├── main.{ext}       # Entry point with basic API
  │   └── requirements/deps
  ├── frontend/            # Frontend service (if specified)
  │   ├── Dockerfile
  │   ├── package.json
  │   ├── src/             # Source code with basic component
  │   └── public/
  └── docs/
      └── setup.md         # Detailed setup instructions

🚀 QUICK START WORKFLOWS:`,
	Example: `  # Full-stack web application
  strap create --be=fastapi --fe=react --db=postgres --name=webapp
  cd webapp && docker-compose up

  # REST API with caching
  strap create --be=express --db=redis --name=api
  cd api && docker-compose up

  # Frontend SPA only
  strap create --fe=vue --name=frontend
  cd frontend && docker-compose up

  # Go microservice with MongoDB
  strap create --be=gin --db=mongo --name=service
  cd service && docker-compose up

  # Force creation in existing directory
  strap create --be=fastapi --fe=angular --force

  # Multiple microservices (create in separate directories)
  mkdir services && cd services
  strap create --be=fastapi --db=postgres --name=user-service
  strap create --be=gin --db=redis --name=auth-service
  strap create --fe=react --name=web-client`,
	RunE: func(cmd *cobra.Command, args []string) error {
		config := interfaces.CLIConfig{
			Backend:     backend,
			Frontend:    frontend,
			Database:    database,
			ProjectName: projectName,
			Force:       force,
		}
		
		// Create generator and execute project generation
		gen := generator.NewGenerator()
		if err := gen.Generate(config); err != nil {
			// Enhanced error handling with user-friendly messages
			return handleGenerationError(err)
		}
		
		return nil
	},
}

var (
	backend     string
	frontend    string
	database    string
	projectName string
	force       bool
)

// Supported technology options
var (
	supportedBackends  = []string{"fastapi", "express", "gin"}
	supportedFrontends = []string{"react", "vue", "angular"}
	supportedDatabases = []string{"mongo", "postgres", "mysql", "redis"}
)



// handleGenerationError provides user-friendly error messages based on error type
func handleGenerationError(err error) error {
	errMsg := err.Error()
	
	// Check for specific error patterns and provide helpful messages
	if strings.Contains(errMsg, "validation failed") {
		return fmt.Errorf("❌ %w\n\n📚 Get help:\n  • strap create --help\n  • strap examples\n  • strap workflows", err)
	}
	
	// Check if it's a prerequisite error
	if strings.Contains(errMsg, "prerequisite check failed") {
		return fmt.Errorf("❌ %w\n\n🔧 Required tools:\n  • Docker (https://docs.docker.com/get-docker/)\n  • Docker Compose (https://docs.docker.com/compose/install/)", err)
	}
	
	// Check if it's a directory conflict
	if strings.Contains(errMsg, "not empty") {
		return fmt.Errorf("❌ %w\n\n💡 Solutions:\n  • Use --force to overwrite existing files\n  • Choose a different directory\n  • Create a subdirectory: mkdir myproject && cd myproject", err)
	}
	
	// Check for permission errors
	if strings.Contains(errMsg, "permission denied") {
		return fmt.Errorf("❌ %w\n\n🔐 Permission issue:\n  • Check directory write permissions\n  • Try running in a different directory\n  • On Windows, try running as administrator", err)
	}
	
	// Generic error with comprehensive help
	return fmt.Errorf("❌ Project generation failed: %w\n\n📚 Troubleshooting:\n  • strap create --help (usage examples)\n  • strap examples (common patterns)\n  • strap workflows (development guides)\n  • Check Docker is running: docker --version", err)
}

func init() {
	createCmd.Flags().StringVar(&backend, "be", "", 
		fmt.Sprintf("Backend technology: %s", strings.Join(supportedBackends, ", ")))
	createCmd.Flags().StringVar(&frontend, "fe", "", 
		fmt.Sprintf("Frontend technology: %s", strings.Join(supportedFrontends, ", ")))
	createCmd.Flags().StringVar(&database, "db", "", 
		fmt.Sprintf("Database technology: %s", strings.Join(supportedDatabases, ", ")))
	createCmd.Flags().StringVar(&projectName, "name", "", 
		"Project name (defaults to current directory name)")
	createCmd.Flags().BoolVar(&force, "force", false, 
		"Force creation even if directory is not empty")
	
	// Add detailed flag usage information
	createCmd.Flags().SetAnnotation("be", cobra.BashCompCustom, []string{"__strap_complete_backend"})
	createCmd.Flags().SetAnnotation("fe", cobra.BashCompCustom, []string{"__strap_complete_frontend"})
	createCmd.Flags().SetAnnotation("db", cobra.BashCompCustom, []string{"__strap_complete_database"})
	
	// Add validation for flag values
	createCmd.PreRunE = validateFlags
}

// validateFlags performs comprehensive CLI flag validation before execution
func validateFlags(cmd *cobra.Command, args []string) error {
	var validationErrors []string
	
	// Validate backend flag
	if backend != "" && !contains(supportedBackends, backend) {
		validationErrors = append(validationErrors, 
			fmt.Sprintf("❌ Invalid backend '%s'", backend))
		validationErrors = append(validationErrors, 
			fmt.Sprintf("   💡 Supported backends: %s", strings.Join(supportedBackends, ", ")))
	}
	
	// Validate frontend flag
	if frontend != "" && !contains(supportedFrontends, frontend) {
		validationErrors = append(validationErrors, 
			fmt.Sprintf("❌ Invalid frontend '%s'", frontend))
		validationErrors = append(validationErrors, 
			fmt.Sprintf("   💡 Supported frontends: %s", strings.Join(supportedFrontends, ", ")))
	}
	
	// Validate database flag
	if database != "" && !contains(supportedDatabases, database) {
		validationErrors = append(validationErrors, 
			fmt.Sprintf("❌ Invalid database '%s'", database))
		validationErrors = append(validationErrors, 
			fmt.Sprintf("   💡 Supported databases: %s", strings.Join(supportedDatabases, ", ")))
	}
	
	// Check that at least one service is specified
	if backend == "" && frontend == "" && database == "" {
		validationErrors = append(validationErrors, 
			"❌ At least one service must be specified")
		validationErrors = append(validationErrors, 
			"   💡 Use --be for backend, --fe for frontend, or --db for database")
		validationErrors = append(validationErrors, 
			"   💡 Example: strap create --be=fastapi --fe=react --db=postgres")
	}
	
	// Validate project name if provided
	if projectName != "" {
		if strings.Contains(projectName, " ") {
			validationErrors = append(validationErrors, 
				"❌ Project name cannot contain spaces")
			validationErrors = append(validationErrors, 
				"   💡 Use hyphens or underscores instead: my-project or my_project")
		}
		if strings.HasPrefix(projectName, "-") || strings.HasPrefix(projectName, "_") {
			validationErrors = append(validationErrors, 
				"❌ Project name cannot start with - or _")
		}
	}
	
	// Return combined validation errors with helpful formatting
	if len(validationErrors) > 0 {
		errorMsg := fmt.Sprintf("\n🚨 Validation Errors:\n%s\n\n📚 Need help? Try:\n  • strap create --help\n  • strap examples\n  • strap workflows", 
			strings.Join(validationErrors, "\n"))
		return fmt.Errorf(errorMsg)
	}
	
	return nil
}

// contains checks if a slice contains a specific string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}