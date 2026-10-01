package attempt

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/repo"
)

// admission is the bounded set of records a new workspace, started from the
// target, takes from the launching checkout's HEAD (G-260930-tcc9w): each
// selected record the target lacks, every record an admitted one's
// relationships name that the target lacks, which validation needs, and each
// plan or question the target lacks whose work or blocks names a selected
// member. Nothing else of HEAD: no code, configuration, harness file or
// unrelated record. A selected record the target holds with other bytes is a
// changed input, refused; an admitted path the target holds is never taken.
// Other work admitted must be proposed or abandoned: work further along
// names code the workspace lacks, and its acceptance would reach the target
// with this delivery. An admitted file with uncommitted changes is refused,
// since what is taken is HEAD's.
func admission(p *project.Project, ids []string, head, target string) ([]Member, error) {
	at, err := blobs(p.Root, target, p.RecordDir)
	if err != nil {
		return nil, err
	}
	ah, err := blobs(p.Root, head, p.RecordDir)
	if err != nil {
		return nil, err
	}
	byID := map[string]*project.Record{}
	for _, r := range p.Records {
		byID[r.ID] = r
	}
	queue, seen := slices.Clone(ids), map[string]bool{}
	for _, id := range ids {
		for _, r := range p.Records {
			if slices.Contains(r.Work, id) || slices.Contains(r.Blocks, id) {
				queue = append(queue, r.ID)
			}
		}
	}
	var admitted []Member
	for len(queue) != 0 {
		id := queue[0]
		queue = queue[1:]
		r := byID[id]
		if seen[id] || r == nil {
			continue
		}
		seen[id] = true
		h, t := ah[r.Path], at[r.Path]
		switch {
		case h == "":
			return nil, fmt.Errorf("%s is not committed at HEAD %s; commit it so the new workspace can take it", r.Path, short(head))
		case t == h:
			continue
		case t != "" && slices.Contains(ids, id):
			return nil, fmt.Errorf("%s differs between HEAD %s and the target %s: the workspace starts from the target, so deliver or reconcile that record there first", r.Path, short(head), short(target))
		case t != "":
			continue // the target's copy stands
		}
		if dirty, err := repo.Git(p.Root, "status", "--porcelain", "--", r.Path); err != nil {
			return nil, err
		} else if strings.TrimSpace(dirty) != "" {
			return nil, fmt.Errorf("%s has uncommitted changes in this checkout; commit them so the new workspace takes what you read", r.Path)
		}
		if r.Type == "work" && !slices.Contains(ids, id) && r.Status != "proposed" && r.Status != "abandoned" {
			return nil, fmt.Errorf("the selection names %s, which the target lacks and which is %s at HEAD %s: its code is not on the target, so deliver it first or reconcile the reference", id, r.Status, short(head))
		}
		admitted = append(admitted, Member{ID: r.ID, Path: r.Path, Revision: project.Revision(r.Source), Status: r.Status})
		for _, targets := range [][]string{r.DependsOn, r.Members, r.Blocks, r.RelatesTo, r.Work} {
			queue = append(queue, targets...)
		}
	}
	return admitted, nil
}

// blobs maps each file under dir, project-relative, to its blob at commit.
func blobs(root, commit, dir string) (map[string]string, error) {
	out, err := repo.Git(root, "ls-tree", "-r", "-z", commit, "--", dir)
	if err != nil {
		return nil, fmt.Errorf("the records at %s could not be listed: %v", short(commit), err)
	}
	m := map[string]string{}
	for entry := range strings.SplitSeq(strings.TrimSuffix(out, "\x00"), "\x00") {
		info, path, ok := strings.Cut(entry, "\t")
		if f := strings.Fields(info); ok && len(f) == 3 {
			m[filepath.ToSlash(path)] = f[2]
		}
	}
	return m, nil
}

// admit copies the admitted records from source into the new worktree's
// project, validates the project there and commits them with their
// provenance. A project that does not validate removes the worktree and the
// branch just made for it, since nothing else is on them.
func admit(root, worktree, prefix, branch, source string, admitted []Member, report func(string)) error {
	dir := filepath.Join(worktree, prefix)
	paths := make([]string, len(admitted))
	ids := make([]string, len(admitted))
	for i, m := range admitted {
		paths[i], ids[i] = filepath.FromSlash(m.Path), m.ID
	}
	undo := func(err error) error {
		repo.Git(root, "worktree", "remove", "--force", "--", worktree)
		repo.Git(root, "branch", "-D", branch)
		os.RemoveAll(worktree)
		return err
	}
	if _, err := repo.Git(dir, append([]string{"checkout", source, "--"}, paths...)...); err != nil {
		return undo(fmt.Errorf("the records could not be taken from %s: %v; the new workspace was removed", short(source), err))
	}
	if _, ds := project.Load(dir, dir); len(ds) != 0 {
		var lines []string
		for _, d := range ds {
			lines = append(lines, d.String())
		}
		return undo(fmt.Errorf("the records taken from %s do not validate on the target: %s; the new workspace was removed", short(source), strings.Join(lines, "; ")))
	}
	message := fmt.Sprintf("chore: admit %s from %s\n\n%s\n\nGrove-Admitted-From: %s\n", strings.Join(ids, ", "), short(source), strings.Join(paths, "\n"), source)
	if _, err := repo.Git(dir, "commit", "-q", "-m", message); err != nil {
		return undo(fmt.Errorf("the admitted records could not be committed: %v; the new workspace was removed", err))
	}
	report(fmt.Sprintf("admitted: %s from %s", strings.Join(ids, ", "), short(source)))
	return nil
}
