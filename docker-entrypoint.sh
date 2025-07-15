#!/bin/bash
set -e

# Function to mount EFS with IP address
mount_efs_with_ip() {
    local efs_ip="$1"
    local efs_id="$2"
    local access_point="$3"
    local mount_path="$4"
    
    echo "Mounting EFS with IP address: $efs_ip"
    
    # Create mount options
    MOUNT_OPTIONS="nfsvers=4.1,rsize=1048576,wsize=1048576,hard,timeo=600,retrans=2,noresvport"
    
    # Add access point if provided
    if [ -n "$access_point" ]; then
        MOUNT_OPTIONS="${MOUNT_OPTIONS},accesspoint=${access_point}"
    fi
    
    # Ensure mount path exists
    mkdir -p "$mount_path"
    
    # Mount using IP address
    mount -t nfs4 -o "$MOUNT_OPTIONS" "${efs_ip}:/" "$mount_path"
    
    echo "EFS mounted successfully at $mount_path"
}

# Check if EFS mounting is needed
if [ -n "$EFS_IP_ADDRESS" ] && [ -n "$EFS_FILE_SYSTEM_ID" ]; then
    echo "EFS IP address provided, attempting to mount..."
    
    # Default mount path
    MOUNT_PATH="${EFS_MOUNT_PATH:-/data}"
    
    # Mount EFS using IP address
    mount_efs_with_ip "$EFS_IP_ADDRESS" "$EFS_FILE_SYSTEM_ID" "$EFS_ACCESS_POINT_ID" "$MOUNT_PATH"
    
    # Verify mount
    if mountpoint -q "$MOUNT_PATH"; then
        echo "EFS mount verified at $MOUNT_PATH"
        
        # Create required directories with proper permissions
        mkdir -p "$MOUNT_PATH/workdir" "$MOUNT_PATH/database" "$MOUNT_PATH/config"
        chown -R unity:unity "$MOUNT_PATH"
    else
        echo "ERROR: EFS mount failed"
        exit 1
    fi
else
    echo "No EFS_IP_ADDRESS provided, skipping EFS mount"
fi

# Drop privileges and execute the main command as unity user
exec gosu unity:unity "$@"