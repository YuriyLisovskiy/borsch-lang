SHELL := /bin/bash

.PHONY: all build install uninstall lint-local test

PWD=$(shell pwd)
VERSION=$(shell git describe --tags --dirty --always)

BUILD_TIME=`LC_ALL=uk_UA.utf8 date '+%b %d %Y, %T'`
CONFIG_PACKAGE_NAME="github.com/YuriyLisovskiy/borsch-lang/internal/config"
MAIN_FILE_PATH="internal/cmd/main.go"

#LDFLAGS="-X '${ROOT_PACKAGE_NAME}/cli/build.Time=${BUILD_TIME}'"

all: build test

build:
	@mkdir -p bin
	@go build \
	  -ldflags "-X '${CONFIG_PACKAGE_NAME}.Version=$(shell git describe --tags --dirty --always)' \
	  			-X '${CONFIG_PACKAGE_NAME}.Time=${BUILD_TIME}' \
	  			-X '${CONFIG_PACKAGE_NAME}.License=$(cat LICENSE)'" \
	  -o bin/borsch internal/cmd/main.go

install:
	@bash ./Scripts/install.sh

uninstall:
	@bash ./Scripts/uninstall.sh

test:
	go test ./...
	go run ${MAIN_FILE_PATH} test/вбудовані/типи/тест_дійсний.борщ
	go run ${MAIN_FILE_PATH} test/вбудовані/типи/тест_логічний.борщ
	go run ${MAIN_FILE_PATH} test/вбудовані/типи/тест_рядок.борщ
	go run ${MAIN_FILE_PATH} test/вбудовані/типи/тест_цілий.борщ

lint-local:
	docker run --rm -v $(HOME)/go/pkg:/go/pkg -v $(PWD):/app -w /app golangci/golangci-lint:v1.50.1 golangci-lint run --max-issues-per-linter 100 --max-same-issues 10 -c .golangci.yml
