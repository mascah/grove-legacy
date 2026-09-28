---
id: "G-260928-d8py6"
type: decision
title: "A standing policy may delegate the approval judgment to an LLM"
status: accepted
created: "2026-09-28T19:28:59Z"
updated: "2026-09-28T19:34:00Z"
relates_to: ["G-260925-wh9ax", "G-260925-5wrn8", "G-260926-a8vyj", "G-260921-btyck", "G-260921-rz7bn", "G-260928-c5j9d"]
---

## Decision

## Alternatives

## Reconsideration

## Decision

On 2026-09-28, asked in a shaping conversation what should make a candidate
safe to integrate without a per-candidate human act, the owner answered:
"An LLM decides ideally." The standing policy of
[G-260925-wh9ax](G-260925-wh9ax-delegate-conflict-resolu.md) may therefore
name, as one of its written conditions, a delegated judgment: a bounded
headless call that reads the candidate, its record's outcome, constraints
and acceptance, its reviews and evidence, and returns a structured verdict,
approve or wait for the owner, with reasons. The judgment is one condition
beside the deterministic ones; `never`, `verify` and the review's closing
line stay. Its model, effort and spend are keys of the policy, so the owner
writes them once, and each act it leads to is attributed to the policy
revision and the judgment's recorded output, told apart from the owner's
verdict as every delegated act is.
[G-260928-c5j9d](G-260928-c5j9d-judge-a-candidate-agains.md) carries it out.

This extends G-260925-wh9ax, whose alternatives rejected "the review agent
decides by itself, with no written policy". That rejection stands: the
review remains evidence, and the judgment is a separate act the policy
names, with its own record of what it read and returned.

## Alternatives

- **Deterministic proxies only** (`kinds`, `sizes`, `max_lines`, `paths`):
  attributable and cheap, but they cannot say whether a candidate did what
  the record asked, which is what the owner judges by.
- **The implementing session's own review decides**: rejected in
  G-260925-wh9ax and here; a session judging its own work is not
  independent.
- **Always wait for the owner**: the state before G-260925-wh9ax. The owner
  named review as the phase where work sits idle.

## Reconsideration

Reopen the judgment's conditions when a delegated approval lands a candidate
the owner would have sent back, or when its wait verdicts are so frequent
that nothing is delegated. Reopen the delegation when the owner stops
unattended use.
