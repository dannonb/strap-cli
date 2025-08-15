# Design Document

## Overview

This design enhances the microservice bootstrapper CLI tool to provide a more intuitive user experience by automatically inferring project names from the current directory and removing the hard dependency on Docker being actively running during project generation.

The key improvements focus on:
1. **Smart project naming**: Automatically use the current directory name when --name is not provided
2. **Docker-optional generation**: Generate project files without requiring Docker to be running
3. **Enhanced user feedback**: Clear messaging about inferred names and Docker requirements
4. **Flexible directory handling**: Support project creation in any directory structure

## Architecture

### Current Architecture
The current system requires explicit project naming and assumes Docker is available during generation. The CLI validates all prerequisites upfront, including Docker availability.

### Enhanced Architecture
The enhanced system will:
- Add directory name inference logic to the CLI layer
- Modify prerequisite checking to be more granular
- Separate file generation from Docker validation
- Enhance user feedback throughout the process

## Components and Interfaces

### 1. Directory Name Inference Service

**Location**: `internal/config/directory.go`

```go
type DirectoryService interface {
    GetCurrentDirectoryName() (string, error)
    SanitizeProjectName(name string) string
    ValidateDirectoryForProject(path string) error
}

type directoryService struct{}

func (d *directoryService) GetCurrentDirectoryName() (string, error)
func (d *directoryService) SanitizeProjectName(name string) string
func (d *directoryService) ValidateDirectoryForProject(path string) error
```

**Responsibilities**:
- Extract current directory name from filesystem
- Sanitize directory names to be valid project names
- Check if directory is suitable for project creation
- Handle edge cases (root directory, special characters, etc.)

### 2. Enhanced CLI Configuration

**Modifications to**: `internal/cli/create.go`

```go
// Enhanced project name resolution logic
func resolveProjectName(providedName string) (string, error) {
    if providedName != "" {
        return providedName, nil
    }
    
    dirService := config.NewDirectoryService()
    dirName, err := dirService.GetCurrentDirectoryName()
    if err != nil {
        return "", err
    }
    
    sanitized := dirService.SanitizeProjectName(dirName)
    return sanitized, nil
}
```

**Changes**:
- Modify flag validation to make --name optional
- Add project name resolution logic
- Update help text to reflect new behavior
- Add user feedback about inferred names

### 3. Prerequisite Checking Enhancement

**Modifications to**: `internal/generator/generator.go`

```go
type PrerequisiteLevel int

const (
    PrerequisiteGeneration PrerequisiteLevel = iota
    PrerequisiteExecution
)

func (g *Generator) CheckPrerequisites(level PrerequisiteLevel) error {
    switch level {
    case PrerequisiteGeneration:
        return g.checkGenerationPrerequisites()
    case PrerequisiteExecution:
        return g.checkExecutionPrerequisites()
    }
}

func (g *Generator) checkGenerationPrerequisites() error {
    // Only check filesystem permissions, directory access
    // No Docker requirement
}

func (g *Generator) checkExecutionPrerequisites() error {
    // Check Docker availability for running generated projects
}
```

**Changes**:
- Split prerequisite checking into generation vs execution phases
- Remove Docker requirement from generation phase
- Add warnings instead of errors for missing Docker

### 4. Enhanced User Feedback System

**New component**: `internal/feedback/messenger.go`

```go
type Messenger interface {
    ShowProjectNameInference(inferred, directory string)
    ShowDockerWarning()
    ShowGenerationSuccess(projectName, path string)
    ShowNextSteps(config interfaces.CLIConfig)
}

type messenger struct {
    output io.Writer
}

func (m *messenger) ShowProjectNameInference(inferred, directory string)
func (m *messenger) ShowDockerWarning()
func (m *messenger) ShowGenerationSuccess(projectName, path string)
func (m *messenger) ShowNextSteps(config interfaces.CLIConfig)
```

**Responsibilities**:
- Provide consistent, user-friendly messaging
- Show inferred project names clearly
- Display appropriate warnings and next steps
- Format output for better readability

## Data Models

### Enhanced CLIConfig

**Modifications to**: `internal/interfaces/generator.go`

```go
type CLIConfig struct {
    Backend     string
    Frontend    string
    Database    string
    ProjectName string
    Force       bool
    // New fields
    InferredName bool   // Track if name was inferred
    WorkingDir   string // Store working directory path
}
```

### Directory Information

**New model**: `internal/config/models.go`

```go
type DirectoryInfo struct {
    Path         string
    Name         string
    IsEmpty      bool
    HasConflicts bool
    Permissions  os.FileMode
}

type ProjectNameInfo struct {
    Original   string
    Sanitized  string
    Source     string // "provided", "inferred", "fallback"
}
```

## Error Handling

### Directory Name Inference Errors

1. **Root directory**: Provide fallback name "microservice-project"
2. **Invalid characters**: Sanitize by replacing with hyphens
3. **Reserved names**: Add suffix like "-project"
4. **Empty directory name**: Use fallback naming strategy

### Docker Availability Handling

1. **Docker not running**: Show warning, continue with generation
2. **Docker not installed**: Show informational message about installation
3. **Docker permission issues**: Provide guidance on setup

### File Generation Errors

1. **Permission denied**: Clear error message with solutions
2. **Directory conflicts**: Enhanced messaging about --force option
3. **Disk space issues**: Informative error handling

## Testing Strategy

### Unit Tests

1. **Directory Service Tests**:
   - Test directory name extraction
   - Test name sanitization logic
   - Test edge cases (root, special chars, etc.)

2. **CLI Configuration Tests**:
   - Test project name resolution
   - Test flag validation with optional name
   - Test user feedback messaging

3. **Prerequisite Checking Tests**:
   - Test generation vs execution prerequisites
   - Test Docker availability detection
   - Test warning vs error scenarios

### Integration Tests

1. **End-to-End Project Creation**:
   - Test creation with inferred names
   - Test creation without Docker running
   - Test creation in various directory structures

2. **User Experience Tests**:
   - Test feedback messaging
   - Test error handling flows
   - Test next steps guidance

### Manual Testing Scenarios

1. **Directory Name Inference**:
   - Create projects in directories with various names
   - Test in root directory, nested directories
   - Test with special characters in directory names

2. **Docker-Optional Generation**:
   - Generate projects with Docker stopped
   - Generate projects without Docker installed
   - Verify all files are created correctly

## Implementation Phases

### Phase 1: Directory Name Inference
- Implement DirectoryService
- Modify CLI to make --name optional
- Add project name resolution logic
- Update help documentation

### Phase 2: Docker-Optional Generation
- Split prerequisite checking
- Remove Docker requirement from generation
- Add warning messages for missing Docker
- Update error handling

### Phase 3: Enhanced User Feedback
- Implement Messenger component
- Add informative output during generation
- Improve error messages and guidance
- Add next steps information

### Phase 4: Testing and Polish
- Comprehensive test coverage
- User experience testing
- Documentation updates
- Performance optimization