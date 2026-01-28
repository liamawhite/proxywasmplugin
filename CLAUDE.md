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
Runs integration tests using testcontainers-go against Envoy (defaults to v1.37.0).
Requires Docker to be running.

To test against a specific Envoy version:
```bash
ENVOY_VERSION=v1.35.0 make test-integration
```

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

- **Local Testing**: Tests against latest Envoy version (v1.37.0) by default
- **CI Matrix**: GitHub Actions tests against Envoy 1.33.0, 1.34.0, 1.35.0, 1.36.0, and 1.37.0 in parallel
- **Test Scenarios**: 5 HTTP request patterns (GET, POST with query, PUT, DELETE, root path)
- **Validation**: Verifies plugin loads successfully and logs correct request data
- **Environment Variable**: Set `ENVOY_VERSION` to test against a specific version
- **Files**:
  - `integration_test.go` - Main test suite
  - `internal/testutil/` - Helper utilities for container management and version selection
  - `testdata/envoy-config.yaml` - Envoy bootstrap configuration

## CI/CD

### Build and Publish Workflow (.github/workflows/build.yml)

This workflow runs on every push to main, version tags, and all pull requests.

**Jobs:**
1. **Build Job**: Validates compilation and uploads WASM artifact
2. **Test Job**: Matrix strategy runs integration tests against 5 Envoy versions in parallel (v1.33.0-v1.37.0)
   - Each version runs independently for faster feedback
   - Set to `fail-fast: false` so all versions complete even if one fails
3. **Publish Job**: Publishes OCI image to GHCR (only runs on push to main/tags, requires both build and all test jobs to pass)

**Publishing:**
- Only triggers on pushes to `main` or version tags (`v*.*.*`), not on PRs
- Requires build job and all 5 test matrix jobs to pass
- Tags: `main` for main branch pushes, semantic version for tags (e.g., `1.0.0`, `1.0`)

## Deployment

Deploy to Istio using `WasmPlugin` custom resource (see `examples/wasmplugin.yaml`):
```yaml
apiVersion: extensions.istio.io/v1alpha1
kind: WasmPlugin
spec:
  url: oci://ghcr.io/liamawhite/proxywasmplugin:main
```

Logs appear in Envoy proxy logs as:
```
HTTP Request: method=GET path=/api/users host=example.com
```
