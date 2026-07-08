BINARY_NAME=ghdump

.PHONY: build test

build:
	go build -o $(BINARY_NAME) .

test:
	go test -race ./...
