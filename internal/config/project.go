package config

import "microservice-bootstrapper/internal/interfaces"

// ProjectConfig represents the complete project configuration
type ProjectConfig struct {
	Name     string
	Services []interfaces.ServiceConfig
	Database *interfaces.DatabaseConfig
	Network  NetworkConfig
}

// NetworkConfig represents Docker network configuration
type NetworkConfig struct {
	Name   string
	Driver string
}

// DefaultNetworkConfig returns the default network configuration
func DefaultNetworkConfig() NetworkConfig {
	return NetworkConfig{
		Name:   "microservice-network",
		Driver: "bridge",
	}
}