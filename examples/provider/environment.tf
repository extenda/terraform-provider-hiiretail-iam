# Test environment configuration example
provider "hiiretail-iam" {
  environment = "test"  # Use test environment
}

# Example test environment resource
resource "hiiretail-iam_group" "test_group" {
  name = "test-group"
  # ... other group attributes
}