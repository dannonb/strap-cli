package interfaces

// CLIConfig represents the configuration parsed from CLI flags
type CLIConfig struct {
	Backend     string
	Frontend    string
	Database    string
	ProjectName string
	Force       bool
}

// ServiceConfig represents configuration for a single service
type ServiceConfig struct {
	Type        string
	Technology  string
	Port        int
	Environment map[string]string
}

// ProjectGenerator orchestrates the project creation process
type ProjectGenerator interface {
	Generate(config CLIConfig) error
	ValidateConfig(config CLIConfig) error
	CheckPrerequisites() error
}