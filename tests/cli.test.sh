#!/bin/bash

# Source the assertions file
source tests/assert.sh

go run main.go --log-level="error" \
    new \
    --preset="default" \
    --no-prompt \
    --snapshot="stdout" \
    cli

# Run the command and store the output in a variable
output=$(go run ./gen/scaffold-test*/main.go hello)

expected_output="colors=red, green description=This is a test description"

# Call the function to compare output with expected output
assert_output "$output" "$expected_output"