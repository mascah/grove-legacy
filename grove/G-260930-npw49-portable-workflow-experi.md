---
id: "G-260930-npw49"
type: page
title: "Portable workflow experience design"
created: "2026-09-30T01:31:25Z"
updated: "2026-09-30T20:22:07Z"
relates_to: ["G-260930-62nmj", "G-260930-84fnb", "G-260930-60c3d", "G-260930-e8jj7", "G-260930-2qa4a", "G-260930-tcc9w"]
---

## Status and scope

Accepted experience baseline for
[G-260930-62nmj](G-260930-62nmj-design-the-portable-work.md), following the
owner-selected direction in
[G-260930-e8jj7](G-260930-e8jj7-build-a-portable-workflo.md).
On 2026-09-29 (owner's local date), the owner reviewed this page and the
contracts at commit f376a6dab7a5999560b25adaacb93d919c36e7a8 and said:
"I've reviewed the designs and I accept them."

The [brief](brief.md) remains the direction owner. This page owns interaction
examples; the [contracts design](G-260930-84fnb-portable-workflow-contra.md)
owns mechanics. The original sketches were accepted at f376a6d. This revision applies the
owner's selected 2026-09-30 boundaries in
[G-260930-tcc9w](G-260930-tcc9w-bound-delivery-groups-an.md); exact preparation
and recovery mechanics in the contracts remain proposed design for review.
Actual terminal behavior still needs implementation and evaluation. Exact
keys, spacing and command spelling remain implementation details. Completion
follows the [derived-Done decision](G-260930-2qa4a-derive-done-from-recorde.md).

## A journey a person can explain

The example is an existing command-line application. The developer wants
"Hide completed tasks from the default list, with an option to show them."
One requirement is ambiguous: should archived tasks appear with that option?

The developer's mental model is one work item moving from an intention to
an accepted result. The main actions are Shape, Start, Answer, Continue,
Review result and Deliver. Worktrees, provider processes and evidence refs
support those actions and remain available through Details.

There is no fixed wizard the person must approve at every phase. Grove
continues within the assignment until it reaches its bound, needs a real
decision, or cannot continue. Applicable policy can supply approval and
delivery, including configured LLM judgment, without a per-item human act.

## 1. Adopt an existing project

Bare grove in a repository without Grove offers setup in the terminal.
Explicit commands remain noninteractive and consume the same setup
plan/apply operation. Existing initialization stays usable through its
documented migration.

The setup reads what is already present and proposes, rather than invents,
the project's configuration. Instructions already in AGENTS.md or CLAUDE.md
remain intact. A brief can be drafted from supplied project context and
confirmed; inferred product choices remain marked as suggestions.

```text
Set up Grove in tasks

Project context     README + existing project instructions
Work execution      Claude Code
Independent review  Codex
Checks              go test ./...       [edit]
Delivery            Local → main       [change]
Limits              Configure before starting work

Changes: add project records and workflow entrypoints
         keep existing project instructions

[Review setup changes]  [Apply setup]  [Cancel]
```

"Detected" means the executable or configuration was actually inspected.
It does not mean logged in, authorized, affordable, or capable of every
role. Details list missing capabilities with a remedy. Setup can finish
while execution is unavailable; Start explains what remains.

No guessed dollar amount is installed as a spend authorization. The user
chooses effective limits appropriate to the harness and sees what those
limits do and do not bound. Choosing a repository check does not execute
untrusted project commands during discovery.

Delivery offers local integration and the supported hosted path. A
local-only user never has to configure a remote or hosting account. Setup
also exposes approval authority and the optional configured LLM judge;
enabling a provider does not itself delegate judgment or integration.

## 2. Shape and authorize work

"Shape an idea" opens a supported interactive harness with the project brief
and relevant instructions. On return, Grove shows the resulting proposal,
questions and consequential decisions. It does not silently assign them.

```text
Hide completed tasks from the default list                 Proposed

Outcome
  The default list shows current tasks.
  --all also shows completed tasks.

Acceptance
  Existing filtering and ordering still work.
  Tests cover default and --all behavior.

[Start work]  [Discuss]  [Edit]  [Details]
```

The developer can discuss an unresolved product choice before Start.
A question discovered during execution appears through the same Answer
action. Routine technical decisions do not acquire a human approval gate.

Start previews the actual assignment, using project defaults:

```text
Start: Hide completed tasks from the default list

Execution       Claude Code / configured model and effort
Review          Codex / configured review model and effort
Boundary        Return a reviewed result for my judgment
Checks          go test ./...
Limits          Execution: configured monetary cap
                Review: configured time limit; cost not capped
Delivery        Local → main; requires my approval

[Start]  [Change settings]  [Cancel]
```

This example describes different kinds of limits rather than promising a
particular harness version supports them. Real UI uses the capability probe.

For several selected items, Start exposes the delivery boundary:

~~~text
Start selected work: A → B → C

Deliver         Separately (default)              [Together]
Progression     Continue selected work after each delivery
Approval        Configured policy + LLM judge
Limits          One aggregate allowance for this selection
Workspace       Clean up after delivery           [Keep]

Separate: A is reviewed and delivered before B starts.
Together: all selected items must finish before one combined delivery.

[Start]  [Change settings]  [Cancel]
~~~

Mode and actual authority are explicit: without delegated approval, the
preview says it waits for the owner at each delivery. Both modes implement
members sequentially; parallel fan-out is outside this milestone. Existing
assignments keep the boundary they were launched with. Preparation explains
when committed proposal records will be admitted from another source and
reports input conflicts before spending, without asking the user to choose
execution worktree paths.

## 3. Leave and return

The project view emphasizes attention and outcomes. The existing board may
remain the main layout; its cards and work detail use these same facts.
"Needs you" and "Running" are filters over lifecycle plus execution facts,
not independently editable statuses.

```text
tasks                                      1 needs you · 1 running

NEEDS YOU
Hide completed tasks
  Decision: should --all include archived tasks?
  Completed: default filtering and its tests
  Next: answer to continue                       [Open]

RUNNING
Explain invalid configuration
  Implementing · last activity 12 seconds ago     [Open]

RECENT RESULTS
Fix empty task export
  Delivered to main · verification passed         [Open]
```

A result summary is backed by records, checks and repository state.
If a provider does not emit usable progress, show the last observed
activity and that limitation. Elapsed silence is not declared a failure.

A local "since last viewed" cursor may control which changes the catch-up
view highlights; losing it merely shows more history. It never owns work
state. Raw events, costs, worktrees and revisions are a detail view away.

## 4. Answer and continue

```text
Hide completed tasks                         Waiting for your answer

Should --all include archived tasks?
  Recommendation: yes, because --all promises the full list.
  Affects: --all behavior and its acceptance tests.
  Completed: default filtering; default-filter tests passed.

[Answer]  [Discuss]  [Inspect changes]
```

The answer is durably recorded with attribution. The pending assignment's
scope does not expand automatically. A consequential answer can create a
decision and a revised plan; changed scope is brought back to the owner.

After answering, the primary action is Continue with the existing settings,
remaining authority/limits and checkpoint. Completed deliveries stay done;
continuation never resets the selection's budget or restarts those items.
The user may change harness/model/effort at this boundary.
It is a supported continuation, not a raw "restart the prompt" action.

```text
Continue: Hide completed tasks

Keep            Completed filtering and tests
Use             Answer: --all includes archived tasks
Next            Finish --all, verify, obtain independent review
Execution       Codex                         [change]

[Continue]  [Inspect handoff]  [Cancel]
```

The new harness receives current project context and the durable handoff.
Native resume is an optional shortcut within the same compatible provider
and checkout. Fresh sessions and fresh clones do not require it.

Stop ends the current execution while preserving its work. A checkpoint
missing after a crash is shown as interrupted with unverified partial
changes; Grove does not pretend the last milestone was completed. Continue
first reconciles those changes within the user's mandate.

## 5. Move into an interactive session

Discuss opens the selected harness on the chosen work and returns to the
same refreshed view when it exits. It explains whether it is resuming a
native session or opening a fresh one from project context. A delivered
workspace kept for inspection is not offered as an implementation resume;
new or reopened work gets a fresh workspace from the target.

For an actively owned worktree, the user sees Stop and discuss, Wait, or a
supported read-only inspection path. The UI does not start a second writer.
An unavailable session offers fresh continuation; a provider change never
tries to import the other provider's private transcript.

Opening a session is not an unattended execution assignment. Any seeded
turn that spends provider resources is clear before submission.

## 6. Judge a result

```text
Hide completed tasks                              Ready for judgment

Changed
  Default list hides completed and archived tasks.
  --all includes both, preserving existing ordering.

Decisions
  --all includes archived tasks — your answer

Evidence
  Tests passed on the submitted result
  Independent review: no open findings
  Reviewed by Codex / configured model             [Open evidence]

Delivery
  Local → main
  No unresolved conflict at the last check

[Approve and deliver]  [Give feedback]  [View changes]  [More]
```

"Approve and deliver" is one user action authorizing two attributable
operations. It does not collapse their meanings. Approval is retained if
delivery subsequently fails, and the screen says which operation completed.
Approve only remains available. Hosted policy may impose additional waits.

An updated candidate replaces the primary action with Review updated result.
It explains why old approval no longer applies. A target change triggers
fresh delivery verification; conflicts require a newly reviewed candidate
rather than silently resolving into previously approved code.

A non-clean independent review presents its findings and the appropriate
feedback/continue action. Successful tests cannot hide an unmet acceptance
item or a requirement for the owner's judgment. When configured policy and
the LLM judge accept the candidate, the result names that authority and its
evidence, then proceeds to delivery. A judge wait explains which acceptance
item or missing evidence needs attention; a failed call is not approval.

## 7. Deliver and understand the outcome

Local delivery shows the resulting main commit and retained candidate
evidence. It verifies the delivery once, then ordinary views read Done from
the accepted record on the target. Cleanup runs automatically when safe,
unless Keep was selected. Additional files/commits or a live owner preserve
the workspace and explain why. This does not undo Done. A crash after target
advance is recovered by inspecting the receipt, without another code merge
or post-delivery status commit.

For separate mode the same view then shows "A delivered; starting B" from
the updated target, within the authorized selection and remaining limits.
Each delivery gets a fresh workspace. A kept completed workspace is for
inspection; repeated delivery from it is unsupported. Native sessions are
optional, and project context carries continuity after cleanup.

The hosted variant shows the PR and host requirements:

```text
Hide completed tasks                           Waiting for delivery

Grove approval      Accepted candidate — by you
Pull request        #42                          [Open]
Required checks     1 running · 2 passed
Host approval       Required
Next                Wait for checks and host approval

[Refresh]  [Inspect evidence]
```

An external merge of the prepared submission is reconciled through the same
delivery operation. Normal reading uses the accepted record on the remote
target; a closed-unmerged PR or merged label alone cannot establish Done.
The same standing is available through CLI and agent context. Verification
and optional audit results are shown separately, not claimed by every read.

```text
Hide completed tasks                                           Done

Accepted result     Delivered to origin/main
Delivery            PR #42 · resulting commit D   [Open]
Approval            Accepted by configured policy · judgment available
Target checked      origin/main at T · just refreshed
Workspace           Cleaned up · evidence retained

[Inspect evidence]  [Audit delivery]
```

There is no follow-up completion PR. The record states acceptance, and
Grove checks Git to establish delivery. Source details distinguish the raw
record from the derived result without requiring users to learn the schema.

If the host cannot be reached, show standing as of the last observed target
revision and "remote freshness unknown". If no target can be read, show
"Accepted · delivery unknown" and offer refresh. A fresh clone missing
historical audit objects still reads Done from the target's accepted record;
only an explicit evidence inspection/audit reports those objects missing.
Never present missing evidence as a successful proof or as a reason to redo
delivered work. A changed candidate or its governed requirements invalidates
old acceptance for current work.

The [contracts design](G-260930-84fnb-portable-workflow-contra.md) owns the
exact source, evidence, readiness and mutation rules. Direct Markdown reads,
fresh clones and stale worktrees are part of acceptance, including whether
an agent knows to inspect delivery instead of interpreting acceptance as
either unfinished implementation or proof of Done.

## Failures have specific next actions

| Situation | Primary explanation/action |
| --- | --- |
| Unsupported limit or role | Configure a supported setting; nothing launches |
| Missing authentication | Explain provider setup; preserve project and assignment |
| Unanswered question | Answer or discuss; unchanged wait spends nothing |
| Failed check | Read failing check and continue within the assignment |
| Lost process | Show partial work and last verified checkpoint; reconcile before continuing |
| Dirty or replaced checkout | Identify the mismatch and a supported recovery path |
| Changed candidate | Judge the updated result; old approval cannot authorize it |
| Target moved | Verify the new integration result; explain conflicts |
| Host unavailable | Show last observed state and retry inspection |
| Merge succeeded, process interrupted | Read Done from the target receipt and reconcile sequence progress; do not merge again |
| Cleanup cannot remove a workspace | Preserve it, explain why, and retain delivered standing |
| Kept delivered workspace selected for implementation | Explain that new execution starts fresh from the target |
| Missing historical audit objects | Keep normal standing; explain recovery when inspection/audit is requested |

## Proposed adoption evaluation

Use the milestone's external cohort after its declared prerequisites.
Observe people completing the same journey in their own supported projects
without live coaching. Capture unaided completion, interventions, incorrect
beliefs about authority/completion, time to recover context, and whether they
would keep using Grove. Deterministic fixtures exercise recovery states;
real bounded harness trials establish behavior beyond mocks.

The original sketches and the new product boundaries have owner acceptance;
this reconciled written revision is presented for review. Each implementing
boundary includes a small local walkthrough before dependent work builds on
it; the external cohort is not the first usability check. The actual terminal
implementation still needs the owner's judgment. A changed visual arrangement is acceptable when it preserves the journey
and improves those outcomes. Exact keys, spacing and command spelling are
implementation design details, not settled by this page.
