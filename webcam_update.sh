#!/bin/bash

# Clean exit on shutdown signals
trap "echo 'Shutdown signal received! Exiting cleanly.'; exit" SIGTERM SIGINT

# --- Configuration ---
REPO_DIR="/home/teo/bruggi.it"
WEBCAM_DIR="static/webcam"
TIMESTAMP=$(date +"%Y-%m-%d_%H-%M-%S")
RETENTION_COUNT=20
MAX_RETRIES=3
RETRY_DELAY=5

# --- 1. Wait for Internet Time Sync ---
echo "Waiting for internet time sync..."
for i in {1..30}; do
    if timedatectl status | grep -q "System clock synchronized: yes"; then
        echo "Time synchronized: $(date)"
        break
    fi
    echo "Waiting for NTP... ($i/30)"
    sleep 2
done

# --- 2. Navigate and Check Health ---
cd "$REPO_DIR" || { echo "Directory $REPO_DIR not found!"; exit 1; }

# Ensure Sparse Checkout is active for the target dir
git sparse-checkout set "$WEBCAM_DIR"

run_with_retry() {
    local n=1
    until [ $n -gt $MAX_RETRIES ]; do
        "$@" && return 0 || {
            echo "Command failed. Attempt $n/$MAX_RETRIES..."
            ((n++))
            sleep $RETRY_DELAY
        }
    done
    return 1
}

# --- 3. Pull Only the Webcam Dir ---
# This prevents push rejection if the remote has changed (e.g. from another Pi)
echo "Syncing remote changes (Sparse Pull)..."
run_with_retry git pull origin main --rebase

# --- 4. Capture ---
echo "Capturing image..."
# Using run_with_retry for the camera in case the pipeline is temporarily locked
run_with_retry rpicam-jpeg -o capture.jpg -t 2000 -n --width 1920 --height 1080

# Move capture to the webcam folder
mkdir -p "$WEBCAM_DIR"
mv capture.jpg "$WEBCAM_DIR/$TIMESTAMP.jpg"

# --- 5. Retention Policy ---
echo "Pruning old images (keeping last $RETENTION_COUNT)..."
ls -1 "$WEBCAM_DIR"/[0-9]*.jpg 2>/dev/null | sort -r | tail -n +$((RETENTION_COUNT + 1)) | xargs -I {} rm -f {}

# --- 6. Git Push ---
# We ONLY add the webcam directory
git add "$WEBCAM_DIR"
git commit -m "webcam update: $TIMESTAMP" || echo "Nothing to commit"

echo "Pushing to GitHub..."
run_with_retry git push origin main

# --- 7. Final Cleanup and Power Down ---
sync
echo "Workflow complete. Shutting down."
sudo shutdown -h now
