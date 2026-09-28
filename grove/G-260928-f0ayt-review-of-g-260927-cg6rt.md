---
id: "G-260928-f0ayt"
type: review
title: "Review of G-260927-cg6rt at 75675e7"
status: current
created: "2026-09-28T00:53:51Z"
updated: "2026-09-28T00:54:11Z"
work: ["G-260927-cg6rt"]
examined: "75675e7d8c021ae6307ec7bf66bd6109d1fd9567"
---

## Examined

[G-260927-cg6rt](G-260927-cg6rt-print-the-work-guide-in.md) on
`worktree-G-260927-cg6rt`: the change from main `38511c1` to `75675e7`,
against the record's acceptance 1 to 3 and 5, its proposed design, the
[Entrypoint revision](G-260925-m9jcr-entrypoint-revision.md) and
[Shipped document](G-260925-khfe7-shipped-document.md) terms, and the
repository's instructions. Acceptance 4, a real `grove run` attempt, was out
of scope. Two rounds, each by a fresh `grove-reviewer` agent, read-only:
round 1 at `e50b955`, round 2 at `75675e7`. No plan: the record's design
served as one.

## Findings

Round 1 (`e50b955`): nothing blocking, two minor.

- Acceptance 1 (run). `guide work --entrypoint 2` printed 12,995 bytes,
  ending with "The rest of this guide" and its table. The head and all seven
  parts, from a built binary, concatenated and compared with `cmp`, equal
  `docs/work-execution.md`, and so does `--part all`. `--part nope` exits 2,
  prints nothing on stdout and names the parts.
- Acceptance 2 (run, read). The byte-for-byte concatenation puts every `##`
  section in exactly one part, with "When a human decision is missing" in
  `checkpoint`. The table holds only when to read each part.
- Acceptance 3 (run). The shipped-document checks pass. The revision
  decision was still to be written into the record at handoff.
- Acceptance 5 (read). `docs/commands.md`, README.md and CLAUDE.md agree,
  apart from minor 2. `evals/README.md` is untouched.
- Deviations judged sound: one Read in stages row pointing at the closing
  table, keeping one owner of the part list; the shipped-document tests
  reading `guide work --part all`, so they still see the whole file (their
  assertions are unchanged, the invocation is not); `docs/work-review.md`
  naming step 6 by `--part review`; the shorter intro, which lost no
  instruction. The head is 5 bytes under the test's 13,000-byte cap, so an
  edit to steps 1 to 3 will trip it: a working guard.
- Entrypoint revision: keeping 2 holds. The entrypoint asks only for
  `guide work --entrypoint 2`; an older binary prints the whole guide; a
  newer one prints a head whose `--part` the same binary serves.
- Links that now cross parts each map to a part through the table's Holds
  column. Exit codes checked on a built binary.
- Minor 1: `guide wrok --part review` said "--part applies only to guide
  work", hiding the unknown guide name (`internal/cli/cli.go`).
- Minor 2: in `docs/commands.md`, "asks for one" after the new text read as
  asking for a part.
- Knowledge: no new shared concept needs a term; "head" and "part" describe
  one command's output, owned by `docs/commands.md`. Nothing contradicts a
  settled term or decision.

Round 2 (`75675e7`): both resolved, nothing new. `guide wrok --part review`,
`guide --part review` and `guide work extra --part review` report the guide
name; `guide shape --part all` reports `--part`; the test cases fail without
the fix. `go vet`, `gofmt -l`, `check` and `go test -short . ./internal/cli`
passed.

Open findings: none

## Disposition

Both round 1 findings were fixed in `75675e7` and confirmed in round 2.
