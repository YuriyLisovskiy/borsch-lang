all: build

build:
	mkdir -p bin/ && \
	go build -o bin/borsch src/cmd/main.go

install:
	mkdir -p /usr/local/lib/borsch-lang/std
	export BORSCH_STD="/usr/local/lib/borsch-lang/std"
	cp bin/borsch /usr/local/bin/borsch
