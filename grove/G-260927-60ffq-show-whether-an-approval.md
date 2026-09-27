---
id: "G-260927-60ffq"
type: work
title: "Show whether an approval was delegated or the owner's wherever a record's standing is shown"
status: proposed
created: "2026-09-27T22:13:37Z"
updated: "2026-09-27T22:28:18Z"
size: small
relates_to: ["G-260921-btyck", "G-260926-a8vyj", "G-260925-wh9ax", "G-260925-5wrn8", "G-260921-jwk4e"]
---

## Outcome

Wherever a record's standing is summarized, the owner sees at a glance
whether its approval was their own or the standing policy's.

Owner intent, conversation 2026-09-27: "we should easily be able to see if
something was policy approved or human approved."

## Constraints

Observed 2026-09-27 at main `fef102f`:

- Much of this exists. The record model requires a delegated verdict to
  begin `delegated under policy grove.yaml REVISION` and to be "told apart
  from the owner's own wherever the record is shown"
  (`docs/record-model.md:523-527`); decision
  [G-260925-wh9ax](G-260925-wh9ax-delegate-conflict-resolu.md) and the
  [Policy](G-260926-a8vyj-policy.md) term say the same. `update.Delegated`
  (`internal/update/review.go:203`) reads it from the verdict line; its only
  display use is the board's review screen, which says `approved under
  policy` (`internal/tui/review.go:309`). `grove show` prints the body,
  where the verdict line is last.
- Not shown: the board card's standing line (`internal/tui/detail.go:437-475`)
  lists kind, size, priority, candidate, `for`, examined, blocks and
  states, not approval or who gave it; `grove list` columns are ID, TYPE,
  STATUS and TITLE; `versions` shows no approval; `show --json` has no
  derived field for it.
- No delegated verdict exists yet: all 33 verdict lines in this
  repository's records are the owner's, so the distinction has never been
  seen in use.

**Proposed design.** One word from one function. The card's standing line
says `approved` or `approved under policy` whenever `approved` is set, in
review and done alike; `show --json`, and `list --json` where it exists,
carry `approved_by: owner|policy` derived by `update.Delegated`; the review
screen's wording stays. No schema change: the verdict line remains the
fact, and `list`'s text columns are unchanged.

## Acceptance

1. In the board, a record with `approved` set shows `approved` or `approved
   under policy` in its standing line, in review and in done.
2. `show --json` carries the derived field, from the same function the
   review screen and the sweep use.
3. A test with an owner verdict and a delegated verdict covers the card and
   the JSON.
4. `docs/board.md` and `docs/commands.md` say so; owner judgment in the
   terminal.

## Next

Assign: `/grove-work G-260927-60ffq`.
