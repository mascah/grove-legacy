---
id: "G-260921-3qgsf"
type: term
title: "Integration"
status: settled
created: "2026-09-21T05:01:57Z"
updated: "2026-09-21T14:23:34Z"
relates_to: ["G-260921-vr8a8", "G-260921-jatts", "G-260921-btyck"]
formerly: "T-007"
---

## Meaning

An approved [candidate](G-260921-jatts-candidate.md) becoming part of the project's
configured target, such as a merge into the main branch, shown by Git
ancestry rather than by a status.

Since G-260921-9wkjt, done means exactly this for work with a `candidate`: `update`
writes done only where that commit is already an ancestor of HEAD. A done
record without a candidate is older and claims its outcome only in its own
branch; it may never have been merged, and it keeps that meaning.

Since G-260921-jwk4e, `grove integrate` performs it from the target's checkout: a
plain merge of the one branch holding an approved candidate, aborted on
conflict, then done written and committed there, with the branch and its
worktree removed only on request and only where Git agrees. Since G-260925-h8rj5,
a conflict is predicted with `git merge-tree` in objects only and refused
before the merge starts, naming the files; the abort remains for what the
prediction cannot see.

## Relationships

Follows [approval](G-260921-btyck-approval.md). Completes implementation
[work](G-260921-vr8a8-work.md) under the target lifecycle.
