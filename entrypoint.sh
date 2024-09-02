#!/bin/sh

# Function to terminate processes
terminate() {
    echo "Received SIGTERM or SIGINT. Shutting down..."
    kill -TERM "$api_pid" 2>/dev/null
    kill -TERM "$cron_pid" 2>/dev/null
    wait "$api_pid" 2>/dev/null
    wait "$cron_pid" 2>/dev/null
    exit 0
}

# Set up signal handling
trap terminate TERM INT

# Start the API in the background
./api &
api_pid=$!

# Start the cron process in the background
./cron &
cron_pid=$!

# Wait for any process to exit
wait -n

# Exit with status of process that exited first
exit $?
