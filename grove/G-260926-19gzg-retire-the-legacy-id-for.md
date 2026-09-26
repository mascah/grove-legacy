---
id: "G-260926-19gzg"
type: work
title: "Retire the legacy ID form from the binary"
status: active
created: "2026-09-26T16:09:19Z"
updated: "2026-09-26T19:53:31Z"
kind: refactor
size: small
depends_on: ["G-260926-vkv48"]
relates_to: ["G-260926-yvjy6", "G-260926-2da4n", "G-260926-pgj43", "G-260921-czt8x"]
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

## Next

Depends on G-260926-vkv48 because this deletes the command it adds and
the validator change would stop that command loading legacy records.
Assign only after the owner has run the command in nullsec and committed
there, which is the owner's step outside this repository.
