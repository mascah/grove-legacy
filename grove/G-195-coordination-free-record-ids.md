---
id: "G-195"
type: work
title: "Coordination-free record IDs"
status: review
created: "2026-09-26T02:56:11Z"
updated: "2026-09-26T15:12:17Z"
relates_to: ["G-194", "G-004", "G-006", "G-064"]
candidate: "f6f180279905e07427b31851ca32569f92475992"
approved: "f6f180279905e07427b31851ca32569f92475992"
---

## Outcome

`grove new` and `grove convert` issue IDs that need no shared state, so
records created in separate clones of one repository, on different
machines, never collide, while every existing numeric ID keeps working
unchanged. The intent is the owner's, selected in
[G-194](G-194-identify-records-by-creation-dat.md) on 2026-09-25: shape on
a laptop and run attempts on an always-on machine without ID collisions.

## Constraints

In scope:

- The form G-194 selects, `G-YYMMDD-xxxxx`: the UTC date of creation, then
  five lowercase Crockford base32 characters (`0-9a-hjkmnp-tv-z`) from
  `crypto/rand`. Proposed: the validator accepts
  `^G-[0-9]{6}-[0-9a-hjkmnp-tv-z]{5}$` and the legacy `^G-[0-9]{3,}$` with
  canonical padding; only the date form is issued. Five characters keep
  the odds of two clones colliding on one day near 3 in a million; four is
  about 1 in 10,000, which is thin for a heavy user, and six is longer than
  the owner wants.
- Proposed: generation retries, a bounded number of times, a tail that the
  loaded project, any local ref or any worktree already holds, using the
  exact-ID search the floor scan already does with one `git grep`. The write
  lock still serializes publication with `update`.
- Delete the counter file, the allocator lock, the floor scan,
  `create.Allocate`, and the initialization and below-floor notices, with
  their tests. `convert` issues through the same generator.
- Where a list orders by ID (`versions` groups, `deps` rows): legacy IDs
  first in numeric order, then date-form IDs lexically, which is
  chronological. The board keeps its `updated` then `created` order.
- Derived slugs cap at 24 characters instead of 32; an explicit `--slug`
  keeps only its character check.
- The documents that own the contract: the record model's identity,
  placement and allocation sections and its clone paragraph; the commands
  document on `new` and `convert`; the shaping guide's sentence on clone
  counters; and this repository's `CLAUDE.md` contract line on the counter.
  Shipped documents keep their linking rules.

Out of scope:

- Renaming or renumbering any existing record; nullsec's records;
  `formerly` and [G-069](G-069-migration-map.md).
- Any fetch, push or remote coordination in Grove.
- Changes to the board's layout beyond what the longer ID forces in columns.
- An entrypoint revision: the entrypoints need nothing new of the binary.

Observed evidence, 2026-09-25 at main `e812672`: the numeric shape lives in
`IDPattern` and `validID` in `internal/project/metadata.go`, the allocator
and `idLine` in `internal/create/create.go`, the ID format in
`internal/update/convert.go`, and `compareIDs` in
`internal/versions/versions.go`. The board sorts by `updated` then
`created`. An attempt ID joins the work ID and a timestamp with `.`, and a
branch joins IDs with `-`; neither parses IDs back, so both accept the new
form. The derived slug cap is in `create.Slug`.

## Acceptance

1. In two disposable clones of one repository with no shared state, `new`
   on each issues date-form IDs, and after merging one into the other
   `check` passes. A test injects the tail generator to show that a tail
   already present in a local ref or worktree is not reused and that the
   retries are bounded.
2. Every record in this repository validates unchanged, `check` passes
   with no rename, and a hand-authored `G-1234` or `G-01` still fails;
   `new` and `convert` never issue a numeric ID.
3. `grove/neutral-ids` and `grove/lock` under the common directory are
   neither read nor written; `new` needs no network and creates only the
   write lock; the read-only commands still create nothing.
4. `versions` and `deps` list legacy IDs before date-form IDs, and the
   latter in date order; the board's order is unchanged.
5. A derived slug is at most 24 characters and the longest generated
   filename is 42 characters; an explicit `--slug` behaves as today.
6. The record model, the commands document, the shaping guide and
   `CLAUDE.md` describe the new form and no longer describe the counter,
   the floor or clone collision checks; `grove guide model` prints the
   updated model; every link resolves.
7. The owner opens the board and a directory listing with a few new
   records and judges the names readable in an actual terminal.

## Evidence

Implemented headless on branch `worktree-G-195`
(`.claude/worktrees/worktree-G-195`), based on main `fa4acc1`, from this
record at `sha256:dc126f3f…` and the plan
[G-199](G-199-coordination-free-record-ids.md) as committed in `613cfb1`.
Code and documents: `f5cbdfa`, `653a8b6`, `af6d0be`, `a9a63c2` (reviewed
tip). The candidate is the commit that adds this evidence and the review
record on top of `a9a63c2`, and changes nothing outside `grove/`.

Against acceptance:

1. `create.Issue` draws `G-YYMMDD-` plus five Crockford base32 characters
   from `crypto/rand`, under the write lock, and draws again, at most eight
   times, while the ID is held by the loaded project, an `id:` line on any
   local branch, remote-tracking ref or tag (one `git grep` for that day's
   prefix), or a live record in any worktree. Tests:
   `TestNewInSeparateClonesMergesClean` (two clones, three `new` each, pull,
   `check` OK: 8 records); `TestNewSkipsIDsInRefsAndWorktrees` (injected
   tails held only in branch history, only in another worktree, only here);
   `TestIssueBoundsItsDraws` (exactly 8 draws, then nothing created);
   `TestNewConcurrentAcrossWorktrees` (12 concurrent `new` from a source that
   repeats every tail: 12 distinct IDs). Mutating the grep pattern fails the
   ref test. By hand, two clones of this checkout each issued date-form IDs
   and `check` passed after the pull (197 records).
2. `project.IDForm` is `G-(?:[0-9]{3}|[0-9]{6}-[0-9a-hjkmnp-tv-z]{5})` and
   `G-000` is refused. `check` here: `OK: 193 records`, no rename.
   `G-1234`, `G-01`, `G-0001`, short, uppercase and `l` tails fail
   (`TestStrictRecordMetadata`). `new` and `convert` only issue through
   `Issue`.
3. `neutral-ids`, `grove/lock`, `create.Allocate` and `repo.AllocatorLock`
   are deleted; `TestNewCreatesOnlyTheWriteLock`,
   `TestCoordinationStateStaysUnderTheCommonDirectory` and the convert
   refusal test find only `write.lock`; the read-only commands still create
   nothing (`TestNewCreatesRecordAndReadCommandsLeaveNoState`). No fetch.
4. `project.CompareIDs` (length, then bytes) orders `versions` groups,
   `deps` rows, groups and outside prerequisites, candidate groups and the
   board's blocking and unlocks lists: legacy first, then date form by date
   (`TestOverviewOrdersLegacyThenDateForm`, which fails with a plain sort;
   `TestOrderUsesCreationThenNumericID`). Board card order is unchanged.
5. `create.Slug` caps at 24 (`TestSlug`); a generated filename is at most
   42 characters, as the demonstration showed; `--slug` keeps only its
   character check.
6. The record model (identity, placement, issue, conversion, Identity and
   dates, `new`), the commands reference (ID order), the shaping and work
   guides, README (opening and index) and `CLAUDE.md` (contract and fixture
   lines) describe the date form; `grove --help` too. `grove guide model`
   prints it; the shipped-document test now also checks date-form IDs.
7. Not judged: the owner's.

Decisions taken, within the outcome:

- Legacy IDs validate as exactly three digits, not the proposed `{3,}`,
  because acceptance 2 requires `G-1234` to fail. Every legacy ID here and
  in nullsec (highest G-127) is below G-1000.
- The ID is drawn under the write lock, so two `new` in worktrees of one
  repository serialize and see each other; nothing is reserved, and errors
  say "nothing created" or "nothing converted".
- Tests inject the tail through the package variable `create.tail`; tests
  outside `create` read issued IDs from `new`'s output.
- Contrary to the observed note above, `internal/attempt` parsed IDs back
  (`^[A-Z]+-[0-9]+$`): `run` and attempt listings now use `IDForm`
  (`TestRefusals`, `TestAttemptNames`).
- The board's deps preview and search pad the ID column to the longest ID,
  keeping today's widths as the minimum.
- `evals/run.py` reads IDs by pattern, not as five characters.
- Knowledge: G-006 is `superseded`, linked to G-194; G-004 and G-064 stay
  `accepted`, linked forward to G-194 with a line naming what it revised.

Verification at `a9a63c2`: `go vet ./...` clean; `gofmt -l .` empty;
`go run ./cmd/grove check` OK: 193 records;
`go test -count=1 -timeout 120s ./...` all packages ok; at `653a8b6`
`python3 internal/tui/testdata/terminal.py BINARY` all ok; at `af6d0be`
`python3 evals/run.py selftest` ok.

Review: [G-260926-afe5w](G-260926-afe5w-review-of-g-195-coordina.md), three
rounds, `Open findings: none`.

Limits: `-race` and the Linux container run were not done (no concurrency
primitive changed; the lock is the existing one). With 14-character IDs a
card's version tag clips below about 29-column board columns. The change is
about 1,100 lines, over the standing policy's `max_lines`, so it waits for
the owner.

## Next

In review, awaiting the owner's judgment, acceptance 7 included: open the
board and a directory listing with a few new records (for example, run
`go run ./cmd/grove new page "Probe"` in a disposable clone) in a real
terminal. Owner decisions outside this work: whether G-004 should become
`superseded` too, and the brief's foundation sentence on sequential IDs and
clone collision checks (grove/brief.md), which is the owner's to change.

To accept, in this checkout, then in main's:

```sh
go run ./cmd/grove approve G-195 "VERDICT"
go run ./cmd/grove integrate G-195
```

Verdict on candidate f6f1802, 2026-09-26: approved
