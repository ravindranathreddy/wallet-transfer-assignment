# AI Usage Log

## Tool

Claude Code (Anthropic CLI agent).

## How it was used

- Scoped strictly to boilerplate/project setup: Go module layout, Docker/docker-compose,
  environment configuration, HTTP routing skeleton, and dependency wiring in `cmd/server/main.go`.
- Used as a senior-engineer reviewer for my own designs and implementations — bug hunting,
  race conditions, edge cases, and improvement suggestions.
- Testing strategy was discussed and agreed on jointly before writing tests.
- It did **not** write the domain model, transfer state machine, ledger/idempotency logic,
  or locking strategy, and did not make design decisions on those — those are mine.

## Session transcript (prompts, chronological)

### 2026-07-05

1. "Read the Assignment.md file. Tech stack we will be going with is golang,PostgreSQL. In
   this implementation your role is an Assistant for boilerplate and project set up i.e
   Docker, configuration, routing, dependency wiring . A Senior Engineer while reviewing
   reviewing my designs and implementations, identifying bugs, race conditions, edge cases,
   and suggesting improvements. We will finally discuss on the testing. I prefer you not
   generating code for core logic and making design decisions on my behalf. Also we need to
   send prompts during this, please keep track of this to send."

2. (clarifying-question round, answered) Prompt/session tracking → committed log file in
   repo; HTTP routing → standard library only; DB access → pgx (native driver); migrations →
   golang-migrate.

3. "Why go 1.26.3 is a valid one right?" — correcting an unnecessary edit to go.mod's `go`
   directive (trimming `1.26.3` to `1.26`); confirmed the full patch version is valid and
   intentional since Go 1.21, so the change was reverted.

<!-- Append new prompts below, in order, as the session continues. -->
