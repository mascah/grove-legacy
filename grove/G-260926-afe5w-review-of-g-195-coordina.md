---
id: "G-260926-afe5w"
type: review
title: "Review of G-195 coordination-free record IDs"
status: current
created: "2026-09-26T05:27:39Z"
updated: "2026-09-26T05:27:57Z"
work: ["G-195"]
examined: "a9a63c2"
---

## Examined

Three rounds by the `grove-reviewer` agent, a fresh one per round, read-only,
in `.claude/worktrees/worktree-G-195` on branch `worktree-G-195`:

1. `fa4acc1..653a8b6`: the whole implementation and documents, against
   [G-195](G-195-coordination-free-record-ids.md)'s acceptance, the plan
   [G-199](G-199-coordination-free-record-ids.md) and
   [G-194](G-194-identify-records-by-creation-dat.md). The reviewer ran
   `go vet`, `gofmt`, `check`, the full suite and its own two-clone
   demonstration under `/tmp` (four date-form IDs, `check` OK after the
   pull, only `write.lock` in each `.git/grove`), and checked a hand-authored
   `G-1234` fails with the built binary.
2. `653a8b6..af6d0be`: the fixes, with `python3 evals/run.py selftest`
   (`selftest: ok`), `--help`, the four shipped guides and a tree-wide grep
   for the old scheme.
3. `af6d0be..a9a63c2`: the forward links on G-004 and G-064.

The examined commit is `a9a63c2`; the candidate adds only this record and
G-195's evidence after it.

## Findings

Round 1, `Open findings: 3`:

1. should-fix: `evals/run.py` read an ID as a filename's first five
   characters, breaking the selftest and every eval that seeds records.
2. should-fix: `grove --help`, the README index and a `convert` comment
   still said "next shared ID" or "reserving".
3. should-fix (knowledge): G-006 stayed `accepted` though G-194 retires its
   mechanism.
4. note: the board sorted a candidate group's IDs bytewise.
5. note: no `versions` test orders a mixed set; it uses `CompareIDs`, tested
   in `project` and `deps`.
6. note: a stray blank line in `create.go`'s imports.
7. note: with a 14-character ID a card's version tag clips below about
   29-column board columns (acceptance 7).
8. note: the brief's sequential-ID sentence, the owner's per G-195's Next.
9. note: `TestOwnerLost` failed once in the reviewer's full run and passed
   on rerun; the change does not touch it.

Round 2, `Open findings: 1`: A. should-fix (knowledge): G-004 and G-064,
which G-194 revises, did not link forward to it. Every round-1 fix was
confirmed and no regression found.

Round 3: finding A closed, no regression.

Open findings: none

## Disposition

- 1: fixed in `af6d0be` with a `record_id` helper; selftest passes.
- 2: fixed in `af6d0be`.
- 3: fixed in `af6d0be`: G-006 `superseded`, `relates_to` G-194, with a
  note.
- 4 and 6: fixed in `af6d0be`.
- 5: not done, for the reason given.
- 7 and 8: left to the owner.
- 9: unrelated; the full suite passed at `653a8b6`, `af6d0be` and `a9a63c2`.
- A: fixed in `a9a63c2`: both keep `accepted`, gain G-194 in `relates_to`
  and a line naming what it revised, as G-001 did for G-004.
