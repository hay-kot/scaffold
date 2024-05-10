#!/bin/bash
export SCAFFOLD_NO_CLOBBER="true"
export SCAFFOLD_OUT="gen"
export SCAFFOLD_DIR=".scaffold,.examples"

checkmark="✓"
crossmark="✗"

echo "Running Script Tests"
# run each test script in the tests directory
for test_script in tests/*.test.sh; do
    rm -rf ./gen

    # if exit code of script is 0, print checkmark, else print crossmark
    # and the output indented by 4 spaces
    output=$($test_script 2>&1)

    if [ $? -eq 0 ]; then
        echo "  $checkmark $test_script"
    else
        echo "  $crossmark $test_script"
        # Print each line of the output indented by 4 spaces
        echo "$output" | sed 's/^/      /'
    fi
done