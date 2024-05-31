#!/bin/bash

output=$(go test ./...)
result=$?

if [ $result -ne 0 ]; then
    echo "Some tests failed:"
    echo "$output"
    exit 1
else
    go test ./... -v
    exit 0
fi
