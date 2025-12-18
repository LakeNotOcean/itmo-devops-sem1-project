#!/bin/bash

set -e

SCRIPT_DIR="$(dirname "$0")"

echo "Starting full deployment process"
echo "======================================"

# Load configuration
source "$SCRIPT_DIR/configs/.env"

# сheck dependencies
if ! command -v yc &> /dev/null; then
    echo "Error: Yandex Cloud CLI not found"
    exit 1
fi

if ! command -v jq &> /dev/null; then
    echo "Error: jq not found"
    exit 1
fi

# Stage 1: Create infrastructure
echo -e "\nStage 1: Creating infrastructure"
"$SCRIPT_DIR/helpers/create-infrastructure.sh"
if [ $? -ne 0 ]; then
    echo "Error: Stage 1 failed"
    exit 1
fi

echo "Stage 2: Installing dependencies"
"$SCRIPT_DIR/helpers/install-dependencies.sh" "$VM_IP"
if [ $? -ne 0 ]; then
    echo "Error: Stage 2 failed"
    exit 1
fi

echo "Stage 3: Deploying application"
"$SCRIPT_DIR/helpers/deploy-application.sh" "$VM_IP"
if [ $? -ne 0 ]; then
    echo "Error: Stage 3 failed"
    exit 1
fi

echo -e "\n======================================"
echo "All stages completed successfully!"
echo "======================================"
echo -e "\nSummary:"
echo "  VM ID:        $VM_ID"
echo "  VM IP:        $VM_IP"
echo "  API Endpoint: http://$VM_IP:8080"
echo "  SSH Access:   ssh $REMOTE_USER@$VM_IP"