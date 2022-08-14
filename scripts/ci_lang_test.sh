#!/bin/bash

DIR=$1
for i in $(find $DIR -name '*.борщ'); do
    ./bin/borsch "$i"
done;
