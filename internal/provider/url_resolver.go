package provider

import (
	"fmt"
)

// URLResolver handles the resolution of API URLs based on environment settings
type URLResolver struct {
	environment string
	baseURL    string
}

// NewURLResolver creates a new URLResolver instance
func NewURLResolver(environment, baseURL string) *URLResolver {
	return &URLResolver{
		environment: environment,
		baseURL:    baseURL,
	}
}

// ResolveURL returns the appropriate URL based on configuration
func (r *URLResolver) ResolveURL() string {
	// Base URL takes precedence if specified
	if r.baseURL != "" {
		return r.baseURL
	}

	// Default to live environment if not specified
	env := r.environment
	if env == "" {
		env = "live"
	}

	return r.environmentToURL(env)
}

// environmentToURL maps environment names to their corresponding URLs
func (r *URLResolver) environmentToURL(env string) string {
	urls := map[string]string{
		"test": "https://iam-api-test.retailsvc.com/schemas/v1/openapi.json",
		"live": "https://iam-api.retailsvc.com/schemas/v1/openapi.json",
	}

	if url, ok := urls[env]; ok {
		return url
	}

	// This shouldn't happen due to schema validation, but included for safety
	return urls["live"]
}

// ValidateEnvironment checks if the given environment is valid
func (r *URLResolver) ValidateEnvironment() error {
	if r.environment == "" {
		return nil // Empty is valid, will use default
	}

	validEnvs := []string{"test", "live"}
	for _, env := range validEnvs {
		if r.environment == env {
			return nil
		}
	}

	return fmt.Errorf("invalid environment %q, valid values are: test, live", r.environment)
}