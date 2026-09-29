---
id: "G-260929-gm3m4"
type: work
title: "Clean main history with local-first squash integration"
status: proposed
created: "2026-09-29T16:15:25Z"
updated: "2026-09-29T16:15:25Z"
---

## Outcome

`main` reads as a clean Conventional Commits history, one commit per
integrated work item, so any changelog or release tool (release-please,
git-cliff, conventional-changelog) produces useful notes, whether the
project is local only, private on GitHub, or public.

Owner intent, conversation 2026-09-29: "fix the commit problem more than
anything so that things are not noisy no matter what tool people use", and
eventually set up release-please for Grove itself to distribute releases to
friends. Parked here unshaped: the conversation covered too many threads at
once, and this record keeps them apart.

## Constraints

Where the conversation ended (proposed, not selected):

- **Local first.** Some projects will never be hosted; on private GitHub
  repositories Actions minutes cost money. Grove's local integration, with
  `policy.approve.verify` running checks on the merged result, is the free
  replacement for hosted CI and merging, so it stays. Hosting is optional
  per project, and local output should match what a host would produce.
- **The one change:** `integrate` squashes, writing the commit message from
  the record (type from `kind`, subject from the title) with `Grove-Work`
  and `Grove-Candidate` trailers. Work branches may then be as noisy as
  they like; nothing but the squash reaches `main`.

Observed 2026-09-29:

- Squashing breaks only the post-merge "delivered" checks, which assume
  the candidate is an ancestor of `main`: `internal/update/update.go:351`
  (`done`), `internal/deps/deps.go:358` (dependency delivery),
  `internal/integrate/integrate.go:103-107`, and `just clean-merged`
  (`--merged main`). One fix is "delivered = ancestor, or a commit on the
  ref whose trailer names the candidate". The branch's evidence commits
  then need a ref to stay reachable.
- This repository is trunk-based with short-lived per-item branches in
  worktrees, merged locally by `integrate`, and no PRs. Since 2026-09-24,
  145 of 189 commits made directly on `main` touch only records, all typed
  `docs`. release-please hides `docs` by default, but git-cliff's default
  does not. Their type is an open choice (`chore(grove)` or keep `docs`).
- In keyborg, 7 of 14 `fix` commits on `main` are review-round fixes,
  which a changelog would list as bug fixes.

Other threads from the same conversation, each for its own shaping later,
not this record's scope:

1. Release records and a roadmap view: membership on the release record
   (owner); Grove does not own version numbers or replace release tools
   (owner).
2. Orchestration: one session with subagents versus separate processes;
   test the reviewer's per-call `model` first. G-260929-1w9x4's finding 3
   is wrong: headless subagents work (owner); the denial was a write into
   a sibling worktree.
3. A digest of "what happened since I last looked", against comprehension
   debt.
4. Deferred until a project needs it: a PR mode for protected `main`, and
   deriving `done` instead of writing it.

## Acceptance

To be shaped. At least: one integration leaves one conventional commit on
`main` with the trailers; `done`, dependency delivery and branch cleanup
recognize a squashed candidate; the record-only commit type is decided.

## Next

Parked, unshaped. Shape with `/grove-shape` on this record's Outcome and
Constraints only; take the other threads one at a time.
