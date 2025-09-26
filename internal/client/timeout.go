package client

import (
	"context"
	"time"
)

// TimeoutConfig holds timeout configuration for different operations
type TimeoutConfig struct {
	// Default timeout for operations if not specified
	Default time.Duration
	// Create operation timeout
	Create time.Duration
	// Read operation timeout
	Read time.Duration
	// Update operation timeout
	Update time.Duration
	// Delete operation timeout
	Delete time.Duration
	// List operation timeout
	List time.Duration
}

// DefaultTimeoutConfig provides default timeouts for operations
var DefaultTimeoutConfig = TimeoutConfig{
	Default: 30 * time.Second,
	Create:  60 * time.Second,
	Read:    30 * time.Second,
	Update:  60 * time.Second,
	Delete:  60 * time.Second,
	List:    60 * time.Second,
}

// withTimeout wraps an operation with a timeout based on operation type
func withTimeout(ctx context.Context, timeout time.Duration, operation func(context.Context) error) error {
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	done := make(chan error, 1)
	go func() {
		done <- operation(ctx)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}

// getTimeout returns the appropriate timeout for the given operation type
func (c *TimeoutConfig) getTimeout(opType string) time.Duration {
	switch opType {
	case "create":
		return c.Create
	case "read":
		return c.Read
	case "update":
		return c.Update
	case "delete":
		return c.Delete
	case "list":
		return c.List
	default:
		return c.Default
	}
}