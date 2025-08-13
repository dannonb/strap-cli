#!/bin/bash

# Release preparation script for Microservice Bootstrapper
set -e

# Configuration
BINARY_NAME="strap"
REPO_URL="https://github.com/your-org/microservice-bootstrapper"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
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

# Check if we're in a git repository
check_git_repo() {
    if ! git rev-parse --git-dir > /dev/null 2>&1; then
        log_error "Not in a git repository"
        exit 1
    fi
}

# Check if working directory is clean
check_clean_working_dir() {
    if ! git diff-index --quiet HEAD --; then
        log_error "Working directory is not clean. Please commit or stash changes."
        exit 1
    fi
}

# Check if we're on main/master branch
check_main_branch() {
    local current_branch=$(git branch --show-current)
    if [[ "$current_branch" != "main" && "$current_branch" != "master" ]]; then
        log_warning "Not on main/master branch (current: $current_branch)"
        read -p "Continue anyway? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    fi
}

# Get version from user input or git tags
get_version() {
    if [[ -n "$1" ]]; then
        VERSION="$1"
    else
        # Get the latest tag
        local latest_tag=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
        log_info "Latest tag: $latest_tag"
        
        echo "Enter new version (e.g., v1.0.0):"
        read -r VERSION
    fi
    
    # Validate version format
    if [[ ! "$VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9]+)?$ ]]; then
        log_error "Invalid version format. Use semantic versioning (e.g., v1.0.0)"
        exit 1
    fi
    
    log_info "Release version: $VERSION"
}

# Check if tag already exists
check_tag_exists() {
    if git tag -l | grep -q "^$VERSION$"; then
        log_error "Tag $VERSION already exists"
        exit 1
    fi
}

# Run tests
run_tests() {
    log_info "Running tests..."
    if ! make test; then
        log_error "Tests failed"
        exit 1
    fi
    log_success "Tests passed"
}

# Run linter
run_linter() {
    log_info "Running linter..."
    if ! make lint; then
        log_error "Linting failed"
        exit 1
    fi
    log_success "Linting passed"
}

# Build for all platforms
build_all() {
    log_info "Building for all platforms..."
    if ! make clean; then
        log_error "Clean failed"
        exit 1
    fi
    
    if ! make build-all VERSION="$VERSION"; then
        log_error "Build failed"
        exit 1
    fi
    log_success "Build completed"
}

# Create release packages
create_packages() {
    log_info "Creating release packages..."
    if ! make release VERSION="$VERSION"; then
        log_error "Package creation failed"
        exit 1
    fi
    log_success "Release packages created"
}

# Generate checksums
generate_checksums() {
    log_info "Generating checksums..."
    cd dist
    
    # Generate SHA256 checksums
    if command -v sha256sum >/dev/null 2>&1; then
        sha256sum *.tar.gz *.zip > checksums.txt
    elif command -v shasum >/dev/null 2>&1; then
        shasum -a 256 *.tar.gz *.zip > checksums.txt
    else
        log_warning "No checksum utility found, skipping checksum generation"
        cd ..
        return
    fi
    
    cd ..
    log_success "Checksums generated"
}

# Create git tag
create_tag() {
    log_info "Creating git tag..."
    
    # Create annotated tag
    git tag -a "$VERSION" -m "Release $VERSION"
    log_success "Tag $VERSION created"
}

# Generate release notes
generate_release_notes() {
    log_info "Generating release notes..."
    
    local previous_tag=$(git describe --tags --abbrev=0 HEAD^ 2>/dev/null || echo "")
    local release_notes_file="dist/RELEASE_NOTES.md"
    
    cat > "$release_notes_file" << EOF
# Release $VERSION

## What's New

<!-- Add release highlights here -->

## Changes

EOF
    
    if [[ -n "$previous_tag" ]]; then
        echo "### Commits since $previous_tag:" >> "$release_notes_file"
        git log --pretty=format:"- %s (%h)" "$previous_tag..HEAD" >> "$release_notes_file"
    else
        echo "### All commits:" >> "$release_notes_file"
        git log --pretty=format:"- %s (%h)" >> "$release_notes_file"
    fi
    
    cat >> "$release_notes_file" << EOF

## Installation

### Download Pre-built Binaries

Choose the appropriate binary for your platform:

- **Linux (x64)**: \`$BINARY_NAME-linux-amd64.tar.gz\`
- **Linux (ARM64)**: \`$BINARY_NAME-linux-arm64.tar.gz\`
- **macOS (Intel)**: \`$BINARY_NAME-darwin-amd64.tar.gz\`
- **macOS (Apple Silicon)**: \`$BINARY_NAME-darwin-arm64.tar.gz\`
- **Windows (x64)**: \`$BINARY_NAME-windows-amd64.zip\`
- **Windows (ARM64)**: \`$BINARY_NAME-windows-arm64.zip\`

### Verification

Verify the download with the provided checksums:

\`\`\`bash
# Download checksums
curl -L $REPO_URL/releases/download/$VERSION/checksums.txt

# Verify (Linux/macOS)
sha256sum -c checksums.txt

# Verify (Windows)
certutil -hashfile $BINARY_NAME-windows-amd64.zip SHA256
\`\`\`

### Installation

See [INSTALL.md](INSTALL.md) for detailed installation instructions.

## Full Changelog

**Full Changelog**: $REPO_URL/compare/$previous_tag...$VERSION
EOF
    
    log_success "Release notes generated: $release_notes_file"
}

# Display summary
display_summary() {
    echo
    log_success "Release preparation completed!"
    echo
    echo "📦 Release artifacts created in dist/:"
    ls -la dist/
    echo
    echo "🏷️  Git tag created: $VERSION"
    echo
    echo "📝 Next steps:"
    echo "   1. Review the release notes in dist/RELEASE_NOTES.md"
    echo "   2. Push the tag: git push origin $VERSION"
    echo "   3. Create a GitHub release with the generated artifacts"
    echo "   4. Update documentation if needed"
    echo
    echo "🚀 GitHub release command:"
    echo "   gh release create $VERSION dist/*.tar.gz dist/*.zip dist/checksums.txt --notes-file dist/RELEASE_NOTES.md"
}

# Main execution
main() {
    local version_arg="$1"
    
    log_info "Starting release preparation for $BINARY_NAME"
    
    # Pre-flight checks
    check_git_repo
    check_clean_working_dir
    check_main_branch
    
    # Get version
    get_version "$version_arg"
    check_tag_exists
    
    # Quality checks
    run_tests
    run_linter
    
    # Build and package
    build_all
    create_packages
    generate_checksums
    
    # Git operations
    create_tag
    
    # Documentation
    generate_release_notes
    
    # Summary
    display_summary
}

# Help function
show_help() {
    echo "Usage: $0 [VERSION]"
    echo
    echo "Prepare a release of the Microservice Bootstrapper CLI"
    echo
    echo "Arguments:"
    echo "  VERSION    Release version (e.g., v1.0.0). If not provided, will prompt."
    echo
    echo "Options:"
    echo "  -h, --help    Show this help message"
    echo
    echo "Examples:"
    echo "  $0 v1.0.0     Prepare release v1.0.0"
    echo "  $0            Interactive mode - will prompt for version"
}

# Parse arguments
case "${1:-}" in
    -h|--help)
        show_help
        exit 0
        ;;
    *)
        main "$1"
        ;;
esac