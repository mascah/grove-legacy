---
id: "G-260928-0pbpf"
type: review
title: "Review of G-260928-4qv1m: rewritten copies"
status: current
created: "2026-09-28T17:36:00Z"
updated: "2026-09-28T17:36:18Z"
work: ["G-260928-4qv1m"]
examined: "905ece29ce1bb09d92235ef61595104ed6347b97"
---

## Examined

[G-260928-4qv1m](G-260928-4qv1m-rewritten-copy.md) on
`worktree-G-260928-4qv1m`: the change from main `ac43184` against the
record's acceptance 1 to 6 (7 is the owner's), its plan
[G-260928-zyqn9](G-260928-zyqn9-rewritten-copy-plan.md), the settled terms
[Approval](G-260921-btyck-approval.md), [Candidate](G-260921-jatts-candidate.md)
and [Integration](G-260921-3qgsf-integration.md), and the repository's
instructions. Two rounds, each by a fresh `grove-reviewer` agent, read-only:
round 1 at `ff18761`, round 2 at `905ece2`. Both ran `go vet ./...`,
`gofmt -l .`, `grove check` and `go test -short` on the changed packages
(all pass), and each reproduced the incident with a built binary in a
disposable repository: integrate by fast-forward, keep the branch, rebase
main onto a new commit; then `integrate`, `resolve`, the printed commands,
and `run --dry-run` of dependent work before and after the printed
`update`. Round 1 also drove the real board through a pty.

## Findings

Round 1 (`ff18761`), three:

1. Medium: `integrate` and `resolve` refused any branch every commit of
   which the target held as a copy, including one landed by a hand
   cherry-pick or rebase-merge whose record on the target was still in
   review; the printed delete commands would then leave that record in
   review with nothing left to integrate (reproduced: base binary merged
   and wrote done, candidate refused).
2. Knowledge: the prerequisite repair printed `--set approved=COPY`,
   carrying an approval to another commit, which the settled Approval term
   says does not carry over; no term or decision records an exception.
3. Low: `resolve` refused "no checkout is on branch …" before the
   rewritten-copy check, telling an owner who had removed the worktree as
   told to add one again.

Observation, not counted: a merge commit on the branch always counts as
without a copy, so a branch that went through `resolve` (which leaves such
a merge), was integrated and then had its target rebased gets only the
partial message.

Round 2 (`905ece2`): each fix verified by reading and by running,
`update --set candidate=… --unset approved --commit` on a done record
accepted and `check` OK, sweep's handling of a refused resolution read (it
waits, writing nothing); no new finding.

## Disposition

All three fixed in `905ece2`: the rewritten-copy refusal and the board's
comparison apply only to work the target holds as done (integrate reads the
target checkout's record, resolve the committed target tip, the board the
on-target state), with `TestIntegrateMergesACopyNotYetDone`, a first call
in `TestResolveRefusesARewrittenCopy` and a tail of
`TestRewrittenCopyExplainedOnOpen`; the repair prints `--unset approved`
where the record is approved, and the docs say why; `resolve` explains a
rewritten copy before asking for a checkout, with `git branch -D` alone
when none is on the branch. The merge ceiling stays as documented in
`internal/versions/copies.go`.

Open findings: none
