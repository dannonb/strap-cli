# Test build script to verify the build system works correctly (Windows)
param(
    [switch]$Help
)

# Helper functions
function Write-Info {
    param([string]$Message)
    Write-Host "[INFO] $Message" -ForegroundColor Blue
}

function Write-Success {
    param([string]$Message)
    Write-Host "[SUCCESS] $Message" -ForegroundColor Green
}

function Write-Warning {
    param([string]$Message)
    Write-Host "[WARNING] $Message" -ForegroundColor Yellow
}

function Write-Error {
    param([string]$Message)
    Write-Host "[ERROR] $Message" -ForegroundColor Red
}

function Show-Help {
    Write-Host "Build System Test Script"
    Write-Host ""
    Write-Host "Usage: .\scripts\test-build.ps1 [-Help]"
    Write-Host ""
    Write-Host "Options:"
    Write-Host "  -Help    Show this help message"
    Write-Host ""
    Write-Host "This script tests the build system functionality including:"
    Write-Host "  • Current platform builds"
    Write-Host "  • Multi-platform builds"
    Write-Host "  • Release packaging"
    Write-Host "  • Version information embedding"
    Write-Host "  • Cleanup operations"
}

function Test-CurrentBuild {
    Write-Info "Testing current platform build..."
    
    try {
        .\build.ps1 -Target build
        if ($LASTEXITCODE -eq 0) {
            Write-Success "Current platform build successful"
            
            # Test the binary
            if (Test-Path "build\strap.exe") {
                Write-Info "Testing binary functionality..."
                $output = & "build\strap.exe" --version 2>&1
                if ($LASTEXITCODE -eq 0) {
                    Write-Success "Binary works correctly"
                    return $true
                } else {
                    Write-Error "Binary execution failed"
                    return $false
                }
            } else {
                Write-Error "Binary not found at build\strap.exe"
                return $false
            }
        } else {
            Write-Error "Current platform build failed"
            return $false
        }
    } catch {
        Write-Error "Current platform build failed: $($_.Exception.Message)"
        return $false
    }
}

function Test-MultiplatformBuild {
    Write-Info "Testing multi-platform build..."
    
    try {
        .\build.ps1 -Target build-all
        if ($LASTEXITCODE -eq 0) {
            Write-Success "Multi-platform build successful"
            
            # Check if all expected binaries exist
            $expectedPlatforms = @("linux-amd64", "linux-arm64", "darwin-amd64", "darwin-arm64", "windows-amd64", "windows-arm64")
            $missingPlatforms = @()
            
            foreach ($platform in $expectedPlatforms) {
                $binaryName = "strap"
                if ($platform -match "windows") {
                    $binaryName = "strap.exe"
                }
                
                $binaryPath = "dist\strap-$platform\$binaryName"
                if (-not (Test-Path $binaryPath)) {
                    $missingPlatforms += $platform
                }
            }
            
            if ($missingPlatforms.Count -eq 0) {
                Write-Success "All platform binaries created successfully"
                return $true
            } else {
                Write-Error "Missing binaries for platforms: $($missingPlatforms -join ', ')"
                return $false
            }
        } else {
            Write-Error "Multi-platform build failed"
            return $false
        }
    } catch {
        Write-Error "Multi-platform build failed: $($_.Exception.Message)"
        return $false
    }
}

function Test-ReleasePackaging {
    Write-Info "Testing release packaging..."
    
    try {
        .\build.ps1 -Target release
        if ($LASTEXITCODE -eq 0) {
            Write-Success "Release packaging successful"
            
            # Check if packages were created
            $packages = Get-ChildItem dist\ -Filter "*.tar.gz", "*.zip" -ErrorAction SilentlyContinue
            if ($packages.Count -gt 0) {
                Write-Success "Release packages created: $($packages.Count) files"
                Write-Info "Package list:"
                $packages | Sort-Object Name | ForEach-Object { Write-Host "  $($_.Name)" }
                return $true
            } else {
                Write-Error "No release packages found"
                return $false
            }
        } else {
            Write-Error "Release packaging failed"
            return $false
        }
    } catch {
        Write-Error "Release packaging failed: $($_.Exception.Message)"
        return $false
    }
}

function Test-VersionInfo {
    Write-Info "Testing version information..."
    
    try {
        # Build with custom version
        $testVersion = "v1.0.0-test"
        
        .\build.ps1 -Target build -Version $testVersion
        if ($LASTEXITCODE -eq 0) {
            Write-Success "Build with custom version successful"
            
            # Test version output
            $versionOutput = & "build\strap.exe" --version 2>&1
            if ($versionOutput -match $testVersion) {
                Write-Success "Version information correctly embedded"
                return $true
            } else {
                Write-Warning "Version information may not be correctly embedded"
                Write-Info "Version output: $versionOutput"
                return $false
            }
        } else {
            Write-Error "Build with custom version failed"
            return $false
        }
    } catch {
        Write-Error "Build with custom version failed: $($_.Exception.Message)"
        return $false
    }
}

function Test-Cleanup {
    Write-Info "Testing cleanup..."
    
    try {
        .\build.ps1 -Target clean
        if ($LASTEXITCODE -eq 0) {
            Write-Success "Cleanup successful"
            
            # Check if directories were removed
            $buildExists = Test-Path "build"
            $distExists = Test-Path "dist"
            
            if (-not $buildExists -and -not $distExists) {
                Write-Success "Build artifacts cleaned up"
                return $true
            } else {
                Write-Warning "Some build artifacts may still exist"
                return $false
            }
        } else {
            Write-Error "Cleanup failed"
            return $false
        }
    } catch {
        Write-Error "Cleanup failed: $($_.Exception.Message)"
        return $false
    }
}

function Main {
    if ($Help) {
        Show-Help
        return
    }
    
    Write-Info "Starting build system tests..."
    Write-Host ""
    
    # Clean up first
    try {
        .\build.ps1 -Target clean | Out-Null
    } catch {
        # Ignore cleanup errors
    }
    
    $failedTests = @()
    
    # Run tests
    if (-not (Test-CurrentBuild)) {
        $failedTests += "current_build"
    }
    Write-Host ""
    
    if (-not (Test-MultiplatformBuild)) {
        $failedTests += "multiplatform_build"
    }
    Write-Host ""
    
    if (-not (Test-ReleasePackaging)) {
        $failedTests += "release_packaging"
    }
    Write-Host ""
    
    if (-not (Test-VersionInfo)) {
        $failedTests += "version_info"
    }
    Write-Host ""
    
    if (-not (Test-Cleanup)) {
        $failedTests += "cleanup"
    }
    Write-Host ""
    
    # Summary
    if ($failedTests.Count -eq 0) {
        Write-Success "All build system tests passed! 🎉"
        Write-Host ""
        Write-Info "Build system is ready for production use."
        exit 0
    } else {
        Write-Error "Some tests failed: $($failedTests -join ', ')"
        Write-Host ""
        Write-Info "Please fix the issues before using the build system."
        exit 1
    }
}

# Execute main function
Main