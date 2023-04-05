#!/bin/bash

TEST_DIR=$1
for i in $(find $TEST_DIR -name 'тест_*.борщ'); do
    go run "$i"
    exit_code=$?
    [ $exit_code -eq 0 ] || exit $exit_code
    echo "$i" - Success
done;
