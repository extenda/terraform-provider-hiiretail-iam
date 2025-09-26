# Feature Specification: Multi-Tenant Support

**Feature Branch**: `002-tenant-support`  
**Created**: 2025-09-26  
**Status**: Draft  
**Input**: User description: "Add support for test and live tenants without requiring explicit URLs"

## Overview
Enhance the HiiRetail IAM provider to support simplified configuration for test and live tenants. Instead of requiring users to know and specify complete API URLs, users should simply specify whether they want to use a test or live tenant, making the provider more user-friendly and less error-prone.

## User Scenarios & Testing

### Primary User Story
As a HiiRetail IAM administrator, I want to easily switch between test and live environments in my Terraform configurations without needing to remember specific API URLs, so that I can manage resources across environments more efficiently and with less chance of error.

### Acceptance Scenarios
1. **Given** I want to manage IAM resources in the test environment  
   **When** I configure the provider with `environment = "test"`  
   **Then** all API calls are directed to `https://iam-api.retailsvc-test.com`

2. **Given** I want to manage IAM resources in the live environment  
   **When** I configure the provider with `environment = "live"`  
   **Then** all API calls are directed to `https://iam-api.retailsvc.com`

3. **Given** I have existing configurations using explicit URLs  
   **When** I continue using the `base_url` parameter  
   **Then** the provider respects my explicit URL configuration (backward compatibility)

### Edge Cases
- What happens if neither environment nor base_url is specified?
  - Should default to "live" environment with clear documentation
- What happens if both environment and base_url are specified?
  - base_url should take precedence with a warning log message
- What happens with invalid environment values?
  - Should fail fast with clear error message about valid options

## Requirements

### Functional Requirements
- **FR-001**: System MUST support a new `environment` provider configuration option
- **FR-002**: System MUST map "test" environment to `https://iam-api.retailsvc-test.com`
- **FR-003**: System MUST map "live" environment to `https://iam-api.retailsvc.com`
- **FR-004**: System MUST maintain backward compatibility with existing `base_url` parameter
- **FR-005**: System MUST provide clear validation errors for invalid environment values
- **FR-006**: System MUST default to "live" environment if neither environment nor base_url is specified
- **FR-007**: System MUST prefer explicit base_url over environment if both are specified

### Provider Schema
```hcl
provider "hiiretail-iam" {
  environment = optional(string) # Optional: "test" or "live". Defaults to "live"
  base_url   = optional(string) # Optional: Override environment with explicit URL
  token      = string          # Required: Your HiiRetail API token
}
```

### Migration Impact
- Existing configurations using base_url continue to work without changes
- New configurations can use simplified environment parameter
- Documentation needs updating to reflect new configuration options

## Non-Functional Requirements
- **Error Handling**:
  - Clear validation messages for invalid environment values
  - Warning logs when both environment and base_url are specified
  - Helpful upgrade guidance in error messages
- **Documentation**:
  - Clear examples for both configuration methods
  - Migration guide for existing users
  - Best practices for environment selection
- **Testing**:
  - Unit tests for URL mapping logic
  - Acceptance tests for both environment options
  - Backward compatibility tests

## Review & Acceptance Checklist

### Implementation Quality
- [ ] Clean separation of URL resolution logic
- [ ] Comprehensive validation messages
- [ ] Full test coverage
- [ ] Updated documentation
- [ ] Backward compatible changes

### Documentation Completeness
- [ ] Provider configuration reference updated
- [ ] Example configurations for both methods
- [ ] Clear explanation of precedence rules
- [ ] Migration guidance included

## Test Scenarios
1. Configuration validation
   - Test environment specification
   - Live environment specification
   - Invalid environment values
   - Missing environment (default case)
   - Both environment and base_url specified

2. URL Resolution
   - Test environment URL mapping
   - Live environment URL mapping
   - Explicit base_url override
   - Default environment behavior

3. Backward Compatibility
   - Existing base_url configurations
   - Mixed environment and base_url scenarios