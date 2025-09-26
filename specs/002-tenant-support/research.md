# Research: Multi-Tenant Support

**Feature Branch**: `002-tenant-support`  
**Date**: 2025-09-26  
**Input**: spec.md environment configuration requirements

## Provider Configuration Patterns

### Terraform Provider Best Practices Review

1. **Azure Provider**
   - Uses `environment` parameter for cloud selection (public, government, china)
   - Environment names map to predefined endpoints
   - Allows override with explicit endpoints
   - Example:
     ```hcl
     provider "azurerm" {
       environment = "public"  # or "china", "german", "usgovernment"
     }
     ```

2. **AWS Provider**
   - Uses named endpoints for different regions
   - Supports endpoint configuration for testing
   - Example:
     ```hcl
     provider "aws" {
       endpoints {
         s3 = "http://localhost:4566"  # for testing
       }
     }
     ```

3. **Google Provider**
   - Uses `user_project_override` for test vs production
   - Supports custom endpoints for testing
   - Example:
     ```hcl
     provider "google" {
       user_project_override = true
       endpoint = "https://test-endpoint.googleapis.com"  # optional
     }
     ```

### Common Patterns Found

1. **Environment Selection**
   - Most providers use string enum for environment selection
   - Common names: "production", "staging", "test"
   - Usually defaults to production/live environment
   - Clear validation messages for invalid values

2. **URL Resolution**
   - Environment names map to well-known URLs
   - Providers handle URL construction internally
   - Custom URLs allowed for testing/development
   - Base URLs often have version paths included

3. **Backward Compatibility**
   - Maintain support for existing parameters
   - Explicit URLs take precedence over environment selection
   - Clear deprecation warnings when needed
   - Migration guides provided

## Validation Approaches

1. **Environment Values**
   - Use string validation with allowed values
   - Clear error messages listing valid options
   - Case-insensitive comparison recommended
   - Example from AWS provider:
     ```go
     if err := validateEnvironment(env); err != nil {
       return fmt.Errorf("invalid environment %q, valid values are: %s", env, strings.Join(validEnvs, ", "))
     }
     ```

2. **URL Construction**
   - Central URL mapping function
   - Version path handling
   - Validation of constructed URLs
   - Example pattern:
     ```go
     func resolveBaseURL(env string) string {
       urls := map[string]string{
         "test": "https://api-test.example.com/v1",
         "live": "https://api.example.com/v1",
       }
       return urls[env]
     }
     ```

## Testing Strategies

1. **Unit Tests**
   - Test environment validation
   - Test URL resolution
   - Test precedence rules
   - Test default behaviors

2. **Acceptance Tests**
   - Test actual API interactions
   - Verify correct endpoint selection
   - Test backward compatibility
   - Example from Azure provider:
     ```go
     func TestAccProvider_Environment(t *testing.T) {
       // Test cases for environment selection
     }
     ```

## Documentation Patterns

1. **Provider Configuration**
   - Clear examples for each use case
   - Tables of valid environment values
   - Explicit precedence rules
   - Migration guides

2. **Argument Reference**
   - Document default values
   - Note deprecated parameters
   - Show environment-to-URL mapping
   - Example:
     ```markdown
     * `environment` - (Optional) The HiiRetail environment to use. Valid values are "test" or "live". Defaults to "live".
     * `base_url` - (Optional) Override the API endpoint URL. If specified, takes precedence over `environment`.
     ```

## Implementation Recommendations

1. **URL Resolution**
   - Create a dedicated URL resolver
   - Handle both environment and explicit URLs
   - Include version path management
   - Log selected environment/URL

2. **Schema Updates**
   - Add environment parameter
   - Maintain base_url parameter
   - Add proper validations
   - Document precedence

3. **Testing Focus**
   - URL resolution logic
   - Environment validation
   - Backward compatibility
   - Error messages

4. **Migration Support**
   - No breaking changes
   - Clear upgrade path
   - Comprehensive examples

## Implementation Strategy

Based on our research and planning, we recommend the following implementation strategy:

1. **Schema Enhancement**
   ```go
   type ProviderConfig struct {
     Environment *string `tfsdk:"environment"`
     BaseURL    *string `tfsdk:"base_url"`
     Token      string  `tfsdk:"token"`
   }
   ```

2. **URL Resolution Pattern**
   ```go
   type URLResolver struct {
     environment string
     baseURL    string
   }

   func (r *URLResolver) ResolveURL() string {
     if r.baseURL != "" {
       return r.baseURL // Explicit URL takes precedence
     }
     return r.environmentToURL()
   }
   ```

3. **Implementation Phases**
   - Schema update with backward compatibility
   - URL resolution with environment mapping
   - Validation and error handling
   - Documentation and testing

4. **Migration Support**
   - Keep base_url support indefinitely
   - Clear documentation for both approaches
   - Environment-specific examples
   - Testing guidelines

## Open Questions

1. Should we support more environments in the future (e.g., "staging")?
2. Should we deprecate base_url in a future version?
3. Should we add environment-specific token validation?

## Conclusions

1. **Recommended Approach**
   - Follow Azure provider pattern for environment selection
   - Maintain backward compatibility with base_url
   - Implement clear validation and error messages
   - Provide comprehensive documentation

2. **Key Considerations**
   - Keep URL construction internal
   - Validate environments early
   - Log selected environment
   - Support testing scenarios

3. **Next Steps**
   - Document URL mapping
   - Define validation rules
   - Plan backward compatibility
   - Create migration guide