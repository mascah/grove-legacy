---
id: "G-260926-vkv48"
type: work
title: "Rename legacy records to the date form"
status: active
created: "2026-09-26T16:09:19Z"
updated: "2026-09-26T16:26:03Z"
kind: refactor
size: large
relates_to: ["G-260926-yvjy6", "G-260926-2da4n", "G-260926-pgj43", "G-260921-r491p", "G-260921-gtydy", "G-260921-czt8x", "G-260921-905y3"]
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
  commit, writing `created` when it was missing, as G-260921-r491p did in `ed14d92`;
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
- [G-260921-czt8x](G-260921-czt8x-identity-and-path-migrat.md) gains a second table, legacy ID to
  date-form ID and path, written after the rewrite so its old IDs survive
  it, with a sentence saying that commit subjects and branch names in Git
  history keep the legacy IDs and this table resolves them. The README's
  sentence naming G-260921-czt8x as the map stays true.
- The commands reference and the record model's conversion section
  describe the command, marked one-time and slated for removal by
  G-260926-19gzg; `grove guide model` prints it. `CLAUDE.md`'s rule against
  hand renumbering stands unchanged: the command is the sanctioned path.
- Before the command runs here: remove the merged `worktree-G-260926-pgj43` worktree
  and branch and the merged remote branches `worktree-G-260922-jtsed` and
  `worktree-G-260922-g6e7p`.

Out of scope:

- Any write to nullsec (G-260921-905y3: evidence only). The owner runs the command
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
   IDs finds them only in G-260921-czt8x's map table and in test string literals;
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

## Evidence

Headless attempt `G-260926-vkv48.20260926T162220Z`, branch
`worktree-G-260926-vkv48` in `.claude/worktrees/worktree-G-260926-vkv48`,
base `main` `a9f8fce`, started from this record at
`sha256:a68aadb2fc39f30456f7cfc88a09c18cd8c7f9480628934874de52ce7ee981e5`.
Plan [G-260926-yzd8f](G-260926-yzd8f-plan-for-renaming-legacy.md), review
[G-260926-w0std](G-260926-w0std-review-of-the-renumber-c.md) (three rounds,
examined `24a43fa`, "Open findings: none"). The candidate differs from
`24a43fa` only in this record, the plan and the review.

Commits: `80e3397` the command (`internal/renumber`, `grove renumber`,
`update.Set`, docs); `ade21b3`, `fd04146` the slug rule; `ef1cdfc` the
194 renames alone, content unchanged; `b3d6c92` the rewrite the command
made plus nullsec's numbers restored by hand; `ac69921` references outside
`grove/` and the map table; `3bdcf25`, `502d4bf`, `24a43fa` review fixes.

**Decisions taken, and why.**

- A slug leaves out every ID its title cites, with a following `'s`:
  102 titles cite IDs, so slugs from the rewritten titles were mostly ID
  and a cut at the cap left fragments like `g-260` that read as legacy IDs.
  The date-form record G-260926-afe5w, whose filename cited a legacy ID, took
  the same kind of name. The Constraints said `create.Slug(title)`; the
  plan records the adjustment, and readability is acceptance 8.
- The renames are committed alone before the rewrite: in one commit 3
  records here (18 in nullsec) fell under Git's 50% rename similarity and
  lost `git log --follow`. The record model says to do the same.
- Nullsec's own record numbers in its cutover records (G-260921-905y3,
  G-260922-r1dhw, G-260922-9d399, G-260925-dzxm6, G-260926-pgj43) are not
  this repository's IDs and keep their text.
- Outside `grove/`, citations were rewritten; examples (the shipped guides,
  `CLAUDE.md`'s invocations, `docs/board.md` and `docs/commands.md` UI and
  command examples, the README's invocations, `--help`, `metadata.go`'s
  diagnostic, `renumber.go`'s own comment), fixtures and test string
  literals keep theirs for G-260926-19gzg. Test comments that describe
  fixtures were kept; the 36 test comment lines and 80 non-test Go comments
  citing records were rewritten. The README's `show` demo, which runs
  against this repository, was rewritten.
- The command skips a running or orphaned attempt of date-form work (its
  owner still appends) and refuses one of legacy work.
- A stale link in G-260923-twv25 to a misspelt filename of G-260922-08wxx now
  names the record's file.

**Acceptance.**

1. At `24a43fa`: `check` OK, 198 records; no `G-NNN-…` file under
   `grove/`. A case-sensitive search for three-digit IDs outside test files
   and the map page finds only examples and fixtures (list above), in
   records only non-records (`G-000`, `G-300`, `G-999`, the never-created
   `G-172`) and nullsec's numbers. Markdown links: no new broken link; the
   two deleted `docs/prompts/*.txt` in G-260919-nddsf were broken on `main`
   too, and links into `../skills` and `../bench` resolve only from the main
   checkout. `versions` lists no legacy ID in a clone holding only this
   branch; here the 193 legacy rows are all `main`'s until the merge. The
   board opens (below).
2. Every new ID's date is its `created` date; days are non-decreasing in
   legacy order; `ef1cdfc`, `b3d6c92` change no `updated` or `created`
   line; the longest filename is 42 characters.
3. `grove attempts` lists 40: the 39 finished under their new IDs and this
   running one. Through `internal/tui/testdata/terminal.py`'s `Session` on
   the real repository, `A` then `o` on a finished attempt opened card
   G-260925-5wrn8, renamed from a legacy ID.
4. `internal/renumber/renumber_test.go` covers every listed case, plus an
   escaped ID in the events and a cited ID in a title;
   `TestRenumberThroughTheCLI` checks the JSON lines, the refused rerun
   (exit 1) and the usage error (exit 2).
5. All 194 renamed files keep their pre-rename commits under
   `git log --follow`; on the board, `v` on G-260925-5wrn8 lists history
   back to the shaping commit `214dd3d`.
6. A clone of nullsec at `37548aa` under `/tmp`, with the binary from this
   branch: the command completed, printed 127 map lines, wrote `created` on
   the 126 records that lacked it (dates from 2026-09-15 UTC), `check` OK,
   no legacy token left under `grove/`. Nullsec's legacy order is not date
   order (10 inversions), so its map is not non-decreasing by day.
7. `docs/record-model.md` (Renumbering, printed by `grove guide model`) and
   `docs/commands.md` describe the command; their links resolve.
8. The owner's.

Verification at `502d4bf`, rerun unchanged since except one Markdown
record: `go vet ./...` clean, `gofmt -l .` empty, `go run ./cmd/grove
check` OK, `go test -count=1 -timeout 120s ./...` all ok (the terminal
checks included).

**Limits.**

- The attempt store under the common directory was rewritten before the
  merge: until integration, `main`'s board and `attempts` show attempts
  under IDs its records do not have yet. A copy of the store before the run
  is at `/tmp/vkv48-attempts-backup` and does not outlive `/tmp`. This
  attempt's own log keeps legacy IDs.
- The merged remote branches of G-260922-jtsed and G-260922-g6e7p were not
  deleted: that is a push. `worktree-G-260926-pgj43` and its branch were
  already gone when the step came.
- The command reads local branches only, and does not look at detached
  worktrees.

## Next

In review, candidate below. The change touches `CLAUDE.md`, `grove.yaml`
and `.github/`, so under the standing policy it waits for the owner. Judge
the names (acceptance 8) with the board and `ls grove` in this worktree,
then:

```sh
cd /Users/mascah/GitHub/mascah/grove/.claude/worktrees/worktree-G-260926-vkv48
go run ./cmd/grove approve G-260926-vkv48 "VERDICT"
cd /Users/mascah/GitHub/mascah/grove
go run ./cmd/grove integrate G-260926-vkv48 --cleanup
just install
git branch -r --merged origin/main   # the two merged worktree- branches, then git push origin --delete them
```

Then in nullsec, following its own instructions: clean its 14 merged
branches and 3 stale worktrees, copy `.git/grove/attempts` if it has one,
run `grove renumber > map.jsonl`, commit the renames alone with the old
content and then the rewrite (the record model says how), repair the 372
comment references and 18 instruction-file references from the map, and
commit. Only then assign
[G-260926-19gzg](G-260926-19gzg-retire-the-legacy-id-for.md).
