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

4. "Lets commit this and dive into design"

5. "lets push the commit to remote" — push failed (no GitHub credentials in this
   environment); asked how to authenticate.

6. (clarifying-question round, answered) Git auth → "I'll run it myself" (user pushed the
   branch from their own terminal).

7. "pushed, let's start the design walkthrough. I have already included my design and
   thoughts in design.md" — reviewed the filled-in design.md and the flowchart as a
   senior-engineer pass: flagged a DDL syntax bug (missing comma before a table constraint
   in `ledger_entries`), an unexplained `currency` column inconsistent with the stated
   out-of-scope note, the idempotency check-then-insert race and whether it's handled,
   request-hash canonicalization, the `201`-for-business-failures status code choice, and a
   few smaller schema/testing notes. No code or schema was written on the user's behalf —
   findings only.

8. "1,2 syntax and Bug is fixed once check. 4. Request-hash is on (fromWalletId, toWalletId,
   amount) 5. Pending ones to be retried by client as a scope of this assignment. I am not
   implementing retry mechanism 6. thats a good one added UNIQUE to
   idempotency_records.transfer_id 9. Please add it as a note 8. Yes, it is intentional to
   keep FAILED as a terminal status 7. Agreed ledger-correctness tests has to assert this
   3. If a request fails unique-violation on idempotency_records this means, idempotency_key
   exists, next we have to validate the request and check transfer status and if status is
   pending we should attempt the wallet transfer." — resolved each review point; added the
   201-for-business-failures clarification as a note in design.md at the user's request.

9. "Yes add it" — added a note under design.md section 6 (Flow) spelling out the
   idempotency check-then-insert race and the unique-violation-as-signal handling, per the
   user's confirmed answer to review point #3.

<!-- Append new prompts below, in order, as the session continues. -->
