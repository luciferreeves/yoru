BINARY_NAME=yoru
BUILD_PATH=bin/$(BINARY_NAME)
MAIN_PATH=$(BINARY_NAME)/main.go

VERSION ?= $(shell git describe --tags 2>/dev/null || echo "dev")
BUILD_DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS := -ldflags "-X $(BINARY_NAME)/build.Version=$(VERSION) -X $(BINARY_NAME)/build.Date=$(BUILD_DATE)"

.PHONY: clean modules build run dev all

clean:
	@echo "Cleaning up..."
	@rm -rf bin
	@echo "Cleanup complete."

modules:
	@echo "Tidying Go modules..."
	@go mod tidy
	@echo "Go modules tidied."

build:
	@echo "Building version $(VERSION) at $(BUILD_DATE)..."
	@mkdir -p bin
	@go build $(LDFLAGS) -o $(BUILD_PATH) $(MAIN_PATH) || true
	@echo "Build complete."

run:
	@if [ ! -f $(BUILD_PATH) ]; then echo "Binary not found. Building binary..."; $(MAKE) -s build; fi
	@echo "Running..."
	@$(BUILD_PATH) || true

dev:
	@echo "Running in development mode with version $(VERSION)..."
	@go run $(LDFLAGS) $(MAIN_PATH) || true

all: clean modules build run

.SILENT: