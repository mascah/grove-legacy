// Package standing derives each work record's standing (G-260930-2qa4a,
// G-260930-gj9d7): what its record says, and whether an acceptance it holds
// reached the target. Done is never read from a status. An acceptance
// reaches the target only through a delivery, so an applicable acceptance is
// done when the target tip's own copy of the record, at the same path, is
// accepted for the same candidate, and otherwise awaits delivery. A reading
// consults nothing else, no history, trailers or refs/grove, so it costs the
// same however old the repository is and however many branches or
// deliveries it has. integrate proves each delivery once, as it makes it
// (Prove), and Audit proves them again when a person asks. Nothing here
// writes, fetches or pushes.
package standing

import (
	"bytes"
	"cmp"
	"context"
	"errors"
	"fmt"
	"path"
	"regexp"
	"slices"
	"strconv"
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
	Accepted  = "accepted" // accepted, not yet on the target
	Done      = "done"
	Unknown   = "unknown" // accepted; the target could not be read
	Abandoned = "abandoned"
)

// Standing is one work record's recorded and derived facts.
type Standing struct {
	ID       string `json:"id"`
	Recorded string `json:"recorded"` // the record's own status, as its file says
	State    string `json:"state"`
	Legacy   bool   `json:"legacy,omitempty"` // done is schema 3's claim, not verified
	// Target and Tip are the branch and commit the acceptance was judged at.
	Target string `json:"target,omitempty"`
	Tip    string `json:"tip,omitempty"`
	Why    string `json:"why,omitempty"`
}

// Text is the standing as one phrase: "done: delivered to main".
func (s *Standing) Text() string {
	switch {
	case s.Legacy:
		return "done (schema 3 claim)"
	case s.State == Done:
		return "done: delivered to " + s.Target
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
// and its review and acceptance, through cleanup and garbage collection,
// for Audit, and for Retired to tell a delivered branch. It is outside
// refs/heads, so it is evidence, never a work branch.
func Ref(submitted string) string { return "refs/grove/submitted/" + submitted }

// Delivery is a squash delivery a branch holds: its commit on the target and
// the work that commit names.
type Delivery struct {
	Commit string
	Work   []string
}

// Holds says which delivery retires branch, the opening of every refusal.
func (d *Delivery) Holds(branch string) string {
	return fmt.Sprintf("%s holds the delivery %s of %s", branch, short(d.Commit), strings.Join(d.Work, ", "))
}

// Retired is the delivery tip holds, or nil: a retained submission tip
// contains and the target does not, whose squash commit is on the target's
// first parents. A workspace holding one is retired from execution
// (G-260930-tcc9w): kept for inspection, and no command acts on it again,
// whichever work the delivery names and whichever would come next. Every
// command that would act asks here, and none narrows the answer:
//
//   - run refuses it named with --branch, and moves the default name on to
//     -2, -3, …;
//   - approve and feedback refuse in its checkout, and so does resolve,
//     whose feedback that is, before anything is written;
//   - integrate never delivers from it and refuses work in review or
//     accepted only there; its rerun for work that reads done only cleans up.
//
// A hand edit there is preserved and never delivered. A submission retained
// whose target never advanced, or one the target contains as an ordinary
// merge, retires nothing. It walks the target's first parents, so no
// reading asks.
// ponytail: finds submissions through refs/grove alone, so a clone that
// fetched the branch without them cannot tell; fetch refs/grove/* with it.
func Retired(ctx context.Context, root, target, tip string) (*Delivery, error) {
	out, err := repo.Command(ctx, root, "for-each-ref", "--merged="+tip, "--format=%(objectname)", "refs/grove/submitted/").Output()
	if err != nil {
		return nil, fmt.Errorf("the retained submissions in %s could not be read: %v", tip, err)
	}
	for s := range strings.FieldsSeq(string(out)) {
		if repo.Command(ctx, root, "merge-base", "--is-ancestor", s, "refs/heads/"+target).Run() == nil {
			continue
		}
		out, err := repo.Command(ctx, root, "log", "-1", "--first-parent", "--format=%H%n%(trailers:key=Grove-Work,valueonly)", "--fixed-strings", "--grep=Grove-Submitted: "+s, "refs/heads/"+target).Output()
		if err != nil {
			return nil, fmt.Errorf("%s's history could not be read: %v", target, err)
		}
		if lines := strings.Fields(string(out)); len(lines) != 0 {
			return &Delivery{Commit: lines[0], Work: lines[1:]}, nil
		}
	}
	return nil, nil
}

// Of is a record's standing from its file alone, before the target is read:
// an accepted record whose acceptance applies stays Unknown until Judge.
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

// Judge is the standing of every work record in records against copies,
// the target tip's own records by project-relative path, read at commit
// tip; nil copies means the target could not be read. It starts no process,
// for a caller that already holds the target's records, as the board does.
func Judge(target, tip string, copies map[string]*project.Record, records []*project.Record) map[*project.Record]*Standing {
	out := map[*project.Record]*Standing{}
	for _, r := range records {
		if r.Type != "work" {
			continue
		}
		s := Of(r)
		out[r] = s
		if s.State != Unknown {
			continue
		}
		s.Target, s.Tip = target, tip
		switch t := copies[r.Path]; {
		case target == "":
			s.Why = "no target is configured in grove.yaml"
		case copies == nil:
			s.Why = "the target branch " + target + " cannot be read here"
		case t != nil && t.ID == r.ID && t.Status == "accepted" && sameCommit(t.Candidate, r.Candidate):
			s.State, s.Why = Done, ""
		default:
			s.State, s.Why = Accepted, ""
		}
	}
	return out
}

// Each is Judge with the target's copies read from root's repository, in one
// git cat-file whatever the number of records or versions of them, and none
// when no acceptance applies. Records may be drawn from several versions of
// the same work.
func Each(ctx context.Context, root, target string, records []*project.Record) map[*project.Record]*Standing {
	var paths []string
	for _, r := range records {
		if r.Type == "work" && Of(r).State == Unknown && !slices.Contains(paths, r.Path) {
			paths = append(paths, r.Path)
		}
	}
	if target == "" || len(paths) == 0 {
		return Judge(target, "", nil, records)
	}
	tip, copies := read(ctx, root, target, paths)
	return Judge(target, tip, copies, records)
}

// Inspect is Each by ID, for the records of one version of the project.
func Inspect(ctx context.Context, root, target string, records []*project.Record) map[string]*Standing {
	out := map[string]*Standing{}
	for r, s := range Each(ctx, root, target, records) {
		out[r.ID] = s
	}
	return out
}

// read is the target tip's commit and its own record at each of paths, in
// one git cat-file; a path the target lacks, or holds unparsable, has no
// copy. Copies are nil when the target branch cannot be read.
func read(ctx context.Context, root, target string, paths []string) (string, map[string]*project.Record) {
	var in bytes.Buffer
	fmt.Fprintf(&in, "refs/heads/%s^{commit}\n", target)
	for _, p := range paths {
		// ponytail: a path holding a newline misreads; record paths never do.
		fmt.Fprintf(&in, "refs/heads/%s:./%s\n", target, p)
	}
	out, err := run(ctx, root, in.Bytes(), "cat-file", "--batch")
	if err != nil {
		return "", nil
	}
	objects, err := batch(out, len(paths)+1)
	if err != nil || objects[0].kind != "commit" {
		return "", nil
	}
	copies := map[string]*project.Record{}
	for i, p := range paths {
		if o := objects[i+1]; o.kind == "blob" {
			if r, ds := project.ParseRecord(p, o.data); len(ds) == 0 {
				copies[p] = r
			}
		}
	}
	return objects[0].id, copies
}

// Proof is what Audit or Prove found of one record the target holds
// accepted: proved by the target containing the candidate or by a squash
// delivery Git verifies, not proved, or not auditable in this clone.
type Proof struct {
	ID          string `json:"id"`
	Proved      bool   `json:"proved"`
	Unauditable bool   `json:"unauditable,omitempty"` // the evidence is not in this clone
	// Delivered is the target's commit holding the result: the candidate
	// for an ordinary merge, or the squash commit, whose submitted tip is
	// Submitted.
	Delivered string `json:"delivered,omitempty"`
	Submitted string `json:"submitted,omitempty"`
	Why       string `json:"why,omitempty"`
}

// Text is the proof as one phrase.
func (p *Proof) Text() string {
	switch {
	case p.Proved && p.Submitted != "":
		return "proved: squashed as " + short(p.Delivered) + " from submitted tip " + short(p.Submitted)
	case p.Proved:
		return "proved: the target contains candidate " + short(p.Delivered)
	case p.Unauditable:
		return "cannot be audited here: " + p.Why
	}
	return "not proved: " + p.Why
}

// Audit proves the delivery of every work record the target tip holds under
// recordDir, the project-relative record root, whose acceptance applies: the
// target contains its candidate, or a squash delivery on the target names it
// and verifies. records are the checkout's, which with the target's may
// change after a candidate. It walks the target's history and reads
// refs/grove, so it runs only when a person asks (grove check --deliveries),
// never on a reading.
func Audit(ctx context.Context, root, target, recordDir string, records []*project.Record) ([]*Proof, error) {
	if target == "" {
		return nil, errors.New("no target is configured in grove.yaml")
	}
	listed, err := run(ctx, root, nil, "ls-tree", "-r", "-z", "--name-only", "refs/heads/"+target, "--", "./"+cmp.Or(recordDir, "."))
	if err != nil {
		return nil, fmt.Errorf("the target branch %s cannot be read here", target)
	}
	var paths []string
	for _, p := range strings.Split(listed, "\x00") {
		if strings.HasSuffix(p, ".md") {
			paths = append(paths, p)
		}
	}
	_, copies := read(ctx, root, target, paths)
	if copies == nil {
		return nil, fmt.Errorf("the target branch %s cannot be read here", target)
	}
	var open []*project.Record
	for _, p := range paths {
		if t := copies[p]; t != nil {
			records = append(records, t)
			if t.Type == "work" && Of(t).State == Unknown {
				open = append(open, t)
			}
		}
	}
	if len(open) == 0 {
		return nil, nil
	}
	head, err := repo.GitContext(ctx, root, "rev-parse", "--is-shallow-repository", "--verify", "-q", "refs/heads/"+target+"^{commit}")
	f := strings.Fields(head)
	if err != nil || len(f) != 2 {
		return nil, fmt.Errorf("the target branch %s cannot be read here", target)
	}
	claims, err := deliveries(ctx, root, f[1])
	if err != nil {
		return nil, err
	}
	return prove(ctx, root, f[1], f[0] == "false", open, records, claims)
}

// Prove proves the one squash commit delivered, which integrate just made,
// for members, the target's copies of what it delivered: scoped to that
// commit and the branch it came from, with no walk of the target's history.
// records are every record read, any of which may change after a candidate.
func Prove(ctx context.Context, root, delivered string, members, records []*project.Record) ([]*Proof, error) {
	claims, err := deliveries(ctx, root, "--no-walk", delivered)
	if err != nil {
		return nil, err
	}
	return prove(ctx, root, delivered, true, members, records, claims)
}

// prove decides each of open, the target's accepted copies, at tip from
// claims, the delivery claims found there; complete is false in a shallow
// clone, which cannot show a delivery absent.
func prove(ctx context.Context, root, tip string, complete bool, open, records []*project.Record, claims []delivery) ([]*Proof, error) {
	_, _, prefix, err := repo.IdentifyContext(ctx, root)
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
	proofs := make([]*Proof, len(open))
	for i, r := range open {
		p := &Proof{ID: r.ID}
		proofs[i] = p
		c := full[r.Candidate]
		_, off := graph[c]
		switch {
		case c == "":
			p.Unauditable, p.Why = true, "candidate "+short(r.Candidate)+" is not in this repository; fetch the branch or the refs/grove evidence that holds it"
		case !off:
			p.Proved, p.Delivered = true, c
		default:
			for _, d := range claims {
				if d.candidate == c {
					k := &check{record: r, candidate: c, delivery: d}
					pending = append(pending, k)
					mine[r] = append(mine[r], k)
				}
			}
			switch {
			case len(mine[r]) != 0:
			case complete:
				p.Why = "no delivery on the target names candidate " + short(r.Candidate) + ", yet the target's record is accepted for it"
			default:
				p.Unauditable, p.Why = true, "this clone's history is shallow, so the target may hold a delivery it cannot see; fetch the full history"
			}
		}
	}
	if err := verify(ctx, root, prefix, records, pending, full, graph); err != nil {
		return nil, err
	}
	// Proved by any verified claim; otherwise a claim missing its submitted
	// tip cannot be audited here, and one that fails is not proof.
	for i, r := range open {
		p := proofs[i]
		var missing, failed string
		for _, k := range mine[r] {
			switch {
			case p.Proved:
			case k.ok:
				p.Proved, p.Delivered, p.Submitted = true, k.delivery.commit, k.delivery.submitted
			case k.missing:
				missing = cmp.Or(missing, k.why)
			default:
				failed = cmp.Or(failed, k.why)
			}
		}
		switch {
		case p.Proved || missing == "" && failed == "":
		case missing != "":
			p.Unauditable, p.Why = true, missing
		default:
			p.Why = "a delivery claim names this candidate but does not verify: " + failed
		}
	}
	return proofs, nil
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
	missing   bool // the submitted tip is absent: not auditable, not failed
	why       string
}

var hexCommit = regexp.MustCompile("^[0-9a-f]{40}([0-9a-f]{24})?$")

// deliveries lists the commits git log revs lists whose trailers name a
// candidate and a submitted tip. The trailers only locate: verify decides.
func deliveries(ctx context.Context, root string, revs ...string) ([]delivery, error) {
	args := append([]string{"log", "-z", "--fixed-strings", "--grep=Grove-Submitted: ", "--format=%H%x1f%P%x1f%T%x1f%(trailers:only,unfold)"}, revs...)
	out, err := repo.GitContext(ctx, root, args...)
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
// names can fail every record's audit.
func verify(ctx context.Context, root, prefix string, records []*project.Record, checks []*check, full map[string]string, graph map[string][]string) error {
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
	out, err := run(ctx, root, blobs.Bytes(), "cat-file", "--batch")
	if err != nil {
		return err
	}
	sources, err := batch(out, len(live))
	if err != nil {
		return err
	}
	for i, c := range live {
		if at, ds := project.ParseRecord(c.record.Path, sources[i].data); sources[i].kind != "blob" || len(ds) != 0 || at.ID != c.record.ID || at.Status != "accepted" || !sameCommit(at.Approved, c.candidate) {
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
		if trees[i] != c.delivery.tree {
			c.why = "delivery " + short(c.delivery.commit) + " is not what merging " + short(c.delivery.submitted) + " into its parent gives"
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
// slop can list a commit the tip holds, which an audit then reports not
// proved, never proved; merge-base --is-ancestor per candidate is exact if
// that is ever seen.
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

// object is one git cat-file --batch answer; kind is "" where it is missing.
type object struct {
	id, kind string
	data     []byte
}

// batch parses n answers of git cat-file --batch.
func batch(out string, n int) ([]object, error) {
	data := []byte(out)
	objects := make([]object, n)
	for i := 0; i < n && len(data) != 0; i++ {
		header, rest, _ := bytes.Cut(data, []byte("\n"))
		data = rest
		f := strings.Fields(string(header))
		if len(f) != 3 || !hexCommit.MatchString(f[0]) {
			continue // "NAME missing" or "NAME ambiguous"
		}
		size, err := strconv.Atoi(f[2])
		if err != nil || size >= len(rest) {
			return nil, errors.New("git cat-file --batch: unexpected output")
		}
		objects[i] = object{f[0], f[1], rest[:size]}
		data = rest[size+1:]
	}
	return objects, nil
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
