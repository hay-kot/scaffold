#!/bin/bash

# Function to compare output with expected output and print result
assert_output() {
    local output="$1"
    local expected_output="$2"

    if [ "$output" = "$expected_output" ]; then
        echo "Test passed: output matches the expected string."
        exit 0
    else
        echo "Expected: $expected_output"
        echo "Got:      $output"
        echo "Test failed: output does not match the expected string."
        exit 1
    fi
}

# Function to assert snapshots
assert_snapshot() {
    local snapshot_name="$1"
    local output="$2"
    local snapshots_dir="tests/snapshots"
    local snapshot_file="$snapshots_dir/$snapshot_name"

    # Create snapshots directory if it doesn't exist
    mkdir -p "$snapshots_dir"

    # Check if snapshot file exists
    if [ -f "$snapshot_file" ]; then
        # Compare current output with snapshot
        diff_result=$(diff -u "$snapshot_file" <(echo "$output"))
        if [ $? -eq 0 ]; then
            echo "Test passed: output matches the snapshot '$snapshot_name'."
            exit 0
        else
            echo "Test failed: output does not match the snapshot '$snapshot_name'. Diff:"
            echo "$diff_result"
            exit 1
        fi
    else
        # Snapshot does not exist, create it
        echo "$output" > "$snapshot_file"
        echo "Snapshot '$snapshot_name' created."
        exit 0
    fi
}