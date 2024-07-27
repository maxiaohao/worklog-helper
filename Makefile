BINARY_NAME=worklog-helper
BUILD_DIR=bin

all: build

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) main.go

run: build
	$(BUILD_DIR)/$(BINARY_NAME)

clean:
	rm -rf $(BUILD_DIR)

test:
	go test ./...

deps:
	go mod tidy

help:
	@echo "Makefile commands:"
	@echo "  all       - Build all"
	@echo "  build     - Build the app"
	@echo "  run       - Build and run"
	@echo "  clean     - Clean the build directory"
	@echo "  test      - Run tests"
	@echo "  deps      - Install dependencies"
	@echo "  help      - Show this help message"
