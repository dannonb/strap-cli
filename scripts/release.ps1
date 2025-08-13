# Release preparation script for Microservice Bootstrapper (Windows)
param(
    [string]$Version = "",
    [switch]$Help
)

# Configuration
$BinaryName = "strap"
$RepoUrl = "https://github.com/your-org/microservice-bootstrapper"

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
    Write-Host "Release Preparation Script for Microservice Bootstrapper"
    Write-Host ""
    Write-Host "Usage: .\scripts\release.ps1 [-Version <version>] [-Help]"
    Write-Host ""
    Write-Host "Parameters:"
    Write-Host "  -Version    Release version (e.g., v1.0.0). If not provided, will prompt."
    Write-Host "  -Help       Show this help message"
    Write-Host ""
    Write-Host "Examples:"
    Write-Host "  .\scripts\release.ps1 -Version v1.0.0"
    Write-Host "  .\scripts\release.ps1"
}

function Test-GitRepository {
    try {
        git rev-parse --git-dir | Out-Null
        return $true
    } catch {
        Write-Error "Not in a git repository"
        return $false
    }
}

function Test-CleanWorkingDirectory {
    try {
        $status = git status --porcelain
        if ($status) {
            Write-Error "Working directory is not clean. Please commit or stash changes."
            return $false
        }
        return $true
    } catch {
        Write-Error "Failed to check git status"
        return $false
    }
}

function Test-MainBranch {
    try {
        $currentBranch = git branch --show-current
        if ($currentBranch -notin @("main", "master")) {
            Write-Warning "Not on main/master branch (current: $currentBranch)"
            $response = Read-Host "Continue anyway? (y/N)"
            return $response -match "^[Yy]$"
        }
        return $true
    } catch {
        Write-Error "Failed to check current branch"
        return $false
    }
}

function Get-ReleaseVersion {
    param([string]$InputVersion)
    
    if ($InputVersion) {
        return $InputVersion
    }
    
    try {
        $latestTag = git describe --tags --abbrev=0 2>$null
        if ($latestTag) {
            Write-Info "Latest tag: $latestTag"
        }
    } catch {
        Write-Info "No previous tags found"
    }
    
    $version = Read-Host "Enter new version (e.g., v1.0.0)"
    return $version
}

function Test-VersionFormat {
    param([string]$Version)
    
    if ($Version -notmatch "^v\d+\.\d+\.\d+(-[a-zA-Z0-9]+)?$") {
        Write-Error "Invalid version format. Use semantic versioning (e.g., v1.0.0)"
        return $false
    }
    return $true
}

function Test-TagExists {
    param([string]$Version)
    
    try {
        $existingTag = git tag -l $Version
        if ($existingTag) {
            Write-Error "Tag $Version already exists"
            return $false
        }
        return $true
    } catch {
        Write-Error "Failed to check existing tags"
        return $false
    }
}

function Invoke-Tests {
    Write-Info "Running tests..."
    try {
        .\build.ps1 -Target test
        if ($LASTEXITCODE -ne 0) {
            Write-Error "Tests failed"
            return $false
        }
        Write-Success "Tests passed"
        return $true
    } catch {
        Write-Error "Failed to run tests"
        return $false
    }
}

function Invoke-Build {
    param([string]$Version)
    
    Write-Info "Building for all platforms..."
    try {
        .\build.ps1 -Target clean
        .\build.ps1 -Target build-all -Version $Version
        if ($LASTEXITCODE -ne 0) {
            Write-Error "Build failed"
            return $false
        }
        Write-Success "Build completed"
        return $true
    } catch {
        Write-Error "Failed to build"
        return $false
    }
}

function New-ReleasePackages {
    param([string]$Version)
    
    Write-Info "Creating release packages..."
    try {
        .\build.ps1 -Target release -Version $Version
        if ($LASTEXITCODE -ne 0) {
            Write-Error "Package creation failed"
            return $false
        }
        Write-Success "Release packages created"
        return $true
    } catch {
        Write-Error "Failed to create packages"
        return $false
    }
}

function New-Checksums {
    Write-Info "Generating checksums..."
    try {
        Push-Location dist
        
        $files = Get-ChildItem -Filter "*.tar.gz", "*.zip"
        $checksums = @()
        
        foreach ($file in $files) {
            $hash = Get-FileHash -Path $file.Name -Algorithm SHA256
            $checksums += "$($hash.Hash.ToLower())  $($file.Name)"
        }
        
        $checksums | Out-File -FilePath "checksums.txt" -Encoding UTF8
        
        Pop-Location
        Write-Success "Checksums generated"
        return $true
    } catch {
        Pop-Location
        Write-Error "Failed to generate checksums"
        return $false
    }
}

function New-GitTag {
    param([string]$Version)
    
    Write-Info "Creating git tag..."
    try {
        git tag -a $Version -m "Release $Version"
        Write-Success "Tag $Version created"
        return $true
    } catch {
        Write-Error "Failed to create git tag"
        return $false
    }
}

function New-ReleaseNotes {
    param([string]$Version)
    
    Write-Info "Generating release notes..."
    try {
        $previousTag = ""
        try {
            $previousTag = git describe --tags --abbrev=0 HEAD^ 2>$null
        } catch {
            # No previous tag
        }
        
        $releaseNotesPath = "dist\RELEASE_NOTES.md"
        
        $content = @"
# Release $Version

## What's New

<!-- Add release highlights here -->

## Changes

"@
        
        if ($previousTag) {
            $content += "`n### Commits since $previousTag:`n"
            $commits = git log --pretty=format:"- %s (%h)" "$previousTag..HEAD"
            $content += $commits -join "`n"
        } else {
            $content += "`n### All commits:`n"
            $commits = git log --pretty=format:"- %s (%h)"
            $content += $commits -join "`n"
        }
        
        $content += @"

## Installation

### Download Pre-built Binaries

Choose the appropriate binary for your platform:

- **Linux (x64)**: ``$BinaryName-linux-amd64.tar.gz``
- **Linux (ARM64)**: ``$BinaryName-linux-arm64.tar.gz``
- **macOS (Intel)**: ``$BinaryName-darwin-amd64.tar.gz``
- **macOS (Apple Silicon)**: ``$BinaryName-darwin-arm64.tar.gz``
- **Windows (x64)**: ``$BinaryName-windows-amd64.zip``
- **Windows (ARM64)**: ``$BinaryName-windows-arm64.zip``

### Verification

Verify the download with the provided checksums:

``````bash
# Download checksums
curl -L $RepoUrl/releases/download/$Version/checksums.txt

# Verify (Windows)
certutil -hashfile $BinaryName-windows-amd64.zip SHA256
``````

### Installation

See [INSTALL.md](INSTALL.md) for detailed installation instructions.

## Full Changelog

**Full Changelog**: $RepoUrl/compare/$previousTag...$Version
"@
        
        $content | Out-File -FilePath $releaseNotesPath -Encoding UTF8
        Write-Success "Release notes generated: $releaseNotesPath"
        return $true
    } catch {
        Write-Error "Failed to generate release notes"
        return $false
    }
}

function Show-Summary {
    param([string]$Version)
    
    Write-Host ""
    Write-Success "Release preparation completed!"
    Write-Host ""
    Write-Host "📦 Release artifacts created in dist/:"
    Get-ChildItem dist\ | Format-Table Name, Length, LastWriteTime
    Write-Host ""
    Write-Host "🏷️  Git tag created: $Version"
    Write-Host ""
    Write-Host "📝 Next steps:"
    Write-Host "   1. Review the release notes in dist\RELEASE_NOTES.md"
    Write-Host "   2. Push the tag: git push origin $Version"
    Write-Host "   3. Create a GitHub release with the generated artifacts"
    Write-Host "   4. Update documentation if needed"
    Write-Host ""
    Write-Host "🚀 GitHub release command:"
    Write-Host "   gh release create $Version dist\*.tar.gz dist\*.zip dist\checksums.txt --notes-file dist\RELEASE_NOTES.md"
}

# Main execution
function Main {
    param([string]$InputVersion)
    
    if ($Help) {
        Show-Help
        return
    }
    
    Write-Info "Starting release preparation for $BinaryName"
    
    # Pre-flight checks
    if (-not (Test-GitRepository)) { return }
    if (-not (Test-CleanWorkingDirectory)) { return }
    if (-not (Test-MainBranch)) { return }
    
    # Get and validate version
    $releaseVersion = Get-ReleaseVersion $InputVersion
    if (-not (Test-VersionFormat $releaseVersion)) { return }
    if (-not (Test-TagExists $releaseVersion)) { return }
    
    Write-Info "Release version: $releaseVersion"
    
    # Quality checks
    if (-not (Invoke-Tests)) { return }
    
    # Build and package
    if (-not (Invoke-Build $releaseVersion)) { return }
    if (-not (New-ReleasePackages $releaseVersion)) { return }
    if (-not (New-Checksums)) { return }
    
    # Git operations
    if (-not (New-GitTag $releaseVersion)) { return }
    
    # Documentation
    if (-not (New-ReleaseNotes $releaseVersion)) { return }
    
    # Summary
    Show-Summary $releaseVersion
}

# Execute main function
Main $Version