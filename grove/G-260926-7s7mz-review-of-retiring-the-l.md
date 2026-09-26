---
id: "G-260926-7s7mz"
type: review
title: "Review of retiring the legacy ID form"
status: current
created: "2026-09-26T20:08:45Z"
updated: "2026-09-26T20:09:06Z"
work: ["G-260926-19gzg"]
examined: "db568f3c7d701ec5d470435f411f3d57c86a307f"
---

## Examined

Work [G-260926-19gzg](G-260926-19gzg-retire-the-legacy-id-for.md), no plan
(small; its Constraints are the plan), on branch `worktree-G-260926-19gzg`,
base `main` `2bd50b9`. Two rounds by independent `grove-reviewer` agents,
each fresh, read-only, writing only under `/tmp`: round 1 examined
`4ff88f3`, round 2 the fix commit `db568f3` over the same base. Each ran
`go test -count=1 -timeout 120s ./...`, `go vet ./...`, `gofmt -l .` and
`check`, all passing; round 1 also ran `terminal.py` (5 scenarios, exit 0),
a hand-authored `G-001` record in a clone (refused by `check`, exit 1),
`new` and `convert` there (date-form IDs), and the Constraints' searches.

## Findings

Round 1 (at `4ff88f3`), acceptance 1 to 6 met:

1. Medium, test coverage: restoring the three-digit alternative in `IDForm`
   passed every test, since the `zero_id` case was deleted and the bulk
   rewrite turned `TestAttemptNames`' legacy case into a date-form one.
2. Medium, knowledge: renaming the `init --check` verdict `legacy` to
   `unrevised` left term
   [G-260925-m9jcr](G-260925-m9jcr-entrypoint-revision.md) saying "revision
   1, legacy". The reviewer judged the rename within the record's authority
   (acceptance 3 forbids the word in `--help`; nothing parses the token; no
   entrypoint revision is needed) and asked that the handoff name it.
3. Low-medium, the record's text: G-260926-vkv48's rename rewrote this
   record's own acceptance 1 and example list to real date-form records
   (`G-260919-6mpmw` was `G-001` at `a9f8fce`), so read literally they
   cannot be met.

Round 2 (at `db568f3`): 1 and 2 resolved, verified by rerunning the
mutation (both new cases fail with the alternation restored) and reading
the term; no regression; 3 open until the handoff records it.

## Disposition

1. Fixed in `db568f3`: a `retired numeric id` case in
   `TestStrictRecordMetadata` and a refused `G-001.…` attempt name in
   `TestAttemptNames`, their literals split so the Constraints' search
   finds no fixture.
2. Fixed in `db568f3`: the term says "unrevised"; the handoff names the
   rename.
3. Recorded in the work record's Evidence with the original reading; no
   code change.

Open findings: none
