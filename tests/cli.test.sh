#!/bin/bash 
set -euo pipefail 
export SCAFFOLD_NO_CLOBBER="true" 
export SCAFFOLD_OUT="gen" 
export SCAFFOLD_DIR=".scaffold,.examples" 

rm -rf gen/*

go run main.go --log-level="error" new --preset="default" --no-prompt --snapshot="stdout" cli

# Run the command and store the output in a variable
output=$(go run ./gen/scaffold-test*/main.go hello)

echo "Output: '$output'"  
# Define the expected output
expected_output="colors=red, green description=This is a test description"

# Compare the actual output with the expected output
if [ "$output" = "$expected_output" ]; then
    echo "Test passed: Output matches the expected string."
    exit 0
else
    echo "Test failed: Output does not match the expected string."
    exit 1
fi
