#!/bin/bash

# Source the assert_snapshot function
source tests/assert.sh

# Your script continues as before...
output=$(go run main.go --log-level="error" \
    new \
    --preset="default" \
    --no-prompt \
    --snapshot="stdout" \
    cli)

# Call the function to assert the snapshot
assert_snapshot "cli.snapshot.txt" "$output"