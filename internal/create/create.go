// Package create issues coordination-free IDs and writes new records.
// The reader stays Git-free; Git access goes through internal/repo.
package create

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/repo"
)

// bodies are the skeletons new writes; the rest of a type is project.Types.
var bodies = map[string]string{
	"work":     "## Outcome\n\n## Constraints\n\n## Acceptance\n\n## Next\n",
	"question": "## Question\n\n## Next\n",
	"decision": "## Decision\n\n## Alternatives\n\n## Reconsideration\n",
	"term":     "## Meaning\n\n## Relationships\n",
	"plan":     "## Design\n\n## Steps\n",
	"review":   "## Examined\n\n## Findings\n\n## Disposition\n",
	"page":     "",
}

var (
	slugPattern = regexp.MustCompile(`^[a-z0-9-]+$`)
	idLine      = regexp.MustCompile(`^id:\s*["']?(` + project.IDForm + `)["']?\s*$`)
)

// New issues an ID for kind, creates the record without overwriting
// anything, and reloads the project so an unreadable result fails loudly.
// It returns the created file's path relative to the project root.
func New(p *project.Project, kindName, title, slug string, now time.Time) (string, error) {
	k := project.Type(kindName)
	if k == nil {
		return "", fmt.Errorf("record type must be work, question, decision, term, plan, review, or page")
	}
	title = strings.TrimSpace(title)
	if title == "" {
		return "", fmt.Errorf("a nonempty title is required")
	}
	// A term's title is its identity, so new can collide where other types
	// cannot. Refuse before locking; a term that appears after this load is
	// refused again under the write lock, before anything is written.
	if err := definedTerm(p, kindName, title); err != nil {
		return "", err
	}
	if slug == "" {
		slug = Slug(title)
	} else if !ValidSlug(slug) {
		return "", fmt.Errorf("slug must contain only lowercase ASCII letters, digits, and hyphens")
	}
	common, showPrefix, err := repo.CommonDir(p.Root)
	if err != nil {
		return "", err
	}
	// The ID is drawn under the write lock that serializes publication with
	// update, so a second new in any worktree sees the first one's file.
	unlock, err := repo.WriteLock(common)
	if err != nil {
		return "", fmt.Errorf("nothing created: %w", err)
	}
	defer unlock()
	current, ds := project.Load(p.Root, p.Root)
	if len(ds) != 0 {
		return "", fmt.Errorf("nothing created: the project no longer validates:\n%s", diagnostics(ds))
	}
	if current.RecordDir != p.RecordDir {
		return "", fmt.Errorf("nothing created: the record root changed from %s to %s since it was read", p.RecordDir, current.RecordDir)
	}
	if err := definedTerm(current, kindName, title); err != nil {
		return "", fmt.Errorf("nothing created: %w", err)
	}
	if !bytes.Equal(current.Config, p.Config) {
		return "", fmt.Errorf("nothing created: grove.yaml changed since it was read; inspect it and retry")
	}
	id, err := Issue(current, showPrefix, now)
	if err != nil {
		return "", fmt.Errorf("nothing created: %w", err)
	}
	// Every type creates flat in the record root.
	relative := path.Join(filepath.ToSlash(p.RecordDir), id+"-"+slug+".md")
	stamp := now.UTC().Format("2006-01-02T15:04:05Z")
	status := ""
	if len(k.Statuses) != 0 {
		status = "status: " + k.Statuses[0] + "\n"
	}
	// %q emits Go escapes, a subset of YAML double-quoted escapes.
	content := fmt.Sprintf("---\nid: %q\ntype: %s\ntitle: %q\n%screated: %q\nupdated: %q\n---\n\n%s",
		id, kindName, title, status, stamp, stamp, bodies[kindName])
	full := filepath.Join(p.Root, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		return "", fmt.Errorf("nothing created: %w", err)
	}
	f, err := os.OpenFile(full, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return "", fmt.Errorf("nothing created: %w", err)
	}
	if _, err := f.WriteString(content); err != nil {
		f.Close()
		return "", fmt.Errorf("%s: %w", relative, err)
	}
	if err := f.Close(); err != nil {
		return "", fmt.Errorf("%s: %w", relative, err)
	}
	if _, ds := project.Load(p.Root, p.Root); len(ds) != 0 {
		return "", fmt.Errorf("created %s but the project no longer validates:\n%s", relative, diagnostics(ds))
	}
	return relative, nil
}

func definedTerm(p *project.Project, kindName, title string) error {
	for _, r := range p.Records {
		if kindName == "term" && r.Type == "term" && project.TermKey(r.Title) == project.TermKey(title) {
			return fmt.Errorf("the term %s is already defined by %s in %s", title, r.ID, r.Path)
		}
	}
	return nil
}

func diagnostics(ds []project.Diagnostic) string {
	lines := make([]string, len(ds))
	for i, d := range ds {
		lines[i] = d.String()
	}
	return strings.Join(lines, "\n")
}

// ValidSlug reports whether a caller-supplied slug is acceptable in a filename.
func ValidSlug(slug string) bool { return slugPattern.MatchString(slug) }

// Slug derives a short filename fragment from a title per the record model.
func Slug(title string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(title) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		case b.Len() > 0 && !strings.HasSuffix(b.String(), "-"):
			b.WriteByte('-')
		}
	}
	s := b.String()
	if len(s) > 24 {
		s = s[:24]
	}
	s = strings.TrimRight(s, "-")
	if s == "" {
		return "record"
	}
	return s
}

// tries bounds the draws for one ID. With 32^5 tails a day, needing more
// means a broken random source, not bad luck.
const tries = 8

const crockford = "0123456789abcdefghjkmnpqrstvwxyz"

// tail draws the random part of an ID from crypto/rand; tests replace it.
var tail = func() string {
	b := make([]byte, 5)
	rand.Read(b) // never fails: crypto/rand crashes the program instead
	for i := range b {
		b[i] = crockford[b[i]&31]
	}
	return string(b)
}

// Issue draws an ID for now, the UTC date and a random tail, that no record
// in p, on any local ref or in any worktree already holds (G-260926-2da4n). It reads
// no shared state and writes nothing; callers hold the write lock.
func Issue(p *project.Project, showPrefix string, now time.Time) (string, error) {
	day := project.NeutralPrefix + "-" + now.UTC().Format("060102") + "-"
	used, err := usedIDs(p.Root, p.RecordDir, showPrefix, day)
	if err != nil {
		return "", err
	}
	for _, r := range p.Records {
		used[r.ID] = true
	}
	for range tries {
		if id := day + tail(); !used[id] {
			return id, nil
		}
	}
	return "", fmt.Errorf("every one of %d random IDs for %s was already in use", tries, strings.TrimSuffix(day, "-"))
}

// usedIDs collects the IDs starting with day in committed records on every
// local ref and live records in every worktree. Lines that merely look like
// IDs only count as used.
// ponytail: one git grep over all refs per new; narrow the refs if
// repositories with many refs make this slow.
func usedIDs(root, recordDir, showPrefix, day string) (map[string]bool, error) {
	used := map[string]bool{}
	note := func(line string) {
		if m := idLine.FindStringSubmatch(line); m != nil {
			used[m[1]] = true
		}
	}
	refs, err := repo.Git(root, "for-each-ref", "--format=%(objectname)", "refs/heads", "refs/remotes", "refs/tags")
	if err != nil {
		return nil, err
	}
	if trees := strings.Fields(refs); len(trees) != 0 {
		// git grep only pre-filters; note validates every line. day is
		// letters, digits and hyphens, so it is its own pattern.
		args := append([]string{"grep", "-h", "-I", "-e", "^id:.*" + day}, trees...)
		cmd := repo.Command(context.Background(), root, append(args, "--", recordDir)...)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		out, err := cmd.Output()
		var exit *exec.ExitError
		if err != nil && !(errors.As(err, &exit) && exit.ExitCode() == 1) { // 1 means no matches
			return nil, fmt.Errorf("git grep: %s", strings.TrimSpace(stderr.String()))
		}
		for _, line := range strings.Split(string(out), "\n") {
			note(line)
		}
	}
	worktrees, err := repo.Worktrees(root)
	if err != nil {
		return nil, err
	}
	for _, w := range worktrees {
		tree := filepath.Join(w.Path, filepath.FromSlash(showPrefix), recordDir)
		err := filepath.WalkDir(tree, func(p string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() || filepath.Ext(p) != ".md" {
				return nil
			}
			data, err := os.ReadFile(p)
			if err != nil {
				return err
			}
			for _, line := range strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n") {
				note(line)
			}
			return nil
		})
		// A worktree without the record folder is normal; anything else would
		// silently miss an ID in use and risk issuing it again.
		if err != nil && !(errors.Is(err, fs.ErrNotExist) && strings.HasPrefix(err.Error(), "lstat "+tree)) {
			return nil, fmt.Errorf("cannot scan worktree records: %w", err)
		}
	}
	return used, nil
}
