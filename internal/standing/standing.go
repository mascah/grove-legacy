// Package standing derives each work record's standing (G-260930-2qa4a): what
// its record says, and whether an acceptance it holds was delivered to the
// target, from Git evidence read at the target's tip. Done is never read
// from a status: an applicable acceptance is done when the tip contains its
// candidate, or a squash delivery of it that Git verifies. Missing or
// insufficient evidence is unknown, never done. Nothing here writes, fetches
// or pushes; every consumer that decides from completion reads this result.
package standing

import (
	"bytes"
	"cmp"
	"context"
	"errors"
	"fmt"
	"path"
	"regexp"
	"strings"

	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/repo"
)

// States a Standing can have. Review includes an acceptance that no longer
// applies; Done includes schema 3's claim, marked Legacy.
const (
	Proposed  = "proposed"
	Active    = "active"
	Review    = "review"
	Accepted  = "accepted" // accepted, verified not yet delivered
	Done      = "done"
	Unknown   = "unknown" // accepted; delivery could not be established either way
	Abandoned = "abandoned"
)

// Standing is one work record's recorded and derived facts.
type Standing struct {
	ID       string `json:"id"`
	Recorded string `json:"recorded"` // the record's own status, as its file says
	State    string `json:"state"`
	Legacy   bool   `json:"legacy,omitempty"` // done is schema 3's claim, not verified
	// Target and Tip are the branch and commit delivery was judged at.
	Target string `json:"target,omitempty"`
	Tip    string `json:"tip,omitempty"`
	// Delivered is the target's commit holding the result: the candidate
	// for an ordinary merge, or the squash commit.
	Delivered string `json:"delivered,omitempty"`
	// Submitted is the tip a squash delivered, retained under Ref.
	Submitted string `json:"submitted,omitempty"`
	Why       string `json:"why,omitempty"`
}

// Text is the standing as one phrase: "done at 1234567 on main".
func (s *Standing) Text() string {
	switch {
	case s.Legacy:
		return "done (schema 3 claim)"
	case s.State == Done && s.Submitted != "":
		return "done: squashed as " + short(s.Delivered) + " on " + s.Target
	case s.State == Done:
		return "done: " + s.Target + " contains " + short(s.Delivered)
	case s.State == Accepted:
		return "accepted, awaiting delivery to " + s.Target
	case s.State == Unknown:
		return "accepted; delivery unknown: " + s.Why
	case s.State == Review && s.Recorded == "accepted":
		return "awaiting judgment: " + s.Why
	case s.State == Review:
		return "awaiting judgment"
	}
	return s.State
}

// Ref retains a squash delivery's submitted tip, and with it the candidate
// and its review and acceptance, through cleanup and garbage collection.
// It is outside refs/heads, so it is evidence, never a work branch.
func Ref(submitted string) string { return "refs/grove/submitted/" + submitted }

// Of is a record's standing from its file alone, before any delivery is
// examined: an accepted record whose acceptance applies stays Unknown until
// Inspect judges it.
func Of(r *project.Record) *Standing {
	s := &Standing{ID: r.ID, Recorded: r.Status, State: r.Status}
	switch r.Status {
	case "done":
		s.State, s.Legacy = Done, true
		s.Why = "schema 3's claim, kept by migration and not verified"
	case "accepted":
		if r.ApprovedContext != project.AcceptanceContext(r) {
			s.State, s.Why = Review, "the record changed since candidate "+short(r.Candidate)+" was accepted, so the acceptance no longer applies"
		} else {
			s.State, s.Why = Unknown, "delivery not examined"
		}
	}
	return s
}

// Inspect judges every work record in records, the project whose files sit
// under root, at the target branch's tip in root's repository. It returns
// one Standing per work record, by ID. An error is only a failure to read
// Git at all; each record's own unknowns are in its Standing.
func Inspect(ctx context.Context, root, target string, records []*project.Record) (map[string]*Standing, error) {
	out := map[string]*Standing{}
	var open []*project.Record // applicable acceptances
	for _, r := range records {
		if r.Type != "work" {
			continue
		}
		s := Of(r)
		out[r.ID] = s
		if s.State == Unknown {
			s.Target = target
			open = append(open, r)
		}
	}
	if len(open) == 0 {
		return out, nil
	}
	unknown := func(why string) (map[string]*Standing, error) {
		for _, r := range open {
			out[r.ID].Why = why
		}
		return out, nil
	}
	if target == "" {
		return unknown("no target is configured in grove.yaml")
	}
	_, _, prefix, err := repo.IdentifyContext(ctx, root)
	if err != nil {
		return unknown("not in a readable Git repository")
	}
	tip, err := repo.GitContext(ctx, root, "rev-parse", "--verify", "-q", "refs/heads/"+target+"^{commit}")
	if tip = strings.TrimSpace(tip); err != nil || tip == "" {
		return unknown("the target branch " + target + " cannot be read here")
	}
	shallow, err := repo.GitContext(ctx, root, "rev-parse", "--is-shallow-repository")
	if err != nil {
		return nil, err
	}
	complete := strings.TrimSpace(shallow) == "false"
	reach, err := lines(ctx, root, nil, "rev-list", tip)
	if err != nil {
		return nil, err
	}
	reachable := map[string]bool{}
	for _, c := range reach {
		reachable[c] = true
	}
	names := []string{}
	for _, r := range open {
		names = append(names, r.Candidate)
	}
	full, err := resolve(ctx, root, names)
	if err != nil {
		return nil, err
	}
	claims, err := deliveries(ctx, root, tip)
	if err != nil {
		return nil, err
	}
	var pending []*check
	for _, r := range open {
		s := out[r.ID]
		s.Tip = tip
		c := full[r.Candidate]
		switch {
		case c == "":
			s.Why = "candidate " + short(r.Candidate) + " is not in this repository; fetch the branch or the refs/grove evidence that holds it"
		case reachable[c]:
			s.State, s.Delivered, s.Why = Done, c, ""
		default:
			claimed := false
			for _, d := range claims {
				if d.candidate == c {
					claimed = true
					pending = append(pending, &check{record: r, candidate: c, delivery: d})
				}
			}
			if !claimed {
				absent(s, complete)
			}
		}
	}
	if err := verify(ctx, root, prefix, records, pending); err != nil {
		return nil, err
	}
	// A record is done by any verified claim; otherwise a claim whose
	// submitted tip is missing leaves it unknown, and claims that all fail
	// are no delivery at all.
	for _, r := range open {
		s := out[r.ID]
		var missing, failed string
		for _, p := range pending {
			switch {
			case p.record != r || s.State == Done:
			case p.ok:
				s.State, s.Delivered, s.Submitted, s.Why = Done, p.delivery.commit, p.delivery.submitted, ""
			case p.missing:
				missing = cmp.Or(missing, p.why)
			default:
				failed = cmp.Or(failed, p.why)
			}
		}
		switch {
		case s.State != Unknown || missing == "" && failed == "":
		case missing != "":
			s.Why = missing
		default:
			absent(s, complete)
			s.Why = strings.TrimPrefix(s.Why+"; "+failed, "; ")
		}
	}
	return out, nil
}

// absent is an acceptance whose delivery the tip does not hold: accepted,
// awaiting delivery, only when the history examined is complete.
func absent(s *Standing, complete bool) {
	if complete {
		s.State, s.Why = Accepted, ""
	} else {
		s.Why = "this clone's history is shallow, so the target may hold a delivery it cannot see; fetch the full history"
	}
}

// delivery is a commit on the target whose message claims a squash delivery.
type delivery struct {
	commit, parent, tree string
	candidate, submitted string
}

// check is one claimed delivery of one record.
type check struct {
	record    *project.Record
	candidate string
	delivery  delivery
	ok        bool
	missing   bool // the submitted tip is absent: unknown, not absent
	why       string
}

var hexCommit = regexp.MustCompile("^[0-9a-f]{40}([0-9a-f]{24})?$")

// deliveries lists the commits reachable from tip whose trailers name a
// candidate and a submitted tip. The trailers only locate: verify decides.
func deliveries(ctx context.Context, root, tip string) ([]delivery, error) {
	out, err := repo.GitContext(ctx, root, "log", "-z", "--fixed-strings", "--grep=Grove-Submitted: ", "--format=%H%x1f%P%x1f%T%x1f%(trailers:only,unfold)", tip)
	if err != nil {
		return nil, err
	}
	var ds []delivery
	for _, entry := range strings.Split(out, "\x00") {
		f := strings.SplitN(strings.TrimPrefix(entry, "\n"), "\x1f", 4)
		if len(f) != 4 {
			continue
		}
		d := delivery{commit: f[0], tree: f[2]}
		if parents := strings.Fields(f[1]); len(parents) != 0 {
			d.parent = parents[0]
		}
		for _, l := range strings.Split(f[3], "\n") {
			key, value, _ := strings.Cut(l, ": ")
			switch key {
			case "Grove-Candidate":
				d.candidate = strings.TrimSpace(value)
			case "Grove-Submitted":
				d.submitted = strings.TrimSpace(value)
			}
		}
		if d.parent != "" && hexCommit.MatchString(d.candidate) && hexCommit.MatchString(d.submitted) {
			ds = append(ds, d)
		}
	}
	return ds, nil
}

// verify decides each claimed delivery from Git objects: the submitted tip
// is present; the candidate is its ancestor; only the group's records change
// between them; the record there is accepted for the candidate; and the
// claimed commit's first parent merged with the submitted tip gives exactly
// the claimed commit's tree.
func verify(ctx context.Context, root, prefix string, records []*project.Record, checks []*check) error {
	if len(checks) == 0 {
		return nil
	}
	var names []string
	for _, c := range checks {
		names = append(names, c.delivery.submitted)
	}
	present, err := resolve(ctx, root, names)
	if err != nil {
		return err
	}
	var live []*check
	for _, c := range checks {
		if present[c.delivery.submitted] == "" {
			c.missing = true
			c.why = "delivery " + short(c.delivery.commit) + " names submitted tip " + short(c.delivery.submitted) + ", which is not in this repository; fetch refs/grove/* from where it was delivered"
			continue
		}
		live = append(live, c)
	}
	if len(live) == 0 {
		return nil
	}
	var diffs, merges, blobs bytes.Buffer
	for _, c := range live {
		fmt.Fprintf(&diffs, "%s %s\n", c.delivery.submitted, c.candidate)
		fmt.Fprintf(&merges, "%s %s\n", c.delivery.parent, c.delivery.submitted)
		fmt.Fprintf(&blobs, "%s:%s\n", c.delivery.submitted, path.Join(prefix, c.record.Path))
	}
	// ponytail: one rev-list per claimed squash delivery, which Git cannot
	// batch; a cache keyed by the immutable pair would bound it if the
	// number of deliveries on a target ever makes loading slow.
	for _, c := range live {
		// rev-list C ^S is empty exactly when C is an ancestor of S.
		out, err := repo.GitContext(ctx, root, "rev-list", "-n1", c.candidate, "^"+c.delivery.submitted)
		if err != nil {
			return err
		}
		if strings.TrimSpace(out) != "" {
			c.why = "delivery " + short(c.delivery.commit) + " names submitted tip " + short(c.delivery.submitted) + ", which does not contain the candidate"
		}
	}
	changed, err := diffTree(ctx, root, live, diffs.Bytes())
	if err != nil {
		return err
	}
	trees, err := mergeTrees(ctx, root, len(live), merges.Bytes())
	if err != nil {
		return err
	}
	sources, err := catBlobs(ctx, root, len(live), blobs.Bytes())
	if err != nil {
		return err
	}
	for i, c := range live {
		if c.why != "" {
			continue
		}
		allowed := map[string]bool{}
		for _, o := range records {
			if o.Type == "work" && sameCommit(o.Candidate, c.candidate) {
				allowed[path.Join(prefix, o.Path)] = true
			}
		}
		for _, f := range changed[i] {
			if !allowed[f] {
				c.why = "submitted tip " + short(c.delivery.submitted) + " changes " + f + " after the candidate, which its acceptance does not cover"
				break
			}
		}
		if c.why != "" {
			continue
		}
		if at, ds := project.ParseRecord(c.record.Path, sources[i]); sources[i] == nil || len(ds) != 0 || at.ID != c.record.ID || at.Status != "accepted" || !sameCommit(at.Approved, c.candidate) {
			c.why = "the record at submitted tip " + short(c.delivery.submitted) + " is not accepted for candidate " + short(c.candidate)
			continue
		}
		if trees[i] != c.delivery.tree {
			c.why = "delivery " + short(c.delivery.commit) + " is not what merging " + short(c.delivery.submitted) + " into its parent gives"
			continue
		}
		c.ok = true
	}
	return nil
}

// resolve maps each name to its full commit, or "" when the object is not
// here, in one git cat-file.
func resolve(ctx context.Context, root string, names []string) (map[string]string, error) {
	var in bytes.Buffer
	for _, n := range names {
		fmt.Fprintf(&in, "%s^{commit}\n", n)
	}
	got, err := lines(ctx, root, in.Bytes(), "cat-file", "--batch-check=%(objectname) %(objecttype)")
	if err != nil {
		return nil, err
	}
	full := map[string]string{}
	for i, n := range names {
		if i < len(got) {
			if f := strings.Fields(got[i]); len(f) == 2 && f[1] == "commit" {
				full[n] = f[0]
				continue
			}
		}
		full[n] = ""
	}
	return full, nil
}

// diffTree lists, per check in order, the paths changed from its candidate
// to its submitted tip, in one git diff-tree. Each group starts with the
// submitted tip's own ID, which is how the groups are told apart.
func diffTree(ctx context.Context, root string, checks []*check, in []byte) ([][]string, error) {
	out, err := run(ctx, root, in, "diff-tree", "--stdin", "-r", "--name-only", "-z")
	if err != nil {
		return nil, err
	}
	result := make([][]string, len(checks))
	i := -1
	for _, f := range strings.Split(strings.TrimSuffix(out, "\x00"), "\x00") {
		if next := i + 1; next < len(checks) && f == checks[next].delivery.submitted {
			i = next
			continue
		}
		if i >= 0 && f != "" {
			result[i] = append(result[i], f)
		}
	}
	return result, nil
}

// mergeTrees gives each merge's tree, or "" for a conflict, in one git
// merge-tree: per merge "status\0tree\0" then conflicted names and "\0".
func mergeTrees(ctx context.Context, root string, n int, in []byte) ([]string, error) {
	out, err := run(ctx, root, in, "merge-tree", "--write-tree", "--stdin", "--no-messages", "--name-only", "-z")
	if err != nil {
		return nil, err
	}
	f := strings.Split(out, "\x00")
	trees := make([]string, n)
	for i, j := 0, 0; i < n && j+1 < len(f); i++ {
		clean, tree := f[j], f[j+1]
		j += 2
		for j < len(f) && f[j] != "" {
			j++
		}
		j++
		if clean == "1" {
			trees[i] = tree
		}
	}
	return trees, nil
}

// catBlobs reads each blob, nil where it is missing, in one git cat-file.
func catBlobs(ctx context.Context, root string, n int, in []byte) ([][]byte, error) {
	out, err := run(ctx, root, in, "cat-file", "--batch")
	if err != nil {
		return nil, err
	}
	data := []byte(out)
	blobs := make([][]byte, n)
	for i := 0; i < n && len(data) != 0; i++ {
		header, rest, _ := bytes.Cut(data, []byte("\n"))
		f := strings.Fields(string(header))
		if len(f) != 3 || f[1] != "blob" {
			data = rest
			continue
		}
		var size int
		if _, err := fmt.Sscanf(f[2], "%d", &size); err != nil || size > len(rest) {
			return nil, errors.New("git cat-file --batch: unexpected output")
		}
		blobs[i] = rest[:size]
		data = rest[min(size+1, len(rest)):]
	}
	return blobs, nil
}

func run(ctx context.Context, root string, in []byte, args ...string) (string, error) {
	cmd := repo.Command(ctx, root, args...)
	cmd.Stdin = bytes.NewReader(in)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git %s: %v: %s", args[0], err, strings.TrimSpace(stderr.String()))
	}
	return string(out), nil
}

func lines(ctx context.Context, root string, in []byte, args ...string) ([]string, error) {
	out, err := run(ctx, root, in, args...)
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.TrimSuffix(out, "\n"), "\n"), nil
}

func sameCommit(a, b string) bool {
	return a != "" && b != "" && (strings.HasPrefix(a, b) || strings.HasPrefix(b, a))
}

func short(commit string) string { return commit[:min(len(commit), 7)] }
