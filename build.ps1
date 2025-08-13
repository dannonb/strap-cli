# Microservice Bootstrapper Build Script for Windows
param(
    [string]$Target = "build",
    [string]$Version = "",
    [switch]$Help
)

# Configuration
$BinaryName = "strap"
$BuildDir = "build"
$DistDir = "dist"

# Get version information
if ($Version -eq "") {
    try {
        $Version = git describe --tags --always --dirty 2>$null
        if ($LASTEXITCODE -ne 0) { $Version = "dev" }
    } catch {
        $Version = "dev"
    }
}

try {
    $Commit = git rev-parse --short HEAD 2>$null
    if ($LASTEXITCODE -ne 0) { $Commit = "unknown" }
} catch {
    $Commit = "unknown"
}

$BuildDate = (Get-Date).ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ")

# Build flags
$LdFlags = "-ldflags `"-X main.Version=$Version -X main.Commit=$Commit -X main.BuildDate=$BuildDate -s -w`""

# Supported platforms
$Platforms = @(
    @{OS="linux"; ARCH="amd64"},
    @{OS="linux"; ARCH="arm64"},
    @{OS="darwin"; ARCH="amd64"},
    @{OS="darwin"; ARCH="arm64"},
    @{OS="windows"; ARCH="amd64"},
    @{OS="windows"; ARCH="arm64"}
)

function Show-Help {
    Write-Host "Microservice Bootstrapper Build Script"
    Write-Host ""
    Write-Host "Usage: .\build.ps1 [-Target <target>] [-Version <version>] [-Help]"
    Write-Host ""
    Write-Host "Targets:"
    Write-Host "  build      - Build for current platform"
    Write-Host "  build-all  - Build for all supported platforms"
    Write-Host "  release    - Create release packages for all platforms"
    Write-Host "  test       - Run tests"
    Write-Host "  clean      - Clean build artifacts"
    Write-Host "  version    - Show version information"
    Write-Host ""
    Write-Host "Examples:"
    Write-Host "  .\build.ps1"
    Write-Host "  .\build.ps1 -Target build-all"
    Write-Host "  .\build.ps1 -Target release -Version v1.0.0"
}

function Build-Current {
    Write-Host "Building $BinaryName for current platform..."
    
    if (!(Test-Path $BuildDir)) {
        New-Item -ItemType Directory -Path $BuildDir | Out-Null
    }
    
    $OutputName = "$BinaryName.exe"
    $Command = "go build $LdFlags -o $BuildDir\$OutputName .\cmd"
    
    Write-Host "Executing: $Command"
    Invoke-Expression $Command
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "Build successful: $BuildDir\$OutputName"
    } else {
        Write-Error "Build failed"
        exit 1
    }
}

function Build-All {
    Write-Host "Building $BinaryName for all platforms..."
    
    # Clean first
    Clean-Artifacts
    
    if (!(Test-Path $DistDir)) {
        New-Item -ItemType Directory -Path $DistDir | Out-Null
    }
    
    foreach ($Platform in $Platforms) {
        $OS = $Platform.OS
        $ARCH = $Platform.ARCH
        $OutputName = $BinaryName
        
        if ($OS -eq "windows") {
            $OutputName = "$BinaryName.exe"
        }
        
        $PlatformDir = "$DistDir\$BinaryName-$OS-$ARCH"
        if (!(Test-Path $PlatformDir)) {
            New-Item -ItemType Directory -Path $PlatformDir | Out-Null
        }
        
        Write-Host "Building for $OS/$ARCH..."
        
        $env:GOOS = $OS
        $env:GOARCH = $ARCH
        
        $Command = "go build $LdFlags -o $PlatformDir\$OutputName .\cmd"
        Invoke-Expression $Command
        
        if ($LASTEXITCODE -ne 0) {
            Write-Error "Failed to build for $OS/$ARCH"
            exit 1
        }
    }
    
    # Reset environment variables
    Remove-Item Env:GOOS -ErrorAction SilentlyContinue
    Remove-Item Env:GOARCH -ErrorAction SilentlyContinue
    
    Write-Host "All builds completed successfully!"
}

function Create-Release {
    Write-Host "Creating release packages..."
    
    Build-All
    
    Get-ChildItem $DistDir -Directory | ForEach-Object {
        $PlatformName = $_.Name
        Write-Host "Packaging $PlatformName..."
        
        if ($PlatformName -match "windows") {
            $ZipPath = "$DistDir\$PlatformName.zip"
            Compress-Archive -Path $_.FullName -DestinationPath $ZipPath -Force
        } else {
            $TarPath = "$DistDir\$PlatformName.tar.gz"
            tar -czf $TarPath -C $DistDir $_.Name
        }
    }
    
    Write-Host "Release packages created in $DistDir\"
}

function Run-Tests {
    Write-Host "Running tests..."
    go test -v .\...
}

function Clean-Artifacts {
    Write-Host "Cleaning build artifacts..."
    
    if (Test-Path $BuildDir) {
        Remove-Item -Recurse -Force $BuildDir
    }
    
    if (Test-Path $DistDir) {
        Remove-Item -Recurse -Force $DistDir
    }
}

function Show-Version {
    Write-Host "Version: $Version"
    Write-Host "Commit: $Commit"
    Write-Host "Build Date: $BuildDate"
}

# Main execution
if ($Help) {
    Show-Help
    exit 0
}

switch ($Target.ToLower()) {
    "build" { Build-Current }
    "build-all" { Build-All }
    "release" { Create-Release }
    "test" { Run-Tests }
    "clean" { Clean-Artifacts }
    "version" { Show-Version }
    default {
        Write-Error "Unknown target: $Target"
        Show-Help
        exit 1
    }
}