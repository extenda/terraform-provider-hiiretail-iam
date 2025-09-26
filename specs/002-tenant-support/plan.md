# Implementation Plan: Multi-Tenant Support

**Branch**: `002-tenant-support` | **Date**: 2025-09-26 | **Spec**: [/specs/002-tenant-support/spec.md](/specs/002-tenant-support/spec.md)
**Input**: Feature specification and research from `/specs/002-tenant-support/`

## Summary
Enhance the HiiRetail IAM provider with environment-based configuration for test and live tenants, following established provider patterns while maintaining backward compatibility with existing configurations.

## Technical Context
**Language/Version**: Go 1.21+ (Terraform Plugin Framework requirement)  
**Primary Dependencies**: 
- Terraform Plugin Framework v1.4.0+
- Existing provider codebase  
**Storage**: N/A (API-driven)  
**Testing**: 
- Unit tests for URL resolution
- Acceptance tests for environment selection
- Backward compatibility tests  
**Target Platform**: Cross-platform (Terraform provider)  
**Project Type**: Provider Enhancement  
**Performance Goals**: No impact on existing performance  
**Constraints**: 
- Must maintain backward compatibility
- Must validate environment values
- Must log environment selection  
**Scale/Scope**: Provider configuration enhancement

## Constitution Check

### I. Provider Foundations
- ✓ Uses Plugin Framework schema for configuration
- ✓ Maintains backward compatibility
- ✓ Clear validation rules
- ✓ Proper logging of configuration

### II. Testing Fundamentals
- ✓ Unit tests for URL resolution
- ✓ Acceptance tests for environments
- ✓ Backward compatibility tests
- ✓ Error message validation

### III. Schema Updates
- ✓ New environment parameter
- ✓ Maintained base_url parameter
- ✓ Clear documentation
- ✓ Proper validations

### IV. Error Handling & Diagnostics
- ✓ Environment validation messages
- ✓ URL resolution logging
- ✓ Clear upgrade guidance
- ✓ Warning messages for conflicts

### V. Documentation & Examples
- ✓ Updated provider documentation
- ✓ Configuration examples
- ✓ Migration guide
- ✓ Best practices guide

## Project Structure

### Source Code Changes
```
internal/
├── provider/
│   ├── provider.go        # Add environment config
│   ├── provider_test.go   # Add environment tests
│   └── url_resolver.go    # New URL resolution logic
```

### Documentation Updates
```
docs/
└── index.md              # Update provider configuration docs
examples/
└── provider/
    ├── environment.tf    # Environment-based example
    └── migration.tf      # Migration example
```

## Phase 0: Research & Analysis
1. **Review similar providers** → `research.md`:
   - ✓ Azure environment pattern
   - ✓ AWS endpoint configuration
   - ✓ Testing approaches
   - ✓ Documentation patterns

2. **Identify impact areas**:
   - Provider schema
   - URL resolution
   - Configuration validation
   - Documentation updates

## Phase 1: Design & Contracts

1. **Provider Schema Update**:
   ```hcl
   provider "hiiretail-iam" {
     environment = optional(string) # "test" or "live"
     base_url   = optional(string) # Override URL
     token      = string          # API token
   }
   ```

2. **URL Resolution Logic**:
   ```go
   type URLResolver struct {
     environment string
     baseURL    string
   }

   func (r *URLResolver) ResolveURL() string {
     if r.baseURL != "" {
       return r.baseURL
     }
     return r.environmentToURL()
   }
   ```

3. **Validation Rules**:
   - Environment values: "test", "live"
   - Environment validation before URL resolution
   - Base URL format validation
   - Precedence documentation

4. **Error Messages**:
   - Invalid environment: List valid options
   - URL resolution: Clear error context
   - Configuration conflicts: Warning with precedence
   - Migration guidance: Clear upgrade path

## Phase 2: Task Planning Approach

1. **Schema Updates**:
   - Add environment parameter
   - Update schema documentation
   - Add validation functions
   - Update provider tests

2. **URL Resolution**:
   - Create URL resolver
   - Add environment mapping
   - Add logging
   - Add resolver tests

3. **Documentation**:
   - Update provider docs
   - Add migration guide
   - Create examples
   - Document best practices

4. **Testing**:
   - Add unit tests
   - Add acceptance tests
   - Test backward compatibility
   - Test error messages

## Implementation Steps

### Task Categories:

1. **Provider Enhancement**:
   - Schema updates
   - Configuration handling
   - URL resolution
   - Validation logic

2. **Testing**:
   - Unit tests
   - Acceptance tests
   - Migration tests
   - Error handling tests

3. **Documentation**:
   - Provider configuration
   - Environment options
   - Migration guide
   - Examples

4. **Quality Assurance**:
   - Code review
   - Test coverage
   - Documentation review
   - Migration testing

## Risk Assessment

### Potential Issues:
1. **Breaking Changes**:
   - Mitigation: Maintain base_url support
   - Test existing configurations
   - Clear migration guide

2. **URL Resolution**:
   - Mitigation: Comprehensive testing
   - Proper error handling
   - Clear logging

3. **Configuration Conflicts**:
   - Mitigation: Clear precedence rules
   - Warning messages
   - Documentation

### Validation Points:
1. **Functionality**:
   - Environment selection works
   - URL resolution correct
   - Backward compatibility maintained

2. **Error Handling**:
   - Clear error messages
   - Proper validation
   - Helpful guidance

3. **Documentation**:
   - Clear configuration guide
   - Good examples
   - Migration guidance

## Success Criteria

1. **Provider Configuration**:
   - Environment selection works
   - Base URL still supported
   - Clear validation messages

2. **Testing**:
   - All tests pass
   - Good coverage
   - No regressions

3. **Documentation**:
   - Updated and clear
   - Good examples
   - Migration guide

4. **Compatibility**:
   - No breaking changes
   - Clear upgrade path
   - Good error messages