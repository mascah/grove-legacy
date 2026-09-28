---
id: "G-260928-r1hkh"
type: work
title: "Show approval and merge standing on review cards and bound the changes list"
status: proposed
created: "2026-09-28T19:28:59Z"
updated: "2026-09-28T19:34:01Z"
kind: feature
size: small
relates_to: ["G-260927-60ffq", "G-260925-h8rj5", "G-260925-dzxm6", "G-260921-jwk4e", "G-260920-svpbc", "G-260920-z8vfp"]
---

## Outcome

In the Review column the owner sees on each card whether its candidate is
approved, and by whom, and whether it currently merges cleanly into the
target; and the changes list under a candidate stays one line per file.

Owner intent, shaping conversation 2026-09-28: "Items in the review column
should show whether they've been approved or not. They should also show an
indicator if they're unable to be merged due to conflicts"; the changes
section "has a ton of noise junking up the sidebar from the back references".

## Constraints

Observed at main `6fbb888`, 2026-09-28:

- A card (`cardBox`, `internal/tui/view.go`, around 38-77) shows the ID, a
  tag, the title and one metadata row. Tags come from `model.go` (around
  1147-1159: `⑂ N states`, `uncommitted`, `not on TARGET`) and `attemptTag`
  (`● running`, `● orphaned`). `approval()` (`review.go`, around 308)
  returns `approved` or `approved under policy` and is used only by the
  detail's standing line and the Review block
  ([G-260927-60ffq](G-260927-60ffq-show-whether-an-approval.md));
  `conflicted()` (`review.go`, around 595) only by the footer hint and `m`.
- A merge prediction costs three Git processes per candidate and runs only
  when a card is open, a preview is read or `deps` runs, never during the
  board load ([G-260925-h8rj5](G-260925-h8rj5-predict-whether-a-candid.md)
  acceptance 4; [G-260920-svpbc](G-260920-svpbc-show-a-work-item-s-linea.md),
  [G-260920-z8vfp](G-260920-z8vfp-keep-a-full-board-load-f.md)).
- `changesSection` (`review.go`, around 387-423) draws each file as one
  clipped row and, under it, `described by …` from `describedBy` (around
  433-468, [G-260925-dzxm6](G-260925-dzxm6-search-record-bodies-and.md)),
  word-wrapped without bound, so one file can take many rows in a narrow
  sidebar.

Proposed design, labelled proposed:

- Card tags `approved` and `approved under policy` from `approval()`, read
  from the record, in Review and Done.
- A `conflicts` tag on Review cards computed after the board has drawn: one
  prediction per card in the column, cancelled by any re-read or key that
  starts another read, shown when ready with the target commit it read, and
  never part of the load. Where the prediction cannot run, no tag.
- One row per file: the described-by row becomes a count on the file's row
  (`described by 3`), and the names move to the head of the diff that Enter
  opens on the file.

Out of scope: a dedicated code review screen, which the owner may shape
after seeing this; changing what `deps` or the Review block say.

## Acceptance

1. A Review card of an approved candidate shows its tag with the owner's and
   the policy's told apart; a Done card keeps it.
2. A Review card whose candidate conflicts with the target shows `conflicts`
   once the after-load read completes; a re-read during it cancels it; the
   board load itself runs no prediction, asserted the way the load tests
   count processes.
3. The changes list shows one row per file at every width the fit tests
   cover, and the diff head names the records that describe the file.
4. `docs/board.md` says so; the owner judges the cards in a terminal.

## Next

Assign: `/grove-work G-260928-r1hkh`. Touches `internal/tui` in places
separable from [G-260928-csg91](G-260928-csg91-make-the-board-s-prompts.md),
[G-260928-y50a4](G-260928-y50a4-edit-any-record-in-your.md) and
[G-260928-63124](G-260928-63124-show-each-member-s-state.md), which needs no
edge.

No plan needed: small, and the proposed design above names each change and where it lands in `internal/tui` (card tags, one after-load prediction read per Review card, the changes list); implementation follows it.
