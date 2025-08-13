package version

import (
	"fmt"
	"runtime"
)

var (
	// Version is the current version of the application
	Version = "0.1.0"
	// BuildDate is the date when the binary was built
	BuildDate = "unknown"
	// GitCommit is the git commit hash
	GitCommit = "unknown"
)

// SetBuildInfo updates the version information with build-time values
func SetBuildInfo(version, commit, buildDate string) {
	if version != "" && version != "dev" {
		Version = version
	}
	if commit != "" && commit != "unknown" {
		GitCommit = commit
	}
	if buildDate != "" && buildDate != "unknown" {
		BuildDate = buildDate
	}
}

// Info returns formatted version information
func Info() string {
	return fmt.Sprintf(`🚀 Microservice Bootstrapper v%s

📋 Build Information:
  Version:    %s
  Built:      %s
  Commit:     %s
  Go Version: %s
  Platform:   %s/%s`,
		Version, Version, BuildDate, GitCommit, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}

// Short returns a short version string
func Short() string {
	return fmt.Sprintf("v%s", Version)
}

// Detailed returns detailed version information with additional context
func Detailed() string {
	return fmt.Sprintf(`%s

🔧 System Information:
  Go Version:     %s
  OS:             %s
  Architecture:   %s
  Compiler:       %s
  
📦 Build Details:
  Build Date:     %s
  Git Commit:     %s
  
🌟 Features:
  • Multi-technology project generation
  • Docker and Docker Compose integration
  • Comprehensive boilerplate code
  • Development-ready configuration`,
		Info(), runtime.Version(), runtime.GOOS, runtime.GOARCH, runtime.Compiler, BuildDate, GitCommit)
}