# Simple Makefile for a Go project

# Build the application
all: build

build:
	@echo "Building..."
	@templ generate

	@go build -o bin/main cmd/api/main.go

build-prod:
	@echo "Building PROD..."
	@templ generate

	@go build -a -ldflags "-s -w" -o bin/main cmd/api/main.go
	#upx --best --lzma bin/main
	#upx --ultra-brute bin/main

# Run the application
run:
	@go run cmd/api/main.go



# Test the application
test:
	@echo "Testing..."
	@go test ./tests -v

# Quality check
lint:
	@golangci-lint run --enable-all


# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f bin/main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
	    air; \
	    echo "Watching...";\
	else \
	    read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
	    if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
	        go install github.com/air-verse/air@latest; \
	        air; \
	        echo "Watching...";\
	    else \
	        echo "You chose not to install air. Exiting..."; \
	        exit 1; \
	    fi; \
	fi

.PHONY: all build run test clean
