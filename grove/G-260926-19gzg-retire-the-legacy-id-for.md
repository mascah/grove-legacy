---
id: "G-260926-19gzg"
type: work
title: "Retire the legacy ID form from the binary"
status: accepted
created: "2026-09-26T16:09:19Z"
updated: "2026-09-30T13:18:01Z"
kind: refactor
size: small
depends_on: ["G-260926-vkv48"]
relates_to: ["G-260926-yvjy6", "G-260926-2da4n", "G-260926-pgj43", "G-260921-czt8x"]
candidate: "33f78f2d7c228804a1c0e4c2fd78161213f24685"
approved: "33f78f2d7c228804a1c0e4c2fd78161213f24685"
approved_by: owner
approved_context: "sha256:49ffeeb1bed094892856219c84940bf3dcf1e2b02b6ca48bf38e1127a459cac1"
---

## Outcome

The binary, its tests and its documents know one record ID form, the date
form, and nothing in this repository outside Git history and the migration
map speaks of legacy IDs. The one-time renaming command is gone. The intent
is the owner's, decided in
[G-260926-yvjy6](G-260926-yvjy6-retire-legacy-ids-by-ren.md) on 2026-09-26.

## Constraints

In scope:

- `IDForm` in `internal/project/metadata.go` accepts only
  `G-[0-9]{6}-[0-9a-hjkmnp-tv-z]{5}`; the diagnostic names only that form.
  `CompareIDs` and its legacy-first rule go, and every caller orders IDs as
  strings, which is date order. `grove --help` and the attempt-name pattern
  follow, where they still mention or allow the other form.
- Test string literals that hold fixture IDs become date-form IDs by one
  deterministic rule, proposed as `G-NNN` to `G-260101-00NNN`, across the 46
  test files, `internal/tui/testdata/terminal.py` and `evals/`. Comment
  lines citing real records were already repaired by
  [G-260926-vkv48](G-260926-vkv48-rename-legacy-records-to.md) and are not
  fixtures. The seven tests the measurement found are then fixed or
  deleted: the `zero_id` case of `TestStrictRecordMetadata`,
  `TestOverviewOrdersLegacyThenDateForm` and
  `TestOrderUsesCreationThenNumericID` (delete: nothing legacy is left to
  order first), and the fixed widths in `TestDepsCLI`,
  `TestRenderedRowsHoldOnlyGlamourStyles`, `TestSearchReachesEveryRecord`
  and the alias literal in `TestAliasesAreIncludedAndChargedOnce`.
- The shipped guides and the model use date-form example IDs (`G-260920-svpbc` and
  `G-260920-z8vfp` in the work guide, `G-260921-w9x25` in the shaping guide, `G-260919-6mpmw`, `G-260919-rt9h9`
  and `G-999` in the record model), and so do `CLAUDE.md`'s invocation
  examples and the evals' README. The shipped-document test keeps
  requiring that the guides name only their own example IDs.
- The record model's identity paragraph, the commands reference's sentence
  on ID order, and `CLAUDE.md`'s contract line stop describing legacy IDs
  as valid. The brief's foundation sentence on stable sequential IDs and
  clone collision checks becomes: one neutral namespace of date-form IDs
  issued without coordination; the owner directed this edit in
  G-260926-yvjy6.
- Delete the renaming command, its tests and its documentation.
- A search for `legacy` and for `G-[0-9]{3}` across `internal/`, `cmd/`,
  `docs/`, the README, `CLAUDE.md`, the brief and `evals/` finds nothing
  but [G-260921-czt8x](G-260921-czt8x-identity-and-path-migrat.md)'s tables and records' historical
  prose.

Out of scope:

- Any record's ID or filename: G-260926-vkv48 did that.
- Nullsec: it must already be renamed when this runs, since afterwards the
  installed binary refuses its legacy records.
- Any change to how IDs are issued.

Observed evidence, 2026-09-26 at main `a691f5f`, in a disposable clone:
with the legacy alternative removed from `IDForm` and one regex replace
over 53 files, `go test -short ./...` failed exactly the seven tests named
above and nothing else; `go build` and `go vet` were clean. The
attempt-name pattern already composes `IDForm`.

## Acceptance

1. `check` passes here; a hand-authored `G-260919-6mpmw` record fails `check` on its
   ID; `new` and `convert` issue the date form as before.
2. `go test -count=1 -timeout 120s ./...`, `go vet ./...`, `gofmt -l .` and
   `python3 internal/tui/testdata/terminal.py BINARY` pass, with no fixture
   or golden output holding a legacy ID.
3. `grove guide work|shape|review|model` and `grove --help` print no legacy
   ID and no word "legacy"; the shipped-document test passes with the
   date-form examples.
4. The searches under Constraints find only G-260921-czt8x and historical prose in
   records; `deps` and `versions` order date-form IDs by date.
5. The renaming command is absent from `grove --help`, the commands
   reference and the record model.
6. The brief's foundation sentence reads as G-260926-yvjy6 directs, and the
   record model, commands reference and `CLAUDE.md` describe one form.

## Evidence

Headless assignment `G-260926-19gzg --interaction headless`, branch
`worktree-G-260926-19gzg` in `.claude/worktrees/worktree-G-260926-19gzg`,
base `main` `2bd50b9`, started from this record at
`sha256:220fcd987cc873b02738e8587a4ed89659400dabd2832a418c95b3db7aa2df95`.
No plan: the work is small and the Constraints name every change. Nullsec
was confirmed renamed first (read only: no `G-NNN-` file under its
`grove/`, the three "Grove renumber" commits at its HEAD).

Commits: `eaa0f84` code and tests; `4ff88f3` documents; `db568f3` review
fixes. Review [G-260926-7s7mz](G-260926-7s7mz-review-of-retiring-the-l.md),
two rounds, examined `db568f3`, "Open findings: none" once this Evidence
records finding 3. The candidate differs from `db568f3` only in this
record and that review.

**This record's own text.** G-260926-vkv48's rename rewrote the legacy IDs
in this record too, so its Constraints and Acceptance cite real records
where the owner wrote examples. At `a9f8fce` they read: acceptance 1, "a
hand-authored `G-001` record fails `check`"; the guides' examples were
`G-030` and `G-031` (work), `G-037` (shape), and `G-001`, `G-003` and
`G-999` (model). The work was done and judged against that reading.

**Decisions taken, and why.**

- The `init --check` verdict `legacy` (a managed entrypoint with no
  revision line, not an ID) is now `unrevised`, in `entrypoints.go`,
  `init.go`, `attempt.go`, `--help`, `docs/commands.md` and term
  [G-260925-m9jcr](G-260925-m9jcr-entrypoint-revision.md). Acceptance 3
  forbids the word in `--help` and the Constraints' search covers
  `internal/` and `docs/`. Nothing parses the token and no entrypoint names
  it, so no new entrypoint revision. It is a visible change to `init
  --check` output: the owner may prefer another word.
- Example IDs: `G-260925-7k2qm` (the model's existing example) and
  `G-260925-8m3xd` in the work guide, `G-260925-7k2qm` in the shaping
  guide, `G-260924-2b8rc` and `G-260925-7k2qm` in the model's layout and
  elsewhere; none is a record. `G-999` went with the sentence that named
  it. The shaping guide's "answer G-NNN" placeholder reads "answer
  G-YYMMDD-xxxxx".
- `update.Set` existed only for `renumber` and went with it. Titles,
  fixture paths and test names that said "legacy" without meaning an ID
  were reworded, so the search finds nothing.
- Two refusal cases, `"G-"+"001"`, split so the search finds no fixture,
  keep a test failing if the three-digit alternative returns (review
  finding 1).

**Acceptance.**

1. `check` here: OK, 199 records. In a clone, a hand-authored `G-001`
   record: `grove/G-001-hand.md:2: id: expected a canonical ID, e.g.
   G-260925-7k2qm`, exit 1; `new` issued `G-260926-d3g1p`, `convert`
   `G-260926-avq3f`.
2. At `db568f3`: `go test -count=1 -timeout 120s ./...` all ok; `go vet
   ./...` clean; `gofmt -l .` empty; `python3
   internal/tui/testdata/terminal.py BINARY` (built from `4ff88f3`, TUI
   unchanged since) all 5 scenarios ok, exit 0. No fixture holds `G-NNN`:
   fixtures follow `G-NNN` to `G-260101-00NNN`. Of the named tests,
   `zero_id`, `TestOverviewOrdersLegacyThenDateForm` and
   `TestOrderUsesCreationThenNumericID` are deleted; `TestDepsCLI` and
   `TestSearchReachesEveryRecord` take the wider column,
   `TestRenderedRowsHoldOnlyGlamourStyles` looks for `00093-` (the longer
   link target wraps at widths 20 and 40), and the alias literal is
   `g-260101-00001.MD`. `TestRecordProblems` also asserted the old
   diagnostic text and follows it.
3. `grove guide work|shape|review|model` and `grove --help`: no match for
   `legacy`, `renumber` or a three-digit ID; the shipped-document test
   passes with the new example lists.
4. Case-insensitive `legacy` and `\bG-[0-9]{3}\b` over `internal/`,
   `cmd/`, `docs/`, the README, `CLAUDE.md`, the brief, `evals/` and the
   root Go files: nothing (the README's row naming "a three-digit `G-` ID"
   points to G-260921-czt8x, as it should). `versions` rows and `deps`
   groups list IDs in string order, which is date order.
5. `grove renumber` is an unknown command (exit 2); `--help`,
   `docs/commands.md` and `docs/record-model.md` do not name it.
6. The brief reads "One neutral `G-` namespace of date-form IDs issued
   without coordination"; the record model's identity paragraph and ID
   order, the commands reference and `CLAUDE.md`'s contract line (now
   citing G-260926-yvjy6) describe one form.

**Limits.** `evals/run.py` was edited (its ID pattern and fixture literals)
but not run: it spends model budget. The TUI and deps unit tests keep
typed `W-`/`S-` IDs, which predate this work and are not `G-NNN`.

## Next

In review, candidate in the commit that sets it. The diff is far over the
standing policy's `max_lines`, so it waits for the owner, who also judges
the `unrevised` rename above.

```sh
cd /Users/mascah/GitHub/mascah/grove/.claude/worktrees/worktree-G-260926-19gzg
go run ./cmd/grove approve G-260926-19gzg "VERDICT"
cd /Users/mascah/GitHub/mascah/grove
go run ./cmd/grove integrate G-260926-19gzg --cleanup
just install
```

After `just install`, the installed `grove` refuses any legacy record, in
nullsec too (already renamed there).

Verdict on candidate 33f78f2, 2026-09-26: approved

Migrated to schema 4, 2026-09-30: status done with approval of its candidate became status accepted by owner.
