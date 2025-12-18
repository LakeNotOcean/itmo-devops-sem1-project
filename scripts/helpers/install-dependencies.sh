#!/bin/bash

set -e

SCRIPT_DIR="$(dirname "$0")"
source "$SCRIPT_DIR/configs/.env"
source "$SCRIPT_DIR/utils.sh"

echo "Installing dependencies on $VM_IP"

# check SSH connection
check_ssh $SSH_USER $VM_IP

echo "Installing Docker..."
ssh -i "$SSH_KEY" "$REMOTE_USER@$REMOTE_HOST" bash <<'ENDSSH'
    set -e

    # Update system
    echo "Updating packages..."
    sudo apt-get update

    # Add Docker's official GPG key:
    sudo apt update
    sudo apt install ca-certificates curl
    sudo install -m 0755 -d /etc/apt/keyrings
    sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
    sudo chmod a+r /etc/apt/keyrings/docker.asc

    # Add the repository to Apt sources:
    sudo tee /etc/apt/sources.list.d/docker.sources <<EOF
    Types: deb
    URIs: https://download.docker.com/linux/ubuntu
    Suites: $(. /etc/os-release && echo "${UBUNTU_CODENAME:-$VERSION_CODENAME}")
    Components: stable
    Signed-By: /etc/apt/keyrings/docker.asc
    EOF

    sudo apt update
    sudo apt install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

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

echo "Dependencies installed successfully"