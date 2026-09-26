---
id: "G-260926-yvjy6"
type: decision
title: "Retire legacy IDs by renaming every record to the date form"
status: accepted
created: "2026-09-26T16:09:19Z"
updated: "2026-09-26T16:10:26Z"
relates_to: ["G-260926-2da4n", "G-260919-4h6pn", "G-260921-gtydy", "G-260921-r491p", "G-260926-pgj43", "G-260921-905y3", "G-260921-czt8x", "G-260923-gsthp", "G-260926-vkv48", "G-260926-19gzg"]
---

## Decision and authority

On 2026-09-26, in an interactive shaping session held after
[G-260926-pgj43](G-260926-pgj43-coordination-free-record.md) merged, the owner chose to
rename every legacy `G-NNN` record in this repository and in nullsec to the
date form that [G-260926-2da4n](G-260926-2da4n-identify-records-by-crea.md) selected,
and to retire the legacy form from the validator, rather than keep the two
forms coexisting as G-260926-2da4n selected the day before. In their words: "I'd
prefer not having mixed history that has confusing intentions for LLMs going
forward. We don't have any external users of this project yet, I've only
just 'converted' nullsec from the legacy format to the G-XXX format, but not
done any development since, so breaking changes are safe to make."

With the evidence below in front of them, the owner also settled the choices
the rename raises:

1. Filename slugs are re-derived from the title at the 24-character cap, so
   the longest filename stays 42 characters.
2. Every mention is rewritten, historical literals included: a quoted branch
   name such as `worktree-G-260926-pgj43` in a record's evidence reads under the new
   ID even though Git history keeps the old one. The migration map is the key
   for reading old commits.
3. `formerly` is left untouched: it keeps naming the typed ID or path a record
   replaced in [G-260921-r491p](G-260921-r491p-reconcile-all-grove-cont.md), and the map carries the
   legacy ID to its date-form ID.
4. The legacy form leaves the validator, ordering code, tests, help text and
   documents in the same effort; test fixtures are converted in bulk. The
   owner questioned a larger estimate for this, and a measurement (below)
   showed it small.
5. Finished attempts stay viewable in the board, so the attempt directories
   under the Git common directory are renamed and rewritten with the records.

Selected: an ID is a date-form ID and nothing else. A record's date-form ID
takes the record's `created` date; the sequence number a legacy ID carried is
history that the map preserves, not identity. Nullsec's records are renamed
the same way in nullsec, by the owner, with the command this repository
ships; nullsec stays outside this repository's write scope
([G-260921-905y3](G-260921-905y3-cut-nullsec-over-to-this.md)).

## What this revises

- [G-260926-2da4n](G-260926-2da4n-identify-records-by-crea.md): its bullets that
  existing numeric IDs stay valid, are never renumbered and coexist in the
  validator, and its rejected alternative "Renumbering existing records to
  the new form". The date form itself, the random tail and the slug cap
  stand. G-260926-2da4n named "if the first release wants to fix one ID form and
  retire the legacy pattern" as a reason to reopen; this is that.
- [G-260919-4h6pn](G-260919-4h6pn-use-shared-sequential-id.md) is superseded: no sequential ID remains an
  identity. Its short-filename choice lives on in G-260926-2da4n's slug cap.
- [G-260921-gtydy](G-260921-gtydy-keep-identity-and-placem.md): stable identity and placement stand for
  ordinary edits. This is a second explicitly authorized one-time migration
  of the same shape as the one G-260921-gtydy authorized for G-260921-r491p, not automatic
  behaviour, and the last: after it, one ID form exists and there is
  nothing left to migrate.
- The brief's foundation sentence on "stable sequential IDs" and clone
  collision checks, which G-260926-2da4n left to the owner, now has its answer: one
  neutral namespace of date-form IDs issued without coordination. The owner
  directs that edit here; the retirement work carries it.
- The owner's standing policy that Grove keeps no backward compatibility
  before its first release (G-260921-r491p, 2026-09-21) applies unchanged: an old
  commit stays inspectable with the CLI in that commit.

## Evidence

Observed 2026-09-26 at main `a691f5f`, by inspection and in a disposable
clone under the scratch directory:

- 193 records carry legacy IDs and 1 the date form. Inside `grove/`: 728
  Markdown links between records, 217 frontmatter reference lines and
  4,625 legacy-ID tokens in all. Outside it: 86 in `docs/`, 11 in the
  README, 32 in `CLAUDE.md`, 30 in the brief, 63 under `evals/`, 3 in CI
  and hook configuration, 85 in non-test Go comments, and 1,856 fixture
  IDs in 46 test files plus 61 in `internal/tui/testdata/terminal.py`.
- Every legacy record has `created`, and legacy ID order agrees with
  `created` order with no inversion, so date-form IDs keep the order across
  days. 87 of the 193 slugs exceed 24 characters.
- In the clone, removing the legacy alternative from `IDForm`
  (`internal/project/metadata.go`) and one regex replace over the fixtures,
  `G-260919-6mpmw` to `G-260101-00001` and so on, left 7 failing tests in 5 packages:
  a case expecting `G-000` refused, two tests of legacy-first ordering, a
  fixed-width header in `TestDepsCLI`, a 20-column render, a clipped search
  column and a case-folded alias literal. Nothing else failed.
- Attempts live under `<common>/grove/attempts/<WORK.TIMESTAMP>/`; 39 exist
  here. `attempt.json` holds `work`, the selected IDs and record paths;
  `owner.log`, `result.json` and `events.jsonl` (about 80 MB in all) quote
  IDs. The listing matches names against `IDForm` and silently skips the
  rest (`internal/attempt/attempt.go`), and the board finds a card's
  attempts through `work`. Card history uses `git log --follow`
  (`internal/versions/history.go`), so a renamed record keeps its history.
- `create.Issue` already takes the time whose date it uses, and
  `create.Slug` already caps at 24.
- Nullsec (`../nullsec`, read only, installed `grove` at `a691f5f`): 127
  legacy records, `check` OK; 126 have no `created`, since `convert` writes
  none, and every one's `formerly` source has a first commit in its history,
  dated 2026-09-14 to 2026-09-22. 185 links, 114 frontmatter references and
  1,827 tokens in records; 372 in Rust and SQL comments; 18 in its README,
  `AGENTS.md` and `CLAUDE.md`. 14 merged local branches and 3 linked
  worktrees still hold pre-conversion records. No attempts.

## Alternatives

- **Coexistence, as G-260926-2da4n selected.** One alternation in the validator, no
  rename. Rejected by the owner: two forms in records, documents, code
  comments and history give a reader two intentions to reconcile, and no
  adopter depends on the legacy form.
- **Rename the records but keep the legacy form valid.** Avoids the fixture
  change. Rejected: measured at 7 tests, and a stale branch could then
  restore a legacy record that `check` would accept.
- **Rename only this repository.** Rejected: nullsec is the pilot adopter
  and its 127 records would keep the second form alive.
- **Do nothing.** Rejected for the same reason as coexistence.

## Reconsideration

Reopen if an adopter outside the owner's projects holds legacy records
before the external preview ([G-260923-gsthp](G-260923-gsthp-prepare-grove-for-extern.md)); none
does today. The renaming command is a one-time tool and leaves the binary
before that preview.
