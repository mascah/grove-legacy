---
id: "G-260926-bz4n5"
type: plan
title: "Coordination-free record IDs"
status: current
created: "2026-09-26T05:02:04Z"
updated: "2026-09-26T05:02:30Z"
work: ["G-260926-pgj43"]
---

## Design

Plan for [G-260926-pgj43](G-260926-pgj43-coordination-free-record.md), from its revision
`sha256:dc126f3f…` at base `fa4acc1`, under the form
[G-260926-2da4n](G-260926-2da4n-identify-records-by-crea.md) selects.

- **Validator.** `project.IDPattern` becomes
  `^G-([0-9]{3}|[0-9]{6}-[0-9a-hjkmnp-tv-z]{5})$`, and `validID` still
  refuses `G-000`. Legacy IDs are exactly three digits, not `{3,}` as G-260926-pgj43
  proposes: its acceptance 2 requires a hand-authored `G-1234` to fail, and
  today's `{3,}` with canonical padding accepts it. Every legacy ID in this
  repository and in nullsec is below `G-1000`, and no numeric ID is issued
  again, so nothing that exists stops validating. The date digits are not
  checked as a calendar date; the validator checks shape, as today.
- **Comparator.** `versions`' `compareIDs` (length, then bytes) already
  orders legacy before date form, and date form by date then tail, because
  every legacy ID is shorter. It moves to `project.CompareIDs` and `deps`
  uses it wherever it orders rows, groups and outside prerequisites, instead
  of a plain string sort, which would put `G-260925-…` before `G-300`.
- **Generator.** `create.ID(root, recordDir, showPrefix, taken, now)` draws
  a tail of five Crockford base32 characters from `crypto/rand` through the
  package variable `tail`, which tests replace. One `git grep` for
  `^id:.*G-YYMMDD-` over every local ref (heads, remotes, tags), plus a walk
  of every worktree's record tree, collects the IDs already used that day;
  a drawn ID found there or in the loaded project is drawn again, at most
  eight times, then `new` fails naming the bound. With 32^5 tails, eight
  misses mean a broken source, not bad luck.
- **Locking.** The ID is drawn under the write lock, after the reload `new`
  already does there, so two `new` in worktrees of one repository serialize
  and the second sees the first's file. The allocator lock, the counter, the
  floor, `create.Allocate` and their notices are deleted; nothing is
  reserved, so the "reserved but not created" wording goes too. `convert`
  draws its ID the same way, under the lock it already takes.
- **Attempts.** Contrary to G-260926-pgj43's observed note, `internal/attempt`
  parses IDs back: `idPattern` (`^[A-Z]+-[0-9]+$`) refuses a date-form ID
  in `grove run`, and `attemptPattern` would hide its attempts. Both derive
  from `project.IDPattern`.
- **Slug.** `create.Slug` caps at 24; the longest generated filename is
  `G-YYMMDD-xxxxx-` plus 24 plus `.md`, 42 characters.
- **Tests.** Tests that expected `new` to issue `G-260919-6mpmw` read the ID from its
  output instead. Acceptance 1's two-clone check is a CLI test with two
  independent clones; the tail injection test lives in `internal/create`.

## Steps

1. `project`: pattern, `validID`, `CompareIDs`; tests for both forms and
   `G-1234`, `G-01`, `G-000`, bad tails.
2. `create`: generator, slug cap, delete allocator; `repo.AllocatorLock`
   deleted; `convert` switched. Tests: injected tails (taken in a ref, in a
   worktree, bounded retries), no `neutral-ids` or `lock` file, slug length.
3. `versions`, `deps`, `attempt` use the shared pattern and comparator.
4. Fix the tests that assumed sequential issue.
5. Documents: record model (identity, placement, allocation, clones),
   commands (`new`, `convert`), shaping guide, work guide's counter notice,
   README's opening sentence, `CLAUDE.md`'s counter lines. The brief is the
   owner's.
6. Verification per `CLAUDE.md`, two-clone demonstration, independent
   review, handoff.
