package provider

import (
	"context"
	"os"

	"github.com/extenda/terraform-provider-hiiretail-iam/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &HiiRetailProvider{}
)

// HiiRetailProvider is the provider implementation.
type HiiRetailProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// HiiRetailProviderModel describes the provider data model.
type HiiRetailProviderModel struct {
	OpenAPISchema types.String `tfsdk:"openapi_schema"`
	Environment   types.String `tfsdk:"environment"`
	BaseURL      types.String `tfsdk:"base_url"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &HiiRetailProvider{
			version: version,
		}
	}
}

// Metadata returns the provider type name.
func (p *HiiRetailProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "hiiretail"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *HiiRetailProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"openapi_schema": schema.StringAttribute{
				Optional:    true,
				Sensitive:  false,
				Description: "OpenAPI schema URL for the HiiRetail IAM API",
			},
			"environment": schema.StringAttribute{
				Optional:    true,
				Description: "The HiiRetail environment to use. Valid values are 'test' or 'live'. Defaults to 'live'.",
				Validators: []validator.String{
					stringvalidator.OneOf("test", "live"),
				},
			},
			"base_url": schema.StringAttribute{
				Optional:    true,
				Description: "Override the API endpoint URL. If specified, takes precedence over environment.",
			},
		},
	}
}

// Configure prepares a HiiRetail API client for data sources and resources.
func (p *HiiRetailProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config HiiRetailProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create a URL resolver with environment settings
	resolver := NewURLResolver(
		config.Environment.ValueString(),
		config.BaseURL.ValueString(),
	)

	// Use the resolver to determine the API URL
	apiURL := config.OpenAPISchema.ValueString()
	if apiURL == "" {
		apiURL = resolver.ResolveURL()
	}

	// Validate environment if specified
	if !config.Environment.IsNull() {
		if err := resolver.ValidateEnvironment(); err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("environment"),
				"Invalid Environment Configuration",
				err.Error(),
			)
			return
		}
	}

	// Get the API token from the environment
	token := os.Getenv("HIIRETAIL_TOKEN")
	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("openapi_schema"),
			"Missing API Token Configuration",
			"The HIIRETAIL_TOKEN environment variable must be set.",
		)
		return
	}

	// Initialize a new HiiRetail client using the configuration
	var c client.IClient = client.NewClient(apiURL, token)

	resp.DataSourceData = c
	resp.ResourceData = c
}

// DataSources defines the data sources implemented in the provider.
func (p *HiiRetailProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

// Resources defines the resources implemented in the provider.
func (p *HiiRetailProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewGroupResource,
	}
}