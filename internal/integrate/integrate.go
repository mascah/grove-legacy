// Package integrate delivers an accepted candidate to the configured target
// as one squash commit (G-260929-gm3m4), a sequence of separately reported
// facts: acceptance found, evidence retained, delivery made or refused,
// delivery proved, cleanup done or kept. Every refusal happens before the
// target moves, and nothing after the delivery undoes it; no record is
// written, since Done is derived from the acceptance the delivery brings to
// the target (G-260930-2qa4a, G-260930-gj9d7). A candidate shared by several
// work records on its branch (G-260925-wc2pz) is delivered as their group:
// every member must be accepted. Rerun after an interruption, it reads a
// delivery already made as every reader does, from the target's copy of the
// record, and goes on to cleanup, never a second commit.
package integrate

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/repo"
	"github.com/mascah/grove/internal/standing"
	"github.com/mascah/grove/internal/update"
	"github.com/mascah/grove/internal/versions"
)

// Request names the work to integrate from Root, the target's checkout. Cwd
// is the process's directory, which cleanup keeps out of any removed
// worktree.
type Request struct {
	Root, ID, Cwd string
	Cleanup       bool
	// Expect, when set, is the target commit a delegated integration verified
	// the merge against: a target at any other commit is refused unmerged.
	Expect string
	// Policy, when set, attributes the integration to a standing policy
	// (G-260925-wh9ax): the delivery's message names it and its revert.
	Policy string
}

// Run integrates the work, reporting each fact to report as it holds. The
// error is the refusal or failure that stopped the sequence; facts already
// reported stand.
func Run(req Request, now time.Time, report func(fact string)) error {
	root := req.Root
	ctx := context.Background()
	p, ds := project.Load(root, root)
	if len(ds) != 0 {
		return fmt.Errorf("the project is not valid; fix it before integrating:\n%s", diagnostics(ds))
	}
	if p.Target == "" {
		return errors.New("integration needs target: BRANCH in grove.yaml, the branch accepted work is delivered to")
	}
	branch, err := update.Branch(root)
	if err != nil {
		return fmt.Errorf("this checkout's branch could not be read: %v", err)
	}
	if branch != p.Target {
		return fmt.Errorf("integrate runs in the checkout of the target %s; this one is on %s", p.Target, orNoBranch(branch))
	}
	if dirty, err := repo.Git(root, "status", "--porcelain", "--untracked-files=no"); err != nil {
		return err
	} else if dirty != "" {
		return fmt.Errorf("the checkout of %s has uncommitted changes; commit or set them aside before delivering", p.Target)
	}
	// Every record, not only this one: the branch's other members of the
	// group are read in the same pass.
	res, err := versions.Inspect(root, "")
	if err != nil {
		return err
	}
	from, r, err := accepted(ctx, root, res, req.ID, p.Target)
	if errors.As(err, new(noBranch)) {
		// No branch holds the work at all: the target's own record may
		// say why.
		switch s := standing.Inspect(ctx, root, p.Target, p.Records)[req.ID]; {
		case s != nil && s.State == standing.Done && s.Legacy:
			return fmt.Errorf("%s is done under schema 3 here; there is nothing to deliver", req.ID)
		case s != nil && s.State == standing.Done:
			report(fmt.Sprintf("delivered: %s is already %s", req.ID, s.Text()))
			return nil
		}
	} else if err != nil {
		if s := standing.Inspect(ctx, root, p.Target, p.Records)[req.ID]; s != nil && s.State == standing.Done {
			return fmt.Errorf("%v (%s is %s)", err, req.ID, s.Text())
		}
	}
	if err != nil {
		return err
	}
	name, submitted := strings.TrimPrefix(from.Ref, "refs/heads/"), from.Commit
	// The branch's own acceptance decides, read against the target's copy as
	// every reader reads it: a delivery already made, as by an interrupted
	// run or someone else, leaves only cleanup, and work reopened and
	// accepted again on a kept branch is a new delivery.
	switch s := standing.Each(ctx, root, p.Target, []*project.Record{r})[r]; s.State {
	case standing.Done:
		report(fmt.Sprintf("delivered: %s is already %s", req.ID, s.Text()))
		if !req.Cleanup {
			return nil
		}
		// The squash commit whose submitted tip is the branch's, if one is.
		// ponytail: walks the target's first parents until it finds one, only
		// for a cleanup of work already delivered.
		d, err := repo.Git(root, "log", "-1", "--first-parent", "--format=%H", "--fixed-strings", "--grep=Grove-Submitted: "+submitted, "HEAD")
		if err != nil {
			return err
		}
		if strings.TrimSpace(d) == "" {
			submitted = ""
		}
		return cleanup(root, req.Cwd, name, submitted, worktreeOf(res, from.Ref), report)
	case standing.Unknown:
		return fmt.Errorf("%s's delivery cannot be decided (%s); nothing was delivered", req.ID, s.Why)
	}
	if _, err := repo.Git(root, "merge-base", "--is-ancestor", r.Candidate, submitted); err != nil {
		return fmt.Errorf("branch %s does not contain candidate %s, which it names; repair the record before integrating", name, r.Candidate)
	}
	group := update.Group(branchRecords(res, from), r)
	var waiting, paths []string
	for _, m := range group {
		paths = append(paths, m.Path)
		if s := standing.Of(m); m.Status != "accepted" || s.State != standing.Unknown {
			waiting = append(waiting, m.ID+" is "+strings.TrimPrefix(s.Text(), "accepted; delivery unknown: delivery not examined"))
		}
	}
	if waiting != nil {
		where := "its checkout"
		if w := worktreeOf(res, from.Ref); w != "" {
			where = w
		}
		return fmt.Errorf("candidate %s on %s is shared by %s, and delivering it delivers all of them, but %s; judge each first (grove approve ID VERDICT in %s); %s is unchanged", short(r.Candidate), name, ids(group), strings.Join(waiting, ", "), where, p.Target)
	}
	// The squash carries every commit the branch holds, so other work whose
	// candidate it contains would arrive with it. The target's own copy of
	// that work decides, as for every reader: done for this candidate, it
	// arrived by its own delivery; done for any other, the branch is behind
	// the target; not done, it needs its own delivery first, or, accepted on
	// this branch alone, to go with this work as one candidate; unaccepted, it
	// would arrive unaccepted.
	for _, o := range branchRecords(res, from) {
		if o.Type != "work" || o.Candidate == "" || slices.Contains(group, o) || o.Status == "done" || o.Status == "abandoned" {
			continue
		}
		carried, err := ancestor(root, o.Candidate, submitted)
		if err != nil {
			return err
		}
		landed, err := ancestor(root, o.Candidate, "HEAD")
		if err != nil {
			return err
		}
		if !carried || landed {
			continue
		}
		read := []*project.Record{o}
		if i := slices.IndexFunc(p.Records, func(t *project.Record) bool { return t.ID == o.ID }); i >= 0 {
			read = append(read, p.Records[i])
		}
		st := standing.Each(ctx, root, p.Target, read)
		carry := fmt.Sprintf("delivering %s would also carry %s's candidate %s", name, o.ID, short(o.Candidate))
		switch {
		case st[o].State == standing.Done:
			continue
		case len(read) == 2 && st[read[1]].State == standing.Done:
			return fmt.Errorf("%s, but %s already holds %s done, from candidate %s: %s is behind it; merge %s into %s and hand the work off again; %s is unchanged", carry, p.Target, o.ID, short(read[1].Candidate), name, p.Target, name, p.Target)
		case o.Status != "accepted":
			return fmt.Errorf("%s, which is %s without an acceptance; hand it off and judge it with %s, or move it off %s; %s is unchanged", carry, o.Status, ids(group), name, p.Target)
		}
		switch f, _, err := accepted(ctx, root, res, o.ID, p.Target); {
		case err != nil:
			return fmt.Errorf("%s, accepted and not yet delivered to %s; integrate %s first, which is refused now: %v; %s is unchanged", carry, p.Target, o.ID, err, p.Target)
		case f != from:
			return fmt.Errorf("%s, accepted and not yet delivered to %s; integrate %s first, from %s; %s is unchanged", carry, p.Target, o.ID, strings.TrimPrefix(f.Ref, "refs/heads/"), p.Target)
		default:
			return fmt.Errorf("%s, accepted and not yet delivered to %s, and no other branch delivers %s without %s; reopen each with grove feedback, hand them off as one candidate and judge them again; %s is unchanged", carry, p.Target, o.ID, ids(group), p.Target)
		}
	}
	if others, err := versions.Others(ctx, root, r.Candidate, submitted, paths...); err != nil {
		return err
	} else if len(others) != 0 {
		return fmt.Errorf("commits after candidate %s on %s change %s: the tip %s is a new candidate; accept it before integrating", short(r.Candidate), name, strings.Join(others, ", "), short(submitted))
	}
	for _, m := range group {
		report(fmt.Sprintf("acceptance: candidate %s of %s accepted by %s on branch %s at %s%s", short(r.Candidate), m.ID, m.ApprovedBy, name, short(submitted), verdict(m)))
	}

	before, err := head(root)
	if err != nil {
		return err
	}
	if req.Expect != "" && before != req.Expect {
		return fmt.Errorf("%s moved from %s, where the merge was verified, to %s; nothing was delivered and %s stays accepted", p.Target, short(req.Expect), short(before), req.ID)
	}
	// A conflict is refused before anything is written (G-260925-h8rj5):
	// merge-tree performs it in objects only and names the files.
	ms, err := versions.PredictContext(ctx, root, before, []string{submitted})
	if err != nil {
		return fmt.Errorf("the delivery of %s into %s could not be prepared: %v; nothing was delivered", name, p.Target, err)
	}
	if ms[0].Outcome == "conflict" {
		conflict := ms[0].Text(p.Target)
		where := worktreeOf(res, from.Ref)
		if where == "" {
			where = "a checkout of " + name
		}
		return fmt.Errorf("delivery of %s into %s refused: it %s; nothing was delivered, %s is unchanged at %s and %s stays accepted. Next: grove resolve %s records that as feedback and starts one attempt to merge %s at %s, resolve and hand off a new candidate; or by hand, in %s, git merge %s, resolve the conflicts and commit, then hand that commit to review as the new candidate",
			name, p.Target, conflict, p.Target, short(before), req.ID, req.ID, p.Target, short(before), where, p.Target)
	}
	if ms[0].Outcome == "integrated" {
		return fmt.Errorf("%s already holds %s at %s, yet its own copy of %s is not accepted for candidate %s: it changed there since; inspect it with grove show %s", p.Target, name, short(before), req.ID, short(r.Candidate), req.ID)
	}
	tree, err := repo.Git(root, "merge-tree", "--write-tree", "--no-messages", before, submitted)
	if err != nil {
		return fmt.Errorf("the delivery of %s into %s could not be prepared: %v; nothing was delivered", name, p.Target, err)
	}
	if current, err := repo.Git(root, "rev-parse", before+"^{tree}"); err != nil {
		return err
	} else if strings.TrimSpace(current) == strings.TrimSpace(tree) {
		return fmt.Errorf("%s at %s already holds everything %s would deliver, yet its own copy of %s is not accepted for candidate %s. Nothing was delivered; inspect it with grove show %s", p.Target, short(before), name, req.ID, short(r.Candidate), req.ID)
	}
	// The evidence is retained before the target moves, so a delivery that
	// succeeds can always be verified, whatever happens to the branch; and
	// only once nothing refuses it, so no refusal leaves a ref behind.
	if _, err := repo.Git(root, "update-ref", standing.Ref(submitted), submitted); err != nil {
		return fmt.Errorf("the submitted tip %s could not be retained: %v; nothing was delivered", short(submitted), err)
	}
	report("retained: " + standing.Ref(submitted))
	message, err := Message(root, p.RecordDir, before, r.Candidate, submitted, group, req.Policy)
	if err != nil {
		return err
	}
	delivered, err := repo.Git(root, "commit-tree", strings.TrimSpace(tree), "-p", before, "-m", message)
	if err != nil {
		return fmt.Errorf("the delivery commit could not be made: %v; nothing was delivered", err)
	}
	delivered = strings.TrimSpace(delivered)
	// Only a fast-forward from the commit prepared against, so a target that
	// moved meanwhile is refused rather than overwritten.
	if out, err := repo.Command(ctx, root, "merge", "--ff-only", "-q", delivered).CombinedOutput(); err != nil {
		return fmt.Errorf("%s could not be advanced to the delivery %s: %s; nothing was delivered, and the evidence stays retained", p.Target, short(delivered), strings.TrimSpace(string(out)))
	}
	report(fmt.Sprintf("delivery: squash commit %s on %s (was %s)", short(delivered), p.Target, short(before)))
	if req.Policy != "" {
		report(fmt.Sprintf("delivered under %s; to reverse it: git revert %s", req.Policy, delivered))
	}
	// The one commit just made is proved as an audit proves it, and each
	// member then reads done as every reader reads it.
	after, ds := project.Load(root, root)
	if len(ds) != 0 {
		return fmt.Errorf("delivered as %s, but the project there does not load: %s", short(delivered), diagnostics(ds))
	}
	var members []*project.Record
	for _, m := range group {
		i := slices.IndexFunc(after.Records, func(o *project.Record) bool { return o.ID == m.ID && o.Type == "work" })
		if i < 0 {
			return fmt.Errorf("delivered as %s, but %s is not a work record there; inspect it before anything else", short(delivered), m.ID)
		}
		members = append(members, after.Records[i])
	}
	proofs, err := standing.Prove(ctx, root, delivered, members, after.Records)
	if err != nil {
		return fmt.Errorf("delivered as %s, but it could not be proved: %v; grove check --deliveries audits it", short(delivered), err)
	}
	read := standing.Each(ctx, root, p.Target, members)
	for i, m := range members {
		if !proofs[i].Proved {
			return fmt.Errorf("delivered as %s, but %s is %s; inspect it before anything else", short(delivered), m.ID, proofs[i].Text())
		}
		if s := read[m]; s.State != standing.Done {
			return fmt.Errorf("delivered as %s, but %s reads %s there; inspect it before anything else", short(delivered), m.ID, s.Text())
		}
		report(fmt.Sprintf("done: %s is %s, %s", m.ID, read[m].Text(), proofs[i].Text()))
	}
	if !req.Cleanup {
		return nil
	}
	return cleanup(root, req.Cwd, name, submitted, worktreeOf(res, from.Ref), report)
}

// types orders Conventional Commit types from most to least significant.
var types = []string{"feat", "fix", "perf", "refactor", "revert", "docs", "test", "build", "ci", "style", "chore"}

// Message is the delivery's commit message: a Conventional Commit whose
// type is the most significant among the candidate's own commits since the
// target (those changing anything outside the record root, else all), "!"
// when any is breaking, and whose subject is the first member's title; the
// body lists the members, then the trailers that locate the evidence. It
// never names the delivery's own commit, so nothing refers to itself.
func Message(root, records, base, candidate, submitted string, group []*project.Record, policy string) (string, error) {
	subjects, err := repo.Git(root, "log", "--no-merges", "--format=%s", base+".."+candidate, "--", ":/", ":(exclude)"+records)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(subjects) == "" {
		if subjects, err = repo.Git(root, "log", "--no-merges", "--format=%s", base+".."+candidate); err != nil {
			return "", err
		}
	}
	kind, breaking := "chore", false
	for l := range strings.SplitSeq(subjects, "\n") {
		prefix, _, ok := strings.Cut(l, ": ")
		if !ok {
			continue
		}
		if strings.HasSuffix(prefix, "!") {
			breaking, prefix = true, strings.TrimSuffix(prefix, "!")
		}
		if i := strings.IndexByte(prefix, '('); i > 0 {
			prefix = prefix[:i]
		}
		if i, j := slices.Index(types, prefix), slices.Index(types, kind); i >= 0 && i < j {
			kind = prefix
		}
	}
	if breaking {
		kind += "!"
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%s: %s\n\n", kind, subject(group[0].Title))
	for _, m := range group {
		fmt.Fprintf(&b, "- %s %s\n", m.ID, m.Title)
	}
	if policy != "" {
		fmt.Fprintf(&b, "\nIntegrated under %s.\n", policy)
	}
	b.WriteString("\n")
	for _, m := range group {
		fmt.Fprintf(&b, "Grove-Work: %s\n", m.ID)
	}
	full, err := repo.Git(root, "rev-parse", candidate+"^{commit}")
	if err != nil {
		return "", err
	}
	fmt.Fprintf(&b, "Grove-Candidate: %s\nGrove-Submitted: %s\n", strings.TrimSpace(full), submitted)
	return b.String(), nil
}

// subject is a title as a commit subject: lowercase first, unless its first
// word is an acronym.
func subject(title string) string {
	first, n := utf8.DecodeRuneInString(title)
	if next, _ := utf8.DecodeRuneInString(title[n:]); n == len(title) || unicode.IsUpper(next) {
		return title
	}
	return string(unicode.ToLower(first)) + title[n:]
}

// accepted finds the one branch holding the record accepted, its acceptance
// applicable. The target itself never counts: what it holds is what
// delivery produces, not a submission. Of different candidates, the one
// containing every other is the latest acceptance, and candidates none of
// which contains the others are refused, naming each. Of branches accepting
// the same
// candidate, as when later work was based on this work's branch, it delivers
// from the one whose commits after the candidate change only the records
// sharing it, as delivery requires, or
// among several such, the one every other contains.
func accepted(ctx context.Context, root string, res *versions.Result, id, target string) (*versions.Source, *project.Record, error) {
	var judged, ok []*versions.Version
	for i := range res.Groups {
		if res.Groups[i].ID != id {
			continue
		}
		for j := range res.Groups[i].Versions {
			v := &res.Groups[i].Versions[j]
			if v.Source.Kind != "committed" || v.Source.Ref == "refs/heads/"+target || v.Record == nil || v.Record.Type != "work" || v.Record.Status != "review" && v.Record.Status != "accepted" {
				continue
			}
			judged = append(judged, v)
			if standing.Of(v.Record).State == standing.Unknown {
				ok = append(ok, v)
			}
		}
	}
	names := func(vs []*versions.Version) string {
		var out []string
		for _, v := range vs {
			out = append(out, strings.TrimPrefix(v.Source.Ref, "refs/heads/"))
		}
		return strings.Join(out, ", ")
	}
	switch {
	case len(ok) == 1:
		return ok[0].Source, ok[0].Record, nil
	case len(ok) > 1:
		// A later acceptance of work that contains an earlier one supersedes
		// it; acceptances that do not line up are named, not chosen between.
		var newest []*versions.Version
		for _, v := range ok {
			all := true
			for _, u := range ok {
				in, err := ancestor(root, u.Record.Candidate, v.Record.Candidate)
				all = all && err == nil && in
			}
			if all {
				newest = append(newest, v)
			}
		}
		if len(newest) == 0 {
			var at []string
			for _, v := range ok {
				at = append(at, strings.TrimPrefix(v.Source.Ref, "refs/heads/")+" for "+short(v.Record.Candidate))
			}
			return nil, nil, fmt.Errorf("%s is accepted for candidates on several branches (%s), none containing the others; integrate needs one: merge the branches, or reopen the acceptance that no longer stands", id, strings.Join(at, ", "))
		}
		ok = newest
		if len(ok) == 1 {
			return ok[0].Source, ok[0].Record, nil
		}
		var clean []*versions.Version
		for _, v := range ok {
			var paths []string
			for _, m := range update.Group(branchRecords(res, v.Source), v.Record) {
				paths = append(paths, m.Path)
			}
			others, err := versions.Others(ctx, root, v.Record.Candidate, v.Source.Commit, paths...)
			if err == nil && len(others) == 0 {
				clean = append(clean, v)
			}
		}
		if len(clean) == 1 {
			return clean[0].Source, clean[0].Record, nil
		}
		if v := earliest(root, clean); v != nil {
			return v.Source, v.Record, nil
		}
		return nil, nil, fmt.Errorf("%s is accepted on several branches (%s); integrate needs one", id, names(ok))
	case len(judged) != 0:
		v := judged[0]
		where := "its checkout"
		if w := worktreeOf(res, v.Source.Ref); w != "" {
			where = w
		}
		return nil, nil, fmt.Errorf("%s is in review on %s but not accepted, or its acceptance no longer applies: run grove approve %s VERDICT in %s first", id, names(judged), id, where)
	}
	return nil, nil, noBranch(id)
}

// earliest is the one of vs, all accepting the same candidate, whose branch
// tip every other contains, or nil.
func earliest(root string, vs []*versions.Version) *versions.Version {
	for _, v := range vs {
		all := true
		for _, u := range vs {
			in, err := ancestor(root, v.Source.Commit, u.Source.Commit)
			all = all && err == nil && in
		}
		if all {
			return v
		}
	}
	return nil
}

// noBranch is accepted's answer when no branch holds the work in review or
// accepted.
type noBranch string

func (id noBranch) Error() string {
	return "no branch holds " + string(id) + " accepted; nothing to integrate"
}

// branchRecords is every record the committed source holds.
func branchRecords(res *versions.Result, from *versions.Source) []*project.Record {
	var out []*project.Record
	for _, g := range res.Groups {
		for _, v := range g.Versions {
			if v.Source == from && v.Record != nil {
				out = append(out, v.Record)
			}
		}
	}
	return out
}

func ids(records []*project.Record) string {
	var out []string
	for _, r := range records {
		out = append(out, r.ID)
	}
	return strings.Join(out, ", ")
}

// worktreeOf is the registered checkout of a branch, or "".
func worktreeOf(res *versions.Result, ref string) string {
	for _, s := range res.Sources {
		if s.Kind == "live" && s.Ref == ref {
			return s.Worktree
		}
	}
	return ""
}

// verdict quotes the verdict paragraph the approval appended, when the body
// still holds one for this candidate.
func verdict(r *project.Record) string {
	prefix := "Verdict on candidate " + short(r.Candidate) + ","
	found := ""
	for _, l := range strings.Split(string(r.Source), "\n") {
		if strings.HasPrefix(l, prefix) {
			found = strings.TrimRight(l, "\r")
		}
	}
	if found == "" {
		return ""
	}
	return " (" + found + ")"
}

// cleanup removes the branch's worktree, then the branch, through Git's own
// refusals: a worktree with changes or untracked files is kept, and so is a
// branch whose deletion would lose a commit. A worktree holding cwd is kept
// too, and one holding ignored files, which git worktree remove would delete.
// submitted is the tip a squash delivery retained, "" after an ordinary merge.
// Whatever is kept is reported and makes the result an error, since the
// caller asked for a cleanup that did not fully happen.
func cleanup(root, cwd, name, submitted, worktree string, report func(string)) error {
	kept := false
	keep := func(what, reason string) {
		kept = true
		report(fmt.Sprintf("cleanup: kept %s: %s", what, reason))
	}
	// Only a branch nothing is lost by deleting: after a squash, one still at
	// the submitted tip, since one that moved on holds a later edit, and only
	// once the evidence ref retains it, written here before anything is
	// removed, since a clone that fetched the branch without refs/grove
	// verifies the delivery from the branch alone; after an ordinary merge,
	// which retains nothing, one whose tip the target contains, since a merge
	// of the candidate alone leaves the handoff and acceptance after it on the
	// branch only.
	tip, err := repo.Git(root, "rev-parse", "-q", "--verify", "refs/heads/"+name)
	tip = strings.TrimSpace(tip)
	if err == nil && submitted == "" {
		if in, err := ancestor(root, tip, "HEAD"); err != nil || !in {
			keep("worktree and branch "+name, "the target does not contain its tip "+short(tip)+", whose commits nothing else retains")
			return errors.New("cleanup incomplete; the integration stands")
		}
		submitted = tip
	} else if err == nil && tip == submitted {
		if _, err := repo.Git(root, "update-ref", standing.Ref(submitted), submitted); err != nil {
			keep("worktree and branch "+name, "its tip "+short(tip)+" could not be retained under "+standing.Ref(submitted)+": "+err.Error())
			return errors.New("cleanup incomplete; the integration stands")
		}
	}
	if err != nil || tip != submitted {
		keep("worktree and branch "+name, "the branch is not at a tip that a squash delivery retained")
		return errors.New("cleanup incomplete; the integration stands")
	}
	if worktree != "" {
		if rel, err := filepath.Rel(worktree, cwd); cwd != "" && err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			keep("worktree "+worktree, "it holds this process's working directory")
		} else if ignored, err := ignored(worktree); err != nil {
			keep("worktree "+worktree, err.Error())
		} else if len(ignored) != 0 {
			keep("worktree "+worktree, "it holds ignored files ("+strings.Join(ignored, ", ")+"); remove them or the worktree by hand")
		} else if _, err := repo.Git(root, "worktree", "remove", "--", worktree); err != nil {
			keep("worktree "+worktree, err.Error())
		} else {
			report("cleanup: removed worktree " + worktree)
		}
	}
	// A branch a kept worktree has checked out stays with it; update-ref
	// deletes only at the tip checked above, should it move meanwhile.
	if kept {
		keep("branch "+name, "its worktree is kept")
	} else if _, err := repo.Git(root, "update-ref", "-d", "refs/heads/"+name, submitted); err != nil {
		keep("branch "+name, "it moved past the delivered tip "+short(submitted)+", or could not be deleted: "+err.Error())
	} else {
		report("cleanup: deleted branch " + name)
	}
	if kept {
		return errors.New("cleanup incomplete; the integration stands")
	}
	return nil
}

// ignored lists a checkout's ignored files, which git worktree remove deletes
// without --force; a .env or a build directory is not Grove's to delete.
func ignored(worktree string) ([]string, error) {
	out, err := repo.Git(worktree, "status", "--porcelain", "--ignored", "-z")
	if err != nil {
		return nil, err
	}
	var files []string
	for _, entry := range strings.Split(strings.TrimSuffix(out, "\x00"), "\x00") {
		if strings.HasPrefix(entry, "!! ") {
			files = append(files, entry[3:])
		}
	}
	return files, nil
}

// ancestor reports whether commit is an ancestor of ref in root's repository.
func ancestor(root, commit, ref string) (bool, error) {
	_, err := repo.Git(root, "merge-base", "--is-ancestor", commit, ref)
	var exit *exec.ExitError
	if errors.As(err, &exit) && exit.ExitCode() == 1 {
		return false, nil
	}
	return err == nil, err
}

func head(root string) (string, error) {
	out, err := repo.Git(root, "rev-parse", "HEAD")
	return strings.TrimSpace(out), err
}

func orNoBranch(branch string) string {
	if branch == "" {
		return "no branch"
	}
	return branch
}

func short(commit string) string { return commit[:min(len(commit), 7)] }

func diagnostics(ds []project.Diagnostic) string {
	lines := make([]string, len(ds))
	for i, d := range ds {
		lines[i] = d.String()
	}
	return strings.Join(lines, "\n")
}
