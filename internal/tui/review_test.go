package tui

import (
	"context"
	"errors"
	standings "github.com/mascah/grove/internal/standing"
	"slices"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/charmbracelet/x/ansi"

	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/sweep"
	"github.com/mascah/grove/internal/versions"
)

// reviewFixture: W-001 is in review on branch feature, held by the checkout
// feat, with a review record; main still has it proposed, superseded by
// feature's change. The target is main, whose checkout is ".".
func reviewFixture(fx fixture, approved bool) *fake {
	var vs []versions.Version
	for _, s := range []*versions.Source{fx.cMain, fx.main} {
		v := version(s, "W-001", "Inspect records", "proposed")
		v.Older = "branch feature changed it since"
		vs = append(vs, v)
	}
	for _, s := range []*versions.Source{fx.cFeat, fx.feat} {
		w := version(s, "W-001", "Inspect records", "review")
		w.Record.Candidate = "abcdef1"
		w.Record.Source = []byte("---\nid: W-001\n---\n\n## Outcome\n\nA board.\n" + strings.Repeat("\nfiller\n", 30) + "\n## Evidence\n\nIt works.\n\n## Next\n\nJudge it.\n")
		if approved {
			accept(w.Record, "owner")
		}
		review := version(s, "W-006", "W-001 review", "current")
		review.Record.Type, review.Record.Work, review.Record.Examined = "review", []string{"W-001"}, "abcdef1"
		vs = append(vs, w, review)
	}
	res := result(fx.main, fx.sources(), vs...)
	res.Target = "main"
	return &fake{res: res, actions: true,
		history: func(context.Context, string, string) ([]versions.Commit, error) { return nil, nil },
		changes: func(target, candidate, tip, path string) (*versions.Changes, error) {
			merge := &versions.Merge{Target: strings.Repeat("a", 40), Commit: candidate, Outcome: "conflict", Conflicts: []string{"internal/x.go"}}
			return &versions.Changes{Base: "base000", Merge: merge, Files: []versions.Change{{Path: "internal/x.go", Added: 12, Removed: 3}, {Path: "grove/work/W-001.md", Added: 5, Removed: 1}, {Path: "bin.dat", Added: -1, Removed: -1}}}, nil
		},
		diff: func(from, to, path string) (string, error) {
			return "diff --git a/" + path + " b/" + path + "\n@@ -1,2 +1,3 @@\n-old\n+new \x1b]0;evil\a\n+\tmore\n", nil
		},
	}
}

// openReview opens W-001's detail from the Review column and delivers the
// history and changes reads.
func openReview(t *testing.T, f *fake, w, h int) *Model {
	t.Helper()
	m := open(t, f, w, h)
	press(m, "right", "right")
	if m.cardID != "W-001" || m.col != 2 {
		t.Fatalf("W-001 should be the Review column's card: col %d card %s", m.col, m.cardID)
	}
	deliverAll(m, press(m, "enter"))
	return m
}

// A prediction read against a target commit other than the board's reading
// of the target says so (G-260925-h8rj5): one of the two is stale, and r re-reads.
func TestReviewMergePredictionNamesAMovedTarget(t *testing.T) {
	t.Parallel()
	f := reviewFixture(newFixture(), true)
	f.changes = func(target, candidate, tip, path string) (*versions.Changes, error) {
		return &versions.Changes{Base: "base000", Merge: &versions.Merge{Target: strings.Repeat("b", 40), Commit: candidate, Outcome: "clean", Conflicts: []string{}}}, nil
	}
	s := plain(openReview(t, f, 200, 36))
	if want := "accepted · only the record changed since it · merges cleanly into main at bbbbbbb, which moved since the branch left it; the board read main at aaaaaaa (r re-reads)"; !strings.Contains(s, want) {
		t.Fatalf("review detail lacks %q:\n%s", want, s)
	}
}

// A merge Git could not predict says so and keeps the changed files.
func TestReviewMergeUnpredicted(t *testing.T) {
	t.Parallel()
	f := reviewFixture(newFixture(), true)
	f.changes = func(target, candidate, tip, path string) (*versions.Changes, error) {
		return &versions.Changes{Base: "base000", Unpredicted: "git merge-tree: too old", Files: []versions.Change{{Path: "internal/x.go", Added: 1}}}, nil
	}
	s := plain(openReview(t, f, 200, 36))
	for _, want := range []string{"only the record changed since it · the merge into main could not be predicted: git merge-tree: too old", "internal/x.go  +1 −0"} {
		if !strings.Contains(s, want) {
			t.Fatalf("review detail lacks %q:\n%s", want, s)
		}
	}
}

// The detail of a candidate in review leads with its standing and where the
// actions run, starts at the Evidence, lists the changed files, and shows a
// file's diff escaped; a narrow terminal keeps every row its width.
func TestReviewDetailShowsStandingChangesAndDiffs(t *testing.T) {
	t.Parallel()
	fx := newFixture()
	f := reviewFixture(fx, false)
	m := openReview(t, f, 120, 36)
	s := plain(m)
	for _, want := range []string{
		"W-001 · review", "candidate abcdef1 · not on main",
		"Review: candidate abcdef1 · not yet accepted · only the record changed since it · conflicts with main at aaaaaaa in", "┃ internal/x.go ",
		"a approve and f feedback run on branch feature in /repo/feat · i integrate runs into main in /repo",
		"## Evidence", "It works.",
		"review     W-006  W-001 review  current", "examined abcdef1 = candidate",
		"Changes against main from base000", "internal/x.go  +12 −3", "grove/work/W-001.md  +5 −1", "bin.dat  binary",
		"a approve  f feedback  i integrate",
	} {
		if !strings.Contains(s, want) {
			t.Fatalf("review detail lacks %q:\n%s", want, s)
		}
	}
	if strings.Contains(s, "## Outcome") {
		t.Fatalf("the content should start at the Evidence:\n%s", s)
	}
	if strings.Join(f.reads, ";") != "changes main abcdef1 a grove/work/W-001.md" {
		t.Fatalf("reads: %v", f.reads)
	}
	// Tab reaches the changes after the linked records; Enter shows the diff.
	press(m, "tab", "down")
	s = plain(m)
	if !strings.Contains(s, "> internal/x.go  +12 −3") {
		t.Fatalf("Tab and ↓ should reach the changed file:\n%s", s)
	}
	deliverAll(m, press(m, "enter"))
	s = plain(m)
	for _, want := range []string{"Diff of internal/x.go (Esc returns to the content)", "-old", `+new \x1b]0;evil\a`, "+    more"} {
		if !strings.Contains(s, want) {
			t.Fatalf("diff view lacks %q:\n%s", want, s)
		}
	}
	if strings.Contains(m.render(), "\x1b]0;") {
		t.Fatal("the diff's control sequence reached the screen")
	}
	if strings.Join(f.reads, ";") != "changes main abcdef1 a grove/work/W-001.md;diff base000 abcdef1 internal/x.go" {
		t.Fatalf("reads: %v", f.reads)
	}
	press(m, "esc")
	if s = plain(m); m.diff != "" || !strings.Contains(s, "▶ Content") && !strings.Contains(s, "  Content") || strings.Contains(s, "Diff of") {
		t.Fatalf("Esc returns to the content:\n%s", s)
	}
	press(m, "esc")
	if m.screen != boardScreen {
		t.Fatal("Esc then returns to the board")
	}
	// Narrow: every row is the terminal's width with the review block.
	for _, size := range [][2]int{{99, 30}, {80, 24}, {40, 12}} {
		m := openReview(t, reviewFixture(fx, true), size[0], size[1])
		for _, row := range strings.Split(m.render(), "\n") {
			if got := ansi.StringWidth(row); got != size[0] {
				t.Fatalf("%v: row is %d cells: %q", size, got, ansi.Strip(row))
			}
		}
		if s := plain(m); size[1] >= 16 && !strings.Contains(s, "Review: candidate abcdef1 · accepted") {
			t.Fatalf("%v: review block missing:\n%s", size, s)
		}
		// One row a file, whatever describes it (G-260928-r1hkh).
		side := strings.Split(ansi.Strip(strings.Join(m.sidebar(m.shown(m.group()), size[0], 100), "\n")), "\n")
		from := slices.IndexFunc(side, func(r string) bool { return strings.HasPrefix(r, "Changes against") })
		to := slices.IndexFunc(side, func(r string) bool { return strings.HasPrefix(r, "Timeline on") })
		if from < 0 || to-from-1 != 3 {
			t.Fatalf("%v: three files in %d rows:\n%s", size, to-from-1, strings.Join(side, "\n"))
		}
	}
}

// Esc from the detail cancels a changes or diff read as it cancels history.
func TestReviewReadsStopWithTheDetail(t *testing.T) {
	t.Parallel()
	f := reviewFixture(newFixture(), false)
	m := open(t, f, 120, 36)
	press(m, "right", "right", "enter")
	deliverAll(m, press(m, "esc")) // the history read is cancelled: the detail's own
	if m.pending != "" && m.pending != "inspect" {
		t.Fatalf("pending %q", m.pending)
	}
	m = openReview(t, f, 120, 36)
	press(m, "tab", "down")
	if cmd := press(m, "enter"); cmd == nil || m.pending != "diff" {
		t.Fatalf("Enter on a change starts the diff read: pending %q", m.pending)
	}
	press(m, "esc")
	if m.pending != "" || m.cancel != nil {
		t.Fatalf("Esc must cancel the diff read: pending %q", m.pending)
	}
}

// a and f open a prompt; Enter runs the action in the branch's checkout,
// shows its facts on the result screen, and re-reads the board; keys wait
// while it runs; Esc cancels a prompt; an empty verdict is refused.
func TestReviewApproveAndFeedbackFromTheBoard(t *testing.T) {
	t.Parallel()
	fx := newFixture()
	f := reviewFixture(fx, false)
	m := openReview(t, f, 120, 36)
	press(m, "a")
	if s := plain(m); m.prompt == nil || !strings.Contains(s, "Verdict on W-001: ▏") || !strings.Contains(s, "Enter approves W-001 on branch feature with it; Esc cancels") {
		t.Fatalf("a should open the verdict prompt:\n%s", s)
	}
	press(m, "enter")
	if s := plain(m); m.prompt == nil || !strings.Contains(s, "type the verdict first, or Esc") {
		t.Fatalf("an empty verdict is refused:\n%s", s)
	}
	typeText(m, "Good ✓")
	press(m, "backspace")
	press(m, "esc")
	if s := plain(m); m.prompt != nil || !strings.Contains(s, "cancelled; nothing was written") || len(f.acts) != 0 {
		t.Fatalf("Esc cancels the prompt:\n%s", s)
	}
	press(m, "a")
	typeText(m, "Ship it")
	cmd := press(m, "enter")
	if cmd == nil || m.pending != "act" || m.prompt != nil || !strings.Contains(plain(m), "Approving…") {
		t.Fatalf("Enter should start the action: pending %q\n%s", m.pending, plain(m))
	}
	press(m, "q", "esc", "a")
	if m.done || !strings.Contains(plain(m), "Approving… Keys wait for it.") {
		t.Fatalf("keys wait while an action runs:\n%s", plain(m))
	}
	next := deliver(m, cmd)
	if m.screen != resultScreen || next == nil || m.pending != "inspect" {
		t.Fatalf("the outcome shows and the board is re-read: screen %d pending %q", m.screen, m.pending)
	}
	s := plain(m)
	for _, want := range []string{"Approved W-001", "approved: W-001 in /repo/feat", "Re-reading the board…"} {
		if !strings.Contains(s, want) {
			t.Fatalf("result lacks %q:\n%s", want, s)
		}
	}
	press(m, "esc") // waits for the re-read, so the record shown next is current
	if m.screen != resultScreen || m.pending != "inspect" || !strings.Contains(plain(m), "Esc waits for the re-read") {
		t.Fatal("Esc must wait for the re-read, and the hints say so")
	}
	deliverAll(m, next)
	if s := plain(m); m.screen != resultScreen || !strings.Contains(s, "The board has been re-read. Esc returns to the record.") {
		t.Fatalf("after the re-read:\n%s", s)
	}
	if f.inspects != 2 || strings.Join(f.acts, ";") != "approve /repo/feat W-001 Ship it" {
		t.Fatalf("inspects %d acts %v", f.inspects, f.acts)
	}
	press(m, "esc")
	if m.screen != detailScreen || m.openID() != "W-001" {
		t.Fatalf("Esc returns to the record: screen %d", m.screen)
	}
	// Feedback, with the continuation in its facts.
	press(m, "f")
	if s := plain(m); !strings.Contains(s, "Feedback on W-001: ▏") || !strings.Contains(s, "Enter records it, returning it to active on branch feature") {
		t.Fatalf("f should open the feedback prompt:\n%s", s)
	}
	typeText(m, "Needs the empty case")
	deliverAll(m, press(m, "enter"))
	if s := plain(m); m.screen != resultScreen || !strings.Contains(s, "Feedback recorded on W-001") || !strings.Contains(s, "next: /grove-work W-001") {
		t.Fatalf("feedback outcome:\n%s", s)
	}
	if f.acts[1] != "feedback /repo/feat W-001 Needs the empty case" {
		t.Fatalf("acts %v", f.acts)
	}
	// A failed action shows why, and the facts it returned.
	press(m, "esc", "a")
	f.fail = errors.New("this checkout does not hold candidate abcdef1")
	typeText(m, "x")
	deliverAll(m, press(m, "enter"))
	if s := plain(m); !strings.Contains(s, "NOT DONE: this checkout does not hold candidate abcdef1") {
		t.Fatalf("a failure is shown:\n%s", s)
	}
}

// i needs an approval and the target's checkout, confirms the merge and then
// the cleanup, and runs in the target's checkout.
func TestReviewIntegrateFromTheBoard(t *testing.T) {
	t.Parallel()
	fx := newFixture()
	m := openReview(t, reviewFixture(fx, false), 120, 36)
	press(m, "i")
	if s := plain(m); m.prompt != nil || !strings.Contains(s, "accept candidate abcdef1 first (a)") {
		t.Fatalf("i before approval:\n%s", s)
	}
	f := reviewFixture(fx, true)
	m = openReview(t, f, 120, 36)
	press(m, "a")
	if s := plain(m); m.prompt != nil || !strings.Contains(s, "candidate abcdef1 is already accepted; i integrates it") {
		t.Fatalf("a after approval:\n%s", s)
	}
	press(m, "i")
	if s := plain(m); m.prompt == nil || !strings.Contains(s, "Squash branch feature onto main, delivering W-001? y/n   (runs in /repo)") {
		t.Fatalf("i should confirm the merge:\n%s", s)
	}
	press(m, "n")
	if s := plain(m); m.prompt != nil || !strings.Contains(s, "cancelled; nothing was merged") || len(f.acts) != 0 {
		t.Fatalf("n cancels:\n%s", s)
	}
	press(m, "i", "y")
	if s := plain(m); m.prompt == nil || !strings.Contains(s, "Remove branch feature and its worktree? y/n, n keeps them   (/repo/feat)") {
		t.Fatalf("y should ask about cleanup:\n%s", s)
	}
	deliverAll(m, press(m, "n"))
	if s := plain(m); m.screen != resultScreen || !strings.Contains(s, "Integration of W-001") || !strings.Contains(s, "merge: fast-forward") {
		t.Fatalf("integration outcome:\n%s", s)
	}
	press(m, "esc", "i", "y")
	deliverAll(m, press(m, "y"))
	if strings.Join(f.acts, ";") != "integrate /repo W-001 keep=true;integrate /repo W-001 keep=false" {
		t.Fatalf("acts %v", f.acts)
	}
	// Other keys during a y/n prompt do nothing.
	press(m, "esc", "i", "x", "q")
	if m.prompt == nil || m.done {
		t.Fatal("a y/n prompt ignores other keys")
	}
	press(m, "esc")
	if m.prompt != nil {
		t.Fatal("Esc cancels a y/n prompt")
	}
}

// The actions need a clean checkout of the branch, and the target's
// checkout: the review block says what is missing, and the key says so too.
func TestReviewActionsNeedTheRightCheckout(t *testing.T) {
	t.Parallel()
	fx := newFixture()
	f := reviewFixture(fx, true)
	for i := range f.res.Groups[0].Versions { // W-001's checkout copy is modified
		if v := &f.res.Groups[0].Versions[i]; v.Source == fx.feat {
			v.Change = "modified"
		}
	}
	m := openReview(t, f, 120, 36)
	if s := plain(m); !strings.Contains(s, "a approve and f feedback: the record has uncommitted changes in checkout feat; commit or discard them before judging") {
		t.Fatalf("review block:\n%s", s)
	}
	press(m, "f")
	if s := plain(m); m.prompt != nil || !strings.Contains(s, "uncommitted changes in checkout feat") {
		t.Fatalf("f on a modified record:\n%s", s)
	}
	// No checkout of the branch, and no checkout of the target.
	f = reviewFixture(fx, true)
	f.res.Sources = []*versions.Source{fx.cFeat, fx.cMain}
	f.res.Groups[0].Versions = f.res.Groups[0].Versions[:0:0]
	for _, s := range []*versions.Source{fx.cMain, fx.cFeat} {
		v := version(s, "W-001", "Inspect records", map[bool]string{true: "proposed", false: "review"}[s == fx.cMain])
		if s == fx.cMain {
			v.Older = "branch feature changed it since"
		} else {
			v.Record.Candidate = "abcdef1"
			accept(v.Record, "owner")
		}
		f.res.Groups[0].Versions = append(f.res.Groups[0].Versions, v)
	}
	f.res.Groups = f.res.Groups[:1]
	m = openReview(t, f, 120, 36)
	if s := plain(m); !strings.Contains(s, "a approve and f feedback: no checkout is on branch feature; git worktree add one, or run grove approve there") || !strings.Contains(s, "integrate: no checkout is on the target main") {
		t.Fatalf("review block without checkouts:\n%s", s)
	}
	press(m, "a")
	if s := plain(m); m.prompt != nil || !strings.Contains(s, "no checkout is on branch feature") {
		t.Fatalf("a without a checkout:\n%s", s)
	}
	press(m, "i")
	if s := plain(m); m.prompt != nil || !strings.Contains(s, "no checkout is on the target main") {
		t.Fatalf("i without the target's checkout:\n%s", s)
	}
	if len(f.acts) != 0 {
		t.Fatalf("nothing should have run: %v", f.acts)
	}
	// Not in review: the keys explain.
	m = open(t, linkedFixture(newFixture()), 120, 36)
	press(m, "right")
	deliverAll(m, press(m, "enter"))
	press(m, "a")
	if s := plain(m); m.prompt != nil || !strings.Contains(s, "W-001 is not in review: nothing to approve") || strings.Contains(s, "Review: candidate") {
		t.Fatalf("a on active work:\n%s", s)
	}
}

// Each changed file is one row that counts the other records linking it or
// naming it in a code span, and its diff names them; nothing more is read
// for them (G-260928-r1hkh). A path under the project's prefix, which Git
// gives with its slash, is matched as a project path, one outside it only
// by a code span, and a rename by either side.
func TestReviewListsRecordsDescribingEachFile(t *testing.T) {
	t.Parallel()
	for _, prefix := range []string{"", "proj/"} {
		fx := newFixture()
		f := reviewFixture(fx, false)
		for _, s := range []*versions.Source{fx.cMain, fx.main, fx.cFeat, fx.feat} {
			f.res.Groups = append(f.res.Groups,
				versions.Group{ID: "Q-002", Versions: []versions.Version{withBody(version(s, "Q-002", "Where?", "open"), "In `x.go`, and [bin](../../bin.dat).\n")}},
				versions.Group{ID: "D-002", Versions: []versions.Version{withBody(version(s, "D-002", "Keep x", "accepted"), "[x](../../internal/x.go), [old](../../old/y.go), `Cargo.toml`\n")}})
		}
		for i := range f.res.Groups {
			for j := range f.res.Groups[i].Versions {
				if v := &f.res.Groups[i].Versions[j]; v.Record.ID == "W-001" {
					*v = withBody(*v, "Changes [x](../../internal/x.go).\n")
					v.Record.Candidate = "abcdef1"
				}
			}
		}
		f.res.Prefix = prefix
		changes := f.changes
		f.changes = func(target, candidate, tip, path string) (*versions.Changes, error) {
			c, err := changes(target, candidate, tip, path)
			for i := range c.Files {
				c.Files[i].Path = prefix + c.Files[i].Path
			}
			c.Files = append(c.Files, versions.Change{Path: prefix + "old/y.go → " + prefix + "new/y.go", Added: 1}, versions.Change{Path: "Cargo.toml", Added: 1})
			return c, err
		}
		m := openReview(t, f, 160, 50)
		s := plain(m)
		want := []string{"internal/x.go  +12 −3  described by 2 ", "grove/work/W-001.md  +5 −1  ", "bin.dat  binary  described by 1 ",
			"new/y.go  +1 −0  described by 1 ", "Cargo.toml  +1 −0  described by 1 "}
		at := 0
		for _, w := range want {
			i := strings.Index(s[at:], w)
			if i < 0 {
				t.Fatalf("prefix %q: the review lacks %q in order:\n%s", prefix, w, s)
			}
			at += i + len(w)
		}
		if strings.Contains(s, "W-001 link") || strings.Contains(s, "Q-002 code span") || strings.Join(f.reads, ";") != "changes main abcdef1 a grove/work/W-001.md" {
			t.Fatalf("prefix %q: the open record is not listed, the names wait for the diff, and nothing more is read: %v\n%s", prefix, f.reads, s)
		}
		press(m, "tab", "tab") // past the linked review to the first file
		for _, c := range []struct {
			moves int
			head  string
		}{
			{0, "Described by Q-002 code span, D-002 link."}, {1, "No record names this file."}, {1, "Described by Q-002 link."},
			{1, "Described by D-002 link."}, {1, "Described by D-002 code span."},
		} {
			for range c.moves {
				press(m, "down")
			}
			deliverAll(m, press(m, "enter"))
			if s := plain(m); !strings.Contains(s, c.head) {
				t.Fatalf("prefix %q: the diff of %s lacks %q:\n%s", prefix, m.diff, c.head, s)
			}
		}
	}
}

// Records sharing a candidate on the branch (G-260925-wc2pz) are named where the
// action covers them, and their files are not changes after the candidate.
func TestReviewNamesTheGroupSharingACandidate(t *testing.T) {
	t.Parallel()
	fx := newFixture()
	f := reviewFixture(fx, true)
	for _, s := range []*versions.Source{fx.cFeat, fx.feat} {
		o := version(s, "W-003", "Build on it", "review")
		o.Record.Candidate = "abcdef1"
		accept(o.Record, "owner")
		done := version(s, "W-004", "Reopened by hand", "proposed") // same commit, not in review
		done.Record.Candidate = "abcdef1"
		f.res.Groups = append(f.res.Groups, versions.Group{ID: "W-003", Versions: []versions.Version{o}}, versions.Group{ID: "W-004", Versions: []versions.Version{done}})
	}
	m := openReview(t, f, 160, 36)
	f.mu.Lock()
	reads := strings.Join(f.reads, "\n")
	f.mu.Unlock()
	if !strings.Contains(reads, "changes main abcdef1 a grove/work/W-001.md,grove/work/W-003.md") {
		t.Fatalf("the changes read leaves out both records. files:\n%s", reads)
	}
	press(m, "f")
	if s := plain(m); !strings.Contains(s, "Enter records it, returning it and W-003, which share its candidate, to active on branch feature") {
		t.Fatalf("f names the group:\n%s", s)
	}
	press(m, "esc")
	press(m, "i")
	if s := plain(m); !strings.Contains(s, "Squash branch feature onto main, delivering W-001 and W-003? y/n") || strings.Contains(s, "W-003 and W-004") {
		t.Fatalf("i names the group:\n%s", s)
	}
}

// m on a candidate that conflicts with the target opens the resolve line,
// over the branch checkout's run: defaults, and runs Conflict with the fact
// the board showed (G-260925-dz10z); without a conflict it says there is nothing to
// resolve.
func TestReviewResolveAConflictFromTheBoard(t *testing.T) {
	t.Parallel()
	fx := newFixture()
	f := reviewFixture(fx, true)
	fx.feat.Run = project.RunDefaults{BudgetUSD: "7", PermissionMode: "auto"}
	m := openReview(t, f, 200, 36)
	if s := plain(m); !strings.Contains(s, "m resolves the conflict: feedback and one attempt on its branch") || !strings.Contains(s, "m resolve") {
		t.Fatalf("the Review block and hints offer m:\n%s", s)
	}
	press(m, "m")
	s := plain(m)
	if m.prompt == nil || !strings.Contains(s, "Resolve W-001 ▏") || !strings.Contains(s, "it conflicts with main at aaaaaaa in internal/x.go: Enter records that as feedback and launches one attempt") || !strings.Contains(s, "$7, mode auto, to the handoff") {
		t.Fatalf("m opens the resolve line:\n%s", s)
	}
	typeText(m, "--until plan")
	press(m, "enter")
	if m.prompt == nil || !strings.Contains(plain(m), "--until does not apply") {
		t.Fatalf("a bound is refused on the line:\n%s", plain(m))
	}
	for range "--until plan" {
		press(m, "backspace")
	}
	typeText(m, "--resume")
	press(m, "enter")
	if m.prompt == nil || !strings.Contains(plain(m), "--resume does not apply") {
		t.Fatalf("a resume is refused on the line:\n%s", plain(m))
	}
	for range "--resume" {
		press(m, "backspace")
	}
	typeText(m, "--effort xhigh")
	deliverAll(m, press(m, "enter"))
	if s := plain(m); m.screen != resultScreen || !strings.Contains(s, "Resolution attempt of W-001") || !strings.Contains(s, "attempt: started") {
		t.Fatalf("the outcome:\n%s", s)
	}
	if strings.Join(f.acts, ";") != "conflict /repo/. W-001 abcdef1 aaaaaaa budget=7 mode=auto effort=xhigh" {
		t.Fatalf("acts %v", f.acts)
	}

	// A clean merge has nothing to resolve, and m is not offered.
	f = reviewFixture(newFixture(), false)
	f.changes = func(target, candidate, tip, path string) (*versions.Changes, error) {
		return &versions.Changes{Base: "base000", Merge: &versions.Merge{Target: strings.Repeat("a", 40), Commit: candidate, Outcome: "clean", Conflicts: []string{}}}, nil
	}
	m = openReview(t, f, 200, 36)
	if strings.Contains(plain(m), "m resolve") {
		t.Fatalf("m is offered without a conflict:\n%s", plain(m))
	}
	press(m, "m")
	if m.prompt != nil || !strings.Contains(plain(m), "merges cleanly into main at aaaaaaa, which moved since the branch left it: nothing to resolve") {
		t.Fatalf("m without a conflict:\n%s", plain(m))
	}
}

// A candidate whose branch merged the target names the merge, what it
// merged, the candidate before it, and the files it resolved.
func TestReviewNamesTheResolution(t *testing.T) {
	t.Parallel()
	f := reviewFixture(newFixture(), false)
	f.changes = func(target, candidate, tip, path string) (*versions.Changes, error) {
		return &versions.Changes{Base: "base000",
			Merge:      &versions.Merge{Target: strings.Repeat("a", 40), Commit: candidate, Outcome: "fast-forward", Conflicts: []string{}},
			Files:      []versions.Change{{Path: "internal/x.go", Added: 12, Removed: 3}, {Path: "bin.dat", Added: -1, Removed: -1}},
			Resolution: &versions.Resolution{Merge: "abcdef1", Target: strings.Repeat("a", 40), Previous: "9876543", Files: []versions.Resolved{{Path: "internal/x.go"}, {Path: "gone.go", Kept: "target"}}}}, nil
	}
	m := openReview(t, f, 200, 40)
	s := plain(m)
	for _, want := range []string{"Resolution: merge abcdef1 of main at aaaaaaa into candidate 9876543, whose reviews stay comparable · resolved: internal/x.go, gone.go (took main's side, dropping the branch's change)", "resolved in merge abcdef1"} {
		if !strings.Contains(s, want) {
			t.Fatalf("lacks %q:\n%s", want, s)
		}
	}
	if strings.Count(s, "resolved in merge") != 1 {
		t.Fatalf("only the resolved file is marked:\n%s", s)
	}
}

// A verdict given under the standing policy is told apart from the owner's
// in the Review block (G-260925-5wrn8).
func TestReviewNamesADelegatedApproval(t *testing.T) {
	t.Parallel()
	f := reviewFixture(newFixture(), true)
	for _, g := range f.res.Groups {
		for _, v := range g.Versions {
			if v.Record.ID == "W-001" && v.Record.Approved != "" {
				v.Record.Source = append(v.Record.Source, "\nVerdict on candidate abcdef1, 2026-09-25: delegated under policy grove.yaml sha256:x: review W-006.\n"...)
				v.Record.ApprovedBy = "policy sha256:x"
			}
		}
	}
	if s := plain(openReview(t, f, 200, 36)); !strings.Contains(s, "Review: candidate abcdef1 · accepted under policy ·") {
		t.Fatalf("review detail:\n%s", s)
	}
}

// The card's standing line says who approved, owner or policy, in review
// and in done alike (G-260927-60ffq).
func TestStandingLineNamesWhoApproved(t *testing.T) {
	t.Parallel()
	for _, delegated := range []bool{false, true} {
		want := "candidate abcdef1 · accepted · not on main"
		f := reviewFixture(newFixture(), true)
		if delegated {
			want = "candidate abcdef1 · accepted under policy · not on main"
			for _, g := range f.res.Groups {
				for _, v := range g.Versions {
					if v.Record.Approved != "" {
						v.Record.Source = append(v.Record.Source, "\nVerdict on candidate abcdef1, 2026-09-25: delegated under policy grove.yaml sha256:x: review W-006.\n"...)
						v.Record.ApprovedBy = "policy sha256:x"
					}
				}
			}
		}
		m := openReview(t, f, 200, 36)
		if s := plain(m); !strings.Contains(s, want) {
			t.Fatalf("review detail lacks %q:\n%s", want, s)
		}
		g := m.group()
		v := m.shown(g)
		m.standing = map[*project.Record]*standings.Standing{v.Record: {State: standings.Done, Target: "main"}}
		if got := m.detailMeta(g, v); !strings.Contains(got, want) {
			t.Fatalf("done standing line lacks %q: %s", want, got)
		}
	}
}

// rewrittenFixture is G-260928-4qv1m's incident: W-001 is done on main and
// in review, approved, on feature, whose commits main holds as rewritten
// copies, and each changed the record since they split.
func rewrittenFixture() *fake {
	fx := newFixture()
	fx.cFeat.Commit, fx.feat.Commit = strings.Repeat("b", 40), strings.Repeat("b", 40)
	var vs []versions.Version
	for _, s := range []*versions.Source{fx.cMain, fx.main} {
		v := version(s, "W-001", "Inspect records", "done")
		v.OnTarget = true
		vs = append(vs, v)
	}
	for _, s := range []*versions.Source{fx.cFeat, fx.feat} {
		v := version(s, "W-001", "Inspect records", "review")
		v.Record.Candidate = "c0ffee1"
		accept(v.Record, "owner")
		vs = append(vs, v)
	}
	res := result(fx.main, fx.sources(), vs...)
	res.Target = "main"
	return &fake{res: res, actions: true, copies: func(context.Context, string, string) (*versions.Copies, error) {
		return &versions.Copies{Commits: 4}, nil
	}}
}

func TestRewrittenCopyExplainedOnOpen(t *testing.T) {
	t.Parallel()
	f := rewrittenFixture()
	m := open(t, f, 200, 60)
	if len(f.reads) != 0 || !onRow(plain(m), "W-001", "⑂ 2 states") {
		t.Fatalf("the board load compared by patch: %v\n%s", f.reads, plain(m))
	}
	cmd := press(m, "right", "right", "enter")
	if cmd == nil || m.pending != "copies" || !strings.Contains(plain(m), "Checking whether branch feature") || !strings.Contains(m.copiesText(m.group(), " "), "Checking whether branch feature is a rewritten copy of work on main…") {
		t.Fatalf("opening the card should compare feature with main: pending %q\n%s", m.pending, plain(m))
	}
	if next := deliver(m, cmd); next != nil || strings.Join(f.reads, ",") != "copies a b" {
		t.Fatalf("reads %v", f.reads)
	}
	for _, want := range []string{"Branch feature is a rewritten copy of work already on main", "git worktree remove /repo/feat,", "git branch -D feature"} {
		if !strings.Contains(m.copiesText(m.group(), " "), want) || !strings.Contains(plain(m), "Branch feature is a rewritten copy") {
			t.Fatalf("the detail lacks %q:\n%s", want, plain(m))
		}
	}
	// The shown record is main's done; each review action names the branch in review.
	for _, k := range []string{"a", "f", "i", "m"} {
		press(m, k)
		if want := "W-001 is done on branch main and in review on branch feature: its states diverge, and v shows both and how to settle them"; m.alert != want || m.prompt != nil {
			t.Fatalf("%s: alert %q", k, m.alert)
		}
	}
	press(m, "v")
	if screen := plain(m); !strings.Contains(screen, "Diverging: 2 current states.") || !strings.Contains(screen, "Branch feature is a rewritten copy") || len(f.reads) != 1 {
		t.Fatalf("the versions screen, reads %v:\n%s", f.reads, screen)
	}

	// Partial: the commits main lacks are named, and nothing says delete.
	f.copies = func(context.Context, string, string) (*versions.Copies, error) {
		return &versions.Copies{Commits: 3, Missing: []string{strings.Repeat("d", 40)}}, nil
	}
	deliverAll(m, press(m, "esc", "r")) // the re-read card compares again
	if text := m.copiesText(m.group(), " "); !strings.Contains(text, "2 of the 3 commits of branch feature that main lacks have a copy there") ||
		!strings.Contains(text, "this one has not: ddddddd") || strings.Contains(text, "branch -D") || !strings.Contains(plain(m), "2 of the 3 commits") {
		t.Fatalf("partial: %s\n%s", text, plain(m))
	}

	// Work the target does not hold as done has no branch to clear: nothing is compared.
	for i := range f.res.Groups[0].Versions[:2] {
		f.res.Groups[0].Versions[i].Record.Status = "active"
	}
	reads := len(f.reads)
	deliverAll(m, press(m, "r"))
	if len(f.reads) != reads || m.copiesText(m.group(), " ") != "" {
		t.Fatalf("compared a branch of work not done on main: %v", f.reads[reads:])
	}
}

// Like history, the comparison never makes a key wait.
func TestRewrittenCopyReadYieldsToEveryKey(t *testing.T) {
	t.Parallel()
	f := rewrittenFixture()
	var cancelled []error
	f.copies = func(ctx context.Context, _, _ string) (*versions.Copies, error) {
		<-ctx.Done()
		cancelled = append(cancelled, ctx.Err())
		return nil, ctx.Err()
	}
	m := open(t, f, 160, 50)
	cmd := press(m, "right", "right", "enter")
	for i, key := range []string{"esc", "r", "q"} {
		if m.pending != "copies" || cmd == nil {
			t.Fatalf("before %s: pending %q", key, m.pending)
		}
		reply := make(chan tea.Msg, 1)
		go func() { reply <- cmd() }()
		next := press(m, key)
		if m.Update(<-reply); len(cancelled) != i+1 || !errors.Is(cancelled[i], context.Canceled) || len(m.copies) != 0 {
			t.Fatalf("%s during the comparison: cancelled %v, held %d", key, cancelled, len(m.copies))
		}
		switch key {
		case "esc":
			cmd = press(m, "enter")
		case "r":
			cmd = deliver(m, next) // the re-read card asks again
		}
	}
	m.reads.close()
}

// A Review card says who approved its candidate, the owner or the policy,
// and a Done card keeps it; once the board has drawn, one prediction per
// Review card, never part of the load, says whether it conflicts with the
// target, and a re-read during it cancels it (G-260928-r1hkh).
func TestReviewCardsShowApprovalAndConflicts(t *testing.T) {
	t.Parallel()
	for _, delegated := range []bool{false, true} {
		approved := "accepted"
		f := reviewFixture(newFixture(), true)
		if delegated {
			approved = "accepted under policy"
			for _, g := range f.res.Groups {
				for _, v := range g.Versions {
					if v.Record.Approved != "" {
						v.Record.Source = append(v.Record.Source, "\nVerdict on candidate abcdef1, 2026-09-25: delegated under policy grove.yaml sha256:x: review W-006.\n"...)
						v.Record.ApprovedBy = "policy sha256:x"
					}
				}
			}
		}
		var predicted []string
		block := make(chan struct{})
		b := f.backend()
		b.Predict = func(ctx context.Context, _, target string, commits []string) ([]versions.Merge, error) {
			f.mu.Lock()
			predicted = append(predicted, target+" "+strings.Join(commits, " "))
			f.mu.Unlock()
			select {
			case <-block:
			case <-ctx.Done():
				return nil, ctx.Err()
			}
			return []versions.Merge{{Target: strings.Repeat("c", 40), Commit: commits[0], Outcome: "conflict", Conflicts: []string{"internal/x.go"}}}, nil
		}
		m := New(t.Context(), "/repo/.", b)
		m.Update(tea.WindowSizeMsg{Width: 160, Height: 36})
		next := deliver(m, m.Init())
		if f.inspects != 1 || len(predicted) != 0 || next == nil || m.pending != "predict" {
			t.Fatalf("the load predicts nothing, then the board asks for one: inspects %d predicted %v pending %q", f.inspects, predicted, m.pending)
		}
		if s := plain(m); !strings.Contains(s, approved+" · not on main") || strings.Contains(s, "conflicts") {
			t.Fatalf("the card is %s and no conflict is known yet:\n%s", approved, s)
		}
		reply := make(chan tea.Msg, 1)
		go func() { reply <- next() }()
		refresh := press(m, "r")
		if _, cmd := m.Update(<-reply); cmd != nil || strings.Contains(plain(m), "conflicts") {
			t.Fatal("a re-read cancels the prediction and ignores its reply")
		}
		close(block)
		deliverAll(m, deliver(m, refresh))
		if s := flat(m); !strings.Contains(s, "conflicts main@ccccccc "+approved+" not on main") || strings.Join(predicted, ";") != "refs/heads/main abcdef1;refs/heads/main abcdef1" {
			t.Fatalf("after the re-read the prediction runs again and shows: %v\n%s", predicted, s)
		}
	}
	fx := newFixture()
	done := version(fx.cMain, "W-001", "Finished", "review")
	done.Record.Candidate = "abcdef1"
	accept(done.Record, "owner")
	m := open(t, &fake{res: result(fx.main, fx.sources(), done)}, 160, 36)
	m.standing = map[*project.Record]*standings.Standing{done.Record: {State: standings.Done, Target: "main"}}
	if s := plain(m); !onRow(s, "W-001", "accepted") {
		t.Fatalf("a Done card keeps its approval:\n%s", s)
	}
}

// Once the board has drawn, one sweep plan per re-read, never part of the
// load, tags each Review card with its act and names act and reason in the
// Review block; S asks, then sweeps in the target's checkout, and the board
// is re-read and planned again (G-260928-dtrnw).
func TestReviewCardsShowTheSweepAndSRunsIt(t *testing.T) {
	t.Parallel()
	f := reviewFixture(newFixture(), false)
	var plans, sweeps []string
	b := f.backend()
	b.SweepPlan = func(_ context.Context, root string) ([]sweep.Item, error) {
		f.mu.Lock()
		defer f.mu.Unlock()
		plans = append(plans, root)
		return []sweep.Item{{ID: "W-001", Act: sweep.Integrate, Why: "merges cleanly; verify, then approve and integrate \x1b[31m"}}, nil
	}
	b.Sweep = func(_ context.Context, root string) ([]string, error) {
		f.mu.Lock()
		defer f.mu.Unlock()
		sweeps = append(sweeps, root)
		return []string{"W-001: done: W-001 done at commit 1234567"}, nil
	}
	m := New(t.Context(), "/repo/.", b)
	m.Update(tea.WindowSizeMsg{Width: 160, Height: 36})
	next := deliver(m, m.Init())
	if f.inspects != 1 || len(plans) != 0 || m.pending != "plan" {
		t.Fatalf("the load plans nothing, then the board asks for one: plans %v pending %q", plans, m.pending)
	}
	deliverAll(m, next)
	if s := plain(m); len(plans) != 1 || plans[0] != "/repo" || !strings.Contains(s, "sweep: integrate · not on main") || !strings.Contains(s, "S sweep") {
		t.Fatalf("plans %v:\n%s", plans, s)
	}
	press(m, "S")
	if s := plain(m); !strings.Contains(s, "Sweep every candidate in review under the policy: 1 integrate? y/n   (runs in /repo)") || len(sweeps) != 0 {
		t.Fatalf("S asks first:\n%s", s)
	}
	press(m, "n")
	if m.prompt != nil || m.notice != "cancelled; nothing was swept" || len(sweeps) != 0 {
		t.Fatalf("n cancels: %q", m.notice)
	}
	deliverAll(m, press(m, "S", "y"))
	if s := plain(m); len(sweeps) != 1 || sweeps[0] != "/repo" || !strings.Contains(s, "Sweep under the policy") || !strings.Contains(s, "W-001: done: W-001 done at commit 1234567") {
		t.Fatalf("sweeps %v:\n%s", sweeps, s)
	}
	if f.inspects != 2 || len(plans) != 2 {
		t.Fatalf("the board is re-read and planned again: inspects %d plans %v", f.inspects, plans)
	}
	press(m, "esc", "right", "right")
	deliverAll(m, press(m, "enter"))
	if s := flat(m); !strings.Contains(s, "Sweep: integrate: merges cleanly; verify, then approve and integrate") || strings.Contains(m.render(), "\x1b[31m") {
		t.Fatalf("the Review block names the plan, escaped:\n%s", s)
	}

	// Without a policy the plan says so, and S refuses.
	b.SweepPlan = func(context.Context, string) ([]sweep.Item, error) {
		return nil, errors.New("grove.yaml has no policy: nothing is automatic")
	}
	m = New(t.Context(), "/repo/.", b)
	m.Update(tea.WindowSizeMsg{Width: 160, Height: 36})
	deliverAll(m, m.Init())
	press(m, "S")
	if m.prompt != nil || m.alert != "nothing to sweep: grove.yaml has no policy: nothing is automatic" || len(sweeps) != 1 {
		t.Fatalf("S refuses: %q", m.alert)
	}
}

// accept makes r the acceptance of its candidate by, as grove approve writes it.
func accept(r *project.Record, by string) {
	r.Status, r.Approved, r.ApprovedBy, r.ApprovedContext = "accepted", r.Candidate, by, project.AcceptanceContext(r)
}
