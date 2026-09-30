---
id: "G-260930-r2k4g"
type: work
title: "Make the complete Grove workflow usable without coaching"
status: proposed
created: "2026-09-30T01:05:46Z"
updated: "2026-09-30T03:18:04Z"
kind: feature
size: large
depends_on: ["G-260930-gwnb1", "G-260929-04svs", "G-260930-4742q"]
relates_to: ["G-260930-e8jj7", "G-260930-60c3d", "G-260930-62nmj", "G-260930-npw49", "G-260930-84fnb", "G-260930-2qa4a"]
---

## Outcome

A developer can adopt and use the complete supported Grove workflow without
learning its storage schema, branch selection rules or provider session
mechanics. On returning, they can explain what happened, what needs their
attention and what can happen next from the product itself.

Owner direction: [G-260930-e8jj7](G-260930-e8jj7-build-a-portable-workflo.md); milestone: [G-260930-60c3d](G-260930-60c3d-complete-the-portable-gr.md).
This work owns the assembled onboarding, CLI/TUI experience and adoption
instructions. Runtime guarantees belong to the component records it consumes.

## Scope and constraints

Implement the reviewed journey across setup, shaping, assignment, active
work, questions, continuation, result review and delivery. Give the primary
path sensible project defaults and explicit authority; advanced harness,
model/effort, checks and delivery settings remain inspectable and changeable.

Present work outcomes and decision needs before attempts, source revisions,
raw events or worktree paths. Derive presentation from shared facts and
operations; do not add a second editable work status in the UI. Surface
candidate freshness, verification, review, approval provenance and delivery
standing with links to their evidence. A summary explains its source and
uncertainty rather than making missing evidence look successful.

Configuration diagnostics explain what to fix and where. Preserve existing
instructions and show proposed configuration changes before applying them.
A setup walkthrough does not authorize implementation or merge. Reconnecting,
inspecting and unchanged waits do not launch work.

## Dependencies

- [G-260930-gwnb1](G-260930-gwnb1-prove-a-complete-workflo.md) supplies exercised continuation, recovery and evidence.
- [G-260929-04svs](G-260929-04svs-open-a-record-in-an-inte.md) supplies the
  transition from looking at work into helping it interactively.
- [G-260930-4742q](G-260930-4742q-deliver-accepted-work-th.md) supplies the additional host waits and delivery facts.
  The complete experience must incorporate those states rather than ship a
  local-only presentation under a portable-workflow claim.

## Acceptance

1. Starting with an ordinary existing repository and the adoption
   instructions, a new user can configure either supported harness, checks,
   limits and either supported delivery path, with actionable diagnostics
   for missing or unsupported capabilities.
2. The primary journey covers creating/shaping work, authorizing a bounded
   assignment, leaving and returning, answering, stopping/continuing,
   judging and delivering. The user need not select internal paths,
   lookup native session IDs or understand record frontmatter.
3. CLI and TUI consume the shared standing already delivered by the local
   contract and agree with agent context about current work, waits, candidate
   changes, evidence and actions. Accepted-but-undelivered, verified Done and
   unknown delivery/freshness are understandable. Raw source is inspectable
   without being misrepresented as the derived result; no completion PR or
   independent UI status store is added.
4. A result/catch-up view explains changed behavior, decisions, checks,
   review findings and what remains, including work that failed or stopped.
   A user can inspect the evidence without finding the original chat.
5. Terminal lifecycle and connected-workflow checks cover the primary path
   and recovery states; onboarding preserves existing project content.
   The owner judges the actual terminal experience against the design.
6. Published adoption instructions and supported-environment claims match
   what was exercised. The milestone's independent-user trial owns the
   final evidence that these pieces work together without author coaching.

## Next

Needs its declared prerequisites. Implement the reviewed experience and
documentation against their shared operations, then hand off a reproducible
adoption journey and remaining limitations to [G-260930-60c3d](G-260930-60c3d-complete-the-portable-gr.md).

The [experience sketches](G-260930-npw49-portable-workflow-experi.md)
and [shared contracts](G-260930-84fnb-portable-workflow-contra.md) are the
owner-accepted design baseline at f376a6d. Completion follows
[G-260930-2qa4a](G-260930-2qa4a-derive-done-from-recorde.md). The actual terminal
implementation still requires evaluation; exact arrangement and interface
spelling may be refined within the accepted journey.
