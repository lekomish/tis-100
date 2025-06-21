SOURCE_CODE=tis-100
BINARY_NAME=tis-100
BUILD_DIR=bin

.PHONY: help lint test pre-commit clean build run

.DEFAULT_GOAL := help

help:
	@echo "Usage: make <target>"
	@echo "Targets:"
	@echo "  help              					Display this help message"
	@echo "  lint              					Run linters"
	@echo "  test              					Run tests"
	@echo "  pre-commit        					Run pre-commit checks"
	@echo "  clean             					Remove build artifacts"
	@echo "  build             					Build the application"
	@echo "  run ARGS=<diskpartition> 	Run the application"

lint:
	go mod tidy
	gofmt -s -w .
	golangci-lint run

test:
	go test -race -timeout 30s ./...

pre-commit:
	pre-commit run --all-files --hook-stage pre-push

clean:
	rm -rf $(BUILD_DIR)

build:
	mkdir -p $(BUILD_DIR)
	go build -o ./$(BUILD_DIR)/$(BINARY_NAME) ./cmd/$(SOURCE_CODE)

run: clean build
	./$(BUILD_DIR)/$(BINARY_NAME) $(ARGS)
