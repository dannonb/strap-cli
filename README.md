# Microservice Bootstrapper

A CLI tool written in Go that generates complete microservice project structures with Docker orchestration.

## Project Structure

```
strap-cli/
├── cmd/                    # Application entry points
│   └── main.go            # Main application entry point
├── internal/              # Private application code
│   ├── cli/               # CLI command handling
│   │   ├── root.go        # Root command and CLI setup
│   │   └── create.go      # Create command implementation
│   ├── config/            # Configuration structures
│   │   └── project.go     # Project configuration models
│   ├── interfaces/        # Core interfaces
│   │   ├── generator.go   # Project generator interface
│   │   ├── template.go    # Template engine interface
│   │   └── filesystem.go  # File system manager interface
│   └── version/           # Version information
│       └── version.go     # Version constants and info
├── pkg/                   # Public library code
│   └── errors/            # Custom error types
│       └── errors.go      # Error type definitions
├── go.mod                 # Go module definition
└── README.md             # This file
```

## Building

```bash
go build -o strap cmd/main.go
```

## Usage

```bash
# Show help
./strap --help

# Show version
./strap version

# Create a new project (implementation in progress)
./strap create --be=fastapi --fe=react --db=postgres
```

## Development Status

This project is currently under development. Task 1 (project structure and core interfaces) has been completed.