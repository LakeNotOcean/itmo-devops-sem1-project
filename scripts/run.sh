#!/bin/bash

set -e

SCRIPT_DIR="$(dirname "$0")"
STATE_FILE="$SCRIPT_DIR/configs/.scripts_varaibles"

echo "Starting full deployment process"
echo "======================================"

# Load configuration
source "$SCRIPT_DIR/configs/.env"

# Check dependencies
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
"$SCRIPT_DIR/create-infrastructure.sh"
if [ $? -ne 0 ]; then
    echo "Error: Stage 1 failed"
    exit 1
fi

# Load state
if [ -f "$STATE_FILE" ]; then
    source "$STATE_FILE"
    echo "State loaded: VM_IP=$VM_IP"
else
    echo "Error: State file not found"
    exit 1
fi

# Stage 2: Install dependencies
echo -e "\nStage 2: Installing dependencies"
"$SCRIPT_DIR/install-dependencies.sh" "$VM_IP"
if [ $? -ne 0 ]; then
    echo "Error: Stage 2 failed"
    exit 1
fi

# Stage 3: Deploy application
echo -e "\nStage 3: Deploying application"
"$SCRIPT_DIR/deploy-application.sh" "$VM_IP"
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
echo "  SSH Access:   ssh ubuntu@$VM_IP"