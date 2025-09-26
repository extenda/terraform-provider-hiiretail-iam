# Tasks: IAM Group Resource

**Input**: Design documents from `/specs/001-add-iam-group/`
**Prerequisites**: plan.md (required), research.md, data-model.md, contracts/

## Phase 3.1: Setup
- [ ] T001 [P] Create provider directory structure
      ```
      internal/
      ├── provider/
      ├── group/
      └── client/
      ```
- [ ] T002 Initialize Go module and add Terraform Plugin Framework dependencies
- [ ] T003 [P] Configure testing infrastructure with testify and test fixtures
- [ ] T004 [P] Add OpenAPI client wrapper in internal/client/client.go

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
- [ ] T005 [P] Acceptance test: Create IAM group in internal/group/resource_test.go
- [ ] T006 [P] Acceptance test: Read IAM group in internal/group/resource_test.go
- [ ] T007 [P] Acceptance test: Update IAM group description in internal/group/resource_test.go
- [ ] T008 [P] Acceptance test: Delete IAM group in internal/group/resource_test.go
- [ ] T009 [P] Acceptance test: Import existing group in internal/group/resource_test.go
- [ ] T010 [P] Unit test: Group name validation in internal/group/resource_test.go
- [ ] T011 [P] Unit test: Error handling in internal/group/resource_test.go

## Phase 3.3: Core Implementation
- [ ] T012 Define resource schema in internal/group/resource.go
- [ ] T013 Implement Create function with OpenAPI integration
- [ ] T014 Implement Read function with state updates
- [ ] T015 Implement Update function for description changes
- [ ] T016 Implement Delete function with error handling
- [ ] T017 Add Import support
- [ ] T018 Implement name validation logic
- [ ] T019 Add diagnostic message handling
- [ ] T020 Register resource in provider.go

## Phase 3.4: Documentation & Examples
- [ ] T021 [P] Create resource documentation in docs/resources/group.md
- [ ] T022 [P] Add example configurations in examples/resources/group/resource.tf
- [ ] T023 [P] Write import guide with examples
- [ ] T024 [P] Create troubleshooting guide

## Phase 3.5: Polish
- [ ] T025 [P] Add retry logic for API calls
- [ ] T026 Performance testing (<10s operations)
- [ ] T027 [P] Implement proper logging
- [ ] T028 [P] Add provider acceptance tests
- [ ] T029 Remove duplication and optimize code

## Dependencies
- Tests (T005-T011) must be written and failing before implementation (T012-T020)
- T012 must be completed before T013-T017
- T020 requires T012-T019
- Documentation can be done in parallel with implementation

## Parallel Execution Example
```bash
# Launch test creation tasks in parallel:
Task: "Write create group acceptance test"
Task: "Write read group acceptance test"
Task: "Write update group acceptance test"
Task: "Write delete group acceptance test"
Task: "Write import acceptance test"
```

## Notes
- All [P] tasks can be executed in parallel
- Follow TDD: tests must fail before implementation
- Verify OpenAPI spec conformance in tests
- Document state handling in comments
- Follow HashiCorp provider best practices

## Task Generation Rules
1. From Spec:
   - Each CRUD operation → acceptance test + implementation
   - Each validation rule → unit test + implementation
   - Each error case → test + handler

2. From Constitution:
   - Resource documentation required
   - Example configurations needed
   - Error handling comprehensive
   - State management explicit

3. From OpenAPI:
   - Map endpoints to resource functions
   - Validate request/response schemas
   - Handle API-specific errors
   - Follow API conventions

4. Task Dependencies:
   - Tests before implementation
   - Schema before CRUD
   - Error handling with each operation
   - Documentation after implementation verified

## Validation Checklist
- [x] All operations have corresponding tests
- [x] All errors have handling and tests
- [x] Documentation tasks complete coverage
- [x] Example configurations included
- [x] Constitution requirements addressed