package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"microservice-bootstrapper/internal/generator"
	"microservice-bootstrapper/internal/interfaces"
)

func TestGenerationPerformance(t *testing.T) {
	// Test that project generation completes within reasonable time limits
	tests := []struct {
		name        string
		config      interfaces.CLIConfig
		maxDuration time.Duration
	}{
		{
			name: "simple backend project",
			config: interfaces.CLIConfig{
				Backend:     "fastapi",
				ProjectName: "perf-test-1",
				Force:       true,
			},
			maxDuration: 5 * time.Second,
		},
		{
			name: "simple frontend project",
			config: interfaces.CLIConfig{
				Frontend:    "react",
				ProjectName: "perf-test-2",
				Force:       true,
			},
			maxDuration: 5 * time.Second,
		},
		{
			name: "full stack project",
			config: interfaces.CLIConfig{
				Backend:     "gin",
				Frontend:    "vue",
				Database:    "postgres",
				ProjectName: "perf-test-3",
				Force:       true,
			},
			maxDuration: 10 * time.Second,
		},
		{
			name: "complex full stack project",
			config: interfaces.CLIConfig{
				Backend:     "express",
				Frontend:    "angular",
				Database:    "mongo",
				ProjectName: "perf-test-4",
				Force:       true,
			},
			maxDuration: 10 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory for this test
			tempDir := t.TempDir()
			originalDir, _ := os.Getwd()
			defer os.Chdir(originalDir)

			// Change to temp directory
			os.Chdir(tempDir)

			// Measure generation time
			start := time.Now()
			
			gen := generator.NewGenerator()
			err := gen.Generate(tt.config)
			
			duration := time.Since(start)

			if err != nil {
				t.Fatalf("Failed to generate project: %v", err)
			}

			if duration > tt.maxDuration {
				t.Errorf("Generation took too long: %v (max: %v)", duration, tt.maxDuration)
			}

			t.Logf("Generation completed in %v", duration)

			// Verify project was actually created
			projectPath := filepath.Join(tempDir, tt.config.ProjectName)
			if _, err := os.Stat(projectPath); os.IsNotExist(err) {
				t.Error("Project directory was not created")
			}
		})
	}
}

func BenchmarkProjectGeneration(b *testing.B) {
	// Benchmark different project generation scenarios
	configs := []interfaces.CLIConfig{
		{
			Backend:     "fastapi",
			ProjectName: "bench-fastapi",
			Force:       true,
		},
		{
			Frontend:    "react",
			ProjectName: "bench-react",
			Force:       true,
		},
		{
			Backend:     "gin",
			Frontend:    "vue",
			Database:    "postgres",
			ProjectName: "bench-fullstack",
			Force:       true,
		},
	}

	for _, config := range configs {
		b.Run(config.ProjectName, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				// Create temporary directory for this benchmark
				tempDir := b.TempDir()
				originalDir, _ := os.Getwd()
				
				// Change to temp directory
				os.Chdir(tempDir)

				// Generate project
				gen := generator.NewGenerator()
				err := gen.Generate(config)
				if err != nil {
					b.Fatalf("Failed to generate project: %v", err)
				}

				// Restore original directory
				os.Chdir(originalDir)
			}
		})
	}
}

func TestConcurrentGeneration(t *testing.T) {
	// Test that multiple projects can be generated concurrently without issues
	// Note: This test may have cleanup issues on Windows due to file locking
	const numConcurrent = 3 // Reduced for Windows compatibility
	
	// Create channels for synchronization
	done := make(chan error, numConcurrent)
	
	// Start concurrent generations
	for i := 0; i < numConcurrent; i++ {
		go func(index int) {
			// Create temporary directory for this goroutine
			tempDir := t.TempDir()
			originalDir, _ := os.Getwd()
			defer func() {
				// Ignore errors on directory change back - Windows file locking issue
				_ = os.Chdir(originalDir)
			}()

			// Change to temp directory
			os.Chdir(tempDir)

			// Create unique config for each goroutine
			config := interfaces.CLIConfig{
				Backend:     "fastapi",
				ProjectName: fmt.Sprintf("concurrent-test-%d", index),
				Force:       true,
			}

			// Generate project
			gen := generator.NewGenerator()
			err := gen.Generate(config)
			done <- err
		}(i)
	}

	// Wait for all generations to complete
	for i := 0; i < numConcurrent; i++ {
		err := <-done
		if err != nil {
			t.Errorf("Concurrent generation %d failed: %v", i, err)
		}
	}
	
	// Note: Temp directory cleanup may fail on Windows due to file locking
	// This is a known issue with Windows and doesn't affect the actual functionality
}

func TestMemoryUsage(t *testing.T) {
	// Test that project generation doesn't consume excessive memory
	// This is a basic test - in production you might use more sophisticated memory profiling
	
	config := interfaces.CLIConfig{
		Backend:     "express",
		Frontend:    "angular",
		Database:    "mysql",
		ProjectName: "memory-test",
		Force:       true,
	}

	// Create temporary directory
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	// Change to temp directory
	os.Chdir(tempDir)

	// Generate multiple projects to test memory usage
	const numProjects = 10
	
	for i := 0; i < numProjects; i++ {
		gen := generator.NewGenerator()
		err := gen.Generate(config)
		if err != nil {
			t.Fatalf("Failed to generate project %d: %v", i, err)
		}
		
		// Clean up the generated project to avoid disk space issues
		projectPath := filepath.Join(tempDir, config.ProjectName)
		os.RemoveAll(projectPath)
	}

	// If we get here without running out of memory, the test passes
	t.Log("Memory usage test completed successfully")
}

func TestLargeProjectGeneration(t *testing.T) {
	// Test generation of projects with all possible services
	config := interfaces.CLIConfig{
		Backend:     "express",
		Frontend:    "angular",
		Database:    "postgres",
		ProjectName: "large-project-test",
		Force:       true,
	}

	// Create temporary directory
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	// Change to temp directory
	os.Chdir(tempDir)

	// Measure generation time for large project
	start := time.Now()
	
	gen := generator.NewGenerator()
	err := gen.Generate(config)
	
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("Failed to generate large project: %v", err)
	}

	// Large projects should still complete within reasonable time
	maxDuration := 15 * time.Second
	if duration > maxDuration {
		t.Errorf("Large project generation took too long: %v (max: %v)", duration, maxDuration)
	}

	t.Logf("Large project generation completed in %v", duration)

	// Verify all expected components were created
	projectPath := filepath.Join(tempDir, config.ProjectName)
	
	expectedFiles := []string{
		"docker-compose.yml",
		"README.md",
		".env.example",
		".gitignore",
		"backend/Dockerfile",
		"backend/server.js",
		"backend/package.json",
		"frontend/Dockerfile",
		"frontend/package.json",
		"frontend/angular.json",
	}

	for _, file := range expectedFiles {
		filePath := filepath.Join(projectPath, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("Expected file %s does not exist in large project", file)
		}
	}
}