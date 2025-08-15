package feedback

import (
	"fmt"
	"io"
	"os"

	"microservice-bootstrapper/internal/interfaces"
)

// Messenger provides consistent user communication throughout the project generation process
type Messenger interface {
	ShowProjectNameInference(inferred, directory string)
	ShowDockerWarning()
	ShowGenerationSuccess(projectName, path string)
	ShowNextSteps(config interfaces.CLIConfig)
	ShowDirectoryInferenceWarning(issue, suggestion string)
	ShowFallbackNameUsage(fallbackName, reason string)
	ShowErrorRecoveryGuidance(errorType, guidance string)
	ShowDirectoryAccessError(path string, errorType string, solutions []string)
	ShowDockerRecoverySteps(issue string, steps []string)
	ShowPermissionError(path string, platform string)
	ShowDiskSpaceError(path string, availableSpace string)
}

// messenger implements the Messenger interface
type messenger struct {
	output io.Writer
}

// NewMessenger creates a new messenger instance that writes to stdout
func NewMessenger() Messenger {
	return &messenger{
		output: os.Stdout,
	}
}

// NewMessengerWithOutput creates a new messenger instance with custom output writer
func NewMessengerWithOutput(output io.Writer) Messenger {
	return &messenger{
		output: output,
	}
}

// ShowProjectNameInference displays information about project name inference
func (m *messenger) ShowProjectNameInference(inferred, directory string) {
	fmt.Fprintf(m.output, "📁 Using directory name as project name: %s\n", inferred)
	if inferred != directory {
		fmt.Fprintf(m.output, "   (sanitized from directory: %s)\n", directory)
	}
}

// ShowDockerWarning displays a warning when Docker is not available
func (m *messenger) ShowDockerWarning() {
	fmt.Fprintf(m.output, "\n⚠️  Docker Warning:\n")
	fmt.Fprintf(m.output, "   Docker is not running or not available on your system.\n")
	fmt.Fprintf(m.output, "   Project files will be generated successfully, but you'll need Docker to run the services.\n")
	fmt.Fprintf(m.output, "\n💡 To install Docker:\n")
	fmt.Fprintf(m.output, "   • Visit: https://docs.docker.com/get-docker/\n")
	fmt.Fprintf(m.output, "   • Install Docker Desktop for your platform\n")
	fmt.Fprintf(m.output, "   • Start Docker and try running: docker --version\n\n")
}

// ShowDirectoryInferenceWarning displays a warning when directory inference has issues
func (m *messenger) ShowDirectoryInferenceWarning(issue, suggestion string) {
	fmt.Fprintf(m.output, "\n⚠️  Directory Name Warning:\n")
	fmt.Fprintf(m.output, "   %s\n", issue)
	fmt.Fprintf(m.output, "   💡 %s\n\n", suggestion)
}

// ShowFallbackNameUsage displays information when a fallback name is used
func (m *messenger) ShowFallbackNameUsage(fallbackName, reason string) {
	fmt.Fprintf(m.output, "\n📝 Using fallback project name: %s\n", fallbackName)
	fmt.Fprintf(m.output, "   Reason: %s\n", reason)
	fmt.Fprintf(m.output, "   💡 Use --name to specify a custom project name\n\n")
}

// ShowErrorRecoveryGuidance displays guidance for error recovery
func (m *messenger) ShowErrorRecoveryGuidance(errorType, guidance string) {
	fmt.Fprintf(m.output, "\n🔧 %s Recovery:\n", errorType)
	fmt.Fprintf(m.output, "   %s\n\n", guidance)
}

// ShowDirectoryAccessError displays specific guidance for directory access issues
func (m *messenger) ShowDirectoryAccessError(path string, errorType string, solutions []string) {
	fmt.Fprintf(m.output, "\n❌ Directory Access Error: %s\n", errorType)
	fmt.Fprintf(m.output, "   Path: %s\n", path)
	fmt.Fprintf(m.output, "\n🔧 Solutions:\n")
	for _, solution := range solutions {
		fmt.Fprintf(m.output, "   • %s\n", solution)
	}
	fmt.Fprintf(m.output, "\n")
}

// ShowDockerRecoverySteps displays comprehensive Docker troubleshooting steps
func (m *messenger) ShowDockerRecoverySteps(issue string, steps []string) {
	fmt.Fprintf(m.output, "\n🐳 Docker Issue: %s\n", issue)
	fmt.Fprintf(m.output, "\n🔧 Recovery Steps:\n")
	for i, step := range steps {
		fmt.Fprintf(m.output, "   %d. %s\n", i+1, step)
	}
	fmt.Fprintf(m.output, "\n💡 After following these steps, try the command again.\n\n")
}

// ShowPermissionError displays permission-specific error guidance
func (m *messenger) ShowPermissionError(path string, platform string) {
	fmt.Fprintf(m.output, "\n🔐 Permission Error\n")
	fmt.Fprintf(m.output, "   Path: %s\n", path)
	
	switch platform {
	case "windows":
		fmt.Fprintf(m.output, "\n🔧 Windows Solutions:\n")
		fmt.Fprintf(m.output, "   • Run terminal as administrator\n")
		fmt.Fprintf(m.output, "   • Right-click folder → Properties → Security → Edit permissions\n")
		fmt.Fprintf(m.output, "   • Try a different directory (e.g., C:\\Users\\[USERNAME]\\Documents)\n")
		fmt.Fprintf(m.output, "   • Disable antivirus temporarily if it's blocking file creation\n")
	case "linux", "darwin":
		fmt.Fprintf(m.output, "\n🔧 Unix/Linux Solutions:\n")
		fmt.Fprintf(m.output, "   • Check ownership: ls -la %s\n", path)
		fmt.Fprintf(m.output, "   • Fix ownership: sudo chown $USER %s\n", path)
		fmt.Fprintf(m.output, "   • Fix permissions: chmod 755 %s\n", path)
		fmt.Fprintf(m.output, "   • Try your home directory: cd ~ && mkdir myproject\n")
	default:
		fmt.Fprintf(m.output, "\n🔧 General Solutions:\n")
		fmt.Fprintf(m.output, "   • Check directory permissions\n")
		fmt.Fprintf(m.output, "   • Try a different directory with write access\n")
		fmt.Fprintf(m.output, "   • Run with appropriate privileges\n")
	}
	fmt.Fprintf(m.output, "\n")
}

// ShowDiskSpaceError displays disk space specific guidance
func (m *messenger) ShowDiskSpaceError(path string, availableSpace string) {
	fmt.Fprintf(m.output, "\n💾 Insufficient Disk Space\n")
	fmt.Fprintf(m.output, "   Path: %s\n", path)
	if availableSpace != "" {
		fmt.Fprintf(m.output, "   Available: %s\n", availableSpace)
	}
	fmt.Fprintf(m.output, "   Required: ~100MB for project generation\n")
	
	fmt.Fprintf(m.output, "\n🔧 Disk Space Solutions:\n")
	fmt.Fprintf(m.output, "   • Free up disk space by deleting unnecessary files\n")
	fmt.Fprintf(m.output, "   • Clean temporary files and caches\n")
	fmt.Fprintf(m.output, "   • Choose a different directory/drive with more space\n")
	fmt.Fprintf(m.output, "   • Remove unused Docker images: docker system prune\n")
	fmt.Fprintf(m.output, "   • Check disk usage: df -h (Linux/Mac) or dir (Windows)\n\n")
}

// ShowGenerationSuccess displays success message after project generation
func (m *messenger) ShowGenerationSuccess(projectName, path string) {
	fmt.Fprintf(m.output, "\n✅ Successfully created microservice project '%s'!\n", projectName)
	if path != "." {
		fmt.Fprintf(m.output, "   📂 Location: %s\n", path)
	}
}

// ShowNextSteps displays next steps for the user based on their configuration
func (m *messenger) ShowNextSteps(config interfaces.CLIConfig) {
	fmt.Fprintf(m.output, "\n📋 Next steps:\n")
	
	// Step 1: Navigate to project directory (if not current directory)
	// If name was inferred, we're already in the right directory
	if !config.InferredName && config.ProjectName != "" {
		fmt.Fprintf(m.output, "   1. cd %s\n", config.ProjectName)
	} else {
		fmt.Fprintf(m.output, "   1. Review the generated files in your current directory\n")
	}
	
	// Step 2: Environment setup
	fmt.Fprintf(m.output, "   2. cp .env.example .env\n")
	fmt.Fprintf(m.output, "   3. Edit .env file with your configuration\n")
	
	// Step 3: Docker commands
	fmt.Fprintf(m.output, "   4. docker-compose up -d\n")
	fmt.Fprintf(m.output, "   5. docker-compose ps\n")
	
	// Show service URLs
	m.showServiceURLs(config)
	
	// Additional helpful information
	fmt.Fprintf(m.output, "\n📚 Helpful commands:\n")
	fmt.Fprintf(m.output, "   • View logs: docker-compose logs -f\n")
	fmt.Fprintf(m.output, "   • Stop services: docker-compose down\n")
	fmt.Fprintf(m.output, "   • Rebuild: docker-compose up --build\n")
	
	fmt.Fprintf(m.output, "\n🚀 Happy coding!\n")
}

// showServiceURLs displays the URLs for each configured service
func (m *messenger) showServiceURLs(config interfaces.CLIConfig) {
	fmt.Fprintf(m.output, "\n🔗 Service URLs:\n")
	
	if config.Backend != "" {
		port := getDefaultPort("backend", config.Backend)
		fmt.Fprintf(m.output, "   Backend (%s): http://localhost:%d\n", config.Backend, port)
	}
	
	if config.Frontend != "" {
		port := getDefaultPort("frontend", config.Frontend)
		fmt.Fprintf(m.output, "   Frontend (%s): http://localhost:%d\n", config.Frontend, port)
	}
	
	if config.Database != "" {
		port := getDefaultPort("database", config.Database)
		fmt.Fprintf(m.output, "   Database (%s): localhost:%d\n", config.Database, port)
	}
}

// getDefaultPort returns the default port for a service type and technology
// This duplicates logic from generator but keeps the feedback package independent
func getDefaultPort(serviceType, technology string) int {
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