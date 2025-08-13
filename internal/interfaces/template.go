package interfaces

// Template represents a template with its files and configuration
type Template struct {
	Files       map[string]string // filename -> template content
	Directories []string
	Variables   map[string]interface{}
}

// TemplateData represents data passed to templates during processing
type TemplateData struct {
	ProjectName string
	Services    []ServiceConfig
	Backend     *ServiceTemplateData
	Frontend    *ServiceTemplateData
	Database    *DatabaseConfig
	Ports       PortConfig
	Environment map[string]string
}

// ServiceTemplateData represents service-specific template data
type ServiceTemplateData struct {
	Type        string
	Technology  string
	Port        int
	Environment map[string]string
}

// PortConfig represents port configuration for all services
type PortConfig struct {
	Backend  int
	Frontend int
	Database int
}

// DatabaseConfig represents database service configuration
type DatabaseConfig struct {
	Type        string
	Port        int
	Volume      string
	Environment map[string]string
}

// TemplateEngine manages and processes template files
type TemplateEngine interface {
	ProcessTemplate(templateName string, data interface{}) ([]byte, error)
	GetTemplate(service, technology string) (Template, error)
	GetProjectTemplate() (Template, error)
	GenerateDockerCompose(config CLIConfig) ([]byte, error)
	ValidateComposeConfiguration(config CLIConfig) error
}