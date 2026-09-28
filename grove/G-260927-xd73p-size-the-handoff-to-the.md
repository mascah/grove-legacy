---
id: "G-260927-xd73p"
type: work
title: "Size the handoff to the record"
status: done
created: "2026-09-27T22:13:36Z"
updated: "2026-09-28T16:39:37Z"
size: small
depends_on: ["G-260927-cg6rt"]
relates_to: ["G-260921-9wkjt", "G-260925-wh9ax", "G-260921-btyck", "G-260927-ngkbz"]
candidate: "dc8149c9b5d2dfc7da8c3e93f19043896a2f49a9"
approved: "dc8149c9b5d2dfc7da8c3e93f19043896a2f49a9"
---

## Outcome

A small piece of work hands off with the facts a judge needs and no more, so
the fixed cost of Grove's handoff is proportional to the work, while the
candidate discipline and the review a record requires do not shrink.

Owner intent, conversation 2026-09-27: "we should definitely right size the
handoff." Which records get the compact form, and what it holds, is
proposed below.

## Constraints

Observed 2026-09-27 at main `fef102f`:

- Step 8 of the work guide (3,860 bytes) asks one handoff shape of every
  record: branch, base and candidate; starting revisions; behavior against
  each acceptance item; decisions and why; verification commands, results
  and commit; each review record with findings and dispositions; issues and
  limits; the integrator's commands. Step 6 requires an independent review
  at each boundary and allows a self-check only for "a small
  documentation-only change".
- `size` is an optional work field, `small`, `medium` or `large`
  (`docs/record-model.md:165,393`), shown on the board card
  (`internal/tui/detail.go:445`) and mentioned nowhere in the work guide.
- The smallest recent attempt, G-260927-ngkbz (`size: small`, a README
  framing change), ran 46 turns, $2.99 and 10 minutes headless and produced
  a review record ([G-260927-917tb](G-260927-917tb-review-of-g-260927-ngkbz.md))
  and a full Evidence and Next.
- In this repository since 2026-09-19, 357 of 672 commits are `docs(G-…)`
  record edits, 107 of them status flips; in 48 interactive work sessions
  agents edited records 67 times and other files 46 times. This
  repository's work is unusually documentation-heavy, so these are an upper
  bound, not a rate.
- The standing policy in `grove.yaml` already judges by diff size
  (`max_lines: 300`).
- The record model says "The record's Evidence and Next carry the handoff
  the work guide describes", and the [Approval](G-260921-btyck-approval.md)
  term binds approval to one candidate; neither changes.
- **Depends on [G-260927-cg6rt](G-260927-cg6rt-print-the-work-guide-in.md):**
  both rewrite step 8, and the staging moves it into an on-demand part, so
  editing it in both at once would redo or conflict with that work.

**Proposed design.** Step 8 names two shapes. Full: as today. Compact, when
the record's `size` is `small` or its own Next says so: branch, base and
candidate; one line of evidence per acceptance item; the verification
commands, their results and the commit they ran at; the review's closing
line and the disposition of each finding written into the record's
Evidence, with a separate review record only where the record or plan asks
for one; the integrator's two commands. What never shrinks: the status
change committed alone with the candidate, the independent review where
step 6 requires it, and "an implementation session never writes done". The
review guide's return (`grove guide review`) is unchanged; step 6 says where
the closing line goes when no review record is written. A record without
`size` gets the full shape.

## Acceptance

1. The work guide states both shapes, when each applies and what never
   shrinks; the record model's Review paragraph stays true.
2. A `size: small` record carried through `grove run` after the change
   hands off in the compact shape: Evidence and Next together under 40
   lines, no review record unless required, and the owner judges the
   candidate from the record alone (owner judgment, recorded with the
   attempt's turns, cost and duration next to G-260927-ngkbz's 46 turns,
   $2.99 and 10 minutes).
3. The shipped-document checks pass; `docs/work-review.md` and
   `docs/commands.md` need no change, or are reconciled.

## Evidence

Compact handoff (`size: small`). Branch `worktree-G-260927-xd73p`, base
`bb39668`; implementation `1539e3e` and `a75fa09`; the candidate is the
commit adding this Evidence. No plan needed: the design above is the plan.

1. Step 8 of the work guide names a full and a compact shape: compact for
   `size: small` or a Next that asks, full otherwise; the step 6 review, the
   candidate discipline and "never writes done" never shrink. The record
   model's Review paragraph still holds; its `size` lines and "Not enforced
   by software" list were reconciled.
2. Not judgeable in this attempt: it needs a later `grove run` of a
   `size: small` record on a target holding this change.
3. `grove check`, `go test ./` and `./internal/cli` pass; the added lines
   hold no link, record ID or path. `docs/work-review.md` and
   `docs/commands.md` need no change.

Decision beyond the proposed design: a compact handoff still writes a review
record where `policy:` delegates approval, since `grove sweep` reads the
closing line only from one; this repository's policy does, hence
[G-260928-b4aat](G-260928-b4aat-review-of-g-260927-xd73p.md). Compact also
keeps a line per limit and the examined commit (review round 1).

Verification at `a75fa09`: `gofmt -l .` empty, `go vet ./...` ok,
`go run ./cmd/grove check` OK 219 records, `go test -count=1 -timeout 120s
./...` all ok. Review: two rounds, 3 findings fixed, `Open findings: none`.

## Next

In review. Judge from this checkout, then integrate from the target's:
`go run ./cmd/grove approve G-260927-xd73p "VERDICT"`, then
`go run ./cmd/grove integrate G-260927-xd73p`. After integration, judge
acceptance 2 on the next `size: small` attempt against G-260927-ngkbz's 46
turns, $2.99 and 10 minutes.

Verdict on candidate dc8149c, 2026-09-28: delegated under policy grove.yaml sha256:182036ce84beda7043a09222a7e22419798a46d27e4511f5e747a5032a5d7dcb: review G-260928-b4aat examined a75fa09 with no open finding; merged with main at 7459594, verification passed (go test -count=1 -timeout 120s ./...; go vet ./...; go run ./cmd/grove check); attempt G-260927-xd73p.20260928T162049Z produced it for 2.13 USD

Integrated under policy grove.yaml sha256:182036ce84beda7043a09222a7e22419798a46d27e4511f5e747a5032a5d7dcb as merge 6a657a034ae2160a37fa80097497d8f3fece27ce on main (was 74595944efb4b9a3381ea9763c49767666cb7f93); to reverse it: git revert -m 1 6a657a034ae2160a37fa80097497d8f3fece27ce
