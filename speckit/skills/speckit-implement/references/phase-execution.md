# Phase Execution Guide

## Purpose

Use this file as the execution guardrail for `speckit-implement`.
It defines how to move through `tasks.md` without crossing phases or faking progress.

## Phase Cursor Rules

- Treat the first unfinished Phase as the active execution context
- Read later phases only for architectural foresight, not for early implementation
- Ignore completed tasks except as context

## Task Batching Rules

- Prefer batching tasks that touch the same module or highly related files
- Avoid mixing unrelated technical areas in one batch unless `tasks.md` explicitly groups them
- Keep each batch small enough to validate before moving on

## Status Writeback Rules

- Mark `[✅️]` only after the code has actually been written
- Update `tasks.md` regularly within the active phase
- Do not postpone all writeback until the end of the project
- If a task is blocked, report it explicitly instead of pretending partial completion

## Skill Loading Rules

- If the active task clearly falls under another installed skill, load that skill before implementation
- Reuse already loaded skills within the same phase when appropriate
- If a missing skill becomes obvious mid-implementation, pause, load it, and review the code already written

## Blocked State Guidance

When blocked, report:

1. Which task is blocked
2. Why it is blocked
3. What has already been completed
4. What user decision or missing input is needed
