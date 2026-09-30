# Command reference

`grove --help` gives every command's usage. This document owns what it does
not say about the commands below. The contract they read and write, the
configuration, types, fields, statuses and what `check` and `update` refuse,
is [the record model](record-model.md) (`grove guide model`), and why it is
so is [the record design](record-design.md). [The board](board.md) has its
own document. Each command's acceptance, evidence and limits belong to the
work record that delivered it.

## Records

`list`, `show`, `brief` and `check` read the selected checkout's live files,
uncommitted records included, name the project on stderr, and change
nothing. They load the whole record set before resolving relationships; a
concurrent direct edit can invalidate a read, since the reader promises no
transactional snapshot. `--project DIR` can come before or after the
command, and `--help` needs no project.

- `list` prints ID, type, status, standing and title, ordered by `created` with
  undated records last, then by ID; the order implies no urgency. One-line
  output escapes control characters, so a multiline title cannot break the
  table. `--status VALUE`, repeatable, keeps records whose status equals any
  given value; a value that is no type's status, or an empty one, is a usage
  error that lists the statuses, and a status no record holds prints the
  header alone. STANDING is work's [derived standing](#standing), `-` for
  other types.
- `show ID` writes the file's original bytes to stdout and the project and
  file to stderr, with work's standing there too; `--json` prints
  `{id, path, revision, source}`, plus `standing` for work and
  `approved_by` (`owner` or `policy`) while the record is accepted.
- `new TYPE TITLE [--slug SLUG]` refuses an unknown type before anything
  else, writes the record with a body skeleton, prints the root-relative
  path, and fails without deleting the file if the project no longer
  validates.
- `update` prints `{id, path, revision, changed}`. A failure before the file
  is replaced (a refused request, a stale revision, an invalid result)
  leaves its bytes unchanged; once it is replaced, a later failure
  (directory sync, final validation, output) is reported as an applied
  update. `--commit`, after a change, commits the record's file alone with a
  generated message and adds `commit` (`null` when nothing changed); a
  commit Git refuses is reported as an applied, uncommitted update. No
  update makes a record `done`, which schema 4 derives.
- `convert` prints one JSON line, `{from, from_path, id, path}`, the caller's
  durable old-to-new mapping. A rerun neither duplicates a record nor remaps
  an identity, since a converted or missing source is refused before any ID
  is drawn, and an existing target file refuses the run. It never rewrites
  body prose, Markdown links (the moved file's own included), `examined`, or
  anything outside the record root, and never modifies or removes the
  original; set `work`, `status` or `examined` afterwards with `update`.

## Judging and integrating

`approve ID VERDICT` records the acceptance: `status: accepted`,
`approved`, `approved_by: owner` (the sweep writes `policy sha256:…`) and
`approved_context`, and appends `Verdict on candidate X, DATE: VERDICT` as
the body's last paragraph. Authority is who ran the command, never what the
verdict says. `feedback ID TEXT` appends `Feedback on candidate X,
DATE: TEXT`, and each other member of the group it reopens gets `Reopened
with ID's feedback on candidate X, DATE`; nothing earlier is removed.
`resolve ID` is that feedback, generated, naming the target commit and the
conflicting files.

`integrate ID` runs in the target's clean checkout and delivers the one
branch holding the work accepted, its acceptance applicable, as one squash
commit ([G-260929-gm3m4](../grove/G-260929-gm3m4-clean-main-history-with.md)).
It reads the branch tip once, the submitted commit S, and refuses before
anything changes: a branch lacking the candidate, a later commit changing
more than the records, a shared candidate whose members are not all
accepted, a squash that would carry the candidate of other unfinished work
on the branch, such as a member reopened by feedback and not handed off
again, and a conflict, which `git merge-tree` predicts in objects only and
names with the next action. Then it

1. retains S as `refs/grove/submitted/S`, the evidence, before the target
   moves;
2. writes S merged onto the target tip B as a commit whose only parent is
   B, with a Conventional Commit message: the most significant type among
   the candidate's own commits, `!` if any is breaking, the first member's
   title, the members, and the trailers `Grove-Work: ID` (one per member),
   `Grove-Candidate: C` and `Grove-Submitted: S`;
3. fast-forwards the target to it, so a target that moved meanwhile is
   refused, never overwritten;
4. reads the standing again and reports each member done only once it
   verifies.

It writes no record: the delivered record is the branch's, accepted.
Rerun after an interruption, it finds the delivery through the same
verifier and goes on to cleanup, never a second commit. `--cleanup` removes
the worktree only where Git agrees and it holds no ignored files, and the
branch only while its tip is still S. The retained ref keeps S from
garbage collection and is local: to take the evidence to another clone,
push or fetch `refs/grove/*` with the branches. A clone without S reads
the delivery as unknown, never as done.

### Standing

Done is derived, never written. A work record's standing is `proposed`,
`active`, `review` (awaiting judgment, or an acceptance that no longer
applies, since the title or body outside `## Next` and Grove's own
paragraphs changed), `accepted` (applicable, not delivered), `done`,
`abandoned`, or `unknown` with a reason. Accepted work is done when the
configured target contains its candidate, or contains a commit whose
trailers name it and the retained submission, whose parent is on the
target, and whose tree is exactly the submission merged onto that parent.
A delivery commit that is forged or altered is not done, and one whose
submitted tip this repository lacks is unknown; nor is one the target no longer contains after a
rewrite, while a later revert does not undo it. Every consumer reads this
one verifier: `list`, `show`, `deps`, `context`, `run`, the sweep and the
board. A `done` record is schema 3's claim, kept by migration and labelled
so, never as verified.

## Versions

`versions [ID] [--json]` shows one row per version of each record: its
committed version on every local branch tip and its live version in every
registered worktree, grouped by ID in [ID order](record-model.md#identity-and-dates)
with each source's own status, so main can
see a feature branch's progress without switching or merging. Live rows say
how the file compares with that checkout's HEAD (`unchanged`, `modified`,
`renamed`, `added`, `deleted`). Every row ends with a selector that binds the
repository, source, commit, configuration, record path, and content revision;
identical bytes in two sources get two selectors. A committed source's
`grove.yaml` is checked for the form of `brief` only, never that the file
exists, since only a live checkout must hold the brief. Sources are listed on
stderr with any diagnostics; an invalid or unreadable source makes the result
incomplete and the exit code 1 while valid sources still print. A live
project must be inside its registered checkout: a project location reached
through a symlink, or belonging to another nested repository, is invalid,
while a missing one is simply absent.

The `CURRENT` column says whether a version is `yes`, current, or `older`.
A version is older when another has different bytes and the record at their
common commit has the older one's bytes, so only the other side changed it
since they split. The common commit is the two commits' merge base, or a
checkout's HEAD for its own uncommitted edit. A revert is a change like any
other. Several current contents are a divergence, shown as they are. A pair
whose common commit cannot be read, or whose several common commits hold
different versions, stays unordered: both current, with a note on stderr.
Reverts carried across merges can make versions older than each other in a
cycle; then a version is current when everything newer than it, through any
chain, is also older than it, and a note says so. A branch whose current
state removes the record gets a committed `deleted` row. Dates, status order,
and branch names never decide, so the answer is the same from every checkout,
and no branch is special. [G-260921-ms6ev](../grove/G-260921-ms6ev-derive-a-project-wide-cu.md) owns this.

With an integration target named in `grove.yaml` (`target: main` here),
stderr names it and the `TARGET` column says whether that branch holds the
version's bytes (`yes`, `no`; `-` without a target). The target labels only
and decides nothing. `--json` adds each version's exact source text,
`current`, `older` (why it is older), and `on_target` (null without a
target), each record's `notes`, and the result's `target` and target
`notes`. The command writes nothing: no refs, index, worktrees, records, or
coordination state.

## Workspace

`workspace --source SELECTOR [--json]` takes one selector from `versions`,
checks that the version is still exactly what was selected, and prints the
absolute project directory of the existing checkout holding it, ready for
`--project`:

```sh
grove --project "$(grove workspace --source SELECTOR)" show G-260925-7k2qm
```

A live selection resolves to its worktree. A committed selection resolves to
the one worktree that has that branch checked out, provided its live copy of
the record still has the committed bytes; otherwise, or when no or several
worktrees hold the branch, it refuses and says to run `versions` again and
select a live version. Any change since selection (branch or HEAD moved,
worktree moved or removed, attached or detached state, configuration, record
path, or content) refuses with the reason, and the selected checkout is
re-read once more just before its path is returned. The command never creates
a worktree, switches a branch, launches anything, claims ownership, or edits a
record; the directory it prints is a location, not write authority, and a
later `update` performs its own revision check. The board uses `versions` and
`workspace` in process, and selecting a version never creates a worktree for
a branch without one.

## Context

`context WORK_ID... [--json] [--interaction interactive|headless]
[--max-bytes N] [--include PATH]...` assembles context for explicitly selected
work in one checkout, so an assignment is a few IDs rather than a composed
prompt. [The work guide](work-execution.md) is the workflow that uses it. It
prints the IDs ordered prerequisites first (ties keep the requested order;
prerequisites are never added to the selection), Git identity, and two kinds
of content that it keeps apart:

- **Sources, read in full**, each once with its `sha256:` revision:
  `grove.yaml`, the selected work records, and every `--include PATH`. That is
  the whole default. An include is a required project-relative file of any
  type; one file reached by several names (letter case, a hard link) is
  included and charged once. A page cannot be selected; a related page is
  listed with status `-` and read only through `--include PATH`.
- **Observations, listed and not read**: the transitive `depends_on`
  prerequisites, questions blocking the selected work or a prerequisite
  (resolved ones too), plans and reviews whose `work` names a selected ID, and
  the records the selected work names in `relates_to` or `members` or links
  to, each with title, status, path, revision, why it is listed, and whether
  its source is included (`source` is the `sources[]` path that holds it,
  which differs from `path` when it was included under another name). Every
  link in the selected records' bodies is listed with the project path it
  resolves to. Links are taken from parsed Markdown, so code, fences, images,
  and HTML are never links; Markdown escapes and entities are decoded before
  the URL is. A link is never opened, so a listed path is not checked, not
  even for existence. URLs, fragment-only links, other projects, absolute
  paths, and Git metadata are listed with that reason and no path.

Retrieval is staged with existing commands: `show ID` prints a listed record,
and `--include PATH` adds a listed file, such as the plan the record names as
current, with its revision. The command picks no plan or review by itself, and
it never follows links inside an included document. A listed record's revision
says which version was seen; a listing is not a reading of its constraints.

It refuses rather than degrade: unknown, duplicate, or non-work IDs, a
malformed link destination, an included file that is missing, a symlink, not
regular, or not UTF-8, a source that does not fit `--max-bytes` (default
262144, counting source bytes, at most 8388608), or a checkout that changed
while it was read. Nothing is truncated or summarized to fit. A refusal
prints nothing to stdout. Exit 0 means context was assembled, never that work
is ready, authorized, or integrated: a prerequisite's `done` is its recorded
status, and only one with a `candidate` claims a merge, where it was written.
`--interaction` records whether a person can answer (default `interactive`);
it is the caller's declaration, passed through for the guide to act on. The
command writes nothing and starts nothing, and needs Git only inside a
repository. Text output (format 3) escapes terminal controls, fences each
source under one line with its path, revision and why it is included, prints
an included record's path and revision only there unless it was included
under another spelling, shortens each link's reason to a word (`listed`,
`included`, `included as SOURCE` for a file included under another
spelling, `external`, `fragment`, `absolute`, `outside`, `Git metadata`)
and the scope notice to two sentences; `--json`
has the exact source strings, the full reasons and notice (`format_version`
3, whose JSON is version 2's; version 1 read prerequisites, related records,
and linked documents in full):

```text
{format_version, root, interaction, selected[], order[],
 git: null | {checkout, common_dir, ref, head},
 records[]: {id, path, type, title, status, revision, roles[], selected, included, source},
 requirements[]: {work, prerequisite, status, selected},
 questions[]: {id, status, blocks[]},
 sources[]: {path, revision, reasons[], content},
 references[]: {from, target, path, reason},
 scope_notice, source_bytes, max_bytes}
```

## Dependencies

`deps [WORK_ID...] [--json]` shows how work depends on other work, from one
checkout's records, so the owner can choose what to do next or assign
together without asking an agent to reconstruct the order. Only `depends_on`
orders: `members`, `priority` and question `blocks` never do, and an absent
`depends_on` means no declared prerequisite, not readiness.

Without IDs it lists the unfinished work (proposed, active, review), one row
each with `GROUP`, `LAYER`, `NEEDS` (its `depends_on`), `UNLOCKS` (the
unfinished work that names it) and `DELIVERY` (below). Work in one group needs other work in it,
directly or through anything else; a separate group is unrelated. A layer is
one more than the deepest work of its group that it needs, so layer 0 needs
no unfinished work. Equal layers list in ID order and have no declared
order, which is not
evidence that the work can proceed in parallel. The prerequisites of that
work which are not unfinished (done, abandoned) are listed after it.

With IDs it previews that selection: the IDs as given, their order as
[`context`](#context) prints it, every transitive prerequisite outside the
selection with the selected work that needs it, listed and never added, the
open questions blocking any of them, and the pairs with no declared order.
Unknown, duplicate and non-work IDs are refused as `context` refuses them.
The preview adds no work, starts nothing and authorizes nothing.

Each item's delivery is a fact, not its status alone: proposed and active
work awaits implementation; a candidate in review or schema 3's done says
whether HEAD of this checkout and the target branch contain it, by Git
ancestry, or that Git cannot read it here; accepted work gives its
[standing](#standing), with the delivering commit once done; done without a
candidate says its delivery is unrecorded; abandoned work will not be delivered, and a note names the work
that still needs it. Done work whose candidate HEAD lacks, where HEAD holds
exactly one commit with the same patch, a rewritten copy as after a rebase,
names that copy in its delivery, and a note gives the repair: in the
target's checkout, `grove update ID --set candidate=COPY --commit`, with
`--unset approved` where the record is approved, since an approval is of
one commit, then a note under its verdict. The copy counts as delivered only once the record names it. A candidate in review also says what merging it into
the target's current tip would do, performed with `git merge-tree` in
objects only: `integrated`, a fast-forward, a clean merge although the
target moved since the branch left it, or a conflict in named files, with
the target commit it read, stale once the target moves. A selection with
two or more candidates in review not on the target merges them in its
order, each onto the ones before, and a note names where the first conflict
lands and in which files; Grove states that order but never chooses one, and
a clean order is not evidence that the changes work together. Git before
2.38 cannot perform such a merge, and the delivery says the merge could not
be predicted. Notes also report what the current view (see
[Versions](#versions)) says about each listed record: an older version here,
diverging current versions, other `depends_on` elsewhere (never merged into
this checkout's), uncommitted changes here, a record that changed while it
was read, and an incomplete inspection, which prints everything, says so on
stderr and exits 1. `--json` prints `{checkout: {root, ref, head}, target,
selected, order, items[]: {id, title, status, revision, candidate, outside,
layer, group, needs, unlocks, needed_by, delivery, merge}, questions[]: {id,
title, blocks}, notes, merge_order}`, where `merge` is `{target, commit,
outcome, conflicts}` or null and `merge_order` lists `{id, target, commit,
outcome, conflicts}` up to the first conflict or is null, with `selected`
null for the overview. The command writes
nothing but the objects a merge writes, and no ref.
[G-260925-g39ga](../grove/G-260925-g39ga-see-work-dependencies-an.md) owns it, and
[G-260925-h8rj5](../grove/G-260925-h8rj5-predict-whether-a-candid.md) the merge prediction; the board's
`g` shows the same interpretation ([Dependencies](board.md#dependencies)).

## Attempts

Grove starts an agent only through `run` or the board's `R`, one assigned
work, or one explicit selection of work, per attempt. `run ID...
[--dry-run | --expect DIGEST] [--budget USD] [--permission-mode MODE]
[--until plan] [--model MODEL] [--effort LEVEL] [--branch NAME]
[--worktree DIR] [--resume]` starts one bounded implementation
attempt of proposed or active work as a Grove-owned `claude -p "/grove-work
ID... --interaction headless"` process that outlives the terminal
([G-260923-tnn5e](../grove/G-260923-tnn5e-run-attempts-as-a-grove.md),
[G-260921-h46pb](../grove/G-260921-h46pb-run-one-bounded-implemen.md)), the IDs passed as given. It creates `worktree-` plus the IDs joined by `-`
(`worktree-G-260925-7k2qm`, `worktree-G-260925-7k2qm-G-260925-8m3xd`) under
`.claude/worktrees/` from this checkout's HEAD, or reuses the branch's
registered worktree so a next attempt continues from preserved partial work,
then starts an owner process in its own session that runs the provider there
with `--output-format stream-json`, `--max-budget-usd`, `--permission-mode`
and `--permission-prompts none`, its stdout and stderr written straight to
files under the Git common directory (`.git/grove/attempts/ATTEMPT/`, shared
by every worktree, never committed). Budget and permission mode are required,
from the flags or from the launching checkout's `grove.yaml` `run:` defaults
([record model](record-model.md#configuration-and-discovery),
[G-260924-ecs9m](../grove/G-260924-ecs9m-default-an-attempt-s-bud.md)), which may also
set the model and effort; a flag overrides its default for one launch, and
`attempt.json` records the resolved values alike. Without either source `run`
is refused as a usage error: Grove itself sets no default spend or profile.
Grove starts one process and never retries; subagents the provider starts
share the budget.

Several IDs are one selection
([G-260925-7c8g9](../grove/G-260925-7c8g9-execute-an-explicitly-se.md),
[G-260925-wc2pz](../grove/G-260925-wc2pz-review-an-explicitly-sel.md)): still one
process, one worktree and one budget over all of it, never one process per
ID. `run` orders the members as [`deps`](#dependencies) does and adds
nothing: a prerequisite outside the selection is listed with its delivery at
the base (the launching HEAD, or the reused branch's tip), never implemented.
A member **waits**, and the agent does not start it, when an open question
blocks it, when an outside prerequisite is not delivered at the base
(proposed, active, review or abandoned, or done with a candidate the base
lacks; done without a candidate is reported as unrecorded delivery), or when
a selected prerequisite waits. Where the base holds exactly one rewritten
copy of such a candidate, the wait names it and the `update` that records
it, as [`deps`](#dependencies) does; it still waits until the record names
the copy. One ID is a selection of one, so work whose
prerequisite is undelivered is refused too. Bounded by `--until plan`, only
an open question stops a member: a plan needs its prerequisites named, not
delivered. The agent implements the members one at a
time in that order; a new question, an outside blocker or a failure stops
that member and every member that needs it while the rest continue, and
budget exhaustion or `stop` ends everything with what is committed kept. The
complete members enter review together on one shared candidate, judged per
member and integrated as a group ([Work lifecycle](record-model.md#work-lifecycle));
a started member left incomplete holds the whole branch out of review,
since integrating it would carry the unfinished code, while one that never
started keeps its wait and does not. `--dry-run` checks everything a launch
checks and prints the assignment without writing or starting anything: the
IDs, order, each member's status, revision and whether it can start or
waits and why, the outside prerequisites, the base, worktree, bounds, review
boundary, continuation policy, and a digest, sha256 over the IDs as given,
each member's revision, the base and the resolved options. `--expect DIGEST`
refuses a launch whose assignment no longer has that digest and prints what
it would run now. `attempt.json` records the selection with its digest and
member revisions; an attempt from before selections reads as a selection of
its one work.

`--resume`, off by default and the owner's choice for a relaunch after a
question is answered or a plan is ready, continues the previous attempt's
session instead of starting fresh ([G-260928-kehya](../grove/G-260928-kehya-resume-an-attempt-s-sess.md)).
It resumes the newest finished attempt of the same selection on the same
branch and worktree, whatever its bound, so the attempt after a `--until plan`
one continues that session. The command gains `--resume SESSION
--fork-session`, so the new attempt keeps its own session id and events and
cost, and the source's transcript stays as it ended; the prompt is unchanged,
so the resumed session runs the guide again and rereads what changed. The
source's ID is recorded as `resumed_from`, and appears in `Requested:`
(`, resuming ATTEMPT`) and in the `--dry-run` preview and its digest. It is
refused, before anything is written, when the branch has no worktree at the
one that would be used, no finished attempt of the selection ran there, or the
newest one never started its provider's session; `resolve` does not take it.
A session the provider no longer holds ends the attempt at once with the
provider's own error and no cost.

Three options shape one launch, recorded in `attempt.json` and reported as
the attempt's `Requested:` fact ([G-260924-5b6pz](../grove/G-260924-5b6pz-bound-an-attempt-at-its.md)).
`--until plan` adds the bound to the assignment (`/grove-work ID --until plan
--interaction headless`): the attempt stops at a committed plan with the
record's status as it found it and the continuation in its Next, as the work
guide's step 4 says, and launching again without the bound, after reading the
plan, is the implementation. `--model` and `--effort` pass through to the
provider, which owns their values (`claude --help`); a preparation attempt and
the implementation after it can differ in both. The launch also records the
sha256 of the worktree's `.claude/agents/grove-reviewer.md`, the reviewer
definition `init` writes and step 6 reviews through, or `none` where there is
none, which the launch warns of: without it the attempt has no independent
reviewer, and work whose record requires one stays active. It records the
sha256 of the worktree's `.claude/skills/grove-work/SKILL.md` as `skill`,
and in `differs_from_template` which of that skill and the reviewer are not
byte for byte the launching `grove`'s templates: a custom, older or newer
file keeps its own digest and never borrows the template's identity, and the
attempt's `Entrypoints:` fact says which (an attempt from before this says
they were not recorded). A skill the harness prefers over the worktree's,
such as a personal one of the same name, is not seen. `grove_version` is the
launching `grove`'s
[version](#version-and-guide) line. The `grove` the agent itself runs is not
recorded: `PATH` or the project's instructions choose it, and it need not be
the launching one. An attempt started
by hand, such as an interactive `/grove-work`, writes no attempt files: it is
visible only as its branch, its worktree and the checkpoint in the work's
Next.

`grove --help` lists what `run` refuses. The worktree holds only what is
committed, so `run` refuses a worktree without `.claude/skills/grove-work/SKILL.md`,
the skill the prompt names, which `init` writes and you commit
([G-260925-3pj9a](../grove/G-260925-3pj9a-launch-attempts-only-whe.md)); a new branch
is checked in HEAD before it is created. It likewise refuses a worktree
whose marked `grove-work` skill or reviewer definition carries an
[entrypoint revision](#entrypoint-revisions) the launching `grove` does not
serve, unrevised included, which would stop or contradict the guide after the
spend began
([G-260925-p2k54](../grove/G-260925-p2k54-keep-installed-harness-e.md)). An open
question that blocks the work is the wait the headless guide persists, so
rerunning with nothing changed refuses the same way: a selection none of
whose members can start is refused, naming each wait, whether the launching
checkout or the reused branch holds it. The launch prints the selection as
the preview does. A running or orphaned attempt whose selection shares any
member refuses the launch, so overlapping selections cannot both own a
record. After feedback reopens a group, a selection that leaves out a member
still sharing the candidate on that branch is refused: the group runs again
together, on the same branch (`--branch` names it when the IDs are given in
another order than the branch's name). Every refusal comes before a write, except
that what an existing branch holds is checked in its checkout, so a branch
that had no worktree keeps the one `run` made.

`attempts [ID]` lists attempts newest first, or those whose selection
includes ID, with every member in the WORK column. With ID it ends with
their total ([G-260927-dx0yn](../grove/G-260927-dx0yn-retain-per-attempt-proce.md)):
the attempts, the cost and turns of their result events, and the minutes
from start to finish. A selection's attempt counts in full for each member;
an unfinished attempt is in none of the sums, and one without a result
event is not in the cost or turns, each said beside the total. Turns sum
`num_turns` over a run's result events, one per query of a resumed session;
an attempt that finished before Grove summed them holds only its last
result's, and the total's turns show `≥`. `attempt ATTEMPT [--json]`
prints one attempt's launch, what it requested, event counts (parsed bounded:
a line over 1 MiB is counted, not read), the provider's init fields with the
model that actually ran, the result event's fields with its cost split by
model where the provider reports one (a subagent on another model shows
apart), the result
(exit, the worktree's HEAD and whether it holds uncommitted or untracked
changes), what its commits changed, the record as the branch holds it, whether the record on the target
changed since launch, what the owner's sweep after a handoff did ([after an
attempt](#after-an-attempt)), and the file paths; no provider text is
printed. For a
selection it prints each member at launch and, once ended, one line per member
as the branch holds it: awaiting judgment (review, with the candidate),
active with its checkpoint in its Next, not started and why, or waiting on a
question; and whether each member's record on the target changed.
`Changed:` is `git diff --name-only BASE HEAD`, recorded in `result.json`
when the attempt finishes, since its worktree and branch may be gone when it
is read, split into the files under the record root, which `grove.yaml` at
the base names, and the rest, with the time from the start to the committer
time of the first commit that touched a file outside it, following first
parents so that the commits a merge brings in are not the attempt's own;
an attempt finished
before this says it was not recorded, and one whose HEAD or `grove.yaml` Git
could not read says why. `Shape:` is derived from all of `events.jsonl` each
time the attempt is read, a line over 1 MiB skipped and counted, which makes
its counts lower bounds (`≥`); the board does not read it. It counts every
tool call, subagents' included, the Bash calls one of whose segments, split
at `;`, `&`, `|` and newlines without parsing quotes, runs `grove` (a program
of that name by any path, or `go run` of a `cmd/grove` package), and `grove
guide NAME` runs by NAME. Known misses: `grove` through a wrapper, `$(…)`,
`env` or a leading variable assignment is not seen; a heredoc's line that
starts with `grove` is; a guide read as a file is not a guide print. While
the attempt runs these are so far.
Liveness is the owner's file lock, never a pid: `running` while it is held,
`finished` once `result.json` exists, `orphaned` when the owner is gone but
the provider's process group lives, `interrupted` when nothing is left and no
result was written (a machine restart reads so; no reboot recovery is
selected).

### Resolving a conflict

`resolve ID` runs from any checkout, for work in review whose candidate
conflicts with the target
([G-260925-dz10z](../grove/G-260925-dz10z-update-a-conflicting-can.md)). It predicts the merge
as [`deps`](#dependencies) does, against the target's current tip. It
records feedback in the branch's checkout, as `feedback` would, with
generated text naming the target commit and the conflicting files. That
text is the attempt's whole mandate: merge that commit, never a rebase or a
later tip; resolve those files; rerun the verification; and hand off the
merge as the new candidate, or stop at a choice the record does not settle.
It then starts one attempt, as `run` would, on the candidate's branch in
that checkout. The records sharing the candidate reopen and run with it.
The assignment the agent receives is the ordinary
`/grove-work ID --interaction headless`: the mandate travels in the record,
where the work guide reads feedback.

Every refusal comes before the feedback is written:

- the target is not named;
- no branch, or several, hold the record in review;
- the target holds the work as done and every commit of the branch that it
  lacks has a copy there with the same patch, a rewritten copy, which it
  explains as `integrate` does, whether or not a checkout is on the branch;
- no checkout is on that branch;
- the candidate merges without a conflict;
- an attempt of any member is running or orphaned;
- no budget or permission mode is supplied;
- every member would wait once active, as `run` refuses it: an open
  question blocks it, or it needs a prerequisite the branch does not hold;
- the branch's checkout lacks the grove-work skill, or holds a skill or
  reviewer marked with an entrypoint revision this `grove` does not serve;
- the provider executable is missing or its `--version` fails.

The board's `m` passes the prediction it showed, so a candidate or target
commit that changed since it is refused too. If the attempt still fails to
start after the feedback, as when another launch takes the work meanwhile,
the error says the feedback stands and gives the
`grove run … --branch … --worktree …` that launches it. A second `resolve`
finds the work active and is refused. One operation is one attempt; nothing
retries, and a target that moves during it shows in the next prediction. The
[board](board.md#judging-a-candidate) names the resolution in the Review
block.

`stop ATTEMPT` sends SIGINT through the owner so the provider ends its turn,
SIGKILL to its process group after 15 s, and writes the result; an orphaned
attempt is stopped directly and reconciled. Stop touches neither the worktree
nor the record. A result is facts, never acceptance: the record's own status
on the branch, which the headless guide sets, is the handoff, and a process
exit or a `result` event proves nothing about it. `GROVE_CLAUDE` names
another executable, for fakes.

## Sweep

`sweep` acts on every candidate in review under the owner's standing
policy, the `policy:` mapping the [record model](record-model.md) describes
([G-260925-5wrn8](../grove/G-260925-5wrn8-resolve-approve-and-inte.md), under decision
[G-260925-wh9ax](../grove/G-260925-wh9ax-delegate-conflict-resolu.md)). It runs in the
target's checkout, whose committed `grove.yaml` holds the policy, and every
act is attributed to that file's revision, as `policy grove.yaml
sha256:…`. Without a policy it is refused: nothing is automatic. The owner,
a scheduler such as cron, or the board's `S` runs one, and so does the
owner process of an attempt that hands work off ([after an
attempt](#after-an-attempt)); Grove starts no resident process, and a wait
that has not changed does not retry. One sweep of a repository runs at a
time: while one holds `.git/grove/sweep.lock`, `sweep` and `S` are refused,
having done nothing, and an attempt's owner waits for it.

For each candidate, in ID order, it decides one act and prints it with the
reason; `--dry-run` stops there and writes nothing:

- **skip**, when the target already holds the branch's tip;
- **wait**, for the owner, when the candidate is shared, on several
  branches or without a checkout, already approved, blocked by an open
  question, has an attempt running, or fails a condition of the policy;
- **resolve**, for a conflict, when `resolve` is present, the record holds
  no earlier resolution feedback naming the same target commit, and the
  policy's `budget` still covers the attempt's. The attempt is `resolve`'s,
  with the permission mode, model and effort of the target's `run:`, never
  the candidate branch's (without a permission mode there it waits), and
  its feedback begins `delegated under policy grove.yaml sha256:…,
  budget N USD`;
- **approve**, or **integrate** with `integrate: true`, for a clean merge
  that meets the policy.

Work already accepted waits: delivering it is the owner's `integrate`.

To approve, it predicts the merge again, since an earlier act of the same
sweep may have moved the target, and merges the branch's tip into that
target commit in a temporary worktree outside every checkout, and runs each
`verify` command there with `sh -c` in the project's directory; the first
failure leaves the target and the record unchanged and prints the command's
last output. Once they pass, it approves in the branch's checkout with the
verdict `delegated under policy grove.yaml sha256:…: review ID examined X
with no open finding; merged with TARGET at T, verification passed
(COMMANDS); attempt A produced it for N USD` (or that no Grove attempt is
recorded as producing it), with `approved_by: policy sha256:…`, the
policy's revision; then, with `integrate: true`, it delivers as `integrate`
does, refused if the target moved from the verified commit. The delivery's
message says `Integrated under policy …`, and the sweep prints the `git
revert` that reverses it. A delivery refused after a delegated approval
leaves the record accepted, waiting for the owner's `integrate`. Such an
acceptance is told apart from the owner's wherever the record's standing is
shown: the board's Review block reads `accepted under policy` where the
owner's reads `accepted`, and `show --json` gives `approved_by` as
`policy`, or `owner`.

### After an attempt

The owner process of an attempt, once it has written `result.json`, sweeps
the work the attempt handed off, alone
([G-260928-dtrnw](../grove/G-260928-dtrnw-run-sweep-from-a-finishi.md)): the
members its result holds in review with a candidate, in the checkout Git
lists on the target branch, under that checkout's committed policy, as
`sweep` decides and acts. It appends each fact, with its time, to the
attempt's `sweep.log`, which `attempt` prints as `Sweep:` lines and the
board's attempt screen shows. It writes nothing at all when the attempt
handed nothing off, the launch had no target, or the target's `grove.yaml`
has no policy. It reports, and does nothing else, when no checkout is on
the target, that checkout has uncommitted changes to tracked files or
cannot be read, or the plan is refused, as for a `grove.yaml` that does not
parse. While another sweep runs it says so and waits for it, since nobody
would sweep its work later. A resolution it starts is an attempt of its
own, whose end sweeps again: a candidate that still conflicts with the same
target commit then waits for the owner. The owner has released the
attempt's lock before it sweeps, so the attempt reads as finished
meanwhile, and the environment variable that made it the owner is gone
before any `verify` command runs.

## Init

`init` runs at the top of a Git checkout (or `--project /absolute/path`). It
writes `grove.yaml` (`records: grove`, `brief: grove/brief.md`, and no `run:`
launch defaults), the record
root, a placeholder brief that states no intent, the `grove-work` and
`grove-shape` entrypoints for Claude Code (`.claude/skills/`) and Codex
(`.agents/skills/`), and Claude Code's `grove-reviewer` agent definition
(`.claude/agents/grove-reviewer.md`), each marked as managed. The reviewer is
read-only, inherits the session's model at `high` effort, and loads the
review guide the work guide's step 6 dispatches (`grove guide review`), as
the skills load theirs. Codex has no
equivalent, so a Codex session reports its missing independent reviewer as
step 6 says. It prints one line per path:
`created`; `kept` for an existing `grove.yaml`, whose own `records` and
`brief` it then follows, for the brief and the record root, and for an
entrypoint without the marker, which is yours; `unchanged`; or `updated` for
a marked entrypoint whose template changed in the binary, which is the
managed update. A `grove.yaml` that is not a schema 4 configuration, or a
directory, symlink, or unreadable file at a managed path, is a conflict: init
prints every reason, writes nothing, and exits 1. It never reads or writes
`AGENTS.md` or `CLAUDE.md`: the entrypoints defer to them for how the CLI is
invoked, and `init` ends with a note on stderr naming what they should say if
the entrypoints' defaults are not wanted: how `grove` is invoked, and how
work and proposal branches are named; that `grove guide model` describes
`grove.yaml`'s keys, such as `target`; and to commit what it wrote: an
attempt's worktree holds only committed files, and `run` refuses one without
the `grove-work` skill.

### Entrypoint revisions

The files `init` writes are thin: the work, shaping and review instructions,
including what an assignment or a shaping request may hold, are the guides
the binary prints, so installing another `grove` changes what a fresh
session follows without rewriting them. What an entrypoint needs of the
binary is its *entrypoint revision*, an integer apart from the release
version and the record schema, written in each managed file as
`grove entrypoint revision N` and passed as
`grove guide NAME --entrypoint N`. From revision 2 an entrypoint holds
nothing the guides evolve (the assignment grammar and the review brief are
the guides'), so a guide change never needs a new revision; one changes only
when an entrypoint needs something an older `grove` lacks, or a newer one
stops serving what an older entrypoint asks. This `grove` writes and serves
revision 2. What `init` wrote before revisions existed, marked files with no
revision line, is revision 1, `unrevised`: each generation kept its own
assignment grammar and review brief (the earliest rejects `--until plan`,
and has no reviewer), and nothing in the file says which, so it is not
served. Refresh it with `init`.

`guide --entrypoint N` refuses a revision it does not serve, printing
nothing and saying how to repair it, and a `grove` from before revisions
refuses the option; either way the entrypoint stops at its first command,
before any work. A plain `guide` always prints, since a person reads it the
same way, and an unrevised skill loads its guide that way; so the work and
shaping guides' Inputs tell a session loaded through a managed skill with no
revision line to stop and name the repair. An unrevised reviewer loads no guide
and is reached only through `init --check` and `run`.

`init --check` diagnoses each managed path and writes nothing, one line
each: `current` (this `grove`'s template), `compatible` (a served revision
in other text, older or edited), `unrevised` (no revision line), `incompatible`
(a stated revision this `grove` does not serve, newer or older, or not a
number), `missing`, `custom` (no marker: yours, listed and never judged), or
`conflict` (not a file `init` could replace). It exits 1 when any path is
unrevised, incompatible, missing or a conflict. Different text alone never
makes a file incompatible.

Upgrading is installing the new `grove`, then in the target's checkout
`init --check`, `init`, and committing what it updated; `init` rewrites only
marked files and keeps unmarked ones, the configuration, the brief and the
records. A linked worktree holds its own branch's copy, which `init` in
another checkout never touches: run `init` in that worktree and commit it
there, or merge the target after the refresh is committed. A session that
already loaded an entrypoint or a guide keeps what it read, so start a new
session (or reload skills) after refreshing. `run` checks the worktree it
launches in against the launching `grove`, and refuses an unserved revision
before any spend; the `grove` the session itself runs (`PATH` or the
project's instructions) is not checked, so where it differs, the session's
first `guide` is what stops it.

Check which `grove` answers before `init`: the predecessor also has an
`init`, which would write its own scaffolding instead. Codex runs each
command through a login shell, so it sees the profile's `PATH`, not the
caller's: put the build directory on the login `PATH` ahead of any other
`grove`, or name the executable in the target's `AGENTS.md` or `CLAUDE.md`.
When the wrong `grove` answers, the entrypoints stop and say so rather than
act.

### Migrate

Every command but `migrate` refuses a `grove.yaml` at `schema_version: 3`.
`migrate` reads that project and prints one line per work record that
changes, `ID  FROM -> TO  WHY`, a `reconcile:` line per problem, and the
counts; it writes nothing. Proposed, active, abandoned and unapproved review
work keep their meaning. An approval in review, or on done work whose
candidate the configured target contains, becomes the acceptance it was,
`approved_by` read from the last verdict on that candidate (`owner`, or the
policy revision a sweep's verdict names). Every other done record keeps
schema 3's claim, labelled so wherever standing is shown and never
presented as verified; a candidate Git cannot read is a problem to
reconcile, never a guess. `migrate --commit`, under the write lock, in a
checkout whose project has no uncommitted changes and no problems, first
records `refs/grove/schema-3/BRANCH` at HEAD, the way back, then writes
`grove.yaml` at 4 and each changed record, keeping its ID, path and body
and appending one `Migrated to schema 4` paragraph, and commits them alone.
A result that would not validate is restored. Old commits stay readable
with the binary of their time; a branch still at 3 is refused as a source
until it is migrated too, or rebased onto a migrated target.

## Version and guide

`version` names the executable and the workflow content it ships, in one
line that `attempt.json` also records as `grove_version`:

```text
grove VERSION (COMMIT[, vcs OTHER][, modified]) guides sha256:GUIDES content sha256:CONTENT
```

- `VERSION` is the release stamp, else the module version Go recorded: for
  `go build` in a checkout, a tag or pseudo-version from that checkout, with
  `+dirty` for uncommitted changes; for `go install …@VERSION`, that tag or
  pseudo-version; `(devel)` for `go run` and a tree without `.git`; and
  `(version unknown)` without build information.
- `COMMIT` is the stamped commit, else Go's `vcs.revision`, else
  `commit unknown`, as `go run`, a build from a source archive without
  `.git` and a `go install` print it. Where Go recorded a different commit
  from the stamp, `vcs OTHER` follows it. `modified` is Go's `vcs.modified`:
  the checkout Go found had uncommitted changes, stamped or not. Build from a
  primary checkout or a clone: for a linked worktree that lies inside its
  repository, Go only recognises the enclosing checkout's `.git` directory
  and takes that checkout's pseudo-version, commit and cleanliness instead
  (observed with go 1.26.2 and 1.26.8), which a stamp then shows as `vcs OTHER`; the
  `content` digest still names the worktree's own files.
- `guides` digests the work, shaping and review guides and the record model,
  which `guide` prints. `content` digests the workflow the binary ships into
  a session or a project: those four and the harness entrypoints `init`
  writes, the reviewer definition among them. It leaves out the `grove.yaml` and
  placeholder brief `init` creates once and never rewrites. Neither digest
  describes a project's installed files, which may be custom or from another
  release; an attempt records those itself ([Attempts](#attempts)).

The predecessor rejects `version` as an unknown command, so the line tells
the two apart.

A release build stamps the version and the full commit through the linker,
with `-trimpath` so no machine path enters the binary; Go embeds no build
time. The same command serves a clone and an extracted source archive:

```sh
go build -trimpath -ldflags "-X github.com/mascah/grove.version=v0.1.0 -X github.com/mascah/grove.commit=$(git rev-parse HEAD)" -o bin/grove ./cmd/grove
bin/grove version        # grove v0.1.0 (<that commit>) guides sha256:… content sha256:…
go version -m bin/grove  # Go's own record: -trimpath=true, and vcs.* from a clone
```

`bin/` is ignored, so an artifact leaves the tree clean for the next one;
`-o grove` would write into this repository's `grove/` record directory.
From an archive, pass the commit the archive was made from instead of
`git rev-parse`. With `-trimpath`, Go leaves `-ldflags` out of the build
settings, so `version` is where the stamp is read. A stamp is a claim the
builder makes: `version` prints it as given, and only the builder's inputs
make it true. Without a stamp, the line falls back as above and never
presents a release. Every artifact of one release, whatever its platform,
prints the same version, commit, `guides` and `content`, with neither `vcs`
nor `modified`; any difference means one was built from other inputs.

`go install …@v0.1.0` resolves a pushed tag, `@latest` the newest one and
`@COMMIT` only a pushed commit; upgrading an installed binary is running the
same line again, and nothing self-updates.

`guide work`, `guide shape` and `guide review` print the
[work](work-execution.md), [shaping](work-shaping.md) and
[review](work-review.md) guides the binary carries, so the workflow version
is the executable version and no copy is edited elsewhere. `guide work`
prints the work guide's head: everything through step 3, ending with a table
of the parts after it and the step at which each is read. `--part NAME`
prints one part alone (`prepare`, `implement`, `review`, `checkpoint`,
`handoff`, `judge` or `invocation`), `--part all` the whole guide, and any
other name is a usage error that lists them. The head and the parts, in
order, are the file byte for byte, so the file stays the one owner and a
session reads a step when it reaches it rather than carrying all of it from
the start. `--entrypoint N`
is how an entrypoint `init` wrote asks for a guide
([Entrypoint revisions](#entrypoint-revisions)). `guide model`
prints the [record model](record-model.md) the guides cite, the contract that
binary validates, so a session in any project reads it without Grove's
repository. Since it ships into projects whose own `G-` IDs are live, it links
to no record.
