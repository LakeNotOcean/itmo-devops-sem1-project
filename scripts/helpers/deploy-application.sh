#!/bin/bash

set -e

SCRIPT_DIR="$(dirname "$0")"
source "$SCRIPT_DIR/../../configs/.env"
source "$SCRIPT_DIR/utils.sh"
SSH_COMMON_OPTS="-o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o ConnectTimeout=30"

echo "Deploying application to $VM_IP"

# check SSH connection
set +e
wait_ssh $SSH_USER $VM_IP
set -e

echo "Creating archive..."
tar czf deploy.tar.gz \
    cmd \
    configs \
    internal \
    scripts \
    go.mod \
    go.sum \
    docker-compose.yml \
    Dockerfile

echo "Copying archive to server..."
scp $SSH_COMMON_OPTS \
    -i "$SSH_KEY_PRIVATE" \
    deploy.tar.gz "$SSH_USER@$VM_IP:deploy.tar.gz"

echo "Extracting and running on server..."
ssh $SSH_COMMON_OPTS \
    -t -i "$SSH_KEY_PRIVATE" \
    "$SSH_USER@$VM_IP" bash <<ENDSSH
    set -e
    cd ~
    REMOTE_DIR="$REMOTE_DIR"
    
    echo 'Stopping containers...'
    docker compose down || true
    
    echo 'Cleaning or creating directory...'
    if [ -d "\$REMOTE_DIR" ]; then
        rm -rf "\$REMOTE_DIR"
    fi
    mkdir -p "\$REMOTE_DIR"
    
    echo 'Extracting archive...'
    tar xzf deploy.tar.gz -C "\$REMOTE_DIR"
    rm deploy.tar.gz

    echo 'Entering application directory...'
    cd "\$REMOTE_DIR"
    
    echo 'Loading environment variables...'
    set -a
    source ./configs/.env
    set +a
    
    echo 'Running build script...'
    chmod +x ./scripts/prepare.sh
    ./scripts/prepare.sh
    
    echo 'Starting containers...'
    sudo docker compose --env-file ./configs/.env up -d --build
    
    echo 'Done! Containers are running.'
ENDSSH

rm deploy.tar.gz

echo "Application deployment completed successfully"