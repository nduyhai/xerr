# Makefile for xerr library
# Note: This Makefile requires GNU Make, which is not installed by default on Windows.
# Windows users may need to install GNU Make or use alternative build tools.

# Go parameters
GOCMD=go
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
GOLINT=golangci-lint
GOIMPORTS=goimports

# Package path
PACKAGE=./...

.PHONY: all test test-coverage clean deps lint fmt goimports verify help

all: test lint fmt

# Run tests
test:
	$(GOTEST) -v $(PACKAGE)

# Run tests with coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out $(PACKAGE)
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Clean artifacts
clean:
	if exist coverage.out del coverage.out
	if exist coverage.html del coverage.html

# Install dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Run linter
lint:
	$(GOLINT) run

# Format code
fmt:
	$(GOCMD) fmt $(PACKAGE)

# Run goimports
goimports:
	@where $(GOIMPORTS) >nul 2>&1 || go install golang.org/x/tools/cmd/goimports@latest
	$(GOIMPORTS) -w .

# Verify dependencies
verify:
	$(GOMOD) verify

# Show help
help:
	@echo Make targets:
	@echo   all          - Run tests, lint, and format code
	@echo   test         - Run tests
	@echo   test-coverage - Run tests with coverage report
	@echo   clean        - Clean artifacts
	@echo   deps         - Install dependencies
	@echo   lint         - Run linter
	@echo   fmt          - Format code
	@echo   goimports    - Run goimports to format code and update imports
	@echo   verify       - Verify dependencies
	@echo   help         - Show this help
