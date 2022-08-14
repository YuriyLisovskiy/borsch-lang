.PHONY: all rebuild build install uninstall clean test

APP_NAME=borsch
BINARY=bin/${APP_NAME}
LIB_DIR=/usr/local/lib/${APP_NAME}-lang
BIN_DIR=/usr/local/bin

BUILD_TIME=`LC_ALL=uk_UA.utf8 date '+%b %d %Y, %T'`
LDFLAGS=-ldflags "-X 'github.com/YuriyLisovskiy/borsch-lang/Borsch/cli/build.Time=${BUILD_TIME}'"

all: build

rebuild: clean build

build:
	mkdir -p bin/ && \
	go build ${LDFLAGS} -o ${BINARY} Borsch/cli/main.go

install:
	mkdir -p ${LIB_DIR}
	cp -R Lib/ ${LIB_DIR}
	export BORSCH_LIB="${LIB_DIR}/Lib"
	cp ${BINARY} ${BIN_DIR}/${APP_NAME}

uninstall:
	rm -rf ${LIB_DIR}
	rm ${BIN_DIR}/${APP_NAME}

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

test:
	go run Borsch/cli/main.go Test/вбудовані/типи/дійсний.борщ
	go run Borsch/cli/main.go Test/вбудовані/типи/логічний.борщ
	go run Borsch/cli/main.go Test/вбудовані/типи/рядок.борщ
	go run Borsch/cli/main.go Test/вбудовані/типи/цілий.борщ
