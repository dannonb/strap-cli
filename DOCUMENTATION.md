# Documentation Index

This document provides an overview of all available documentation for the Microservice Bootstrapper CLI tool.

## Quick Start

- **[README.md](README.md)** - Main project overview, quick start guide, and basic usage
- **[INSTALL.md](INSTALL.md)** - Installation instructions for all platforms
- **[WORKFLOWS.md](WORKFLOWS.md)** - Comprehensive workflow guide for the enhanced user experience

## Key Features

### 🎯 Smart Project Naming
- Automatically uses current directory name when `--name` is not specified
- Handles special characters and invalid names gracefully
- Provides fallback naming strategies for edge cases

### 🐳 Docker-Optional Generation
- Generate project files even when Docker isn't running
- Clear warnings about Docker requirements vs recommendations
- Full project structure creation without Docker dependency

### 💬 Enhanced User Experience
- Intuitive directory-based workflow
- Clear feedback and error messages
- Comprehensive troubleshooting guidance

## Documentation Structure

### User Guides
- **[README.md](README.md)** - Overview and quick start
- **[WORKFLOWS.md](WORKFLOWS.md)** - Detailed workflow patterns and best practices
- **[INSTALL.md](INSTALL.md)** - Installation and setup instructions

### Development Documentation
- **[BUILD_TEST.md](BUILD_TEST.md)** - Build system and testing information
- **Spec Documents** (`.kiro/specs/improved-user-experience/`)
  - `requirements.md` - Feature requirements
  - `design.md` - Technical design document
  - `tasks.md` - Implementation task list

## Command Reference

### Basic Commands
```bash
strap --help          # Show all available commands
strap version         # Display version information
strap examples        # Show common usage examples
strap workflows       # Show development workflow guides
```

### Project Creation
```bash
# Directory-based creation (recommended)
mkdir my-app && cd my-app
strap create --be=fastapi --fe=react --db=postgres

# Explicit naming
strap create --be=gin --db=redis --name=api-service

# Force creation in non-empty directory
strap create --be=express --fe=vue --force
```

## Workflow Patterns

### 1. Quick Start (Directory-Based)
```bash
mkdir my-project && cd my-project
strap create --be=fastapi --fe=react --db=postgres
docker-compose up
```

### 2. Microservices Architecture
```bash
mkdir microservices && cd microservices

mkdir user-service && cd user-service
strap create --be=fastapi --db=postgres
cd ..

mkdir auth-service && cd auth-service
strap create --be=gin --db=redis
cd ..
```

### 3. Docker-Optional Development
```bash
# Generate without Docker
mkdir my-api && cd my-api
strap create --be=fastapi --db=postgres
# ⚠️ Docker warnings (but files are generated)

# Install Docker later and run
docker-compose up
```

## Docker Requirements

### For Project Generation
- **NOT REQUIRED** - All files are generated successfully
- Warnings shown if Docker isn't available
- Full template processing works without Docker

### For Running Generated Projects
- **Docker** - Required for running microservices
- **Docker Compose** - Required for multi-service orchestration
- **Installation** - [Docker's official site](https://docs.docker.com/get-docker/)

## Supported Technologies

### Backends
- **FastAPI** (Python) - Modern, fast web framework
- **Express.js** (Node.js) - Minimal and flexible web framework
- **Gin** (Go) - High-performance HTTP web framework

### Frontends
- **React** - Popular JavaScript library for building user interfaces
- **Vue.js** - Progressive JavaScript framework
- **Angular** - Platform for building mobile and desktop web applications

### Databases
- **PostgreSQL** - Advanced open source relational database
- **MySQL** - Popular open source relational database
- **MongoDB** - Document-oriented NoSQL database
- **Redis** - In-memory data structure store (cache/database)

## Getting Help

### Built-in Help
```bash
strap --help           # General help
strap create --help    # Detailed create command help
strap examples         # Common usage patterns
strap workflows        # Development workflow guides
```

### Documentation
- Read the appropriate guide from the list above
- Check the troubleshooting sections in [WORKFLOWS.md](WORKFLOWS.md)
- Review error messages for specific guidance

### Community
- Report issues on [GitHub](https://github.com/dannonb/strap-cli/issues)
- Check existing issues for solutions
- Contribute improvements via pull requests

## Migration Guide

### From Previous Versions
The new version is backward compatible:

```bash
# Old way (still works)
strap create --be=fastapi --name=my-project

# New way (simpler)
mkdir my-project && cd my-project
strap create --be=fastapi
```

### Key Changes
- `--name` parameter is now optional
- Directory name is used automatically when `--name` is not specified
- Docker is not required for project generation (only for running)
- Enhanced error messages and user feedback
- Improved workflow patterns for better developer experience

## Best Practices

### Project Organization
- Use descriptive directory names (they become project names)
- Organize related services in a common parent directory
- Follow consistent naming conventions across your projects

### Development Workflow
- Start with directory creation and navigation
- Generate project structure first
- Install Docker when ready to run services
- Use version control for generated projects

### Docker Management
- Install Docker when ready to run services
- Use `docker-compose up` to start all services
- Leverage Docker for consistent development environments
- Consider Docker Desktop for easier management on Windows/Mac

---

For the most up-to-date information, always refer to the built-in help commands and the latest documentation in this repository.