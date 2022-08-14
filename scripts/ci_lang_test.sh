#!/bin/bash

INTERPRETER_BIN=$1
TEST_DIR=$2
for i in $(find $TEST_DIR -name '*.борщ'); do
    $INTERPRETER_BIN "$i"
done;
