---
id: "G-260930-4mykk"
type: page
title: "Why G-260929-gm3m4 churned through review"
created: "2026-09-30T19:13:17Z"
updated: "2026-09-30T19:13:23Z"
relates_to: ["G-260929-gm3m4", "G-260930-c0g5p", "G-260930-gj9d7", "G-260930-2qa4a"]
---

A retrospective written at the owner's request on 2026-09-30, when
[G-260929-gm3m4](G-260929-gm3m4-clean-main-history-with.md) was merged
after its seventh independent review. It is written by the implementing
session and is its own account, for the owner to review and analyze; it
decides nothing.

## What happened

All on 2026-09-30, on branch `worktree-G-260929-gm3m4` from `556f362`,
between 13:09 and 18:34 UTC: 29 commits with the handoff, seven independent reviews and the
owner's own review, recorded in
[G-260930-c0g5p](G-260930-c0g5p-review-of-g-260929-gm3m4.md). The owner
counts more than ten fix rounds including earlier conversations; this page
covers what the branch and the review record show.

| Review | Examined | Open findings | What they were about |
| --- | --- | --- | --- |
| 1 | `af6c8df` | 10 | kept and re-accepted branches, feedback after a squash, an altered claim delivered twice, stale terms, small correctness items |
| 2 | `e17ddd4` | 5 | the continuation merge base that round 1's fix added |
| 3 | `0e5e8cd` | 4 | the continuation base under a criss-cross merge, its dependence on a local ref, a ref left by a refused delivery, a missing test |
| 4 | `9c9fed2` | 3 | a date-ordered history walk misreading delivery, the continuation's cost on every read, cleanup deleting an unretained branch |
| Owner | — | redesign | "I need grove to NOT be a complete performance hog": re-proving deliveries on every read is dropped ([G-260930-gj9d7](G-260930-gj9d7-prove-delivery-once-at-i.md)) |
| 5 | `2cb7bdd` | 5 | accepted work carried inside another delivery, revert wording, `context` reading the target twice, a stale model sentence, the audit reading the checkout's records |
| 6 | `aaf2d5a` | 2 | round 5's carry refusal deadlocking, and wrong advice for carried work in review |
| 7 | `281a8f1` | 2 | carried work delivered elsewhere for another candidate; a superseded acceptance chosen over a later one |

Of 31 findings, about 20 concern what a branch, a kept branch or a carried
candidate means once delivery is a squash; about 5 are cost; about 5 are
documents disagreeing with each other or the code; the rest are tests and
small correctness.

## Why it churned

### 1. Squash delivery removed the fact everything else was built on

Before schema 4, "the target contains the candidate" (Git ancestry) answered
many questions at once for free: is the work delivered, which branches are
merged, what a later branch carries, whether a branch can be deleted, what
a revert undoes. A squash commit is not an ancestor-linked copy of the
branch, so each of those questions had to be answered again separately,
from trailers, the evidence ref or the target's record. Every place that
had leaned on ancestry became a seam, and most findings sit in those seams:
kept branches, re-delivery, carried work, cleanup, reverts, criss-cross
merges. This is inherent in the owner's choice of squash delivery, not a
defect of one fix; it is where the cost of that choice shows.

### 2. The first design made every reading prove history

The first implementation re-proved every past delivery on every read, from
history walks, merge bases and a continuation base for kept branches. Each
of those mechanisms had its own edge cases (criss-cross merges, commit-date
ordering, dependence on a local ref) and a cost that grew with history. The
owner's cost constraint ("I didn't authorize adding more git work per
load") was real from the start but was not stated as an acceptance
criterion or measured before building, so it surfaced only at the owner's
review after four rounds, and the redesign threw most of those rounds'
fixes away.

### 3. Fixes added mechanisms, and the next round reviewed the mechanism

Twice a fix introduced a new mechanism that became the next round's main
subject:

- Round 1's kept-branch finding was fixed with a continuation merge base
  (`e17ddd4`). Round 2 found it re-applied target changes; round 3 found it
  wrong under a criss-cross merge and dependent on a local ref; round 4
  found its cost. The redesign removed it.
- Round 5's carried-work finding was fixed with a refusal naming what to
  integrate first (`f1fa1b7`). Round 6 found it deadlocked; the fix
  (`a4453d1`) patched the two reported scenarios; round 7 found their
  neighbours. `e144aa6` finally replaced the case-by-case refusals with one
  rule: the target's own copy of the carried work decides.

The implementing session fixed the reproduced scenario, tested that
scenario, and dispatched the next review, rather than first stating the
rule the scenario was an instance of and testing the states around it. An
adversarial reviewer explores exactly those neighbouring states, so each
narrow fix produced the next round's findings.

### 4. Refusals that prescribe recoveries multiply the claims to verify

`integrate` names a next action for every refusal. Each piece of advice is
a claim over a state space that grows combinatorially: the record's status
on the branch, the target's copy, how many branches accept the work, and
how their candidates relate. Reviewers verified advice by following it,
which is how the circular advice was found. Precise advice is valuable, but
every prescribed recovery is another behaviour that must be right.

### 5. One work item carried too much

The work covered schema 4, migration, squash delivery, derived standing in
about ten consumers, the audit, the board, and the reconciliation of the
model, the design, the command reference, the work guide, five terms, the
contract, the brief, the README and CLAUDE.md. Every review examined the
whole combined change against the whole acceptance, so a fix round in one
corner reopened the review of all of it, and findings could come from
anywhere.

### 6. The pass bar was zero findings of any weight

Each review had to end `Open findings: none`, knowledge and low findings
included, and the implementing session's review prompts listed edge cases
to "check hard". On a change this size, an adversarial reviewer rarely
reaches zero; the round cap, not the findings, became the effective stop.
Several late findings were about refusal wording in rare multi-branch
cases, which the reviewer itself confirmed could not deliver unaccepted
work.

### 7. The same rule lives in many documents

A change to what Done means had to be restated consistently in the model,
the design, the command reference, the work guide, the terms, the contract,
the brief, the README and CLAUDE.md, and the two shipped documents have
size budgets. Drift between them was itself a finding in several rounds,
and the budgets forced trims late in the full-suite run.

### 8. The work spanned several conversations

Context was compacted and carried across sessions by summary. The
implementing session worked from summaries of earlier rounds and verified
each fix against the named scenario, without its own adversarial pass
before dispatching the next reviewer.

## What might have prevented it

Options for the owner to weigh, not decisions:

- **Settle cost and the reading rule before building.** State the Git work a
  read may do as acceptance, and measure the design against it before
  implementation, as [G-260930-gj9d7](G-260930-gj9d7-prove-delivery-once-at-i.md)
  now does.
- **Split large work.** Schema and migration, the squash delivery, derived
  standing, the audit, and document reconciliation could have been separate
  work items, each reviewed on its own diff.
- **Rule before fix.** For each finding, write the rule it violates in the
  document that owns it, derive the code from the rule, and test the states
  around it (a table over the branch's copy, the target's copy and the
  candidates), before asking for review.
- **A severity bar for passing.** Pass a round with no open correctness
  finding of medium or higher; record low and wording findings in the
  work's Next; scope later rounds to the fix and the rule it claims.
- **Fewer prescriptive recoveries.** State the fact that refuses the
  delivery and at most one next action.
- **Reduce the states squash creates.** Most carried-work cases arise from
  basing work on another unmerged work branch, and from branches kept after
  delivery. Requiring work to branch from the target, or cleanup by default
  after a squash, would remove whole classes of cases; each is a trade the
  owner would have to choose.

## State at merge

- The carried-work rule in `e144aa6` was merged without an independent
  review, by the owner's choice after review 7. Its tests cover every
  scenario from reviews 5 to 7 and follow each recovery through to a
  delivery; the full suite passed at `e144aa6`.
- The audit's history walk (`grove check --deliveries`) still orders by
  commit date and can report a delivery not proved under clock skew, never
  proved; a `ponytail:` comment in `internal/standing` names the exact
  alternative.
