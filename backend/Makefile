# Variables
APP_NAME = app
BUILD_DIR = build
SRC_DIR = .
GO_FILES = $(wildcard $(SRC_DIR)/*.go)

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build: $(GO_FILES)
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(SRC_DIR)

# Run tests
.PHONY: test
test:
	go test $(SRC_DIR)/...

# Clean build artifacts
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)

# Run the application
.PHONY: run
run: build
	./$(BUILD_DIR)/$(APP_NAME)

# Help message
.PHONY: help
help:
	@echo "Makefile commands:"
	@echo "  make build   - Build the application"
	@echo "  make test    - Run tests"
	@echo "  make clean   - Clean build artifacts"
	@echo "  make run     - Run the application"
	@echo "  make help    - Show this help message"