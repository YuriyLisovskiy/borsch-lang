.PHONY: all build install uninstall test

APP_NAME=borsch
BUILD_TIME=`LC_ALL=uk_UA.utf8 date '+%b %d %Y, %T'`
LDFLAGS="-X '${ROOT_PACKAGE_NAME}/cli/build.Time=${BUILD_TIME}'"

all: build test

build:
	@mkdir -p bin
	@go build -ldflags ${LDFLAGS} -o bin/${APP_NAME} Borsch/cli/main.go

install:
	@bash ./Scripts/install.sh

uninstall:
	@bash ./Scripts/uninstall.sh

test:
	@go run Borsch/cli/main.go Test/вбудовані/типи/дійсний.борщ
	@go run Borsch/cli/main.go Test/вбудовані/типи/логічний.борщ
	@go run Borsch/cli/main.go Test/вбудовані/типи/рядок.борщ
	@go run Borsch/cli/main.go Test/вбудовані/типи/цілий.борщ
