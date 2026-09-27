---
id: "G-260927-jds8s"
type: question
title: "What mandate covers rerunning the retrieval cases after the context framing trim?"
status: open
created: "2026-09-27T22:13:37Z"
updated: "2026-09-27T22:28:19Z"
blocks: ["G-260927-n4wvk"]
relates_to: ["G-260923-659zw", "G-260925-pbx81"]
---

## Question

[G-260927-n4wvk](G-260927-n4wvk-trim-the-framing-context.md) acceptance 4
reruns the retrieval cases of
[G-260925-pbx81](G-260925-pbx81-evaluate-whether-agents.md) as a `with` row
after the framing trim. What a paid step spends is the owner's to set, and
the previous mandate ([G-260923-659zw](G-260923-659zw-what-mandate-should-the.md))
covered G-260923-p5pt6's runs only.

Options:

1. **Reuse the last configuration** (recommended): `--model claude-opus-5-5
   --runs 5 --budget 5 --permission-mode auto --config-dir
   ~/.cache/grove-evals/claude`, cases `listed-constraint` and
   `code-constraint`. The cap is 2 cases × 5 runs × $5 = $50; the without
   row ([G-260925-khwkq](G-260925-khwkq-without-row-both-constra.md)) is the
   comparison.
2. **Three runs per case**, same otherwise: cap $30, patterns only.
3. **No rerun**: accept G-260927-n4wvk on its deterministic checks and owner
   judgment alone, and drop its acceptance 4.

Who can answer: the owner. The answer is the mandate; the implementing
session runs `python3 evals/run.py run --case …` with exactly those values
and reads the report into a review record.

## Next

Waits for the owner. Answer here, set `status=resolved`, then assign
G-260927-n4wvk.
