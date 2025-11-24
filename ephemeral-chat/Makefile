.PHONY: build run test test-timeout test-expiration clean

# Build the application
build:
	go build -o ephemeral-chat-server .

# Run the server
run: build
	./ephemeral-chat-server

# Run basic API tests
test:
	./test_api.sh

# Test presence timeout functionality
test-timeout:
	./test_timeout.sh

# Test message expiration
test-expiration:
	./test_expiration.sh

# Run all tests
test-all: test test-expiration

# Clean build artifacts
clean:
	rm -f ephemeral-chat-server

# Install dependencies
deps:
	go mod tidy
	go mod download

# Development server with auto-reload (requires air)
dev:
	air

# Format code
fmt:
	go fmt ./...