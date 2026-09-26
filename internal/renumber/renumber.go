// Package renumber is the one-time rename of legacy G-NNN records to the
// date form (G-260926-yvjy6). It leaves the binary once adopters have run it.
package renumber

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/mascah/grove/internal/attempt"
	"github.com/mascah/grove/internal/create"
	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/repo"
	"github.com/mascah/grove/internal/update"
)

// Renamed is one line of the map, as convert prints one conversion.
type Renamed struct{ From, FromPath, ID, Path string }

var (
	legacy = regexp.MustCompile(`^` + project.NeutralPrefix + `-[0-9]{3}$`)
	token  = regexp.MustCompile(project.NeutralPrefix + `-[0-9]+`)
)

// Run renames every legacy record in root's project to a date-form ID drawn
// for its creation date and rewrites every reference under the record root
// and in the attempt store. Every refusal comes before the first write.
func Run(root string) ([]Renamed, error) {
	p, ds := project.Load(root, root)
	if len(ds) != 0 {
		return nil, fmt.Errorf("the project is not valid; fix it before renumbering:\n%s", diagnostics(ds))
	}
	if !slices.ContainsFunc(p.Records, func(r *project.Record) bool { return legacy.MatchString(r.ID) }) {
		return nil, fmt.Errorf("nothing renumbered: no record has a legacy ID")
	}
	common, showPrefix, err := repo.CommonDir(p.Root)
	if err != nil {
		return nil, err
	}
	unlock, err := repo.WriteLock(common)
	if err != nil {
		return nil, err
	}
	defer unlock()
	refused := func(err error) ([]Renamed, error) { return nil, fmt.Errorf("nothing renumbered: %w", err) }
	if p, ds = project.Load(root, root); len(ds) != 0 {
		return refused(fmt.Errorf("the project no longer validates:\n%s", diagnostics(ds)))
	}
	if err := otherBranches(p); err != nil {
		return refused(err)
	}
	attempts, err := attempt.List(p.Root, "")
	if err != nil {
		return refused(err)
	}
	live := map[string]bool{}
	for _, v := range attempts {
		if v.Status != attempt.Running && v.Status != attempt.Orphaned {
			continue
		}
		for _, m := range v.Launch.Members() {
			if legacy.MatchString(m.ID) {
				return refused(fmt.Errorf("attempt %s of %s is %s; stop it first", filepath.Base(v.Dir), m.ID, v.Status))
			}
		}
		live[v.Dir] = true
	}

	// Plan every rename before writing anything.
	var renames []Renamed
	created := map[string]string{} // record path -> created to add
	issued := *p
	issued.Records = slices.Clone(p.Records)
	for _, r := range p.Records {
		if !legacy.MatchString(r.ID) {
			continue
		}
		when, err := creation(p.Root, r)
		if err != nil {
			return refused(err)
		}
		if r.Created == nil {
			created[r.Path] = strconv.Quote(when.Format("2006-01-02T15:04:05Z"))
		}
		// ponytail: Issue scans every ref per record; one-time, so not worth batching.
		id, err := create.Issue(&issued, showPrefix, when)
		if err != nil {
			return refused(err)
		}
		issued.Records = append(issued.Records, &project.Record{ID: id})
		to := path.Join(path.Dir(r.Path), id+"-"+create.Slug(r.Title)+".md")
		if _, err := os.Lstat(filepath.Join(p.Root, filepath.FromSlash(to))); !os.IsNotExist(err) {
			return refused(fmt.Errorf("%s already exists", to))
		}
		renames = append(renames, Renamed{From: r.ID, FromPath: r.Path, ID: id, Path: to})
	}
	slices.SortFunc(renames, func(a, b Renamed) int { return strings.Compare(a.From, b.From) })

	names, ids := map[string]string{}, map[string]string{}
	for _, m := range renames {
		names[path.Base(m.FromPath)] = path.Base(m.Path)
		ids[m.From] = m.ID
	}
	incomplete := func(err error) ([]Renamed, error) {
		return renames, fmt.Errorf("%w (the renumber is incomplete; the map is printed, inspect with Git)", err)
	}
	for _, m := range renames {
		from := filepath.Join(p.Root, filepath.FromSlash(m.FromPath))
		if value, ok := created[m.FromPath]; ok {
			source, err := os.ReadFile(from)
			if err == nil {
				source, err = update.Set(source, "created", value)
			}
			if err == nil {
				err = os.WriteFile(from, source, 0o644)
			}
			if err != nil {
				return incomplete(fmt.Errorf("%s: %w", m.FromPath, err))
			}
		}
		if err := os.Rename(from, filepath.Join(p.Root, filepath.FromSlash(m.Path))); err != nil {
			return incomplete(err)
		}
	}
	if err := rewriteTree(filepath.Join(p.Root, p.RecordDir), names, ids, nil); err != nil {
		return incomplete(err)
	}
	if _, ds := project.Load(root, root); len(ds) != 0 {
		return incomplete(fmt.Errorf("the project no longer validates:\n%s", diagnostics(ds)))
	}
	dir := filepath.Join(common, "grove", "attempts")
	if err := rewriteTree(dir, names, ids, live); err != nil {
		return incomplete(err)
	}
	for _, v := range attempts {
		name := filepath.Base(v.Dir)
		id, rest, _ := strings.Cut(name, ".")
		if to, ok := ids[id]; ok {
			if err := os.Rename(v.Dir, filepath.Join(dir, to+"."+rest)); err != nil {
				return incomplete(err)
			}
		}
	}
	return renames, nil
}

// otherBranches refuses while a local branch other than HEAD's and the
// target holds records: the board and versions read every local branch and
// would show both forms.
func otherBranches(p *project.Project) error {
	current, err := update.Branch(p.Root)
	if err != nil {
		return err
	}
	out, err := repo.Git(p.Root, "for-each-ref", "--format=%(refname:short)", "refs/heads")
	if err != nil {
		return err
	}
	var holding []string
	for b := range strings.FieldsSeq(out) {
		if b == current || b == p.Target {
			continue
		}
		files, err := repo.Git(p.Root, "ls-tree", "-r", "--name-only", "refs/heads/"+b, "--", filepath.ToSlash(p.RecordDir))
		if err != nil {
			return err
		}
		if strings.TrimSpace(files) != "" {
			holding = append(holding, b)
		}
	}
	if len(holding) != 0 {
		return fmt.Errorf("local branches other than this one and the target hold records: %s; merge or delete them first", strings.Join(holding, ", "))
	}
	return nil
}

// creation is r's created time, else the UTC author date of the first
// commit of its formerly path, else of its own path.
func creation(root string, r *project.Record) (time.Time, error) {
	if r.Created != nil {
		return r.Created.UTC(), nil
	}
	first := func(args ...string) (time.Time, bool, error) {
		out, err := repo.Git(root, append([]string{"log", "--format=%at"}, args...)...)
		if err != nil {
			return time.Time{}, false, err
		}
		lines := strings.Fields(out)
		if len(lines) == 0 {
			return time.Time{}, false, nil
		}
		unix, err := strconv.ParseInt(lines[len(lines)-1], 10, 64)
		return time.Unix(unix, 0).UTC(), err == nil, err
	}
	if r.Formerly != "" {
		if t, ok, err := first("--", r.Formerly); ok || err != nil {
			return t, err
		}
	}
	if t, ok, err := first("--follow", "--", r.Path); ok || err != nil {
		return t, err
	}
	return time.Time{}, fmt.Errorf("%s has no created and no commit to date it from", r.Path)
}

// rewriteTree rewrites every regular file under dir, skipping the
// directories in skip.
func rewriteTree(dir string, names, ids map[string]string, skip map[string]bool) error {
	err := filepath.WalkDir(dir, func(p string, e fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if e.IsDir() && skip[p] {
			return filepath.SkipDir
		}
		if !e.Type().IsRegular() {
			return nil
		}
		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		if changed := Rewrite(data, names, ids); !bytes.Equal(changed, data) {
			return os.WriteFile(p, changed, 0o644) // an existing file keeps its mode
		}
		return nil
	})
	if os.IsNotExist(err) && skip != nil {
		return nil // no attempt store
	}
	return err
}

// Rewrite replaces, in one pass, each old basename and then each old ID as
// a whole token: an ID preceded or followed by a letter or digit is part of
// another word, while a hyphen bounds one, so worktree-G-030 and G-030-G-031
// are rewritten. names and ids map old to new; a basename starts with its ID.
func Rewrite(data []byte, names, ids map[string]string) []byte {
	base := map[string]string{} // old ID -> its old basename
	for name := range names {
		id, _, _ := strings.Cut(strings.TrimPrefix(name, project.NeutralPrefix+"-"), "-")
		base[project.NeutralPrefix+"-"+id] = name
	}
	var out []byte
	last := 0
	for _, m := range token.FindAllIndex(data, -1) {
		id := string(data[m[0]:m[1]])
		to, ok := ids[id]
		if !ok {
			continue
		}
		end := m[1]
		if name := base[id]; name != "" && bytes.HasPrefix(data[m[0]:], []byte(name)) {
			to, end = names[name], m[0]+len(name)
		} else if m[0] > 0 && alnum(data[m[0]-1]) || m[1] < len(data) && alnum(data[m[1]]) {
			continue
		}
		out = append(append(out, data[last:m[0]]...), to...)
		last = end
	}
	if out == nil {
		return data
	}
	return append(out, data[last:]...)
}

func alnum(c byte) bool {
	return c >= '0' && c <= '9' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

func diagnostics(ds []project.Diagnostic) string {
	lines := make([]string, len(ds))
	for i, d := range ds {
		lines[i] = d.String()
	}
	return strings.Join(lines, "\n")
}
