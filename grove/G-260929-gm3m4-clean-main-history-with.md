---
id: "G-260929-gm3m4"
type: work
title: "Deliver accepted work locally with clean history and retained evidence"
status: proposed
created: "2026-09-29T16:15:25Z"
updated: "2026-09-30T01:16:09Z"
kind: feature
size: large
depends_on: ["G-260930-62nmj"]
relates_to: ["G-260930-e8jj7", "G-260930-60c3d", "G-260930-62nmj", "G-260930-4742q", "G-260930-r2k4g", "G-260921-3qgsf", "G-260921-jatts", "G-260921-btyck"]
---

## Outcome

An accepted change can be delivered locally with useful Conventional Commits
history and durable evidence connecting it to the reviewed candidate.
A project can use its preferred release tooling without exposing every
agent checkpoint and review repair as a separate product change.

Owner intent, 2026-09-29: "fix the commit problem more than anything so that
things are not noisy no matter what tool people use." The portable adoption
direction [G-260930-e8jj7](G-260930-e8jj7-build-a-portable-workflo.md) retains this record as the local delivery
component of [G-260930-60c3d](G-260930-60c3d-complete-the-portable-gr.md).

## Scope and constraints

Keep a complete local path without a remote, credentials, hosted CI or a
resident service. Squash delivery is in scope. A single work item normally
produces one meaningful product-change commit; a shared candidate is
delivered as its approved group without duplicating the code change.

Preserve exact review/approval attribution and verification of the result
against the actual target. A squash changes commit identity. The design
must establish correspondence between approved candidate, integration
inputs and resulting content, retain the evidence objects, and make
completion, dependency delivery and cleanup use that same fact.

A work ID, a matching patch or unverified trailers alone are insufficient
authority to close work. Changed candidates or targets require appropriate
reconsideration/reverification. Partial failure after integration must be
recoverable without a second delivery or destruction of evidence.

Choose how work metadata and completion evidence reach Git so that they do
not recreate checkpoint noise on the target. Message formatting and metadata
placement are proposed design choices to resolve in the journey/contracts
design; deriving a Conventional Commit type from the current kind field is
not assumed correct. Do not replace release tools or select version numbers.

Hosted PR delivery belongs to [G-260930-4742q](G-260930-4742q-deliver-accepted-work-th.md); this record owns the reusable
local delivery/correspondence behavior that it will consume. No host API,
push or protected-branch bypass is part of local integration.

## Observed evidence

The original 2026-09-29 shaping found 145 of 189 direct main commits since
2026-09-24 touched records only, and 7 of 14 fix commits in the keyborg pilot
were review-round fixes. Those are dated observations, not new measurements.

At main c8fa070ef9ff, internal/update/update.go checks candidate ancestry for
done, internal/deps checks delivery through ancestry, internal/integrate
performs a plain merge followed by separate done commits, and just
clean-merged selects branches by Git's merged relation. The earlier proposal
suggested Grove-Work/Grove-Candidate trailers and retained refs. That is a
candidate design, not sufficient proof of delivery and not a claim that
these are the only affected callers.

The settled [Integration](G-260921-3qgsf-integration.md) concept and the
implemented lifecycle must be reconciled when this work lands. The current
mechanism remains in force until then.

Other threads retained from the original conversation: release membership
was owner-selected, while a release/roadmap feature and release-please setup
remain separate future shaping; Grove does not own version numbers.
A catch-up digest is now in [G-260930-r2k4g](G-260930-r2k4g-make-the-complete-grove.md); hosted delivery is in
[G-260930-4742q](G-260930-4742q-deliver-accepted-work-th.md). The owner's correction that headless subagents work (the
observed denial was a sibling-worktree write) remains relevant to the
separate-process versus subagent comparison in the design.

## Dependencies

Depends on [G-260930-62nmj](G-260930-62nmj-design-the-portable-work.md): candidate/delivery correspondence, metadata
placement and the visible approval/delivery journey determine which
mechanisms and lifecycle contracts change.

## Acceptance

1. Local integration of one approved change and of a shared approved group
   produces the designed clean product history and traceable delivery
   evidence. Record-only bookkeeping does not reintroduce the noise this
   outcome removes; the owner judges representative main/changelog output.
2. Review, approval and verification are bound to the candidate and actual
   integration inputs. A moved target or altered candidate cannot use
   stale evidence; forged or mismatched delivery claims are rejected.
3. Done, dependency delivery, current-view presentation, policy integration
   and cleanup agree about a squashed result through the same validated
   correspondence. Relevant candidate/review objects survive cleanup and
   ordinary Git garbage collection under the documented retention design.
4. Refusal and partial-failure tests cover dirty targets, conflicts, shared
   groups, metadata failure after delivery and safe retry. Evidence survives
   and a retry does not integrate the same result twice.
5. An existing local project can migrate deliberately without losing work
   identity or historical approval/review evidence. Delivered behavior is
   documented in the model, guides, integration term and cleanup commands.
6. A bounded real local trial shows useful history, inspectable evidence and
   supported recovery without hosted infrastructure.

## Next

Needs [G-260930-62nmj](G-260930-62nmj-design-the-portable-work.md). Use its delivery and metadata contracts to plan the
affected lifecycle, policy, views and retention changes. This is shaped
local-delivery work; it no longer waits on an unspecified future shaping pass.
