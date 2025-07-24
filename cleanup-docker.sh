#!/bin/bash

echo "🧹 Cleaning up Docker Compose environment for Shuffle..."

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

# Step 1: Stop all Docker Compose services
print_status "Stopping Docker Compose services..."
if docker compose down; then
    print_success "Docker Compose services stopped"
else
    print_warning "Failed to stop some services (might be already stopped)"
fi

# Step 2: Remove all containers (including stopped ones)
print_status "Removing all Docker containers..."
CONTAINERS=$(docker ps -aq --filter "name=shuffle")
if [ -n "$CONTAINERS" ]; then
    docker rm -f $CONTAINERS
    print_success "Removed Shuffle containers"
else
    print_warning "No Shuffle containers found"
fi

# Step 3: Remove Docker networks
print_status "Removing Docker networks..."
NETWORKS=$(docker network ls --filter "name=shuffle" -q)
if [ -n "$NETWORKS" ]; then
    docker network rm $NETWORKS 2>/dev/null || true
    print_success "Removed Shuffle networks"
else
    print_warning "No Shuffle networks found"
fi

# Step 4: Remove Docker volumes
print_status "Removing Docker volumes..."
VOLUMES=$(docker volume ls --filter "name=shuffle" -q)
if [ -n "$VOLUMES" ]; then
    docker volume rm $VOLUMES 2>/dev/null || true
    print_success "Removed Shuffle volumes"
else
    print_warning "No Shuffle volumes found"
fi

# Step 5: Clean up data directories
print_status "Cleaning up data directories..."
sudo rm -rf /tmp/shuffle-database/* 2>/dev/null || true
sudo rm -rf /tmp/shuffle-apps/* 2>/dev/null || true
sudo rm -rf /tmp/shuffle-files/* 2>/dev/null || true
print_success "Cleaned up data directories"

# Step 6: Recreate clean directories
print_status "Recreating clean directories..."
mkdir -p /tmp/shuffle-database /tmp/shuffle-apps /tmp/shuffle-files
chmod 755 /tmp/shuffle-database /tmp/shuffle-apps /tmp/shuffle-files
print_success "Recreated clean directories"

# Step 7: Optional - Remove unused Docker images
read -p "Do you want to remove unused Docker images? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    print_status "Removing unused Docker images..."
    docker image prune -f
    print_success "Removed unused Docker images"
fi

# Step 8: Optional - Remove all Shuffle-related images
read -p "Do you want to remove all Shuffle Docker images? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    print_status "Removing Shuffle Docker images..."
    docker images | grep -E "(shuffle|opensearch)" | awk '{print $3}' | xargs docker rmi -f 2>/dev/null || true
    print_success "Removed Shuffle Docker images"
fi

# Step 9: Show final status
print_status "Final cleanup verification..."
echo "Docker containers:"
docker ps -a --filter "name=shuffle"
echo ""
echo "Docker volumes:"
docker volume ls --filter "name=shuffle"
echo ""
echo "Docker networks:"
docker network ls --filter "name=shuffle"
echo ""

print_success "Cleanup completed! You can now run 'docker compose up -d' for a fresh installation."
print_status "Tip: Run './test_system.sh' before starting Docker Compose to verify system readiness."