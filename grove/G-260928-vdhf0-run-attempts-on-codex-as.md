---
id: "G-260928-vdhf0"
type: decision
title: "Run attempts on Codex as a second provider"
status: accepted
created: "2026-09-28T19:28:59Z"
updated: "2026-09-30T01:16:11Z"
relates_to: ["G-260923-tnn5e", "G-260928-917h8", "G-260924-59f5k", "G-260925-42j50", "G-260928-y2p5h", "G-260928-n4f1q", "G-260930-e8jj7", "G-260930-62nmj"]
---

## Decision

On 2026-09-28, in a shaping conversation, the owner chose to add Codex as a
second [provider](G-260928-917h8-provider.md) for attempts, answering "Yes" to
reopening the disposition of
[G-260925-42j50](G-260925-42j50-codex-eval-row-pattern-p.md) ("a Codex provider
is not worth shaping now", 2026-09-25). The owner's stated intent, same
conversation: Grove is the control plane over harnesses, and an attempt's
harness, model, effort and cap should be the owner's choice per launch, so
that one phase can run on Codex and the next on Claude Code.

[G-260923-tnn5e](G-260923-tnn5e-run-attempts-as-a-grove.md) stands in
everything but its pinning to `claude -p`: an attempt is still a Grove-owned
process under an on-demand owner process, with events written to files, a
required spend bound and permission setting, Grove-owned worktree, identity,
duplicate-start refusal, stop and owner-loss handling, and no resident
service. What changes: the process may be `codex exec`, with the cap and
permission expressed in that harness's terms, and the provider is recorded
per attempt. [G-260928-y2p5h](G-260928-y2p5h-run-an-attempt-on-codex.md)
carries it out.

## Alternatives

- **Claude Code only**, the 2026-09-25 disposition: the simplest runner, and
  per-phase harness choice impossible. Set aside by the owner's changed
  intent.
- **A generic provider interface first, with one implementation**: an
  abstraction nothing exercises. Rejected: the seam is drawn by building the
  second provider.
- **Every installed harness**: the owner named opencode, which is not
  installed here and has no evidence. A third provider is shaped when the
  owner names it.

## Reconsideration

Reopen when Codex cannot be capped or stopped reliably as a Grove-owned
process, when its sandbox cannot be confined to the worktree (G-260925-42j50
observed reads outside the clone), or when the work guide behaves
differently enough under Codex that one workflow cannot serve both.

## Clarification, 2026-09-29

The owner selected [G-260930-e8jj7](G-260930-e8jj7-build-a-portable-workflo.md). The two-provider direction and
on-demand ownership constraints stand. Substantial runner replacement and
refactoring are permitted; preserving the existing package shape or exact
Claude command bytes is not required. [G-260930-62nmj](G-260930-62nmj-design-the-portable-work.md) owns the proposed
journey and contracts; the real two-provider trial tests the resulting
boundary. No public plugin framework is selected.
