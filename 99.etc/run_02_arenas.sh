#!/bin/bash
# This script runs the Go arenas example with the required experiment flag.

# Ensure we are running from the project root or the script handles paths correctly
# Here we use the script's directory to locate the go file.

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
TARGET_FILE="$SCRIPT_DIR/02_arenas.go"

echo "___________________________________________________"
echo "Running: $TARGET_FILE"
echo "With:    GOEXPERIMENT=arenas"
echo "___________________________________________________"
echo ""

GOEXPERIMENT=arenas go run "$TARGET_FILE"
