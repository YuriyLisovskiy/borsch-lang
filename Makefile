.PHONY: all install uninstall test

all: install test

install:
	@bash ./Scripts/install.sh

uninstall:
	@bash ./Scripts/uninstall.sh

test:
	@go run Borsch/cli/main.go Test/вбудовані/типи/дійсний.борщ
	@go run Borsch/cli/main.go Test/вбудовані/типи/логічний.борщ
	@go run Borsch/cli/main.go Test/вбудовані/типи/рядок.борщ
	@go run Borsch/cli/main.go Test/вбудовані/типи/цілий.борщ
