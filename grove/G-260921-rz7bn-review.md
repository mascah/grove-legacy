---
id: "G-260921-rz7bn"
type: term
title: "Review"
status: settled
created: "2026-09-21T05:01:56Z"
updated: "2026-09-21T14:23:33Z"
relates_to: ["G-260921-sth8q", "G-260921-jatts", "G-260921-btyck"]
formerly: "T-005"
---

## Meaning

An examination of a [candidate](G-260921-jatts-candidate.md) against the work's
acceptance and evidence, recorded with what it looked at. An independent agent
review looks for defects the implementer missed; a human review leads with
changed behavior, decisions, open issues, and checks before any diff.

A review is evidence, never a verdict by itself: it is not
[approval](G-260921-btyck-approval.md), and a self-check is not an independent review.
The word also names the work status between active and accepted, in which a
candidate awaits human judgment, and the standing of accepted work whose
acceptance no longer applies; say "review record" or "Review status" when
the difference matters.

## Relationships

A review record names its work and the commit it examined; work in Review
status names its `candidate`, so the two can be compared. Feedback that
starts another [attempt](G-260921-sth8q-attempt.md) keeps earlier reviews.
