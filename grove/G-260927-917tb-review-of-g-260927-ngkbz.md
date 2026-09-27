---
id: "G-260927-917tb"
type: review
title: "Review of G-260927-ngkbz at cc87c6b"
status: current
created: "2026-09-27T17:23:34Z"
updated: "2026-09-27T17:23:52Z"
work: ["G-260927-ngkbz"]
examined: "cc87c6bb9fb91cb054552af3497b6a1301d82072"
---

## Examined

[G-260927-ngkbz](G-260927-ngkbz-frame-the-readme-for-som.md) on `worktree-G-260927-ngkbz`: the change
from main `90e82d4` to `cc87c6b`, against the record's acceptance, its
six-section design and out-of-scope list, and the repository's instructions.
Two rounds, each by a fresh `grove-reviewer` agent, read-only: round 1 at
`558bce8`, round 2 at `cc87c6b`. No plan: the record's Next says why. The
candidate adds only this record and G-260927-ngkbz's evidence to `cc87c6b`.

## Findings

Round 1 (`558bce8`): nothing consequential, two minor.

- Acceptance met, each with evidence:
  - Order (read). The headings run intro, Install, Set up a project, The
    loop, Commands, Develop Grove. `go run` and record IDs appear only at
    README lines 122–123, inside Develop Grove.
  - `grove.yaml` (run). In a disposable project after `init`, the README's
    yaml block written verbatim passed `check` (`OK: 0 records`). `target`,
    `run:` (`budget`, `permission_mode`, `model`, `effort`) and `policy:`
    agree with the record model's "Configuration and discovery".
  - Links (run). A script resolved every relative path and anchor in
    README.md, CLAUDE.md and docs/commands.md. The evidence-record,
    legacy-ID, evals and CLAUDE.md rows are gone or moved into Develop
    Grove.
  - `init` note (run). The note names `grove guide model`, which prints the
    heading it names. `init_test.go` asserts it, and `Init` agrees.
  - CLAUDE.md (read). Its README line matches.
  - Checks (run). `go vet`, `gofmt -l`, `check` and `go test -short
    ./internal/cli` passed.
- Checked (read): what `init` writes, the board keys (`a`, `f`, `i`, `R`),
  `feedback` returning work to `active`, the upgrade path, `just install`,
  and the Go version.
- Minor 1: the `target` bullet said only that the board marks work and that
  `done` is written there. `integrate`, `resolve` and `sweep` refuse without
  a target (`grove --help`; `internal/integrate`, `internal/attempt`,
  `internal/sweep`).
- Minor 2: G-260921-czt8x said `W-001` appears "in the README and the
  record model", but the README no longer has it.
- Knowledge: no new concept without a term, and no contradiction of a
  settled term or accepted decision.
- Noted, not findings:
  - `@latest` resolves to `v0.1.0`, which predates the new `init` note, so
    adopters see it only after the next tag. A release is out of scope.
  - The install-from-clone and `@COMMIT` lines were dropped, as the design's
    Install item intends.

Round 2 (`cc87c6b`): both minor findings resolved, and nothing new.

- Minor 1: resolved. The bullet names `integrate`, `resolve` and `sweep`,
  each confirmed against `--help` and the code's `Target == ""` refusal. It
  stays one line, as the design asks.
- Minor 2: resolved. The record now says "in the record model" only.
- Nothing else changed since `558bce8` (`git diff --stat`). `check`,
  `go test -short ./internal/cli`, `go vet` and `gofmt -l` passed.

Open findings: none

## Disposition

Both minor findings were fixed in `cc87c6b`. This review is evidence, not
approval.
