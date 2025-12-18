#!/bin/bash

# change or add varaible
# update_env_var "$FILE_NAME" "VAR" "$VAR"
update_env_var() {
    local file_name="$1"
    local key="$2"
    local value="$3"
    
    if grep -q "^$key=" "$file_name"; then
        sed -i "s/^$key=.*/$key=$value/" "$file_name"
    else
        echo "$key=$value" >> "$file_name"
    fi
}

# check ssh
check_ssh() {
    $SSH_USER=$1
    $VM_IP=$2

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
}
