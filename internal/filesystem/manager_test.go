package filesystem

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewManager(t *testing.T) {
	manager := NewManager()
	if manager == nil {
		t.Fatal("NewManager() returned nil")
	}

	// Verify it implements the interface
	_, ok := manager.(*Manager)
	if !ok {
		t.Fatal("NewManager() did not return a *Manager")
	}
}

func TestManager_CreateDirectory(t *testing.T) {
	manager := NewManager().(*Manager)
	tempDir := t.TempDir()

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "create single directory",
			path:    filepath.Join(tempDir, "test-dir"),
			wantErr: false,
		},
		{
			name:    "create nested directories",
			path:    filepath.Join(tempDir, "nested", "deep", "directory"),
			wantErr: false,
		},
		{
			name:    "create directory that already exists",
			path:    tempDir, // temp dir already exists
			wantErr: false,
		},
		{
			name:    "empty path should fail",
			path:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.CreateDirectory(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateDirectory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify directory was created
				if _, err := os.Stat(tt.path); os.IsNotExist(err) {
					t.Errorf("Directory %s was not created", tt.path)
				}
			}
		})
	}
}

func TestManager_WriteFile(t *testing.T) {
	manager := NewManager().(*Manager)
	tempDir := t.TempDir()

	tests := []struct {
		name    string
		path    string
		content []byte
		wantErr bool
	}{
		{
			name:    "write file to existing directory",
			path:    filepath.Join(tempDir, "test.txt"),
			content: []byte("test content"),
			wantErr: false,
		},
		{
			name:    "write file to non-existing directory",
			path:    filepath.Join(tempDir, "new-dir", "test.txt"),
			content: []byte("test content"),
			wantErr: false,
		},
		{
			name:    "write empty file",
			path:    filepath.Join(tempDir, "empty.txt"),
			content: []byte(""),
			wantErr: false,
		},
		{
			name:    "empty path should fail",
			path:    "",
			content: []byte("content"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.WriteFile(tt.path, tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify file was created with correct content
				content, err := os.ReadFile(tt.path)
				if err != nil {
					t.Errorf("Failed to read written file: %v", err)
					return
				}

				if string(content) != string(tt.content) {
					t.Errorf("File content = %s, want %s", string(content), string(tt.content))
				}
			}
		})
	}
}

func TestManager_FileExists(t *testing.T) {
	manager := NewManager().(*Manager)
	tempDir := t.TempDir()

	// Create a test file
	testFile := filepath.Join(tempDir, "exists.txt")
	err := os.WriteFile(testFile, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "existing file",
			path: testFile,
			want: true,
		},
		{
			name: "existing directory",
			path: tempDir,
			want: true,
		},
		{
			name: "non-existing file",
			path: filepath.Join(tempDir, "does-not-exist.txt"),
			want: false,
		},
		{
			name: "empty path",
			path: "",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := manager.FileExists(tt.path)
			if got != tt.want {
				t.Errorf("FileExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_IsDirectoryEmpty(t *testing.T) {
	manager := NewManager().(*Manager)
	tempDir := t.TempDir()

	// Create an empty directory
	emptyDir := filepath.Join(tempDir, "empty")
	err := os.Mkdir(emptyDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create empty directory: %v", err)
	}

	// Create a non-empty directory
	nonEmptyDir := filepath.Join(tempDir, "non-empty")
	err = os.Mkdir(nonEmptyDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create non-empty directory: %v", err)
	}
	err = os.WriteFile(filepath.Join(nonEmptyDir, "file.txt"), []byte("content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create file in non-empty directory: %v", err)
	}

	// Create a test file (not a directory)
	testFile := filepath.Join(tempDir, "file.txt")
	err = os.WriteFile(testFile, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "empty directory",
			path: emptyDir,
			want: true,
		},
		{
			name: "non-empty directory",
			path: nonEmptyDir,
			want: false,
		},
		{
			name: "file (not directory)",
			path: testFile,
			want: false,
		},
		{
			name: "non-existing path",
			path: filepath.Join(tempDir, "does-not-exist"),
			want: false,
		},
		{
			name: "empty path",
			path: "",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := manager.IsDirectoryEmpty(tt.path)
			if got != tt.want {
				t.Errorf("IsDirectoryEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_CleanupOnFailure(t *testing.T) {
	manager := NewManager().(*Manager)
	tempDir := t.TempDir()

	// Create some files and directories using the manager
	testDir := filepath.Join(tempDir, "test-project")
	testFile1 := filepath.Join(testDir, "file1.txt")
	testFile2 := filepath.Join(testDir, "subdir", "file2.txt")

	// Create directories and files
	err := manager.CreateDirectory(testDir)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	err = manager.WriteFile(testFile1, []byte("content1"))
	if err != nil {
		t.Fatalf("Failed to write test file1: %v", err)
	}

	err = manager.WriteFile(testFile2, []byte("content2"))
	if err != nil {
		t.Fatalf("Failed to write test file2: %v", err)
	}

	// Verify files exist before cleanup
	if !manager.FileExists(testFile1) {
		t.Fatal("Test file1 should exist before cleanup")
	}
	if !manager.FileExists(testFile2) {
		t.Fatal("Test file2 should exist before cleanup")
	}

	tests := []struct {
		name     string
		basePath string
		wantErr  bool
	}{
		{
			name:     "cleanup within base path",
			basePath: tempDir,
			wantErr:  false,
		},
		{
			name:     "empty base path should fail",
			basePath: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.CleanupOnFailure(tt.basePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("CleanupOnFailure() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify files were cleaned up
				if manager.FileExists(testFile1) {
					t.Error("Test file1 should have been cleaned up")
				}
				if manager.FileExists(testFile2) {
					t.Error("Test file2 should have been cleaned up")
				}
			}
		})
	}
}

func TestManager_CleanupOnFailure_OutsideBasePath(t *testing.T) {
	manager := NewManager().(*Manager)
	tempDir := t.TempDir()

	// Create a file outside the base path we'll use for cleanup
	outsideDir := t.TempDir() // Different temp directory
	outsideFile := filepath.Join(outsideDir, "outside.txt")

	// Manually add to created paths to simulate creation
	manager.createdPaths = append(manager.createdPaths, outsideFile)

	// Create the file
	err := os.WriteFile(outsideFile, []byte("content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create outside file: %v", err)
	}

	// Create a file inside the base path
	insideFile := filepath.Join(tempDir, "inside.txt")
	err = manager.WriteFile(insideFile, []byte("content"))
	if err != nil {
		t.Fatalf("Failed to create inside file: %v", err)
	}

	// Cleanup with tempDir as base path
	err = manager.CleanupOnFailure(tempDir)
	if err != nil {
		t.Errorf("CleanupOnFailure() error = %v", err)
	}

	// File outside base path should still exist
	if !manager.FileExists(outsideFile) {
		t.Error("File outside base path should not have been cleaned up")
	}

	// File inside base path should be cleaned up
	if manager.FileExists(insideFile) {
		t.Error("File inside base path should have been cleaned up")
	}
}

func TestManager_Integration(t *testing.T) {
	manager := NewManager().(*Manager)
	tempDir := t.TempDir()

	// Test a complete workflow: create directories, write files, then cleanup
	projectDir := filepath.Join(tempDir, "my-project")
	backendDir := filepath.Join(projectDir, "backend")
	frontendDir := filepath.Join(projectDir, "frontend")

	// Create project structure
	err := manager.CreateDirectory(backendDir)
	if err != nil {
		t.Fatalf("Failed to create backend directory: %v", err)
	}

	err = manager.CreateDirectory(frontendDir)
	if err != nil {
		t.Fatalf("Failed to create frontend directory: %v", err)
	}

	// Write some files
	err = manager.WriteFile(filepath.Join(projectDir, "docker-compose.yml"), []byte("version: '3'"))
	if err != nil {
		t.Fatalf("Failed to write docker-compose.yml: %v", err)
	}

	err = manager.WriteFile(filepath.Join(backendDir, "main.go"), []byte("package main"))
	if err != nil {
		t.Fatalf("Failed to write main.go: %v", err)
	}

	err = manager.WriteFile(filepath.Join(frontendDir, "package.json"), []byte("{}"))
	if err != nil {
		t.Fatalf("Failed to write package.json: %v", err)
	}

	// Verify everything was created
	if !manager.FileExists(projectDir) {
		t.Error("Project directory should exist")
	}
	if !manager.FileExists(filepath.Join(projectDir, "docker-compose.yml")) {
		t.Error("docker-compose.yml should exist")
	}
	if !manager.FileExists(filepath.Join(backendDir, "main.go")) {
		t.Error("main.go should exist")
	}
	if !manager.FileExists(filepath.Join(frontendDir, "package.json")) {
		t.Error("package.json should exist")
	}

	// Test cleanup
	err = manager.CleanupOnFailure(tempDir)
	if err != nil {
		t.Errorf("CleanupOnFailure() error = %v", err)
	}

	// Verify everything was cleaned up
	if manager.FileExists(filepath.Join(projectDir, "docker-compose.yml")) {
		t.Error("docker-compose.yml should have been cleaned up")
	}
	if manager.FileExists(filepath.Join(backendDir, "main.go")) {
		t.Error("main.go should have been cleaned up")
	}
	if manager.FileExists(filepath.Join(frontendDir, "package.json")) {
		t.Error("package.json should have been cleaned up")
	}
}