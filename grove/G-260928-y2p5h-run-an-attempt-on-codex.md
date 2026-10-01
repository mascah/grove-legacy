---
id: "G-260928-y2p5h"
type: work
title: "Run bounded work through portable Claude and Codex execution"
status: proposed
created: "2026-09-28T19:29:00Z"
updated: "2026-10-01T02:22:09Z"
kind: feature
size: large
relates_to: ["G-260928-vdhf0", "G-260928-917h8", "G-260923-tnn5e", "G-260924-59f5k", "G-260925-42j50", "G-260925-04ccr", "G-260925-3pj9a", "G-260925-p2k54", "G-260923-895zb", "G-260928-kehya", "G-260930-e8jj7", "G-260930-60c3d", "G-260930-62nmj", "G-260930-gwnb1", "G-260930-84fnb", "G-260930-2qa4a", "G-260930-tcc9w"]
depends_on: ["G-260930-62nmj", "G-260929-gm3m4", "G-260930-0s29t", "G-261001-62n4p"]
---

## Outcome

A project's bounded work can run on Claude Code or Codex through a common
Grove execution boundary, with its actual harness capabilities and
configuration visible. Changing harness does not require changing the
project's work model or learning another workflow.

The owner selected two providers on 2026-09-28 in
[G-260928-vdhf0](G-260928-vdhf0-run-attempts-on-codex-as.md), and reframed
this work around portable adoption on 2026-09-29 in [G-260930-e8jj7](G-260930-e8jj7-build-a-portable-workflo.md).
This record supplies execution for [G-260930-60c3d](G-260930-60c3d-complete-the-portable-gr.md); it no longer requires
preserving the current command bytes or package shape.

## Scope and constraints

Separate provider-specific command construction, event decoding, capabilities,
permissions, limits and native-session behavior. Consume the common workspace
eligibility and ownership boundary delivered by
[G-261001-62n4p](G-261001-62n4p-isolate-workspace-eligib.md); do not repeat its
lifecycle redesign inside provider extraction. Preserve duplicate-start
refusal, stop and result reconciliation through that boundary.
Extract or replace code where the design warrants it. Exercise the boundary
with both providers; a public plugin SDK and a third provider are out of scope.

Keep the shared Grove work guide and thin adapters. Necessary guide changes
for the two-provider contract are in scope. A provider's successful exit
does not claim reviewed or accepted work. Separate review belongs to
[G-260928-n4f1q](G-260928-n4f1q-review-a-candidate-as-a.md); complete durable
handoff and the combined trial belong to [G-260930-gwnb1](G-260930-gwnb1-prove-a-complete-workflo.md).
Selection progression across deliveries belongs to
[G-260930-yfh91](G-260930-yfh91-continue-an-authorized-s.md); do not build a
second orchestrator inside either provider. One provider execution serves
the current delivery group, with common operations owning the next launch.

Provider capabilities are explicit. Unsupported and unreported are distinct
from zero or success. A time limit is not a dollar cap, and permission modes
must not be equated by name alone. Preserve the no-default-spend mandate and
[G-260925-04ccr](G-260925-04ccr-never-run-gpt-6-astra-un.md)'s model constraint.
Keep usable existing project configuration through an explicit migration
where contracts change.

## Observed evidence

Observed 2026-09-28 at main `6fbb888`:

- The runner composes `claude -p PROMPT --output-format stream-json --verbose
  --session-id … --max-budget-usd … --permission-mode … --permission-prompts
  none [--model] [--effort]` (`internal/attempt/attempt.go`; the
  `attempt.json` of `G-260928-4qv1m.20260928T170244Z`); `GROVE_CLAUDE` only
  swaps the executable. The activity reader (`internal/attempt/activity.go`)
  fills provider-neutral `Metrics` from Claude's stream, built so that
  "another harness's reader fills the same fields"
  ([G-260923-895zb](G-260923-895zb-make-attempts-easy-to-sc.md)). Liveness by
  file lock, stop by SIGINT then SIGKILL of the process group, `result.json`
  at the end (`docs/commands.md`, Attempts).
- Codex 0.157.0 `exec`: `--json` events as JSONL, `-m MODEL`,
  `-c model_reasoning_effort=…`, `-s read-only|workspace-write|
  danger-full-access`, `--dangerously-bypass-approvals-and-sandbox`,
  `-C DIR`, `--output-schema FILE`, `-o` last message, `resume`; no budget
  flag; `CODEX_HOME` is its configuration directory. The eval runner
  (`evals/run.py`, [G-260924-59f5k](G-260924-59f5k-run-the-eval-pair-on-cod.md))
  composes `codex exec --json` with a time cap and a plan-usage cap, reads
  its events for model and turns, and has a fake `codex` for its selftest.
- [G-260925-42j50](G-260925-42j50-codex-eval-row-pattern-p.md): on the eval
  fixture Codex read outside the clone in every run and the owner's Grove
  checkout in five; the guide's pattern matched Claude's. The Codex
  adapter is `.agents/skills/grove-work/SKILL.md` and its interactive
  invocation `$grove-work`; Codex has no agent definitions, so a Codex
  session reports its missing independent reviewer (`docs/commands.md`,
  Init) and work whose record requires a review stays active until
  [G-260928-n4f1q](G-260928-n4f1q-review-a-candidate-as-a.md).
- `run` refuses a worktree without `.claude/skills/grove-work/SKILL.md` and
  records the Claude files' digests and entrypoint revisions
  ([G-260925-3pj9a](G-260925-3pj9a-launch-attempts-only-whe.md),
  [G-260925-p2k54](G-260925-p2k54-keep-installed-harness-e.md)); a Codex
  attempt needs the `.agents` file checked instead.
- [G-260925-04ccr](G-260925-04ccr-never-run-gpt-6-astra-un.md) binds any
  default model: never gpt-6-astra unless the owner names it.

These are dated observations, including the earlier statement about Codex
agent definitions; recheck supported versions and capabilities at
implementation. They are not portable product requirements.

## Dependencies

Depends on [G-260930-62nmj](G-260930-62nmj-design-the-portable-work.md): the accepted assignment/capability/continuation
contracts determine this boundary. Coding directly to the old CLI shape
would pre-empt that design and risk redoing the runner.

Also depends on [G-260929-gm3m4](G-260929-gm3m4-clean-main-history-with.md):
it replaces the work record contract and the shared standing used by
internal/attempt selection, facts and resolution, plus context and policy.
That migration is delivered; consume its cheap standing API and keep
acceptance/delivery out of provider events.

Also depends on [G-260930-0s29t](G-260930-0s29t-use-one-fresh-workspace.md):
it changes workspace admission and lifetime before this extraction touches
the same launch/selection code. Repeated execution on an old squashed branch
is not behavior to preserve; before-delivery continuation remains required.

Also depends on [G-261001-62n4p](G-261001-62n4p-isolate-workspace-eligib.md):
the owner authorized a narrower refactor and imported-branch demonstration
before introducing another provider. It separates delivery's workspace
checks from provider execution and inventories the old rules removed.

## Acceptance

1. Both providers execute a bounded delivery group in a Grove-owned worktree.
   Dry run and inspection explain provider, model/effort, effective
   permissions and limits, entrypoint compatibility and resume behavior.
   CLI and TUI read the same execution facts.
2. Both providers have deterministic coverage for launch, duplicate refusal,
   successful/failed exit, missing executable or incompatible entrypoint,
   stop escalation, lost owner, partial output and reconnect. Unsupported
   combinations fail before spending or changing the assignment.
3. Provider-specific events preserve raw attribution and supply common
   activity/result facts without leaking Claude-only assumptions into the
   work model. Missing usage or cost is reported honestly.
4. Existing Claude workflows retain their required behavior, with migration
   for deliberate interface changes. No byte-for-byte command-preservation
   constraint prevents a cleaner implementation. Both providers consume the
   same workspace eligibility/ownership interface; neither parses the other's
   events or introduces its own branch-retirement or cleanup rules.
5. Real bounded smoke trials on both providers record versions, permission
   behavior, effective limits, stopping and observed activity. The owner
   supplies execution configuration and any paid usage mandate.
6. Implemented contracts are reconciled in the owning guides, configuration
   model and command documentation. The complete cross-provider journey is
   tested by [G-260930-gwnb1](G-260930-gwnb1-prove-a-complete-workflo.md), not inferred from these smoke trials.

## Next

Needs the declared design, delivered local integration, local workspace and
shared workspace-boundary prerequisites. Plan extraction against the reconciled
[contracts](G-260930-84fnb-portable-workflow-contra.md) and
[G-260930-tcc9w](G-260930-tcc9w-bound-delivery-groups-an.md). Preserve
ownership, limits, source safety and pre-delivery continuation; do not retain
obsolete kept-branch execution merely because a test covers it. The workspace
prerequisite's local walkthrough must establish its boundary before extraction.
Recheck dated provider observations before relying on flags or capabilities.
This remains proposed and unassigned.
