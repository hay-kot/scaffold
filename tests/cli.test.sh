go run main.go --log-level="error" \
    new \
    --preset="default" \
    --no-prompt \
    --snapshot="stdout" \
    cli

# Run the command and store the output in a variable
output=$(go run ./gen/scaffold-test*/main.go hello)

expected_output="colors=red, green description=This is a test description"

if [ "$output" = "$expected_output" ]; then
    echo "Test passed: output matches the expected string."
    exit 0
else
    echo "Expected: $expected_output"
    echo "Got:      $output"
    echo "Test failed: output does not match the expected string."
    exit 1
fi
