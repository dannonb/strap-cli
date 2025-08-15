package interfaces

// CLIConfig represents the configuration parsed from CLI flags
type CLIConfig struct {
	Backend      string
	Frontend     string
	Database     string
	ProjectName  string
	Force        bool
	InferredName bool   // Track if name was inferred from directory
	WorkingDir   string // Store working directory path
}

// ServiceConfig represents configuration for a single service
type ServiceConfig struct {
	Type        string
	Technology  string
	Port        int
	Environment map[string]string
}

// PrerequisiteLevel defines different levels of prerequisite checking
type PrerequisiteLevel int

const (
	// PrerequisiteGeneration checks only what's needed for file generation
	PrerequisiteGeneration PrerequisiteLevel = iota
	// PrerequisiteExecution checks what's needed for running the generated project
	PrerequisiteExecution
)

// ProjectGenerator orchestrates the project creation process
type ProjectGenerator interface {
	Generate(config CLIConfig) error
	GenerateWithMessenger(config CLIConfig, messenger Messenger) error
	ValidateConfig(config CLIConfig) error
	CheckPrerequisites() error
	CheckPrerequisitesWithLevel(level PrerequisiteLevel) error
}

// Messenger provides consistent user communication throughout the project generation process
type Messenger interface {
	ShowProjectNameInference(inferred, directory string)
	ShowDockerWarning()
	ShowGenerationSuccess(projectName, path string)
	ShowNextSteps(config CLIConfig)
}