---
id: "G-261001-fg05z"
type: decision
title: "Simplify workflow boundaries before extending portable execution"
status: accepted
created: "2026-10-01T02:17:08Z"
updated: "2026-10-01T02:22:08Z"
relates_to: ["G-260930-tcc9w", "G-260930-84fnb", "G-260930-0s29t", "G-261001-mbjwz", "G-261001-62n4p"]
---

## Decision

On 2026-09-30 (owner's local date), after reviewing the repeated repairs in
G-260929-gm3m4 and G-260930-0s29t, the owner agreed to the five recommendations
in the architectural assessment and said, "I want to do all the things you
recommended."

- Retain selective rebuilding, schema 4, cheap completion reads and one
  workspace per delivery. Keep G-260930-0s29t unmerged pending the owner's
  judgment of its demonstrated behavior; agreement with these recommendations
  is not that judgment or approval of its candidate.
- State the workspace rules by operation, with implementation ownership and
  compatibility limits visible. Judge simplification by fewer independently
  maintained rules and removed obsolete recoveries, not moved files.
- Settle the missing-evidence boundary before portable execution relies on
  it. A clone can read Done without audit refs; eligibility to reuse an
  imported implementation branch is a different question. Its technical
  mechanism remains to be designed and demonstrated.
- Give workspace ownership and eligibility a narrow shared boundary used by
  delivery and execution, ahead of the two-provider implementation.
- Distinguish findings that block acceptance from recorded follow-ups in the
  review workflow and delegated approval contract. Preserve correctness,
  incomplete-review refusal, candidate binding and explicit authority.

[The contracts](G-260930-84fnb-portable-workflow-contra.md) own the design;
[the review work](G-261001-mbjwz-separate-blocking-review.md) and
[the workspace work](G-261001-62n4p-isolate-workspace-eligib.md) own acceptance
and next actions. This authorizes follow-through on those recommendations,
not a provider launch, paid trial, merge or push. Existing commands retain
their current contract until their replacement is delivered.

## Alternatives

Continuing provider work immediately would build on the unsettled workspace
boundary. Repeated unrestricted reviews of the same broad change would not
settle that boundary. A whole-application rewrite would discard established
record, Git and terminal behavior without first demonstrating a simpler
replacement. The owner chose targeted simplification and smaller outcomes.

## Reconsideration

Reconsider if a bounded local demonstration cannot explain the workflow, or
the shared boundary requires another authoritative lifecycle store. Preserve
the candidate and evidence while comparing concrete alternatives. No review
threshold permits an unmet acceptance item to be relabelled a follow-up.
