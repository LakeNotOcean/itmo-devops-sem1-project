#!/bin/bash

set -e

SCRIPT_DIR="$(dirname "$0")"
source "$SCRIPT_DIR/configs/.env"
source "$SCRIPT_DIR/utils.sh"

echo "Creating infrastructure..."

# ssh key check
if [ ! -f "$SSH_KEY_PUB" ]; then
    echo "Error: SSH key not found"
    exit 1
fi

echo "Checking for existing VM..."
EXISTING_VM=$(yc compute instance list --folder-id="$YC_FOLDER_ID" --format=json | jq -r --arg vm_name "$VM_NAME" '.[] | select(.name==$vm_name) | .id')

# we won't create an VM if it already exists   
if [ "$EXISTING_VM" ]; then
    echo "VN exists already!"
    exit 0
fi

echo "Creating network..."
NETWORK_ID=$(yc vpc network list --folder-id="$YC_FOLDER_ID" --format=json | jq -r --arg network_name "$NETWORK_NAME" '.[] | select(.name==$network_name) | .id')
if [ -z "$NETWORK_ID" ]; then
    NETWORK_ID=$(yc vpc network create \
        --name="$NETWORK_NAME" \
        --folder-id="$YC_FOLDER_ID" \
        --format=json | jq -r '.id')
    echo "Network created: $NETWORK_ID"
else
    echo "Using existing network: $NETWORK_ID"
fi

echo "Creating subnet..."
SUBNET_ID=$(yc vpc subnet list --folder-id="$YC_FOLDER_ID" --format=json | jq -r --arg subnet_name "$SUBNET_NAME" '.[] | select(.name==$subnet_name) | .id')
if [ -z "$SUBNET_ID" ]; then
    SUBNET_ID=$(yc vpc subnet create \
        --name="$SUBNET_NAME" \
        --folder-id="$YC_FOLDER_ID" \
        --network-id="$NETWORK_ID" \
        --zone="$YC_ZONE" \
        --range="$SUBNET_CIDR" \
        --format=json | jq -r '.id')
    echo "Subnet created: $SUBNET_ID"
else
    echo "Using existing subnet: $SUBNET_ID"
fi

echo "Preparing configuration..."
SSH_KEY_CONTENT=$(cat "$SSH_KEY_PUB")
CLOUD_CONFIG=$(cat <<EOF
#cloud-config
users:
  - name: ubuntu
    groups: sudo
    shell: /bin/bash
    sudo: ['ALL=(ALL) NOPASSWD:ALL']
    ssh-authorized-keys:
      - $SSH_KEY_CONTENT
EOF
)

CLOUD_CONFIG_FILE=$(mktemp)
echo "$CLOUD_CONFIG" > "$CLOUD_CONFIG_FILE"


echo "Creating virtual machine..."
PREEMPTIBLE_ARG=""
[ "$USE_PREEMPTIBLE" = "true" ] && PREEMPTIBLE_ARG="--preemptible"

VM_ID=$(yc compute instance create \
    --name="$VM_NAME" \
    --folder-id="$YC_FOLDER_ID" \
    --zone="$YC_ZONE" \
    --platform="$PLATFORM_ID" \
    --cores="$VM_CORES" \
    --memory="${VM_MEMORY}GB" \
    --create-boot-disk size="${VM_DISK_SIZE}GB",image-id="$VM_IMAGE",type="network-hdd" \
    --network-interface subnet-id="$SUBNET_ID",nat-ip-version=ipv4 \
    --metadata-from-file user-data="$CLOUD_CONFIG_FILE" \
    --metadata "enable-oslogin=true" \
    $PREEMPTIBLE_ARG \
    --format=json | jq -r '.id')

rm -f "$CLOUD_CONFIG_FILE"
echo "VM created: $VM_ID"

echo "Waiting for IP address..."
sleep 20

VM_IP=$(yc compute instance get "$VM_ID" --folder-id="$YC_FOLDER_ID" --format=json | jq -r '.network_interfaces[0].primary_v4_address.one_to_one_nat.address')
if [ -z "$VM_IP" ] || [ "$VM_IP" = "null" ]; then
    echo "Error: Failed to get IP address"
    exit 1
fi

echo "VM IP: $VM_IP"

echo "Waiting for SSH availability..."
DELAY=20
ATTEMPTS=50
for ((i=1; i<=$ATTEMPTS; i++)); do
    echo "  Attempt $i/100..."
    
    if ssh -o StrictHostKeyChecking=no \
           -o ConnectTimeout=5 \
           -o BatchMode=yes \
           -o PasswordAuthentication=no \
           -o UserKnownHostsFile=/dev/null \
           "$SSH_USER@$VM_IP" "echo OK" &>/dev/null; then
        echo "SSH is available!"
        break
    fi
    
    if [[ $i -eq $MAX_ATTEMPTS ]]; then
        echo "SSH is not available"
        exit 1
    fi
    
    echo "Waiting $DELAY seconds..."
    sleep $DELAY
done

# update env
declare -A config=(
    ["VM_ID"]="$VM_ID"
    ["VM_IP"]="$VM_IP"
    ["NETWORK_ID"]="$NETWORK_ID"
    ["SUBNET_ID"]="$SUBNET_ID"
)

for key in "${!config[@]}"; do
    value="${config[$key]}"
    update_env_var "$FILE_NAME" "$key" "$value"
done


echo "Infrastructure created successfully!" 