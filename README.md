# Proxy WASM Plugin

An Istio WASM plugin built with Go 1.24 and [proxy-wasm-go-sdk](https://github.com/proxy-wasm/proxy-wasm-go-sdk) that logs inbound HTTP requests.

## Features

- Logs HTTP request method, path, and host for all inbound requests
- Minimal performance overhead
- Published as OCI artifact to GitHub Container Registry (GHCR)
- Compatible with Istio service mesh

## Prerequisites

- Go 1.24 or later
- Docker (for building OCI images)
- Istio installed in your Kubernetes cluster (for deployment)

## Building

### Build WASM Plugin

```bash
make build
```

This compiles the Go code to a WASM binary (`plugin.wasm`) using the WASI target.

### Build OCI Image

```bash
make docker-build
```

### Clean Build Artifacts

```bash
make clean
```

## Usage with Istio

The plugin is published to GHCR and can be deployed to Istio using the `WasmPlugin` custom resource.

### Deploy the Plugin

Apply the example configuration:

```bash
kubectl apply -f examples/wasmplugin.yaml
```

Or create your own `WasmPlugin` resource:

```yaml
apiVersion: extensions.istio.io/v1alpha1
kind: WasmPlugin
metadata:
  name: http-logger
  namespace: istio-system
spec:
  selector:
    matchLabels:
      istio: ingressgateway
  url: oci://ghcr.io/liamawhite/proxywasmplugin:latest
```

### Verify Logs

Once deployed, the plugin will log HTTP requests. Check the Envoy proxy logs:

```bash
kubectl logs -n istio-system -l istio=ingressgateway -c istio-proxy
```

You should see log entries like:

```
HTTP Request: method=GET path=/api/users host=example.com
HTTP Request: method=POST path=/api/login host=example.com
```

## Development

### Project Structure

```
.
├── main.go                    # Plugin implementation
├── go.mod                     # Go module definition
├── Makefile                   # Build automation
├── .github/workflows/
│   ├── build.yml             # CI build workflow
│   └── publish.yml           # GHCR publishing workflow
└── examples/
    └── wasmplugin.yaml       # Example Istio WasmPlugin CR
```

### How It Works

The plugin implements the proxy-wasm ABI using the Go SDK:

1. **VM Context**: Initializes the plugin
2. **Plugin Context**: Manages plugin lifecycle
3. **HTTP Context**: Intercepts HTTP requests
4. **OnHttpRequestHeaders**: Extracts and logs request details

The plugin logs the following for each HTTP request:
- HTTP method (`:method` pseudo-header)
- Request path (`:path` pseudo-header)
- Host (`:authority` or `host` header)

## Publishing New Versions

The GitHub Actions workflow automatically publishes to GHCR:

- **Pushes to `main`**: Published as `latest` tag
- **Git tags** (e.g., `v1.0.0`): Published with the version tag

To publish a new version:

```bash
git tag v1.0.0
git push origin v1.0.0
```

## License

See LICENSE file for details.
