package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// DirectoryService handles directory name inference and project name sanitization
type DirectoryService interface {
	GetCurrentDirectoryName() (string, error)
	SanitizeProjectName(name string) string
	ValidateDirectoryForProject(path string) error
}

// directoryService implements DirectoryService
type directoryService struct{}

// NewDirectoryService creates a new DirectoryService instance
func NewDirectoryService() DirectoryService {
	return &directoryService{}
}

// GetCurrentDirectoryName returns the name of the current working directory
func (d *directoryService) GetCurrentDirectoryName() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		// Enhanced error handling with fallback strategy
		return d.handleCurrentDirectoryAccessError(err)
	}

	dirName := filepath.Base(wd)
	
	// Handle edge cases with specific error reporting
	if dirName == "." || dirName == "/" || dirName == "\\" {
		return d.handleRootDirectoryCase(wd)
	}

	// Handle empty or problematic directory names
	if dirName == "" {
		return d.handleEmptyDirectoryName(wd)
	}

	return dirName, nil
}

// handleCurrentDirectoryAccessError handles cases where we can't access the current directory
func (d *directoryService) handleCurrentDirectoryAccessError(err error) (string, error) {
	fallbackName := "microservice-project"
	
	// Try to provide more specific guidance based on the error
	if strings.Contains(err.Error(), "permission denied") {
		return fallbackName, fmt.Errorf("cannot access current directory due to permission restrictions. "+
			"Using fallback name '%s'. "+
			"🔧 Recovery Steps:\n"+
			"  • Navigate to a directory with proper permissions\n"+
			"  • Use --name to specify a project name explicitly\n"+
			"  • On Windows: Try running as administrator\n"+
			"  • On Linux/Mac: Check directory ownership (ls -la)\n"+
			"Original error: %w", fallbackName, err)
	}
	
	if strings.Contains(err.Error(), "no such file or directory") {
		return fallbackName, fmt.Errorf("current directory no longer exists or is inaccessible. "+
			"Using fallback name '%s'. "+
			"🔧 Recovery Steps:\n"+
			"  • Navigate to a valid directory (cd /path/to/valid/directory)\n"+
			"  • Create a new directory (mkdir myproject && cd myproject)\n"+
			"  • Use --name to specify a project name explicitly\n"+
			"  • Check if you're on a network drive that disconnected\n"+
			"Original error: %w", fallbackName, err)
	}
	
	if strings.Contains(err.Error(), "network") || strings.Contains(err.Error(), "remote") {
		return fallbackName, fmt.Errorf("network error accessing current directory. "+
			"Using fallback name '%s'. "+
			"🔧 Recovery Steps:\n"+
			"  • Check network connection to remote directory\n"+
			"  • Navigate to a local directory\n"+
			"  • Use --name to specify a project name explicitly\n"+
			"  • Ensure network drives are properly mounted\n"+
			"Original error: %w", fallbackName, err)
	}
	
	// Generic fallback with comprehensive guidance
	return fallbackName, fmt.Errorf("unable to determine current directory name. "+
		"Using fallback name '%s'. "+
		"🔧 Recovery Steps:\n"+
		"  • Use --name to specify a project name explicitly\n"+
		"  • Navigate to a different directory\n"+
		"  • Check directory permissions and accessibility\n"+
		"  • Try: mkdir myproject && cd myproject && strap create ...\n"+
		"Original error: %w", fallbackName, err)
}

// handleRootDirectoryCase handles the case where we're in a root directory
func (d *directoryService) handleRootDirectoryCase(path string) (string, error) {
	fallbackName := "microservice-project"
	
	return fallbackName, fmt.Errorf("cannot use root directory '%s' as project name. "+
		"Using fallback name '%s'. "+
		"🔧 Root Directory Solutions:\n"+
		"  • Create a project directory: mkdir myproject && cd myproject\n"+
		"  • Navigate to your home directory: cd ~ && mkdir myproject && cd myproject\n"+
		"  • Use --name to specify a project name explicitly\n"+
		"  • Navigate to a development folder: cd ~/Development\n"+
		"⚠️  Creating projects in root directory is not recommended for security and organization", 
		path, fallbackName)
}

// handleEmptyDirectoryName handles the case where directory name is empty
func (d *directoryService) handleEmptyDirectoryName(path string) (string, error) {
	fallbackName := "microservice-project"
	
	return fallbackName, fmt.Errorf("directory name is empty or invalid for path '%s'. "+
		"Using fallback name '%s'. "+
		"🔧 Empty Directory Name Solutions:\n"+
		"  • Navigate to a properly named directory\n"+
		"  • Create a new directory: mkdir myproject && cd myproject\n"+
		"  • Use --name to specify a project name explicitly\n"+
		"  • Check if you're in a special system directory\n"+
		"💡 Example: mkdir my-microservice && cd my-microservice && strap create --be=fastapi", 
		path, fallbackName)
}

// SanitizeProjectName sanitizes a directory name to be a valid project name
func (d *directoryService) SanitizeProjectName(name string) string {
	if name == "" {
		return "microservice-project"
	}

	original := name
	
	// Convert to lowercase
	sanitized := strings.ToLower(name)

	// Replace spaces and underscores with hyphens
	sanitized = strings.ReplaceAll(sanitized, " ", "-")
	sanitized = strings.ReplaceAll(sanitized, "_", "-")

	// Remove or replace invalid characters (keep only alphanumeric and hyphens)
	reg := regexp.MustCompile(`[^a-z0-9-]`)
	sanitized = reg.ReplaceAllString(sanitized, "")

	// Remove consecutive hyphens
	reg = regexp.MustCompile(`-+`)
	sanitized = reg.ReplaceAllString(sanitized, "-")

	// Remove leading and trailing hyphens
	sanitized = strings.Trim(sanitized, "-")

	// Handle empty result after sanitization
	if sanitized == "" {
		return d.handleEmptySanitizedName(original)
	}

	// Handle reserved names and add suffix if needed
	reservedNames := map[string]bool{
		"con": true, "prn": true, "aux": true, "nul": true,
		"com1": true, "com2": true, "com3": true, "com4": true,
		"com5": true, "com6": true, "com7": true, "com8": true,
		"com9": true, "lpt1": true, "lpt2": true, "lpt3": true,
		"lpt4": true, "lpt5": true, "lpt6": true, "lpt7": true,
		"lpt8": true, "lpt9": true, "node": true, "src": true,
		"test": true, "tests": true, "build": true, "dist": true,
	}

	if reservedNames[sanitized] {
		sanitized = sanitized + "-project"
	}

	// Ensure it doesn't start with a number
	if len(sanitized) > 0 && sanitized[0] >= '0' && sanitized[0] <= '9' {
		sanitized = "project-" + sanitized
	}

	// Limit length to reasonable size
	if len(sanitized) > 50 {
		sanitized = sanitized[:50]
		// Remove trailing hyphen if truncation created one
		sanitized = strings.TrimSuffix(sanitized, "-")
	}

	return sanitized
}

// handleEmptySanitizedName handles the case where sanitization results in an empty string
func (d *directoryService) handleEmptySanitizedName(original string) string {
	fallbackName := "microservice-project"
	
	// Log the issue for debugging (in a real implementation, you might use a proper logger)
	fmt.Printf("⚠️  Warning: Directory name '%s' contains only invalid characters. Using fallback name '%s'.\n", 
		original, fallbackName)
	fmt.Printf("   💡 Suggestion: Use --name to specify a valid project name explicitly.\n")
	
	return fallbackName
}

// ValidateDirectoryForProject checks if a directory is suitable for project creation
func (d *directoryService) ValidateDirectoryForProject(path string) error {
	// Check if path exists
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return d.handleDirectoryNotExist(path)
		}
		return d.handleDirectoryAccessError(path, err)
	}

	// Check if it's actually a directory
	if !info.IsDir() {
		return d.handleNotADirectory(path, info)
	}

	// Check if we have write permissions
	if err := d.validateWritePermissions(path); err != nil {
		return err
	}

	// Check for potential issues with directory structure
	if err := d.validateDirectoryStructure(path); err != nil {
		return err
	}

	return nil
}

// handleDirectoryNotExist handles the case where the target directory doesn't exist
func (d *directoryService) handleDirectoryNotExist(path string) error {
	dirName := filepath.Base(path)
	return fmt.Errorf("directory does not exist: %s. "+
		"🔧 Directory Creation Solutions:\n"+
		"  • Create the directory: mkdir %s\n"+
		"  • Create with parent directories: mkdir -p %s\n"+
		"  • Navigate to an existing directory\n"+
		"  • Use a different path that exists\n"+
		"💡 Quick fix: mkdir %s && cd %s", 
		path, dirName, path, dirName, dirName)
}

// handleDirectoryAccessError handles directory access errors with specific guidance
func (d *directoryService) handleDirectoryAccessError(path string, err error) error {
	if strings.Contains(err.Error(), "permission denied") {
		return fmt.Errorf("permission denied accessing directory %s. "+
			"🔧 Permission Solutions:\n"+
			"  • Check directory permissions: ls -la %s\n"+
			"  • Change permissions: chmod 755 %s (if you own it)\n"+
			"  • Run with appropriate privileges (sudo on Linux/Mac)\n"+
			"  • On Windows: Run as administrator or check folder properties\n"+
			"  • Navigate to a directory you have access to\n"+
			"Original error: %w", path, filepath.Dir(path), path, err)
	}
	
	if strings.Contains(err.Error(), "network") || strings.Contains(err.Error(), "remote") {
		return fmt.Errorf("network error accessing directory %s. "+
			"🔧 Network Directory Solutions:\n"+
			"  • Check network connection\n"+
			"  • Ensure network drive is properly mounted\n"+
			"  • Try accessing the path manually first\n"+
			"  • Use a local directory instead\n"+
			"  • Reconnect to the network share\n"+
			"Original error: %w", path, err)
	}
	
	if strings.Contains(err.Error(), "too many levels") || strings.Contains(err.Error(), "name too long") {
		return fmt.Errorf("path too long or too deeply nested: %s. "+
			"🔧 Path Length Solutions:\n"+
			"  • Use a shorter path\n"+
			"  • Navigate to a directory closer to the root\n"+
			"  • Create project in a simpler directory structure\n"+
			"  • Use --name with a shorter project name\n"+
			"Original error: %w", path, err)
	}
	
	return fmt.Errorf("failed to access directory %s. "+
		"🔧 General Directory Solutions:\n"+
		"  • Verify the path exists: ls %s\n"+
		"  • Check if you have access permissions\n"+
		"  • Try a different directory\n"+
		"  • Ensure the path is not corrupted or locked\n"+
		"Original error: %w", path, path, err)
}

// handleNotADirectory handles the case where the path exists but is not a directory
func (d *directoryService) handleNotADirectory(path string, info os.FileInfo) error {
	fileType := "file"
	if info.Mode()&os.ModeSymlink != 0 {
		fileType = "symbolic link"
	} else if info.Mode()&os.ModeDevice != 0 {
		fileType = "device"
	}
	
	return fmt.Errorf("path '%s' exists but is a %s, not a directory. "+
		"Suggestion: Choose a different path or remove the existing %s", 
		path, fileType, fileType)
}

// validateWritePermissions checks if we have write permissions in the directory
func (d *directoryService) validateWritePermissions(path string) error {
	testFile := filepath.Join(path, ".strap-write-test")
	file, err := os.Create(testFile)
	if err != nil {
		if strings.Contains(err.Error(), "permission denied") {
			return fmt.Errorf("no write permission in directory %s. "+
				"🔧 Write Permission Solutions:\n"+
				"  • Check directory permissions: ls -la %s\n"+
				"  • Change permissions: chmod 755 %s (if you own it)\n"+
				"  • On Windows: Run as administrator or check folder properties\n"+
				"  • On Linux/Mac: sudo chown $USER %s (if needed)\n"+
				"  • Choose a different directory with write access\n"+
				"  • Try your home directory: cd ~ && mkdir myproject\n"+
				"Original error: %w", path, path, path, path, err)
		}
		
		if strings.Contains(err.Error(), "no space left") || strings.Contains(err.Error(), "disk full") {
			return fmt.Errorf("insufficient disk space in directory %s. "+
				"🔧 Disk Space Solutions:\n"+
				"  • Check available space: df -h %s\n"+
				"  • Free up disk space (need ~100MB for project)\n"+
				"  • Clean temporary files and caches\n"+
				"  • Choose a different directory/drive with more space\n"+
				"  • Remove unused files or applications\n"+
				"Original error: %w", path, path, err)
		}
		
		if strings.Contains(err.Error(), "read-only") || strings.Contains(err.Error(), "readonly") {
			return fmt.Errorf("directory %s is read-only. "+
				"🔧 Read-Only Directory Solutions:\n"+
				"  • Choose a writable directory\n"+
				"  • Change directory permissions if you own it\n"+
				"  • On Windows: Right-click → Properties → uncheck Read-only\n"+
				"  • On Linux/Mac: chmod u+w %s\n"+
				"  • Try your home directory: cd ~ && mkdir myproject\n"+
				"Original error: %w", path, path, err)
		}
		
		if strings.Contains(err.Error(), "device full") || strings.Contains(err.Error(), "quota") {
			return fmt.Errorf("storage quota exceeded or device full in directory %s. "+
				"🔧 Storage Quota Solutions:\n"+
				"  • Check disk quota: quota -u $USER\n"+
				"  • Free up space in your quota\n"+
				"  • Contact system administrator about quota limits\n"+
				"  • Use a different directory or storage location\n"+
				"Original error: %w", path, err)
		}
		
		return fmt.Errorf("cannot write to directory %s. "+
			"🔧 General Write Solutions:\n"+
			"  • Verify directory permissions and ownership\n"+
			"  • Check available disk space: df -h\n"+
			"  • Ensure directory is not locked by another process\n"+
			"  • Try a different directory with known write access\n"+
			"  • Example: cd ~ && mkdir myproject && cd myproject\n"+
			"Original error: %w", path, err)
	}
	
	file.Close()
	if removeErr := os.Remove(testFile); removeErr != nil {
		// Log warning but don't fail - the write test succeeded
		fmt.Printf("⚠️  Warning: Could not clean up test file %s: %v\n", testFile, removeErr)
		fmt.Printf("   💡 This may indicate the directory has unusual permissions\n")
	}

	return nil
}

// validateDirectoryStructure performs additional validation on directory structure
func (d *directoryService) validateDirectoryStructure(path string) error {
	// Check if directory is too deeply nested (potential path length issues)
	if len(path) > 200 {
		return fmt.Errorf("directory path is very long (%d characters): %s. "+
			"Suggestion: Use a shorter path to avoid potential issues with file system limits", 
			len(path), path)
	}
	
	// Check for problematic characters in path (excluding colon for Windows drive letters)
	// Only check the directory name part, not the full path
	dirName := filepath.Base(path)
	if strings.ContainsAny(dirName, "<>\"|?*") {
		return fmt.Errorf("directory name contains invalid characters: %s. "+
			"Suggestion: Use a directory name without special characters like < > \" | ? *", dirName)
	}
	
	return nil
}