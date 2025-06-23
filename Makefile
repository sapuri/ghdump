BINARY_NAME=ghdump

.PHONY: build

build:
	go build -o $(BINARY_NAME) .
