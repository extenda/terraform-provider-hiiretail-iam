# Terraform Provider Constitution

## Core Principles

### I. Provider Foundations
Every feature must align with the core Terraform Provider Development Framework (Plugin Framework) principles:
- Use Plugin Framework APIs and paradigms consistently
- Follow HashiCorp's provider development patterns
- Support proper state management and Terraform lifecycle
- Implement data source/resource distinction correctly

### II. Testing Fundamentals
Test-driven development is mandatory for all provider features:
- Acceptance tests required for every resource and data source
- Unit tests for complex validation or computation logic
- Test fixtures must be self-contained and cleaned up
- Test cases must cover create, read, update, delete flows
- Mock external API calls in unit tests

### III. Schema Contracts
Provider schema changes follow strict versioning rules:
- Breaking changes require major version bump
- Schema must be backwards compatible within major version
- Required fields clearly documented with constraints
- Computed fields marked appropriately
- Default values carefully considered and documented
- All attributes must have clear descriptions

### IV. Error Handling & Diagnostics
Implement robust error handling and user feedback:
- Clear diagnostic messages with actionable context
- Expected error cases explicitly handled and tested
- Unexpected errors properly wrapped with context
- Import behavior documented and tested
- Resource state handled safely during failures
- Follow HashiCorp's diagnostics best practices

### V. Documentation & Examples
Documentation is a primary project deliverable:
- Resource/data source docs with all attributes defined
- Example configurations demonstrating common use cases
- Import documentation with step-by-step instructions
- Release notes detailing all changes
- Troubleshooting guide for common issues
- Provider website docs kept in sync with codebase

## Security & Compliance

### Authentication & Authorization
- Support secure authentication methods only
- Sensitive fields marked in schema
- Credentials handled according to HashiCorp guidelines
- Audit logging for significant operations
- Rate limiting and retry handling implemented

### API Integration
- Implement proper API error handling
- Follow API best practices and conventions
- Version API clients appropriately
- Handle API deprecations gracefully
- Document API version compatibility

## Development Workflow

### Feature Development Process
1. Write acceptance test first (red)
2. Implement minimal provider code to pass test (green)
3. Refactor while maintaining test coverage
4. Document the feature comprehensively
5. Submit for review with test evidence

### Code Review Requirements
- Acceptance tests must pass
- Documentation is complete and accurate 
- Breaking changes clearly identified
- Error handling meets standards
- State management verified
- Security implications considered

### Release Process
- Semantic versioning strictly followed
- Release notes complete and accurate
- Documentation updated
- Migration guide for breaking changes
- Security advisories if applicable

## Governance

Constitution changes require:
1. Documentation of rationale
2. Review of impact on existing resources
3. Migration plan for breaking changes
4. Update of affected documentation
5. Version bump according to semantic versioning

**Version**: 1.0.0 | **Ratified**: 2025-09-25 | **Last Amended**: 2025-09-25

<!--
Sync Impact Report:
Initial constitution creation establishing core principles for provider development.

Version: 1.0.0 (Initial creation)

Core Principles Added:
- Provider Foundations
- Testing Fundamentals
- Schema Contracts
- Error Handling & Diagnostics
- Documentation & Examples

Sections Added:
- Security & Compliance
- Development Workflow
- Governance

Templates to Update:
✅ .specify/templates/plan-template.md
✅ .specify/templates/spec-template.md
✅ .specify/templates/tasks-template.md
-->