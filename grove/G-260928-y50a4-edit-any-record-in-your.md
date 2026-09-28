---
id: "G-260928-y50a4"
type: work
title: "Edit any record in your editor from the board, page like vim, and read a report alone"
status: active
created: "2026-09-28T19:28:59Z"
updated: "2026-09-28T20:18:06Z"
kind: feature
size: small
relates_to: ["G-260924-wp2pe", "G-260923-895zb", "G-260924-nqkkh"]
---

## Outcome

From the board the owner can open any record's file in their own editor, in
the checkout that holds it; move through long content with vim's paging
keys; and select an attempt's final report without the activity beside it.

Owner intent, shaping conversation 2026-09-28: "Needs more vim integration.
I love being able to pop open a question record and slice through it with
vim"; "I can't copy and paste text from the final report without also
selecting event activity text."

## Constraints

Observed at main `6fbb888`, 2026-09-28, in `internal/tui`:

- `e` exists only for an open question (`detail.go`, around 174;
  `attempts.go`, around 602 and 620; `answer.go`). The editor launch is
  `tea.ExecProcess` of `$VISUAL`, else `$EDITOR`, else `vi` (`answer.go`,
  around 186-191), on the one checkout on the branch of the shown version
  (`checkoutOf`), with a revision check before and after
  ([G-260924-wp2pe](G-260924-wp2pe-answer-a-blocking-questi.md)).
  `docs/board.md` says "only questions are edited": that was G-260924-wp2pe's
  scope, not a rule with reasons elsewhere.
- Movement is ↑/↓, `j`/`k`, PgUp/PgDn (`scrollKey`, `model.go`, around 909).
  No `gg`, `G` or half-page keys.
- On the attempt screen the report and the activity are joined by hand into
  one row list from 100 columns (`attemptRows`, `attempts.go`, around
  997-1015) and scroll together; `attemptKey` (around 610-632) handles only
  `x`, `o`, `e` and `d`. The detail's `w` (`detail.go`, around 152-158)
  hides its sidebar for the session
  ([G-260923-895zb](G-260923-895zb-make-attempts-easy-to-sc.md) chose the
  side-by-side layout).

Proposed design, labelled proposed:

- `e` on any record's detail opens its file in the editor in its checkout,
  the same way as a question, without the `## Answer` heading and the
  resolve prompt; a question keeps both. After the editor the board re-reads
  and says whether the file changed; an edited record shows as `uncommitted`
  as today. Refused as today: no checkout on the branch, two, or a file
  changed since the read; also refused on a timeline commit and on a
  version that is not live.
- `gg`, `G`, Ctrl-d and Ctrl-u wherever `j` and `k` move; `/` stays search.
- `w` on the attempt screen hides the activity so the report takes the
  width, for the session; below 100 columns the stacked layout keeps the
  report first, as today.

## Acceptance

1. `e` on a work, term, decision, plan, review or page detail opens the file
   in the editor in its checkout; the pseudo-terminal script covers one
   non-question record end to end; every refusal writes nothing.
2. `gg` and `G` reach the top and bottom of a detail, an attempt and the
   dependency list; Ctrl-d and Ctrl-u move half a page; the key table lists
   them.
3. `w` on an attempt hides the activity, and a selection across the report
   takes no activity text; the bounded-activity test covers both states.
4. `docs/board.md` says so; the owner judges in a terminal with their editor.

## Next

Assign: `/grove-work G-260928-y50a4`. Touches `internal/tui` in places
separable from [G-260928-csg91](G-260928-csg91-make-the-board-s-prompts.md),
[G-260928-r1hkh](G-260928-r1hkh-show-approval-and-merge.md) and
[G-260928-63124](G-260928-63124-show-each-member-s-state.md), which needs no
edge.

No plan needed: small, and the proposed design above names each change and where it lands in `internal/tui` (`answer.go`'s editor launch, `scrollKey` and `moved`, `attemptRows`); implementation follows it.
