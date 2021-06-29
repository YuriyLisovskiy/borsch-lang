BINARY=bin/borsch

BUILD_TIME=`LC_ALL=uk_UA.utf8 date '+%b %d %Y, %T'`
LDFLAGS=-ldflags "-X 'github.com/YuriyLisovskiy/borsch/Borsch/cli/build.Time=${BUILD_TIME}'"

all: build

build:
	mkdir -p bin/ && \
	go build ${LDFLAGS} -o ${BINARY} lang/cli/main.go

install:
	mkdir -p /usr/local/lib/borsch-lang/Lib
	export BORSCH_STD="/usr/local/lib/borsch-lang/Lib"
	cp ${BINARY} /usr/local/bin/borsch

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
