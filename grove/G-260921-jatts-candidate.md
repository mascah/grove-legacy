---
id: "G-260921-jatts"
type: term
title: "Candidate"
status: settled
created: "2026-09-21T05:01:56Z"
updated: "2026-09-21T14:23:32Z"
relates_to: ["G-260921-sth8q", "G-260921-rz7bn", "G-260921-btyck", "G-260921-3qgsf"]
formerly: "T-004"
---

## Meaning

The specific result of an [attempt](G-260921-sth8q-attempt.md) that is offered for
judgment: a Git commit on a work branch, together with the evidence gathered
at that commit. It is exact on purpose, so that what was examined, what was
approved, and what gets integrated can be shown to be the same thing.

A branch name is not a candidate, because it moves. A candidate that changes
is a new candidate.

## Relationships

A [review](G-260921-rz7bn-review.md) examines a candidate. [Approval](G-260921-btyck-approval.md)
is of a candidate. [Integration](G-260921-3qgsf-integration.md) puts an approved
candidate into the target.

Since [G-260929-gm3m4](G-260929-gm3m4-clean-main-history-with.md) the target
receives it as one squash commit, not the candidate itself: the commit's
trailers name the candidate and the submitted tip, retained under
`refs/grove/submitted/`, and its tree must be exactly that tip merged onto
the commit's parent, which is what shows the integrated result is the
candidate that was accepted.
