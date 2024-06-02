#!/bin/bash

# Source the assert_snapshot function
source tests/assert.sh

# Your script continues as before...
output=$($1 --log-level="error" \
    --run-hooks="always" \
    new \
    --preset="default" \
    --no-prompt \
    --snapshot="stdout" \
    hooks)

# Call the function to assert the snapshot
assert_snapshot "hooks.snapshot.txt" "$output"
