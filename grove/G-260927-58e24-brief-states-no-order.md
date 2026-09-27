---
id: "G-260927-58e24"
type: work
title: "The brief states no order among work records"
status: proposed
created: "2026-09-27T20:13:00Z"
updated: "2026-09-27T20:13:00Z"
---

## Outcome

This repository's brief follows decision G-260927-y3pc2 and states no order
among work records. Owner's intent, 2026-09-27: asked to fix the brief's
"Suggested sequence" once G-260927-k4xwq's review found it at odds with the
decision.

## Constraints

Observed at `main` `013839c`: the brief's "Suggested sequence" names the
preview phase's three pieces of work (G-260923-p5pt6, G-260923-895zb,
G-260923-gsthp), says none is a prerequisite of another, and says "Their
order is proposed: the owner has not selected one". The decision allows
prose to repeat an order the edges show and a brief to set direction for
outcomes not yet shaped as work, never to state an order the edges lack.

- Keep the heading: G-260919-rt9h9 links `brief.md#suggested-sequence`.
- Change only that section, plus the decision's sentence that called it the
  owner's to change. G-260927-k4xwq is done, and its record keeps what was
  true at its candidate.
- The same text was first written as `142bfd2` on G-260927-k4xwq's branch
  after the owner had approved `80c93fc`, so it was never integrated. This
  record carries it instead.

## Acceptance

- The section states no order among work records. It points to their
  `depends_on` edges and G-260927-y3pc2, and keeps the preview phase's three
  pieces of work, the owner-judges-readiness line, and the history links.
- The heading and every link in the brief still resolve, and `grove check`
  passes.

## Next

Active on `worktree-G-260927-58e24`.
