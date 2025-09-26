package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	p := New("test")()
	assert.NotNil(t, p)
}

func TestProviderSchema(t *testing.T) {
	p := &HiiRetailProvider{}
	response := &provider.SchemaResponse{}

	p.Schema(context.Background(), provider.SchemaRequest{}, response)

	assert.NotNil(t, response.Schema)
	assert.NotNil(t, response.Schema.Attributes)

	// Check openapi_schema attribute
	assert.Contains(t, response.Schema.Attributes, "openapi_schema")
	openAPISchema := response.Schema.Attributes["openapi_schema"].(schema.StringAttribute)
	assert.False(t, openAPISchema.Required)
	assert.False(t, openAPISchema.Sensitive)

	// Check environment attribute
	assert.Contains(t, response.Schema.Attributes, "environment")
	environmentAttr := response.Schema.Attributes["environment"].(schema.StringAttribute)
	assert.False(t, environmentAttr.Required)
	assert.NotEmpty(t, environmentAttr.Description)
	assert.NotEmpty(t, environmentAttr.Validators)

	// Check base_url attribute
	assert.Contains(t, response.Schema.Attributes, "base_url")
	baseURLAttr := response.Schema.Attributes["base_url"].(schema.StringAttribute)
	assert.False(t, baseURLAttr.Required)
	assert.NotEmpty(t, baseURLAttr.Description)
}

func TestProviderConfigure_MissingToken(t *testing.T) {
	p := &HiiRetailProvider{}
	ctx := context.Background()

	// Get provider schema
	schemaResp := &provider.SchemaResponse{}
	p.Schema(ctx, provider.SchemaRequest{}, schemaResp)

	// Create empty config
	schema := tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"openapi_schema": tftypes.String,
			"environment":   tftypes.String,
			"base_url":      tftypes.String,
		},
	}
	config := tfsdk.Config{
		Schema: schemaResp.Schema,
		Raw: tftypes.NewValue(schema, map[string]tftypes.Value{
			"openapi_schema": tftypes.NewValue(tftypes.String, nil),
			"environment":   tftypes.NewValue(tftypes.String, nil),
			"base_url":      tftypes.NewValue(tftypes.String, nil),
		}),
	}

	// Create request with empty config
	req := provider.ConfigureRequest{
		Config: config,
	}
	resp := &provider.ConfigureResponse{
		Diagnostics: diag.Diagnostics{},
	}

	p.Configure(ctx, req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "HIIRETAIL_TOKEN environment variable must be set")
}

func TestProviderConfigure_ValidConfig(t *testing.T) {
	tests := []struct {
		name        string
		config      map[string]tftypes.Value
		envToken    string
		wantErr     bool
		errContains string
	}{
		{
			name: "Valid config with default settings",
			config: map[string]tftypes.Value{
				"openapi_schema": tftypes.NewValue(tftypes.String, "https://iam-api.retailsvc.com/schemas/v1/openapi.json"),
				"environment":    tftypes.NewValue(tftypes.String, nil),
				"base_url":      tftypes.NewValue(tftypes.String, nil),
			},
			envToken: "test-token",
			wantErr:  false,
		},
		{
			name: "Valid config with test environment",
			config: map[string]tftypes.Value{
				"openapi_schema": tftypes.NewValue(tftypes.String, nil),
				"environment":    tftypes.NewValue(tftypes.String, "test"),
				"base_url":      tftypes.NewValue(tftypes.String, nil),
			},
			envToken: "test-token",
			wantErr:  false,
		},
		{
			name: "Valid config with custom base_url",
			config: map[string]tftypes.Value{
				"openapi_schema": tftypes.NewValue(tftypes.String, nil),
				"environment":    tftypes.NewValue(tftypes.String, nil),
				"base_url":      tftypes.NewValue(tftypes.String, "https://custom-api.example.com"),
			},
			envToken: "test-token",
			wantErr:  false,
		},
		{
			name: "Invalid environment value",
			config: map[string]tftypes.Value{
				"openapi_schema": tftypes.NewValue(tftypes.String, nil),
				"environment":    tftypes.NewValue(tftypes.String, "invalid"),
				"base_url":      tftypes.NewValue(tftypes.String, nil),
			},
			envToken:    "test-token",
			wantErr:     true,
			errContains: "invalid environment",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HiiRetailProvider{}
			ctx := context.Background()

			// Get provider schema
			schemaResp := &provider.SchemaResponse{}
			p.Schema(ctx, provider.SchemaRequest{}, schemaResp)

			// Create config
			schema := tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"openapi_schema": tftypes.String,
					"environment":   tftypes.String,
					"base_url":      tftypes.String,
				},
			}
			config := tfsdk.Config{
				Schema: schemaResp.Schema,
				Raw:    tftypes.NewValue(schema, tt.config),
			}

			// Create request with config
			req := provider.ConfigureRequest{
				Config: config,
			}
			resp := &provider.ConfigureResponse{
				Diagnostics: diag.Diagnostics{},
			}

			// Set HIIRETAIL_TOKEN environment variable for testing
			t.Setenv("HIIRETAIL_TOKEN", tt.envToken)

			provider := &HiiRetailProvider{}
			provider.Configure(context.Background(), req, resp)

			if tt.wantErr {
				assert.True(t, resp.Diagnostics.HasError())
				if tt.errContains != "" {
					assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), tt.errContains)
				}
			} else {
				assert.False(t, resp.Diagnostics.HasError())
				assert.NotNil(t, resp.ResourceData)
				assert.NotNil(t, resp.DataSourceData)
			}
		})
	}
}