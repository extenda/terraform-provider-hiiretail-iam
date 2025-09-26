# Migration example from explicit URL to environment
provider "hiiretail-iam" {
  # Before:
  # openapi_schema = "https://iam-api.retailsvc.com/schemas/v1/openapi.json"

  # After:
  environment = "live"  # Use live environment (or "test" for test environment)
}

# Example using custom URL (if needed)
provider "hiiretail-iam" {
  alias    = "custom"
  base_url = "https://custom-api.example.com"
}

# Using the default provider (live environment)
resource "hiiretail-iam_group" "production_group" {
  name = "production-group"
  # ... other group attributes
}

# Using the custom provider
resource "hiiretail-iam_group" "custom_group" {
  provider = hiiretail-iam.custom
  name     = "custom-group"
  # ... other group attributes
}