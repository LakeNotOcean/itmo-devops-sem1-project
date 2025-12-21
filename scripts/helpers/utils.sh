#!/bin/bash

SSH_COMMON_OPTS="-o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o ConnectTimeout=30"

# change or add varaible
# update_env_var "$FILE_NAME" "VAR" "$VAR"
update_env_var() {
    local file_name="$1"
    local key="$2"
    local value="$3"
    
    if grep -q "^$key=" "$file_name"; then
        sed -i "s/^$key=.*/$key=$value/" "$file_name"
    else
        printf "%s=%s\n" "$key" "$value" >> "$file_name"
    fi
}

# check ssh
check_ssh() {
    local SSH_USER="$1"
    local VM_IP="$2"

    echo "Waiting for SSH availability..."
    DELAY=20
    ATTEMPTS=50
    
    echo "SSH_USER=$SSH_USER"
    echo "VM_IP=$VM_IP"
    echo "SSH_KEY_PRIVATE=$SSH_KEY_PRIVATE"
    
    if [ -z "$SSH_KEY_PRIVATE" ]; then
        echo "Error: SSH_KEY_PRIVATE is not set"
        return 1
    fi
    
    for ((i=1; i<=$ATTEMPTS; i++)); do
        echo "  Attempt $i/$ATTEMPTS..."
        
        SSH_OUTPUT=$(ssh $SSH_COMMON_OPTS \
            -o LogLevel=ERROR \
            -i "$SSH_KEY_PRIVATE" \
            "$SSH_USER@$VM_IP" "echo OK" 2>&1)
        SSH_EXIT_CODE=$?
        
        echo "  SSH exit code: $SSH_EXIT_CODE"
        
        if [ $SSH_EXIT_CODE -eq 0 ]; then
            echo "SSH is available!"
            return 0
        fi

        echo "  SSH error output: $SSH_OUTPUT"
        
        if [ $i -eq $ATTEMPTS ]; then
            echo "Error: SSH is not available after $ATTEMPTS attempts"
            echo "Last error: $SSH_OUTPUT"
            return 1
        fi
        
        echo "Waiting $DELAY seconds..."
        sleep $DELAY
    done
}