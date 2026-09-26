---
id: "G-260925-gymkr"
type: question
title: "What mandate should the G-260925-pbx81 without-row runs have?"
status: resolved
created: "2026-09-25T19:36:20Z"
updated: "2026-09-25T20:14:11Z"
blocks: ["G-260925-pbx81"]
relates_to: ["G-260923-659zw", "G-260924-frzeg", "G-260925-04ccr"]
---

## Question

[G-260925-pbx81](G-260925-pbx81-evaluate-whether-agents.md) acceptance 2 runs its two cases
under "an owner mandate naming a Claude model, runs, budget and permission
mode", and its Next says the owner supplies it before any paid run. The
headless `/grove-work G-260925-pbx81` session of 2026-09-25 was given none, so it
built the cases (`b2c7ef0` on `worktree-G-260925-pbx81`, `selftest: ok`) and stops
here. Under [G-260925-04ccr](G-260925-04ccr-never-run-gpt-6-astra-un.md), an item left
blank does not take the recommendation: it leaves this question open and the
runs unstarted. Answer each:

1. **Model** (`--model`, a Claude model by full ID). Recommendation:
   `claude-opus-5-5`, as [G-260924-frzeg](G-260924-frzeg-baseline-runs-the-missin.md)
   ran, so the rows compare with the pair's baseline.
2. **Runs per case** (`--runs`). Recommendation: 5, as G-260924-frzeg; the `with`
   row after [G-260925-dzxm6](G-260925-dzxm6-search-record-bodies-and.md) uses the same number.
3. **Budget per run** (`--budget`, USD). Recommendation: $5, as G-260923-659zw set
   for G-260924-frzeg, whose runs cost $0.27 to $0.32. These runs read more records,
   so expect more, not above $1. The cap printed first is runs × 2 × budget:
   $50 at the recommendation.
4. **Permission mode** (`--permission-mode`). Recommendation: `auto`, as
   G-260923-659zw chose and G-260924-frzeg ran with no denials.
5. **Config directory** (`--config-dir`). Recommendation:
   `~/.cache/grove-evals/claude`, the G-260924-frzeg login, which the runner still
   has to accept as clean at run time.

The command the answer authorizes, from `worktree-G-260925-pbx81` as it stands (the
`without` row, guides digest before G-260925-dzxm6):

```sh
python3 evals/run.py run --case listed-constraint --case code-constraint \
    --runs RUNS --budget USD --model MODEL --permission-mode MODE \
    --config-dir DIR --out ~/.cache/grove-evals/runs/2026-09-25-G-260925-pbx81-without
```

It never passes `--harness codex` and never names `gpt-6-astra` (G-260925-pbx81).

Who can answer: the owner. Answer in this record, set `status=resolved`,
then assign `/grove-work G-260925-pbx81` again on `worktree-G-260925-pbx81`.

## Next

Waits for the owner.

## Answer
Use the G-260924-frzeg settings
