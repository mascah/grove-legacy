---
id: "G-260928-wp37y"
type: review
title: "Review of G-260927-fwcrr: merge of main at 2633d8d"
status: current
created: "2026-09-28T16:56:36Z"
updated: "2026-09-28T16:56:50Z"
work: ["G-260927-fwcrr"]
examined: "ecca7cf4754f0996e885593fd38545d0125bfbd0"
---

## Examined

[G-260927-fwcrr](G-260927-fwcrr-make-cli-refusals-name-t.md) on
`worktree-G-260927-fwcrr`: the resolution of the merge of main `2633d8d`
into the branch at `7e10a56` (previous candidate `16af7ee`), scoped to the
resolution as the feedback asks: `git show --remerge-diff`, against
main's changes to `docs/record-model.md` since `bb39668`, the record's
acceptance 2, 3 and 5, and the [Shipped document](G-260925-khfe7-shipped-document.md)
rules. Two rounds, each by a fresh `grove-reviewer` agent, read-only:
round 1 at `1ed1f16` (later amended), round 2 at `ecca7cf`. Both ran
`grove check` and the short tests of the root, `cli`, `project` and
`update` packages (pass); the full uncached suite ran separately at
`ecca7cf` (see the work record's Evidence).

## Findings

Round 1 (`1ed1f16`), one: main added `approved_by` to the model's
`show --json` sentence, which the branch had moved to `docs/commands.md`'s
`show ID` bullet; the resolution left that bullet listing only
`{id, path, revision, source}`, relying on the policy paragraph that does
not say when the key appears. Main's two other changes (a small size
selects the compact handoff; a review record before Review only where the
handoff calls for one) were carried correctly, the model was 11,604 bytes
and shipped-document clean, and every `record-model.md#` anchor resolved.

Round 2 (`ecca7cf`): the fix verified against `internal/cli/cli.go`'s
`show --json` code and `update.Delegated`, the new `#judging-and-integrating`
link resolves, the remerge-diff touches only the three documents, and
nothing regressed. Noted as not blocking: the linked section explains the
verdict, and the delegated prefix is explained one section later (Sweep).

## Disposition

The round 1 finding was fixed by amending the merge, so the whole
resolution shows in `git show --remerge-diff ecca7cf`: the `show ID` bullet
now adds `approved_by` while `approved` is set, derived from the latest
verdict on the candidate, linking Judging and integrating.

Open findings: none
