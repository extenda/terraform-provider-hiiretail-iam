# HiiRetail IAM Provider

The HiiRetail IAM provider allows you to manage IAM (Identity and Access Management) resources in your HiiRetail environment.

## Example Usage

```hcl
# Configure the HiiRetail IAM Provider with environment selection
provider "hiiretail-iam" {
  environment = "test"  # Optional: Use test environment (defaults to live)
}

# Configure the HiiRetail IAM Provider with custom URL
provider "hiiretail-iam" {
  base_url = "https://custom-api.example.com"  # Optional: Override API URL
}

# Minimal configuration using environment variables
provider "hiiretail-iam" {}
```

## Authentication

The provider expects an API token to be set via the `HIIRETAIL_TOKEN` environment variable:

```sh
export HIIRETAIL_TOKEN="your-api-token-here"
```

## Provider Configuration

### Optional

* `environment` - (Optional) The HiiRetail environment to use. Valid values are:
  * `test` - Use the test environment API endpoints
  * `live` - Use the live environment API endpoints (default)

* `base_url` - (Optional) Override the API endpoint URL. If specified, this takes precedence over the `environment` setting.

* `openapi_schema` - (Optional) OpenAPI schema URL for the HiiRetail IAM API. If not specified, the URL will be determined based on the selected environment.

## Environment Selection

The provider supports multiple environments through the `environment` parameter:

### Live Environment (Default)
```hcl
provider "hiiretail-iam" {
  # environment = "live" # This is the default
}
```

### Test Environment
```hcl
provider "hiiretail-iam" {
  environment = "test"
}
```

### Custom Environment
```hcl
provider "hiiretail-iam" {
  base_url = "https://custom-api.example.com"
}
```

## URL Resolution

The provider uses the following logic to determine the API endpoint URL:

1. If `base_url` is specified, it is used directly
2. Otherwise, if `environment` is specified, the corresponding environment URL is used
3. If neither is specified, the live environment URL is used

## Migration Guide

### Upgrading to Multi-Environment Support

If you're upgrading from a previous version that didn't support environments, your existing configurations will continue to work without changes. The provider maintains backward compatibility:

- Existing configurations without `environment` or `base_url` will default to the live environment
- Existing configurations with explicit `openapi_schema` URLs will continue to work as before

To take advantage of the new environment support:

1. Remove any hardcoded `openapi_schema` URLs from your provider configurations
2. Add the `environment` parameter to specify your target environment:

```hcl
# Before
provider "hiiretail-iam" {
  openapi_schema = "https://iam-api.retailsvc.com/schemas/v1/openapi.json"
}

# After
provider "hiiretail-iam" {
  environment = "live"  # or "test" for test environment
}
```

## Import

The provider does not require any special configuration for importing existing resources. Use the standard Terraform import commands with the appropriate environment configuration.