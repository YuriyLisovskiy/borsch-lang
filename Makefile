all: build

build:
	mkdir -p bin/ && \
	go build -o bin/borsch src/cmd/main.go
