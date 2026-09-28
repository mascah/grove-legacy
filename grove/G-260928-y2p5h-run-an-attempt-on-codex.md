---
id: "G-260928-y2p5h"
type: work
title: "Run an attempt on Codex"
status: proposed
created: "2026-09-28T19:29:00Z"
updated: "2026-09-28T19:34:02Z"
kind: feature
size: large
relates_to: ["G-260928-vdhf0", "G-260928-917h8", "G-260923-tnn5e", "G-260924-59f5k", "G-260925-42j50", "G-260925-04ccr", "G-260925-3pj9a", "G-260925-p2k54", "G-260923-895zb", "G-260928-kehya"]
---

## Outcome

## Constraints

## Acceptance

## Next

## Outcome

The owner launches an attempt on Codex the way they launch one on Claude
Code, with `--provider codex` or a `provider:` default under `run:`, giving
the model, effort and cap in that harness's terms; Grove owns, watches,
stops and reconciles it as any attempt; the board and `attempt` show its
provider, model, cost where reported, activity and final report; and a
launch on Claude Code composes exactly what it composes today.

Decision [G-260928-vdhf0](G-260928-vdhf0-run-attempts-on-codex-as.md), the
owner on 2026-09-28. Term: [Provider](G-260928-917h8-provider.md).

## Constraints

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

Proposed design, labelled proposed:

- `--provider claude|codex` on `run` and `resolve`, `run: provider:` in
  `grove.yaml`, recorded in `attempt.json` and shown wherever the model is
  shown; the board's launch line takes it like any flag.
- For Codex: `codex exec --json -C WORKTREE -m MODEL -c
  model_reasoning_effort=EFFORT` with the permission mode mapped to a
  sandbox and approval setting the record model documents, and the prompt
  `$grove-work IDs --interaction headless`. Codex has no dollar budget, so
  `--budget` is refused for it and a `--max-minutes` cap is required in its
  place under the same no-default-spend rule; the eval runner's caps are
  the precedent. `run:` may default the cap per provider.
- A Codex events reader filling `Metrics`, the last message as the report,
  the model from its events, cost where reported and otherwise "not
  reported" with the reason.
- The entrypoint check and `differs_from_template` cover the Codex adapter.
- Stop and orphan handling through the same owner process; whether Codex
  ends its turn on SIGINT is verified, and the difference documented.

Out of scope: a reviewer for Codex attempts (G-260928-n4f1q); a third
provider; guide changes for Codex, which become their own record if the
trial shows the adapter failing.

## Acceptance

1. `grove run ID --provider codex --max-minutes N --permission-mode MODE`
   starts a Grove-owned Codex process in the worktree; `attempts`,
   `attempt` and the board show provider, model, cap, activity and the
   final report; `--dry-run` prints the composed command; a fake `codex`
   covers the lifecycle as the fake `claude` does; a launch without
   `--provider` composes byte for byte what it composes today.
2. Stop, orphan and interrupted states reconcile for a Codex attempt as for
   Claude; a process that ignores SIGINT is killed after the grace period
   and the result says so.
3. The permission-to-sandbox mapping and the cap are in the record model
   and `docs/commands.md`; `check` refuses an unknown provider and a budget
   for Codex.
4. One real `--until plan` attempt on Codex on real work, here or in
   keyborg, under a cap the owner names at assignment, reported with cost,
   duration, whether it stayed inside the worktree, and whether a Claude
   implementation attempt could follow from its plan.
5. The entrypoint revision check covers the Codex adapter.

## Next

Assign: `/grove-work G-260928-y2p5h`. This changes the command composition
in `internal/attempt/attempt.go`, which
[G-260928-kehya](G-260928-kehya-resume-an-attempt-s-sess.md) also touches;
no order is declared between them, and the owner should not launch them
alongside each other without reading both. G-260928-n4f1q names this record
in its `depends_on`.
