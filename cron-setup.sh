#!/bin/bash

# Setup cron job for daily fetch at 7:00 AM

PROJECT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
SCRIPT_PATH="$PROJECT_DIR/daily-fetch.sh"

# Check if crontab already has this entry
EXISTING=$(crontab -l 2>/dev/null | grep -F "$SCRIPT_PATH")

if [ -n "$EXISTING" ]; then
    echo "Cron job already exists:"
    echo "$EXISTING"
    exit 0
fi

# Add the cron job
# 0 7 * * * = every day at 7:00 AM
CRON_LINE="0 7 * * * cd $PROJECT_DIR && GITHUB_TOKEN=\$GITHUB_TOKEN $SCRIPT_PATH >> $PROJECT_DIR/cron.log 2>&1"

echo "Adding cron job:"
echo "$CRON_LINE"
echo ""

# Add to crontab
(crontab -l 2>/dev/null; echo "$CRON_LINE") | crontab -

echo "✅ Cron job added successfully!"
echo "The script will run daily at 7:00 AM"
echo "Output will be logged to: $PROJECT_DIR/cron.log"
