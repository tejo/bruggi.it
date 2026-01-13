#!/bin/bash
# Clean exit on shutdown signals
trap "echo 'Shutdown signal received! Exiting cleanly.'; exit" SIGTERM SIGINT

# --- 1. Argument Parsing ---
SHUTDOWN_ENABLED=true
for arg in "$@"; do
    case $arg in
        shutdown=false|--no-shutdown)
            SHUTDOWN_ENABLED=false
            shift
            ;;
    esac
done

# --- Configuration ---
REPO_DIR="/home/teo/bruggi.it"
WEBCAM_DIR="static/webcam"
WITTY_PATH="/home/teo/wittypi"
TIMESTAMP=$(date +"%Y-%m-%d_%H-%M-%S")
RETENTION_COUNT=20
MAX_RETRIES=3
RETRY_DELAY=5

# --- 2. Wait for Internet Time Sync ---
echo "Waiting for internet time sync..."
for i in {1..30}; do
    if timedatectl status | grep -q "System clock synchronized: yes"; then
        echo "Time synchronized: $(date)"
        break
    fi
    echo "Waiting for NTP... ($i/30)"
    sleep 2
done

# --- 3. Operating Hours Check (06:30 to 20:00) ---
# We use HHMM format (e.g., 0630 or 2000)
CURRENT_TIME=$(date +"%H%M")
START_TIME="0630"
END_TIME="2000"

# Note: we use 10# to force Bash to treat the string as a decimal number (avoids octal errors with leading zeros)
if (( 10#$CURRENT_TIME < 10#$START_TIME )) || (( 10#$CURRENT_TIME > 10#$END_TIME )); then
    echo "Current time ($CURRENT_TIME) is outside operating hours ($START_TIME - $END_TIME)."
    echo "Skipping execution."
    # Skip to the shutdown logic at the end
else
    # --- 4. Main Workflow ---
    echo "Within operating hours. Starting capture workflow..."

    cd "$REPO_DIR" || { echo "Directory $REPO_DIR not found!"; exit 1; }

    # Ensure Sparse Checkout is active
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

    echo "Syncing remote changes..."
    run_with_retry git pull origin main --rebase

    echo "Capturing image..."
    run_with_retry rpicam-jpeg -o capture.jpg -t 2000 -n --width 1920 --height 1080 --autofocus-mode manual --lens-position 0

    if [ -f "capture.jpg" ]; then
        echo "Getting temperature from Witty Pi..."
        TEMP=$(echo "13" | sudo "$WITTY_PATH/wittyPi.sh" | grep "Current temperature" | awk '{print $4}')

        if [ -n "$TEMP" ]; then
            if command -v exiftool >/dev/null 2>&1; then
                exiftool -UserComment="WittyPi_Temp: $TEMP" -overwrite_original capture.jpg
            fi
        fi
    fi

    mkdir -p "$WEBCAM_DIR"
    mv capture.jpg "$WEBCAM_DIR/$TIMESTAMP.jpg"

    echo "Pruning old images..."
    ls -1 "$WEBCAM_DIR"/[0-9]*.jpg 2>/dev/null | sort -r | tail -n +$((RETENTION_COUNT + 1)) | xargs -I {} rm -f {}

    git add "$WEBCAM_DIR"
    git commit -m "webcam update: $TIMESTAMP (Temp: $TEMP)" || echo "Nothing to commit"
    echo "Pushing to GitHub..."
    run_with_retry git push origin main
fi

# --- 5. Final Cleanup and Power Down ---
sync

if [ "$SHUTDOWN_ENABLED" = true ]; then
    echo "Workflow complete/terminated. Shutting down."
    sudo shutdown -h now
else
    echo "Workflow complete/terminated. Shutdown skipped (shutdown=false)."
fi
