# Docker Setup and Cleanup Guide for Shuffle

## 🐛 Issue Fixed

The original error was:
```
ERROR: for opensearch  Cannot create container for service opensearch: failed to mount local volume: mount /tmp/shuffle-database:/var/lib/docker/volumes/shuffle_shuffle-database/_data, flags: 0x1000: no such file or directory
```

This occurred because the required directories `/tmp/shuffle-database`, `/tmp/shuffle-apps`, and `/tmp/shuffle-files` didn't exist on the host system before Docker tried to mount them.

## ✅ Solution

1. **Created required directories**: The setup script now creates all necessary directories with proper permissions
2. **Fixed Docker socket permissions**: Ensured non-root access to Docker daemon
3. **Added validation**: Scripts validate Docker and Docker Compose installation
4. **Created cleanup functionality**: Complete cleanup script for fresh installations

## 🚀 Quick Start

### Option 1: Automated Setup (Recommended)
```bash
# Run the setup script
./setup-docker.sh

# Start services
docker compose up -d

# Check status
docker compose ps
```

### Option 2: Manual Setup
```bash
# Create required directories
mkdir -p /tmp/shuffle-database /tmp/shuffle-apps /tmp/shuffle-files
chmod 755 /tmp/shuffle-database /tmp/shuffle-apps /tmp/shuffle-files

# Fix Docker permissions (if needed)
sudo chmod 666 /var/run/docker.sock

# Start services
docker compose up -d
```

## 🧹 Complete Cleanup

To completely remove all Shuffle containers, volumes, networks, and data:

```bash
./cleanup-docker.sh
```

This script will:
- Stop all Docker Compose services
- Remove Shuffle containers
- Remove Docker networks
- Remove Docker volumes
- Clean data directories
- Optionally remove Docker images
- Recreate clean directories

## 📋 Available Scripts

| Script | Purpose | Usage |
|--------|---------|-------|
| `setup-docker.sh` | Initial setup and validation | `./setup-docker.sh` |
| `cleanup-docker.sh` | Complete cleanup | `./cleanup-docker.sh` |
| `test_system.sh` | System readiness test | `./test_system.sh` |

## 🌐 Access Points

After successful deployment:

- **Frontend**: http://localhost:3001
- **Backend API**: http://localhost:5001  
- **OpenSearch**: http://localhost:9200

## 🔧 Configuration

Environment variables are configured in `.env` file:

- `DB_LOCATION`: Database storage location (default: `/tmp/shuffle-database`)
- `SHUFFLE_APP_HOTLOAD_LOCATION`: Apps directory (default: `/tmp/shuffle-apps`)
- `SHUFFLE_FILE_LOCATION`: Files directory (default: `/tmp/shuffle-files`)
- `FRONTEND_PORT`: Frontend port (default: `3001`)
- `BACKEND_PORT`: Backend port (default: `5001`)

## 🐳 Docker Commands Reference

```bash
# Start services
docker compose up -d

# Stop services
docker compose down

# View logs
docker compose logs -f

# Check status
docker compose ps

# Restart specific service
docker compose restart <service_name>

# Pull latest images
docker compose pull
```

## 🔍 Troubleshooting

### Common Issues

1. **Permission denied accessing Docker socket**
   ```bash
   sudo chmod 666 /var/run/docker.sock
   ```

2. **Directory doesn't exist errors**
   ```bash
   mkdir -p /tmp/shuffle-database /tmp/shuffle-apps /tmp/shuffle-files
   ```

3. **Port conflicts**
   - Modify ports in `.env` file
   - Check for conflicting services: `netstat -tulpn | grep :3001`

4. **Out of disk space**
   - Clean unused Docker resources: `docker system prune -f`
   - Remove old images: `docker image prune -f`

### Logs and Debugging

```bash
# View all logs
docker compose logs

# View specific service logs
docker compose logs opensearch
docker compose logs backend

# Follow logs in real-time
docker compose logs -f --tail=100
```

## 🔄 Update Process

To update to latest versions:

```bash
# Stop current services
docker compose down

# Pull latest images
docker compose pull

# Start with new images
docker compose up -d
```

## 🛡️ Security Notes

- Change default passwords in `.env` file
- Use strong passwords for `SHUFFLE_OPENSEARCH_PASSWORD`
- Update `SHUFFLE_ENCRYPTION_MODIFIER` to a unique value
- Consider using environment-specific configurations for production

## 📊 Monitoring

Monitor resource usage:

```bash
# Container stats
docker stats

# Disk usage
docker system df

# Volume usage
docker volume ls
```