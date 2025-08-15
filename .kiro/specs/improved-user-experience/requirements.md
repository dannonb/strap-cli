# Requirements Document

## Introduction

This feature enhances the user experience of the microservice bootstrapper CLI tool by making it more intuitive and removing unnecessary dependencies. Users should be able to create projects in their current directory without specifying a name, and the tool should work without requiring Docker to be running.

## Requirements

### Requirement 1

**User Story:** As a developer, I want to create a project in my current directory without specifying a name, so that I can quickly bootstrap projects with minimal setup.

#### Acceptance Criteria

1. WHEN a user runs the create command without the --name flag THEN the system SHALL use the current directory name as the project name
2. WHEN a user runs the create command in a directory with an existing name THEN the system SHALL use that directory name for the project
3. WHEN the current directory name contains invalid characters THEN the system SHALL sanitize the name to be valid for project naming
4. WHEN a user explicitly provides a --name flag THEN the system SHALL use the provided name instead of the directory name

### Requirement 2

**User Story:** As a developer, I want to use the CLI tool without having Docker running, so that I can generate project scaffolding even when Docker is not available on my system.

#### Acceptance Criteria

1. WHEN Docker is not running THEN the system SHALL still generate all project files successfully
2. WHEN Docker is not available THEN the system SHALL display a warning message about Docker being needed for running the generated project
3. WHEN generating project files THEN the system SHALL NOT require Docker to be running for template processing
4. WHEN Docker is unavailable THEN the system SHALL include clear instructions in the generated README about Docker requirements

### Requirement 3

**User Story:** As a developer, I want clear feedback about what the tool is doing, so that I understand the project creation process and any limitations.

#### Acceptance Criteria

1. WHEN creating a project THEN the system SHALL display the inferred project name to the user
2. WHEN Docker is not available THEN the system SHALL display a clear warning message
3. WHEN project creation is complete THEN the system SHALL display next steps for the user
4. WHEN using directory name as project name THEN the system SHALL inform the user of this behavior

### Requirement 4

**User Story:** As a developer, I want the tool to work in any directory structure, so that I can organize my projects however I prefer.

#### Acceptance Criteria

1. WHEN running in an empty directory THEN the system SHALL create the project structure in that directory
2. WHEN running in a directory with existing files THEN the system SHALL warn about potential conflicts
3. WHEN the directory name is not suitable for a project name THEN the system SHALL provide a fallback naming strategy
4. WHEN creating files THEN the system SHALL respect existing directory permissions and structure