#!/bin/bash
# Script to update the management console with a new version
# Usage: update_management_console.sh /path/to/new/files /path/to/destination

set -e

# Check for the first argument (source directory)
if [ ! -z "$1" ]; then
    SOURCE_DIR="$1"
else
    echo "Error: Source directory not specified"
    echo "Usage: $0 /path/to/new/files /path/to/destination"
    exit 1
fi

# Get the destination directory
if [ ! -z "$2" ]; then
    DEST_DIR="$2"
else
    echo "Error: Destination directory not specified"
    echo "Usage: $0 /path/to/new/files /path/to/destination"
    exit 1
fi

# Check if the source directory exists
if [ ! -d "$SOURCE_DIR" ]; then
    echo "Error: Source directory '$SOURCE_DIR' does not exist"
    exit 1
fi
BACKUP_DIR=$(mktemp -d)
TIMESTAMP=$(date +%Y%m%d%H%M%S)
# Try to use a persistent location for logs, but fall back to /tmp if we don't have permissions
# First try to use a directory where we're likely to have write permissions
LOG_DIR="${DEST_DIR}/logs"
# Create the logs directory if it doesn't exist - this should work since it's in our app directory
mkdir -p "$LOG_DIR" 2>/dev/null

# If that fails or isn't writable, try to use a system log directory
if [ ! -w "$LOG_DIR" ]; then
    LOG_DIR="/var/log"
    # We might not have permission to create directories in /var/log
    if [ ! -w "$LOG_DIR" ]; then
        LOG_DIR="/tmp"
    fi
fi

LOG_FILE="${LOG_DIR}/update_management_console_${TIMESTAMP}.log"
# Verify we can write to the log file
touch "$LOG_FILE" 2>/dev/null || LOG_FILE="/tmp/update_management_console_${TIMESTAMP}.log"

echo "Starting update process at $(date)" | tee -a "$LOG_FILE"
echo "Source directory: $SOURCE_DIR" | tee -a "$LOG_FILE"
echo "Destination directory: $DEST_DIR" | tee -a "$LOG_FILE"
echo "Creating backup at: $BACKUP_DIR" | tee -a "$LOG_FILE"
echo "Logging output to: $LOG_FILE" | tee -a "$LOG_FILE"

# Function to clean up on error
function cleanup_on_error() {
    echo "Error occurred: $1" | tee -a "$LOG_FILE"
    echo "Restoring backup from $BACKUP_DIR" | tee -a "$LOG_FILE"
    
    # Restore files
    echo "Restoring files..." | tee -a "$LOG_FILE"
    if ! sudo rsync -a "$BACKUP_DIR/" "$DEST_DIR/"; then
        echo "Fatal: Failed to restore backup!" | tee -a "$LOG_FILE"
        exit 2
    fi
    
    # Start service
    echo "Restarting managementconsole service..." | tee -a "$LOG_FILE"
    if ! sudo systemctl start managementconsole; then
        echo "Fatal: Failed to restart service after restore!" | tee -a "$LOG_FILE"
        exit 3
    fi
    
    echo "Backup restored and service restarted" | tee -a "$LOG_FILE"
    rm -rf "$BACKUP_DIR"
    exit 1
}

# Step 1: Create backup of current files
echo "Creating backup of current files..." | tee -a "$LOG_FILE"
if ! sudo rsync -a "$DEST_DIR/" "$BACKUP_DIR/"; then
    echo "Failed to create backup" | tee -a "$LOG_FILE"
    exit 1
fi

# Step 2: Stop the service
echo "Stopping managementconsole service..." | tee -a "$LOG_FILE"
if ! sudo systemctl stop managementconsole; then
    echo "Warning: Failed to stop service, continuing anyway" | tee -a "$LOG_FILE"
fi

# Step 3: Copy new files
echo "Copying new files to $DEST_DIR..." | tee -a "$LOG_FILE"
if ! sudo rsync -a "$SOURCE_DIR/" "$DEST_DIR/"; then
    cleanup_on_error "Failed to copy new files"
fi

# Step 4: Start the service
echo "Starting managementconsole service..." | tee -a "$LOG_FILE"
if ! sudo systemctl start managementconsole; then
    cleanup_on_error "Failed to start service with new files"
fi

# Verify service is running
echo "Verifying service is running..." | tee -a "$LOG_FILE"
sleep 3
if ! systemctl is-active --quiet managementconsole; then
    cleanup_on_error "Service failed to start properly"
fi

echo "Update completed successfully at $(date)" | tee -a "$LOG_FILE"

echo "Removing update directory" | tee -a "$LOG_FILE"
rm -rf "$SOURCE_DIR"

# echo "Removing backup directory" | tee -a "$LOG_FILE"
# rm -rf "$BACKUP_DIR"

exit 0