.PHONY: build clean test test-integration test-all docker-build

# Output file
OUTPUT := plugin.wasm

# Build the WASM plugin using Go 1.24
build:
	@echo "Building WASM plugin..."
	GOOS=wasip1 GOARCH=wasm CGO_ENABLED=0 go build -buildmode=c-shared -o $(OUTPUT) .
	@echo "WASM plugin built successfully: $(OUTPUT)"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(OUTPUT)
	@echo "Clean complete"

# Run unit tests
test:
	@echo "Running unit tests..."
	go test -v -short ./...

# Run integration tests (requires plugin.wasm to be built first)
test-integration: build
	@echo "Running integration tests..."
	go test -v -timeout 10m ./...

# Run all tests (unit + integration)
test-all: test test-integration
	@echo "All tests complete"

# Build OCI image with WASM artifact
docker-build: build
	@echo "Building OCI image..."
	docker build -t proxywasmplugin:latest .
	@echo "OCI image built successfully"
