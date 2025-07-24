#!/bin/bash

echo "🚀 Setting up Shuffle Docker Environment..."

# Function to print colored output
print_status() {
    echo -e "\033[1;34m$1\033[0m"
}

print_success() {
    echo -e "\033[1;32m✅ $1\033[0m"
}

print_warning() {
    echo -e "\033[1;33m⚠️  $1\033[0m"
}

print_error() {
    echo -e "\033[1;31m❌ $1\033[0m"
}

# Step 1: Check if Docker is installed
print_status "Checking Docker installation..."
if command -v docker &> /dev/null; then
    print_success "Docker is installed"
    docker --version
else
    print_error "Docker is not installed. Please install Docker first."
    exit 1
fi

# Step 2: Check if Docker Compose is available
print_status "Checking Docker Compose..."
if docker compose version &> /dev/null; then
    print_success "Docker Compose is available"
    docker compose version
else
    print_error "Docker Compose is not available. Please install Docker Compose plugin."
    exit 1
fi

# Step 3: Check if .env file exists
print_status "Checking environment configuration..."
if [ -f ".env" ]; then
    print_success "Environment file (.env) exists"
else
    print_warning "Environment file (.env) not found"
    echo "Creating default .env file..."
    cat > .env << EOF
# Shuffle Environment Configuration

# Backend Configuration
BACKEND_HOSTNAME=shuffle-backend
BACKEND_PORT=5001

# Frontend Configuration
FRONTEND_PORT=3001
FRONTEND_PORT_HTTPS=3443

# External hostname (change this to your domain or IP)
OUTER_HOSTNAME=localhost

# Database Configuration
DB_LOCATION=/tmp/shuffle-database
SHUFFLE_OPENSEARCH_PASSWORD=StrongShufflePassword321!

# File and App Locations
SHUFFLE_APP_HOTLOAD_LOCATION=/tmp/shuffle-apps
SHUFFLE_FILE_LOCATION=/tmp/shuffle-files

# Proxy Configuration (leave empty if not using proxy)
HTTP_PROXY=
HTTPS_PROXY=
SHUFFLE_PASS_WORKER_PROXY=false
SHUFFLE_PASS_APP_PROXY=false

# Shuffle Configuration
SHUFFLE_DEFAULT_USERNAME=admin
SHUFFLE_DEFAULT_PASSWORD=password
SHUFFLE_DEFAULT_APIKEY=

# Database/Search Configuration
SHUFFLE_ELASTIC=true
SHUFFLE_OPENSEARCH_URL=https://shuffle-opensearch:9200
SHUFFLE_OPENSEARCH_USERNAME=admin
SHUFFLE_OPENSEARCH_CERTIFICATE_FILE=

# Organization Configuration
ORG_ID=shuffle
ENVIRONMENT_NAME=Shuffle

# Worker Configuration
SHUFFLE_SWARM_CONFIG=run
SHUFFLE_WORKER_IMAGE=ghcr.io/shuffle/shuffle-worker:latest

# Logging Configuration
SHUFFLE_LOGS_DISABLED=false
SHUFFLE_STATS_DISABLED=false

# Security Configuration
SHUFFLE_ENCRYPTION_MODIFIER=shuffle_default_encryption_key_change_me

# License Configuration (for our new system)
SHUFFLE_LICENSE_STANDALONE=false

# Debug Configuration
SHUFFLE_DEBUG=false
EOF
    print_success "Created default .env file"
fi

# Step 4: Create required directories
print_status "Creating required directories..."
mkdir -p /tmp/shuffle-database /tmp/shuffle-apps /tmp/shuffle-files
chmod 755 /tmp/shuffle-database /tmp/shuffle-apps /tmp/shuffle-files
print_success "Created and configured directories"

# Step 5: Validate Docker Compose configuration
print_status "Validating Docker Compose configuration..."
if docker compose config > /dev/null 2>&1; then
    print_success "Docker Compose configuration is valid"
else
    print_error "Docker Compose configuration has errors"
    docker compose config
    exit 1
fi

# Step 6: Check Docker daemon
print_status "Checking Docker daemon..."
if docker info > /dev/null 2>&1; then
    print_success "Docker daemon is running"
else
    print_warning "Docker daemon may not be running or accessible"
    print_status "Trying to fix Docker socket permissions..."
    sudo chmod 666 /var/run/docker.sock 2>/dev/null || true
fi

cursor/fix-docker-compose-volume-error-and-add-cleanup-script-e0be
# Step 6.5: Test the system readiness
print_status "Testing system readiness..."
if [ -f "./test_system.sh" ]; then
    if ./test_system.sh > /dev/null 2>&1; then
        print_success "System test passed"
    else
        print_warning "System test failed - some components may need attention"
        print_status "You can run './test_system.sh' to see detailed results"
    fi
else
    print_warning "System test script not found"
fi

=======
 main
# Step 7: Show final instructions
echo ""
print_success "Setup completed successfully!"
echo ""
echo "📋 Next steps:"
echo "1. Review the .env file and modify settings as needed"
echo "2. Start the services: docker compose up -d"
echo "3. Check status: docker compose ps"
echo "4. View logs: docker compose logs -f"
echo ""
echo "🌐 Access URLs:"
echo "   Frontend: http://localhost:3001"
echo "   Backend API: http://localhost:5001"
echo "   OpenSearch: http://localhost:9200"
echo ""
echo "🧹 To clean up everything: ./cleanup-docker.sh"
echo "📊 To test system readiness: ./test_system.sh"