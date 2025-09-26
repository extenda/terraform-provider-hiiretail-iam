package client

import (
	"context"
	"errors"
	"net/http"
	"time"
)

type RetryConfig struct {
	MaxRetries      int
	InitialInterval time.Duration
	MaxInterval     time.Duration
	Multiplier      float64
	MaxElapsedTime  time.Duration
}

var DefaultRetryConfig = RetryConfig{
	MaxRetries:      3,
	InitialInterval: 100 * time.Millisecond,
	MaxInterval:     10 * time.Second,
	Multiplier:      2.0,
	MaxElapsedTime:  30 * time.Second,
}

// isRetryable determines if an error should be retried
func isRetryable(resp *http.Response, err error) bool {
	if err != nil {
		// Retry on network errors
		return true
	}

	if resp == nil {
		return false
	}

	// Retry on rate limit errors (429), gateway errors (502, 503, 504)
	switch resp.StatusCode {
	case http.StatusTooManyRequests,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return true
	}

	return false
}

// withRetry executes the given operation with exponential backoff retry logic
func withRetry(ctx context.Context, config RetryConfig, operation func() (*http.Response, error)) (*http.Response, error) {
	var resp *http.Response
	var err error
	var nextInterval = config.InitialInterval
	startTime := time.Now()

	for retries := 0; retries <= config.MaxRetries; retries++ {
		// Check if we've exceeded the maximum elapsed time
		if time.Since(startTime) > config.MaxElapsedTime {
			if err != nil {
				return nil, errors.New("max elapsed time exceeded: " + err.Error())
			}
			return nil, errors.New("max elapsed time exceeded")
		}

		// If this isn't our first try, wait before retrying
		if retries > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(nextInterval):
				// Increase the interval for the next iteration
				nextInterval = time.Duration(float64(nextInterval) * config.Multiplier)
				if nextInterval > config.MaxInterval {
					nextInterval = config.MaxInterval
				}
			}
		}

		resp, err = operation()
		
		// If the operation was successful or we shouldn't retry, return the result
		if err == nil && (resp == nil || !isRetryable(resp, err)) {
			return resp, nil
		}

		// If this was our last try, return the error
		if retries == config.MaxRetries {
			if err != nil {
				return nil, errors.New("max retries exceeded: " + err.Error())
			}
			return resp, errors.New("max retries exceeded")
		}
	}

	return resp, err
}