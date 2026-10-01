---
id: "G-261001-hcb0h"
type: plan
title: "Plan review summaries and explicit follow-up delegation"
status: current
created: "2026-10-01T02:31:58Z"
updated: "2026-10-01T02:32:55Z"
work: ["G-261001-mbjwz"]
---

## Design

Implements G-261001-mbjwz from main bf86caa. The owner's approved record
bounds the change to reports, their existing policy consumer and the guides.

A new review ends with one strict, human-readable, versioned line:

```text
Review v1: complete; blockers=0; follow-ups=1
```

The state is complete or incomplete, and counts are nonnegative decimal
integers. A shared parser in internal/project reads only the terminal line;
malformed, missing, unsupported-version or overflowed summaries wait. The
legacy terminal `Open findings: none` remains a complete zero-findings
report. Legacy nonzero counts never become follow-ups. A modern report
cannot fall back to a legacy clean line appended beneath it.

The report body retains numbered findings, evidence and dispositions. A
blocker is an unmet acceptance item, correctness/authority defect or missing
required verification; misleading documentation can block. Follow-ups are
improvements that do not invalidate acceptance. The implementer preserves
the independent reviewer's categories and summary.

`policy.approve.allow_followups` is an optional boolean, false by default.
Every current applicable review must be complete and blocker-free. A report
with follow-ups additionally requires this explicit opt-in; existing policy
configurations remain strict. No project configuration is opted in by this
implementation. Conflicting reviews cannot outvote a blocker. Existing
candidate checks, required verification, prohibited paths and size bounds
remain in force.

Preview and delegated verdict name the review and its counts without claiming
there are no findings when follow-ups exist. Review summaries are validated
by their consumer, not a schema migration of historical records. Failed
review output remains incomplete; successful process exit is not a report.

## Steps

1. Add failing policy/consumer tests for a complete report with follow-ups,
   default refusal and explicit opt-in; include blockers in another review,
   incomplete/malformed reports, stale evidence and legacy nonzero reports.
2. Implement the shared terminal-summary parser, boolean configuration and
   policy checks. Exercise one real fixture preview and delegated approval.
3. Reconcile the review guide's format and model's configuration, and the work
   guides, command/help text and the contracts' observed/proposed boundary.
   Repair reviews focus on the changed rule and affected callers; preserve
   the cap and acceptance requirements.
4. Run focused tests, then repository final checks once. Commit and obtain
   one independent review scoped to this change; repair consequential findings
   with their neighboring-state evidence and record follow-ups separately.
5. Hand G-261001-mbjwz into review with its candidate. Its delivery is the
   prerequisite for G-261001-62n4p; do not silently merge this new candidate.

## Bounded adjustment

The review guide owns the terminal report grammar, and the model names that
guide while owning delegation configuration. This removes duplicate format
rules and keeps the model within its existing 12 KB budget.

## Execution

Steps 1–4 are implemented at `4a2f64d`; the bounded adjustment above is the
only plan adjustment. The [work record](G-261001-mbjwz-separate-blocking-review.md)
holds verification, its terminal-test limit and the independent review.
Step 5 hands off this candidate for owner judgment; delivery remains the
prerequisite for G-261001-62n4p.
