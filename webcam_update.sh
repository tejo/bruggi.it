#!/bin/bash

# Configuration
REPO_DIR="/home/teo/bruggi.it"
MAX_RETRIES=5
RETRY_DELAY=10 # seconds

# Function to run commands with a retry limit
run_with_retry() {
    local n=1
    until [ $n -ge $MAX_RETRIES ]
    do
        "$@" && break || {
            echo "Command '$1' failed. Attempt $n/$MAX_RETRIES..."
            ((n++))
            sleep $RETRY_DELAY
        }
    done

    if [ $n -eq $MAX_RETRIES ]; then
        echo "Command '$1' failed after $MAX_RETRIES attempts. Exiting."
        exit 1
    fi
}

# 1. Navigate to directory
cd "$REPO_DIR" || exit 1

# 2. Sync with remote (Retry if network is flaky)
run_with_retry git pull origin main

# 3. Capture image (Retry if camera is busy/initializing)
run_with_retry rpicam-jpeg -o capture.jpg -t 1000

# 4. Process image
run_with_retry ./bin/bruggi-arm64 -update-webcam  capture.jpg

# 5. Git workflow
git add dist/ static/
# We don't retry commit; if there's nothing to change, it returns 1
git commit -m "webcam capture: $(date)" || echo "Nothing to commit"

# 6. Push changes (Retry for network)
run_with_retry git push origin main

# 7. remove image
rm -f capture.jpg
