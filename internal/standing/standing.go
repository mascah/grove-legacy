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
	"os/exec"
	"path"
	"regexp"
	"slices"
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
	each, err := Each(ctx, root, target, records)
	if err != nil {
		return nil, err
	}
	out := map[string]*Standing{}
	for r, s := range each {
		out[r.ID] = s
	}
	return out, nil
}

// Each is Inspect by record rather than ID, for records drawn from several
// versions of the same work, such as a board's branches and checkouts.
func Each(ctx context.Context, root, target string, records []*project.Record) (map[*project.Record]*Standing, error) {
	out := map[*project.Record]*Standing{}
	var open []*project.Record // applicable acceptances
	for _, r := range records {
		if r.Type != "work" {
			continue
		}
		s := Of(r)
		out[r] = s
		if s.State == Unknown {
			s.Target = target
			open = append(open, r)
		}
	}
	if len(open) == 0 {
		return out, nil
	}
	unknown := func(why string) (map[*project.Record]*Standing, error) {
		for _, r := range open {
			out[r].Why = why
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
	// Git work here is a fixed number of processes, whatever the number of
	// deliveries, records or versions of them: one rev-parse, one log, one
	// cat-file and one rev-list, then one per batch in verify. Only a
	// delivery whose tree is not the plain merge costs more (bases).
	head, err := repo.GitContext(ctx, root, "rev-parse", "--is-shallow-repository", "--verify", "-q", "refs/heads/"+target+"^{commit}")
	f := strings.Fields(head)
	if err != nil || len(f) != 2 {
		return unknown("the target branch " + target + " cannot be read here")
	}
	complete, tip := f[0] == "false", f[1]
	claims, err := deliveries(ctx, root, tip)
	if err != nil {
		return nil, err
	}
	names := []string{}
	for _, r := range open {
		names = append(names, r.Candidate)
	}
	for _, d := range claims {
		names = append(names, d.submitted)
	}
	full, err := resolve(ctx, root, names)
	if err != nil {
		return nil, err
	}
	candidates := map[string]bool{}
	var heads []string
	for _, r := range open {
		if c := full[r.Candidate]; !candidates[c] {
			candidates[c] = true
			heads = append(heads, c)
		}
	}
	for _, d := range claims {
		if candidates[d.candidate] {
			heads = append(heads, full[d.submitted])
		}
	}
	graph, err := unmerged(ctx, root, tip, heads)
	if err != nil {
		return nil, err
	}
	var pending []*check
	mine := map[*project.Record][]*check{}
	shared := map[string]*check{} // one check per claim, record ID and path
	for _, r := range open {
		s := out[r]
		s.Tip = tip
		c := full[r.Candidate]
		_, off := graph[c]
		switch {
		case c == "":
			s.Why = "candidate " + short(r.Candidate) + " is not in this repository; fetch the branch or the refs/grove evidence that holds it"
		case !off:
			s.State, s.Delivered, s.Why = Done, c, ""
		default:
			for _, d := range claims {
				if d.candidate != c {
					continue
				}
				key := d.commit + "\x00" + r.ID + "\x00" + r.Path
				if shared[key] == nil {
					shared[key] = &check{record: r, candidate: c, delivery: d}
					pending = append(pending, shared[key])
				}
				mine[r] = append(mine[r], shared[key])
			}
			if len(mine[r]) == 0 {
				absent(s, complete)
			}
		}
	}
	if err := verify(ctx, root, prefix, records, pending, full, graph, &bases{claims: claims, verified: map[string]bool{}}); err != nil {
		return nil, err
	}
	// A record is done by any verified claim; otherwise a claim that is
	// missing its submitted tip, or that fails, leaves it unknown: the claim
	// is not a delivery, but nor is its failure proof of none, so nothing
	// may deliver it again until someone reconciles it.
	for _, r := range open {
		s := out[r]
		var missing, failed string
		for _, p := range mine[r] {
			switch {
			case s.State == Done:
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
			s.Why = "a delivery claim on " + s.Target + " names this candidate but does not verify: " + failed + "; reconcile it before delivering again"
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

// check is one claimed delivery of one record, shared by every version of
// it with the same ID and path.
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
// is present; the candidate is its ancestor; only records change
// between them; the record there is accepted for the candidate; and the
// claimed commit's first parent merged with the submitted tip gives exactly
// the claimed commit's tree. full resolves each submitted tip and graph is
// unmerged's. Each batch is one process holding only the checks still
// passing: a claim already failed is never merged or read, so nothing it
// names can fail every record's reading.
func verify(ctx context.Context, root, prefix string, records []*project.Record, checks []*check, full map[string]string, graph map[string][]string, earlier *bases) error {
	failed := func(c *check) bool { return c.why != "" }
	for _, c := range checks {
		switch d := c.delivery; {
		case full[d.submitted] == "":
			c.missing = true
			c.why = "delivery " + short(d.commit) + " names submitted tip " + short(d.submitted) + ", which is not in this repository; fetch refs/grove/* from where it was delivered"
		case !contains(graph, full[d.submitted], c.candidate):
			c.why = "delivery " + short(d.commit) + " names submitted tip " + short(d.submitted) + ", which does not contain the candidate"
		}
	}
	live := slices.DeleteFunc(slices.Clone(checks), failed)
	if len(live) == 0 {
		return nil
	}
	var diffs bytes.Buffer
	for _, c := range live {
		fmt.Fprintf(&diffs, "%s %s\n", c.delivery.submitted, c.candidate)
	}
	changed, err := diffTree(ctx, root, live, diffs.Bytes())
	if err != nil {
		return err
	}
	// Records may change after the candidate, the group's handoff and
	// acceptance among them; what the records say later cannot matter, so
	// any record read counts, never only those naming it now.
	allowed := map[string]bool{}
	for _, o := range records {
		allowed[path.Join(prefix, o.Path)] = true
	}
	for i, c := range live {
		for _, f := range changed[i] {
			if !allowed[f] {
				c.why = "submitted tip " + short(c.delivery.submitted) + " changes " + f + " after the candidate, which its acceptance does not cover"
				break
			}
		}
	}
	if live = slices.DeleteFunc(live, failed); len(live) == 0 {
		return nil
	}
	var blobs bytes.Buffer
	for _, c := range live {
		fmt.Fprintf(&blobs, "%s:%s\n", c.delivery.submitted, path.Join(prefix, c.record.Path))
	}
	sources, err := catBlobs(ctx, root, len(live), blobs.Bytes())
	if err != nil {
		return err
	}
	for i, c := range live {
		if at, ds := project.ParseRecord(c.record.Path, sources[i]); sources[i] == nil || len(ds) != 0 || at.ID != c.record.ID || at.Status != "accepted" || !sameCommit(at.Approved, c.candidate) {
			c.why = "the record at submitted tip " + short(c.delivery.submitted) + " is not accepted for candidate " + short(c.candidate)
		}
	}
	if live = slices.DeleteFunc(live, failed); len(live) == 0 {
		return nil
	}
	var merges bytes.Buffer
	for _, c := range live {
		fmt.Fprintf(&merges, "%s %s\n", c.delivery.parent, c.delivery.submitted)
	}
	trees, err := mergeTrees(ctx, root, len(live), merges.Bytes())
	if err != nil {
		return err
	}
	for i, c := range live {
		// A delivery's tree is judged once, however many records it holds.
		d := c.delivery
		ok, seen := earlier.verified[d.commit]
		if !seen {
			if ok = trees[i] == d.tree; !ok {
				if ok, err = earlier.rebased(ctx, root, d); err != nil {
					return err
				}
			}
			earlier.verified[d.commit] = ok
		}
		if !ok {
			c.why = "delivery " + short(d.commit) + " is not what merging " + short(d.submitted) + " into its parent gives"
			continue
		}
		c.ok = true
	}
	return nil
}

// unmerged maps each commit reachable from heads and not from tip to its
// parents, in one git rev-list: a head missing from it is on the target, and
// every path from a head to a commit in it stays in it.
// ponytail: the ^tip walk stops by commit date, so clock skew past Git's
// slop can list a commit the tip holds, which then reads as awaiting
// delivery or unknown, never done; merge-base --is-ancestor per candidate
// is exact if that is ever seen.
func unmerged(ctx context.Context, root, tip string, heads []string) (map[string][]string, error) {
	var in bytes.Buffer
	for _, h := range heads {
		if h != "" {
			fmt.Fprintln(&in, h)
		}
	}
	out, err := lines(ctx, root, in.Bytes(), "rev-list", "--parents", "--stdin", "^"+tip)
	if err != nil {
		return nil, err
	}
	graph := map[string][]string{}
	for _, l := range out {
		if f := strings.Fields(l); len(f) != 0 {
			graph[f[0]] = f[1:]
		}
	}
	return graph, nil
}

// contains reports whether commit, which is in graph, is head or an
// ancestor of it.
func contains(graph map[string][]string, head, commit string) bool {
	seen := map[string]bool{}
	for next := []string{head}; len(next) != 0; {
		c := next[len(next)-1]
		next = next[:len(next)-1]
		if c == commit {
			return true
		}
		if !seen[c] {
			seen[c] = true
			next = append(next, graph[c]...)
		}
	}
	return false
}

// Base is the submitted tip of the newest squash delivery reachable from
// tip that the next submission continues, or "": a branch kept after a
// squash delivery continues from what was delivered, so its next submission
// merges from there, not from where it first left the target. Only a
// retained submission submitted contains and tip lacks can be one, which
// one for-each-ref rules out on the ordinary path.
func Base(ctx context.Context, root, tip, submitted string) (string, error) {
	retained, err := repo.GitContext(ctx, root, "for-each-ref", "--format=%(objectname)", "--merged="+submitted, "--no-merged="+tip, "refs/grove/submitted/")
	if err != nil || strings.TrimSpace(retained) == "" {
		return "", err
	}
	claims, err := deliveries(ctx, root, tip)
	if err != nil {
		return "", err
	}
	return (&bases{claims: claims, verified: map[string]bool{}}).of(ctx, root, tip, submitted)
}

// bases finds and checks the earlier deliveries a later one merged from.
type bases struct {
	claims   []delivery
	verified map[string]bool // by delivery commit: its tree is its merge
}

// of is the newest claim that parent contains whose tree verifies and whose
// submitted tip descends from Git's own merge base of parent and submitted
// and is contained in submitted: it replaces that fork point, and nothing
// the branch took from the target since is older than it, so merging from
// it can neither undo the target's later changes nor skip the branch's.
// ponytail: one ancestry check per earlier claim among the branch's own
// commits, newest first; bounded by the deliveries that branch continues.
func (b *bases) of(ctx context.Context, root, parent, submitted string) (string, error) {
	// Unrelated histories have nothing to continue, and several best merge
	// bases (a criss-cross) are left to Git's own merge: the earlier
	// submission cannot be shown to replace them all.
	out, err := repo.GitContext(ctx, root, "merge-base", "--all", parent, submitted)
	fork := strings.Fields(out)
	if err != nil || len(fork) != 1 {
		return "", nil
	}
	own, err := lines(ctx, root, nil, "rev-list", "--ancestry-path", submitted, "^"+fork[0])
	if err != nil {
		return "", err
	}
	after := map[string]bool{}
	for _, c := range own {
		after[c] = true
	}
	for _, d := range b.claims {
		if d.submitted == submitted || !after[d.submitted] {
			continue
		}
		if on, err := ancestor(ctx, root, d.commit, parent); err != nil || !on {
			continue
		}
		ok, err := b.tree(ctx, root, d)
		if err != nil {
			return "", err
		}
		if ok {
			return d.submitted, nil
		}
	}
	return "", nil
}

// tree reports whether d's tree is its submitted tip merged onto its
// parent, from where they split or from the delivery it continued.
func (b *bases) tree(ctx context.Context, root string, d delivery) (bool, error) {
	if ok, seen := b.verified[d.commit]; seen {
		return ok, nil
	}
	b.verified[d.commit] = false // a cycle proves nothing
	trees, err := mergeTrees(ctx, root, 1, []byte(d.parent+" "+d.submitted+"\n"))
	if err != nil {
		return false, err
	}
	ok := trees[0] == d.tree
	if !ok {
		ok, err = b.rebased(ctx, root, d)
	}
	b.verified[d.commit] = ok
	return ok, err
}

// rebased reports whether d's tree is its submitted tip merged onto its
// parent from the earlier delivery it continues.
func (b *bases) rebased(ctx context.Context, root string, d delivery) (bool, error) {
	base, err := b.of(ctx, root, d.parent, d.submitted)
	if err != nil || base == "" {
		return false, err
	}
	trees, err := mergeTrees(ctx, root, 1, []byte(base+" -- "+d.parent+" "+d.submitted+"\n"))
	return err == nil && trees[0] == d.tree, err
}

func ancestor(ctx context.Context, root, commit, of string) (bool, error) {
	_, err := repo.GitContext(ctx, root, "merge-base", "--is-ancestor", commit, of)
	var exit *exec.ExitError
	if errors.As(err, &exit) && exit.ExitCode() == 1 {
		return false, nil
	}
	return err == nil, err
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
	out, err := run(ctx, root, in, "diff-tree", "--stdin", "--always", "-r", "--name-only", "-z")
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
// A claim can name histories unrelated to its parent, which Git merges from
// an empty base here rather than refusing the whole batch.
func mergeTrees(ctx context.Context, root string, n int, in []byte) ([]string, error) {
	out, err := run(ctx, root, in, "merge-tree", "--write-tree", "--stdin", "--allow-unrelated-histories", "--no-messages", "--name-only", "-z")
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
