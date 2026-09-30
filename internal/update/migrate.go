package update

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/repo"
)

// Migration is schema 3's conversion to 4 (G-260930-2qa4a), planned from a
// checkout's live files before anything is written.
type Migration struct {
	Root, Branch, Head string
	Changes            []Migrated // the work records that change, by ID
	Kept               []Migrated // done records kept as schema 3's claim
	Problems           []string   // what must be reconciled before writing
	Unchanged          int        // every other record
}

// Migrated is one work record's classification.
type Migrated struct {
	ID, Path, From, To, Why string
	by, context             string
	source                  []byte
}

var schemaLine = regexp.MustCompile(`(?m)^schema_version:[ \t]*3((?:[ \t]+#.*)?[ \t]*)$`)

// PlanMigration classifies root's records, which must be a Git checkout at
// schema 3. Proposed, active, abandoned and unapproved review work keep their
// meaning. An approval in review, or on done work whose candidate the
// configured target contains, becomes the acceptance it was, its authority
// read from the last verdict on that candidate. Every other done record keeps
// schema 3's claim, which is never presented as verified.
func PlanMigration(root string) (*Migration, error) {
	p, ds := project.LoadSchema3(root)
	if len(ds) != 0 {
		return nil, fmt.Errorf("the project is not a valid schema 3 project; fix it before migrating:\n%s", diagnostics(ds))
	}
	m := &Migration{Root: p.Root}
	if m.Root == "" {
		m.Root = root
	}
	branch, err := Branch(root)
	if err != nil {
		return nil, fmt.Errorf("migration needs a Git checkout: %v", err)
	}
	head, err := repo.Git(root, "rev-parse", "HEAD")
	if err != nil {
		return nil, fmt.Errorf("migration needs a commit to return to: %v", err)
	}
	m.Branch, m.Head = branch, strings.TrimSpace(head)
	if !schemaLine.Match(p.Config) {
		m.Problems = append(m.Problems, "grove.yaml: schema_version: 3 is not on a line of its own, which migration rewrites; write it so")
	}
	for _, r := range p.Records {
		if r.Type != "work" || r.Status != "review" && r.Status != "done" {
			m.Unchanged++
			continue
		}
		c := Migrated{ID: r.ID, Path: r.Path, From: r.Status, source: r.Source}
		switch {
		case r.Approved == "" && r.Status == "review":
			m.Unchanged++
			continue
		case r.Approved == "":
			c.To, c.Why = "done", "kept as schema 3's claim: no approval recorded"
			m.Kept = append(m.Kept, c)
			continue
		case r.Status == "done" && p.Target == "":
			c.To, c.Why = "done", "kept as schema 3's claim: no target configured to verify its delivery"
			m.Kept = append(m.Kept, c)
			continue
		case r.Status == "done":
			on, err := isAncestor(root, r.Candidate, "refs/heads/"+p.Target)
			if err != nil {
				m.Problems = append(m.Problems, fmt.Sprintf("%s: whether %s holds candidate %s could not be read (%v); fetch it, or reconcile the record", r.ID, p.Target, short(r.Candidate), err))
				continue
			}
			if !on {
				c.To, c.Why = "done", "kept as schema 3's claim: "+p.Target+" does not contain candidate "+short(r.Candidate)
				m.Kept = append(m.Kept, c)
				continue
			}
		}
		c.To, c.by, c.context = "accepted", authority(r), project.AcceptanceContext(r)
		c.Why = "approved " + short(r.Candidate) + " by " + c.by
		if r.Status == "done" {
			c.Why += "; " + p.Target + " contains it"
		}
		m.Changes = append(m.Changes, c)
	}
	return m, nil
}

// authority is who gave the last verdict on r's candidate: a sweep writes
// "delegated under policy grove.yaml REVISION", and anything else is the
// owner's.
func authority(r *project.Record) string {
	prefix := "Verdict on candidate " + short(r.Candidate) + ", "
	by := Owner
	for l := range strings.SplitSeq(string(r.Source), "\n") {
		if strings.HasPrefix(l, prefix) {
			by = Owner
			_, verdict, _ := strings.Cut(l, ": ")
			if rest, ok := strings.CutPrefix(verdict, "delegated under policy grove.yaml "); ok && len(rest) >= 71 {
				by = Policy(rest[:71])
			}
		}
	}
	return by
}

// Migrate writes a plan without problems under the write lock, in a
// checkout whose project has no uncommitted changes: it records
// refs/grove/schema-3/BRANCH at HEAD, the way back, then writes grove.yaml
// at schema 4 and each changed record, which keeps its ID, path and body
// with one paragraph appended, and commits them alone. A project that would
// not validate at schema 4 is restored and nothing is committed.
func Migrate(root string, now time.Time) (*Migration, string, error) {
	common, _, err := repo.CommonDir(root)
	if err != nil {
		return nil, "", err
	}
	unlock, err := repo.WriteLock(common)
	if err != nil {
		return nil, "", err
	}
	defer unlock()
	m, err := PlanMigration(root)
	if err != nil {
		return nil, "", err
	}
	if len(m.Problems) != 0 {
		return m, "", fmt.Errorf("nothing was migrated; reconcile first:\n%s", strings.Join(m.Problems, "\n"))
	}
	if dirty, err := repo.Git(root, "status", "--porcelain", "--", "."); err != nil {
		return m, "", err
	} else if dirty != "" {
		return m, "", errors.New("the project has uncommitted changes; commit or set them aside before migrating")
	}
	if m.Branch == "" {
		return m, "", errors.New("migration commits on a branch, and this checkout is on none")
	}
	if _, err := repo.Git(root, "update-ref", "refs/grove/schema-3/"+m.Branch, m.Head); err != nil {
		return m, "", err
	}
	day := now.UTC().Format("2006-01-02")
	stamp := strconv.Quote(now.UTC().Truncate(time.Second).Format("2006-01-02T15:04:05Z"))
	writes := map[string][]byte{}
	original := map[string][]byte{}
	config, err := os.ReadFile(filepath.Join(root, "grove.yaml"))
	if err != nil {
		return m, "", err
	}
	original["grove.yaml"] = config
	writes["grove.yaml"] = schemaLine.ReplaceAll(config, []byte("schema_version: 4${1}"))
	for _, c := range m.Changes {
		source, err := Edit(c.source, []change{set("status", "accepted"), set("approved_by", value(c.by)), set("approved_context", value(c.context)), set("updated", stamp)})
		if err == nil {
			source, err = appended(source, fmt.Sprintf("Migrated to schema 4, %s: status %s with approval of its candidate became status accepted by %s.", day, c.From, c.by))
		}
		if err != nil {
			return m, "", fmt.Errorf("%s: %w; nothing was written", c.Path, err)
		}
		path := filepath.FromSlash(c.Path)
		original[path], writes[path] = c.source, source
	}
	restore := func() {
		for path, source := range original {
			_ = os.WriteFile(filepath.Join(root, path), source, 0o644)
		}
	}
	for path, source := range writes {
		if err := os.WriteFile(filepath.Join(root, path), source, 0o644); err != nil {
			restore()
			return m, "", fmt.Errorf("writing %s: %w; the files were restored", path, err)
		}
	}
	if _, ds := project.Load(root, root); len(ds) != 0 {
		restore()
		return m, "", fmt.Errorf("the migrated project would not validate, so the files were restored:\n%s", diagnostics(ds))
	}
	paths := []string{}
	for path := range writes {
		paths = append(paths, ":(literal)"+path)
	}
	if _, err := repo.Git(root, append([]string{"add", "--"}, paths...)...); err != nil {
		return m, "", fmt.Errorf("%w; the files are migrated but nothing was committed", err)
	}
	message := fmt.Sprintf("chore: migrate records to schema 4\n\n%d work records became accepted and %d done records keep schema 3's claim; the way back is refs/grove/schema-3/%s.", len(m.Changes), len(m.Kept), m.Branch)
	if _, err := repo.Git(root, append([]string{"commit", "-q", "-m", message, "--"}, paths...)...); err != nil {
		return m, "", fmt.Errorf("%w; the files are migrated and staged but nothing was committed", err)
	}
	commit, err := repo.Git(root, "rev-parse", "HEAD")
	return m, strings.TrimSpace(commit), err
}

// value quotes a string field as update does for anything but a plain word.
func value(s string) string {
	if word.MatchString(s) {
		return s
	}
	return strconv.Quote(s)
}

func isAncestor(root, commit, ref string) (bool, error) {
	_, err := repo.Git(root, "merge-base", "--is-ancestor", commit, ref)
	var exit *exec.ExitError
	if errors.As(err, &exit) && exit.ExitCode() == 1 {
		return false, nil
	}
	return err == nil, err
}
