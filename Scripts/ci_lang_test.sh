#!/bin/bash

INTERPRETER_BIN=$1
TEST_DIR=$2
for i in $(find $TEST_DIR -name 'тест_*.борщ'); do
    $INTERPRETER_BIN run --file "$i"
    exit_code=$?
    [ $exit_code -eq 0 ] || exit $exit_code
    echo "$i" - Success
done;
