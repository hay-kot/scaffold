#!/bin/bash

# Source the assert_snapshot function
source tests/assert.sh

# accept bin as first argument


# Your script continues as before...
output=$($1 --log-level="error" \
    --output-dir=":memory:" \
    new \
    --preset="default" \
    --no-prompt \
    --snapshot="stdout" \
    nested)

# Call the function to assert the snapshot
assert_snapshot "nested.snapshot.txt" "$output"
