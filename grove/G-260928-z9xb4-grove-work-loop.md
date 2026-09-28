---
id: "G-260928-z9xb4"
type: page
title: "Grove work loop"
created: "2026-09-28T19:28:59Z"
updated: "2026-09-28T19:28:59Z"
---


## The loop as it runs today

Observed at main `6fbb888`, 2026-09-28, from [the work guide](../docs/work-execution.md),
[the board](../docs/board.md) and [the commands](../docs/commands.md). Nothing
here is proposed: those three documents own the workflow, and this page is
redrawn from them when it changes. The owner asked for it on 2026-09-28 as
an artifact to reason about workflow changes against. GitHub renders the
Mermaid; the board shows it as a code block.

### From a proposal to the target

```mermaid
flowchart TB
  S["Shaping: /grove-shape TOPIC
interactive, or headless on worktree-shape-SLUG
writes proposed work, questions, decisions, terms"]
  S --> A{"Owner assigns"}
  A -->|"R on the board, or grove run ID..."| RUN["Grove-owned attempt
claude -p /grove-work IDs --interaction headless
one process, one worktree, one budget
budget, mode, model, effort from grove.yaml run:"]
  A -->|"/grove-work IDs in a terminal"| INT["Interactive session
no attempt files: visible as branch, worktree, Next"]
  RUN --> W["One attempt (next diagram)"]
  INT --> W
  W -->|"stops at the plan (--until plan)"| PR["plan ready: read it, then R"]
  PR -->|"R without the bound"| RUN
  W -->|"stops on a question"| Q["waiting on question
e answers it in $EDITOR, resolves, commits"]
  Q -->|"question answered: R again"| RUN
  W -->|"hands off"| REV{"Review
record status=review, candidate=COMMIT"}
  REV -->|"a approve, then i integrate"| DONE["done on the target
written after the merge, committed alone"]
  REV -->|"f feedback: back to active"| A
  REV -->|"conflict with the target: m or grove resolve"| RES["feedback naming the target commit
one attempt merges it and hands off again"]
  RES --> W
  REV -->|"grove sweep under policy:"| SW["skip, wait, resolve, or
approve after verify in a temp worktree, then integrate
attributed: delegated under policy"]
  SW --> DONE
```

Sweep runs only when a person or a scheduler runs `grove sweep`; a finishing
attempt does not start one. Several IDs in one launch are one attempt on one
branch and hand off together on one shared candidate.

### Inside one attempt

```mermaid
flowchart TB
  C1["1 Assemble context
grove context IDs: selected records in full, the rest listed"]
  C1 --> C2["2 Inspect the real state
git status, branches, worktrees, versions, checkpoints in Next"]
  C2 -->|"nothing can start and the wait is recorded"| E0(["return the wait, write nothing"])
  C2 --> C3["3 Execution checkout
reuse the branch's worktree, or a new worktree-ID from the base"]
  C3 --> C4["4 Prepare
read or write the plan record, commit it"]
  C4 -->|"--until plan"| E1(["checkpoint the continuation in Next
status unchanged: plan ready"])
  C4 --> C5["5 Implement through evidence
status=active; units one at a time in deps order
routine technical choices are the session's"]
  C5 -->|"a consequential choice the record does not make"| E2(["question with blocks, checkpoint, commit
finish the units that do not depend on it"])
  C5 --> C6["6 Review gate
fresh grove-reviewer subagent, read-only
findings with evidence; at most 3 fix rounds"]
  C6 -->|"findings"| C5
  C6 -->|"cap reached"| E3(["stay active with the open findings"])
  C6 --> C7["7 Checkpoint
Next holds branch, revision, evidence, pending judgments"]
  C7 --> C8["8 Hand off
Evidence and Next reconciled; review record ends Open findings: none or N
status=review candidate=COMMIT, that change committed alone"]
  C8 --> E4(["the owner, or the policy, judges the candidate"])
```

Budget exhaustion or a stop ends everything at any step with what is
committed kept. A started unit left incomplete holds every unit on the branch
out of review; a unit that never started keeps its wait. The reviewer is the
only subagent; every other step runs in the one session.
