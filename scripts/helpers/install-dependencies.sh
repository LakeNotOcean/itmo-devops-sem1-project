#!/bin/bash

set -e

if [ $# -lt 1 ]; then
    echo "Error: Please specify server IP address"
    echo "Usage: $0 <vm_ip>"
    exit 1
fi

VM_IP="$1"
echo "Installing dependencies on $VM_IP"

# Check SSH connection
echo "Checking SSH connection..."
if ! ssh -o StrictHostKeyChecking=no -o ConnectTimeout=5 ubuntu@"$VM_IP" "echo SSH OK" &>/dev/null; then
    echo "Error: Failed to connect via SSH"
    exit 1
fi

# Install Docker
echo "Installing Docker..."
ssh -o StrictHostKeyChecking=no ubuntu@"$VM_IP" bash <<'ENDSSH'
set -e

# Update system
echo "Updating packages..."
sudo apt-get update

# Install dependencies
echo "Installing Docker dependencies..."
sudo apt-get install -y \
    ca-certificates \
    curl \
    gnupg \
    lsb-release

# Add Docker GPG key
echo "Adding Docker repository..."
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg

# Add Docker repository
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# Install Docker
echo "Installing Docker CE..."
sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin

# Configure Docker
echo "Configuring Docker..."
sudo usermod -aG docker ubuntu
sudo systemctl enable docker
sudo systemctl start docker

# Verify installation
echo "Verifying Docker installation..."
docker --version
docker compose version
ENDSSH

if [ $? -eq 0 ]; then
    echo "Docker installed successfully"
else
    echo "Error installing Docker"
    exit 1
fi

# Additional packages (optional)
echo "Installing additional packages..."
ssh -o StrictHostKeyChecking=no ubuntu@"$VM_IP" bash <<'ENDSSH'
set -e

# Useful utilities
sudo apt-get install -y \
    git \
    htop \
    net-tools \
    tree \
    ncdu \
    zip \
    unzip

echo "Additional packages installed"
ENDSSH

echo "Dependencies installed successfully"