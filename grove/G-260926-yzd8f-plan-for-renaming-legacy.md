---
id: "G-260926-yzd8f"
type: plan
title: "Plan for renaming legacy records to the date form"
status: current
created: "2026-09-26T16:25:26Z"
updated: "2026-09-26T16:26:00Z"
work: ["G-260926-vkv48"]
---

## Design

Work: [G-260926-vkv48](G-260926-vkv48-rename-legacy-records-to.md), under
[G-260926-yvjy6](G-260926-yvjy6-retire-legacy-ids-by-ren.md).

`grove renumber` lives in a new package `internal/renumber`, since it needs
both `attempt` (liveness of attempts) and `update` (the frontmatter editor),
and `attempt` already imports `update`. The CLI prints its map as `convert`
prints one conversion.

`renumber.Run(root, now)`:

1. Load; refuse on diagnostics. Refuse when no record has a legacy ID
   (three digits), writing nothing.
2. Take the write lock, reload, and refuse again on diagnostics.
3. Refuse while a local branch other than HEAD's and the target holds a
   file under the record root (`git ls-tree`), naming the branches.
4. Refuse while an attempt that is running or orphaned selects a legacy ID;
   a live attempt of date-form work is left untouched, since its owner still
   appends to its files.
5. Per legacy record, in ID order: its time is `created`, else the UTC
   author date of the first commit of `formerly`'s path, else of the
   record's own path (`git log --follow`); none refuses. Draw the ID with
   `create.Issue` at that time, avoiding IDs drawn earlier in the run; the
   slug is `create.Slug(title)`; the new path must not exist.
6. Write: add `created` where it was missing through `update`'s editor,
   rename, then rewrite every file under the record root, old basenames
   first, then old IDs as whole tokens (not after a letter or digit, not
   before one; a `-` either side is a token boundary, so `worktree-G-260920-svpbc`
   and `G-260920-svpbc-G-260920-z8vfp` are rewritten). Only IDs in the map are touched, so
   `G-000`, `G-1000` and date-form IDs are left alone.
7. Reload and require it to validate; then rename every attempt directory
   whose name starts with a mapped ID and rewrite each regular file of
   every attempt that is not live.
8. Return the map, one `{from, from_path, id, path}` per record, legacy
   order.

Rename detection: `git log --follow` needs the rename commit's similarity
over 50%. If the one-commit rename loses records whose content changes a
lot, commit the renames with unchanged content first (index only), then the
rewrite, and say so in the docs; measured at the step.

Outside the record root, the reference repair is by hand from the map, as
the record's Constraints list; test string literals, guide examples and
`CLAUDE.md`'s invocation examples stay. Mentions of nullsec's own numbers,
as in the nullsec cutover plan and review, are not this repository's IDs
and are restored after the run.

**Adjusted at the rehearsal (bounded, technical).** 102 titles cite IDs
("G-260920-svpbc card lineage plan"), so a slug from the title as rewritten was
mostly ID (`G-260919-q3mg5-g-260919-byjdx-g-260919.md`) and a cut at the
cap left fragments like `g-260` that read as legacy IDs. The slug leaves
out every ID a title cites, with a following `'s`: `card-lineage-plan`.
For the same reason a date-form record whose filename cites a legacy ID
(G-260926-afe5w) keeps its ID and takes that slug, since
records are never renamed by hand here. Readability is the owner's
judgment (acceptance 8).

## Steps

1. [x] Package, CLI wiring, test (acceptance 4). Docs: commands reference and
   record model (acceptance 7).
2. [x] Remove the merged `worktree-G-260926-pgj43` worktree and branch (found already
   removed when this step came). The merged
   remote branches `worktree-G-260922-jtsed` and `worktree-G-260922-g6e7p` need a push to
   delete, which this headless attempt does not make; the command reads
   local branches only, so they do not block it. Left to the owner.
3. [x] Run the command here; restore nullsec mentions; repair references
   outside the record root; G-260921-czt8x's second table.
4. [x] Evidence: acceptance 1-3, 5; nullsec in a disposable clone (6).
5. [x] Independent review, then hand off.
