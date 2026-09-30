---
id: "G-260928-y2p5h"
type: work
title: "Run bounded work through portable Claude and Codex execution"
status: proposed
created: "2026-09-28T19:29:00Z"
updated: "2026-09-30T03:17:59Z"
kind: feature
size: large
relates_to: ["G-260928-vdhf0", "G-260928-917h8", "G-260923-tnn5e", "G-260924-59f5k", "G-260925-42j50", "G-260925-04ccr", "G-260925-3pj9a", "G-260925-p2k54", "G-260923-895zb", "G-260928-kehya", "G-260930-e8jj7", "G-260930-60c3d", "G-260930-62nmj", "G-260930-gwnb1", "G-260930-84fnb", "G-260930-2qa4a"]
depends_on: ["G-260930-62nmj", "G-260929-gm3m4"]
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
permissions, limits and native-session behavior from common ownership,
worktree handling, duplicate-start refusal, stop and result reconciliation.
Extract or replace code where the design warrants it. Exercise the boundary
with both providers; a public plugin SDK and a third provider are out of scope.

Keep the shared Grove work guide and thin adapters. Necessary guide changes
for the two-provider contract are in scope. A provider's successful exit
does not claim reviewed or accepted work. Separate review belongs to
[G-260928-n4f1q](G-260928-n4f1q-review-a-candidate-as-a.md); complete durable
handoff and the combined trial belong to [G-260930-gwnb1](G-260930-gwnb1-prove-a-complete-workflo.md).

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
Those callers currently branch on raw status. Refactoring the provider
boundary against the old lifecycle would duplicate that migration. Consume
the delivered standing API and keep acceptance/delivery out of provider events.

## Acceptance

1. Both providers execute a bounded selection in a Grove-owned worktree.
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
   constraint prevents a cleaner implementation.
5. Real bounded smoke trials on both providers record versions, permission
   behavior, effective limits, stopping and observed activity. The owner
   supplies execution configuration and any paid usage mandate.
6. Implemented contracts are reconciled in the owning guides, configuration
   model and command documentation. The complete cross-provider journey is
   tested by [G-260930-gwnb1](G-260930-gwnb1-prove-a-complete-workflo.md), not inferred from these smoke trials.

## Next

Needs the design and local-delivery prerequisites. At assignment, plan
provider extraction/replacement against their delivered standing and accepted
assignment/capability contracts, retaining meaningful lifecycle regressions.
Recheck the dated provider observations before relying on flags or features.

The [contracts design](G-260930-84fnb-portable-workflow-contra.md) specifies
the accepted boundaries. Completion is selected by
[G-260930-2qa4a](G-260930-2qa4a-derive-done-from-recorde.md); the owner accepted
the full design at f376a6d, recorded in the design prerequisite. No new
schema or provider API is implemented yet.
