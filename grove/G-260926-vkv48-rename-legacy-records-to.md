---
id: "G-260926-vkv48"
type: work
title: "Rename legacy records to the date form"
status: proposed
created: "2026-09-26T16:09:19Z"
updated: "2026-09-26T16:11:02Z"
kind: refactor
size: large
relates_to: ["G-260926-yvjy6", "G-194", "G-195", "G-052", "G-064", "G-069", "G-041"]
---

## Outcome

Every record in this repository carries a date-form ID whose date is the
record's creation date, a filename derived from that ID and a 24-character
slug, and no reference anywhere in the repository or in its attempt store
names a legacy `G-NNN` ID, except the migration map that says what each one
became. The binary ships the command that did it, so the owner can run the
same rename in nullsec. The intent is the owner's, decided in
[G-260926-yvjy6](G-260926-yvjy6-retire-legacy-ids-by-ren.md) on 2026-09-26.

## Constraints

In scope:

- A one-time command, proposed as `grove renumber`, that under the write
  lock loads the project, refuses when `check` fails, and for each record
  with a legacy ID: takes the date from `created`, else from the first
  commit of the path `formerly` names, else from the record's own first
  commit, writing `created` when it was missing, as G-052 did in `ed14d92`;
  draws the ID through `create.Issue` with that time; derives the slug from
  the title with `create.Slug`; renames the file and rewrites the `id` line,
  leaving `updated` alone so the board's order holds. It then rewrites every
  file under the record root, old basename before old ID, whole tokens only,
  in frontmatter, links and prose alike, and the attempts directory under
  the Git common directory: directory names and every file in them. Stdout
  is one JSON line per record, `{from, from_path, id, path}`, as `convert`
  prints, and that is the map.
- The command refuses while another local branch holds records, since the
  board and `versions` read every local branch and would show both forms,
  and refuses a project with no legacy record as a no-op that writes
  nothing. Reruns are therefore safe.
- Outside the record root the caller repairs references from the map, as
  after `convert`. In this repository this work does that: `docs/board.md`
  and `docs/commands.md` where they cite records, the README, `CLAUDE.md`,
  the brief's links, `evals/`, `.github/workflows/ci.yml`, `lefthook.yml`,
  comments in non-test Go files, and comment lines in test files. String
  literals in tests are fixtures, not references, and keep their legacy IDs
  until [G-260926-19gzg](G-260926-19gzg-retire-the-legacy-id-for.md)
  converts them. The shipped guides' example IDs stay for the same reason.
- [G-069](G-069-migration-map.md) gains a second table, legacy ID to
  date-form ID and path, written after the rewrite so its old IDs survive
  it, with a sentence saying that commit subjects and branch names in Git
  history keep the legacy IDs and this table resolves them. The README's
  sentence naming G-069 as the map stays true.
- The commands reference and the record model's conversion section
  describe the command, marked one-time and slated for removal by
  G-260926-19gzg; `grove guide model` prints it. `CLAUDE.md`'s rule against
  hand renumbering stands unchanged: the command is the sanctioned path.
- Before the command runs here: remove the merged `worktree-G-195` worktree
  and branch and the merged remote branches `worktree-G-081` and
  `worktree-G-089`.

Out of scope:

- Any write to nullsec (G-041: evidence only). The owner runs the command
  there; Next says how.
- The validator, ordering code, test fixtures, guide examples and the
  brief's foundation sentence: G-260926-19gzg.
- Deleting the command: G-260926-19gzg, after nullsec has run it.
- Rewriting Git history, `formerly`, and the meaning of any record.

Observed evidence, 2026-09-26 at main `a691f5f`, is in G-260926-yvjy6: the
counts of records, links, references and tokens by area; the attempt store's
layout and the listing's silent skip of names that fail `IDForm`; history by
`git log --follow`; `create.Issue(p, prefix, now)` and the 24-character
`create.Slug`; every legacy record's `created` in ID order across 8 days
with at most 48 on one day; and nullsec's 126 records without `created`
whose `formerly` sources all have first commits. Tail collisions inside one
day are `Issue`'s to avoid, as today.

## Acceptance

1. In this repository after the command and the reference repair: `check`
   passes; no file under `grove/` is named `G-NNN-…`; a search for legacy
   IDs finds them only in G-069's map table and in test string literals;
   every Markdown link resolves; `versions` lists no legacy ID; the board
   opens.
2. Each new ID's date is its record's `created` UTC date; listing the map in
   legacy order shows the new IDs non-decreasing across days; `git diff`
   shows no `updated` line changed; every filename is at most 42 characters.
3. `grove attempts` lists all 39 attempts under their renamed work IDs, and
   the board's attempts view opens the renamed card from one of them.
4. A test renames a fixture project holding a record with `created`, a
   converted record without `created` but with a `formerly` source in
   history, links, `relates_to`, `work` and an attempt directory, and
   checks the dates, the rewritten references, the renamed attempt, the
   untouched `updated`, the JSON map, the refusal on a second local branch
   with records, and the no-op rerun.
5. The card of a renamed record shows its history from before the rename.
6. In a disposable clone of nullsec, as evidence and not a write: the
   command completes, writes `created` on the 126 records that lack it from
   their `formerly` sources' first commits, prints 127 map lines, and
   `check` passes there.
7. The commands reference and record model describe the command, every
   link in them resolves, and `grove guide model` prints the updated model.
8. The owner opens the board and the directory listing and judges the
   names readable.

## Next

Assign. The change is one commit of about 200 renamed files and touches
`CLAUDE.md`, so under the standing policy it waits for the owner's approval
in any case. After integration: `just install`, then in nullsec, following
its own instructions, the owner cleans its 14 merged branches and 3 stale
worktrees, runs the command, repairs the 372 comment references and 18
instruction-file references from the map, commits, and only then assigns
[G-260926-19gzg](G-260926-19gzg-retire-the-legacy-id-for.md).
