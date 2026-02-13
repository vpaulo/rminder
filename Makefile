# Simple Makefile for a Go project

# Build the application
all: build

build:
	@echo "Building..."
	# @templ generate

	@go build -o bin/rminder cmd/rminder/main.go

build-prod:
	@echo "Building PROD..."
	# @templ generate

	@go build -a -ldflags "-s -w" -o bin/rminder cmd/rminder/main.go
	#upx --best --lzma bin/rminder
	#upx --ultra-brute bin/rminder

# Run the application
run:
	@go run cmd/rminder/main.go

package-rminder: build-prod
	@echo "Packaging rminder..."
	@mkdir -p package/rminder/usr/local/bin
	@cp bin/rminder package/rminder/usr/local/bin/rminder
	@mkdir -p package/rminder/var/lib/rminder
	@dpkg-deb --build package/rminder

package-rminder-caddy:
	@mkdir -p package/rminder-caddy/usr/local/bin
	@curl "https://caddyserver.com/api/download?os=linux&arch=amd64&idempotency=61754418096649" -o package/rminder-caddy/usr/local/bin/rminder-caddy
	@chmod +x package/rminder-caddy/usr/local/bin/rminder-caddy
	@mkdir -p package/rminder-caddy/var/lib/rminder-caddy
	@dpkg-deb --build package/rminder-caddy package/rminder-caddy.deb

package: package-rminder package-rminder-caddy

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
	@rm -f bin/rminder

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
