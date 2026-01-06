#!/bin/bash

# Clean exit on shutdown signals from Witty Pi
trap "echo 'Shutdown signal received! Exiting cleanly.'; exit" SIGTERM SIGINT

# --- Configuration ---
REPO_DIR="/home/teo/bruggi.it"
WEBCAM_DIR="static/webcam"
TIMESTAMP=$(date +"%Y-%m-%d_%H-%M-%S")
RETENTION_COUNT=20
MAX_RETRIES=3
RETRY_DELAY=5

# --- 1. Wait for Internet Time Sync ---
# Since Witty Pi RTC sync is disabled, we must wait for NTP 
# to prevent Git SSL errors and incorrect timestamps.
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

if [ "$SYNCED" = false ]; then
    echo "Warning: Time sync failed. Proceeding with local time."
fi

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

# --- 3. Self-Healing Git Check ---
if ! git rev-parse --is-inside-work-tree > /dev/null 2>&1; then
    echo "Git corruption detected. Re-cloning..."
    cd ..
    rm -rf bruggi.it
    # Note: Replace with your actual SSH URL for passwordless push
    git clone --filter=blob:none --sparse git@github.com:TEO_USER/bruggi.it.git
    cd bruggi.it
    git sparse-checkout set static
fi

# --- 4. Capture & Sync ---
# Using a 2-second warmup (-t 2000) for better white balance in 0°C light
echo "Capturing image..."
run_with_retry rpicam-jpeg -o capture.jpg -t 2000 -n --width 1920 --height 1080

# Sync with GitHub before adding new files
run_with_retry git fetch origin
run_with_retry git reset --hard origin/main

# Move capture to the webcam folder
mkdir -p "$WEBCAM_DIR"
cp capture.jpg "$WEBCAM_DIR/$TIMESTAMP.jpg"

# --- 5. Retention Policy ---
echo "Pruning old images (keeping last $RETENTION_COUNT)..."
ls -1 "$WEBCAM_DIR"/[0-9]*.jpg 2>/dev/null | sort -r | tail -n +$((RETENTION_COUNT + 1)) | xargs -I {} rm -f {}

# --- 6. Git Push ---
git add -A "$WEBCAM_DIR/"
git commit -m "webcam update: $TIMESTAMP" || echo "Nothing to commit"

if ! run_with_retry git push origin main; then
    echo "Push failed. Attempting force-sync..."
    git pull --rebase origin main
    git push origin main
fi

# --- 7. Final Cleanup and Power Down ---
rm -f capture.jpg
sync  # Critical: Flush all buffered data to the SD card before power cut
echo "Workflow complete. Notifying Witty Pi and shutting down."

# The Witty Pi daemon will see the shutdown and physically cut power
# sudo shutdown -h now
