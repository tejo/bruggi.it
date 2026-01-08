#!/bin/bash

# Clean exit on shutdown signals
trap "echo 'Shutdown signal received! Exiting cleanly.'; exit" SIGTERM SIGINT

# --- Configuration ---
REPO_DIR="/home/teo/bruggi.it"
WEBCAM_DIR="static/webcam"
WITTY_PATH="/home/pi/wittypi"  # Ensure this path is correct
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

# --- 3. Sync Remote Changes ---
echo "Syncing remote changes (Sparse Pull)..."
run_with_retry git pull origin main --rebase

# --- 4. Capture & Metadata Injection ---
echo "Capturing image..."
run_with_retry rpicam-jpeg -o capture.jpg -t 2000 -n --width 1920 --height 1080 --autofocus-mode manual --lens-position 0

if [ -f "capture.jpg" ]; then
    # Get Temperature from Witty Pi
    # The 'awk' command isolates the numeric value (e.g., 24.5)
    TEMP=$(sudo "$WITTY_PATH/wittyPi.sh" get_temperature | awk '{print $3}')
    
    if [ -n "$TEMP" ]; then
        echo "Recording Witty Pi Temperature: $TEMP°C"
        # Check if exiftool is installed before trying to use it
        if command -v exiftool >/dev/null 2>&1; then
            exiftool -UserComment="WittyPi_Temp: $TEMP C" -overwrite_original capture.jpg
        else
            echo "Warning: exiftool not found. Skipping metadata injection. Install with:  sudo apt update && sudo apt install libimage-exiftool-perl -y"
        fi
    fi
fi

# Move capture to the webcam folder
mkdir -p "$WEBCAM_DIR"
mv capture.jpg "$WEBCAM_DIR/$TIMESTAMP.jpg"

# --- 5. Retention Policy ---
echo "Pruning old images (keeping last $RETENTION_COUNT)..."
ls -1 "$WEBCAM_DIR"/[0-9]*.jpg 2>/dev/null | sort -r | tail -n +$((RETENTION_COUNT + 1)) | xargs -I {} rm -f {}

# --- 6. Git Push ---
git add "$WEBCAM_DIR"
git commit -m "webcam update: $TIMESTAMP (Temp: $TEMP C)" || echo "Nothing to commit"

echo "Pushing to GitHub..."
run_with_retry git push origin main

# --- 7. Final Cleanup and Power Down ---
sync
echo "Workflow complete. Shutting down."
sudo shutdown -h now
