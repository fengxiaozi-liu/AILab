# Code Review Rubric

## Usage

Use this file as the generic review baseline for `speckit-code-review`.
Keep language- or framework-specific checks in separate skills or references.

## Review Dimensions

### 1. Completeness

- Claimed files or modules actually exist
- Core paths are implemented, not placeholders or TODOs
- Changes form a closed loop across code, tests, config, and docs where needed
- If `spec.md` exists, key requirements map to real implementation

### 2. Correctness

- Main flow, branches, boundaries, and failure paths match intent
- Errors are returned, transformed, or surfaced correctly
- Resource lifecycle is handled safely
- State transitions and concurrency behavior are coherent for the project context

### 3. Architecture and Maintainability

- Responsibilities sit at the right layer
- Dependencies flow in a clear direction
- Naming and intent are understandable
- Abstractions are justified rather than speculative
- Complexity is not pushed outward to callers without reason

### 4. Performance

- No obvious repeated expensive work
- No unbounded reads or scans without reason
- No avoidable blocking, serialization, or heavy allocations
- External calls and data access patterns are proportionate to the workload

### 5. Security

- Inputs are validated where needed
- Authorization, tenancy, ownership, or role checks are present where required
- Secrets or sensitive data are not exposed in code or logs
- No obvious injection paths through string concatenation or unsafe evaluation

### 6. Testing

- Core success path is covered
- Important failure and boundary cases are covered
- Tests protect against likely regressions from this change

## Severity Guide

- `BLOCKER`: core requirement missing, severe correctness break, or serious security issue
- `MAJOR`: likely bug, architecture risk, meaningful performance issue, or important test gap
- `MINOR`: maintainability, naming, structure, or low-risk design issue
- `INFO`: optional improvement or observation

## Report Shape

Prefer a report with:

1. Review mode
2. Findings ordered by severity
3. Open assumptions or unknowns
4. Coverage summary
5. Next-step recommendation
