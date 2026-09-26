---
id: "G-260926-a8vyj"
type: term
title: "Policy"
status: proposed
created: "2026-09-26T03:29:15Z"
updated: "2026-09-26T03:29:23Z"
relates_to: ["G-260921-jatts", "G-260921-rz7bn", "G-260921-btyck", "G-260921-3qgsf", "G-260925-5wrn8", "G-260925-wh9ax", "G-260926-kfcpp"]
---

## Meaning

The owner's standing instruction, written as `policy:` in `grove.yaml`, for
what may happen to a [candidate](G-260921-jatts-candidate.md) in review with no
per-candidate human act: start one [resolution](G-260926-kfcpp-resolution.md)
attempt for a conflict, and approve and integrate a candidate that meets
the written conditions once its merged result passed the policy's
verification. Its absence means nothing is automatic. `grove sweep` applies
it ([G-260925-5wrn8](G-260925-5wrn8-resolve-approve-and-inte.md)), and each act it takes is
**delegated**: attributed to the policy's `grove.yaml` revision and the
evidence it relied on, told apart from the owner's own verdict, and
reversible by an ordinary revert. The delegation itself is decision
[G-260925-wh9ax](G-260925-wh9ax-delegate-conflict-resolu.md).

## Relationships

A delegated [approval](G-260921-btyck-approval.md) is approval by someone the owner
delegated to, which that term already allows. A [review](G-260921-rz7bn-review.md)
stays evidence: the policy reads its closing line, `Open findings: none`,
and the review gives no verdict. [Integration](G-260921-3qgsf-integration.md) under
a policy is the ordinary merge and `done`, preceded by the verification.
