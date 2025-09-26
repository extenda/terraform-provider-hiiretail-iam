package provider

import (
	"context"
	"os"
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
	assert.Contains(t, response.Schema.Attributes, "openapi_schema")

	openAPISchema := response.Schema.Attributes["openapi_schema"].(schema.StringAttribute)
	assert.False(t, openAPISchema.Required)
	assert.False(t, openAPISchema.Sensitive)
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
		},
	}
	config := tfsdk.Config{
		Schema: schemaResp.Schema,
		Raw:    tftypes.NewValue(schema, map[string]tftypes.Value{
			"openapi_schema": tftypes.NewValue(tftypes.String, nil),
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
	p := &HiiRetailProvider{}
	ctx := context.Background()

	// Get provider schema
	schemaResp := &provider.SchemaResponse{}
	p.Schema(ctx, provider.SchemaRequest{}, schemaResp)

	// Create config with valid values
	schema := tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"openapi_schema": tftypes.String,
		},
	}
	config := tfsdk.Config{
		Schema: schemaResp.Schema,
		Raw: tftypes.NewValue(schema, map[string]tftypes.Value{
			"openapi_schema": tftypes.NewValue(tftypes.String, "https://iam-api.retailsvc.com/schemas/v1/openapi.json"),
		}),
	}

	// Create request with valid config
	req := provider.ConfigureRequest{
		Config: config,
	}
	resp := &provider.ConfigureResponse{
		Diagnostics: diag.Diagnostics{},
	}

	// Set HIIRETAIL_TOKEN environment variable for testing
	_ = os.Setenv("HIIRETAIL_TOKEN", "test-token")
	
	p.Configure(ctx, req, resp)

	assert.False(t, resp.Diagnostics.HasError())
}