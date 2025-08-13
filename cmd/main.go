package main

import (
	"fmt"
	"os"

	"microservice-bootstrapper/internal/cli"
)

// Version information set during build
var (
	Version   = "dev"
	Commit    = "unknown"
	BuildDate = "unknown"
)

func main() {
	// Set version information for CLI
	cli.SetVersionInfo(Version, Commit, BuildDate)
	
	if err := cli.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}