#!/bin/bash
# CalVer Version Management Script
# Format: YYYY.MM.MICRO
# Usage:
#   ./scripts/version.sh              # Show current version
#   ./scripts/version.sh bump         # Bump micro version
#   ./scripts/version.sh set <version> # Set specific version

set -e

VERSION_FILE="VERSION"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"
VERSION_PATH="$ROOT_DIR/$VERSION_FILE"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get current version
get_current_version() {
    if [ -f "$VERSION_PATH" ]; then
        cat "$VERSION_PATH"
    else
        echo "0.0.0"
    fi
}

# Validate CalVer format (YYYY.MM.MICRO)
validate_version() {
    local version=$1
    if [[ ! $version =~ ^[0-9]{4}\.[0-9]{1,2}\.[0-9]+$ ]]; then
        echo -e "${RED}Error: Invalid CalVer format. Expected YYYY.MM.MICRO (e.g., 2026.01.0)${NC}"
        exit 1
    fi
}

# Calculate next version
bump_version() {
    local current_version=$(get_current_version)
    local year=$(date +%Y)
    local month=$(date +%-m)  # No leading zero

    # Parse current version
    IFS='.' read -r current_year current_month current_micro <<< "$current_version"

    # If same month/year, increment micro
    if [ "$current_year" = "$year" ] && [ "$current_month" = "$month" ]; then
        new_micro=$((current_micro + 1))
        new_version="${year}.${month}.${new_micro}"
    else
        # New month, reset micro to 0
        new_version="${year}.${month}.0"
    fi

    echo "$new_version"
}

# Set version
set_version() {
    local version=$1
    validate_version "$version"
    echo "$version" > "$VERSION_PATH"
    echo -e "${GREEN}Version set to: $version${NC}"
}

# Show current version
show_version() {
    local version=$(get_current_version)
    echo -e "${GREEN}Current version: $version${NC}"
}

# Update version in files
update_version_in_files() {
    local version=$1

    echo -e "${YELLOW}Updating version in project files...${NC}"

    # Update package.json if it exists
    if [ -f "$ROOT_DIR/package.json" ]; then
        if command -v jq &> /dev/null; then
            jq --arg version "$version" '.version = $version' "$ROOT_DIR/package.json" > "$ROOT_DIR/package.json.tmp"
            mv "$ROOT_DIR/package.json.tmp" "$ROOT_DIR/package.json"
            echo -e "${GREEN}✓ Updated package.json${NC}"
        else
            echo -e "${YELLOW}⚠ jq not installed, skipping package.json update${NC}"
        fi
    fi

    # Update pubspec.yaml if it exists (Flutter)
    if [ -f "$ROOT_DIR/mobile/pubspec.yaml" ]; then
        # Extract build number from version (YYYY.MM.MICRO -> YYYYMMMICRO as build number)
        IFS='.' read -r year month micro <<< "$version"
        build_number="${year}$(printf "%02d" $month)$(printf "%03d" $micro)"

        sed -i.bak "s/^version:.*/version: $version+$build_number/" "$ROOT_DIR/mobile/pubspec.yaml"
        rm -f "$ROOT_DIR/mobile/pubspec.yaml.bak"
        echo -e "${GREEN}✓ Updated mobile/pubspec.yaml${NC}"
    fi

    # Update ARCHITECTURE.md
    if [ -f "$ROOT_DIR/docs/ARCHITECTURE.md" ]; then
        sed -i.bak "s/^\*\*Version:\*\*.*/\*\*Version:\*\* $version/" "$ROOT_DIR/docs/ARCHITECTURE.md"
        sed -i.bak "s/^**Document Version**:.*/\*\*Document Version\*\*: $version/" "$ROOT_DIR/docs/ARCHITECTURE.md"
        rm -f "$ROOT_DIR/docs/ARCHITECTURE.md.bak"
        echo -e "${GREEN}✓ Updated docs/ARCHITECTURE.md${NC}"
    fi

    # Update openapi.yaml
    if [ -f "$ROOT_DIR/docs/openapi.yaml" ]; then
        sed -i.bak "s/version:.*/version: $version/" "$ROOT_DIR/docs/openapi.yaml"
        rm -f "$ROOT_DIR/docs/openapi.yaml.bak"
        echo -e "${GREEN}✓ Updated docs/openapi.yaml${NC}"
    fi
}

# Main script logic
case "${1:-show}" in
    show)
        show_version
        ;;
    bump)
        new_version=$(bump_version)
        echo -e "${YELLOW}Bumping version: $(get_current_version) → $new_version${NC}"
        set_version "$new_version"
        update_version_in_files "$new_version"
        echo -e "${GREEN}✓ Version bumped successfully!${NC}"
        ;;
    set)
        if [ -z "$2" ]; then
            echo -e "${RED}Error: Please provide a version (e.g., 2026.01.0)${NC}"
            exit 1
        fi
        set_version "$2"
        update_version_in_files "$2"
        echo -e "${GREEN}✓ Version updated successfully!${NC}"
        ;;
    *)
        echo "CalVer Version Management"
        echo ""
        echo "Usage:"
        echo "  $0                Show current version"
        echo "  $0 bump          Bump micro version (or reset if new month)"
        echo "  $0 set <version> Set specific version (YYYY.MM.MICRO)"
        echo ""
        echo "Current version: $(get_current_version)"
        ;;
esac
