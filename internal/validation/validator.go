package validation

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"microservice-bootstrapper/internal/interfaces"
	"microservice-bootstrapper/pkg/errors"
)

// Validator handles input validation and prerequisite checking
type Validator struct {
	supportedBackends  []string
	supportedFrontends []string
	supportedDatabases []string
}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	return &Validator{
		supportedBackends:  []string{"fastapi", "express", "gin"},
		supportedFrontends: []string{"react", "vue", "angular"},
		supportedDatabases: []string{"mongo", "postgres", "mysql", "redis"},
	}
}

// ValidateConfig performs comprehensive validation of CLI configuration
func (v *Validator) ValidateConfig(config interfaces.CLIConfig) error {
	var validationErrors []error

	// Validate individual technology options
	if err := v.validateBackend(config.Backend); err != nil {
		validationErrors = append(validationErrors, err)
	}

	if err := v.validateFrontend(config.Frontend); err != nil {
		validationErrors = append(validationErrors, err)
	}

	if err := v.validateDatabase(config.Database); err != nil {
		validationErrors = append(validationErrors, err)
	}

	// Validate project name
	if err := v.validateProjectName(config.ProjectName); err != nil {
		validationErrors = append(validationErrors, err)
	}

	// Validate service combinations
	if err := v.validateServiceCombinations(config); err != nil {
		validationErrors = append(validationErrors, err)
	}

	// Validate port conflicts
	if err := v.validatePortConflicts(config); err != nil {
		validationErrors = append(validationErrors, err)
	}

	// Return combined validation errors
	if len(validationErrors) > 0 {
		return v.combineValidationErrors(validationErrors)
	}

	return nil
}

// validateBackend validates the backend technology option
func (v *Validator) validateBackend(backend string) error {
	if backend == "" {
		return nil // Optional field
	}

	if !v.contains(v.supportedBackends, backend) {
		return errors.ValidationError{
			Field: "backend",
			Message: fmt.Sprintf("unsupported backend '%s'. Supported options: %s. "+
				"Suggestion: Use one of the supported backend technologies or omit --be flag",
				backend, strings.Join(v.supportedBackends, ", ")),
		}
	}

	return nil
}

// validateFrontend validates the frontend technology option
func (v *Validator) validateFrontend(frontend string) error {
	if frontend == "" {
		return nil // Optional field
	}

	if !v.contains(v.supportedFrontends, frontend) {
		return errors.ValidationError{
			Field: "frontend",
			Message: fmt.Sprintf("unsupported frontend '%s'. Supported options: %s. "+
				"Suggestion: Use one of the supported frontend technologies or omit --fe flag",
				frontend, strings.Join(v.supportedFrontends, ", ")),
		}
	}

	return nil
}

// validateDatabase validates the database technology option
func (v *Validator) validateDatabase(database string) error {
	if database == "" {
		return nil // Optional field
	}

	if !v.contains(v.supportedDatabases, database) {
		return errors.ValidationError{
			Field: "database",
			Message: fmt.Sprintf("unsupported database '%s'. Supported options: %s. "+
				"Suggestion: Use one of the supported database technologies or omit --db flag",
				database, strings.Join(v.supportedDatabases, ", ")),
		}
	}

	return nil
}

// validateProjectName validates the project name format
func (v *Validator) validateProjectName(projectName string) error {
	if projectName == "" {
		return nil // Will use default from current directory
	}

	// Check for valid project name format (alphanumeric, hyphens, underscores)
	validName := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validName.MatchString(projectName) {
		return errors.ValidationError{
			Field: "project-name",
			Message: fmt.Sprintf("invalid project name '%s'. Project names can only contain "+
				"letters, numbers, hyphens, and underscores. "+
				"Suggestion: Use a name like 'my-project' or 'my_service'", projectName),
		}
	}

	// Check length constraints
	if len(projectName) < 2 {
		return errors.ValidationError{
			Field: "project-name",
			Message: "project name must be at least 2 characters long. " +
				"Suggestion: Use a more descriptive name",
		}
	}

	if len(projectName) > 50 {
		return errors.ValidationError{
			Field: "project-name",
			Message: "project name must be 50 characters or less. " +
				"Suggestion: Use a shorter, more concise name",
		}
	}

	// Check for reserved names
	reservedNames := []string{"docker", "compose", "node_modules", "venv", "env", "build", "dist"}
	for _, reserved := range reservedNames {
		if strings.EqualFold(projectName, reserved) {
			return errors.ValidationError{
				Field: "project-name",
				Message: fmt.Sprintf("project name '%s' is reserved and cannot be used. "+
					"Suggestion: Choose a different name", projectName),
			}
		}
	}

	return nil
}

// validateServiceCombinations validates that at least one service is specified and combinations make sense
func (v *Validator) validateServiceCombinations(config interfaces.CLIConfig) error {
	// Ensure at least one service is specified
	if config.Backend == "" && config.Frontend == "" && config.Database == "" {
		return errors.ValidationError{
			Field: "services",
			Message: "at least one service must be specified (--be, --fe, or --db). " +
				"Suggestion: Add --be=fastapi for a backend service, --fe=react for a frontend, " +
				"or --db=postgres for a database",
		}
	}

	// Validate specific combinations
	if config.Frontend != "" && config.Backend == "" && config.Database == "" {
		// Frontend-only is valid, but warn about common patterns
		// This is just informational, not an error
	}

	if config.Database != "" && config.Backend == "" {
		// Database without backend is unusual but valid for some use cases
		// This is just informational, not an error
	}

	return nil
}

// validatePortConflicts checks for potential port conflicts
func (v *Validator) validatePortConflicts(config interfaces.CLIConfig) error {
	usedPorts := make(map[int]string)

	// Get default ports for each service
	if config.Backend != "" {
		port := v.getDefaultPort("backend", config.Backend)
		if service, exists := usedPorts[port]; exists {
			return errors.ValidationError{
				Field: "ports",
				Message: fmt.Sprintf("port conflict: both %s and backend (%s) use port %d. "+
					"Suggestion: This is handled automatically by Docker Compose, but you may need to "+
					"adjust ports in the generated docker-compose.yml if running multiple projects",
					service, config.Backend, port),
			}
		}
		usedPorts[port] = fmt.Sprintf("backend (%s)", config.Backend)
	}

	if config.Frontend != "" {
		port := v.getDefaultPort("frontend", config.Frontend)
		if service, exists := usedPorts[port]; exists {
			return errors.ValidationError{
				Field: "ports",
				Message: fmt.Sprintf("port conflict: both %s and frontend (%s) use port %d. "+
					"Suggestion: This is handled automatically by Docker Compose, but you may need to "+
					"adjust ports in the generated docker-compose.yml if running multiple projects",
					service, config.Frontend, port),
			}
		}
		usedPorts[port] = fmt.Sprintf("frontend (%s)", config.Frontend)
	}

	if config.Database != "" {
		port := v.getDefaultPort("database", config.Database)
		if service, exists := usedPorts[port]; exists {
			return errors.ValidationError{
				Field: "ports",
				Message: fmt.Sprintf("port conflict: both %s and database (%s) use port %d. "+
					"Suggestion: This is handled automatically by Docker Compose, but you may need to "+
					"adjust ports in the generated docker-compose.yml if running multiple projects",
					service, config.Database, port),
			}
		}
		usedPorts[port] = fmt.Sprintf("database (%s)", config.Database)
	}

	return nil
}

// CheckPrerequisites verifies that required tools and environment are available
func (v *Validator) CheckPrerequisites() error {
	var prerequisiteErrors []error

	// Check Docker
	if err := v.checkDocker(); err != nil {
		prerequisiteErrors = append(prerequisiteErrors, err)
	}

	// Check Docker Compose
	if err := v.checkDockerCompose(); err != nil {
		prerequisiteErrors = append(prerequisiteErrors, err)
	}

	// Check disk space (basic check)
	if err := v.checkDiskSpace(); err != nil {
		prerequisiteErrors = append(prerequisiteErrors, err)
	}

	// Check write permissions
	if err := v.checkWritePermissions(); err != nil {
		prerequisiteErrors = append(prerequisiteErrors, err)
	}

	// Return combined prerequisite errors
	if len(prerequisiteErrors) > 0 {
		return v.combinePrerequisiteErrors(prerequisiteErrors)
	}

	return nil
}

// CheckGenerationPrerequisites verifies only what's needed for file generation
func (v *Validator) CheckGenerationPrerequisites() error {
	var prerequisiteErrors []error

	// Check disk space (basic check)
	if err := v.checkDiskSpace(); err != nil {
		prerequisiteErrors = append(prerequisiteErrors, err)
	}

	// Check write permissions
	if err := v.checkWritePermissions(); err != nil {
		prerequisiteErrors = append(prerequisiteErrors, err)
	}

	// Return combined prerequisite errors
	if len(prerequisiteErrors) > 0 {
		return v.combinePrerequisiteErrors(prerequisiteErrors)
	}

	return nil
}

// CheckExecutionPrerequisites verifies what's needed for running the generated project
func (v *Validator) CheckExecutionPrerequisites() error {
	var prerequisiteErrors []error

	// Check Docker
	if err := v.checkDocker(); err != nil {
		prerequisiteErrors = append(prerequisiteErrors, err)
	}

	// Check Docker Compose
	if err := v.checkDockerCompose(); err != nil {
		prerequisiteErrors = append(prerequisiteErrors, err)
	}

	// Return combined prerequisite errors
	if len(prerequisiteErrors) > 0 {
		return v.combinePrerequisiteErrors(prerequisiteErrors)
	}

	return nil
}

// checkDocker verifies Docker is installed and running
func (v *Validator) checkDocker() error {
	// Check if Docker command exists
	if err := v.checkCommand("docker", "--version"); err != nil {
		return v.handleDockerNotInstalled(err)
	}

	// Check if Docker daemon is running
	if err := v.checkCommand("docker", "info"); err != nil {
		return v.handleDockerNotRunning(err)
	}

	return nil
}

// handleDockerNotInstalled provides specific guidance when Docker is not installed
func (v *Validator) handleDockerNotInstalled(err error) error {
	if strings.Contains(err.Error(), "executable file not found") {
		return fmt.Errorf("Docker is not installed or not found in PATH. "+
			"🐳 Docker Installation Guide:\n"+
			"  • Download Docker Desktop: https://docs.docker.com/get-docker/\n"+
			"  • Windows: Install Docker Desktop for Windows\n"+
			"  • Mac: Install Docker Desktop for Mac\n"+
			"  • Linux: Install Docker Engine or Docker Desktop\n"+
			"  • After installation, restart your terminal\n"+
			"  • Verify installation: docker --version\n"+
			"  • Ensure Docker is in your system PATH\n"+
			"Original error: %w", err)
	}
	
	if strings.Contains(err.Error(), "permission denied") {
		return fmt.Errorf("Docker command found but permission denied. "+
			"🐳 Docker Permission Solutions:\n"+
			"  • Linux: Add user to docker group: sudo usermod -aG docker $USER\n"+
			"  • Linux: Log out and back in after adding to group\n"+
			"  • Windows/Mac: Ensure Docker Desktop is properly installed\n"+
			"  • Windows: Try running terminal as administrator\n"+
			"  • Verify Docker Desktop is running and accessible\n"+
			"  • Test with: docker ps\n"+
			"Original error: %w", err)
	}
	
	if strings.Contains(err.Error(), "access denied") || strings.Contains(err.Error(), "denied") {
		return fmt.Errorf("Docker access denied. "+
			"🐳 Docker Access Solutions:\n"+
			"  • Ensure Docker Desktop is running\n"+
			"  • Check Docker Desktop settings and permissions\n"+
			"  • On Windows: Run as administrator if needed\n"+
			"  • On Linux: Check docker group membership: groups $USER\n"+
			"  • Restart Docker Desktop if recently installed\n"+
			"Original error: %w", err)
	}
	
	return fmt.Errorf("Docker is required but not accessible. "+
		"🐳 General Docker Setup:\n"+
		"  • Install Docker Desktop: https://docs.docker.com/get-docker/\n"+
		"  • Ensure Docker is running and accessible\n"+
		"  • Restart terminal after installation\n"+
		"  • Verify with: docker --version && docker info\n"+
		"Original error: %w", err)
}

// handleDockerNotRunning provides specific guidance when Docker is installed but not running
func (v *Validator) handleDockerNotRunning(err error) error {
	if strings.Contains(err.Error(), "Cannot connect to the Docker daemon") {
		return fmt.Errorf("Docker is installed but the Docker daemon is not running. "+
			"🐳 Docker Startup Solutions:\n"+
			"  • Windows/Mac: Start Docker Desktop application\n"+
			"  • Linux: Start Docker service: sudo systemctl start docker\n"+
			"  • Wait 30-60 seconds for Docker to fully initialize\n"+
			"  • Check Docker Desktop system tray icon (should be running)\n"+
			"  • Verify startup: docker info\n"+
			"  • If still failing, restart Docker Desktop completely\n"+
			"Original error: %w", err)
	}
	
	if strings.Contains(err.Error(), "permission denied") {
		return fmt.Errorf("Docker daemon is running but permission denied. "+
			"🐳 Docker Permission Solutions:\n"+
			"  • Linux: Add user to docker group: sudo usermod -aG docker $USER\n"+
			"  • Linux: Log out and back in after group change\n"+
			"  • Linux: Alternative: sudo docker info (temporary fix)\n"+
			"  • Windows/Mac: Ensure Docker Desktop has proper permissions\n"+
			"  • Windows: Try running terminal as administrator\n"+
			"  • Verify group membership: groups $USER | grep docker\n"+
			"Original error: %w", err)
	}
	
	if strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "context deadline exceeded") {
		return fmt.Errorf("Docker daemon is not responding (timeout). "+
			"🐳 Docker Timeout Solutions:\n"+
			"  • Docker may still be starting - wait 1-2 minutes\n"+
			"  • Check Docker Desktop status in system tray\n"+
			"  • Restart Docker Desktop if it appears stuck\n"+
			"  • Check system resources (CPU/Memory usage)\n"+
			"  • Close other resource-intensive applications\n"+
			"  • Try: docker system prune (if Docker partially works)\n"+
			"Original error: %w", err)
	}
	
	if strings.Contains(err.Error(), "dial unix") || strings.Contains(err.Error(), "socket") {
		return fmt.Errorf("Docker socket connection failed. "+
			"🐳 Docker Socket Solutions:\n"+
			"  • Ensure Docker Desktop is fully started\n"+
			"  • Linux: Check Docker socket: ls -la /var/run/docker.sock\n"+
			"  • Linux: Fix socket permissions: sudo chmod 666 /var/run/docker.sock\n"+
			"  • Restart Docker service: sudo systemctl restart docker\n"+
			"  • Windows/Mac: Restart Docker Desktop\n"+
			"Original error: %w", err)
	}
	
	return fmt.Errorf("Docker is installed but not accessible. "+
		"🐳 General Docker Troubleshooting:\n"+
		"  • Ensure Docker Desktop is running (check system tray)\n"+
		"  • Wait for Docker to fully start (can take 1-2 minutes)\n"+
		"  • Restart Docker Desktop if needed\n"+
		"  • Check Docker status: docker info\n"+
		"  • Verify system resources are sufficient\n"+
		"Original error: %w", err)
}

// checkDockerCompose verifies Docker Compose is available
func (v *Validator) checkDockerCompose() error {
	// Try modern Docker Compose (docker compose)
	modernErr := v.checkCommand("docker", "compose", "version")
	if modernErr == nil {
		return nil
	}

	// Try legacy Docker Compose (docker-compose)
	legacyErr := v.checkCommand("docker-compose", "--version")
	if legacyErr == nil {
		return nil
	}

	return v.handleDockerComposeNotFound(modernErr, legacyErr)
}

// handleDockerComposeNotFound provides specific guidance when Docker Compose is not found
func (v *Validator) handleDockerComposeNotFound(modernErr, legacyErr error) error {
	// Check if it's a Docker daemon issue first
	if strings.Contains(modernErr.Error(), "Cannot connect to the Docker daemon") {
		return fmt.Errorf("Docker Compose check failed because Docker daemon is not running. "+
			"🐳 Docker Daemon Solutions:\n"+
			"  • Start Docker Desktop first\n"+
			"  • Wait for Docker to fully initialize\n"+
			"  • Verify Docker is running: docker info\n"+
			"  • Docker Compose is included with modern Docker installations\n"+
			"  • Try again after Docker is running\n"+
			"Original error: %v", modernErr)
	}
	
	// Check if it's a permission issue
	if strings.Contains(modernErr.Error(), "permission denied") {
		return fmt.Errorf("Docker Compose check failed due to permission issues. "+
			"🐳 Docker Compose Permission Solutions:\n"+
			"  • Linux: Add user to docker group: sudo usermod -aG docker $USER\n"+
			"  • Linux: Log out and back in after group change\n"+
			"  • Windows/Mac: Ensure Docker Desktop has proper permissions\n"+
			"  • Try running terminal as administrator (Windows)\n"+
			"  • Verify Docker permissions: docker ps\n"+
			"Original error: %v", modernErr)
	}
	
	// Check if modern compose is not available but legacy might work
	if strings.Contains(modernErr.Error(), "unknown command") {
		if legacyErr != nil && strings.Contains(legacyErr.Error(), "executable file not found") {
			return fmt.Errorf("Docker Compose is not installed. "+
				"🐳 Docker Compose Installation:\n"+
				"  • Modern Docker (recommended): Update to Docker Desktop latest\n"+
				"  • Modern Docker includes 'docker compose' (no hyphen)\n"+
				"  • Legacy: Install docker-compose separately\n"+
				"  • Installation guide: https://docs.docker.com/compose/install/\n"+
				"  • Verify installation: docker compose version\n"+
				"  • Alternative: docker-compose --version (legacy)\n"+
				"Modern error: %v, Legacy error: %v", modernErr, legacyErr)
		}
		
		return fmt.Errorf("Docker Compose command not recognized. "+
			"🐳 Docker Compose Command Solutions:\n"+
			"  • Update Docker Desktop to latest version\n"+
			"  • Modern syntax: docker compose (no hyphen)\n"+
			"  • Legacy syntax: docker-compose (with hyphen)\n"+
			"  • Check Docker version: docker --version\n"+
			"  • Reinstall Docker Desktop if needed\n"+
			"Modern error: %v", modernErr)
	}
	
	// Check for version compatibility issues
	if strings.Contains(modernErr.Error(), "version") || strings.Contains(modernErr.Error(), "compatibility") {
		return fmt.Errorf("Docker Compose version compatibility issue. "+
			"🐳 Docker Compose Version Solutions:\n"+
			"  • Update Docker Desktop to latest version\n"+
			"  • Check Docker version: docker --version\n"+
			"  • Check Compose version: docker compose version\n"+
			"  • Uninstall and reinstall Docker Desktop if needed\n"+
			"  • Ensure system meets Docker requirements\n"+
			"Modern error: %v, Legacy error: %v", modernErr, legacyErr)
	}
	
	return fmt.Errorf("Docker Compose is required but not found. "+
		"🐳 General Docker Compose Solutions:\n"+
		"  • Install/Update Docker Desktop: https://docs.docker.com/get-docker/\n"+
		"  • Modern Docker includes Compose as 'docker compose'\n"+
		"  • Verify installation: docker compose version\n"+
		"  • For older Docker: install docker-compose separately\n"+
		"  • Installation guide: https://docs.docker.com/compose/install/\n"+
		"  • Restart terminal after installation\n"+
		"Modern error: %v, Legacy error: %v", modernErr, legacyErr)
}

// checkDiskSpace performs a basic disk space check
func (v *Validator) checkDiskSpace() error {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("unable to check current directory: %w", err)
	}

	// Check if we can write to the current directory (basic check)
	testFile := filepath.Join(cwd, ".strap_test_write")
	file, err := os.Create(testFile)
	if err != nil {
		return fmt.Errorf("insufficient disk space or write permissions in current directory. "+
			"Please ensure you have at least 100MB free space and write permissions. "+
			"Error: %w", err)
	}
	file.Close()
	os.Remove(testFile) // Clean up test file

	return nil
}

// checkWritePermissions verifies write permissions in current directory
func (v *Validator) checkWritePermissions() error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("unable to check current directory: %w", err)
	}

	// Test write permissions
	testFile := filepath.Join(cwd, ".strap_permission_test")
	file, err := os.Create(testFile)
	if err != nil {
		return fmt.Errorf("no write permissions in current directory '%s'. "+
			"Please ensure you have write permissions or change to a different directory. "+
			"Error: %w", cwd, err)
	}
	file.Close()
	os.Remove(testFile) // Clean up test file

	return nil
}

// ValidateDirectoryConflicts checks for existing files and handles conflicts
func (v *Validator) ValidateDirectoryConflicts(basePath string, force bool) error {
	// Check if path exists
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		return nil // No conflict, directory doesn't exist
	}

	// Check if it's a directory
	info, err := os.Stat(basePath)
	if err != nil {
		return fmt.Errorf("unable to check path '%s': %w", basePath, err)
	}

	if !info.IsDir() {
		return fmt.Errorf("path '%s' exists but is not a directory. "+
			"Please choose a different project name or remove the existing file", basePath)
	}

	// Check if directory is empty
	if !v.isDirectoryEmpty(basePath) {
		if !force {
			files, _ := v.listDirectoryContents(basePath, 5) // Get first 5 files for display
			fileList := strings.Join(files, ", ")
			if len(files) == 5 {
				fileList += "..."
			}

			return fmt.Errorf("directory '%s' is not empty (contains: %s). "+
				"Use --force to proceed anyway, which may overwrite existing files. "+
				"Suggestion: Choose a different project name or use --force flag", basePath, fileList)
		}

		// Warn about force usage
		fmt.Printf("⚠️  Warning: Directory '%s' is not empty. Proceeding with --force flag. "+
			"Existing files may be overwritten.\n", basePath)
	}

	return nil
}

// Helper functions

// checkCommand checks if a command is available and executable
func (v *Validator) checkCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = nil // Suppress output
	cmd.Stderr = nil // Suppress error output
	return cmd.Run()
}

// contains checks if a slice contains a specific string
func (v *Validator) contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// getDefaultPort returns the default port for a service type and technology
func (v *Validator) getDefaultPort(serviceType, technology string) int {
	switch serviceType {
	case "backend":
		switch technology {
		case "fastapi":
			return 8000
		case "express":
			return 3000
		case "gin":
			return 8080
		default:
			return 8000
		}
	case "frontend":
		switch technology {
		case "react", "vue":
			return 3000
		case "angular":
			return 4200
		default:
			return 3000
		}
	case "database":
		switch technology {
		case "mongo":
			return 27017
		case "postgres":
			return 5432
		case "mysql":
			return 3306
		case "redis":
			return 6379
		default:
			return 5432
		}
	}
	return 8000
}

// isDirectoryEmpty checks if a directory is empty
func (v *Validator) isDirectoryEmpty(path string) bool {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false
	}
	return len(entries) == 0
}

// listDirectoryContents returns a list of files/directories in the given path (up to maxItems)
func (v *Validator) listDirectoryContents(path string, maxItems int) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var files []string
	for i, entry := range entries {
		if i >= maxItems {
			break
		}
		files = append(files, entry.Name())
	}

	return files, nil
}

// combineValidationErrors combines multiple validation errors into a single error
func (v *Validator) combineValidationErrors(validationErrors []error) error {
	var messages []string
	for _, err := range validationErrors {
		messages = append(messages, err.Error())
	}

	return fmt.Errorf("validation failed:\n  - %s", strings.Join(messages, "\n  - "))
}

// combinePrerequisiteErrors combines multiple prerequisite errors into a single error
func (v *Validator) combinePrerequisiteErrors(prerequisiteErrors []error) error {
	var messages []string
	for _, err := range prerequisiteErrors {
		messages = append(messages, err.Error())
	}

	return fmt.Errorf("prerequisite check failed:\n  - %s\n\n"+
		"Please resolve these issues before running the command again",
		strings.Join(messages, "\n  - "))
}