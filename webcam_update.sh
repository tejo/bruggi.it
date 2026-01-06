#!/bin/bash

# Clean exit on shutdown signals from Witty Pi
trap "echo 'Shutdown signal received! Exiting cleanly.'; exit" SIGTERM SIGINT

# Configuration
REPO_DIR="/home/teo/bruggi.it"
WEBCAM_DIR="static/webcam"
TIMESTAMP=$(date +"%Y-%m-%d_%H-%M-%S")
RETENTION_COUNT=20
MAX_RETRIES=3
RETRY_DELAY=5

run_with_retry() {
    local n=1
    until [ $n -gt $MAX_RETRIES ]
    do
        "$@" && return 0 || {
            echo "Command failed. Attempt $n/$MAX_RETRIES..."
            ((n++))
            sleep $RETRY_DELAY
        }
    done
    return 1
}

# 1. Navigate and Check Health
cd "$REPO_DIR" || exit 1

if ! git rev-parse --is-inside-work-tree > /dev/null 2>&1; then
    echo "Git corruption detected. Self-healing..."
    cd ..
    rm -rf bruggi.it
    git clone --filter=blob:none --sparse https://github.com/YOUR_USER/bruggi.it.git
    cd bruggi.it
    git sparse-checkout set static
fi

# 2. Sync
run_with_retry git pull --rebase origin main

# 3. Capture
run_with_retry rpicam-jpeg -o capture.jpg -t 1500 -n

# 4. Copy to Static Folder
mkdir -p "$WEBCAM_DIR"
cp capture.jpg "$WEBCAM_DIR/$TIMESTAMP.jpg"

# 5. RETENTION POLICY: Keep only last 20 images
# We look for files matching the timestamp pattern, sort them,
# and delete those older than the top 20.
echo "Cleaning up old images..."
ls -1 "$WEBCAM_DIR"/[0-9]*.jpg 2>/dev/null | sort -r | tail -n +$((RETENTION_COUNT + 1)) | xargs -I {} rm -f {}

# 6. Git Workflow
# We add everything (including the deletions of old files)
git add -A "$WEBCAM_DIR/"
git commit -m "webcam update: $TIMESTAMP (Retention: $RETENTION_COUNT)" || echo "Nothing to commit"

# 7. Push
if ! run_with_retry git push origin main; then
    git pull --rebase origin main
    git push origin main
fi

# 8. Cleanup Temporary File & Flush to Disk
rm -f capture.jpg
sync
echo "Successfully pushed $TIMESTAMP.jpg and pruned old files."
echo "Workflow complete. Shutting down to save power."
# This tells the Witty Pi to prepare for power-off
sudo shutdown -h now
