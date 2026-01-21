package testutil

// EnvoyVersion represents a specific Envoy version to test against.
type EnvoyVersion struct {
	// Name is a human-readable identifier for the version
	Name string
	// Image is the Docker image name (e.g., "envoyproxy/envoy")
	Image string
	// Tag is the Docker image tag (e.g., "v1.31.2")
	Tag string
	// Source indicates whether this is a standalone Envoy or Istio-bundled version
	Source string
}

// TestVersions defines the matrix of Envoy versions to test against.
// This list includes both standalone Envoy releases and Istio-bundled versions.
// NOTE: proxy-wasm-go-sdk requires Envoy >= 1.33.0 for WASI support
var TestVersions = []EnvoyVersion{
	// Standalone Envoy versions (>= 1.33.0 required for WASI support)
	{
		Name:   "Envoy 1.33.0",
		Image:  "envoyproxy/envoy",
		Tag:    "v1.33.0",
		Source: "standalone",
	},
	{
		Name:   "Envoy 1.34.0",
		Image:  "envoyproxy/envoy",
		Tag:    "v1.34.0",
		Source: "standalone",
	},
	{
		Name:   "Envoy 1.35.0",
		Image:  "envoyproxy/envoy",
		Tag:    "v1.35.0",
		Source: "standalone",
	},
	// Istio-bundled Envoy versions (commented out for initial implementation)
	// Uncomment these once basic tests are working with standalone versions
	// {
	// 	Name:   "Istio 1.26 (Envoy 1.34)",
	// 	Image:  "istio/proxyv2",
	// 	Tag:    "1.26.8",
	// 	Source: "istio",
	// },
	// {
	// 	Name:   "Istio 1.27 (Envoy 1.35)",
	// 	Image:  "istio/proxyv2",
	// 	Tag:    "1.27.5",
	// 	Source: "istio",
	// },
	// {
	// 	Name:   "Istio 1.28 (Envoy 1.36)",
	// 	Image:  "istio/proxyv2",
	// 	Tag:    "1.28.2",
	// 	Source: "istio",
	// },
}

// FullImageName returns the complete Docker image reference (image:tag).
func (v EnvoyVersion) FullImageName() string {
	return v.Image + ":" + v.Tag
}
