# Variables
APP_NAME := web-analyzer
BUILD_DIR := build
MAIN_FILE := ./cmd/webanalyzer/main.go

# Default target
.PHONY: all
all: build

# Install dependencies
.PHONY: deps
deps:
	go mod tidy
	go mod download

# Build the application
.PHONY: build
build: deps
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_FILE)

# Run the application
.PHONY: run
run: build
	$(BUILD_DIR)/$(APP_NAME)

# Clean up build artifacts
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)

# Run tests
.PHONY: test
test:
	go test ./... -v