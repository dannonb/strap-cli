# Microservice Bootstrapper

A CLI tool written in Go that generates complete microservice project structures with Docker orchestration. Create projects quickly using your current directory name or specify a custom name.

## Quick Start

```bash
# Create a project using your current directory name
mkdir my-awesome-app && cd my-awesome-app
strap create --be=fastapi --fe=react --db=postgres

# Or specify a custom project name
strap create --be=express --fe=vue --db=mongo --name=my-api

# Start your services
docker-compose up
```

## Key Features

- **Smart Project Naming**: Automatically uses your current directory name when `--name` is not specified
- **Docker Optional**: Generate project files even when Docker isn't running (with helpful warnings)
- **Multiple Technologies**: Support for FastAPI, Express, Gin, React, Vue, Angular, and various databases
- **Complete Setup**: Includes Docker configuration, boilerplate code, and documentation

## Installation

See [INSTALL.md](INSTALL.md) for detailed installation instructions.

## Usage Patterns

### Directory-Based Project Creation (Recommended)

The most intuitive way to create projects - just navigate to your desired directory:

```bash
# Full-stack web application
mkdir webapp && cd webapp
strap create --be=fastapi --fe=react --db=postgres
docker-compose up

# API service
mkdir api-service && cd api-service
strap create --be=gin --db=redis
docker-compose up

# Frontend application
mkdir my-frontend && cd my-frontend
strap create --fe=vue
docker-compose up
```

### Explicit Project Naming

When you need specific project names regardless of directory:

```bash
# Create with custom name
strap create --be=express --fe=angular --db=mysql --name=enterprise-app

# Multiple services with descriptive names
strap create --be=fastapi --db=postgres --name=user-service
strap create --be=gin --db=redis --name=auth-service
```

### Microservices Architecture

Create multiple related services in organized directories:

```bash
mkdir microservices && cd microservices

# User management service
mkdir user-service && cd user-service
strap create --be=fastapi --db=postgres
cd ..

# Authentication service  
mkdir auth-service && cd auth-service
strap create --be=gin --db=redis
cd ..

# Web client
mkdir web-client && cd web-client
strap create --fe=react
```

## Docker Requirements

### For Project Generation
- **Docker is NOT required** to generate project files
- The tool will create all necessary files and configurations
- You'll receive helpful warnings if Docker isn't available

### For Running Generated Projects
- **Docker is required** to run the generated microservices
- **Docker Compose** is needed for multi-service orchestration
- Install from [Docker's official site](https://docs.docker.com/get-docker/)

### Docker Status Handling
```bash
# Works without Docker running - generates files with warnings
strap create --be=fastapi --fe=react --db=postgres

# Optimal experience with Docker running - full validation
docker info  # Verify Docker is running
strap create --be=fastapi --fe=react --db=postgres
```

## Supported Technologies

- **Backends**: FastAPI (Python), Express.js (Node.js), Gin (Go)
- **Frontends**: React, Vue.js, Angular
- **Databases**: MongoDB, PostgreSQL, MySQL, Redis

## Command Reference

```bash
# Show all available commands and options
strap --help

# Detailed help for project creation
strap create --help

# Show common usage examples
strap examples

# Show development workflow guides
strap workflows

# Display version information
strap version
```

## Project Structure

```
strap-cli/
├── cmd/                    # Application entry points
│   └── main.go            # Main application entry point
├── internal/              # Private application code
│   ├── cli/               # CLI command handling
│   ├── config/            # Configuration and directory services
│   ├── feedback/          # User messaging and feedback
│   ├── generator/         # Project generation logic
│   ├── interfaces/        # Core interfaces
│   ├── template/          # Template processing
│   ├── validation/        # Input validation
│   └── version/           # Version information
├── pkg/                   # Public library code
│   └── errors/            # Custom error types
├── scripts/               # Build and release scripts
├── test/                  # Integration tests
├── go.mod                 # Go module definition
└── README.md             # This file
```

## Development

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Create release packages
make release
```

## Getting Help

- Run `strap --help` for usage information
- Run `strap create --help` for detailed create command help
- Run `strap examples` for common usage patterns
- Run `strap workflows` for development workflow guides
- Check the [documentation](https://github.com/dannonb/strap-cli/docs)
- Report issues on [GitHub](https://github.com/dannonb/strap-cli/issues)