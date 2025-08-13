package config

import (
	"testing"

	"microservice-bootstrapper/internal/interfaces"
)

func TestProjectConfig(t *testing.T) {
	// Test ProjectConfig struct creation and field access
	config := ProjectConfig{
		Name: "test-project",
		Services: []interfaces.ServiceConfig{
			{
				Type:       "backend",
				Technology: "fastapi",
				Port:       8000,
				Environment: map[string]string{
					"NODE_ENV": "development",
				},
			},
		},
		Database: &interfaces.DatabaseConfig{
			Type: "postgres",
			Port: 5432,
			Volume: "test-postgres-data",
			Environment: map[string]string{
				"POSTGRES_DB": "testdb",
			},
		},
		Network: NetworkConfig{
			Name:   "test-network",
			Driver: "bridge",
		},
	}

	// Test field access
	if config.Name != "test-project" {
		t.Errorf("ProjectConfig.Name = %s, want 'test-project'", config.Name)
	}

	if len(config.Services) != 1 {
		t.Errorf("ProjectConfig.Services length = %d, want 1", len(config.Services))
	}

	if config.Services[0].Type != "backend" {
		t.Errorf("ProjectConfig.Services[0].Type = %s, want 'backend'", config.Services[0].Type)
	}

	if config.Services[0].Technology != "fastapi" {
		t.Errorf("ProjectConfig.Services[0].Technology = %s, want 'fastapi'", config.Services[0].Technology)
	}

	if config.Services[0].Port != 8000 {
		t.Errorf("ProjectConfig.Services[0].Port = %d, want 8000", config.Services[0].Port)
	}

	if config.Database == nil {
		t.Fatal("ProjectConfig.Database should not be nil")
	}

	if config.Database.Type != "postgres" {
		t.Errorf("ProjectConfig.Database.Type = %s, want 'postgres'", config.Database.Type)
	}

	if config.Database.Port != 5432 {
		t.Errorf("ProjectConfig.Database.Port = %d, want 5432", config.Database.Port)
	}

	if config.Network.Name != "test-network" {
		t.Errorf("ProjectConfig.Network.Name = %s, want 'test-network'", config.Network.Name)
	}

	if config.Network.Driver != "bridge" {
		t.Errorf("ProjectConfig.Network.Driver = %s, want 'bridge'", config.Network.Driver)
	}
}

func TestNetworkConfig(t *testing.T) {
	// Test NetworkConfig struct creation and field access
	network := NetworkConfig{
		Name:   "custom-network",
		Driver: "overlay",
	}

	if network.Name != "custom-network" {
		t.Errorf("NetworkConfig.Name = %s, want 'custom-network'", network.Name)
	}

	if network.Driver != "overlay" {
		t.Errorf("NetworkConfig.Driver = %s, want 'overlay'", network.Driver)
	}
}

func TestDefaultNetworkConfig(t *testing.T) {
	// Test DefaultNetworkConfig function
	defaultNetwork := DefaultNetworkConfig()

	if defaultNetwork.Name != "microservice-network" {
		t.Errorf("DefaultNetworkConfig().Name = %s, want 'microservice-network'", defaultNetwork.Name)
	}

	if defaultNetwork.Driver != "bridge" {
		t.Errorf("DefaultNetworkConfig().Driver = %s, want 'bridge'", defaultNetwork.Driver)
	}
}

func TestProjectConfigWithMultipleServices(t *testing.T) {
	// Test ProjectConfig with multiple services
	config := ProjectConfig{
		Name: "multi-service-project",
		Services: []interfaces.ServiceConfig{
			{
				Type:       "backend",
				Technology: "fastapi",
				Port:       8000,
				Environment: map[string]string{
					"PYTHONPATH": "/app",
				},
			},
			{
				Type:       "frontend",
				Technology: "react",
				Port:       3000,
				Environment: map[string]string{
					"NODE_ENV": "development",
				},
			},
		},
		Database: &interfaces.DatabaseConfig{
			Type: "mongo",
			Port: 27017,
			Volume: "mongo-data",
			Environment: map[string]string{
				"MONGO_INITDB_DATABASE": "testdb",
			},
		},
		Network: DefaultNetworkConfig(),
	}

	if len(config.Services) != 2 {
		t.Errorf("ProjectConfig.Services length = %d, want 2", len(config.Services))
	}

	// Test backend service
	backend := config.Services[0]
	if backend.Type != "backend" || backend.Technology != "fastapi" || backend.Port != 8000 {
		t.Errorf("Backend service not configured correctly: %+v", backend)
	}

	// Test frontend service
	frontend := config.Services[1]
	if frontend.Type != "frontend" || frontend.Technology != "react" || frontend.Port != 3000 {
		t.Errorf("Frontend service not configured correctly: %+v", frontend)
	}

	// Test database
	if config.Database.Type != "mongo" || config.Database.Port != 27017 {
		t.Errorf("Database not configured correctly: %+v", config.Database)
	}

	// Test network uses default
	if config.Network.Name != "microservice-network" {
		t.Errorf("Network should use default configuration")
	}
}

func TestProjectConfigWithNilDatabase(t *testing.T) {
	// Test ProjectConfig with nil database (database-less project)
	config := ProjectConfig{
		Name: "frontend-only-project",
		Services: []interfaces.ServiceConfig{
			{
				Type:       "frontend",
				Technology: "vue",
				Port:       3000,
				Environment: map[string]string{
					"VUE_APP_ENV": "development",
				},
			},
		},
		Database: nil,
		Network:  DefaultNetworkConfig(),
	}

	if config.Database != nil {
		t.Error("ProjectConfig.Database should be nil for database-less projects")
	}

	if len(config.Services) != 1 {
		t.Errorf("ProjectConfig.Services length = %d, want 1", len(config.Services))
	}

	if config.Services[0].Type != "frontend" {
		t.Errorf("Service type = %s, want 'frontend'", config.Services[0].Type)
	}
}

func TestProjectConfigEnvironmentVariables(t *testing.T) {
	// Test that environment variables are properly handled
	config := ProjectConfig{
		Name: "env-test-project",
		Services: []interfaces.ServiceConfig{
			{
				Type:       "backend",
				Technology: "gin",
				Port:       8080,
				Environment: map[string]string{
					"GIN_MODE":    "debug",
					"GO_ENV":      "development",
					"PORT":        "8080",
					"CUSTOM_VAR":  "custom_value",
				},
			},
		},
		Database: &interfaces.DatabaseConfig{
			Type: "postgres",
			Port: 5432,
			Volume: "postgres-data",
			Environment: map[string]string{
				"POSTGRES_DB":       "myapp",
				"POSTGRES_USER":     "postgres",
				"POSTGRES_PASSWORD": "postgres",
				"PGDATA":           "/var/lib/postgresql/data/pgdata",
			},
		},
		Network: DefaultNetworkConfig(),
	}

	// Test service environment variables
	serviceEnv := config.Services[0].Environment
	expectedServiceEnv := map[string]string{
		"GIN_MODE":   "debug",
		"GO_ENV":     "development",
		"PORT":       "8080",
		"CUSTOM_VAR": "custom_value",
	}

	for key, expectedValue := range expectedServiceEnv {
		if actualValue, exists := serviceEnv[key]; !exists {
			t.Errorf("Service environment missing key: %s", key)
		} else if actualValue != expectedValue {
			t.Errorf("Service environment[%s] = %s, want %s", key, actualValue, expectedValue)
		}
	}

	// Test database environment variables
	dbEnv := config.Database.Environment
	expectedDbEnv := map[string]string{
		"POSTGRES_DB":       "myapp",
		"POSTGRES_USER":     "postgres",
		"POSTGRES_PASSWORD": "postgres",
		"PGDATA":           "/var/lib/postgresql/data/pgdata",
	}

	for key, expectedValue := range expectedDbEnv {
		if actualValue, exists := dbEnv[key]; !exists {
			t.Errorf("Database environment missing key: %s", key)
		} else if actualValue != expectedValue {
			t.Errorf("Database environment[%s] = %s, want %s", key, actualValue, expectedValue)
		}
	}
}