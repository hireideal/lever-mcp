.PHONY: build test clean

BINARY=lever-mcp
BUILD_DIR=.

build:
	go build -o $(BUILD_DIR)/$(BINARY) ./cmd/lever-mcp

test:
	go test ./... -count=1

clean:
	rm -f $(BUILD_DIR)/$(BINARY)
