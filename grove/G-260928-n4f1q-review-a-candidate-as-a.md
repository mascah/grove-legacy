---
id: "G-260928-n4f1q"
type: work
title: "Review candidates independently on the chosen harness"
status: proposed
created: "2026-09-28T19:29:00Z"
updated: "2026-09-30T20:22:08Z"
kind: feature
size: medium
depends_on: ["G-260928-y2p5h"]
relates_to: ["G-260928-917h8", "G-260924-5b6pz", "G-260921-rz7bn", "G-260925-5wrn8", "G-260924-3bapc", "G-260928-c5j9d", "G-260930-e8jj7", "G-260930-60c3d", "G-260930-gwnb1", "G-260930-62nmj", "G-260930-84fnb", "G-260930-2qa4a", "G-260930-tcc9w"]
---

## Outcome

Work can receive an attributable independent review on the harness, model
and effort chosen for review, regardless of the implementation harness.
The person judging the result sees what was examined, the findings and
verification, and any review that could not be completed.

Owner intent, 2026-09-28: "codex for the plan, claude code for the
implementation, and codex for the review." The portable adoption direction
[G-260930-e8jj7](G-260930-e8jj7-build-a-portable-workflo.md) makes this part of [G-260930-60c3d](G-260930-60c3d-complete-the-portable-gr.md).

## Scope and constraints

Provide a separately launched, bounded independent review through the common
execution boundary, available to the CLI, board and work workflow.
The working session's existing independent subagent route can remain where
it meets the same evidence and independence requirements. No universal
fresh-process rule is imposed on preparation and implementation.

A review mandate identifies the exact candidate, outcome/acceptance,
relevant context, previous findings and allowed verification. Isolation
protects the implementation result; checks that require writes need a
controlled disposable environment rather than mutation of the candidate.
Grove records work/examined/provenance and preserves failures without
inventing a clean review. Independent review informs approval; delegated
approval judgment is [G-260928-c5j9d](G-260928-c5j9d-judge-a-candidate-agains.md)'s
separate outcome inside the milestone under
[G-260930-tcc9w](G-260930-tcc9w-bound-delivery-groups-an.md). Review is per
delivery group: one member in separate mode, the shared candidate in together
mode. It does not choose the grouping or authorize the next delivery.

Consume the report and policy contract from
[G-261001-mbjwz](G-261001-mbjwz-separate-blocking-review.md), delivered through
the provider prerequisite's dependency chain. Preserve blockers, follow-ups,
completeness and their dispositions; do not design a second severity or
approval gate inside the review runner.

## Observed evidence

Observed 2026-09-28 at main `6fbb888`:

- Step 6 of the work guide dispatches `grove-reviewer` as a fresh subagent
  per gate, Claude only (`.claude/agents/grove-reviewer.md`: `model:
  inherit`, `effort: high`, read-only tools;
  [G-260924-5b6pz](G-260924-5b6pz-bound-an-attempt-at-its.md)). The review
  guide (`grove guide review`) is the brief; the reviewer returns findings
  and a closing line, and the implementing session writes the review
  record. Codex has no equivalent (`docs/commands.md`, Init).
- The implementing session's context peaks after the reviewer, 278k to
  352k tokens at handoff, and the fix rounds and record writing run there
  (G-260924-5b6pz's observations).
- The policy reads `Open findings: none` from a `current` review record
  that examined the candidate (record model).
- `attempt.json` records the reviewer definition's digest; a subagent
  review shows in `Metrics` as a subagent and in the result's per-model
  cost.
- A headless session must own a long command in the foreground with a
  timeout ([G-260924-3bapc](G-260924-3bapc-headless-attempts-run-lo.md)).

These observations describe the earlier implementation and provider versions,
not a lasting claim about which harness can create subagents.

## Dependencies

Depends on [G-260928-y2p5h](G-260928-y2p5h-run-an-attempt-on-codex.md):
this launches and observes review through its provider boundary, including
capabilities, limits and cancellation. A second review-specific runner would
duplicate those responsibilities.

## Acceptance

1. Both providers can perform an independent review with exact candidate
   attribution and a configured review role distinct from implementation.
   CLI, TUI and implementing sessions can invoke the same operation.
2. Findings, verification, provider/configuration and resource usage or its
   absence are recorded and inspectable. Failed, interrupted or malformed
   results remain visibly incomplete and cannot satisfy the review gate.
3. A candidate change prevents old findings from being presented as review
   of the new result. Concurrent changes are refused or isolated under
   the accepted design, without a second implementation owner.
4. The reviewer cannot modify the candidate; allowed checks and any
   disposable test environment are explicit. Self-review is never labelled
   independent. The person can follow a finding to its evidence.
5. Fake-provider tests cover those boundaries. Real bounded reviews on
   both providers of the same candidate are compared with the existing
   subagent route under an owner-specified mandate.
6. Review and work guides, configuration and command documentation describe
   the delivered operation. Review's settled meaning is unchanged.

## Next

Needs G-260928-y2p5h. At assignment, use the accepted role and evidence
contracts to choose the command/configuration details and isolation model.
The combined implementation-to-review trial belongs to [G-260930-gwnb1](G-260930-gwnb1-prove-a-complete-workflo.md).

The accepted [contracts design](G-260930-84fnb-portable-workflow-contra.md) specifies
the exact-candidate review and approval boundary, including controlled
artifact changes after the candidate. Review evidence feeds the recorded
acceptance and shared standing delivered by the local-delivery prerequisite;
a clean review cannot itself establish acceptance or Done. The owner accepted
the full written design at f376a6d, recorded in the design prerequisite.
