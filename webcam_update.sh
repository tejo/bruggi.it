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
# Critical for valid Git commit timestamps and SSL security
echo "Waiting for internet time sync..."
SYNCED=false
for i in {1..30}; do
    if timedatectl status | grep -q "System clock synchronized: yes"; then
        echo "Time synchronized: $(date)"
        SYNCED=true
        break
    fi
    echo "Waiting for NTP... ($i/30)"
    sleep 2
done

# --- 2. Navigate and Check Health ---
cd "$REPO_DIR" || { echo "Directory $REPO_DIR not found!"; exit 1; }

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

# --- 3. Capture ---
# 2-second warmup (-t 2000) helps auto-exposure in cold 0°C light
echo "Capturing image..."
run_with_retry rpicam-jpeg -o capture.jpg -t 2000 -n --width 1920 --height 1080

# Move capture to the webcam folder
mkdir -p "$WEBCAM_DIR"
cp capture.jpg "$WEBCAM_DIR/$TIMESTAMP.jpg"

# --- 4. Retention Policy ---
echo "Pruning old images (keeping last $RETENTION_COUNT)..."
ls -1 "$WEBCAM_DIR"/[0-9]*.jpg 2>/dev/null | sort -r | tail -n +$((RETENTION_COUNT + 1)) | xargs -I {} rm -f {}

# --- 5. Git Push (No Pull) ---
# We stage changes, commit, and push directly.
git add -A "$WEBCAM_DIR/"
git commit -m "webcam update: $TIMESTAMP" || echo "Nothing to commit"

echo "Pushing to GitHub..."
run_with_retry git push origin main

# --- 6. Final Cleanup and Power Down ---
rm -f capture.jpg
sync  # Ensures all data is written to the SD card before power cut
echo "Workflow complete. Shutting down."

# The Witty Pi will see the shutdown and physically cut power
sudo shutdown -h now
