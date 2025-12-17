#!/bin/bash

# utils.sh - common functions

# Output functions
log_info() {
    echo "[INFO] $1"
}

log_success() {
    echo "[SUCCESS] $1"
}

log_warning() {
    echo "[WARNING] $1"
}

log_error() {
    echo "[ERROR] $1"
}

# Command check
check_command() {
    if ! command -v "$1" &> /dev/null; then
        log_error "Command '$1' not found"
        return 1
    fi
    return 0
}

# SSH access check
check_ssh() {
    local ip="$1"
    local user="${2:-ubuntu}"
    local timeout="${3:-5}"
    
    ssh -o StrictHostKeyChecking=no \
        -o ConnectTimeout="$timeout" \
        -o BatchMode=yes \
        "$user@$ip" "echo OK" &>/dev/null
}

# Wait for SSH
wait_for_ssh() {
    local ip="$1"
    local max_attempts="${2:-20}"
    local wait_time="${3:-10}"
    
    for i in $(seq 1 "$max_attempts"); do
        if check_ssh "$ip"; then
            log_success "SSH is available"
            return 0
        fi
        log_info "Attempt $i/$max_attempts..."
        sleep "$wait_time"
    done
    
    log_error "Failed to connect via SSH"
    return 1
}

# Load state
load_state() {
    local state_file="${1:-.deploy_state}"
    if [ -f "$state_file" ]; then
        source "$state_file"
        return 0
    fi
    return 1
}

# Save state
save_state() {
    local state_file="${1:-.deploy_state}"
    local key="$2"
    local value="$3"
    
    if [ ! -f "$state_file" ]; then
        touch "$state_file"
    fi
    
    # Remove old entry if exists
    grep -v "^$key=" "$state_file" > "${state_file}.tmp" 2>/dev/null || true
    mv "${state_file}.tmp" "$state_file"
    
    # Add new entry
    echo "$key=\"$value\"" >> "$state_file"
}