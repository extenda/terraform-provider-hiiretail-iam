package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestURLResolver(t *testing.T) {
	tests := []struct {
		name        string
		environment string
		baseURL     string
		want        string
	}{
		{
			name:        "Default to live environment",
			environment: "",
			baseURL:     "",
			want:        "https://iam-api.retailsvc.com/schemas/v1/openapi.json",
		},
		{
			name:        "Test environment",
			environment: "test",
			baseURL:     "",
			want:        "https://iam-api-test.retailsvc.com/schemas/v1/openapi.json",
		},
		{
			name:        "Live environment",
			environment: "live",
			baseURL:     "",
			want:        "https://iam-api.retailsvc.com/schemas/v1/openapi.json",
		},
		{
			name:        "Base URL takes precedence",
			environment: "test",
			baseURL:     "https://custom-api.example.com",
			want:        "https://custom-api.example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolver := NewURLResolver(tt.environment, tt.baseURL)
			got := resolver.ResolveURL()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestURLResolver_ValidateEnvironment(t *testing.T) {
	tests := []struct {
		name        string
		environment string
		wantErr     bool
	}{
		{
			name:        "Empty environment is valid",
			environment: "",
			wantErr:     false,
		},
		{
			name:        "Test environment is valid",
			environment: "test",
			wantErr:     false,
		},
		{
			name:        "Live environment is valid",
			environment: "live",
			wantErr:     false,
		},
		{
			name:        "Invalid environment",
			environment: "staging",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolver := NewURLResolver(tt.environment, "")
			err := resolver.ValidateEnvironment()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}