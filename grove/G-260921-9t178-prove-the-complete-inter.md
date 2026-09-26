---
id: "G-260921-9t178"
type: work
title: "Prove the complete interactive workflow on real Grove work"
status: done
created: "2026-09-21T00:54:14Z"
updated: "2026-09-22T15:49:17Z"
kind: investigation
size: small
priority: 1
depends_on: ["G-260919-04z88", "G-260921-w9x25", "G-260921-9wkjt"]
relates_to: ["G-260921-tkdwh", "G-260921-407n6", "G-260921-gtydy", "G-260921-ebsby"]
formerly: "W-021"
candidate: "3e6b9d4"
---

## Outcome

Prove Grove's complete interactive shape → prepare → implement → independent
review → human review → integrate loop on actual Grove development.

## Scope and evidence

Use delivered G-260919-04z88/G-260921-w9x25/G-260921-ebsby/G-260921-9wkjt behavior on a newly selected, bounded real change.
Choose a small real assignment with the owner after the prerequisites land.
Do not count a prerequisite's earlier implementation as a trial of behavior that
was not available then, rerun completed work, or invent demonstration records.
This is a connected-workflow trial, not another implementation of those features.

## Acceptance

1. Fresh interactive harness invocation discovers the guides and appropriate
   context, shapes useful work and records actual knowledge changes.
2. Implementation prepares proportionally, escalates only genuine human
   decisions, retains checkpoints and finishes with linked review evidence.
3. The owner can review the result from files/CLI without reconstructing chat,
   request a follow-up or accept it, and complete explicit local integration.
4. Observe fresh-session continuation. Exercise missing-human and feedback
   cases in disposable fixtures if they do not arise naturally.
5. Record invocation, harness/version, source and candidate revisions, actual
   outcomes, owner feedback, observed friction and the next improvement. Separate
   real trials, simulations and unsupported/unexercised harness behavior.
6. Observe useful general knowledge captured without a new schema type, staged
   retrieval of it when relevant, and stable IDs/paths through work completion.
   Do not manufacture new categories or knowledge solely to satisfy the trial.

## Preparation

The prerequisites are delivered on main at `5f07c94`: G-260921-9wkjt's candidate
`fae1e4c` is an ancestor, and G-260919-04z88, G-260921-w9x25 and G-260921-ebsby show as the shaping guide
with both skill adapters, the `brief:` key with plan and review records, and
neutral IDs with `page` records. The owner selected the real change on
2026-09-22: a status filter for `grove list`, observed missing at `5f07c94`.
The trial plan is [G-260922-10dj4](G-260922-10dj4-trial-plan-for-the-inter.md).
The user must supply the actual usability verdict. This record does not itself
authorize merging any trial branch: obtain that work's explicit disposition.
Feed lessons into their owning guide or work record and the adoption milestone.

## Evidence

Trial run 2026-09-22 under plan G-260922-10dj4 on branch `worktree-G-260921-9t178` (worktree
`.claude/worktrees/G-260921-9t178`), rebased from main `5f07c94` onto `e8bcf03` once
the real change was integrated there; started from this record at revision
`43830baa…` and G-260922-10dj4 at `8dedc7c6…`, neither changed by the rebase. The
evidence per acceptance item, with each observation labelled real,
simulation or unexercised, is the review record
[G-260922-08wxx](G-260922-08wxx-trial-evidence-for-the-i.md), examined `840a78b`, the
candidate of the real change [G-260922-w53bz](G-260922-w53bz-filter-grove-list-by-sta.md)
whose loop is main `ff0f3e9..e8bcf03`. In short: items 1, 2, 4 (continuation)
and 6 were met on real sessions; item 3 was met with one gap, the owner's
verdict not quoted at done; item 4's feedback and missing-human cases were
simulated in disposable clones, where the headless shaping call diverged from
the guide by settling product choices itself; item 5 is G-260922-08wxx and this
section. The candidate is the commit that adds this section.

Decisions: the rebase, so that G-260922-w53bz and G-260922-22hyw resolve as links from G-260922-08wxx;
the headless simulation ran `claude -p` with write tools allowed on the
clone, which used Claude Opus 5 rather than the interactive model. Lessons
went to their owners in the same commit: `docs/work-execution.md` step 8 (the
handoff's integration commands include the verdict edit),
`docs/work-shaping.md` headless bounds (a choice the acceptance depends on is
a missing choice), and G-260921-407n6's Next (open improvements).

Verification at the candidate: `go run ./cmd/grove check` → `OK: 78
records`; every link written here and in G-260922-08wxx resolves to a file in this
checkout; no Go source changed, so `go vet`, `gofmt` and the suite are
unaffected. Review: G-260922-08wxx is this session's evidence, not an independent
review; the record and plan ask for none, so that omission is reported here
as open rather than blocking.

Limits: the headless work row and the review cap were not
exercised; the clones were left under the session scratchpad and are not
part of this branch.

## Next

Accepted and integrated. The owner's verdict on candidate `3e6b9d4`, given
in the G-260921-9t178 review session on 2026-09-22: "it went as expected". Merged
fast-forward into main and marked done there. Demo, from any checkout:

```sh
cd .claude/worktrees/G-260921-9t178
go run ./cmd/grove context G-260921-9t178 --include grove/G-260922-08wxx-trial-evidence-for-the-i.md
go run ./cmd/grove show G-260922-08wxx
git diff --stat main
```

Nothing further for this record; open improvements live in G-260921-407n6's Next.
