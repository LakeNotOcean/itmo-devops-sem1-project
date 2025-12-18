#!/bin/bash

set -e

SCRIPT_DIR="$(dirname "$0")"
source "$SCRIPT_DIR/configs/.env"
STATE_FILE="$SCRIPT_DIR/configs/.scripts_varaibles"

echo "Cleaning up infrastructure"

# Load state
if [ -f "$STATE_FILE" ]; then
    source "$STATE_FILE"
    echo "State loaded"
else
    echo "State file not found"
fi

# Delete VM
if [ -n "$VM_ID" ]; then
    echo "Deleting VM $VM_ID..."
    yc compute instance delete "$VM_ID" --folder-id="$YC_FOLDER_ID" --yes 2>/dev/null || true
    echo "VM deleted"
fi

# Delete subnet
if [ -n "$SUBNET_ID" ]; then
    echo "Deleting subnet $SUBNET_ID..."
    yc vpc subnet delete "$SUBNET_ID" --folder-id="$YC_FOLDER_ID" --yes 2>/dev/null || true
    echo "Subnet deleted"
fi

# Delete network
if [ -n "$NETWORK_ID" ]; then
    echo "Deleting network $NETWORK_ID..."
    yc vpc network delete "$NETWORK_ID" --folder-id="$YC_FOLDER_ID" --yes 2>/dev/null || true
    echo "Network deleted"
fi

# Delete state file
rm -f "$STATE_FILE" 2>/dev/null || true
echo "State file removed"

echo "Cleanup completed successfully"