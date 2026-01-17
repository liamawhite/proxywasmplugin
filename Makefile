.PHONY: build clean test docker-build

# Output file
OUTPUT := plugin.wasm

# Build the WASM plugin using Go 1.24
build:
	@echo "Building WASM plugin..."
	GOOS=wasip1 GOARCH=wasm CGO_ENABLED=0 go build -o $(OUTPUT) .
	@echo "WASM plugin built successfully: $(OUTPUT)"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(OUTPUT)
	@echo "Clean complete"

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Build OCI image with WASM artifact
docker-build: build
	@echo "Building OCI image..."
	docker build -t proxywasmplugin:latest .
	@echo "OCI image built successfully"
