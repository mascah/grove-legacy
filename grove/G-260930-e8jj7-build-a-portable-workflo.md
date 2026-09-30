---
id: "G-260930-e8jj7"
type: decision
title: "Build a portable workflow product through selective redesign"
status: accepted
created: "2026-09-30T01:05:19Z"
updated: "2026-09-30T20:22:10Z"
relates_to: ["G-260930-60c3d", "G-260928-vdhf0", "G-260923-tnn5e", "G-260921-3qgsf", "G-260921-gtydy", "G-260930-tcc9w"]
---

## Decision

On 2026-09-29 (owner's local date), the owner selected a portable Grove
workflow product, developed alongside their job, and asked to reconcile the
brief and relevant proposals around the complete adoption journey.

The owner explicitly endorsed permission to redesign substantial parts of
Grove without first declaring the entire implementation disposable. Fresh
requirements and an understandable experience govern the work; the first
complete journey tests how much existing code deserves to remain.

The selected milestone is [G-260930-60c3d](G-260930-60c3d-complete-the-portable-gr.md). The brief remains the sole
current product-direction document; that work owns acceptance and progress.

## Consequences

- Requirements cover adopting, shaping, delegating, returning, continuing,
  judging and delivering a bounded change without the author coaching the
  user. Two harnesses, local delivery and one hosted PR path exercise
  portability. Independent use is part of acceptance.
- Existing proposals are reshaped around this journey. Package boundaries,
  command syntax, record-writing ceremony and UI layout may change where
  the requirements warrant it. Existing regression tests are evidence of
  behavior to assess, not a blanket requirement to retain every mechanism.
- Keep ordinary files and Git, a Go CLI/TUI, and no required resident
  service for core use. Preserve work identity, explicit authority,
  revision-bound review and approval, recoverable work, and attributable
  evidence. Migration must preserve existing projects' knowledge.
- Reopen ancestry-only delivery representation for squash and hosted
  integration. A rewritten result needs verified correspondence to the
  approved candidate; an ID, status or unverified trailer is insufficient.
  Current commands keep their documented contract until implementation
  delivers the replacement.
- [G-260928-vdhf0](G-260928-vdhf0-run-attempts-on-codex-as.md)'s two-provider
  direction stands. Substantial runner refactoring is allowed; preserving
  its present code shape or byte-for-byte command composition is not a
  product requirement. [G-260923-tnn5e](G-260923-tnn5e-run-attempts-as-a-grove.md)'s
  on-demand process ownership remains the starting constraint.
- The original 2026-09-29 decision deferred parallel launch and delegated
  approval judgment. On 2026-09-30 the owner amended that scope in
  [G-260930-tcc9w](G-260930-tcc9w-bound-delivery-groups-an.md): the LLM judge
  is included so unattended progression is tested in the complete workflow.
  Parallel launch remains deferred. A universal workflow builder, public
  plugin SDK, new git host, schedules and fleet management remain outside.

## Alternatives

Continuing the feature queue alone would leave the connective experience
unowned. A full rewrite could remove local coupling but would also require
re-establishing Git safety, lifecycle behavior and project continuity before
showing that the product is easier to use. Selective rebuilding keeps those
choices accountable to the same acceptance journey.

## Reconsideration

Reconsider the implementation strategy if the complete journey requires
duplicated authoritative state, pervasive incompatible assumptions, or a
more costly adaptation than a replacement with an explicit migration.
Reconsider milestone scope on evidence from the design and trials, recording
the owner's consequential choices rather than silently weakening acceptance.
