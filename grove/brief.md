# Grove brief

Direction reconciled 2026-09-30 (owner's local date), through
[G-260930-e8jj7](G-260930-e8jj7-build-a-portable-workflo.md) and
[G-260930-tcc9w](G-260930-tcc9w-bound-delivery-groups-an.md). This is the single current source of product intent.
Work records own acceptance, progress and next actions; the
[record model](../docs/record-model.md) owns the implemented contract.
Selected direction below includes capabilities still proposed as work.

## Purpose and audience

Grove is a portable workflow product for developing software with agents.
A person shapes intent, delegates bounded work, and returns to an
understandable result they can judge and deliver. Project knowledge,
execution evidence and durable handoffs let the next person or supported
harness continue with the right context, without reconstructing the chat.

The first audience is developers like the owner working across several
existing projects, initially through a small external preview. The owner
selected building a serious product alongside their job on 2026-09-29.
Future small-team/internal SDLC use informs explicit authority, configurable
checks and shared evidence; it is not a commitment to an enterprise platform.

The next product milestone is [G-260930-60c3d](G-260930-60c3d-complete-the-portable-gr.md): independent adoption of a
complete, understandable workflow across two harnesses, local delivery and
one hosted PR path. It owns the acceptance journey and trial evidence.
Completing component features or demonstrating the loop with its author
present is insufficient evidence of adoption.

Grove owns the work and its continuity across tools. Harnesses execute
agent activities; a git host can supply code hosting, collaboration, checks
and merge policy. Core local use needs neither a host nor hosted CI.
Operating a git-host service, a public plugin ecosystem, a universal
workflow builder, schedules and fleet management are outside this milestone.
Pricing, licensing, launch/support commitments and a business model are
not selected by this direction.

## Selected foundations

- Go CLI/TUI, ordinary Markdown with YAML frontmatter, grove.yaml,
  configurable record storage and Git. Core inspection needs no service.
- Stable identity and placement: neutral date-form IDs issued without
  coordination, flat creation and recursive discovery. Folders do not
  determine validity or meaning. Classification changes preserve paths.
- Humans and agents use the same records. Deterministic validation,
  retrieval, mutation and lifecycle mechanics belong in software; judgment
  belongs in instructions and attributable human or delegated decisions.
- Branch-local editing with a combined project view. Mutations bind to exact
  sources and revisions. Normal use should not make source selection the
  person's main job.
- Bare grove opens the TUI; explicit commands remain noninteractive. Keep
  the Charm-based terminal experience keyboard-friendly and visually clear.
- Grove-owned shared workflows shipped with the binary and thin harness
  adapters. Configure
  harness/model/effort, verification, authority and delivery while retaining
  coherent meanings for work, review, approval and completion.

The owner permits substantial architectural refactoring and targeted
replacement. Fresh product requirements govern the experience; an
end-to-end trial tests the implementation boundaries. Existing code and
tests are evidence to assess, not blanket constraints on the next design.
A whole-codebase rewrite is not selected. Preserve existing project
identity, knowledge and evidence through a deliberate migration wherever
contracts change; the current implementation remains authoritative until
replacement behavior is delivered.

## Knowledge, work and evidence

| Information | Owner and purpose |
| --- | --- |
| Brief | Purpose, constraints and selected product direction |
| Terms | Domain vocabulary and relationships |
| Decisions | Consequential choices, authority, alternatives and reconsideration |
| Questions | Unresolved human choices, affected work and who can answer |
| Work | Outcome, bounds, acceptance, dependencies and next action |
| Artifacts | Plans, implementation reports and reviews tied to work and revisions |
| Attempts | Particular executions, inputs, owner, progress, result and recovery |

General knowledge needs no predefined semantic category. Terms and linked
artifacts remain discoverable. Small work may carry preparation in its work
record; substantial work gets a linked plan. A new operational type needs
distinct behavior; a general schema-extension engine is not selected.
Completed records stay put and views bound everyday clutter.

Retrieve context by activity: shaping uses purpose and relevant knowledge;
preparation uses assigned work and affected interfaces; implementation uses
its mandate; review uses acceptance and the candidate; continuation uses
the checkpoint and changed inputs. Links support discovery without loading
everything. Context is facts, never authorization.

Durable work, consequential decisions, questions, checkpoints, reviews and
delivery evidence travel with the project. Raw provider logs and native
session handles can remain local. A supported harness change promises
continuation from project context; translating private conversations is
not required. Inspect actual repository state before trusting a checkpoint.

## Workflow and authority

The workflow still moves from Proposed through Active and Review to Done,
with Abandoned an explicit human choice. Under the selected redesigned
record contract, acceptance is recorded distinctly and Done is derived from
that acceptance plus its delivery. The stored facts must remain
truthful before and after integration; a permanent review field with a
contradictory computed Done label is insufficient. Preparation,
implementation and independent review are activities; process state,
failure and waiting are additional facts. A finished attempt does not itself
establish review, acceptance or Done.

Shaping produces proposals. An assignment or Implement action authorizes
bounded execution; creating a proposal, inspecting it or opening a session
does not. Routine technical choices proceed within the mandate.
Consequential missing human choices become durable questions, and unchanged
waits do not spend another attempt. Separate plan approval at every stage
is not a universal gate.

Independent review examines the candidate against acceptance and evidence.
Human judgment leads with changed behavior, decisions, outstanding findings
and verification, with diffs available. Approval names the actual candidate;
a changed candidate needs reconsideration. Feedback preserves earlier
reviews and returns implementation work to Active.

For implementation, Done means the accepted result reached the configured
target. The owner selected
[G-260930-2qa4a](G-260930-2qa4a-derive-done-from-recorde.md): derive Done with
a redesigned record contract, without mandatory post-delivery completion
commits or PRs. Raw records communicate acceptance; Git/Grove establishes
delivery. Every lifecycle consumer uses that same result, and a target that
cannot be read cannot mean success. As
[G-260930-gj9d7](G-260930-gj9d7-prove-delivery-once-at-i.md) amends it,
reading takes the target's own copy of the record as that result, at a
cost that does not grow with history or branches, and delivery is proved
once, when it is made, and again only by an audit someone runs. Squash and
external integration are proved against retained evidence; status, an ID or
an unproved trailer is insufficient.

Schema 4 now implements recorded acceptance, cheap derived completion and
local squash delivery with retained optional audit evidence. Research/design
uses an accepted Git artifact; older done records retain their historical
meaning without fabricated approval or candidate evidence. Missing audit
objects in a fresh clone do not reopen delivered work.

An explicit standing policy may delegate bounded conflict resolution,
approval and integration under written conditions
([G-260925-wh9ax](G-260925-wh9ax-delegate-conflict-resolu.md)).
The accepted option of bounded delegated LLM approval
([G-260928-d8py6](G-260928-d8py6-a-standing-policy-may-de.md)) remains.
Such acts are attributed to policy and evidence; everything outside its
authority awaits human judgment. The owner brought the LLM approval judge
inside the milestone in G-260930-tcc9w so unattended progression is evaluated
as part of the product. Parallel launch remains outside it.

A selected set defaults to separate deliveries, with an explicit option to
deliver together. Separate work is reviewed, approved and delivered before
dependent work begins from the updated target. Together work shares one
candidate and delivery, with acceptance judged per member. One authorized
launch can progress through its selected deliveries under policy and
aggregate limits without manual relaunch between them. It adds no work or
authority and pauses when the mandate or a necessary human decision requires.

Each delivery starts in a fresh target-based workspace with external
prerequisites delivered. Proposal records can be admitted without inheriting
unmerged implementation. Interrupted work and review fixes resume before
delivery; afterward the workspace's execution role ends. Repeated delivery
from an old squashed branch is unsupported. Cleanup is automatic when safe,
with an option to keep the workspace for inspection and with evidence and
additional work preserved. Context continuity relies on durable project
facts rather than a permanent branch or native session.

Interactive and unattended use share the workflow while exposing human
availability, authority and resource bounds. Unattended shaping publishes
reviewable proposals and cannot authorize its own implementation or merge.
No fresh process per activity is universally required; independence and
durable continuation determine the necessary boundaries.

## Experience and portability

A person returning to Grove should see the intended outcome, what happened,
what needs their attention and what can happen next. Branches, worktrees,
attempts, source revisions and native sessions remain inspectable details.
Onboarding establishes project instructions, harness, checks, limits and
delivery method, preserving existing content and explaining missing
capabilities. It must not silently grant execution or integration authority.

CLI and TUI use the same work facts and operations. Reviews and catch-up
summaries expose their supporting evidence and uncertainty. Inspecting,
reconnecting and unchanged waits do not launch work.

Evidence-backed debriefs may propose improvements to knowledge, guidance or
mechanical checks. They do not silently change acceptance or turn every
observation into permanent policy.

The current project-wide view remains based on Git ancestry and merge bases:
timestamps, status ranks and newest branch tips do not choose authority.
Collapse identical observations, retain genuine divergence and label
unintegrated/uncommitted changes. The integration target labels delivery;
it does not choose which record state is current. Changes to delivery representation must
preserve accurate current-state interpretation. The implementation details
belong in the record model and design, not a parallel UI status store.

Claude Code and Codex are the milestone harnesses
([G-260928-vdhf0](G-260928-vdhf0-run-attempts-on-codex-as.md)).
Roles include shaping/preparation, implementation and independent review;
select suitable strength and record actual configuration. Distinguish real
capabilities, permissions, effective limits and unreported usage. A time
limit is not a monetary ceiling. Native resume is an optional capability;
portable continuation relies on durable work.

An on-demand Grove owner process remains the execution starting point
([G-260923-tnn5e](G-260923-tnn5e-run-attempts-as-a-grove.md)), with
duplicate-start refusal, explicit stop and owner-loss handling.
There is no selected machine-reboot survival guarantee. New host or harness
support must be exercised before claimed; interfaces emerge from supported
implementations rather than an untested universal plugin framework.

## Direction, work and history

This brief states no order among work records. Their depends_on edges,
shown by grove deps and the board, are authoritative
([G-260927-y3pc2](G-260927-y3pc2-order-is-an-edge.md)). The milestone owns
membership and acceptance, and each work record owns its Next.

The earlier interactive adoption milestone
[G-260921-407n6](G-260921-407n6-complete-the-interactive.md) was accepted
on 2026-09-22 using this repository's loop in place of its unexercised
nullsec acceptance. Nullsec's cutover was included; keyborg was selected as
another adoption test without its own work record at that time.
The 2026-09-23 preview work on evaluations, Attempts and distribution is
historical preparation, not evidence of the new milestone's completion.

At main 18ba791, inspection showed schema 4, a Go CLI/TUI, Git-backed records,
shared guides, a Grove-owned Claude attempt runner, candidate-bound judgment,
local squash integration, cheap standing, optional delivery audit and
deterministic policy sweeps. Multi-item runs still deliver together and
cleanup is opt-in. The new selection/workspace boundaries, automatic
progression, second provider, LLM judge, hosted path and assembled experience
remain proposed work. Small local trials belong with the implementing
boundaries, ahead of their dependent work; a final polished experience is
not the first test of the architecture.

The predecessor is uninstalled. The archived application at
grove-archive-2026-09-18, last commit be40e46, remains historical, with no
service, credential, deployment or backlog authority over this project.
Sibling repositories are evidence; changes there need their own assignment.
