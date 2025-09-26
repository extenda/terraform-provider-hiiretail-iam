
# Implementation Plan: IAM Group Resource

**Branch**: `001-add-iam-group` | **Date**: 2025-09-25 | **Spec**: [/specs/001-add-iam-group/spec.md](/specs/001-add-iam-group/spec.md)
**Input**: Feature specification from `/specs/001-add-iam-group/spec.md`

## Execution Flow (/plan command scope)
```
1. Load feature spec from Input path
   → If not found: ERROR "No feature spec at {path}"
2. Fill Technical Context (scan for NEEDS CLARIFICATION)
   → Detect Project Type from context (web=frontend+backend, mobile=app+api)
   → Set Structure Decision based on project type
3. Fill the Constitution Check section based on the content of the constitution document.
4. Evaluate Constitution Check section below
   → If violations exist: Document in Complexity Tracking
   → If no justification possible: ERROR "Simplify approach first"
   → Update Progress Tracking: Initial Constitution Check
5. Execute Phase 0 → research.md
   → If NEEDS CLARIFICATION remain: ERROR "Resolve unknowns"
6. Execute Phase 1 → contracts, data-model.md, quickstart.md, agent-specific template file (e.g., `CLAUDE.md` for Claude Code, `.github/copilot-instructions.md` for GitHub Copilot, `GEMINI.md` for Gemini CLI, `QWEN.md` for Qwen Code or `AGENTS.md` for opencode).
7. Re-evaluate Constitution Check section
   → If new violations: Refactor design, return to Phase 1
   → Update Progress Tracking: Post-Design Constitution Check
8. Plan Phase 2 → Describe task generation approach (DO NOT create tasks.md)
9. STOP - Ready for /tasks command
```

**IMPORTANT**: The /plan command STOPS at step 7. Phases 2-4 are executed by other commands:
- Phase 2: /tasks command creates tasks.md
- Phase 3-4: Implementation execution (manual or via tools)

## Summary
Implement a Terraform resource for managing HiiRetail IAM groups using the OpenAPI specification from https://iam-api.retailsvc.com/schemas/v1/openapi.json. The resource will support CRUD operations for IAM groups with name, description, and ID management, following HashiCorp's provider development best practices.

## Technical Context
**Language/Version**: Go 1.21+ (Terraform Plugin Framework requirement)  
**Primary Dependencies**: 
- Terraform Plugin Framework v1.4.0+
- OpenAPI specification from [iam-api.retailsvc.com](https://iam-api.retailsvc.com/schemas/v1/openapi.json)
- Go-OpenAPI/Swag for schema validation  
**Storage**: N/A (API-driven)  
**Testing**: 
- Terraform Acceptance Testing Framework
- Go testing for unit tests
- testify for assertions  
**Target Platform**: Cross-platform (Terraform provider)  
**Project Type**: Terraform Provider (Plugin)  
**Performance Goals**: API operations complete within 10 seconds  
**Constraints**: 
- Follow OpenAPI specification constraints
- State must be safely managed
- Backward compatibility within major versions  
**Scale/Scope**: Single resource implementation

## Constitution Check
*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### I. Provider Foundations
- ✓ Uses Plugin Framework APIs for resource implementation
- ✓ Follows HashiCorp's CRUD pattern for resources
- ✓ Proper state management with id tracking
- ✓ Clear resource type (not a data source)

### II. Testing Fundamentals
- ✓ Acceptance tests planned for all operations
- ✓ Unit tests planned for name validation
- ✓ Test fixtures will be self-contained
- ✓ Full CRUD test coverage planned

### III. Schema Contracts
- ✓ Schema follows OpenAPI specification
- ✓ Required fields clearly marked
- ✓ Computed fields (id) properly handled
- ✓ All attributes documented

### IV. Error Handling & Diagnostics
- ✓ Error cases identified in edge cases
- ✓ Import behavior documented
- ✓ State safety considered

### V. Documentation & Examples
- ✓ Resource documentation planned
- ✓ Example configurations in spec
- ✓ Import documentation required

## Project Structure

### Documentation (this feature)
```
specs/[###-feature]/
├── plan.md              # This file (/plan command output)
├── research.md          # Phase 0 output (/plan command)
├── data-model.md        # Phase 1 output (/plan command)
├── quickstart.md        # Phase 1 output (/plan command)
├── contracts/           # Phase 1 output (/plan command)
└── tasks.md             # Phase 2 output (/tasks command - NOT created by /plan)
```

### Source Code (repository root)
```
internal/
├── provider/
│   ├── provider.go        # Provider definition
│   └── provider_test.go   # Provider tests
├── group/
│   ├── resource.go        # IAM group resource implementation
│   └── resource_test.go   # Resource tests
└── client/
    └── client.go          # OpenAPI client wrapper

examples/
└── resources/
    └── group/
        └── resource.tf     # Example configurations

docs/
└── resources/
    └── group.md           # Resource documentation
```

**Structure Decision**: Standard Terraform provider layout

## Phase 0: Outline & Research
1. **Review OpenAPI Specification**:
   - Analyze group-related endpoints
   - Document request/response schemas
   - Identify required headers and auth
   - Map OpenAPI types to Terraform types

2. **Research Provider Framework Implementation**:
   - Study schema validation patterns
   - Review state tracking best practices
   - Examine error handling approaches
   - Understand import functionality

3. **Consolidate findings** in `research.md`:
   - OpenAPI Integration:
     * Endpoint mappings for CRUD
     * Authentication requirements
     * Schema type mappings
   - Provider Framework:
     * Resource implementation pattern
     * State management approach
     * Error handling strategy
   - Testing Strategy:
     * Acceptance test setup
     * Mock API responses
     * Test case coverage

**Output**: research.md with implementation approach

## Phase 1: Design & Contracts
*Prerequisites: research.md complete*

1. **Design Resource Schema** → `data-model.md`:
   - Map OpenAPI schema to Terraform schema
   - Define validation rules
   - Document state transitions
   - Specify computed fields behavior

2. **Define Provider Contracts**:
   - Map CRUD operations to API endpoints
   - Define error handling contracts
   - Specify state management rules
   - Document import behavior

3. **Design Test Cases**:
   - Basic CRUD acceptance tests
   - Error handling test cases
   - Import functionality tests
   - State management tests

4. **Create Example Configurations**:
   - Basic usage examples
   - Complete attribute examples
   - Import examples
   - Error handling examples

5. **Update Documentation Plan**:
   - Resource documentation structure
   - Attribute descriptions
   - Example usage
   - Import guide
   - Troubleshooting section

**Output**: 
- data-model.md with schema design
- contracts/ with API mappings
- quickstart.md with examples

## Phase 2: Task Planning Approach
*This section describes what the /tasks command will do - DO NOT execute during /plan*

**Task Generation Strategy**:
- Load `.specify/templates/tasks-template.md` as base
- Generate tasks from Phase 1 design docs
- Acceptance test tasks [P] for CRUD operations
- Resource implementation tasks
- Documentation and example tasks [P]

**Ordering Strategy**:
1. Setup tasks:
   - Provider configuration
   - OpenAPI client setup
   - Test infrastructure

2. Test tasks [P]:
   - Basic CRUD test cases
   - Error handling tests
   - Import functionality tests
   - State management tests

3. Implementation tasks:
   - Resource schema definition
   - CRUD operations implementation
   - Validation logic
   - Error handling
   - Import support

4. Documentation tasks [P]:
   - Resource documentation
   - Example configurations
   - Import guide
   - Troubleshooting guide

**Estimated Output**: 15-20 numbered, ordered tasks in tasks.md

**IMPORTANT**: This phase is executed by the /tasks command, NOT by /plan

## Phase 3+: Future Implementation
*These phases are beyond the scope of the /plan command*

**Phase 3**: Task execution (/tasks command creates tasks.md)  
**Phase 4**: Implementation (execute tasks.md following constitutional principles)  
**Phase 5**: Validation (run tests, execute quickstart.md, performance validation)

## Complexity Tracking
*Fill ONLY if Constitution Check has violations that must be justified*

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |


## Progress Tracking
*This checklist is updated during execution flow*

**Phase Status**:
- [x] Phase 0: Research complete (/plan command)
- [x] Phase 1: Design complete (/plan command)
- [x] Phase 2: Task planning complete (/plan command - describe approach only)
- [ ] Phase 3: Tasks generated (/tasks command)
- [ ] Phase 4: Implementation complete
- [ ] Phase 5: Validation passed

**Gate Status**:
- [x] Initial Constitution Check: PASS
- [x] Post-Design Constitution Check: PASS
- [x] All NEEDS CLARIFICATION resolved
- [x] OpenAPI integration approach documented
- [x] Provider Framework compliance verified

---
*Based on Constitution v2.1.1 - See `/memory/constitution.md`*
