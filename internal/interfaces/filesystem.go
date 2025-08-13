package interfaces

// FileSystemManager handles all file system operations
type FileSystemManager interface {
	CreateDirectory(path string) error
	WriteFile(path string, content []byte) error
	FileExists(path string) bool
	IsDirectoryEmpty(path string) bool
	CleanupOnFailure(basePath string) error
}