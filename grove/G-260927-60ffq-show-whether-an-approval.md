---
id: "G-260927-60ffq"
type: work
title: "Show whether an approval was delegated or the owner's wherever a record's standing is shown"
status: done
created: "2026-09-27T22:13:37Z"
updated: "2026-09-28T16:50:19Z"
size: small
relates_to: ["G-260921-btyck", "G-260926-a8vyj", "G-260925-wh9ax", "G-260925-5wrn8", "G-260921-jwk4e"]
candidate: "3af8383a1524aa97d9b0b35c468227aa2162eeb6"
approved: "3af8383a1524aa97d9b0b35c468227aa2162eeb6"
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

## Evidence

Compact handoff (`size: small`). Branch `worktree-G-260927-60ffq`, base
`991545c`; implementation `727557a` and `80ef0f7`; the candidate is the
commit adding this Evidence. No plan: the design above is small and fixed.

1. `approval` (`internal/tui/review.go`) gives `approved`, `approved under
   policy` or nothing from `update.Delegated`; the standing line
   (`detailMeta`) adds it after the candidate in any status, and the Review
   block uses it with unchanged wording.
2. `show --json` adds `approved_by` (`owner` or `policy`) while `approved`
   is set, from `update.Delegated`, as the sweep uses it; `list` has no
   `--json`. On this repository's main: G-260927-xd73p, the first real
   delegated verdict, gives `policy`; G-260927-dx0yn gives `owner`.
3. `TestStandingLineNamesWhoApproved` (card, owner and delegated, review
   and done); `TestApproveAndFeedbackCommands` (JSON: `owner`, absent after
   feedback, `policy`); `TestBoardReviewWorkflow` expects `approved` on the
   done card.
4. `docs/board.md` (Record detail), `docs/commands.md` (sweep),
   `docs/record-model.md` (`show --json`) and `grove --help` say so.
   Owner judgment in the terminal remains.

Verification at `727557a`'s code (unchanged in `80ef0f7`, which rewraps one
doc line): `go test -count=1 -timeout 120s ./...` all ok; `go vet ./...`,
`gofmt -l .` clean; `grove check` OK (224 records);
`internal/tui/testdata/terminal.py` all five ok. At `80ef0f7`: `go test
-short -count=1 . ./internal/cli` ok.

Review [G-260928-fav7h](G-260928-fav7h-review-of-g-260927-60ffq.md)
examined `80ef0f7`: `Open findings: none`.

## Next

Judge the candidate: open the board on a record with each verdict, then
`grove approve G-260927-60ffq VERDICT` in this worktree and `grove
integrate G-260927-60ffq` in main's checkout.

Verdict on candidate 3af8383, 2026-09-28: delegated under policy grove.yaml sha256:182036ce84beda7043a09222a7e22419798a46d27e4511f5e747a5032a5d7dcb: review G-260928-fav7h examined 80ef0f7 with no open finding; merged with main at 4f8f070, verification passed (go test -count=1 -timeout 120s ./...; go vet ./...; go run ./cmd/grove check); attempt G-260927-60ffq.20260928T164039Z produced it for 2.65 USD

Integrated under policy grove.yaml sha256:182036ce84beda7043a09222a7e22419798a46d27e4511f5e747a5032a5d7dcb as merge f800f02f29d55b56261ec0d7b4b5f321afd0be6f on main (was 4f8f07066ea4d45ee691ae262df959930af6ac5a); to reverse it: git revert -m 1 f800f02f29d55b56261ec0d7b4b5f321afd0be6f
