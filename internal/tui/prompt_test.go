package tui

import (
	"errors"
	"regexp"
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"

	"github.com/mascah/grove/internal/project"
)

// fits fails unless every row of the screen is exactly w cells.
func fits(t *testing.T, m *Model, w int, where string) {
	t.Helper()
	for i, row := range strings.Split(m.render(), "\n") {
		if got := ansi.StringWidth(row); got != w {
			t.Fatalf("%s: row %d is %d cells, want %d: %q", where, i, got, w, ansi.Strip(row))
		}
	}
}

// flat is the screen's words in order, box sides and separators dropped, so
// a check holds however the text wraps.
func flat(m *Model) string {
	return strings.Join(strings.Fields(strings.NewReplacer("┃", " ", "│", " ", " · ", " ").Replace(plain(m))), " ")
}

// typedOnScreen joins the prompt's rows from the one starting with lead to
// the one holding the cursor: what the owner sees of the typed text.
func typedOnScreen(m *Model, lead string) string {
	rows := strings.Split(plain(m), "\n")
	for i, r := range rows {
		if strings.HasPrefix(r, lead) {
			var b strings.Builder
			for _, r := range rows[i:] {
				b.WriteString(r)
				if strings.Contains(r, "▏") {
					break
				}
			}
			return b.String()
		}
	}
	return ""
}

// On 80 columns a verdict, a feedback text and a launch line longer than the
// width show every character typed over as many rows as they need, each
// exactly the width, and Enter records exactly what was typed (G-260928-csg91).
func TestLongPromptsShowEverythingTyped(t *testing.T) {
	t.Parallel()
	long := strings.Repeat("abcdefghij ", 12) + "end"
	for _, c := range []struct{ key, lead string }{{"a", "Verdict on W-001: "}, {"f", "Feedback on W-001: "}} {
		f := reviewFixture(newFixture(), false)
		m := openReview(t, f, 80, 24)
		press(m, c.key)
		typeText(m, long)
		fits(t, m, 80, c.key)
		if got := typedOnScreen(m, c.lead); !strings.Contains(got, c.lead+long+"▏") {
			t.Fatalf("%s: the typed text is not all on screen: %q\n%s", c.key, got, plain(m))
		}
		deliverAll(m, press(m, "enter"))
		if len(f.acts) != 1 || !strings.HasSuffix(f.acts[0], " W-001 "+long) {
			t.Fatalf("%s: Enter records exactly what was typed: %q", c.key, f.acts)
		}
	}

	fx := newFixture()
	f := &fake{res: fx.twoBranches()}
	r := &runs{}
	m := openRuns(t, f, r, 80, 24)
	m.openDetail("W-002")
	press(m, "R")
	model := "claude-" + strings.Repeat("x", 90)
	typed := "--budget 2 --permission-mode auto --model " + model
	typeText(m, typed)
	fits(t, m, 80, "launch")
	if got := typedOnScreen(m, "Launch W-002 "); !strings.Contains(got, "Launch W-002 "+typed+"▏") {
		t.Fatalf("the launch line is not all on screen: %q\n%s", got, plain(m))
	}
	settle(m, press(m, "enter"))
	if len(r.launches) != 1 || r.launches[0].Model != model {
		t.Fatalf("Enter launches what was typed: %+v", r.launches)
	}
}

// Each refusal path is drawn in the attention style and stays until Esc:
// before a prompt opens, on Enter with the line left open, and from the
// runner on the result screen, reason first (G-260928-csg91).
func TestRefusalsStayUntilDismissed(t *testing.T) {
	t.Parallel()
	styled := func(m *Model, text string) bool {
		return strings.Contains(m.render(), attention.Render(wrap(text, m.width)[0]))
	}

	// Before the prompt: work in review is not launched.
	m := openRuns(t, reviewFixture(newFixture(), false), &runs{}, 120, 36)
	m.openDetail("W-001")
	press(m, "R")
	why := "W-001 is in review: judge its candidate (a approve, f feedback) before another attempt"
	if m.prompt != nil || !styled(m, why+" · Esc dismisses") {
		t.Fatalf("the refusal is in the attention style:\n%q", m.render())
	}
	press(m, "j", "tab", "tab")
	if !styled(m, why+" · Esc dismisses") {
		t.Fatalf("the refusal stays through other keys:\n%s", plain(m))
	}
	if press(m, "esc"); strings.Contains(plain(m), why) || m.screen != detailScreen {
		t.Fatalf("Esc dismisses it and does nothing else: screen %d\n%s", m.screen, plain(m))
	}

	// On Enter: the line stays open under the refusal.
	fx := newFixture()
	f := &fake{res: fx.twoBranches()}
	for i := range f.res.Sources {
		if s := f.res.Sources[i]; s.Kind == "live" && s.GitDir == f.res.GitDir {
			s.Run = project.RunDefaults{BudgetUSD: "2", PermissionMode: "auto"}
		}
	}
	r := &runs{launchErr: errors.New("attempt W-002.20260923T005900Z of W-002 is running; stop it first")}
	m = openRuns(t, f, r, 120, 36)
	m.openDetail("W-002")
	press(m, "R")
	typeText(m, "--frob")
	press(m, "enter")
	if m.prompt == nil || !styled(m, "unknown option --frob") || strings.Contains(plain(m), "--frob · Esc dismisses") {
		t.Fatalf("Enter's refusal leaves the line open under it:\n%s", plain(m))
	}
	press(m, "backspace")
	if !styled(m, "unknown option --frob") {
		t.Fatalf("editing the line keeps the refusal:\n%s", plain(m))
	}
	for range len("--fro") {
		press(m, "backspace")
	}

	// From the runner: the result screen leads with the reason, and the
	// refusal before it is gone once the action ran.
	settle(m, press(m, "enter"))
	if m.screen != resultScreen {
		t.Fatalf("the runner's refusal shows the result screen: %d\n%s", m.screen, plain(m))
	}
	s := plain(m)
	if !styled(m, "NOT DONE: attempt W-002.20260923T005900Z of W-002 is running; stop it first") || strings.Index(s, "NOT DONE") > strings.Index(s, "The board") {
		t.Fatalf("NOT DONE leads in the attention style:\n%s", s)
	}
	if press(m, "esc"); m.screen != detailScreen || strings.Contains(plain(m), "NOT DONE") {
		t.Fatalf("Esc leaves the refusal: screen %d", m.screen)
	}
}

// A refusal takes the rows it needs, up to half the screen, and the body
// gives them up: G still reaches the last row of the content under it, and
// another screen opening settles it (G-260928-csg91, G-260928-y50a4).
func TestRefusalRowsAndTheBody(t *testing.T) {
	t.Parallel()
	m := openRuns(t, reviewFixture(newFixture(), false), &runs{}, 80, 30)
	m.openDetail("W-001")
	long := strings.Repeat("the reason goes on ", 20) + "and the next action ends it"
	m.alert = long
	press(m, "G")
	fits(t, m, 80, "long refusal")
	if !strings.Contains(flat(m), strings.Join(strings.Fields(long), " ")+" Esc dismisses") {
		t.Fatalf("the whole refusal is drawn:\n%s", plain(m))
	}
	at := regexp.MustCompile(`Content  \d+-(\d+) of (\d+)`).FindStringSubmatch(plain(m))
	if at == nil || at[1] != at[2] {
		t.Fatalf("G reaches the last row under the refusal: %q\n%s", at, plain(m))
	}
	m.alert = strings.Repeat(long+" ", 10)
	if fits(t, m, 80, "longer refusal"); !strings.Contains(plain(m), "\n…") {
		t.Fatalf("a refusal past half the screen ends in …:\n%s", plain(m))
	}
	if press(m, "A"); m.screen != attemptsScreen || m.alert != "" {
		t.Fatalf("another screen settles the refusal: screen %d alert %q", m.screen, m.alert)
	}
}
