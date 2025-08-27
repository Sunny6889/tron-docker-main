#!/bin/bash

TARGET_DIR="/data/FullNode/tron-docker/"
DOCKER_COMPOSE_FILE="docker-compose.fullnode.main.yml"

cd ${TARGET_DIR} || exit
./trond node run-single stop -t full-main -f ${DOCKER_COMPOSE_FILE}

# Wait and check if container is stopped
max_attempts=30  # 2 minutes total (30 * 10 seconds)
attempt=1

while [ $attempt -le $max_attempts ]; do
    if ! docker-compose -f ${DOCKER_COMPOSE_FILE} ps --services --filter "status=running" | grep -q "tron-node-mainnet"; then
        echo "Container successfully stopped"

        # Check if the target file exists
        if [ -f "$DOCKER_COMPOSE_FILE" ]; then
            # Create a backup by renaming the file with a timestamp
            TIMESTAMP=$(date +%Y%m%d%H%M%S)
            BACKUP_FILE="${DOCKER_COMPOSE_FILE}.${TIMESTAMP}.bak"

            # Rename the file
            mv "$DOCKER_COMPOSE_FILE" "$BACKUP_FILE"
            echo "Backed up $DOCKER_COMPOSE_FILE to $BACKUP_FILE"
        else
            echo "No file found at $DOCKER_COMPOSE_FILE, no backup needed"
        fi
        exit 0
    fi
    echo "Attempt $attempt of $max_attempts: Container still running, waiting 10 seconds..."
    sleep 10
    ((attempt++))
done

echo "Failed to stop container after 2 minutes"
exit 1
