---
id: "G-260928-a372k"
type: review
title: "Review of G-260927-fwcrr: refusals and the record model split"
status: current
created: "2026-09-28T16:42:07Z"
updated: "2026-09-28T16:42:21Z"
work: ["G-260927-fwcrr"]
examined: "c0fb4472dcd7dab110d13034b369048a4a95d724"
---

## Examined

[G-260927-fwcrr](G-260927-fwcrr-make-cli-refusals-name-t.md) on
`worktree-G-260927-fwcrr`: the change from main `bb39668` to `c0fb447`,
against the record's acceptance 1 to 5, its plan
[G-260928-rv8z2](G-260928-rv8z2-refusals-and-model-plan.md), the
[Shipped document](G-260925-khfe7-shipped-document.md) term, and the
repository's instructions. Two rounds, each by a fresh `grove-reviewer`
agent, read-only: round 1 at `e420055`, round 2 at `c0fb447`. Both ran
`go vet ./...`, `gofmt -l .`, `grove check` and the full uncached suite
(all pass), and reproduced the refusals with a built binary in a disposable
clone.

## Findings

Round 1 (`e420055`), six:

1. The `approved` status rule fired beside an invalid or missing status and
   prescribed the wrong fix (`not bogus: unset approved`).
2. The list refusal's example `["W-001"]` is an ID the model calls invalid.
3. `new` and `convert` hard-coded the type list, outside the tables and the
   refusal test.
4. Three model statements untrue of the code: a superseded decision naming
   its replacement (unenforced), "`status` cannot be unset" (it can toward a
   page), convert keeping the body's bytes (it strips a BOM).
5. Two command facts (the brief's form check on committed sources; a page
   and `context`) lived only in the unshipped design document.
6. The Shipped document term's "Not shipped" list lacked the record design.

The reviewer judged the plan's correction (no done-stays-done rule is
enforced, so the refusal names the approval rule) a faithful reading of the
record, to be stated for the owner in Evidence.

Round 2 (`c0fb447`): each fix verified by reading and by running; no new
finding. Noted as not blocking: `model_test.go` checks the policy subkeys
by substring and does not tie priority's 1 to 5 to code.

## Disposition

All six fixed in `c0fb447`: the approval rule now applies only to a valid
status, with a single-line test; the example is `["G-260925-7k2qm"]`;
`project.TypeNames()` builds the `new`, `convert` and `check` type lists,
each asserted; the model moves supersession to "Not enforced" and states
unset and convert as the code does; the two facts moved to
`docs/commands.md` (Versions, Context); the term names the record design.

Open findings: none
