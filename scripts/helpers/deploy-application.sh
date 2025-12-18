#!/bin/bash

set -e

SCRIPT_DIR="$(dirname "$0")"
source "$SCRIPT_DIR/configs/.env"
source "$SCRIPT_DIR/utils.sh"

echo "Deploying application to $VM_IP"

# check SSH connection
check_ssh $SSH_USER $VM_IP

echo "Creating archive..."
tar czf deploy.tar.gz \
    cmd \
    configs \
    internal \
    scripts \
    vendor \
    go.mod \
    go.sum \
    docker-compose.yml \
    Dockerfile

echo "Copying archive to server..."
scp -i "$SSH_KEY" deploy.tar.gz "$REMOTE_USER@$REMOTE_HOST:$REMOTE_DIR.tar.gz"

echo "Extracting and running on server..."
ssh -i "$SSH_KEY" "$REMOTE_USER@$REMOTE_HOST" bash <<'ENDSSH'
    set -e
    
    echo 'Stopping containers...'
    cd $REMOTE_DIR 2>/dev/null && docker-compose down || true
    
    echo 'Creating/cleaning directory...'
    mkdir -p $REMOTE_DIR
    cd $REMOTE_DIR
    find . -mindepth 1 -delete
    
    echo 'Extracting archive...'
    tar xzf ../deploy.tar.gz --strip-components=0
    rm ../deploy.tar.gz
    
    echo 'Loading environment variables...'
    set -a
    source configs/.env
    set +a
    
    echo 'Running build script...'
    chmod +x scripts/prepare.sh
    ./scripts/prepare.sh
    
    echo 'Starting containers...'
    docker-compose up -d --build
    
    echo 'Cleaning up...'
    docker system prune -f
    
    echo 'Done! Containers are running.'
    docker-compose ps
ENDSSH

rm deploy.tar.gz

echo "Application deployment completed successfully"