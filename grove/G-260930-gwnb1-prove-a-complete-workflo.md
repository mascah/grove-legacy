---
id: "G-260930-gwnb1"
type: work
title: "Prove a complete workflow across harnesses and local delivery"
status: proposed
created: "2026-09-30T01:05:45Z"
updated: "2026-09-30T03:18:03Z"
kind: feature
size: large
depends_on: ["G-260928-y2p5h", "G-260928-n4f1q", "G-260929-gm3m4"]
relates_to: ["G-260930-e8jj7", "G-260930-60c3d", "G-260930-62nmj", "G-260928-pqhyg", "G-260928-kehya", "G-260930-2qa4a"]
---

## Outcome

One bounded change traverses the complete local workflow across Claude Code
and Codex, including a real question, continuation, independent review and
delivery. The resulting evidence tells us whether the proposed architecture
supports the product promise and which existing code should remain.

Owner direction: [G-260930-e8jj7](G-260930-e8jj7-build-a-portable-workflo.md); milestone: [G-260930-60c3d](G-260930-60c3d-complete-the-portable-gr.md).

## Scope and constraints

This work owns the connective implementation for durable continuation and
the executable end-to-end trial. It consumes provider execution, separate
review and local delivery rather than implementing them again. Preserve
the live mandate, acceptance, answered questions, relevant decisions,
checkout/revisions, completed work, verification and next action in context
a fresh supported harness can use.

Proposed trial: prepare with one harness, stop at a supported checkpoint,
continue with the other, encounter and answer a consequential question,
obtain independent review and deliver locally. Exercise both handoff
directions across bounded scenarios. Also repeat continuation from committed
project context in a fresh clone without private transcripts. Inspect
repository state before relying on any checkpoint.

Mocks establish deterministic failure behavior; real bounded trials test
agent behavior. Execution configuration and any paid usage require the
owner's assignment. A fixture or an explicitly authorized project is the
write scope. Headless shaping does not execute these trials.

## Dependencies

- [G-260928-y2p5h](G-260928-y2p5h-run-an-attempt-on-codex.md) supplies both
  harnesses and the execution boundary.
- [G-260928-n4f1q](G-260928-n4f1q-review-a-candidate-as-a.md) supplies
  attributable independent review on the chosen harness.
- [G-260929-gm3m4](G-260929-gm3m4-clean-main-history-with.md) supplies local
  delivery and verified correspondence for a squashed candidate.

## Acceptance

1. The assigned scenario reaches a reviewed, accepted, locally delivered
   change with an auditable link from mandate to candidate to delivered
   result. Record all interventions, provider/model/effort, effective
   limits, reported usage and capability gaps.
2. Another supported harness continues at a checkpoint without reconstructing
   instructions from the original conversation. Both directions and a
   fresh-clone continuation retain the work's identity and authority.
3. Unanswered questions, unchanged waits, failed checks and interrupted or
   lost processes preserve completed work and identify the next permitted
   action. A retry neither duplicates a live owner nor silently redoes
   completed work. Changed repository state is reconciled explicitly.
4. Changed candidates and moved targets are exercised; stale approval or
   verification cannot authorize delivery. Independent review is never
   replaced by self-review without a visible unmet requirement. Agent trials
   start from both direct Markdown and tool-provided context: after acceptance
   and delivery they inspect the common standing, avoid assuming acceptance
   alone means Done, and do not restart already delivered work. Include a
   fresh clone and missing-evidence case; record confusion as a failure.
5. A linked report distinguishes deterministic checks, observed agent
   behavior and owner judgment. It recommends concrete retention,
   refactoring or replacement based on this journey. Any further work
   required for this outcome is delivered or remains an explicit blocker;
   a successful demonstration does not excuse an unreliable recovery path.

## Next

Needs its declared prerequisites. At assignment, choose the bounded scenarios
and resource mandate, then implement missing continuation behavior and run
the trial. Carry resulting contract changes into the owning records and
guides; the brief changes only for a new product-direction decision.
