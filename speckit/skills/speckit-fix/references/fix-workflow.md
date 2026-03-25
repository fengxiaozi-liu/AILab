# Fix Workflow

## Purpose

Use this file as the operating contract for `speckit-fix`.
The goal is to repair issue-scoped problems in one pass without turning the repair pass into a new feature.

## Required Inputs

- `specs/<feature>/issue.md`
- `specs/<feature>/spec.md`

Optional supporting inputs:

- `specs/<feature>/plan.md`
- `specs/<feature>/tasks.md`
- `specs/<feature>/review.md`

## Repair Rules

- The original `spec.md` remains the highest boundary constraint
- Every repair item must map back to an `ISSUE-xxx`
- If a proposed change cannot be justified by an issue item, do not include it
- If solving the issue requires expanding the feature scope, stop and request a new feature flow
- Treat human-written issue content as the minimal source of truth
- Fill missing metadata, scope, and trace fields during the AI repair flow
- Run clarify, plan, tasks, and implement as one continuous `/fix` execution

## Stage Expectations

### Clarify

- Resolve ambiguous issue statements
- Confirm expected behavior
- Confirm repair boundary
- Complete any missing issue fields that are needed for planning

### Plan

- Identify root cause area
- List impacted modules and files
- Define validation strategy
- State explicit out-of-scope items
- Keep the repair plan as internal reasoning unless the user explicitly asks to persist it

### Tasks

- Break the repair into small, file-targeted tasks
- Keep each task scoped to issue resolution
- Keep the repair tasks as internal reasoning unless the user explicitly asks to persist them

### Implement

- Apply the minimal safe fix
- Validate the affected path
- Update issue status after implementation

## Suggested Outputs

- Updated `specs/<feature>/issue.md`
- Code changes required to resolve the issue-scoped problems
