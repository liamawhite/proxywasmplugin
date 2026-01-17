# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is an Istio WASM plugin built with Go 1.24 and proxy-wasm-go-sdk that logs inbound HTTP requests. **This plugin is for testing and validation purposes only - not for production use.**

## Build Commands

### Build WASM Plugin
```bash
make build
```
Compiles Go code to WASM binary (`plugin.wasm`) using `GOOS=wasip1 GOARCH=wasm CGO_ENABLED=0`.

### Build OCI Image
```bash
make docker-build
```
Builds the WASM plugin and creates a Docker image.

### Run Tests
```bash
make test
# or for verbose output
go test -v ./...
```

### Clean Build Artifacts
```bash
make clean
```

## Architecture

### WASM Plugin Structure

The plugin implements the proxy-wasm ABI through three context types:

1. **vmContext** (main.go:14-20): Entry point, creates plugin contexts
2. **pluginContext** (main.go:22-28): Manages plugin lifecycle, creates HTTP contexts
3. **httpContext** (main.go:30-59): Intercepts HTTP requests via `OnHttpRequestHeaders()`

### HTTP Request Logging

The plugin extracts and logs three fields from HTTP request headers:
- Method: `:method` pseudo-header
- Path: `:path` pseudo-header
- Host: `:authority` pseudo-header (falls back to `host` header)

All logging uses `proxywasm.LogInfof()` and continues request processing with `types.ActionContinue`.

## WASM Compilation Requirements

- **Must use Go 1.24+** with native WASI support
- **Build environment variables**:
  - `GOOS=wasip1` (WASI preview 1 target)
  - `GOARCH=wasm`
  - `CGO_ENABLED=0`
- **Do not use TinyGo** - this project uses standard Go compiler

## CI/CD

### Build Workflow (.github/workflows/build.yml)
Runs on push/PR to validate compilation and uploads WASM artifact.

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
