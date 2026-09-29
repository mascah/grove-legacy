---
id: "G-260929-s0f25"
type: question
title: "What budget may the real-provider trial of G-260928-dtrnw spend?"
status: open
created: "2026-09-29T00:57:23Z"
updated: "2026-09-29T00:57:48Z"
blocks: ["G-260928-dtrnw"]
---

## Question

Acceptance 5 of [G-260928-dtrnw](G-260928-dtrnw-run-sweep-from-a-finishi.md)
asks for one trial with the real provider in a disposable project "under a
budget the owner names at assignment". The assignment, headless, named
none, so the trial was not run. Please answer each item explicitly; a blank
item keeps this question open:

1. **Budget:** the `--budget` in USD for the trial attempt.
2. **Model and effort:** for the trial, or `run:` defaults.
3. **Who runs it:** a resumed `/grove-work G-260928-dtrnw` session, or you.

The trial: a disposable clone with a `policy:` (verify `true`, `max_lines`
small, `integrate: true`) and one small proposed work record; `grove run`
it; the attempt should hand off, and its owner's sweep should integrate it,
with `grove attempt ATTEMPT` showing the `Sweep:` lines.

Recommendation, not a decision: 3 USD, model sonnet, effort medium, run by
the resumed session.

## Next

Answer in `## Answer`, then resolve; G-260928-dtrnw resumes.
