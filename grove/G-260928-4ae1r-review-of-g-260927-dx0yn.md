---
id: "G-260928-4ae1r"
type: review
title: "Review of G-260927-dx0yn attempt shape facts and totals"
status: current
created: "2026-09-28T00:59:46Z"
updated: "2026-09-28T01:00:05Z"
work: ["G-260927-dx0yn"]
examined: "e096c33"
---

## Examined

An independent `grove-reviewer` subagent (read-only, headless session of
2026-09-28) reviewed [G-260927-dx0yn](G-260927-dx0yn-retain-per-attempt-proce.md) in
two rounds on branch `worktree-G-260927-dx0yn`, from `main` at `38511c1`,
against the record's acceptance and plan
[G-260928-j8sc0](G-260928-j8sc0-attempt-shape-facts-and.md). Round 1 covered
`6fbc8ff`, the implementation. Round 2 covered `e096c33`, the fixes. The
reviewer ran vet, gofmt, `grove check`, the short suites of `attempt`,
`cli` and `tui`, the targeted tests, and `grove attempt` and `grove
attempts` on real attempts, one of them running. It hand-counted the
fixture and checked the turn sum against raw result events.

## Findings

Round 1 found nothing blocking and four low findings:

1. An unreadable HEAD at finish left `Changed` nil, which printed as "not
   recorded when it finished", the reason for an older attempt.
2. After a skipped oversized line, the first edit shown may not be the
   first, and nothing said so.
3. A guide or instructions file read as a file was not listed among the
   known misses.
4. No test exercised the paths `ShowContext` passes to `ReadShape`.

Knowledge: no new term is needed. The Attempt term points to
`docs/commands.md` "Attempts", which now defines process calls, `Shape:`
and `Changed:`. Nothing contradicts a settled term or decision. `Shape` also
names the shaping workflow; the record's design chose the label.

Round 2 found each fixed, with no regressions. It noted two things it did
not count as findings: the unreadable-HEAD default has no test, and the new
test does not exercise the `cwd` argument, built as `child.Dir` is.

Open findings: none

## Disposition

All four fixed in `e096c33`:

1. `reconcile` records `Changed{Error: "HEAD was unreadable when it
   finished"}` unless HEAD was read.
2. The skipped-lines note now says the first edit is the first in the lines
   read, and so do the docs.
3. The known misses in the docs include it.
4. `TestRunToResult`'s fake reads the record and edits a file in its
   worktree and asserts the classification. It was checked to fail with a
   broken root path.
