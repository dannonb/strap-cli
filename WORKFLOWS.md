# Workflow Guide

This guide covers the enhanced workflow patterns available in the Microservice Bootstrapper CLI tool, focusing on the improved user experience with directory-based project creation and Docker-optional generation.

## Overview of New Features

### 🎯 Smart Project Naming
- **Automatic**: Uses your current directory name when `--name` is not specified
- **Explicit**: Override with `--name` when you need specific naming
- **Sanitized**: Automatically handles special characters and invalid names

### 🐳 Docker-Optional Generation
- **Generate Anywhere**: Create project files even when Docker isn't running
- **Clear Warnings**: Helpful messages about Docker requirements for running projects
- **Full Functionality**: All templates and configurations are generated successfully

### 💬 Enhanced Feedback
- **Clear Messages**: Know exactly what the tool is doing
- **Helpful Warnings**: Understand Docker requirements vs recommendations
- **Next Steps**: Get guidance on what to do after project creation

## Workflow Patterns

### 1. Quick Start Workflow (Recommended)

The fastest way to get started with a new project:

```bash
# 1. Create and navigate to your project directory
mkdir my-awesome-app && cd my-awesome-app

# 2. Generate your project (uses directory name)
strap create --be=fastapi --fe=react --db=postgres

# 3. Start your services (requires Docker)
docker-compose up
```

**What happens:**
- Project name becomes "my-awesome-app" (from directory name)
- All project files are generated in the current directory
- Docker warnings appear if Docker isn't running (but generation continues)
- Ready-to-run project structure is created

### 2. Directory-First Development

Organize your projects by creating directories first:

```bash
# Create your workspace
mkdir ~/Development/microservices && cd ~/Development/microservices

# Create individual services
mkdir user-service && cd user-service
strap create --be=fastapi --db=postgres
cd ..

mkdir auth-service && cd auth-service  
strap create --be=gin --db=redis
cd ..

mkdir web-client && cd web-client
strap create --fe=react
cd ..

# Each service is properly named and organized
ls -la
# user-service/    (FastAPI + PostgreSQL)
# auth-service/    (Gin + Redis)  
# web-client/      (React)
```

### 3. Explicit Naming Workflow

When you need specific project names regardless of directory structure:

```bash
# Create projects with descriptive names
strap create --be=express --fe=vue --db=mongo --name=ecommerce-api
strap create --be=gin --db=redis --name=cache-service
strap create --fe=angular --name=admin-dashboard

# Useful for:
# - Standardized naming conventions
# - Multiple projects in the same directory
# - Specific naming requirements
```

### 4. Docker-Optional Development

Work without Docker initially, add it later:

```bash
# 1. Generate project without Docker running
mkdir my-api && cd my-api
strap create --be=fastapi --db=postgres
# ⚠️  Warning: Docker not running - project files generated successfully
# ℹ️  Install Docker later to run the generated services

# 2. Develop your application
code .  # Open in your editor
# - Modify the generated code
# - Add your business logic
# - Customize configurations

# 3. Install and start Docker when ready
# Download from https://docs.docker.com/get-docker/
docker --version  # Verify installation

# 4. Run your services
docker-compose up
```

### 5. Rapid Prototyping Workflow

Quickly test different technology combinations:

```bash
# Test different backends
mkdir test-fastapi && cd test-fastapi
strap create --be=fastapi --db=postgres
cd ..

mkdir test-express && cd test-express
strap create --be=express --db=mongo
cd ..

mkdir test-gin && cd test-gin
strap create --be=gin --db=redis
cd ..

# Compare generated structures and choose your preferred stack
```

### 6. Team Development Workflow

Standardized approach for team projects:

```bash
# 1. Team lead creates project structure
mkdir team-project && cd team-project
strap create --be=fastapi --fe=react --db=postgres

# 2. Initialize version control
git init
git add .
git commit -m "Initial project structure"
git remote add origin <repository-url>
git push -u origin main

# 3. Team members clone and start
git clone <repository-url>
cd team-project
docker-compose up

# 4. Development workflow
# - Make changes to code
# - Services auto-reload in development mode
# - Test changes locally
# - Commit and push changes
```

## Directory Name Handling

### Automatic Sanitization

The tool automatically handles various directory name scenarios:

```bash
# Special characters are sanitized
mkdir "My App!" && cd "My App!"
strap create --be=fastapi
# Project name becomes: "my-app"

# Spaces become hyphens
mkdir "user service" && cd "user service"
strap create --be=gin --db=redis
# Project name becomes: "user-service"

# Invalid names get fallbacks
cd /  # Root directory
strap create --be=express
# Project name becomes: "microservice-project" (fallback)
```

### Edge Case Handling

```bash
# Empty or problematic directory names
mkdir "" && cd ""  # If somehow possible
strap create --be=fastapi
# Uses fallback name with helpful message

# Very long directory names
mkdir "this-is-a-very-long-directory-name-that-might-cause-issues"
cd "this-is-a-very-long-directory-name-that-might-cause-issues"
strap create --be=gin
# Truncated and sanitized appropriately
```

## Docker Integration Patterns

### Development Without Docker

```bash
# 1. Generate project
mkdir my-service && cd my-service
strap create --be=fastapi --db=postgres
# ⚠️  Docker not running - files generated successfully

# 2. Manual development setup (alternative to Docker)
# For FastAPI + PostgreSQL:
python -m venv venv
source venv/bin/activate  # Linux/Mac
# venv\Scripts\activate   # Windows
pip install -r backend/requirements.txt

# Install PostgreSQL locally or use cloud database
# Update .env with your database connection
```

### Docker-Ready Development

```bash
# 1. Ensure Docker is running
docker info

# 2. Generate project
mkdir my-service && cd my-service
strap create --be=fastapi --fe=react --db=postgres
# ✅ Docker detected - full validation completed

# 3. Start services immediately
docker-compose up
# All services start successfully
```

### Hybrid Development

```bash
# 1. Generate without Docker
mkdir my-app && cd my-app
strap create --be=express --fe=vue --db=mongo

# 2. Develop frontend locally, backend in Docker
cd frontend
npm install
npm run dev  # Frontend runs locally on port 3000

# In another terminal
docker-compose up backend database  # Only backend services
```

## Error Recovery Patterns

### Directory Access Issues

```bash
# Problem: Permission denied
cd /restricted-directory
strap create --be=fastapi
# ❌ Directory access denied

# Solution 1: Use explicit name
strap create --be=fastapi --name=my-project

# Solution 2: Navigate to accessible directory
cd ~/Development
mkdir my-project && cd my-project
strap create --be=fastapi
```

### Docker Issues

```bash
# Problem: Docker not installed
strap create --be=fastapi --fe=react --db=postgres
# ⚠️  Docker not found - project files generated successfully
# ℹ️  Install Docker from https://docs.docker.com/get-docker/

# Solution: Install Docker later
# 1. Download and install Docker
# 2. Start Docker Desktop
# 3. Run your project: docker-compose up
```

### Directory Conflicts

```bash
# Problem: Directory not empty
cd existing-project
strap create --be=gin --db=redis
# ❌ Directory not empty

# Solution 1: Use force flag
strap create --be=gin --db=redis --force

# Solution 2: Create subdirectory
mkdir v2 && cd v2
strap create --be=gin --db=redis

# Solution 3: Use different name
strap create --be=gin --db=redis --name=gin-service-v2
```

## Best Practices

### Project Organization

```bash
# Organize by domain
mkdir ~/Development/ecommerce && cd ~/Development/ecommerce
mkdir user-service && cd user-service
strap create --be=fastapi --db=postgres
cd ../
mkdir product-service && cd product-service
strap create --be=gin --db=mongo
cd ../
mkdir web-app && cd web-app
strap create --fe=react
```

### Naming Conventions

```bash
# Use descriptive directory names
mkdir user-authentication-service && cd user-authentication-service
strap create --be=fastapi --db=postgres
# Results in project name: "user-authentication-service"

# Or use explicit naming for consistency
strap create --be=fastapi --db=postgres --name=auth-service
```

### Development Workflow

```bash
# 1. Plan your architecture
mkdir microservices-project && cd microservices-project

# 2. Create services incrementally
mkdir api-gateway && cd api-gateway
strap create --be=express
cd ../

mkdir user-service && cd user-service
strap create --be=fastapi --db=postgres
cd ../

# 3. Add frontend when backend is stable
mkdir web-client && cd web-client
strap create --fe=react
```

### Docker Management

```bash
# Start Docker before intensive development
docker info  # Verify Docker is running

# Generate multiple projects efficiently
for service in user-service auth-service notification-service; do
  mkdir $service && cd $service
  strap create --be=fastapi --db=postgres
  cd ..
done

# Start all services
for service in */; do
  cd $service
  docker-compose up -d
  cd ..
done
```

## Troubleshooting

### Common Issues and Solutions

**Issue**: "Directory name inference failed"
```bash
# Solution: Use explicit naming
strap create --be=fastapi --name=my-project
```

**Issue**: "Docker warnings during generation"
```bash
# This is normal - project files are still generated
# Install Docker later when ready to run services
```

**Issue**: "Project name contains invalid characters"
```bash
# The tool automatically sanitizes names
# Check the output message for the sanitized name
# Use --name for explicit control
```

**Issue**: "Permission denied in current directory"
```bash
# Navigate to a writable directory
cd ~/Development
mkdir my-project && cd my-project
strap create --be=fastapi
```

### Getting Help

```bash
# Comprehensive help
strap --help

# Detailed create command help
strap create --help

# Common examples
strap examples

# Workflow guides
strap workflows

# Version information
strap version
```

## Migration from Previous Versions

If you're used to always specifying `--name`, the new workflow is backward compatible:

```bash
# Old way (still works)
strap create --be=fastapi --fe=react --db=postgres --name=my-project

# New way (simpler)
mkdir my-project && cd my-project
strap create --be=fastapi --fe=react --db=postgres
```

Both approaches produce identical results, but the directory-based approach is more intuitive and requires less typing.