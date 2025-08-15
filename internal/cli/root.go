package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"microservice-bootstrapper/internal/version"
)

var rootCmd = &cobra.Command{
	Use:   "strap",
	Short: "Microservice Bootstrapper - Generate microservice projects with Docker",
	Long: `Microservice Bootstrapper is a CLI tool that generates complete microservice 
project structures with backend services, frontend applications, databases, 
and Docker orchestration based on simple command-line flags.

🚀 SUPPORTED TECHNOLOGIES:
  Backends:  FastAPI (Python), Express.js (Node.js), Gin (Go)
  Frontends: React, Vue.js, Angular
  Databases: MongoDB, PostgreSQL, MySQL, Redis

📦 WHAT YOU GET:
  • Complete project structure with proper directory layout
  • Docker and Docker Compose configuration for all services
  • Boilerplate code with basic API endpoints and components
  • Environment configuration and documentation
  • Ready-to-run development environment

🎯 COMMON WORKFLOWS:
  1. Full-stack development: Backend + Frontend + Database
  2. API development: Backend + Database only
  3. Frontend development: Frontend only
  4. Microservice development: Multiple backends with shared database

Use 'strap create --help' for detailed usage examples and options.`,
	Example: `  # Full-stack project with FastAPI, React, and PostgreSQL
  mkdir my-app && cd my-app
  strap create --be=fastapi --fe=react --db=postgres

  # Backend API with Gin and Redis cache (explicit name)
  strap create --be=gin --db=redis --name=api-service

  # Frontend-only React application (uses directory name)
  mkdir my-frontend && cd my-frontend
  strap create --fe=react

  # Multiple services (run in separate directories)
  mkdir user-service && cd user-service
  strap create --be=fastapi --db=postgres
  cd ../auth-service && mkdir auth-service
  strap create --be=gin --db=redis`,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long: `Display detailed version information including build details,
supported technologies, and system information.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version.Info())
		fmt.Println()
		fmt.Println("📦 Supported Technologies:")
		fmt.Println("  Backends:  fastapi, express, gin")
		fmt.Println("  Frontends: react, vue, angular")
		fmt.Println("  Databases: mongo, postgres, mysql, redis")
		fmt.Println()
		fmt.Println("🔗 More information: https://github.com/dannonb/strap-cli")
	},
}

var helpCmd = &cobra.Command{
	Use:   "help [command]",
	Short: "Help about any command",
	Long: `Help provides help for any command in the application.
Simply type strap help [path to command] for full details.

Available help topics:
  strap help create     - Detailed create command help
  strap help examples   - Common usage examples
  strap help workflows  - Development workflow guides`,
}

var examplesCmd = &cobra.Command{
	Use:   "examples",
	Short: "Show common usage examples",
	Long: `Common usage examples for different development scenarios.

🌐 FULL-STACK APPLICATIONS:
  # Modern web app with Python backend (uses directory name)
  mkdir webapp && cd webapp
  strap create --be=fastapi --fe=react --db=postgres

  # Node.js API with Vue.js frontend (explicit name)
  strap create --be=express --fe=vue --db=mongo --name=fullstack

  # Go backend with Angular frontend (uses directory name)
  mkdir enterprise && cd enterprise
  strap create --be=gin --fe=angular --db=mysql

🔌 API SERVICES:
  # REST API with caching (explicit name)
  strap create --be=fastapi --db=redis --name=api-service

  # GraphQL API with database (uses directory name)
  mkdir graphql-api && cd graphql-api
  strap create --be=express --db=postgres

  # High-performance Go API (uses directory name)
  mkdir fast-api && cd fast-api
  strap create --be=gin --db=mongo

🎨 FRONTEND APPLICATIONS:
  # React SPA (uses directory name)
  mkdir react-app && cd react-app
  strap create --fe=react

  # Vue.js application (uses directory name)
  mkdir vue-app && cd vue-app
  strap create --fe=vue

  # Angular enterprise app (uses directory name)
  mkdir angular-app && cd angular-app
  strap create --fe=angular

🏗️ MICROSERVICES ARCHITECTURE:
  # Create multiple services in separate directories
  mkdir microservices && cd microservices
  
  # User management service (uses directory name)
  mkdir user-service && cd user-service
  strap create --be=fastapi --db=postgres
  
  # Authentication service (uses directory name)
  cd ../
  mkdir auth-service && cd auth-service
  strap create --be=gin --db=redis
  
  # Notification service (uses directory name)
  cd ../
  mkdir notification-service && cd notification-service
  strap create --be=express --db=mongo
  
  # Web client (uses directory name)
  cd ../
  mkdir web-client && cd web-client
  strap create --fe=react`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var workflowsCmd = &cobra.Command{
	Use:   "workflows",
	Short: "Show development workflow guides",
	Long: `Development workflow guides for different scenarios.

🚀 QUICK START WORKFLOW:
  1. Create and navigate to your project directory:
     mkdir myapp && cd myapp

  2. Create your project (uses directory name):
     strap create --be=fastapi --fe=react --db=postgres

  3. Start all services:
     docker-compose up

  4. Access your application:
     - Frontend: http://localhost:3000
     - Backend API: http://localhost:8000
     - Database: localhost:5432

🔄 DEVELOPMENT WORKFLOW:
  1. Make changes to your code
  2. Services auto-reload in development mode
  3. Test your changes
  4. Commit and push to version control

📦 PRODUCTION DEPLOYMENT:
  1. Build production images:
     docker-compose -f docker-compose.prod.yml build

  2. Deploy to your infrastructure:
     docker-compose -f docker-compose.prod.yml up -d

🧪 TESTING WORKFLOW:
  1. Run backend tests:
     cd backend && npm test  # or pytest, go test

  2. Run frontend tests:
     cd frontend && npm test

  3. Integration tests:
     docker-compose -f docker-compose.test.yml up --abort-on-container-exit

🔧 CUSTOMIZATION WORKFLOW:
  1. Modify generated Dockerfiles as needed
  2. Update docker-compose.yml for your requirements
  3. Add environment variables to .env files
  4. Extend the boilerplate code with your business logic`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// SetVersionInfo sets the version information from build-time variables
func SetVersionInfo(ver, commit, buildDate string) {
	version.SetBuildInfo(ver, commit, buildDate)
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(examplesCmd)
	rootCmd.AddCommand(workflowsCmd)
	
	// Configure help settings
	rootCmd.SetHelpCommand(helpCmd)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	
	// Customize help template for better formatting
	rootCmd.SetHelpTemplate(`{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`)
}