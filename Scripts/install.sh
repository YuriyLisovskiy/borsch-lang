#!/bin/bash

APP_NAME=borsch
LIB_DIR=/usr/local/lib/${APP_NAME}-lang
BIN_DIR=/usr/local/bin

mkdir -p ${LIB_DIR} && \
cp -R Lib/ ${LIB_DIR} && \
mkdir -p ${BIN_DIR} && \
cp bin/${APP_NAME} ${BIN_DIR}/${APP_NAME} && \
export BORSCH_LIB="${LIB_DIR}/Lib" && \
echo "Append 'BORSCH_LIB=${LIB_DIR}/Lib' to your ~/.bash_profile for using this variable permanently."
