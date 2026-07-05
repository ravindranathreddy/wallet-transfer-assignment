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

10. "Let's start implementing the migrations for the schema, I have already given the schema
    in design.md add these to migrations folder" — transcribed the already-finalized schema
    from design.md into migrations/000001_init_schema.{up,down}.sql via golang-migrate
    naming, then verified up/down against a real Postgres container (constraints, FKs,
    index all matched design.md) before tearing it down. No schema decisions made — the
    schema itself was already finalized by the user.

11. "have implemented domain will add entities,validation at a high level and also add
    request validation in handler.Will add more validations etc as required, also
    implemented request validation in handler. Review these" — reviewed the user's domain
    entities (Wallet/Transfer/LedgerEntry/IdempotencyRecord) and handler request validation
    as a senior-engineer pass: flagged a compile-breaking undefined-error reference, a
    Currency field/param mismatch against the finalized (currency-less) schema, and an
    unused errorResponseDTO/plain-text-vs-JSON inconsistency in error responses. No code
    written on the user's behalf — findings only; `go mod tidy` run to fix uuid's direct/
    indirect require grouping (dependency-wiring housekeeping).

12. "1. written but not saved. 2. removed, 3. implemented but not saved 4. handlers just
    validates request and domain.NewPendingTransfer gives transfer and i do not want
    transfer object with non-positive and with wallet validation. 5. good catch added 6.
    Added. Good one" — confirmed fixes in progress; clarified that domain-level validation
    in NewPendingTransfer is intentional defense-in-depth, not accidental duplication with
    the handler's own checks.

13. "saved" / "saved right now check" — re-checked disk state each time; first pass still
    showed unsaved edits (Wallet.Currency and unused currency params still present), second
    pass confirmed they were gone and the build/vet/gofmt were clean.

14. "remove the currency field from the DTO" — found the user had already removed
    transferResponseDTO.Currency before the request landed; confirmed build still clean.

15. "commit this" — committed the domain model (Wallet, Transfer state machine,
    LedgerEntry, IdempotencyRecord) and CreateTransfer request validation.

16. "commit and push this. let me service and repository" — committed/pushed the domain
    model and handler validation; user went on to implement the service and repository
    layers independently.

17. "I have implemented service and repository. Please review" — reviewed
    TransferService, WalletRepository, TransferRepository, LedgerRepository,
    IdempotencyRepository as a senior-engineer pass. Confirmed the idempotency
    check-then-insert race is handled correctly (matches design.md's documented approach)
    and wallet balance math is race-safe under the FOR UPDATE lock. Flagged one real
    concurrency bug: `processTransfer` read the transfer row without `FOR UPDATE`, so two
    concurrent replays of the same idempotencyKey could both pass the PENDING check before
    either committed — the second would double-attempt the debit/credit and only get saved
    by the `ledger_entries` UNIQUE constraint rolling it back, but it would then surface a
    raw Postgres error instead of the idempotent success response. Also flagged a
    delimiter-collision risk in computeRequestHash's string formatting, and an untruncated
    timestamp inconsistency in IdempotencyRecord.CreatedAt. No code written on the user's
    behalf — findings only.

18. "Locked transfer row in processTransfer, In computeRequestHash removed |,
    IdempotencyRecord.CreatedAt truncated to microsecond. GetWalletByID existence check ...
    Did not get LockForUpdate suggestion" — re-checked the fixes: found the build broken
    (service.go called a new TransferRepository.LockForUpdate method that hadn't been
    added/saved yet); clarified the WalletRepository.LockForUpdate note was a minor
    round-trip efficiency remark (batch the two single-row locks into one `WHERE id =
    ANY($1)` query), not a correctness issue.

19. "forgot to save check now." — rebuilt clean; confirmed the new
    TransferRepository.LockForUpdate (with FOR UPDATE, called before the pending-status
    check) correctly closes the race. Flagged that the delimiter fix in computeRequestHash
    had gone the wrong way — removing the separator entirely (`"%s%s%d"`) made field-boundary
    collisions *more* likely, not less, since now any split point between the two wallet IDs
    collides.

20. "fixed the delimiter, GetTransferByID we moved to lock own right?" — confirmed
    processTransfer now calls the new LockForUpdate instead of GetTransferByID (which is
    unused for now, not necessarily dead — could back a future read endpoint). Ran `go vet`
    on the "fixed" delimiter and caught a real bug: the format string had 5 verbs against 3
    arguments, feeding a string into a `%d` verb — Sprintf doesn't panic on this, it embeds
    an error placeholder and keeps going, so the hash was being computed over garbled text.

21. "Aah check now" — verified the corrected computeRequestHash (length-prefixing each
    field: `len(fromWalletID):fromWalletID,len(toWalletID):toWalletID,amount`), which
    unambiguously encodes field boundaries regardless of wallet-ID contents. go vet, go
    build, and gofmt all clean.

22. "commit and push this, along with the our history" — committing the service/repository
    implementation and this log.

<!-- Append new prompts below, in order, as the session continues. -->
