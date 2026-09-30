---
id: "G-260921-3qgsf"
type: term
title: "Integration"
status: settled
created: "2026-09-21T05:01:57Z"
updated: "2026-09-30T03:18:06Z"
relates_to: ["G-260921-vr8a8", "G-260921-jatts", "G-260921-btyck", "G-260930-e8jj7", "G-260929-gm3m4", "G-260930-4742q", "G-260930-2qa4a"]
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

## Selected redesign, 2026-09-29

The owner selected [G-260930-e8jj7](G-260930-e8jj7-build-a-portable-workflo.md), reopening the ancestry-only
representation for squash and hosted delivery. The current mechanism above
remains implemented; [G-260929-gm3m4](G-260929-gm3m4-clean-main-history-with.md)
and [G-260930-4742q](G-260930-4742q-deliver-accepted-work-th.md) must reconcile this term and the lifecycle contract when
they deliver verified correspondence between an approved candidate and a
transformed integrated result. An unverified trailer or work status alone
will not establish integration.

## Selected completion contract, 2026-09-29

The owner selected [G-260930-2qa4a](G-260930-2qa4a-derive-done-from-recorde.md):
record candidate acceptance and derive Done from verified target delivery,
with a redesigned work record that does not remain misleadingly in review.
This changes representation, not the distinction between approval and
integration.

## Schema 4, 2026-09-30

[G-260929-gm3m4](G-260929-gm3m4-clean-main-history-with.md) implements it and
supersedes the mechanism in Meaning above. `grove integrate` retains the
branch's tip under `refs/grove/submitted/`, writes that tip merged onto the
target's tip as one commit with the target tip as its only parent and a
Conventional Commit message whose trailers name the work, the candidate and
the retained tip, fast-forwards the target to it, and writes no record.
Done is derived: accepted work is done when the target contains its
candidate, or a commit whose trailers, parent and exact tree match the
retained submission. A forged or altered delivery, or one a rewrite of the
target dropped, is not done, and one whose submission a clone lacks is
unknown there; a later revert does not undo
it. A `done` status from schema 3 is kept by `grove migrate` as that
schema's claim and never presented as verified.
