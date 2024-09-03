# Project settings
BINARY_NAME = go-db-lite
BUILD_DIR = bin
CMD_MAIN = cmd/main.go
SRC_FILES = $(shell find internal -type f -name '*.go') $(CMD_MAIN)

# Go settings
GO_CMD = go
GO_FMT = gofmt
GO_TEST = go test
GO_BUILD = $(GO_CMD) build
GO_CLEAN = $(GO_CMD) clean
GO_MOD_TIDY = $(GO_CMD) mod tidy

# Output binary
BIN_PATH = $(BUILD_DIR)/$(BINARY_NAME)

# Default target
.PHONY: all
all: tidy build

# Build the Go binary
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	$(GO_BUILD) -o $(BIN_PATH) $(CMD_MAIN)

# Run the application
.PHONY: run
run: build
	@echo "Running $(BINARY_NAME)..."
	$(BIN_PATH)

# Clean the build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	$(GO_CLEAN)

# Format the Go code
.PHONY: format
format:
	@echo "Formatting Go code..."
	$(GO_FMT) -w $(SRC_FILES)

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	$(GO_TEST) ./...

# Tidy up the Go modules
.PHONY: tidy
tidy:
	@echo "Tidying Go modules..."
	$(GO_MOD_TIDY)

# Rebuild everything from scratch
.PHONY: rebuild
rebuild: clean all

# Convenience shortcut for formatting, tidying, and testing
.PHONY: check
check: format tidy test