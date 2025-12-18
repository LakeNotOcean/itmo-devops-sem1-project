#!/bin/bash

set -e

if [ $# -lt 1 ]; then
    echo "Error: Please specify server IP address"
    echo "Usage: $0 <vm_ip>"
    exit 1
fi

VM_IP="$1"
SCRIPT_DIR="$(dirname "$0")"
source "$SCRIPT_DIR/configs/.env"

echo "Deploying application to $VM_IP"

# Check local files
echo "Checking local files..."
if [ ! -d "$LOCAL_APP_DIR" ]; then
    echo "Error: Local application directory not found: $LOCAL_APP_DIR"
    exit 1
fi

# Create directory on server
echo "Preparing server..."
ssh -o StrictHostKeyChecking=no ubuntu@"$VM_IP" "mkdir -p $APP_DIR"

# Copy files
echo "Copying files..."
cd "$LOCAL_APP_DIR"

# List of files to copy
FILES_TO_COPY=()
[ -d "./cmd" ] && FILES_TO_COPY+=("./cmd")
[ -d "./internal" ] && FILES_TO_COPY+=("./internal")
[ -d "./migrations" ] && FILES_TO_COPY+=("./migrations")
[ -f "./go.mod" ] && FILES_TO_COPY+=("./go.mod")
[ -f "./go.sum" ] && FILES_TO_COPY+=("./go.sum")
[ -f "./Dockerfile" ] && FILES_TO_COPY+=("./Dockerfile")
[ -f "./docker-compose.yaml" ] && FILES_TO_COPY+=("./docker-compose.yaml")
[ -f "./docker-compose.yml" ] && FILES_TO_COPY+=("./docker-compose.yml")
[ -f "./compose.yaml" ] && FILES_TO_COPY+=("./compose.yaml")
[ -f "./compose.yml" ] && FILES_TO_COPY+=("./compose.yml")
[ -f "./.env" ] && FILES_TO_COPY+=("./.env")
[ -f "./.env.example" ] && FILES_TO_COPY+=("./.env.example")

if [ ${#FILES_TO_COPY[@]} -eq 0 ]; then
    echo "Application files not found, creating test files..."
    # Create test files
    cat > Dockerfile << 'EOF'
FROM nginx:alpine
COPY . /usr/share/nginx/html
EXPOSE 80
EOF
    
    cat > docker-compose.yaml << 'EOF'
version: '3.8'
services:
  web:
    build: .
    ports:
      - "8080:80"
EOF
    
    echo "<h1>Hello from Yandex Cloud!</h1>" > index.html
    FILES_TO_COPY=("./Dockerfile" "./docker-compose.yaml" "./index.html")
fi

# Copy files
scp -o StrictHostKeyChecking=no -r "${FILES_TO_COPY[@]}" ubuntu@"$VM_IP":$APP_DIR/
echo "Files copied"

# Start application
echo "Starting application..."
ssh -o StrictHostKeyChecking=no ubuntu@"$VM_IP" bash <<ENDSSH
set -e
cd $APP_DIR

echo "Current directory: \$(pwd)"
echo "Contents:"
ls -la

# Check Docker Compose
if [ -f "docker-compose.yaml" ] || [ -f "docker-compose.yml" ] || [ -f "compose.yaml" ] || [ -f "compose.yml" ]; then
    echo "Starting Docker Compose..."
    docker compose down 2>/dev/null || true
    docker compose up -d --build
    echo "Application started with Docker Compose"
elif [ -f "Dockerfile" ]; then
    echo "Starting Docker container..."
    docker stop app 2>/dev/null || true
    docker rm app 2>/dev/null || true
    docker build -t myapp .
    docker run -d -p 8080:80 --name app myapp
    echo "Application started in Docker container"
else
    echo "Configuration files not found"
    echo "Creating simple web server..."
    python3 -m http.server 8080 &
    echo "Simple Python HTTP server started on port 8080"
fi

# Check running containers
echo "Running containers:"
docker ps
ENDSSH

# Check application availability
echo "Checking application availability..."
sleep 10

if curl -s --connect-timeout 10 "http://$VM_IP:8080" > /dev/null; then
    echo "Application is available at http://$VM_IP:8080"
else
    echo "Warning: Application may not be available or uses different port"
fi

# Display server information
echo -e "\nServer information:"
ssh -o StrictHostKeyChecking=no ubuntu@"$VM_IP" bash <<'ENDSSH'
echo "1. Disk space:"
df -h /

echo -e "\n2. Memory:"
free -h

echo -e "\n3. Running processes:"
ps aux | grep -E "(docker|nginx|python|go)" | grep -v grep || echo "No relevant processes found"

echo -e "\n4. Network ports:"
sudo netstat -tlnp 2>/dev/null | grep :8080 || echo "Port 8080 is not listening"
ENDSSH

echo "Application deployment completed successfully"