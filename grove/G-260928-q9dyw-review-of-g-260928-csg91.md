---
id: "G-260928-q9dyw"
type: review
title: "Review of G-260928-csg91 and G-260928-y50a4: prompts, refusals, edit and paging"
status: current
created: "2026-09-28T21:01:33Z"
updated: "2026-09-28T21:02:05Z"
work: ["G-260928-csg91", "G-260928-y50a4"]
examined: "9deee95e99c817ce9bcd6bb14255a9cd93cf5ead"
---

## Examined

[G-260928-csg91](G-260928-csg91-make-the-board-s-prompts.md) and
[G-260928-y50a4](G-260928-y50a4-edit-any-record-in-your.md) on
`worktree-G-260928-csg91-G-260928-r1hkh-G-260928-y50a4-G-260928-63124`,
base main `453add4`, within the combined diff of the four selected records,
including the helpers they share with
[G-260928-r1hkh](G-260928-r1hkh-show-approval-and-merge.md) and
[G-260928-63124](G-260928-63124-show-each-member-s-state.md) (footer and row
budget, `key()` and `moved()`, `cardBox`, the alert). Three rounds, each by
a fresh read-only `grove-reviewer` agent, against each record's
acceptance and constraints and CLAUDE.md's code rules: round 1 at
`48fd846`, round 2 at `8058c6c`, round 3 at `9deee95`. Each ran
`go vet ./...`, `gofmt -l .`, `grove check`, the short tests of `tui` and
`cli`, and the pseudo-terminal script against a built binary (all pass),
and checked its regressions by probe tests in a scratch copy.

## Findings

Round 1 (`48fd846`), five:

1. Medium: `clampScroll` assumed a one-row footer, so under a refusal `G`
   stopped short of a detail's last rows and the header's count was wrong.
2. Medium: a refusal stayed on screen after the refused action succeeded
   without a prompt (`e` on a non-question, `o`, the chooser), beside the
   success message, contrary to `docs/board.md`.
3. Low: `e`'s timeline refusal said "Esc returns to it", which now takes
   two Escapes.
4. Low: a refusal was cut at three rows with no marker.
5. Low: deps' `p` refusal still used the one-shot notice.

Round 2 (`8058c6c`), one: Low, the dependency trees' scroll still
assumed a one-row footer, so `G` stopped short there under a refusal (the
same defect in a caller round 1's fix did not reach). Noted, not counted:
`moved()`'s page size ignores the footer, and Enter on a linked record
keeps a refusal since the screen kind is unchanged; both predate the
branch or follow the documented rule.

Round 3 (`9deee95`): the fix verified by removing it (the test fails),
`clampScroll`'s fallback keeps the base behaviour, and the pseudo-terminal
script's `whole()` repaint was judged sound: a resize repaints every cell
from the model, so it cannot pass text the model lacks. No findings.

## Disposition

Round 1, in `6a28f16`: scrolling counts the footer's rows (finding 1); a
prompt, another screen, a preview or the editor opening settles a refusal
the key did not set (2); the text now says "so return to it first" (3); a
refusal takes up to half the screen less one row and ends in `…` (4); `p`'s
refusal is an alert (5). Each has a regression (`TestRefusalRowsAndTheBody`,
`TestEditAnyRecord`, `TestDepsPreview`), the scroll one verified failing
without its fix. Round 2, in `b6f92ee`: the trees use the footer's height,
with a regression in `the trees scroll` verified failing without it.

Open findings: none
