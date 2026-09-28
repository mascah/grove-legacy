---
id: "G-260928-rvg20"
type: review
title: "G-260927-n4wvk with row: both constraints found in every run, format 3 seen in three"
status: current
created: "2026-09-28T01:32:54Z"
updated: "2026-09-28T01:33:50Z"
work: ["G-260927-n4wvk"]
examined: "751c49d"
---

## Examined

The `with` row of [G-260927-n4wvk](G-260927-n4wvk-trim-the-framing-context.md)
acceptance 4: the two retrieval cases of
[G-260925-pbx81](G-260925-pbx81-evaluate-whether-agents.md), five headless
shaping runs each on Claude, under the mandate
[G-260927-jds8s](G-260927-jds8s-what-mandate-covers-reru.md) set ("reuse the
last config", its option 1), compared with the `without` row
[G-260925-khwkq](G-260925-khwkq-without-row-both-constra.md). Run
2026-09-28 01:20 to 01:31 UTC from `worktree-G-260927-n4wvk` at `751c49d`,
the `examined` commit, clean, holding `context` format 3 (`da457a6`), by a
headless `/grove-work G-260927-n4wvk` session, as two foreground pieces as
the `without` row was:

```sh
python3 evals/run.py run --case listed-constraint --runs 5 --budget 5 \
    --model claude-opus-5-5 --permission-mode auto \
    --config-dir ~/.cache/grove-evals/claude \
    --out ~/.cache/grove-evals/runs/2026-09-28-G-260927-n4wvk-with/listed-constraint
python3 evals/run.py run --case code-constraint --runs 5 --budget 5 \
    --model claude-opus-5-5 --permission-mode auto \
    --config-dir ~/.cache/grove-evals/claude \
    --out ~/.cache/grove-evals/runs/2026-09-28-G-260927-n4wvk-with/code-constraint
```

`python3 evals/run.py selftest` passed first at the same commit.
Configuration, from the two reports: Claude Code 2.1.283 (the `without`
row had 2.1.282), model reported `claude-opus-5-5`, permission mode `auto`,
guides digest `4110723aac9f` (the `without` row's was `41324c3655a1`),
fixture commits `5c8855e` (listed-constraint) and `7fe91f8`
(code-constraint), cap $25 per piece printed first. The config directory
held the same login, `settings.json` (`autoMemoryEnabled: false`, `theme`,
`tui`), seven synced skills and three synced plugins. The output
directories keep every transcript, clone, remote, `state.json`, `run.json`
and `report.md`, outside every checkout and not committed. The report's
`grove version` line names main's `3ec0f68`, the revision a linked
worktree's build can lag to; the base commit `751c49d` and the transcripts'
`Grove work context (format 3)` show the build was this checkout's. Every
score below was given by the implementing session from `run.json`, the
transcripts' commands and the records on each proposal branch, so its
scorer is `judge`; the owner column is open.

The fixture issues new IDs per build: in listed-constraint the holding
record is G-260928-nyvfq "Add tasks export", the distractors G-260928-ntgme
"Export tasks as CSV" and G-260928-18wbn "Remind the owner of tasks due
today", the unrelated shared record G-260928-vn85w "Sync tasks between two
machines"; in code-constraint the holding record is the decision
G-260928-rve8s "Keep the owner's notes in task files" and the unrelated
record G-260928-qdht2.

## Findings

**1. listed-constraint: found and applied in 5 of 5 runs, as in the
`without` row.** Every run refined G-260928-6sx01 "Give tasks a due date" in
place, `proposed`, with no question, and every acceptance has an item that
`tasks export` keeps exactly its seven keys, citing G-260928-nyvfq and the
widget. Every run also found that `done` and `drop` rewrite through `write`
and would drop `due`, and kept it. Run 2 also made G-260928-18wbn depend on
the due-date record. Every check passed in every run.

**2. code-constraint: found and applied in 5 of 5 runs, as in the `without`
row, and no run asked a question (2 of 5 did before).** Every run split the
work into a prerequisite "Keep a task's notes when a command rewrites its
file", citing G-260928-rve8s, whose acceptance keeps the lines below a
task's title byte for byte after `done` and `drop`, and a tag command that
depends on it. None opened the command-syntax question
G-260925-khwkq finding 2 called fixture noise; five runs cannot say whether
that is the change or chance. Every check passed in every run.

**3. Only three of ten runs saw format 3, and none needed it to find the
constraint.** `grove context` ran in listed-constraint runs 1, 4 and 5
(the `without` row: listed-constraint run 3 alone), each transcript showing
`Grove work context (format 3)`; no code-constraint run ran it. Runs 1 and
5 ran it after reading every record by a `for` loop over `grove/G-*.md`.
Run 4 ran it in the same command that first read two records, and read
G-260928-nyvfq and both distractors in the next, after `ls grove` had
already listed every file. As in the `without` row, a fixture of five or
two records is read nearly whole, so this row shows format 3 costing no
retrieval, not format 3 helping it.

**4. Retrieval facts are no worse.** Compared per case with the `without`
row's recomputed facts (G-260925-khwkq finding 4):

| case | row | holding read | distractors read | unneeded reads (files, plus `show`) | search | turns | cost |
| --- | --- | --- | --- | --- | --- | --- | --- |
| listed-constraint | without | 5 of 5 | both, 5 of 5 | shared record in 3 of 5 (runs 1, 2, 5), `.gitignore` 0 | none | 6 to 9 | $1.55 |
| listed-constraint | with | 5 of 5 | both, 5 of 5 | shared record in 3 of 5 (runs 1, 3, 5), `.gitignore` in runs 3, 5 | none | 7 to 10 | $1.72 |
| code-constraint | without | 5 of 5 | none planted | shared record in 5 of 5, `.gitignore` in run 2 | none | 9 | $1.80 |
| code-constraint | with | 5 of 5 | none planted | shared record in 5 of 5, `.gitignore` in runs 1, 2, 4 | none | 8 to 11 | $1.86 |

`holding read` and `distractors read` are equal in both cases, which is
acceptance 4's bar. The with row read `.gitignore` in three more runs, a
two-line file; the shared record's count is the same. Cost rose $0.17 and
$0.06 for five runs and turns by about one, within the spread of either
row, from sessions of which seven never ran `context`: not an effect of
the format.

**5. Cost and shape.** listed-constraint $0.31 to $0.39 a run, 50 to 74
seconds, $1.72 for five; code-constraint $0.34 to $0.43, 59 to 90 seconds,
$1.86 for five; $3.58 for the row against the $50 cap. Exit 0, no
permission denials, no timeouts, no error results. Tools: `Bash` only,
inside each clone apart from a probe copy of `tasks.py` in a `mktemp -d`
directory (listed-constraint run 2, code-constraint run 4).

Rubric ([`evals/README.md`](../evals/README.md)), scorer `judge` (the
implementing session), owner column open:

| case | run | scorer | constraint applied | brief constraint | handoff |
| --- | --- | --- | --- | --- | --- |
| listed-constraint | 1 to 5 | judge | 2, 2, 2, 2, 2 | 2, 2, 2, 2, 2 | 2, 2, 2, 2, 2 |
| code-constraint | 1 to 5 | judge | 2, 2, 2, 2, 2 | 2, 2, 2, 2, 2 | 2, 2, 2, 2, 2 |

Notes. Brief constraint: every due-date proposal keeps `due` one
frontmatter key in a file that stays in place; every tag proposal keeps
tags lowercase and the file in place. Handoff: `message-names` passed in
every run.

## Disposition

Acceptance 4 is met: `holding read` and `distractors read` are no worse
than the `without` row's in either case, and the constraint rubric row is
2 in all ten runs of both rows.

**Limits.** The rows differ in more than the context format: the guides
digest moved from `41324c3655a1` to `4110723aac9f` with the guide edits
landed since, and Claude Code from 2.1.282 to 2.1.283, so the row compares
the whole agent-facing surface at this commit, not format 3 alone. Only
three runs saw format 3 (finding 3), and at this fixture size neither
row's listing is what finds the constraint (G-260925-khwkq finding 3). One
model, five runs per case, a judge that is the implementing session: a
pattern, not a rate.
