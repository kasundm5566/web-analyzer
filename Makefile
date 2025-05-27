# Variables
APP_NAME := web-analyzer
BUILD_DIR := build
P_DIR := -p
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
	 @if not exist $(BUILD_DIR) (mkdir $(BUILD_DIR))
	 go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_FILE)

# Run the application
.PHONY: run
run: build
	$(BUILD_DIR)/$(APP_NAME)

# Clean up build artifacts
.PHONY: clean
clean:
	@if exist $(BUILD_DIR) (rmdir /s /q $(BUILD_DIR))

# Run tests
.PHONY: test
test:
	go test ./... -v