package sweep

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/mascah/grove/internal/project"
)

// These cases catch accidental approval of incomplete, ambiguous or legacy
// nonzero reports, and any broadening of an existing owner's delegation.
func TestReviewSummaryDelegation(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, report string
		allow, pass  bool
	}{
		{"clean", "Review v1: complete; blockers=0; follow-ups=0", false, true},
		{"follow-ups default", "Review v1: complete; blockers=0; follow-ups=1", false, false},
		{"follow-ups opted in", "Review v1: complete; blockers=0; follow-ups=1", true, true},
		{"blocking documentation", "The command reference permits an unsafe recovery.\nReview v1: complete; blockers=1; follow-ups=0", true, false},
		{"incomplete", "Review v1: incomplete; blockers=0; follow-ups=0", true, false},
		{"missing count", "Review v1: complete; blockers=0", true, false},
		{"negative count", "Review v1: complete; blockers=-1; follow-ups=0", true, false},
		{"overflow", "Review v1: complete; blockers=0; follow-ups=999999999999999999999999999", true, false},
		{"unknown version", "Review v2: complete; blockers=0; follow-ups=0", true, false},
		{"trailing prose", "Review v1: complete; blockers=0; follow-ups=0\nStill checking.", true, false},
		{"legacy clean", "Open findings: none", false, true},
		{"legacy open", "Open findings: 1", true, false},
		{"legacy cannot hide incomplete", "Review v1: incomplete; blockers=0; follow-ups=0\nOpen findings: none", true, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root := t.TempDir()
			write(t, root, "grove/.keep", "")
			opt := ""
			if tc.allow {
				opt = "    allow_followups: true\n"
			}
			write(t, root, "grove.yaml", "schema_version: 4\nrecords: grove\npolicy:\n  approve:\n    verify: ['true']\n"+opt)
			p, ds := project.Load(root, root)
			if len(ds) != 0 {
				t.Fatalf("policy: %v", ds)
			}
			s := &Sweep{Policy: p.Policy}
			r := &project.Record{ID: "G-260101-00001", Candidate: "abcdef1234"}
			rev := &project.Record{ID: "G-260101-00005", Type: "review", Status: "current", Work: []string{r.ID}, Examined: r.Candidate, Source: []byte(tc.report + "\n")}
			got, why := s.review(context.Background(), r, []*project.Record{rev})
			if (got != nil) != tc.pass || got == nil && why == "" {
				t.Fatalf("review = %v, reason %q; want eligible=%v", got, why, tc.pass)
			}
			if tc.pass {
				if tc.allow {
					clean := *rev
					clean.ID = "G-260101-00007"
					clean.Source = []byte(ClosingLine)
					got, _ := s.review(context.Background(), r, []*project.Record{rev, &clean})
					if got == nil || !strings.Contains(got.text(), "G-260101-00005: 1 follow-ups") {
						t.Fatal("a later clean review hid the earlier review's follow-up")
					}
				}
				blocked := *rev
				blocked.ID = "G-260101-00006"
				blocked.Source = []byte("Review v1: complete; blockers=1; follow-ups=0\n")
				for _, reviews := range [][]*project.Record{{rev, &blocked}, {&blocked, rev}} {
					if got, _ := s.review(context.Background(), r, reviews); got != nil {
						t.Fatal("one eligible review outvoted another review's blocker")
					}
				}
			}
		})
	}
}

// A real policy preview and delivery must retain follow-up disclosure, and
// code changed after the examined commit must still invalidate that review.
func TestSweepFollowupsAndStaleReview(t *testing.T) {
	t.Parallel()
	config := strings.Replace(policy, "%s", "'true'", 1)
	config = strings.Replace(config, "  approve:\n", "  approve:\n    allow_followups: true\n", 1)
	root, wt := fixture(t, config, map[string]string{"code.txt": "change\n"}, "Review v1: complete; blockers=0; follow-ups=2")
	s, err := Plan(root)
	if err != nil || len(s.Items) != 1 || s.Items[0].Act != Integrate || !strings.Contains(s.Items[0].Why, "2 follow-ups") {
		t.Fatalf("preview: %+v, %v", s, err)
	}
	t.Logf("policy preview: %s", s.Items[0].Why)
	write(t, wt, "code.txt", "later code\n")
	git(t, wt, "commit", "-qam", "fix: new candidate")
	newCandidate := git(t, wt, "rev-parse", "HEAD")
	r := record(t, wt)
	r.Candidate = newCandidate
	p, _ := project.Load(wt, wt)
	if got, why := s.review(context.Background(), r, p.Records); got != nil || !strings.Contains(why, "no current review") {
		t.Fatalf("changed code still reviewed: %v %s", got, why)
	}
	// Use a second independent fixture for delivery; the changed candidate
	// above stays untouched rather than being reset just to satisfy a test.
	root, _ = fixture(t, config, map[string]string{"code.txt": "change\n"}, "Review v1: complete; blockers=0; follow-ups=2")
	_, facts := sweep(t, root)
	accepted := record(t, root)
	if accepted.Status != "accepted" || !strings.Contains(string(accepted.Source), "2 follow-ups") {
		t.Fatalf("approval hid follow-ups: %s\n%s", accepted.Source, fmt.Sprint(facts))
	}
}
