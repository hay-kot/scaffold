#!/bin/bash 
export SCAFFOLD_NO_CLOBBER="true" 
export SCAFFOLD_OUT="gen" 
export SCAFFOLD_DIR=".scaffold,.examples" 

go run main.go new --test cli
go run ./gen/main.go hello

# Run the command and store the output in a variable
output=$(go run ./gen/main.go hello)

# Define the expected output
expected_output="Hello, your favorite colors are red, green"

# Compare the actual output with the expected output
if [ "$output" = "$expected_output" ]; then
    echo "Test passed: Output matches the expected string."
    exit 0
else
    echo "Test failed: Output does not match the expected string."
    exit 1
fi
