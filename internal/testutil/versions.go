package testutil

import (
	"os"
	"strings"
)

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

// GetEnvoyVersion returns the Envoy version to test against.
// If ENVOY_VERSION environment variable is set, uses that version.
// Otherwise, defaults to the latest supported version (1.37.0).
// This allows GitHub Actions matrix to test multiple versions in parallel.
func GetEnvoyVersion() EnvoyVersion {
	envoyVersion := os.Getenv("ENVOY_VERSION")
	if envoyVersion == "" {
		envoyVersion = "v1.37.0" // Default to latest
	}

	return EnvoyVersion{
		Name:   "Envoy " + strings.TrimPrefix(envoyVersion, "v"),
		Image:  "envoyproxy/envoy",
		Tag:    envoyVersion,
		Source: "standalone",
	}
}

// TestVersions defines the matrix of Envoy versions to test against in CI.
// NOTE: proxy-wasm-go-sdk requires Envoy >= 1.33.0 for WASI support
// In local testing, only the latest version is used (via GetEnvoyVersion).
// In CI, GitHub Actions matrix runs tests against all these versions in parallel.
var TestVersions = []EnvoyVersion{
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
	{
		Name:   "Envoy 1.36.0",
		Image:  "envoyproxy/envoy",
		Tag:    "v1.36.0",
		Source: "standalone",
	},
	{
		Name:   "Envoy 1.37.0",
		Image:  "envoyproxy/envoy",
		Tag:    "v1.37.0",
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
