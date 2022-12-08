.PHONY: all rebuild build install install-lib uninstall clean

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

install-lib:
	mkdir -p ${LIB_DIR}
	cp -R Lib/ ${LIB_DIR}
	export BORSCH_STD="${LIB_DIR}/Lib"
	export BORSCH_LIB="${LIB_DIR}/Lib"

install: install-lib
install: 
	cp ${BINARY} ${BIN_DIR}/${APP_NAME}

uninstall:
	rm -rf ${LIB_DIR}
	rm ${BIN_DIR}/${APP_NAME}

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
