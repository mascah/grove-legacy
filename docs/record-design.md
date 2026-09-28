# Record design

Why the record model is as it is. The contract itself, every key, type,
field, value, status and refusal, is [the record model](record-model.md),
which `grove guide model` prints; command behaviour beyond `grove --help` is
the [command reference](commands.md). This document is not shipped, so it
may cite this repository's records; it states no rule the model does not.

## One schema, no history

Grove keeps no backward compatibility before its first release
([G-260921-ebsby](../grove/G-260921-ebsby-decouple-record-identity.md),
[G-260921-r491p](../grove/G-260921-r491p-reconcile-all-grove-cont.md)).
Schemas 1 and 2, their type folders, typed `W-`/`Q-`/`D-`/`T-`/`P-`/`R-`
IDs and per-type counters were deleted by the conversion to schema 3, whose
old-to-new mapping [G-260921-czt8x](../grove/G-260921-czt8x-identity-and-path-migrat.md)
keeps. A commit from before it is inspected with the CLI built from that
commit; the current CLI reports such a branch in `versions` and the board as
a source it cannot inspect. The files earlier CLIs kept in the Git common
directory, `grove/neutral-ids`, `grove/lock` and `grove/next-ids`, are never
read or written. `schema_version` is this CLI's own number, unrelated to any
other tool's.

## Identity apart from classification

An ID names a record, not what kind of record it is, so reclassifying never
renames and records created in separate clones need no reconciling
([G-260926-2da4n](../grove/G-260926-2da4n-identify-records-by-crea.md)). The full ID is the
canonical identity, not an alias for a hidden random value, so `show` needs
no abbreviated lookup. Its date is the day it was issued and proves nothing
about creation time, priority or execution order; the board orders cards by
their dates.

Only `type` classifies. A missing or unknown type is an error rather than a
page, so a damaged operational record never degrades into valid general
knowledge, and the six operational types keep all their rules wherever they
sit. A page gains nothing from its folder or prose: it is never a work card
and no gate can target it.

Reclassifying grants nothing: an `accepted` decision stays a claim in a
file, as it always was. Because the whole project must validate, work that a
plan or question names cannot stop being work.

`formerly` exists so that a merge from a branch that predates the conversion
cannot quietly restore a second owner of one identity: a restored typed-ID
record fails `check` on its ID, and a second record with one `formerly`
fails on that. `convert` provides no lookup by former ID, batch conversion,
link rewriting, or per-type folders or prefixes; the caller keeps the
mapping `convert` prints and repairs links from it.

## Issuing IDs without coordination

A tail is drawn rather than counted so that no state is shared between
clones. The lock lives in the directory
`git rev-parse --path-format=absolute --git-common-dir` returns, never a
worktree's own `.git` path, since linked worktrees keep private metadata as
well as a shared common directory
([Git's worktree documentation](https://git-scm.com/docs/git-worktree#_details)).
It is an `flock`, which the kernel releases when the holder exits, so it
needs no daemon and no cleanup; it is local coordination state, not a
tracked record, the read-only commands never create it, and it is never
unlinked.

Separate clones share nothing, so the one collision left is two clones
drawing one tail on one day. After a merge `check` reports it as two files
with one ID, an error even when their contents match, since matching IDs
alone cannot prove two independently created records are one item. Genuine
branch copies of one record keep their identity, and are versions to
reconcile, not collisions. A hand-authored ID bypasses the draw, and
creation outside Git is not provided.

## Dates and revisions

Direct editors should update `updated` when they change content, but that is
an authoring convention, not a provable freshness guarantee. Readers never
repair dates or infer them from filenames, modification times or Git, and a
missing creation date is never invented. Safe writes therefore compare
content, not timestamps: `--expect` is the revision a caller read. It is
optional because a person at a shell reads and writes within seconds under
the same write lock, while an agent session's read may be old.

## Files and folders

Filenames stay short even when titles are long, and a title edit does not
extend or rename them: relationships use IDs and survive renaming, while
ordinary Markdown path links need updating when a file moves. Folders
organize without meaning, so closed records stay where they are and
discoverable. A `.md` file that is not a valid record is reported rather
than dropped, so nothing silently falls out of the project.

## Planning metadata

The planning fields are optional so that quick capture stays useful, and an
absent size or priority is not silently turned into an estimate or an
urgency decision. They drive filtering, grouping and presentation, never an
automatic execution policy beyond the handoff shape `size` selects.

- Membership describes decomposition and dependencies order work; the two
  edge types are checked for cycles separately and never combined into one
  precedence graph, so a group may depend on delivery of its own members
  without making membership a prerequisite of each child.
- Grouped work keeps its own outcome and acceptance: completed children do
  not establish the parent's completion or integration, and unfinished
  members can keep a parent from done without blocking work on it.
- Member counts and blocker explanations derive from relationships in the
  selected branch context, never a stored percentage or `blocked` flag.
- An investigation is a kind with an independent size; there is no `spike`
  value, and a size sets no preparation depth; `small` only selects the
  compact handoff the work guide (`grove guide work`) describes.

## Knowledge records and the brief

A term's body gives meaning, relationships and boundaries, not execution
instructions or implementation state; the shaping guide says where those
belong. A question keeps its identity when resolved, the answer in its body
or a linked decision. A superseded decision is an accepted one a later one
replaced, naming its replacement in `relates_to` rather than a dedicated
field; the body says why, and prior versions remain in Git.

Work does not name its plans or reviews: that side is derived, so one plan
can serve several items and nothing has to be kept in step. A review record
is evidence, never approval, and there is no report type, since the work
record's Evidence is the report. Whether reviewed content changed since a
review is a comparison a reader makes between `examined` and the work's
`candidate`, not stored state.

The brief is one file rather than a record because it is the direction, not
an item of work.

## The work lifecycle

Preparation, implementation, independent review and waiting are activities
inside `active`, recorded in the body, never statuses. A failed or
interrupted attempt does not enter Review; it stays active with a
checkpoint. Review means a candidate awaits human judgment, the owner's own
or given in advance as a standing policy
([G-260925-wh9ax](../grove/G-260925-wh9ax-delegate-conflict-resolu.md)), and
the record's Evidence and Next carry the handoff, so a new session can judge
it without the originating chat.

`candidate` names the last implementation commit, and the commit that sets
`review` changes only records, so `git diff --stat CANDIDATE TIP` shows the
handoff alone. A changed candidate is a new value, the earlier one in Git
history, so the reviewed and the approved commit can each be compared to a
review's `examined`. Approval binds one commit
([G-260921-btyck](../grove/G-260921-btyck-approval.md)): a
changed candidate cannot inherit it. It stays per record, so each member of
a group is judged against its own acceptance. Feedback keeps the candidate
so earlier reviews still compare to it, and reopens the whole group, since
the next candidate replaces the shared one. A resolution is a merge, so the
earlier candidate stays an ancestor of the new one.

Done means accepted and merged into the target, for research and design
deliverables too, since those are files
([G-260921-9wkjt](../grove/G-260921-9wkjt-hand-implementation-cand.md)). The
ancestry check keeps a checkout without the code from closing the work, and
the target check exists because on the work branch the candidate is an
ancestor too ([G-260921-jwk4e](../grove/G-260921-jwk4e-review-candidates-and-in.md)).
A done record without a candidate asserts only that its outcome was reached
in its own branch context, as its Evidence says; it is not proof of a merge,
and its prerequisites' delivery is established by Git ancestry or observed
behaviour.

What the model leaves unenforced are rules about people, which software
cannot verify, and the owner edits by hand. The guides carry them, and the
guide's `git diff --stat` check and `approve`'s refusal of a changed tip
catch a reopened record whose candidate was not moved.

Body organization is for readers. The CLI infers no readiness or completion
from headings, populated prose or checked boxes, and directory movement
determines nothing.

## Attempts are not records

An attempt is a process's inputs, raw events and result, not a claim about
the work, so it lives beside the repository, shared by every worktree, and
the work record's own status on the attempt's branch is the only handoff.
The board derives an outcome from the files for display and writes nothing
about an attempt.
