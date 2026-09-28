---
id: "G-260928-63124"
type: work
title: "Show each member's state of a selection on the board"
status: review
created: "2026-09-28T19:28:59Z"
updated: "2026-09-28T21:03:00Z"
kind: feature
size: small
relates_to: ["G-260925-7c8g9", "G-260925-wc2pz", "G-260923-895zb", "G-260928-dbgbw"]
candidate: "393eb56e7f21a91582a57358615d8d90a5fb029f"
---

## Outcome

When a selection of several work items runs or ends, the owner sees on each
member's card and in its detail what that member's state is and what the
next launch would do, without opening the attempt.

Owner intent, shaping conversation 2026-09-28, from keyborg: "It's hard to
understand what is happening across different items, or when one has to
stop for a question"; it was "unclear what would happen" on relaunch, so the
owner ran the three IDs interactively instead.

## Constraints

Observed 2026-09-28:

- keyborg attempt `G-260927-5zhf5.20260928T164846Z` (read with keyborg's
  own `grove attempt`): three members, exit 0, $24.14, 65 minutes. Two
  members ended `active; its Next holds the checkpoint`, one `active,
  waiting on question G-260928-tyrg0`. Nothing entered review, because a
  started unit left incomplete holds every unit on the branch out of review
  (work guide step 8; [G-260925-wc2pz](G-260925-wc2pz-review-an-explicitly-sel.md)).
  The provider's final report said so, and that relaunching resumes without
  redoing the two; the owner instead ran `/grove-work` for the three
  interactively, which handed all three off on one candidate.
- On the board the attempt row is listed by its first ID and a count
  (`G-260927-5zhf5+2`); a member's card is tagged `● running` while it runs
  and nothing after; the per-member states are on the attempt screen only
  (`docs/board.md`, Attempts; `standingOf` in `internal/tui/attempts.go`).
- The outcomes `waiting on question`, `answered since`, `plan ready` and
  `candidate ready` are derived per attempt, not per member
  ([G-260923-895zb](G-260923-895zb-make-attempts-easy-to-sc.md),
  [G-260925-7c8g9](G-260925-7c8g9-execute-an-explicitly-se.md)).

Proposed design, labelled proposed:

- After a selection's attempt ends, each member's card carries its own
  state where it needs the owner: `waiting on G-…` for the blocked member,
  `held by G-…` for a complete member the group holds out of review, `plan
  ready`, `candidate ready`, from the rules the attempt screen already
  applies; the tag names the selection where the card's ID is not its first.
- The work's detail row and the attempt's Next say what `R` does: resumes on
  the branch, does not redo members whose checkpoint the branch confirms,
  and which member it starts.
- Nothing changes in what a selection is or how it hands off
  (G-260925-wc2pz). Parallel launches are
  [G-260928-369c1](G-260928-369c1-launch-independent-membe.md).

## Acceptance

1. With the fake provider, a two-member selection stopped on one member's
   question shows both cards' states as above and the detail's line; after
   the answer, the blocked member shows `question answered: R again` and the
   held member says what `R` does.
2. The attempts list names every member's state in its row where the width
   allows, cut last.
3. `docs/board.md` says so; the owner judges on a real selection.

## Evidence

Branch `worktree-G-260928-csg91-G-260928-r1hkh-G-260928-y50a4-G-260928-63124`,
base main `453add4`; candidate is the commit adding this Evidence, shared by
the selection G-260928-csg91, G-260928-r1hkh, G-260928-y50a4, G-260928-63124.

1. `TestSelectionMembersOnTheBoard`: `held by W-002` and
   `waiting on Q-002 · in W-001+1`, the detail's line, R refused while
   it waits; after the answer `question answered: R again`, the held
   member says `R resumes …`, and R relaunches both IDs on the branch in
   its worktree, checking each member's file on Enter. Members that
   moved on carry no state; `TestSelectionNextOnlyWhereRResumes`. Limit:
   the test renders hand-built results, not a fake-provider run;
   `internal/attempt`'s `TestMemberResults` produces such members.
2. The attempts row names each member's state (`moved on` when it did),
   cut last; asserted at 160 columns, and the list fit test holds.
3. `docs/board.md` Attempts; decision
   [G-260928-dbgbw](G-260928-dbgbw-r-on-a-member-of-an-ende.md). The owner's
   judgment on a real selection is pending.

Review: [G-260928-4q08b](G-260928-4q08b-review-of-g-260928-r1hkh.md),
three rounds, examined `9deee95`: "Open findings: 3", left open at the
gate's cap: R promises a resume that Start refuses when a member was
abandoned on the branch (nothing launches); `docs/board.md` lacks
`moved on`; and the fake-provider limit above.

Verification at `9deee95`: `go vet ./...` and `gofmt -l .` clean;
`grove check` OK (247 records); `go test -count=1 -timeout 120s ./...` all
pass; `terminal.py` all 13 scenarios pass. Under an artificial five-package
load, `internal/attempt`'s `TestOwnerLost` failed once (a kill race in a
package this branch does not touch) and passed on rerun.

## Next

In review with the selection's shared candidate. The owner judges each in
a terminal (`go run ./cmd/grove` in the worktree), then, in this worktree,
`grove approve ID VERDICT` for each of G-260928-csg91, G-260928-r1hkh,
G-260928-y50a4 and G-260928-63124, and in main's checkout
`grove integrate G-260928-csg91`, which merges the group.
