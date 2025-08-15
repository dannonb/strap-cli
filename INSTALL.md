# Installation Guide

This guide covers different ways to install the Microservice Bootstrapper CLI tool.

## Prerequisites

### For Installing and Using the CLI Tool
- **No prerequisites** - The CLI tool runs independently

### For Running Generated Projects
- **Docker**: Required for running generated microservices
- **Docker Compose**: Required for orchestrating multiple services
- **Git**: Optional, for version control of generated projects

### Docker Requirements vs Recommendations

**For Project Generation (File Creation)**:
- Docker is **NOT required** to generate project files
- The tool will create all necessary files and configurations
- You'll receive helpful warnings if Docker isn't available
- All templates and boilerplate code are generated successfully

**For Running Generated Projects**:
- Docker is **required** to run the generated microservices
- Docker Compose is **required** for multi-service orchestration
- Without Docker, you can still:
  - View and edit the generated code
  - Modify configurations
  - Set up manual development environments
  - Use the generated code as a starting point

**Recommended Workflow**:
1. Install the CLI tool (no Docker needed)
2. Generate your project structure (works without Docker)
3. Install Docker when you're ready to run the services
4. Use `docker-compose up` to start your microservices

## Installation Methods

### 1. Download Pre-built Binaries (Recommended)

Download the latest release for your platform from the [releases page](https://github.com/dannonb/strap-cli/releases).

#### Linux/macOS
```bash
# Download and extract (replace with actual release URL)
curl -L https://github.com/dannonb/strap-cli/releases/latest/download/strap-linux-amd64.tar.gz | tar xz

# Move to PATH
sudo mv strap-linux-amd64/strap /usr/local/bin/

# Verify installation
strap --version
```

#### Windows
1. Download `strap-windows-amd64.zip` from the releases page
2. Extract the ZIP file
3. Add the extracted directory to your PATH environment variable
4. Open a new command prompt and run `strap --version`

### 2. Build from Source

#### Prerequisites for Building
- Go 1.21 or later
- Make (Linux/macOS) or PowerShell (Windows)
- Git

#### Linux/macOS
```bash
# Clone the repository
git clone https://github.com/dannonb/strap-cli.git
cd strap-cli

# Build and install
make build
sudo make install

# Or build for all platforms
make build-all
```

#### Windows
```powershell
# Clone the repository
git clone https://github.com/dannonb/strap-cli.git
cd strap-cli

# Build for current platform
.\build.ps1 -Target build

# Or build for all platforms
.\build.ps1 -Target build-all
```

### 3. Package Managers

#### Homebrew (macOS/Linux)
```bash
# Coming soon
brew install strap-cli
```

#### Chocolatey (Windows)
```powershell
# Coming soon
choco install strap-cli
```

#### Scoop (Windows)
```powershell
# Coming soon
scoop install strap-cli
```

## Verification

After installation, verify that the tool is working correctly:

```bash
# Check version
strap --version

# View help
strap --help

# Test project creation using directory name
mkdir test-project && cd test-project
strap create --be=fastapi --fe=react --db=postgres

# Or test with explicit name
strap create --be=fastapi --fe=react --db=postgres --name=my-test-project
```

## Platform-Specific Notes

### Linux
- Requires `sudo` for system-wide installation
- Binary is typically installed to `/usr/local/bin/`
- Ensure `/usr/local/bin/` is in your PATH

### macOS
- Same as Linux
- On Apple Silicon Macs, use the `darwin-arm64` binary
- On Intel Macs, use the `darwin-amd64` binary

### Windows
- No admin privileges required for user installation
- Add installation directory to PATH environment variable
- Use PowerShell or Command Prompt
- Windows Defender may flag the binary initially (this is normal for unsigned binaries)

## Troubleshooting

### Command Not Found
- Ensure the binary is in your PATH
- On Linux/macOS, check with `echo $PATH`
- On Windows, check with `echo %PATH%`

### Permission Denied (Linux/macOS)
- Ensure the binary has execute permissions: `chmod +x /usr/local/bin/strap`
- Use `sudo` for system-wide installation

### Docker Issues

**Important**: Docker is NOT required for project generation, only for running generated projects.

**If you see Docker warnings during project generation**:
- This is normal and expected if Docker isn't running
- All project files will still be generated successfully
- You can install and start Docker later when ready to run the services

**For running generated projects**:
- Ensure Docker is installed and running
- Check with `docker --version` and `docker-compose --version`
- On Linux, ensure your user is in the `docker` group

**Docker Installation**:
- Download from [Docker's official site](https://docs.docker.com/get-docker/)
- Follow platform-specific installation instructions
- Start Docker Desktop (Windows/Mac) or Docker service (Linux)

### Windows Antivirus
- Some antivirus software may flag the binary
- Add an exception for the installation directory
- This is common for unsigned binaries

## Updating

### Pre-built Binaries
1. Download the latest release
2. Replace the existing binary
3. Verify with `strap --version`

### Built from Source
```bash
# Pull latest changes
git pull origin main

# Rebuild and reinstall
make clean
make build
sudo make install
```

## Uninstallation

### Linux/macOS
```bash
sudo rm /usr/local/bin/strap
```

### Windows
1. Remove the binary from your installation directory
2. Remove the directory from your PATH environment variable

## Getting Help

- Run `strap --help` for usage information
- Run `strap create --help` for detailed create command help
- Check the [documentation](https://github.com/dannonb/strap-cli/docs)
- Report issues on [GitHub](https://github.com/dannonb/strap-cli/issues)