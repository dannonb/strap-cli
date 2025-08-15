package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"microservice-bootstrapper/internal/config"
	"microservice-bootstrapper/internal/feedback"
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
  --name   Project name (optional - defaults to current directory name)
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
	Example: `  # Full-stack web application (uses directory name as project name)
  mkdir webapp && cd webapp
  strap create --be=fastapi --fe=react --db=postgres
  docker-compose up

  # REST API with caching (explicit name)
  strap create --be=express --db=redis --name=api
  cd api && docker-compose up

  # Frontend SPA only (uses current directory name)
  mkdir my-frontend && cd my-frontend
  strap create --fe=vue
  docker-compose up

  # Go microservice with MongoDB (explicit name)
  strap create --be=gin --db=mongo --name=service
  cd service && docker-compose up

  # Force creation in existing directory (uses directory name)
  strap create --be=fastapi --fe=angular --force

  # Multiple microservices (create in separate directories)
  mkdir user-service && cd user-service
  strap create --be=fastapi --db=postgres
  cd ../
  mkdir auth-service && cd auth-service
  strap create --be=gin --db=redis
  cd ../
  mkdir web-client && cd web-client
  strap create --fe=react`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create messenger for user feedback
		messenger := feedback.NewMessenger()
		
		// Resolve project name (use provided name or infer from directory)
		resolvedName, originalDir, err := resolveProjectNameWithDetails(projectName)
		if err != nil {
			return fmt.Errorf("failed to resolve project name: %w", err)
		}

		config := interfaces.CLIConfig{
			Backend:      backend,
			Frontend:     frontend,
			Database:     database,
			ProjectName:  resolvedName,
			Force:        force,
			InferredName: projectName == "", // True if name was inferred
			WorkingDir:   originalDir,       // Store original directory name
		}
		
		// Show user feedback about project name inference if name was inferred
		if projectName == "" {
			messenger.ShowProjectNameInference(resolvedName, originalDir)
		}
		
		// Create generator and execute project generation
		gen := generator.NewGenerator()
		if err := gen.GenerateWithMessenger(config, messenger); err != nil {
			// Enhanced error handling with user-friendly messages
			return handleGenerationError(err)
		}
		
		// Show success message and next steps
		messenger.ShowGenerationSuccess(resolvedName, ".")
		messenger.ShowNextSteps(config)
		
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
	
	// Handle directory inference errors with enhanced guidance
	if strings.Contains(errMsg, "directory name inference failed") || strings.Contains(errMsg, "cannot access current directory") {
		return fmt.Errorf("❌ Directory Access Issue:\n%w\n\n🔧 Quick Recovery:\n  • Use --name to specify a project name: strap create --be=fastapi --name=myproject\n  • Navigate to a different directory: cd ~ && mkdir myproject && cd myproject\n  • Create a new project directory: mkdir my-microservice && cd my-microservice\n  • Check directory permissions: ls -la", err)
	}
	
	// Handle project name errors with specific examples
	if strings.Contains(errMsg, "project name error") || strings.Contains(errMsg, "cannot use root directory") {
		return fmt.Errorf("❌ Project Name Issue:\n%w\n\n🔧 Project Name Solutions:\n  • Use explicit name: strap create --be=fastapi --name=my-api\n  • Navigate to a proper directory: cd ~/Development && mkdir myproject && cd myproject\n  • Create a project folder: mkdir user-service && cd user-service\n  • Avoid root directories and special characters", err)
	}
	
	// Handle Docker-related errors with specific guidance
	if strings.Contains(errMsg, "Docker") {
		return handleDockerError(err)
	}
	
	// Handle directory and permission errors with platform-specific guidance
	if strings.Contains(errMsg, "permission denied") || strings.Contains(errMsg, "no write permission") {
		return fmt.Errorf("❌ Permission Denied:\n%w\n\n🔐 Permission Recovery:\n  • Check directory permissions: ls -la .\n  • Try your home directory: cd ~ && mkdir myproject && cd myproject\n  • Windows: Run terminal as administrator\n  • Linux/Mac: Check ownership: ls -la | grep $USER\n  • Fix permissions: chmod 755 . (if you own the directory)\n  • Alternative: Use a different directory with write access", err)
	}
	
	// Handle directory conflicts with clear options
	if strings.Contains(errMsg, "not empty") || strings.Contains(errMsg, "directory conflict") {
		return fmt.Errorf("❌ Directory Not Empty:\n%w\n\n🔧 Directory Conflict Recovery:\n  • Use --force to proceed anyway: strap create --be=fastapi --force\n  • Choose a different name: strap create --be=fastapi --name=myproject-v2\n  • Create a new subdirectory: mkdir v2 && cd v2\n  • Clean the directory first: rm -rf * (⚠️  be careful!)\n  • Use a timestamp: strap create --be=fastapi --name=myproject-$(date +%%Y%%m%%d)", err)
	}
	
	// Handle disk space issues with specific actions
	if strings.Contains(errMsg, "no space left") || strings.Contains(errMsg, "insufficient disk space") {
		return fmt.Errorf("❌ Insufficient Disk Space:\n%w\n\n💾 Disk Space Recovery:\n  • Check available space: df -h .\n  • Free up space (need ~100MB for project)\n  • Clean temporary files: rm -rf /tmp/* ~/.cache/* (Linux/Mac)\n  • Windows: Run Disk Cleanup utility\n  • Choose a different drive/directory with more space\n  • Remove unused Docker images: docker system prune", err)
	}
	
	// Handle validation errors with examples
	if strings.Contains(errMsg, "validation failed") {
		return fmt.Errorf("❌ Configuration Validation Failed:\n%w\n\n📚 Validation Recovery:\n  • See usage examples: strap create --help\n  • Valid backends: --be=fastapi, --be=express, --be=gin\n  • Valid frontends: --fe=react, --fe=vue, --fe=angular\n  • Valid databases: --db=postgres, --db=mysql, --db=mongo, --db=redis\n  • Example: strap create --be=fastapi --fe=react --db=postgres\n  • At least one service required", err)
	}
	
	// Handle prerequisite errors with installation guidance
	if strings.Contains(errMsg, "prerequisite check failed") {
		return fmt.Errorf("❌ System Prerequisites Missing:\n%w\n\n🔧 Prerequisite Recovery:\n  • Install missing tools as indicated above\n  • Ensure tools are in your PATH environment variable\n  • Restart terminal after installations\n  • Verify installations: docker --version && docker compose version\n  • Update system PATH if needed\n  • Try running with explicit paths if needed", err)
	}
	
	// Handle template or generation errors
	if strings.Contains(errMsg, "template") || strings.Contains(errMsg, "generation") {
		return fmt.Errorf("❌ Project Generation Failed:\n%w\n\n🔧 Generation Recovery:\n  • Try with a simpler configuration first\n  • Check if the directory is writable\n  • Ensure sufficient disk space\n  • Try: strap create --be=fastapi --name=test-project\n  • Report this issue if it persists", err)
	}
	
	// Handle network or download errors
	if strings.Contains(errMsg, "network") || strings.Contains(errMsg, "download") || strings.Contains(errMsg, "timeout") {
		return fmt.Errorf("❌ Network/Download Issue:\n%w\n\n🌐 Network Recovery:\n  • Check internet connection\n  • Try again in a few minutes\n  • Check if behind a corporate firewall\n  • Verify DNS resolution\n  • Try from a different network if possible", err)
	}
	
	// Generic error with comprehensive troubleshooting
	return fmt.Errorf("❌ Project Generation Failed:\n%w\n\n🔧 General Troubleshooting:\n  • Try with explicit name: strap create --be=fastapi --name=myproject\n  • Check directory permissions and disk space\n  • Ensure Docker is running: docker info\n  • Use a clean directory: mkdir test-project && cd test-project\n  • See examples: strap create --help\n  • Try minimal config first: strap create --be=fastapi", err)
}

// handleDockerError provides specific guidance for Docker-related errors
func handleDockerError(err error) error {
	errMsg := err.Error()
	
	if strings.Contains(errMsg, "not installed") || strings.Contains(errMsg, "not found") || strings.Contains(errMsg, "executable file not found") {
		return fmt.Errorf("❌ Docker Not Installed:\n%w\n\n🐳 Docker Installation Recovery:\n  • Download Docker Desktop: https://docs.docker.com/get-docker/\n  • Windows: Install Docker Desktop for Windows\n  • Mac: Install Docker Desktop for Mac  \n  • Linux: Install Docker Engine or Docker Desktop\n  • Add Docker to system PATH during installation\n  • Restart terminal after installation\n  • Verify installation: docker --version\n  • Test Docker: docker run hello-world", err)
	}
	
	if strings.Contains(errMsg, "not running") || strings.Contains(errMsg, "Cannot connect") || strings.Contains(errMsg, "daemon") {
		return fmt.Errorf("❌ Docker Not Running:\n%w\n\n🐳 Docker Startup Recovery:\n  • Windows/Mac: Start Docker Desktop application\n  • Check system tray for Docker icon (should show running)\n  • Wait 1-2 minutes for Docker to fully initialize\n  • Linux: Start Docker service: sudo systemctl start docker\n  • Linux: Enable auto-start: sudo systemctl enable docker\n  • Verify Docker is running: docker info\n  • If stuck, restart Docker Desktop completely", err)
	}
	
	if strings.Contains(errMsg, "permission denied") || strings.Contains(errMsg, "access denied") {
		return fmt.Errorf("❌ Docker Permission Denied:\n%w\n\n🐳 Docker Permission Recovery:\n  • Linux: Add user to docker group: sudo usermod -aG docker $USER\n  • Linux: Log out and back in after group change\n  • Linux: Verify group: groups $USER | grep docker\n  • Linux: Alternative: sudo docker info (temporary)\n  • Windows/Mac: Ensure Docker Desktop has admin permissions\n  • Windows: Try running terminal as administrator\n  • Restart Docker Desktop after permission changes", err)
	}
	
	if strings.Contains(errMsg, "Compose") || strings.Contains(errMsg, "compose") {
		return fmt.Errorf("❌ Docker Compose Issue:\n%w\n\n🐳 Docker Compose Recovery:\n  • Modern Docker includes Compose: docker compose version\n  • Try modern syntax: docker compose (no hyphen)\n  • Legacy syntax: docker-compose --version (with hyphen)\n  • Update Docker Desktop to latest version\n  • For older Docker: install docker-compose separately\n  • Installation guide: https://docs.docker.com/compose/install/\n  • Verify: docker compose version", err)
	}
	
	if strings.Contains(errMsg, "timeout") || strings.Contains(errMsg, "context deadline") {
		return fmt.Errorf("❌ Docker Timeout:\n%w\n\n🐳 Docker Timeout Recovery:\n  • Docker may be starting - wait 2-3 minutes\n  • Check system resources (CPU/Memory usage)\n  • Close resource-intensive applications\n  • Restart Docker Desktop if it appears stuck\n  • Check Docker Desktop logs for errors\n  • Try: docker system prune (if Docker partially works)\n  • Increase Docker memory allocation in settings", err)
	}
	
	if strings.Contains(errMsg, "version") || strings.Contains(errMsg, "compatibility") {
		return fmt.Errorf("❌ Docker Version Issue:\n%w\n\n🐳 Docker Version Recovery:\n  • Update Docker Desktop to latest version\n  • Check Docker version: docker --version\n  • Check system compatibility requirements\n  • Uninstall and reinstall Docker if needed\n  • Ensure system meets minimum requirements\n  • Check Docker Desktop release notes", err)
	}
	
	return fmt.Errorf("❌ Docker Issue:\n%w\n\n🐳 General Docker Recovery:\n  • Ensure Docker Desktop is installed and running\n  • Check system tray for Docker icon status\n  • Restart Docker Desktop if needed\n  • Verify Docker status: docker info\n  • Check Docker Desktop logs for errors\n  • Visit Docker documentation: https://docs.docker.com/get-docker/\n  • Try: docker run hello-world (basic test)", err)
}

func init() {
	createCmd.Flags().StringVar(&backend, "be", "", 
		fmt.Sprintf("Backend technology: %s", strings.Join(supportedBackends, ", ")))
	createCmd.Flags().StringVar(&frontend, "fe", "", 
		fmt.Sprintf("Frontend technology: %s", strings.Join(supportedFrontends, ", ")))
	createCmd.Flags().StringVar(&database, "db", "", 
		fmt.Sprintf("Database technology: %s", strings.Join(supportedDatabases, ", ")))
	createCmd.Flags().StringVar(&projectName, "name", "", 
		"Project name (optional - defaults to current directory name)")
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
	
	// Validate project name if provided (now optional)
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

// resolveProjectName resolves the project name using provided name or directory inference
func resolveProjectName(providedName string) (string, error) {
	resolvedName, _, err := resolveProjectNameWithDetails(providedName)
	return resolvedName, err
}

// resolveProjectNameWithDetails resolves the project name and returns both resolved name and original directory
func resolveProjectNameWithDetails(providedName string) (string, string, error) {
	// If name is explicitly provided, use it
	if providedName != "" {
		return providedName, "", nil
	}
	
	// Otherwise, infer from current directory
	dirService := config.NewDirectoryService()
	dirName, err := dirService.GetCurrentDirectoryName()
	if err != nil {
		// Handle directory inference errors gracefully
		return handleDirectoryInferenceError(err)
	}
	
	// Sanitize the directory name to be a valid project name
	sanitized := dirService.SanitizeProjectName(dirName)
	
	// Validate the final project name
	if err := validateResolvedProjectName(sanitized, dirName); err != nil {
		return "", "", err
	}
	
	return sanitized, dirName, nil
}

// handleDirectoryInferenceError handles errors during directory name inference
func handleDirectoryInferenceError(err error) (string, string, error) {
	errMsg := err.Error()
	
	// Extract fallback name from error message if available
	fallbackName := "microservice-project"
	if strings.Contains(errMsg, "Using fallback name '") {
		// Extract the fallback name from the error message
		start := strings.Index(errMsg, "Using fallback name '") + len("Using fallback name '")
		end := strings.Index(errMsg[start:], "'")
		if end > 0 {
			fallbackName = errMsg[start : start+end]
		}
	}
	
	// Provide enhanced error context based on the specific issue
	if strings.Contains(errMsg, "permission denied") {
		return fallbackName, "", fmt.Errorf("directory access denied - using fallback name '%s'. "+
			"🔧 Try: Use --name explicitly, navigate to a writable directory, or run with appropriate permissions. "+
			"Original error: %w", fallbackName, err)
	}
	
	if strings.Contains(errMsg, "no such file or directory") {
		return fallbackName, "", fmt.Errorf("current directory no longer exists - using fallback name '%s'. "+
			"🔧 Try: Navigate to a valid directory (cd ~) or use --name explicitly. "+
			"Original error: %w", fallbackName, err)
	}
	
	if strings.Contains(errMsg, "root directory") {
		return fallbackName, "", fmt.Errorf("cannot use root directory as project name - using fallback name '%s'. "+
			"🔧 Try: Create a project directory (mkdir myproject && cd myproject) or use --name explicitly. "+
			"Original error: %w", fallbackName, err)
	}
	
	// Return the fallback name but preserve the error for user feedback
	// The error will be handled by the caller to show appropriate warnings
	return fallbackName, "", fmt.Errorf("directory name inference issue - using fallback name '%s'. "+
		"🔧 Try: Use --name to specify a project name explicitly. "+
		"Original error: %w", fallbackName, err)
}

// validateResolvedProjectName performs final validation on the resolved project name
func validateResolvedProjectName(sanitized, original string) error {
	if sanitized == "" {
		return fmt.Errorf("project name resolution failed: unable to create valid name from directory '%s'. "+
			"Suggestion: Use --name to specify a project name explicitly", original)
	}
	
	if len(sanitized) < 2 {
		return fmt.Errorf("resolved project name '%s' is too short (from directory '%s'). "+
			"Suggestion: Use --name to specify a longer project name", sanitized, original)
	}
	
	// Check if the sanitized name is very different from original (might confuse users)
	if original != "" && sanitized != original && !strings.Contains(sanitized, strings.ToLower(original)) {
		// This is just a warning, not an error - let it proceed but inform the user
		fmt.Printf("ℹ️  Note: Directory name '%s' was significantly changed to '%s' for project compatibility.\n", 
			original, sanitized)
		fmt.Printf("   💡 Use --name to specify a custom project name if preferred.\n")
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