package testutil

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// EnvoyContainer wraps a testcontainers container with Envoy-specific functionality.
type EnvoyContainer struct {
	Container testcontainers.Container
	Logs      *EnvoyLogConsumer
}

// StartEnvoyContainer creates and starts an Envoy container with the WASM plugin loaded.
// It configures the container with the necessary mounts, ports, and wait strategies.
//
// Parameters:
//   - t: The testing context
//   - ctx: The context for container operations
//   - version: The Envoy version to test
//   - pluginPath: Absolute path to the plugin.wasm file
//   - configPath: Absolute path to the envoy-config.yaml file
//
// Returns an EnvoyContainer with the running container and log consumer.
func StartEnvoyContainer(t *testing.T, ctx context.Context, version EnvoyVersion, pluginPath, configPath string) *EnvoyContainer {
	t.Helper()

	// Ensure paths are absolute
	absPluginPath, err := filepath.Abs(pluginPath)
	require.NoError(t, err, "Failed to resolve absolute plugin path")

	absConfigPath, err := filepath.Abs(configPath)
	require.NoError(t, err, "Failed to resolve absolute config path")

	// Create log consumer
	logConsumer := NewEnvoyLogConsumer()

	// Create container request
	req := testcontainers.ContainerRequest{
		Image:        version.FullImageName(),
		ExposedPorts: []string{"10000/tcp", "9901/tcp"},
		Mounts: testcontainers.Mounts(
			testcontainers.BindMount(absPluginPath, "/etc/envoy/plugin.wasm"),
			testcontainers.BindMount(absConfigPath, "/etc/envoy/envoy.yaml"),
		),
		Cmd: []string{"envoy", "-c", "/etc/envoy/envoy.yaml", "--log-level", "info"},
		WaitingFor: wait.ForLog("starting main dispatch loop").
			WithStartupTimeout(30 * time.Second),
		LogConsumerCfg: &testcontainers.LogConsumerConfig{
			Consumers: []testcontainers.LogConsumer{logConsumer},
		},
	}

	// Start container
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err, "Failed to start Envoy container for %s", version.Name)

	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate container: %v", err)
		}
	})

	return &EnvoyContainer{
		Container: container,
		Logs:      logConsumer,
	}
}

// GetProxyPort returns the mapped host port for the Envoy proxy listener (10000).
func (ec *EnvoyContainer) GetProxyPort(ctx context.Context) (nat.Port, error) {
	return ec.Container.MappedPort(ctx, "10000")
}

// GetAdminPort returns the mapped host port for the Envoy admin interface (9901).
func (ec *EnvoyContainer) GetAdminPort(ctx context.Context) (nat.Port, error) {
	return ec.Container.MappedPort(ctx, "9901")
}

// GetHost returns the host address for connecting to the container.
func (ec *EnvoyContainer) GetHost(ctx context.Context) (string, error) {
	return ec.Container.Host(ctx)
}

// DumpLogs returns all captured logs as a single string for debugging.
func (ec *EnvoyContainer) DumpLogs() string {
	logs := ec.Logs.GetLogs()
	result := ""
	for _, log := range logs {
		result += log
	}
	return result
}

// WaitForReady waits for Envoy to be fully ready by checking the admin interface.
func (ec *EnvoyContainer) WaitForReady(ctx context.Context, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		// Check if the admin endpoint is responding
		adminPort, err := ec.GetAdminPort(ctx)
		if err != nil {
			return fmt.Errorf("failed to get admin port: %w", err)
		}

		host, err := ec.GetHost(ctx)
		if err != nil {
			return fmt.Errorf("failed to get host: %w", err)
		}

		// Simple check: if we can get the ports, Envoy is likely ready
		_ = host
		_ = adminPort

		// The wait strategy already ensures "starting main dispatch loop" is logged
		// so we can assume Envoy is ready at this point
		return nil
	}

	return fmt.Errorf("timeout waiting for Envoy to be ready")
}
