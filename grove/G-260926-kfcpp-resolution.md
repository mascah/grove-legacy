---
id: "G-260926-kfcpp"
type: term
title: "Resolution"
status: proposed
created: "2026-09-26T01:39:11Z"
updated: "2026-09-26T01:39:21Z"
relates_to: ["G-260921-sth8q", "G-260921-jatts", "G-260921-rz7bn", "G-260921-btyck", "G-260921-3qgsf", "G-260925-dz10z", "G-260925-wh9ax"]
---

## Meaning

The merge of a target commit into the branch of a
[candidate](G-260921-jatts-candidate.md) that conflicts with the target, with its
conflicts settled, handed off as a new candidate. It is a merge, never a
rebase, so the earlier candidate and everything reviewed at it stay
ancestors of the new one. What is judged is the resolution: the merged
target commit, the files the merge conflicted on and how each was settled.
Taking one side of a file drops the other side's change, and the
resolution names that.

A **resolution attempt** is one [attempt](G-260921-sth8q-attempt.md) whose whole
mandate is that merge, given as feedback that names the target commit and
the files. `grove resolve` starts one by hand
([G-260925-dz10z](G-260925-dz10z-update-a-conflicting-can.md)); a standing owner policy may
start one ([G-260925-wh9ax](G-260925-wh9ax-delegate-conflict-resolu.md)). This sense is
unrelated to a question being resolved, which means it was answered.

## Relationships

It produces a candidate, which [review](G-260921-rz7bn-review.md) examines and
[approval](G-260921-btyck-approval.md) accepts as for any other. It changes nothing
about [integration](G-260921-3qgsf-integration.md), which merges the new candidate
as usual.
