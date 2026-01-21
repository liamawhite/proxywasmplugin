# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is an Istio WASM plugin built with Go 1.24 and the official proxy-wasm-go-sdk that logs inbound HTTP requests. **This plugin is for testing and validation purposes only - not for production use.**

## Build Commands

### Build WASM Plugin
```bash
make build
```
Compiles Go code to WASM binary (`plugin.wasm`) using `GOOS=wasip1 GOARCH=wasm CGO_ENABLED=0 go build -buildmode=c-shared`.

### Build OCI Image
```bash
make docker-build
```
Builds the WASM plugin and creates a Docker image.

### Run Unit Tests
```bash
make test
```

### Run Integration Tests
```bash
make test-integration
```
Runs integration tests using testcontainers-go against multiple Envoy versions (1.33.0, 1.34.0, 1.35.0).
Requires Docker to be running.

### Run All Tests
```bash
make test-all
```

### Clean Build Artifacts
```bash
make clean
```

## Architecture

### WASM Plugin Structure

The plugin implements the proxy-wasm ABI through three context types:

1. **vmContext**: Entry point, creates plugin contexts
2. **pluginContext**: Manages plugin lifecycle, creates HTTP contexts
3. **httpContext**: Intercepts HTTP requests via `OnHttpRequestHeaders()`

**Important**: The plugin uses `init()` for initialization (not `main()`), which is required for `-buildmode=c-shared`.

### HTTP Request Logging

The plugin extracts and logs three fields from HTTP request headers:
- Method: `:method` pseudo-header
- Path: `:path` pseudo-header
- Host: `:authority` pseudo-header (falls back to `host` header)

All logging uses `proxywasm.LogInfof()` and continues request processing with `types.ActionContinue`.

## WASM Compilation Requirements

- **SDK**: Uses official `github.com/proxy-wasm/proxy-wasm-go-sdk` (not Tetrate fork)
- **Go Version**: Must use Go 1.24+ with native WASI support
- **Build Mode**: Requires `-buildmode=c-shared` for proper WASM reactor module
- **Build Environment Variables**:
  - `GOOS=wasip1` (WASI preview 1 target)
  - `GOARCH=wasm`
  - `CGO_ENABLED=0`
- **Do not use TinyGo** - this project uses standard Go compiler
- **Envoy Requirement**: Requires Envoy >= 1.33.0 for WASI support

## Integration Testing

The project includes comprehensive integration tests using testcontainers-go:

- **Test Matrix**: Validates against Envoy 1.33.0, 1.34.0, and 1.35.0
- **Test Scenarios**: 5 HTTP request patterns per Envoy version
- **Validation**: Verifies plugin loads successfully and logs correct request data
- **Files**:
  - `integration_test.go` - Main test suite
  - `internal/testutil/` - Helper utilities for container management
  - `testdata/envoy-config.yaml` - Envoy bootstrap configuration
- **CI/CD**: Tests run automatically on every PR via GitHub Actions

## CI/CD

### Build Workflow (.github/workflows/build.yml)
- **Build Job**: Validates compilation and uploads WASM artifact
- **Test Job**: Runs integration tests against multiple Envoy versions
- Runs on every push to main and all pull requests
- Tests run in parallel for faster feedback

### Publish Workflow (.github/workflows/publish.yml)
- Triggered by: pushes to `main` or version tags (`v*.*.*`)
- Publishes OCI image to GHCR using Docker
- Tags: `latest` for main branch, semantic version for tags

## Deployment

Deploy to Istio using `WasmPlugin` custom resource (see `examples/wasmplugin.yaml`):
```yaml
apiVersion: extensions.istio.io/v1alpha1
kind: WasmPlugin
spec:
  url: oci://ghcr.io/liamawhite/proxywasmplugin:latest
```

Logs appear in Envoy proxy logs as:
```
HTTP Request: method=GET path=/api/users host=example.com
```
