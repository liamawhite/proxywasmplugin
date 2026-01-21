package testutil

import (
	"sync"

	"github.com/testcontainers/testcontainers-go"
)

// EnvoyLogConsumer captures and stores Envoy container logs for test validation.
// It implements the testcontainers.LogConsumer interface and provides thread-safe
// access to accumulated logs.
type EnvoyLogConsumer struct {
	mu   sync.Mutex
	logs []string
}

// NewEnvoyLogConsumer creates a new log consumer for capturing Envoy output.
func NewEnvoyLogConsumer() *EnvoyLogConsumer {
	return &EnvoyLogConsumer{
		logs: make([]string, 0),
	}
}

// Accept receives log entries from the container and stores them.
// This method is called by testcontainers when new log output is available.
func (c *EnvoyLogConsumer) Accept(log testcontainers.Log) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.logs = append(c.logs, string(log.Content))
}

// GetLogs returns a copy of all captured log entries.
// This method is thread-safe and can be called while logs are still being captured.
func (c *EnvoyLogConsumer) GetLogs() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	// Return a copy to avoid race conditions
	logsCopy := make([]string, len(c.logs))
	copy(logsCopy, c.logs)
	return logsCopy
}

// Clear removes all captured logs.
// Useful for resetting the consumer between test scenarios.
func (c *EnvoyLogConsumer) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.logs = make([]string, 0)
}
