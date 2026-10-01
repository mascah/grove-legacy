package update

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/repo"
	"github.com/mascah/grove/internal/standing"
	"github.com/mascah/grove/internal/versions"
)

// The owner's two dispositions of a candidate (G-260921-jwk4e), each an
// ordinary update of the record committed alone in the checkout of the branch
// that holds it: approval records the acceptance (G-260930-2qa4a), status
// accepted with the candidate, its authority and the acceptance context,
// and quotes the verdict in the body; feedback reopens the work, before or
// after delivery, with the text in the body and no acceptance left behind. Neither touches any other file, ref,
// or worktree.
//
// A candidate can be shared (G-260925-wc2pz): when an explicitly selected set of work
// is handed off together, the records on the branch whose candidate is the
// same commit are one group. Approval stays per record, and ignores the
// other members' own record commits; feedback on any member reopens the
// group, since the next candidate replaces the shared one.

// Owner is approved_by for the owner's own acceptance; a sweep's is Policy's.
const Owner = "owner"

// Policy is approved_by for an acceptance under the grove.yaml revision given.
func Policy(revision string) string { return "policy " + revision }

// Approve records the verdict on the record's candidate as its acceptance by
// by, Owner or a Policy: status accepted, approved the candidate, approved_by
// and the acceptance context, with the verdict appended to the body,
// committed alone. The checkout must hold the candidate with nothing but the
// record changed since it, and the record's file must match HEAD.
func Approve(root, id, verdict, by string, now time.Time) (Result, error) {
	if verdict = strings.TrimSpace(verdict); verdict == "" {
		return Result{}, errors.New("a verdict is required: the owner's words on the candidate")
	}
	p, r, err := judged(root, id)
	if err != nil {
		return Result{}, err
	}
	// A stale acceptance, whose record changed since, is judged again.
	if r.Status == "accepted" && r.ApprovedContext == project.AcceptanceContext(r) {
		return Result{}, fmt.Errorf("%s: candidate %s is already accepted; integrate it", id, r.Candidate)
	}
	if err := holdsCandidate(root, r, Group(p.Records, r)); err != nil {
		return Result{}, err
	}
	return Apply(root, Request{
		ID: id, Expect: project.Revision(r.Source), Commit: true,
		Set:    []Field{{"status", "accepted"}, {"approved", r.Candidate}, {"approved_by", by}, {"approved_context", project.AcceptanceContext(r)}},
		Append: fmt.Sprintf("Verdict on candidate %s, %s: %s", short(r.Candidate), now.UTC().Format("2006-01-02"), verdict),
	}, now, nil)
}

// Feedback returns the record, in review or accepted, delivered or not, to
// active with the owner's text appended to the body, committed alone; an
// acceptance is removed in the same update and the candidate stays, so the
// earlier reviews still compare to it. Every other member of its group in
// review or accepted is returned to active the same way, with a line naming
// this feedback instead of the text, each committed alone and listed in
// Reopened.
func Feedback(root, id, text string, now time.Time) (Result, error) {
	if text = strings.TrimSpace(text); text == "" {
		return Result{}, errors.New("feedback text is required: what the next attempt must change")
	}
	p, r, err := judged(root, id)
	if err != nil {
		return Result{}, err
	}
	var others []*project.Record
	for _, o := range Group(p.Records, r) {
		if o != r && (o.Status == "review" || o.Status == "accepted") {
			others = append(others, o)
		}
	}
	// Every file is checked before the first is written.
	for _, o := range append([]*project.Record{r}, others...) {
		if err := clean(root, o); err != nil {
			return Result{}, err
		}
	}
	if err := holdsCandidate(root, r, nil); err != nil {
		// A squash delivery never puts the candidate in the target, so
		// delivered work reopens where its delivery is.
		if !deliveredHere(root, p, r) {
			return Result{}, err
		}
		if err := clean(root, r); err != nil {
			return Result{}, err
		}
	}
	reopen := func(o *project.Record, line string) (Result, error) {
		req := Request{ID: o.ID, Expect: project.Revision(o.Source), Commit: true, Set: []Field{{"status", "active"}}, Append: line}
		if o.Approved != "" {
			req.Unset = []string{"approved", "approved_by", "approved_context"}
		}
		return Apply(root, req, now, nil)
	}
	day := now.UTC().Format("2006-01-02")
	res, err := reopen(r, fmt.Sprintf("Feedback on candidate %s, %s: %s", short(r.Candidate), day, text))
	if err != nil {
		return res, err
	}
	if branch, _ := Branch(root); p.Target != "" && branch == p.Target {
		res.OnTarget = true
	}
	for _, o := range others {
		done, err := reopen(o, fmt.Sprintf("Reopened with %s's feedback on candidate %s, %s: the next candidate replaces the one this group shared.", id, short(r.Candidate), day))
		if err != nil {
			return res, fmt.Errorf("%s is active again at commit %s, but reopening %s, which shares its candidate, failed: %v; run grove feedback %s with the same text", id, short(res.Commit), o.ID, err, o.ID)
		}
		res.Reopened = append(res.Reopened, done)
	}
	return res, nil
}

// Group is the work in records whose candidate is r's commit, r included,
// ordered by ID (G-260925-wc2pz): the members one shared candidate was handed off
// for. A record without a candidate is a group of none.
func Group(records []*project.Record, r *project.Record) []*project.Record {
	var group []*project.Record
	for _, o := range records {
		if o.Type == "work" && r.Candidate != "" && SameCommit(o.Candidate, r.Candidate) {
			group = append(group, o)
		}
	}
	slices.SortFunc(group, func(a, b *project.Record) int { return strings.Compare(a.ID, b.ID) })
	return group
}

// SameCommit compares two recorded commits, either of which may be
// abbreviated.
func SameCommit(a, b string) bool {
	return a != "" && b != "" && (strings.HasPrefix(a, b) || strings.HasPrefix(b, a))
}

// judged loads the project and returns the work record in review or
// accepted. A checkout holding a delivery is retired (standing.Retired) and
// nothing is judged in it, whichever work the delivery was of: an acceptance
// there is never delivered, and a reopening would continue there.
func judged(root, id string) (*project.Project, *project.Record, error) {
	p, ds := project.Load(root, root)
	if len(ds) != 0 {
		return nil, nil, fmt.Errorf("the project is not valid; fix it before judging a candidate:\n%s", diagnostics(ds))
	}
	i := slices.IndexFunc(p.Records, func(r *project.Record) bool { return r.ID == id })
	if i < 0 {
		return nil, nil, fmt.Errorf("record %s not found in this project", id)
	}
	r := p.Records[i]
	if r.Type != "work" {
		return nil, nil, fmt.Errorf("%s is a %s, not work: only work has a candidate", id, r.Type)
	}
	if r.Status != "review" && r.Status != "accepted" {
		return nil, nil, fmt.Errorf("%s is %s, not in review: there is no candidate awaiting judgment", id, r.Status)
	}
	if p.Target != "" {
		d, err := standing.Retired(context.Background(), root, p.Target, "HEAD")
		if err != nil {
			return nil, nil, err
		}
		if d != nil {
			branch, _ := Branch(root)
			return nil, nil, fmt.Errorf("%s: a delivered workspace is kept for inspection and nothing is judged in it; run %s fresh from %s, reopening it there first (grove feedback) if it was delivered", d.Holds(cmp.Or(branch, "this checkout")), id, p.Target)
		}
	}
	return p, r, nil
}

// holdsCandidate checks that this checkout is the one to judge the record in:
// HEAD contains the candidate, the record's file matches HEAD (these commands
// commit that file, and never someone's uncommitted edit), and, when group
// names the records sharing the candidate, nothing but their files changed
// since it, since a later commit is a new candidate that needs its own
// judgment.
func holdsCandidate(root string, r *project.Record, group []*project.Record) error {
	if _, err := repo.Git(root, "merge-base", "--is-ancestor", r.Candidate, "HEAD"); err != nil {
		var exit *exec.ExitError
		if errors.As(err, &exit) && exit.ExitCode() == 1 {
			return fmt.Errorf("this checkout does not hold candidate %s; run this in the checkout of the branch that has %s in review", r.Candidate, r.ID)
		}
		return fmt.Errorf("candidate %s could not be checked against this checkout's HEAD: %v", r.Candidate, err)
	}
	if err := clean(root, r); err != nil {
		return err
	}
	if group == nil {
		return nil
	}
	var paths []string
	for _, o := range group {
		paths = append(paths, o.Path)
	}
	others, err := versions.Others(context.Background(), root, r.Candidate, "HEAD", paths...)
	if err != nil {
		return err
	}
	if len(others) != 0 {
		tip, err := repo.Git(root, "rev-parse", "HEAD")
		if err != nil {
			return err
		}
		tip = strings.TrimSpace(tip)
		return fmt.Errorf("commits after candidate %s change %s: the tip %s is a new candidate; set candidate=%s, review it, then judge that", short(r.Candidate), strings.Join(others, ", "), short(tip), short(tip))
	}
	return nil
}

// deliveredHere reports accepted work the target holds accepted for the same
// candidate, which only a delivery brings there.
func deliveredHere(root string, p *project.Project, r *project.Record) bool {
	return r.Status == "accepted" && standing.Each(context.Background(), root, p.Target, []*project.Record{r})[r].State == standing.Done
}

// clean refuses a record whose file differs from HEAD in this checkout.
func clean(root string, r *project.Record) error {
	status, err := repo.Git(root, "status", "--porcelain", "-z", "--", ":(literal)"+filepath.FromSlash(r.Path))
	if err != nil {
		return err
	}
	if status != "" {
		return fmt.Errorf("%s has uncommitted changes in this checkout; commit or discard them first, since judging commits that file", r.Path)
	}
	return nil
}

func short(commit string) string { return commit[:min(len(commit), 7)] }

// Delegated reports whether r's acceptance was given under a standing
// policy (G-260925-wh9ax) rather than by the owner.
func Delegated(r *project.Record) bool {
	return strings.HasPrefix(r.ApprovedBy, "policy ")
}
