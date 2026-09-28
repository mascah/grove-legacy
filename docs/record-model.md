# Record model

The record contract this binary validates: configuration, files, types,
fields, statuses, and what `check` and `update` refuse. `grove --help` gives
each command's usage, and the project's brief (`grove brief`) owns its
direction. Every refusal names the rule it enforces, and a refused value or
key comes with what is accepted.

## On-disk contract

### Configuration and discovery

One `grove.yaml` at the project root holds one YAML mapping. `grove init`
writes `schema_version: 3`, `records: grove` and `brief: grove/brief.md`
where none exists, with the record root and a placeholder brief, and keeps
an existing configuration that validates.

| Key | Value |
| --- | --- |
| `schema_version` | Required, exactly `3`; any other or a missing version is refused, never migrated |
| `records` | Required: the record root, a dedicated subdirectory relative to `grove.yaml`, not the project root, never absolute, through `..` or a symlink; it must exist |
| `brief` | Optional: the one project brief, a clean project-relative `.md` path without `..`; not a record, and exempt from discovery (compared without case); a live command needs a regular file there |
| `target` | Optional: the local branch work merges into, such as `main`; surrounding spaces or a `refs/` prefix are refused |
| `run` | Optional launch defaults for `run`, `resolve` and the board's `R`, named as their flags: `budget`, `permission_mode`, `model`, `effort` |
| `policy` | Optional standing delegation to `grove sweep`: `budget`, `resolve`, `approve`, `integrate`; absent, nothing is automatic |

Any other key, at any level, is refused. `target` is compared with branch
names, never passed to Git; every source whose `grove.yaml` names one must
agree, a source naming none has no say, and a conflict or a missing branch
leaves no target, with a note. It labels the current view and never decides
it.

Under `run:`, `budget` is a positive decimal dollar amount and the others
one word each. Without it every launch needs `--budget` and
`--permission-mode`. A launch reads only its own checkout's `grove.yaml`.

Under `policy:`, `budget` bounds the dollars of one sweep's resolution
attempts and is required with `resolve`. `resolve`, a mapping whose
optional `budget` bounds one attempt (else `run:`'s), allows one resolution
attempt per target commit for a candidate that conflicts with the target,
with the permission mode, model and effort of the `run:` beside the policy;
without a permission mode there, it waits.
`approve` takes `verify`, a required list of shell commands; `max_lines`, a
positive count of added plus removed lines, the record's own file excluded;
and `never`, project-relative `path.Match` patterns or `DIR/**`.
`integrate: true` needs `approve`. A malformed policy is a diagnostic named
`policy.KEY` that stops every command. Whatever it says, sweep approves only
a candidate in review, unapproved, that no other record shares; blocked by
no open question; with nothing after it on its branch but its record; named
by a `current` review that examined it, or an earlier commit from which only
records changed, every such review's last line starting `Open findings:`
being `Open findings: none`; changing nothing outside the project, no
`never` path and no binary file, within `max_lines`; and merging cleanly
into the target, the merged result passing every `verify` command in a
temporary worktree.

Without `--project DIR`, the nearest `grove.yaml` upward from the current
directory is used, the search stopping at the Git checkout root, or the
filesystem root outside Git. `--project DIR` names the directory holding it.
Reading needs no Git, and `list`, `show` and `check` read only the selected
checkout's live files.

### Folders and files

Every `.md` file (lowercase extension) beneath the record root, at any
depth, dot folders included, is a record, except the brief; one that is not
a valid record is an error, never skipped. Symlinks are refused; other
extensions are not records. A folder confers no type or status.

`new` writes `ROOT/ID-SLUG.md`, flat, with `--slug` or the title as
lowercase ASCII letters and digits joined by hyphens, at most 24 characters,
else `record`. The name is a convention: ID, type and dates are read from
frontmatter. No command renames, deletes or moves a record (`convert` moves
a document into the root), and none rewrites a body; a verdict or feedback
is appended as a paragraph.

Frontmatter is one YAML mapping between `---` lines, then a free Markdown
body. A BOM and CRLF are kept. Refused: invalid UTF-8, duplicate or unknown
keys, aliases, merge keys, custom tags, a second document, null values, and
scalar coercion (a string field must be a string; `priority` an integer).

### Identity and dates

An ID is `G-`, the UTC issue date as `YYMMDD`, a hyphen and five lowercase
Crockford base32 characters (`0-9a-hjkmnp-tv-z`): `G-260925-7k2qm`. Any
other spelling, such as `W-001`, `G-1234` or an uppercase tail, is invalid.
An ID carries no type and never changes; IDs are unique per checkout,
matched exactly and ordered as strings.

`new` and `convert` need Git. Under an `flock` on `grove/write.lock` in the
Git common directory, which serializes `new`, `convert` and `update` across
worktrees, `new` draws the tail from `crypto/rand` and draws again while the
ID is held by the project, by an `id:` line beneath the record root on any
branch, remote-tracking ref or tag, or by a live record in any worktree.
After eight draws it fails, creating nothing.

`created` and `updated` are optional on every type: quoted UTC
`YYYY-MM-DDTHH:MM:SSZ`, `updated` not before `created`. `new` writes both
equal; `update` keeps `created`, sets `updated` when content changes, and
refuses a clock earlier than either. `show ID --json` gives `revision`,
`sha256:` and the hex digest of the exact bytes, and `update --expect
REVISION` refuses other content.

## Types and fields

| Type | Statuses (first is what `new` writes) | Own fields |
| --- | --- | --- |
| `work` | `proposed`, `active`, `review`, `done`, `abandoned` | `kind`, `size`, `priority`, `members`, `depends_on`, `candidate`, `approved` |
| `question` | `open`, `resolved` | `blocks` |
| `decision` | `proposed`, `accepted`, `rejected`, `superseded` | none |
| `term` | `proposed`, `settled` | none |
| `plan` | `current`, `superseded` | `work` |
| `review` | `current`, `superseded` | `work`, `examined` |
| `page` | none: no lifecycle | none |

Every record may carry `id`, `type`, `title`, `status` (never a page),
`relates_to`, `created`, `updated` and `formerly`, then its type's own
fields.

| Field | Form | Meaning |
| --- | --- | --- |
| `id` | an ID, required | Identity |
| `type` | a type above, required | Classification; missing or unknown is an error, never a page |
| `title` | nonempty string, required | The label; a term's title is the term, unique among terms ignoring case and surrounding space |
| `status` | one of the type's statuses, required but on a page | Lifecycle state |
| `relates_to` | list of record IDs | Related records, with no order or gate |
| `created`, `updated` | quoted UTC timestamps | Chronology |
| `formerly` | string that only `convert` writes | The ID or document path a record replaced; unique without case, and never a live record's ID |
| `kind` | `feature`, `fix`, `refactor`, `investigation`, `tooling`, `release` | What sort of work |
| `size` | `small`, `medium`, `large` | Coarse scope; `small` selects the compact handoff (`grove guide work`), and no size sets another execution rule |
| `priority` | integer 1 (highest) through 5 (lowest) | Importance, independent of dependencies |
| `members` | list of work IDs | Child work in this outcome; order is presentation |
| `depends_on` | list of work IDs | Prerequisites delivered before this work: it needs their result, or building it first or alongside would redo or conflict with them |
| `candidate` | quoted Git commit, 7 to 40 lowercase hex digits | The last implementation commit, offered for judgment |
| `approved` | quoted Git commit | The owner's approval of `candidate` |
| `blocks` | list of work IDs | Work whose outcome needs this question's answer |
| `work` | list of work IDs | The work a plan or review belongs to |
| `examined` | quoted Git commit | The commit a review looked at |

An absent optional field is unspecified; nothing is defaulted. A list may
not repeat an ID or name its own record, and each target must resolve in
this checkout to one record, which is work in every list but `relates_to`.
`depends_on` cycles and `members` cycles are each refused; membership may
nest and be shared. Reverse links, such as a work's plans or dependents, are
derived, never stored. A resolved question stops blocking; an abandoned
prerequisite is not delivered.

## Work lifecycle

Proposed → Active → Review → Done, with Abandoned by explicit human
decision. `check` and `update` enforce:

- `review` requires `candidate`. Records on one branch sharing a
  `candidate` are one group, handed off, reopened and integrated together.
- `approved` must equal `candidate`, and holds only while status is
  `review` or `done`: reopening unsets it in the same update.
- `update` writes `done`, or changes a done record's `candidate`, only when
  the candidate is an ancestor of HEAD and, where `target` is set, on that
  branch; a record that stays done cannot lose its candidate. A done record
  without `candidate` predates this rule: it validates and its other fields
  stay editable, and none is newly written.

`approve`, `feedback`, `integrate`, `resolve` and `sweep` act on this
lifecycle as `grove --help` describes. A verdict `sweep` gave begins
`delegated under policy grove.yaml REVISION`, names the review and
verification it relied on, and is shown apart from the owner's own.

Not enforced, and left to the guides: the order of transitions, so a done
record can be reopened; that Abandoned needs a human decision; that a review
record exists before Review where the work guide's handoff calls for one;
where done is written without a `target`; that a reopened record's
candidate moves to its new commits; and that a `superseded` decision names
its replacement in `relates_to`.

## Reading and writing records

`check` reports every configuration, file, metadata and relationship problem
with its path, line and field, or `OK: N records`. `list` and `show` refuse
an invalid project rather than show part of it. Exit codes are 0 for
success, 1 for failure and 2 for usage.

`update ID` edits only the named frontmatter entries and `updated`. It
accepts `type`, `title`, `status`, `relates_to` and the type's own fields,
never `id`, `created`, `updated` or `formerly`; `type` and `title` cannot
be unset, nor `status` but toward a page. Lists are JSON arrays: `--set 'work=["G-260925-7k2qm"]'`. `--set type=TYPE` reclassifies in
place when the same update meets the new type's whole contract, such as
`--set type=page --unset status`. The project must validate with the
change, or nothing is written.

`convert PATH --type TYPE --title TITLE` makes a record of a Markdown
document outside the record root: its bytes, less a leading BOM, as the
body, the type's first
status, `formerly: "PATH"`. The brief, a missing source, and one that a
record's `formerly` already names (without case) are refused.

## Not records

The brief; attempts, which `grove run` keeps as files under the Git common
directory, never committed; and supporting material, which is prose and
ordinary links in a record's body. There is no report type and no assignee
field.
