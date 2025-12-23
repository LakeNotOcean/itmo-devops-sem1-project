#!/bin/bash

set -e

SCRIPT_DIR="$(dirname "$0")"
source "$SCRIPT_DIR/../../configs/.env"
echo $SCRIPT_DIR
source "$SCRIPT_DIR/utils.sh"
SSH_COMMON_OPTS="-o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o ConnectTimeout=30"

echo "Installing dependencies on $VM_IP"

# check SSH connection
set +e
wait_ssh $SSH_USER $VM_IP
set -e

echo "Installing Docker..."
ssh $SSH_COMMON_OPTS \
    -t -i "$SSH_KEY_PRIVATE" \
    "$SSH_USER@$VM_IP" bash <<'ENDSSH'
    set -e

    export DEBIAN_FRONTEND=noninteractive
    
    if command -v docker >/dev/null 2>&1; then
        echo "Docker is already installed!"
        exit 0
    fi

    # Update system
    echo "Updating packages..."
    sudo apt-get update -y

    # Add Docker's official GPG key:
    sudo apt update -y
    sudo apt install -y ca-certificates curl
    sudo install -m 0755 -d /etc/apt/keyrings
    sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
    sudo chmod a+r /etc/apt/keyrings/docker.asc

    # Add the repository to Apt sources:
    echo "Adding Docker repository..."
    UBUNTU_CODENAME=$(lsb_release -cs)
    echo "Ubuntu codename: $UBUNTU_CODENAME"
    echo "Types: deb" | sudo tee /etc/apt/sources.list.d/docker.sources
    echo "URIs: https://download.docker.com/linux/ubuntu" | sudo tee -a /etc/apt/sources.list.d/docker.sources
    echo "Suites: $UBUNTU_CODENAME" | sudo tee -a /etc/apt/sources.list.d/docker.sources
    echo "Components: stable" | sudo tee -a /etc/apt/sources.list.d/docker.sources
    echo "Signed-By: /etc/apt/keyrings/docker.asc" | sudo tee -a /etc/apt/sources.list.d/docker.sources

    sudo apt update -y
    sudo apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

    # Verify installation
    echo "Verifying Docker installation..."
    docker --version
    docker compose version
    sudo systemctl start docker
    sudo systemctl status docker
ENDSSH

if [ $? -eq 0 ]; then
    echo "Docker installed successfully"
else
    echo "Error installing Docker"
    exit 1
fi

echo "Dependencies installed successfully"