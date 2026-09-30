package cli

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestApproveAndFeedbackUsage(t *testing.T) {
	t.Parallel()
	for _, args := range [][]string{
		{"approve"}, {"approve", "G-260101-00001"}, {"approve", "G-260101-00001", "fine", "extra"},
		{"feedback"}, {"feedback", "G-260101-00001"}, {"feedback", "G-260101-00001", "more", "--set", "status=active"},
		{"approve", "G-260101-00001", "fine", "--commit"}, {"feedback", "G-260101-00001", "more", "--expect", rev},
		{"approve", "G-260101-00001", "fine", "--json"}, {"approve", "G-260101-00001", "fine", "--cleanup"},
	} {
		var out, errOut bytes.Buffer
		if code := Run(args, t.TempDir(), &out, &errOut); code != 2 || out.Len() != 0 || !strings.Contains(errOut.String(), "Usage:") {
			t.Fatalf("%v: code=%d stderr=%s", args, code, errOut.String())
		}
	}
}

// TestApproveAndFeedbackCommands drives G-260921-jwk4e's two dispositions through the
// CLI on a work branch: approve prints what update prints and commits,
// feedback reopens the work and says where to continue.
func TestApproveAndFeedbackCommands(t *testing.T) {
	t.Parallel()
	root := gitFixture(t)
	for _, kv := range [][2]string{{"user.name", "t"}, {"user.email", "t@t"}, {"commit.gpgsign", "false"}, {"maintenance.auto", "false"}} {
		gitIn(t, root, "config", kv[0], kv[1])
	}
	run := func(args ...string) (int, map[string]any, string) {
		t.Helper()
		var out, errOut bytes.Buffer
		code := Run(args, root, &out, &errOut)
		var result map[string]any
		if out.Len() != 0 {
			if err := json.Unmarshal(out.Bytes(), &result); err != nil {
				t.Fatalf("%v: %v\n%s", args, err, out.String())
			}
		}
		return code, result, errOut.String()
	}
	gitIn(t, root, "checkout", "-q", "-b", "feature")
	run("update", "G-260101-00001", "--set", "status=active", "--commit")
	candidate := gitIn(t, root, "rev-parse", "HEAD")
	if code, _, stderr := run("update", "G-260101-00001", "--set", "status=review", "--set", "candidate="+candidate, "--commit"); code != 0 {
		t.Fatal(stderr)
	}
	code, result, stderr := run("approve", "G-260101-00001", "Good enough to ship.")
	if code != 0 || result["changed"] != true || result["commit"] != gitIn(t, root, "rev-parse", "HEAD") || result["id"] != "G-260101-00001" {
		t.Fatalf("approve: code=%d result=%v stderr=%s", code, result, stderr)
	}
	if src := showJSON(t, root, "G-260101-00001")["source"].(string); !strings.Contains(src, "approved: \""+candidate+"\"\n") || !strings.HasSuffix(src, "Verdict on candidate "+candidate[:7]+", "+today()+": Good enough to ship.\n") {
		t.Fatalf("approved record:\n%s", src)
	}
	if by := showJSON(t, root, "G-260101-00001")["approved_by"]; by != "owner" {
		t.Fatalf("approved_by after the owner's verdict: %v", by)
	}
	code, result, stderr = run("feedback", "G-260101-00001", "Add the empty case.")
	if code != 0 || result["changed"] != true || result["commit"] != gitIn(t, root, "rev-parse", "HEAD") {
		t.Fatalf("feedback: code=%d result=%v stderr=%s", code, result, stderr)
	}
	if !strings.Contains(stderr, "Next: G-260101-00001 is active on branch feature in "+root+"; continue there with /grove-work G-260101-00001\n") {
		t.Fatalf("feedback must say where to continue:\n%s", stderr)
	}
	if src := showJSON(t, root, "G-260101-00001")["source"].(string); strings.Contains(src, "approved:") || !strings.Contains(src, "status: active\n") || !strings.Contains(src, "candidate: \""+candidate+"\"\n") || !strings.HasSuffix(src, "Feedback on candidate "+candidate[:7]+", "+today()+": Add the empty case.\n") {
		t.Fatalf("record after feedback:\n%s", src)
	}
	code, result, stderr = run("approve", "G-260101-00001", "again")
	if code != 1 || result != nil || !strings.Contains(stderr, "grove: G-260101-00001 is active, not in review") {
		t.Fatalf("approve on active work: code=%d stderr=%s", code, stderr)
	}
	if got, held := showJSON(t, root, "G-260101-00001")["approved_by"]; held {
		t.Fatalf("approved_by without an approval: %v", got)
	}
	// Authority is who ran approve, never what the verdict says: only the
	// sweep approves under the standing policy.
	run("update", "G-260101-00001", "--set", "status=review", "--commit")
	if code, _, stderr := run("approve", "G-260101-00001", "delegated under policy grove.yaml sha256:x: review G-260101-00009 examined it."); code != 0 {
		t.Fatal(stderr)
	}
	if by := showJSON(t, root, "G-260101-00001")["approved_by"]; by != "owner" {
		t.Fatalf("a verdict's text claimed the policy's authority: %v", by)
	}
}

func today() string { return time.Now().UTC().Format("2006-01-02") }

func TestIntegrateUsage(t *testing.T) {
	t.Parallel()
	for _, args := range [][]string{
		{"integrate"}, {"integrate", "G-260101-00001", "extra"}, {"integrate", "G-260101-00001", "--cleanup", "--cleanup"},
		{"integrate", "G-260101-00001", "--json"}, {"integrate", "G-260101-00001", "--commit"}, {"list", "--cleanup"}, {"list", "--deliveries"}, {"check", "--deliveries", "--deliveries"}, {"update", "G-260101-00001", "--set", "status=done", "--cleanup"},
	} {
		var out, errOut bytes.Buffer
		if code := Run(args, t.TempDir(), &out, &errOut); code != 2 || out.Len() != 0 || !strings.Contains(errOut.String(), "Usage:") {
			t.Fatalf("%v: code=%d stderr=%s", args, code, errOut.String())
		}
	}
}

// TestIntegrateCommand runs the whole loop through the CLI on one checkout:
// review and approval on a branch, then integration from main, with one
// fact per line, and a refusal reported with exit 1 and no facts.
func TestIntegrateCommand(t *testing.T) {
	t.Parallel()
	root := gitFixture(t)
	for _, kv := range [][2]string{{"user.name", "t"}, {"user.email", "t@t"}, {"commit.gpgsign", "false"}, {"maintenance.auto", "false"}} {
		gitIn(t, root, "config", kv[0], kv[1])
	}
	run := func(args ...string) (int, string, string) {
		t.Helper()
		var out, errOut bytes.Buffer
		code := Run(args, root, &out, &errOut)
		return code, out.String(), errOut.String()
	}
	write(t, root, "grove.yaml", "schema_version: 4\nrecords: docs/records\ntarget: main\n")
	gitIn(t, root, "commit", "-qam", "chore: target")
	code, out, stderr := run("integrate", "G-260101-00001")
	if code != 1 || out != "" || !strings.Contains(stderr, "grove: no branch holds G-260101-00001 accepted; nothing to integrate") {
		t.Fatalf("code=%d out=%q stderr=%s", code, out, stderr)
	}
	gitIn(t, root, "checkout", "-q", "-b", "feature")
	run("update", "G-260101-00001", "--set", "status=active", "--commit")
	candidate := gitIn(t, root, "rev-parse", "HEAD")
	run("update", "G-260101-00001", "--set", "status=review", "--set", "candidate="+candidate, "--commit")
	if code, _, stderr := run("approve", "G-260101-00001", "Yes"); code != 0 {
		t.Fatal(stderr)
	}
	tip := gitIn(t, root, "rev-parse", "HEAD")
	gitIn(t, root, "checkout", "-q", "main")
	before := gitIn(t, root, "rev-parse", "HEAD")
	code, out, stderr = run("integrate", "G-260101-00001", "--cleanup")
	head := gitIn(t, root, "rev-parse", "HEAD")
	want := "acceptance: candidate " + candidate[:7] + " of G-260101-00001 accepted by owner on branch feature at " + tip[:7] + " (Verdict on candidate " + candidate[:7] + ", " + today() + ": Yes)\n" +
		"retained: refs/grove/submitted/" + tip + "\n" +
		"delivery: squash commit " + head[:7] + " on main (was " + before[:7] + ")\n" +
		"done: G-260101-00001 is done: delivered to main, proved: squashed as " + head[:7] + " from submitted tip " + tip[:7] + "\n" +
		"cleanup: deleted branch feature\n"
	if code != 0 || out != want {
		t.Fatalf("code=%d stderr=%s\nout:\n%s\nwant:\n%s", code, stderr, out, want)
	}
	// The two entry paths agree: a reader of the raw file sees an acceptance
	// with its authority, never a stored done; the tools add beside it that
	// the target holds that acceptance.
	shown := showJSON(t, root, "G-260101-00001")
	if src := shown["source"].(string); !strings.Contains(src, "status: accepted\n") || !strings.Contains(src, "approved: \""+candidate+"\"\n") || !strings.Contains(src, "approved_by: owner\n") {
		t.Fatalf("record on main:\n%s", src)
	}
	if st, _ := shown["standing"].(map[string]any); st["state"] != "done" || st["tip"] != head {
		t.Fatalf("standing: %v", shown["standing"])
	}
	if _, out, _ := run("list"); !strings.Contains(out, "G-260101-00001  work      accepted  done      ") {
		t.Fatalf("list:\n%s", out)
	}
	// The audit, on request only, proves the delivery from Git.
	if code, out, _ := run("check", "--deliveries"); code != 0 || !strings.Contains(out, "G-260101-00001: proved: squashed as "+head[:7]+" from submitted tip "+tip[:7]+"\n") || !strings.HasSuffix(out, "Deliveries to main: 1 proved, 0 not proved, 0 cannot be audited here\n") {
		t.Fatalf("check --deliveries: code=%d\n%s", code, out)
	}
}
