#!/bin/bash

# Test build script to verify the build system works correctly
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Test current platform build
test_current_build() {
    log_info "Testing current platform build..."
    
    if make build; then
        log_success "Current platform build successful"
        
        # Test the binary
        if [ -f "build/strap" ]; then
            log_info "Testing binary functionality..."
            if ./build/strap --version; then
                log_success "Binary works correctly"
            else
                log_error "Binary execution failed"
                return 1
            fi
        else
            log_error "Binary not found at build/strap"
            return 1
        fi
    else
        log_error "Current platform build failed"
        return 1
    fi
}

# Test multi-platform build
test_multiplatform_build() {
    log_info "Testing multi-platform build..."
    
    if make build-all; then
        log_success "Multi-platform build successful"
        
        # Check if all expected binaries exist
        local expected_platforms=("linux-amd64" "linux-arm64" "darwin-amd64" "darwin-arm64" "windows-amd64" "windows-arm64")
        local missing_platforms=()
        
        for platform in "${expected_platforms[@]}"; do
            local binary_name="strap"
            if [[ "$platform" == *"windows"* ]]; then
                binary_name="strap.exe"
            fi
            
            if [ ! -f "dist/strap-$platform/$binary_name" ]; then
                missing_platforms+=("$platform")
            fi
        done
        
        if [ ${#missing_platforms[@]} -eq 0 ]; then
            log_success "All platform binaries created successfully"
        else
            log_error "Missing binaries for platforms: ${missing_platforms[*]}"
            return 1
        fi
    else
        log_error "Multi-platform build failed"
        return 1
    fi
}

# Test release packaging
test_release_packaging() {
    log_info "Testing release packaging..."
    
    if make release; then
        log_success "Release packaging successful"
        
        # Check if packages were created
        local package_count=$(find dist/ -name "*.tar.gz" -o -name "*.zip" | wc -l)
        if [ "$package_count" -gt 0 ]; then
            log_success "Release packages created: $package_count files"
            log_info "Package list:"
            find dist/ -name "*.tar.gz" -o -name "*.zip" | sort
        else
            log_error "No release packages found"
            return 1
        fi
    else
        log_error "Release packaging failed"
        return 1
    fi
}

# Test version information
test_version_info() {
    log_info "Testing version information..."
    
    # Build with custom version
    local test_version="v1.0.0-test"
    local test_commit="abc1234"
    local test_date="2024-01-01T00:00:00Z"
    
    if make build VERSION="$test_version"; then
        log_success "Build with custom version successful"
        
        # Test version output
        local version_output=$(./build/strap --version)
        if echo "$version_output" | grep -q "$test_version"; then
            log_success "Version information correctly embedded"
        else
            log_warning "Version information may not be correctly embedded"
            log_info "Version output: $version_output"
        fi
    else
        log_error "Build with custom version failed"
        return 1
    fi
}

# Test cleanup
test_cleanup() {
    log_info "Testing cleanup..."
    
    if make clean; then
        log_success "Cleanup successful"
        
        # Check if directories were removed
        if [ ! -d "build" ] && [ ! -d "dist" ]; then
            log_success "Build artifacts cleaned up"
        else
            log_warning "Some build artifacts may still exist"
        fi
    else
        log_error "Cleanup failed"
        return 1
    fi
}

# Main test execution
main() {
    log_info "Starting build system tests..."
    echo
    
    # Clean up first
    make clean >/dev/null 2>&1 || true
    
    local failed_tests=()
    
    # Run tests
    if ! test_current_build; then
        failed_tests+=("current_build")
    fi
    echo
    
    if ! test_multiplatform_build; then
        failed_tests+=("multiplatform_build")
    fi
    echo
    
    if ! test_release_packaging; then
        failed_tests+=("release_packaging")
    fi
    echo
    
    if ! test_version_info; then
        failed_tests+=("version_info")
    fi
    echo
    
    if ! test_cleanup; then
        failed_tests+=("cleanup")
    fi
    echo
    
    # Summary
    if [ ${#failed_tests[@]} -eq 0 ]; then
        log_success "All build system tests passed! 🎉"
        echo
        log_info "Build system is ready for production use."
        return 0
    else
        log_error "Some tests failed: ${failed_tests[*]}"
        echo
        log_info "Please fix the issues before using the build system."
        return 1
    fi
}

# Help function
show_help() {
    echo "Build System Test Script"
    echo
    echo "Usage: $0 [OPTIONS]"
    echo
    echo "Options:"
    echo "  -h, --help    Show this help message"
    echo
    echo "This script tests the build system functionality including:"
    echo "  • Current platform builds"
    echo "  • Multi-platform builds"
    echo "  • Release packaging"
    echo "  • Version information embedding"
    echo "  • Cleanup operations"
}

# Parse arguments
case "${1:-}" in
    -h|--help)
        show_help
        exit 0
        ;;
    *)
        main
        ;;
esac