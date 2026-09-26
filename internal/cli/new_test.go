package cli

import (
	"bytes"
	"context"
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mascah/grove/internal/repo"
)

// gitFixture commits the plain fixture so new can issue IDs.
func gitFixture(t *testing.T) string {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	root := projectFixture(t)
	for _, args := range [][]string{{"init", "-q", "-b", "main"}, {"add", "-A"}, {"-c", "user.name=t", "-c", "user.email=t@t", "-c", "commit.gpgsign=false", "-c", "maintenance.auto=false", "commit", "-q", "-m", "init"}} {
		if out, err := repo.Command(context.Background(), root, args...).CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	return root
}

func TestNewCreatesRecordAndReadCommandsLeaveNoState(t *testing.T) {
	t.Parallel()
	root := gitFixture(t)
	state := filepath.Join(root, ".git", "grove")
	for _, args := range [][]string{{"list"}, {"show", "G-260101-00001"}, {"check"}} {
		var out, errOut bytes.Buffer
		if code := Run(args, root, &out, &errOut); code != 0 {
			t.Fatal(errOut.String())
		}
	}
	if _, err := os.Stat(state); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("read commands must not create coordination state: %v", err)
	}
	second := newID(t, root, "work", "Second thing", "--slug", "second")
	if _, err := os.Stat(filepath.Join(root, "docs/records", second+"-second.md")); err != nil {
		t.Fatal(err)
	}
	var out, errOut bytes.Buffer
	if code := Run([]string{"--project", root, "new", "question", "Why?"}, t.TempDir(), &out, &errOut); code != 0 || !strings.HasSuffix(out.String(), "-why.md\n") {
		t.Fatalf("code=%d stdout=%q stderr=%q", code, out.String(), errOut.String())
	}
	if !strings.HasPrefix(errOut.String(), "Project: ") || strings.Count(errOut.String(), "\n") != 1 {
		t.Fatalf("new prints no notice: %q", errOut.String())
	}
	out.Reset()
	if code := Run([]string{"check"}, root, &out, &errOut); code != 0 || !strings.Contains(out.String(), "4 records") {
		t.Fatalf("created records must validate: %s %s", out.String(), errOut.String())
	}
	out.Reset()
	if code := Run([]string{"list"}, root, &out, &errOut); code != 0 || !strings.Contains(out.String(), second+"  work      proposed  Second thing") {
		t.Fatalf("list must show the created record: %s", out.String())
	}
	// G-260926-pgj43 acceptance 3: the write lock is the only coordination state.
	if entries, err := os.ReadDir(state); err != nil || len(entries) != 1 || entries[0].Name() != "write.lock" {
		t.Fatalf("state after new: %v, %v", entries, err)
	}
}

// G-260926-pgj43 acceptance 1: two clones that share nothing each issue date-form IDs,
// and the merge of one into the other still checks.
func TestNewInSeparateClonesMergesClean(t *testing.T) {
	t.Parallel()
	root := gitFixture(t)
	clone := filepath.Join(t.TempDir(), "clone")
	git := func(dir string, args ...string) {
		t.Helper()
		args = append([]string{"-c", "user.name=t", "-c", "user.email=t@t", "-c", "commit.gpgsign=false", "-c", "maintenance.auto=false"}, args...)
		if out, err := repo.Command(context.Background(), dir, args...).CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	git(root, "clone", "-q", root, clone)
	for _, dir := range []string{root, clone} {
		for _, title := range []string{"Shaped here", "Planned here", "Reviewed here"} {
			newID(t, dir, "work", title)
		}
		git(dir, "add", "-A")
		git(dir, "commit", "-q", "-m", "records")
	}
	git(clone, "pull", "-q", "--no-rebase", "--no-edit", root, "main")
	var out, errOut bytes.Buffer
	if code := Run([]string{"check"}, clone, &out, &errOut); code != 0 || out.String() != "OK: 8 records\n" {
		t.Fatalf("code=%d stdout=%q stderr=%s", code, out.String(), errOut.String())
	}
}

func TestNewUsageAndFailures(t *testing.T) {
	t.Parallel()
	for _, args := range [][]string{{"new"}, {"new", "work"}, {"new", "work", "T", "extra"}, {"new", "work", "T", "--slug"}, {"list", "--slug", "x"}, {"new", "work", "T", "--slug=a", "--slug=b"}} {
		var out, errOut bytes.Buffer
		if code := Run(args, t.TempDir(), &out, &errOut); code != 2 || out.Len() != 0 {
			t.Fatalf("%v: code=%d stderr=%s", args, code, errOut.String())
		}
	}
	root := projectFixture(t)
	before := hashes(t, root)
	for _, args := range [][]string{{"new", "work", "Outside Git"}, {"new", "release", "Bad type"}} {
		var out, errOut bytes.Buffer
		if code := Run(args, root, &out, &errOut); code != 1 || out.Len() != 0 || !strings.Contains(errOut.String(), "grove:") {
			t.Fatalf("%v: code=%d stdout=%q stderr=%q", args, code, out.String(), errOut.String())
		}
	}
	if !reflect.DeepEqual(before, hashes(t, root)) {
		t.Fatal("failed creation changed project files")
	}
}
