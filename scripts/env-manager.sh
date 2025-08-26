#!/bin/bash

# Environment Variable Manager
# This script helps manage environment variables across different environments

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
ENV_FILE="$PROJECT_ROOT/.env"
ENV_TEMPLATE="$PROJECT_ROOT/.env.template"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Function to show usage
show_usage() {
    echo "Environment Variable Manager"
    echo ""
    echo "Usage: $0 <command> [options]"
    echo ""
    echo "Commands:"
    echo "  init                    Initialize .env file from template"
    echo "  validate               Validate current .env file"
    echo "  sync-to-github         Sync environment variables to GitHub secrets/vars"
    echo "  generate-docker-env    Generate Docker environment file"
    echo "  check-consistency      Check consistency across all environments"
    echo "  help                   Show this help message"
    echo ""
}

# Function to initialize .env file from template
init_env() {
    print_info "Initializing .env file from template..."

    if [[ -f "$ENV_FILE" ]]; then
        print_warning ".env file already exists. Do you want to overwrite it? (y/N)"
        read -r response
        if [[ ! "$response" =~ ^[Yy]$ ]]; then
            print_info "Initialization cancelled."
            return 0
        fi
    fi

    if [[ ! -f "$ENV_TEMPLATE" ]]; then
        print_error "Template file $ENV_TEMPLATE not found!"
        return 1
    fi

    cp "$ENV_TEMPLATE" "$ENV_FILE"
    print_success ".env file created from template"
    print_warning "Please edit .env file and fill in the actual values"
}

# Function to validate .env file
validate_env() {
    print_info "Validating .env file..."

    if [[ ! -f "$ENV_FILE" ]]; then
        print_error ".env file not found! Run 'init' command first."
        return 1
    fi

    # Check for required variables
    required_vars=(
        "SLACK_TOKEN"
        "DATABASE_HOST"
        "DATABASE_PORT"
        "DATABASE_USER"
        "DATABASE_PASSWORD"
        "DATABASE_NAME"
        "REDIS_HOST"
        "REDIS_PORT"
        "JWT_SECRET"
    )

    missing_vars=()

    for var in "${required_vars[@]}"; do
        if ! grep -q "^${var}=" "$ENV_FILE" || grep -q "^${var}=$" "$ENV_FILE" || grep -q "^${var}=.*-here$" "$ENV_FILE"; then
            missing_vars+=("$var")
        fi
    done

    if [[ ${#missing_vars[@]} -gt 0 ]]; then
        print_error "Missing or incomplete environment variables:"
        for var in "${missing_vars[@]}"; do
            echo "  - $var"
        done
        return 1
    fi

    print_success "All required environment variables are set"
}

# Function to sync environment variables to GitHub
sync_to_github() {
    print_info "Syncing environment variables to GitHub..."

    if [[ ! -f "$ENV_FILE" ]]; then
        print_error ".env file not found!"
        return 1
    fi

    print_warning "This will update GitHub repository variables."
    print_warning "Make sure you have GitHub CLI installed and authenticated."
    print_warning "Continue? (y/N)"
    read -r response
    if [[ ! "$response" =~ ^[Yy]$ ]]; then
        print_info "Sync cancelled."
        return 0
    fi

    # Check if gh CLI is available
    if ! command -v gh &> /dev/null; then
        print_error "GitHub CLI (gh) is not installed. Please install it first."
        return 1
    fi

    # Variables to sync as repository variables (non-sensitive)
    repo_vars=(
        "DATABASE_HOST"
        "DATABASE_PORT"
        "DATABASE_USER"
        "DATABASE_NAME"
        "REDIS_HOST"
        "REDIS_PORT"
        "APP_PORT"
        "APP_ENV"
    )

    # Variables to sync as secrets (sensitive)
    secrets=(
        "SLACK_TOKEN"
        "DATABASE_PASSWORD"
        "REDIS_PASSWORD"
        "JWT_SECRET"
    )

    # Sync repository variables
    for var in "${repo_vars[@]}"; do
        value=$(grep "^${var}=" "$ENV_FILE" | cut -d'=' -f2- | sed 's/^"//' | sed 's/"$//')
        if [[ -n "$value" ]]; then
            print_info "Setting repository variable: $var"
            gh variable set "$var" --body "$value" || print_warning "Failed to set $var"
        fi
    done

    # Sync secrets
    for var in "${secrets[@]}"; do
        value=$(grep "^${var}=" "$ENV_FILE" | cut -d'=' -f2- | sed 's/^"//' | sed 's/"$//')
        if [[ -n "$value" ]]; then
            print_info "Setting secret: $var"
            echo "$value" | gh secret set "$var" || print_warning "Failed to set $var"
        fi
    done

    print_success "Environment variables synced to GitHub"
}

# Function to generate Docker environment file
generate_docker_env() {
    print_info "Generating Docker environment file..."

    if [[ ! -f "$ENV_FILE" ]]; then
        print_error ".env file not found!"
        return 1
    fi

    # Docker environment file is the same as .env for this project
    print_success "Docker will use the existing .env file"
}

# Function to check consistency across environments
check_consistency() {
    print_info "Checking consistency across environments..."

    # Check if .env file exists
    if [[ ! -f "$ENV_FILE" ]]; then
        print_error ".env file not found!"
        return 1
    fi

    # Check GitHub variables (if gh CLI is available)
    if command -v gh &> /dev/null; then
        print_info "Checking GitHub repository variables..."

        repo_vars=(
            "DATABASE_HOST"
            "DATABASE_PORT"
            "DATABASE_USER"
            "DATABASE_NAME"
            "REDIS_HOST"
            "REDIS_PORT"
        )

        inconsistent_vars=()

        for var in "${repo_vars[@]}"; do
            local_value=$(grep "^${var}=" "$ENV_FILE" | cut -d'=' -f2- | sed 's/^"//' | sed 's/"$//')
            github_value=$(gh variable get "$var" 2>/dev/null || echo "")

            if [[ "$local_value" != "$github_value" ]]; then
                inconsistent_vars+=("$var")
                print_warning "$var: local='$local_value', github='$github_value'"
            fi
        done

        if [[ ${#inconsistent_vars[@]} -gt 0 ]]; then
            print_error "Inconsistent variables found between local and GitHub"
            return 1
        else
            print_success "Local and GitHub variables are consistent"
        fi
    else
        print_warning "GitHub CLI not available, skipping GitHub consistency check"
    fi

    print_success "Consistency check completed"
}

# Main script logic
case "${1:-}" in
    "init")
        init_env
        ;;
    "validate")
        validate_env
        ;;
    "sync-to-github")
        sync_to_github
        ;;
    "generate-docker-env")
        generate_docker_env
        ;;
    "check-consistency")
        check_consistency
        ;;
    "help"|"--help"|"-h")
        show_usage
        ;;
    "")
        print_error "No command specified"
        show_usage
        exit 1
        ;;
    *)
        print_error "Unknown command: $1"
        show_usage
        exit 1
        ;;
esac
