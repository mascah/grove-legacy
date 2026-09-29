package tui

import (
	"cmp"
	"context"
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/mascah/grove/internal/attempt"
	"github.com/mascah/grove/internal/handoff"
	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/sweep"
	"github.com/mascah/grove/internal/update"
	"github.com/mascah/grove/internal/versions"
)

// The review view (G-260921-jwk4e): a work record in review shows the facts a judgment
// needs, its changed files and their diffs on demand, and three explicit
// actions, each confirmed, run through the same functions the CLI offers in
// the checkout the facts name, and for a conflict with the target a fourth,
// m, which records feedback and launches one attempt (G-260925-dz10z). Nothing else
// the board does writes.

type changesRead struct {
	c   *versions.Changes
	err error
}

type diffRead struct {
	text string
	err  error
}

type changesMsg struct {
	gen int
	key string
	c   *versions.Changes
	err error
}

type diffMsg struct {
	gen  int
	key  string
	text string
	err  error
}

type actMsg struct {
	gen   int
	kind  string
	about string // the record, or the attempt a stop was for
	facts []string
	err   error
}

// prompt is the open question on the last row: text to type for a verdict or
// feedback, or y/n for the merge and its cleanup.
type prompt struct {
	kind             string // approve, feedback, integrate, cleanup, launch, stop, resolve, conflict
	text             string
	id, root, branch string   // the record, the checkout the action runs in, the branch judged
	target, wt       string   // integrate: the target and the branch's checkout, if any
	sharing          []string // feedback, integrate: the other records sharing the candidate
	cleanup          bool
	req              *attempt.Request    // launch: what Start is asked for, resolved on Enter
	run              project.RunDefaults // launch: grove.yaml's defaults in this checkout
	attempt          string              // stop: the attempt
	expect           string              // resolve: the question's revision after the editor
	fact             *versions.Merge     // conflict: the prediction shown, which Conflict checks again
	read             []readFile          // launch of a selection: its members as the board read them here
}

// outcome is what an action returned, shown on the result screen until Esc.
type outcome struct {
	title string
	facts []string
	err   string
}

// reviewable reports the detail of a work record in review: where the
// actions apply.
func (m *Model) reviewable() bool {
	if m.screen != detailScreen {
		return false
	}
	r := m.openRecord()
	return r != nil && r.Type == "work" && r.Status == "review"
}

// openRecord is the record the detail shows now, if any.
func (m *Model) openRecord() *project.Record {
	g := m.group()
	if g == nil {
		return nil
	}
	if v := m.shown(g); v != nil {
		return v.Record
	}
	return nil
}

// changesKey names one changes read: the candidate, the tip it is judged
// against, the target, and the paths of the records sharing the candidate.
func (m *Model) changesKey(v *versions.Version) string {
	return v.Record.Candidate + "\x00" + v.Source.Commit + "\x00" + m.res.Target + "\x00" + strings.Join(m.sharing(v, true), "\x00")
}

// sharing is the work records on v's source whose candidate is v's (G-260925-wc2pz),
// v's own included, as their paths, or with paths false, the IDs of the
// others in review:
// one candidate handed off for a selection is judged per record and
// integrated as the group.
func (m *Model) sharing(v *versions.Version, paths bool) []string {
	var out []string
	for i := range m.res.Groups {
		for j := range m.res.Groups[i].Versions {
			o := &m.res.Groups[i].Versions[j]
			if o.Source != v.Source || o.Record == nil || o.Record.Type != "work" || !sameCommit(o.Record.Candidate, v.Record.Candidate) {
				continue
			}
			switch {
			case paths:
				out = append(out, o.Path)
			case o.Record.ID != v.Record.ID && o.Record.Status == "review": // what the action covers
				out = append(out, o.Record.ID)
			}
		}
	}
	slices.Sort(out)
	return out
}

// wantChanges starts reading the shown record's changes when its detail is
// open, it has a candidate, no other read is pending, and they are not held.
func (m *Model) wantChanges() tea.Cmd {
	if m.backend.Changes == nil || m.done || m.screen != detailScreen || !m.free("changes") {
		return nil
	}
	g := m.group()
	if g == nil {
		return nil
	}
	v := m.shown(g)
	if v == nil || v.Record == nil || v.Record.Candidate == "" {
		return nil
	}
	key := m.changesKey(v)
	if _, held := m.changes[key]; held || key == m.reading {
		return nil
	}
	target, candidate, tip, paths := m.res.Target, v.Record.Candidate, v.Source.Commit, m.sharing(v, true)
	cmd := m.read("changes", func(ctx context.Context, gen int) tea.Msg {
		c, err := m.backend.Changes(ctx, m.root, target, candidate, tip, paths...)
		return changesMsg{gen, key, c, err}
	})
	m.reading = key
	return cmd
}

// diffKey names one diff read.
func diffKey(from, to, path string) string { return from + "\x00" + to + "\x00" + path }

// diffOf is the diff the content pane shows, if one is chosen: the file
// against the changes' base.
func (m *Model) diffOf(v *versions.Version) (from, to, path string, ok bool) {
	if m.diff == "" || v == nil || v.Record == nil {
		return "", "", "", false
	}
	read, held := m.changes[m.changesKey(v)]
	if !held || read.c == nil || read.c.Base == "" {
		return "", "", "", false
	}
	return read.c.Base, v.Record.Candidate, m.diff, true
}

// wantDiff starts reading the chosen file's diff when nothing else is pending.
func (m *Model) wantDiff() tea.Cmd {
	if m.backend.Diff == nil || m.done || m.screen != detailScreen || !m.free("diff") {
		return nil
	}
	g := m.group()
	if g == nil {
		return nil
	}
	from, to, path, ok := m.diffOf(m.shown(g))
	if !ok {
		return nil
	}
	key := diffKey(from, to, path)
	if _, held := m.diffs[key]; held || key == m.reading {
		return nil
	}
	cmd := m.read("diff", func(ctx context.Context, gen int) tea.Msg {
		text, err := m.backend.Diff(ctx, m.root, from, to, path)
		return diffMsg{gen, key, text, err}
	})
	m.reading = key
	return cmd
}

// branchOf names the branch a version stands on, or "".
func branchOf(v *versions.Version) string {
	if v == nil || v.Source.Ref == "" {
		return ""
	}
	return strings.TrimPrefix(v.Source.Ref, "refs/heads/")
}

// projectDir is the project directory inside a live source's checkout.
func (m *Model) projectDir(s *versions.Source) string {
	return filepath.Join(s.Worktree, filepath.FromSlash(m.res.Prefix))
}

// judgeRoot finds the checkout in which the shown version can be judged,
// since approval and feedback commit the record's file there, or the reason
// there is none.
func (m *Model) judgeRoot(g *versions.Group, v *versions.Version) (root, branch, why string) {
	lv, branch, why := m.checkoutOf(g, v, "judging")
	if lv == nil {
		return "", branch, why
	}
	return m.projectDir(lv.Source), branch, ""
}

// checkoutOf finds the checkout that holds the shown version to write it:
// the one valid checkout on its branch holding the record. It returns that
// copy, or the reason there is none. doing is "judging", which needs the
// copy to match HEAD, or "answering", whose uncommitted copy is the owner's
// answer so far.
func (m *Model) checkoutOf(g *versions.Group, v *versions.Version, doing string) (lv *versions.Version, branch, why string) {
	judging := doing == "judging"
	if v == nil || v.Record == nil {
		return nil, "", "the current state holds no record to " + map[string]string{"judging": "judge", "answering": "answer", "editing": "edit"}[doing]
	}
	branch = branchOf(v)
	if branch == "" {
		return nil, "", "the current state is on no branch; " + map[bool]string{true: "a and f need", false: "e needs"}[judging] + " a branch checkout"
	}
	where := func(s *versions.Source) string {
		if s.Locator == "." {
			return "checkout ."
		}
		return "checkout " + s.Locator
	}
	dirty := func(s *versions.Source) string {
		return "the record has uncommitted changes in " + where(s) + "; commit or discard them before " + doing + " it"
	}
	if v.Source.Kind == "live" {
		if judging && v.Change != "unchanged" {
			return nil, branch, dirty(v.Source)
		}
		return v, branch, ""
	}
	var found []*versions.Version
	for _, s := range m.res.Sources {
		if s.Kind != "live" || s.Ref != v.Source.Ref || !s.Valid {
			continue
		}
		for i := range g.Versions {
			if cv := &g.Versions[i]; cv.Source == s {
				found = append(found, cv)
			}
		}
	}
	switch {
	case len(found) > 1:
		return nil, branch, fmt.Sprintf("%d checkouts are on branch %s, so which one to write is ambiguous", len(found), branch)
	case len(found) == 0 && judging:
		return nil, branch, "no checkout is on branch " + branch + "; git worktree add one, or run grove approve there"
	case len(found) == 0:
		return nil, branch, "no checkout is on branch " + branch + "; git worktree add one"
	case found[0].Record == nil || judging && found[0].Change != "unchanged":
		return nil, branch, dirty(found[0].Source)
	}
	return found[0], branch, ""
}

// targetRoot finds the target's checkout, where integrate runs.
func (m *Model) targetRoot() (root, why string) {
	if m.res.Target == "" {
		return "", "grove.yaml names no target branch, so nothing can be integrated"
	}
	for _, s := range m.res.Sources {
		if s.Kind == "live" && s.Ref == "refs/heads/"+m.res.Target && s.Valid {
			return m.projectDir(s), ""
		}
	}
	return "", "no checkout is on the target " + m.res.Target
}

// targetTip is the target branch's commit as the board read it, or "".
func (m *Model) targetTip() string {
	for _, s := range m.res.Sources {
		if s.Kind == "committed" && s.Ref == "refs/heads/"+m.res.Target {
			return s.Commit
		}
	}
	return ""
}

// approval is who approved r's candidate, as the header shows it, or ""
// when it is not approved: a delegated verdict is told apart from the
// owner's (G-260925-wh9ax).
func approval(r *project.Record) string {
	switch {
	case r.Approved == "":
		return ""
	case update.Delegated(r):
		return "approved under policy"
	}
	return "approved"
}

// predictMsg is one Review card's predicted merge into the target.
type predictMsg struct {
	gen   int
	key   string
	merge *versions.Merge // nil where no prediction could be made
}

// predictKey names one card's prediction: the target as the board read it
// and the candidate.
func (m *Model) predictKey(candidate string) string { return m.targetTip() + "\x00" + candidate }

// wantPredict starts predicting the next Review card's merge into the
// target while the board shows, one card at a time and never as part of its
// load (G-260928-r1hkh): three Git processes a card, which any other read
// replaces, and a re-read forgets.
func (m *Model) wantPredict() tea.Cmd {
	if m.backend.Predict == nil || m.done || m.screen != boardScreen || m.res == nil || m.res.Target == "" || m.targetTip() == "" || !m.free("predict") {
		return nil
	}
	columns, _ := m.placed()
	for _, c := range columns[reviewColumn] {
		if c.rec == nil || c.rec.Candidate == "" {
			continue
		}
		key := m.predictKey(c.rec.Candidate)
		if key == m.reading {
			return nil
		}
		if _, held := m.predicts[key]; held {
			continue
		}
		target, candidate := "refs/heads/"+m.res.Target, c.rec.Candidate
		cmd := m.read("predict", func(ctx context.Context, gen int) tea.Msg {
			merges, err := m.backend.Predict(ctx, m.root, target, []string{candidate})
			if err != nil || len(merges) != 1 {
				return predictMsg{gen, key, nil}
			}
			return predictMsg{gen, key, &merges[0]}
		})
		m.reading = key
		return cmd
	}
	return nil
}

// planMsg is the sweep's plan, read in the target's checkout.
type planMsg struct {
	gen   int
	items []sweep.Item
	err   error
}

// wantPlan plans a sweep once per re-read, when no other read is pending
// and a Review card has a candidate, never as part of the load and never
// acting (G-260928-dtrnw): the board sweeps only on S.
func (m *Model) wantPlan() tea.Cmd {
	if m.backend.SweepPlan == nil || m.done || m.pending != "" || m.res == nil || m.planned != nil || m.planErr != "" {
		return nil
	}
	root, why := m.targetRoot()
	if why != "" {
		return nil
	}
	columns, _ := m.placed()
	if !slices.ContainsFunc(columns[reviewColumn], func(c card) bool { return c.rec != nil && c.rec.Candidate != "" }) {
		return nil
	}
	return m.read("plan", func(ctx context.Context, gen int) tea.Msg {
		items, err := m.backend.SweepPlan(ctx, root)
		return planMsg{gen, items, err}
	})
}

// planText is the Review block's row for the sweep's plan for id, or ""
// when there is none to show.
func (m *Model) planText(id string) string {
	if m.backend.SweepPlan == nil {
		return ""
	}
	if m.planErr != "" {
		return "Sweep: not planned: " + m.planErr
	}
	if it, ok := m.planned[id]; ok {
		return "Sweep: " + it.Act + ": " + it.Why + " (S sweeps from the board)"
	}
	return ""
}

// askSweep opens the y/n line for S, or says why there is nothing to sweep.
func (m *Model) askSweep() {
	if m.backend.Sweep == nil || m.res == nil {
		return
	}
	root, why := m.targetRoot()
	switch {
	case why != "":
		m.alert = "a sweep runs in the target's checkout: " + why
		return
	case m.planErr != "":
		m.alert = "nothing to sweep: " + m.planErr
		return
	}
	m.prompt = &prompt{kind: "sweep", root: root, target: m.res.Target}
}

// conflictNote is a Review card's predicted conflict with the target,
// naming the target commit it read, or "" when none is predicted or read.
func (m *Model) conflictNote(r *project.Record) string {
	if mg := m.predicts[m.predictKey(r.Candidate)]; r.Candidate != "" && mg != nil && mg.Outcome == "conflict" {
		return "conflicts " + m.res.Target + "@" + short7(mg.Target)
	}
	return ""
}

// reviewRows are the header's Review block: the candidate's standing and
// where each action would run. Facts only; every action is a key away.
func (m *Model) reviewRows(g *versions.Group, v *versions.Version) []string {
	r := v.Record
	parts := []string{"Review: candidate " + short7(r.Candidate)}
	parts = append(parts, cmp.Or(approval(r), "not yet approved"))
	switch read, held := m.changes[m.changesKey(v)]; {
	case m.backend.Changes == nil:
	case !held:
		parts = append(parts, "reading its changes…")
	case read.err != nil:
		parts = append(parts, "changes unreadable (r retries)")
	case len(read.c.After) != 0:
		parts = append(parts, fmt.Sprintf("%d other %s changed since it: the tip is a new candidate", len(read.c.After), plural(len(read.c.After), "file", "files")))
	default:
		parts = append(parts, "only the record changed since it")
	}
	if read, held := m.changes[m.changesKey(v)]; held && read.c != nil && read.c.Merge != nil {
		// The prediction names the target commit it read, which the board's
		// own reading of the target may no longer be (G-260925-h8rj5).
		text := read.c.Merge.Text(m.res.Target)
		if tip := m.targetTip(); tip != "" && tip != read.c.Merge.Target {
			text += "; the board read " + m.res.Target + " at " + short7(tip) + " (r re-reads)"
		}
		parts = append(parts, text)
	} else if held && read.c != nil && read.c.Unpredicted != "" {
		parts = append(parts, "the merge into "+m.res.Target+" could not be predicted: "+read.c.Unpredicted)
	}
	rows := []string{strings.Join(parts, " · ")}
	if read, held := m.changes[m.changesKey(v)]; held && read.c != nil && read.c.Resolution != nil {
		rows = append(rows, resolutionText(read.c.Resolution, m.res.Target))
	}
	var where []string
	if root, branch, why := m.judgeRoot(g, v); why != "" {
		where = append(where, "a approve and f feedback: "+why)
	} else {
		where = append(where, "a approve and f feedback run on branch "+branch+" in "+root)
	}
	if m.backend.Integrate != nil {
		if root, why := m.targetRoot(); why != "" {
			where = append(where, "i integrate: "+why)
		} else {
			where = append(where, "i integrate runs into "+m.res.Target+" in "+root)
		}
	}
	if m.conflicted(v) != nil && m.backend.Conflict != nil {
		where = append(where, "m resolves the conflict: feedback and one attempt on its branch")
	}
	if text := m.planText(g.ID); text != "" {
		rows = append(rows, text)
	}
	return append(rows, strings.Join(where, " · "))
}

// changeEntries lists the changed files as sidebar entries, in Git's order.
func (m *Model) changeEntries(v *versions.Version) []entry {
	if v == nil || v.Record == nil || v.Record.Candidate == "" {
		return nil
	}
	read, held := m.changes[m.changesKey(v)]
	if !held || read.c == nil {
		return nil
	}
	var out []entry
	for i := range read.c.Files {
		out = append(out, entry{change: &read.c.Files[i]})
	}
	return out
}

// changesSection is the sidebar's Changes heading and rows, with item
// drawing each selectable file.
func (m *Model) changesSection(v *versions.Version, w int, heading func(string), item func(string), plain func(string)) {
	if m.backend.Changes == nil || v == nil || v.Record == nil || v.Record.Candidate == "" {
		return
	}
	read, held := m.changes[m.changesKey(v)]
	switch {
	case !held:
		heading("Changes")
		plain("  reading…")
	case read.err != nil:
		heading("Changes")
		plain("  The changes could not be read (r retries): " + read.err.Error())
	case m.res.Target == "":
		heading("Changes")
		plain("  grove.yaml names no target, so there is no base to list files against")
	default:
		heading(fmt.Sprintf("Changes against %s from %s", m.res.Target, short7(read.c.Base)))
		if len(read.c.Files) == 0 {
			plain("  none")
		}
		for _, f := range read.c.Files {
			counts := "binary"
			if f.Added >= 0 {
				counts = fmt.Sprintf("+%d −%d", f.Added, f.Removed)
			}
			// One row a file: the records describing it are counted here and
			// named at the head of its diff.
			if n := len(m.describedBy(f.Path)); n != 0 {
				counts += fmt.Sprintf("  described by %d", n)
			}
			if r := read.c.Resolution; r != nil && slices.ContainsFunc(r.Files, func(x versions.Resolved) bool { return x.Path == f.Path }) {
				counts += "  resolved in merge " + short7(r.Merge)
			}
			item(ansi.Truncate(safe(f.Path), max(w-2-ansi.StringWidth(counts)-2, 8), "…") + "  " + counts)
		}
		if len(read.c.After) != 0 {
			plain("  after the candidate: " + strings.Join(read.c.After, ", "))
		}
	}
}

// describedBy lists the records, other than the open one, that link a
// changed file or name it in a code span (G-260925-dzxm6), each once at its first
// tier, in the inspection's order. It reads the loaded records only: no Git
// process, nothing stored. A file is a path from the repository's top, or a
// rename's two; the prefix, which ends in a slash, makes it a project path,
// and a file outside the project can only be named by a code span.
// ponytail: matched on every frame, about 20 ms for 60 files over 150
// records; cache by changes key if reviews grow larger.
func (m *Model) describedBy(file string) []string {
	var found []string
	for i := range m.res.Groups {
		g := &m.res.Groups[i]
		if g.ID == m.openID() {
			continue
		}
		best := -1
		vs, _ := m.searched(g)
		for _, v := range vs {
			ms := m.mentionsOf(v)
			for _, p := range strings.Split(file, " → ") {
				tier := -1
				if rest, in := strings.CutPrefix(p, m.res.Prefix); in {
					tier, _ = pathTier(rest, ms)
				} else if slices.ContainsFunc(ms, func(mn handoff.Mention) bool { return mn.Span != "" && names(mn.Span, p) }) {
					tier = 2
				}
				if tier >= 0 && (best < 0 || tier < best) {
					best = tier
				}
			}
		}
		if best >= 0 && !slices.ContainsFunc(found, func(f string) bool { return strings.HasPrefix(f, g.ID+" ") }) {
			found = append(found, g.ID+" "+tiers[best])
		}
	}
	return found
}

// diffHead names, above a file's diff, the records that describe the file
// and the merge that resolved it, if one did.
func (m *Model) diffHead(read changesRead, path string, w int) []string {
	file := path
	for _, f := range read.c.Files {
		if f.Path == path || strings.HasSuffix(f.Path, " → "+path) {
			file = f.Path
		}
	}
	text := "No record names this file."
	if found := m.describedBy(file); len(found) != 0 {
		text = "Described by " + strings.Join(found, ", ") + "."
	}
	if r := read.c.Resolution; r != nil && slices.ContainsFunc(r.Files, func(x versions.Resolved) bool { return x.Path == file }) {
		text += " Resolved in merge " + short7(r.Merge) + "."
	}
	return append(wrap(text, w), line("", w))
}

// diffRows renders a diff for the content pane: every line escaped like
// record text, then an accent by its first character that nothing depends
// on. ponytail: rendered on every frame; cache by key and width if diffs
// grow large.
func diffRows(text string, w int) []string {
	var rows []string
	for _, l := range strings.Split(strings.TrimRight(text, "\n"), "\n") {
		row := line(strings.ReplaceAll(l, "\t", "    "), w)
		switch {
		case strings.HasPrefix(l, "+") && !strings.HasPrefix(l, "+++"):
			row = "\x1b[32m" + row + "\x1b[m"
		case strings.HasPrefix(l, "-") && !strings.HasPrefix(l, "---"):
			row = "\x1b[31m" + row + "\x1b[m"
		case strings.HasPrefix(l, "@@"):
			row = "\x1b[36m" + row + "\x1b[m"
		}
		rows = append(rows, row)
	}
	return rows
}

// diffContent is what the content pane shows for the chosen file.
func (m *Model) diffContent(v *versions.Version, w int) []string {
	from, to, path, ok := m.diffOf(v)
	if !ok {
		return wrapAll("The changes are not held any more; r rereads them.", w)
	}
	read, held := m.diffs[diffKey(from, to, path)]
	switch {
	case !held:
		return wrapAll("reading the diff…", w)
	case read.err != nil:
		return wrapAll("The diff could not be read (r retries): "+read.err.Error(), w)
	}
	head := m.diffHead(m.changes[m.changesKey(v)], path, w)
	if strings.TrimSpace(read.text) == "" {
		return append(head, wrapAll("No textual difference.", w)...)
	}
	return append(head, diffRows(read.text, w)...)
}

// startAtEvidence scrolls a review's content to its Evidence heading, so the
// first screen is the handoff.
func (m *Model) startAtEvidence() {
	g := m.group()
	if g == nil {
		return
	}
	v := m.shown(g)
	if v == nil || v.Record == nil || v.Record.Type != "work" || v.Record.Status != "review" {
		return
	}
	w := m.width
	if w >= wideWidth {
		w = w * 11 / 20
	}
	for i, row := range m.contentRows(v, w) {
		if strings.HasPrefix(strings.TrimSpace(ansi.Strip(row)), "## Evidence") {
			m.dscroll = i
			m.clampScroll()
			return
		}
	}
}

// action opens the prompt for a, f or i on a work record in review, or says
// why it cannot.
func (m *Model) action(k string) {
	g := m.group()
	v := m.shown(g)
	if !m.reviewable() {
		if v != nil && v.Record != nil {
			m.alert = m.notInReview(g, v, "nothing to approve, give feedback on, or integrate")
		}
		return
	}
	r := v.Record
	switch k {
	case "a", "f":
		root, branch, why := m.judgeRoot(g, v)
		if why != "" {
			m.alert = why
			return
		}
		if k == "a" && r.Approved != "" {
			m.alert = "candidate " + short7(r.Candidate) + " is already approved; i integrates it"
			return
		}
		m.prompt = &prompt{kind: map[string]string{"a": "approve", "f": "feedback"}[k], id: g.ID, root: root, branch: branch, sharing: m.sharing(v, false)}
	case "i":
		if m.backend.Integrate == nil {
			return
		}
		if r.Approved == "" {
			m.alert = "approve candidate " + short7(r.Candidate) + " first (a)"
			return
		}
		root, why := m.targetRoot()
		if why != "" {
			m.alert = why
			return
		}
		wt, branch, _ := m.judgeRoot(g, v)
		m.prompt = &prompt{kind: "integrate", id: g.ID, root: root, branch: branch, target: m.res.Target, wt: wt, sharing: m.sharing(v, false)}
	}
}

// notInReview says why the shown version v offers no review action: its
// record is not in review, or, on a card whose states diverge, another
// current state is, and the notice names where (G-260928-4qv1m).
func (m *Model) notInReview(g *versions.Group, v *versions.Version, what string) string {
	for _, state := range currentStates(*g) {
		if r := state[0].Record; !m.current() || r == nil || r.Status != "review" || slices.Contains(state, v) {
			continue
		}
		var places []string
		for _, o := range state {
			if o.Source.Kind == "committed" {
				places = append(places, label(o.Source))
			}
		}
		if places == nil {
			places = []string{label(state[0].Source)}
		}
		return g.ID + " is " + v.Record.Status + " on " + label(v.Source) + " and in review on " + strings.Join(places, ", ") + ": its states diverge, and v shows both and how to settle them"
	}
	return g.ID + " is not in review: " + what
}

// conflicted is the shown candidate's conflict with the target, when its
// changes are read and predict one.
func (m *Model) conflicted(v *versions.Version) *versions.Merge {
	if v == nil || v.Record == nil || v.Record.Status != "review" {
		return nil
	}
	if read, held := m.changes[m.changesKey(v)]; held && read.c != nil && read.c.Merge != nil && read.c.Merge.Outcome == "conflict" {
		return read.c.Merge
	}
	return nil
}

// resolveConflict opens the resolve line for a candidate in review that
// conflicts with the target (G-260925-dz10z), or says why not: it records the
// feedback in the branch's checkout and launches one attempt there, with
// the defaults that checkout's grove.yaml sets.
func (m *Model) resolveConflict() {
	g := m.group()
	v := m.shown(g)
	if !m.reviewable() || m.backend.Conflict == nil {
		if v != nil && v.Record != nil {
			m.alert = m.notInReview(g, v, "there is no candidate to resolve")
		}
		return
	}
	read, held := m.changes[m.changesKey(v)]
	fact := m.conflicted(v)
	switch {
	case !held || read.c == nil:
		m.alert = "the candidate's changes are not read yet; m waits for them"
		return
	case fact == nil && read.c.Merge != nil:
		m.alert = "candidate " + short7(v.Record.Candidate) + " " + read.c.Merge.Text(m.res.Target) + ": nothing to resolve"
		return
	case fact == nil:
		m.alert = "no conflict with the target is predicted, so there is nothing to resolve"
		return
	}
	lv, branch, why := m.checkoutOf(g, v, "judging")
	if why != "" {
		m.alert = why
		return
	}
	for _, a := range m.attemptsOf(g.ID) {
		if live(&a) {
			m.alert = fmt.Sprintf("attempt %s of %s is %s; A shows it, x stops it", a.Launch.Attempt, g.ID, a.Status)
			return
		}
	}
	req := attempt.Request{Root: m.root, IDs: []string{g.ID}}
	m.prompt = &prompt{kind: "conflict", id: g.ID, root: m.root, branch: branch, target: m.res.Target, wt: lv.Source.Worktree, req: &req, run: lv.Source.Run, fact: fact}
}

// resolutionText is the Review block's row for a target merge on the branch.
func resolutionText(r *versions.Resolution, target string) string {
	text := "Resolution: merge " + short7(r.Merge) + " of " + target + " at " + short7(r.Target)
	if r.Previous != "" {
		text += " into candidate " + short7(r.Previous) + ", whose reviews stay comparable"
	}
	if len(r.Files) == 0 {
		return text + " · it merged without a conflict"
	}
	var files []string
	for _, f := range r.Files {
		switch f.Kept {
		case "target":
			files = append(files, f.Path+" (took "+target+"'s side, dropping the branch's change)")
		case "branch":
			files = append(files, f.Path+" (kept the branch's side, dropping "+target+"'s change)")
		default:
			files = append(files, f.Path)
		}
	}
	return text + " · resolved: " + strings.Join(files, ", ")
}

// promptKey handles every key while a prompt is open: text goes into it,
// Enter and y/n answer it, Esc cancels it.
func (m *Model) promptKey(msg tea.KeyPressMsg) tea.Cmd {
	p := m.prompt
	switch k := msg.String(); {
	case k == "esc":
		m.prompt, m.alert = nil, ""
		m.notice = "cancelled; nothing was written"
		if p.req != nil {
			m.notice = "cancelled; nothing was launched"
		} else if p.kind == "stop" {
			m.notice = "cancelled; nothing was stopped"
		} else if p.kind == "resolve" {
			m.notice = unresolved(p)
		}
	case msg.Text == "y" && p.kind == "stop":
		return m.act(p)
	case msg.Text == "n" && p.kind == "stop":
		m.prompt = nil
		m.notice = "cancelled; nothing was stopped"
	case msg.Text == "y" && p.kind == "resolve":
		return m.act(p)
	case msg.Text == "y" && p.kind == "sweep":
		return m.act(p)
	case msg.Text == "n" && p.kind == "sweep":
		m.prompt = nil
		m.notice = "cancelled; nothing was swept"
	case msg.Text == "n" && p.kind == "resolve":
		m.prompt = nil
		m.notice = unresolved(p)
	case p.kind == "approve" || p.kind == "feedback" || p.req != nil:
		switch k {
		case "enter":
			if p.req != nil {
				return m.launchKey(p)
			}
			if strings.TrimSpace(p.text) == "" {
				m.alert = "type the " + map[string]string{"approve": "verdict", "feedback": "feedback"}[p.kind] + " first, or Esc"
				return nil
			}
			return m.act(p)
		case "backspace":
			if n := len(p.text); n != 0 {
				_, size := lastRune(p.text)
				p.text = p.text[:n-size]
			}
		default:
			if msg.Text != "" {
				p.text += msg.Text
			}
		}
	case msg.Text == "y" && p.kind == "integrate":
		if p.wt == "" {
			return m.act(p) // no worktree: nothing to ask about
		}
		p.kind = "cleanup"
	case msg.Text == "y" && p.kind == "cleanup":
		p.cleanup = true
		return m.act(p)
	case msg.Text == "n" && p.kind == "cleanup":
		return m.act(p)
	case msg.Text == "n":
		m.prompt = nil
		m.notice = "cancelled; nothing was merged"
	}
	return nil
}

// act runs the prompt's action in the background. Keys wait for it, and a
// key never cancels its Git commands.
func (m *Model) act(p *prompt) tea.Cmd {
	m.prompt, m.alert = nil, "" // what the line was refused for is settled; the result says the rest
	kind, root, id, text, cleanup, req, about, expect, fact := p.kind, p.root, p.id, strings.TrimSpace(p.text), p.cleanup, p.req, p.id, p.expect, p.fact
	switch {
	case kind == "cleanup":
		kind = "integrate"
	case req != nil && kind != "conflict":
		kind = "launch"
	case kind == "stop":
		about = p.attempt
	}
	m.resultBack = m.screen
	cmd := m.read("act", func(ctx context.Context, gen int) tea.Msg {
		var facts []string
		var err error
		switch kind {
		case "launch":
			facts, err = m.backend.Launch(ctx, *req)
		case "conflict":
			facts, err = m.backend.Conflict(ctx, *req, fact)
		case "stop":
			facts, err = m.backend.Stop(ctx, root, about)
		case "approve":
			facts, err = m.backend.Approve(ctx, root, id, text)
		case "feedback":
			facts, err = m.backend.Feedback(ctx, root, id, text)
		case "resolve":
			facts, err = m.backend.Answer(ctx, root, id, expect)
		case "sweep":
			facts, err = m.backend.Sweep(ctx, root)
		default:
			facts, err = m.backend.Integrate(ctx, root, id, cleanup)
		}
		return actMsg{gen, kind, about, facts, err}
	})
	m.acting = kind
	return cmd
}

// acting names the running action for the banner.
func actingText(kind string) string {
	return map[string]string{"approve": "Approving…", "feedback": "Recording the feedback…", "integrate": "Integrating…",
		"launch": "Launching the attempt…", "stop": "Stopping the attempt…", "resolve": "Resolving and committing…",
		"conflict": "Recording the feedback and launching the attempt…", "sweep": "Sweeping under the policy…"}[kind]
}

// promptText is the open prompt as what is typed, led by what it is for,
// and the help after it; a y/n prompt has only the help.
func (m *Model) promptText() (typed, help string) {
	p := m.prompt
	switch p.kind {
	case "approve":
		return "Verdict on " + p.id + ": " + p.text, "Enter approves " + p.id + " on branch " + p.branch + " with it; Esc cancels"
	case "feedback":
		them := "it"
		if p.sharing != nil {
			them = "it and " + strings.Join(p.sharing, ", ") + ", which share its candidate,"
		}
		return "Feedback on " + p.id + ": " + p.text, "Enter records it, returning " + them + " to active on branch " + p.branch + "; Esc cancels"
	case "launch":
		req, _ := p.resolved() // a line that does not parse yet shows what it has so far
		return "Launch " + p.id + " " + p.text, launchText(req) + ", " + p.where() + " · Enter launches; flags such as --until plan or --effort xhigh change it; Esc cancels"
	case "conflict":
		req, _ := p.resolved()
		return "Resolve " + p.id + " " + p.text, "it " + p.fact.Text(p.target) + ": Enter records that as feedback and launches one attempt to merge it, resolve, verify and hand off, " + launchText(req) + ", " + p.where() + "; Esc cancels"
	case "integrate":
		// The question first: a long checkout path is what truncation drops.
		return "", fmt.Sprintf("Merge branch %s into %s and mark %s done? y/n   (runs in %s)", p.branch, p.target, strings.Join(append([]string{p.id}, p.sharing...), " and "), p.root)
	case "stop":
		return "", fmt.Sprintf("Stop attempt %s of %s? Its partial work stays. y/n", p.attempt, p.id)
	case "resolve":
		return "", fmt.Sprintf("Resolve %s and commit it with your answer on branch %s? y/n   (runs in %s)", p.id, p.branch, p.root)
	case "sweep":
		return "", fmt.Sprintf("Sweep every candidate in review under the policy: %s? y/n   (runs in %s)", m.planSummary(), p.root)
	}
	return "", fmt.Sprintf("Also delete branch %s and remove its worktree? y/n   (%s)", p.branch, p.wt)
}

// promptRows draws the open prompt in at most n rows of w cells. A y/n
// prompt is one row. Typed text is broken only at the edge, so every
// character shows, and comes first; the help follows, and is what a lack
// of rows drops, then the typed text's start, never its end.
func (m *Model) promptRows(w, n int) []string {
	typed, help := m.promptText()
	if typed == "" {
		return []string{hot(line(help, w))}
	}
	rows := exact(typed+"▏", w)
	rows = append(rows[max(len(rows)-n, 0):], wrap(help, w)...)
	for i := range rows {
		rows[i] = hot(rows[i])
	}
	return rows[:min(len(rows), max(n, 1))]
}

// footer is the rows under the body: a refusal, then the open prompt or the
// key hints, in at most half the screen; a refusal longer than that ends in
// "…".
func (m *Model) footer(w int, hints string) []string {
	n := max((m.height-2)/2, 1)
	var rows []string
	if m.alert != "" {
		text := m.alert
		if m.prompt == nil {
			text += " · Esc dismisses"
		}
		lines := wrap(text, w)
		if limit := max(n-1, 1); len(lines) > limit {
			lines = append(lines[:limit-1], line("…", w))
		}
		for _, r := range lines {
			rows = append(rows, attention.Render(r))
		}
	}
	if m.prompt == nil {
		return append(rows, line(hints, w))
	}
	return append(rows, m.promptRows(w, max(n-len(rows), 1))...)
}

// attention marks a refusal, in a colour no column or tag uses; its text
// says the same.
var attention = lipgloss.NewStyle().Bold(true).Reverse(true).Foreground(lipgloss.Color("1"))

// resultRows is the result screen: the action's facts, and the refusal or
// failure when there was one.
func (m *Model) resultRows(w int) []string {
	o := m.result
	rows := []string{bold(line(o.title, w))}
	// A refusal leads, in the attention style, before whatever it reported.
	if o.err != "" {
		for _, r := range wrap("NOT DONE: "+o.err, w) {
			rows = append(rows, attention.Render(r))
		}
		rows = append(rows, line("", w))
	}
	for _, f := range o.facts {
		rows = append(rows, wrap("  "+f, w)...)
	}
	if m.pending == "inspect" {
		return append(rows, line("", w), line("Re-reading the board…", w))
	}
	return append(rows, line("", w), line("The board has been re-read. Esc returns to the record.", w))
}

// unresolved says what declining to resolve leaves.
func unresolved(p *prompt) string {
	return p.id + " is not resolved; your edit stays uncommitted in " + p.root + ", and e reopens it to resolve"
}

func short7(commit string) string { return commit[:min(len(commit), 7)] }

// planSummary counts the plan's acts for S's question, in the order a sweep
// takes them, or says the plan is not read.
func (m *Model) planSummary() string {
	if m.planned == nil {
		return "its plan is not read yet"
	}
	var parts []string
	for _, act := range []string{sweep.Integrate, sweep.Approve, sweep.Resolve, sweep.Wait, sweep.Skip} {
		n := 0
		for _, it := range m.planned {
			if it.Act == act {
				n++
			}
		}
		if n != 0 {
			parts = append(parts, fmt.Sprintf("%d %s", n, act))
		}
	}
	if parts == nil {
		return "no candidate in review"
	}
	return strings.Join(parts, ", ")
}
