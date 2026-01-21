package main

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/liamawhite/proxywasmplugin/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestScenario defines a single HTTP request test case.
type TestScenario struct {
	Name        string
	Method      string
	Path        string
	Host        string
	ExpectedLog string
}

// getTestScenarios returns the complete set of HTTP request scenarios to test.
func getTestScenarios() []TestScenario {
	return []TestScenario{
		{
			Name:        "Basic GET",
			Method:      "GET",
			Path:        "/api/users",
			Host:        "example.com",
			ExpectedLog: "method=GET path=/api/users host=example.com",
		},
		{
			Name:        "POST with query",
			Method:      "POST",
			Path:        "/submit?id=123",
			Host:        "api.service.local",
			ExpectedLog: "method=POST path=/submit?id=123 host=api.service.local",
		},
		{
			Name:        "PUT request",
			Method:      "PUT",
			Path:        "/resource/456",
			Host:        "backend.cluster",
			ExpectedLog: "method=PUT path=/resource/456 host=backend.cluster",
		},
		{
			Name:        "DELETE request",
			Method:      "DELETE",
			Path:        "/items/789",
			Host:        "localhost:8080",
			ExpectedLog: "method=DELETE path=/items/789 host=localhost:8080",
		},
		{
			Name:        "Root path",
			Method:      "GET",
			Path:        "/",
			Host:        "test.com",
			ExpectedLog: "method=GET path=/ host=test.com",
		},
	}
}

// TestWASMPluginCompatibility tests the WASM plugin against multiple Envoy versions.
// This integration test verifies that:
// 1. The plugin loads successfully without crashes
// 2. HTTP request headers are logged correctly (method, path, host)
func TestWASMPluginCompatibility(t *testing.T) {
	// Ensure plugin.wasm exists before running tests
	pluginPath := "plugin.wasm"
	require.FileExists(t, pluginPath, "plugin.wasm must be built before running integration tests (run 'make build')")

	// Get absolute path to test data
	configPath := filepath.Join("testdata", "envoy-config.yaml")
	require.FileExists(t, configPath, "envoy-config.yaml must exist in testdata directory")

	ctx := context.Background()

	// Test against each Envoy version
	for _, version := range testutil.TestVersions {
		version := version // Capture for parallel execution
		t.Run(version.Name, func(t *testing.T) {
			t.Parallel() // Run versions in parallel for faster execution

			// Start Envoy container with WASM plugin
			envoyContainer := testutil.StartEnvoyContainer(t, ctx, version, pluginPath, configPath)

			// Wait for Envoy to be fully ready
			err := envoyContainer.WaitForReady(ctx, 5*time.Second)
			require.NoError(t, err, "Envoy failed to become ready")

			// Get proxy port for sending requests
			proxyPort, err := envoyContainer.GetProxyPort(ctx)
			require.NoError(t, err, "Failed to get proxy port")

			host, err := envoyContainer.GetHost(ctx)
			require.NoError(t, err, "Failed to get container host")

			// Verify plugin loaded successfully
			assertPluginLoaded(t, envoyContainer)

			// Run all test scenarios
			scenarios := getTestScenarios()
			for _, scenario := range scenarios {
				scenario := scenario // Capture for closure
				t.Run(scenario.Name, func(t *testing.T) {
					// Send HTTP request to Envoy
					resp := sendRequest(t, host, proxyPort.Port(), scenario)

					// Verify response
					assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected 200 OK response")
					resp.Body.Close()

					// Wait briefly for logs to be captured
					time.Sleep(100 * time.Millisecond)

					// Verify expected log output
					assertLogContains(t, envoyContainer, scenario.ExpectedLog)
				})
			}

			// If test failed, dump all logs for debugging
			if t.Failed() {
				t.Logf("=== Envoy Container Logs ===\n%s\n=== End Logs ===", envoyContainer.DumpLogs())
			}
		})
	}
}

// sendRequest sends an HTTP request to the Envoy proxy with the specified scenario parameters.
func sendRequest(t *testing.T, host, port string, scenario TestScenario) *http.Response {
	t.Helper()

	// Build request URL
	url := fmt.Sprintf("http://%s:%s%s", host, port, scenario.Path)

	// Create HTTP request
	req, err := http.NewRequest(scenario.Method, url, nil)
	require.NoError(t, err, "Failed to create HTTP request")

	// Set Host header to match scenario
	req.Host = scenario.Host

	// Send request
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	require.NoError(t, err, "Failed to send HTTP request")

	return resp
}

// assertPluginLoaded verifies that the WASM plugin loaded successfully without errors.
func assertPluginLoaded(t *testing.T, envoyContainer *testutil.EnvoyContainer) {
	t.Helper()

	allLogs := strings.Join(envoyContainer.Logs.GetLogs(), "\n")

	// Should NOT contain WASM load errors
	assert.NotContains(t, allLogs, "wasm log: panic:", "Plugin panicked during load")
	assert.NotContains(t, allLogs, "Failed to load", "Plugin failed to load")
	assert.NotContains(t, allLogs, "Unable to create WASM", "WASM runtime failed to create plugin")

	// Should contain successful startup
	assert.Contains(t, allLogs, "starting main dispatch loop", "Envoy did not start successfully")
}

// assertLogContains verifies that the expected log pattern appears in the Envoy logs.
func assertLogContains(t *testing.T, envoyContainer *testutil.EnvoyContainer, expectedPattern string) {
	t.Helper()

	allLogs := strings.Join(envoyContainer.Logs.GetLogs(), "\n")

	// The plugin logs with format: "HTTP Request: method=X path=Y host=Z"
	pattern := fmt.Sprintf("HTTP Request: %s", expectedPattern)

	assert.Contains(t, allLogs, pattern,
		"Expected log pattern not found in Envoy output.\nLooking for: %s", pattern)
}
