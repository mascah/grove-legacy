---
id: "G-260924-7x7p7"
type: question
title: "What mandate and login should the G-260924-59f5k Codex eval runs use?"
status: resolved
created: "2026-09-24T22:37:54Z"
updated: "2026-09-24T23:21:55Z"
blocks: ["G-260924-59f5k"]
---

## Question

[G-260924-59f5k](G-260924-59f5k-run-the-eval-pair-on-cod.md)'s paid Codex runs need
the owner's mandate, which its Next says is agreed at assignment, as
[G-260923-659zw](G-260923-659zw-what-mandate-should-the.md) and
[G-260924-y99bx](G-260924-y99bx-how-do-the-eval-runs-log.md) did for Claude. Plan
[G-260924-w07wn](G-260924-w07wn-codex-eval-row-plan.md) builds the runner and its
selftest without them; nothing is spent until this is resolved. Each item
is a runner argument, with the recommendation first.

1. **Login of the clean `CODEX_HOME`.** A new home has no login: observed
   2026-09-24, codex-cli 0.156.1, `codex exec` in one exits 1 with 401 at
   no cost. `/login` needs a browser, so a headless session cannot do it.
   - Recommended: `CODEX_HOME=~/.cache/grove-evals/codex codex login`
     once (ChatGPT sign-in; runs count against the plan's limits and no
     dollar figure exists).
   - Or an API key, `printenv OPENAI_API_KEY | CODEX_HOME=DIR codex login
     --with-api-key`, or `CODEX_API_KEY` exported in the shell that starts
     the runs (API billing; Codex still reports tokens, not dollars).
   The runner refuses a home holding `AGENTS.md`, rules, prompts, hooks,
   memories, authored skills, or `config.toml` settings beyond the
   per-directory trust tables Codex writes itself.
2. **Runs per case:** 5, as G-260924-frzeg.
3. **Model:** `gpt-6-astra`, the default in the owner's
   `~/.codex/config.toml`, or another named.
4. **Reasoning effort:** `high`, the owner's configured default.
5. **Permission mode:** `approve-for-me` (workspace-write sandbox with
   Codex's automatic reviewer on approval requests), the nearest to the
   Claude row's `auto`; or `workspace-write` alone, where a refused
   command is left to the agent, as G-260921-ahbrj ran.
6. **Cap:** 600 seconds per run, killed at the cap; 2 cases × 5 runs ×
   600 s is at most 100 minutes. G-260924-frzeg's Claude runs took 45 to 62
   seconds. Codex has no budget flag, so no dollar cap exists.
7. **A same-digest Claude pair:** yes, recommended. G-260924-frzeg ran at guides
   digest `3f5487904c61`; this branch is at `0c163c41a0f2`, so a Codex
   difference could be the guide, not the harness. Rerunning the pair with
   G-260923-659zw's values (`claude-opus-5-5`, 5 runs, $5 per run, `auto`,
   `~/.cache/grove-evals/claude`) costs about $3 against a $50 cap.

Who can answer: the owner, by doing item 1 and setting this record
`resolved` with the values in Next (any item left blank takes the
recommendation). Then launch `/grove-work G-260924-59f5k` or `$grove-work G-260924-59f5k`
on `worktree-G-260924-59f5k`.

## Next

## Answer
I ran `CODEX_HOME=~/.cache/grove-evals/codex codex login` 
