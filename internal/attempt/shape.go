package attempt

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/repo"
)

// Shape is what an attempt's tool calls went to, derived from the whole of
// events.jsonl when it is read, never stored. Tools counts every tool call,
// subagents' included, as Metrics does; Grove counts the Bash calls one of
// whose segments runs grove. What it misses is documented under Attempts in
// docs/commands.md.
type Shape struct {
	Tools   int            `json:"tool_calls"`
	Grove   int            `json:"grove_calls"`
	Guides  map[string]int `json:"guide_prints,omitempty"` // grove guide NAME runs, by NAME
	Skipped int            `json:"skipped_lines"`          // lines over MaxLine, not read: the counts are lower bounds
}

// Changed is what the attempt's commits changed, base to HEAD, split at the
// record root; paths are repository-relative. Written at finish, since the
// worktree and its branch may be gone when the attempt is read.
type Changed struct {
	RecordRoot string   `json:"record_root"`
	Records    []string `json:"records"` // under the record root
	Other      []string `json:"other"`
	// FirstOther is the committer time of the first commit, base to HEAD,
	// that touched a file outside the record root; zero when none did.
	FirstOther time.Time `json:"first_other_commit,omitzero"`
	Error      string    `json:"error,omitempty"` // Git could not say; the lists are empty
}

// recordRoot is the record folder, repository-relative, as grove.yaml at
// commit names it: the worktree may be gone, the commit is not.
func recordRoot(ctx context.Context, dir, commit, prefix string) (string, error) {
	config, err := repo.GitContext(ctx, dir, "show", commit+":"+path.Join(filepath.ToSlash(prefix), "grove.yaml"))
	if err != nil {
		return "", err
	}
	root := project.RecordRoot([]byte(config))
	if root == "" {
		return "", fmt.Errorf("grove.yaml at %s names no record folder", short(commit))
	}
	return path.Join(filepath.ToSlash(prefix), root), nil
}

// changedFiles lists the files base..head changed in the worktree's
// repository, split at the record root grove.yaml names at base.
func changedFiles(worktree, prefix, base, head string) *Changed {
	c := &Changed{Records: []string{}, Other: []string{}}
	root, err := recordRoot(context.Background(), worktree, base, prefix)
	if err == nil {
		c.RecordRoot = root
		var out string
		if out, err = repo.Git(worktree, "diff", "--name-only", "--no-renames", "-z", base, head); err == nil {
			for f := range strings.SplitSeq(strings.TrimSuffix(out, "\x00"), "\x00") {
				if f == "" {
					continue
				}
				if strings.HasPrefix(f, root+"/") {
					c.Records = append(c.Records, f)
				} else {
					c.Other = append(c.Other, f)
				}
			}
		}
		if err == nil && len(c.Other) != 0 {
			// Oldest last: --reverse applies after -n, so it cannot pick it.
			if out, err = repo.Git(worktree, "log", "--format=%cI", base+".."+head, "--", ":(top)", ":(top,exclude)"+root); err == nil {
				lines := strings.Fields(out)
				if len(lines) != 0 {
					c.FirstOther, err = time.Parse(time.RFC3339, lines[len(lines)-1])
				}
			}
		}
	}
	if err != nil {
		c.Error = err.Error()
	}
	return c
}

// ReadShape reads all of an attempt's events, bounded per line as
// ReadEvents is.
func ReadShape(events string) (Shape, error) {
	s := Shape{}
	err := eachLine(events, func(line []byte, _ int, over, _ bool) {
		if over {
			s.Skipped++
			return
		}
		var ev struct {
			Type    string `json:"type"`
			Message struct {
				Content []struct {
					Type  string `json:"type"`
					Name  string `json:"name"`
					Input struct {
						Command string `json:"command"`
					} `json:"input"`
				} `json:"content"`
			} `json:"message"`
		}
		if json.Unmarshal(line, &ev) != nil || ev.Type != "assistant" {
			return
		}
		for _, b := range ev.Message.Content {
			if b.Type != "tool_use" {
				continue
			}
			s.Tools++
			if b.Name != "Bash" {
				continue
			}
			grove := false
			for _, seg := range segments(b.Input.Command) {
				args, ok := groveArgs(seg)
				if !ok {
					continue
				}
				grove = true
				if len(args) > 1 && args[0] == "guide" && !strings.HasPrefix(args[1], "-") {
					if s.Guides == nil {
						s.Guides = map[string]int{}
					}
					s.Guides[args[1]]++
				}
			}
			if grove {
				s.Grove++
			}
		}
	})
	return s, err
}

// segments splits a shell command at ;, &, |, && and || and newlines,
// without parsing quotes: a heredoc's lines are segments too.
func segments(command string) []string {
	return strings.FieldsFunc(command, func(r rune) bool { return r == ';' || r == '&' || r == '|' || r == '\n' })
}

// groveArgs is a segment's arguments to grove when its program is grove,
// by any path, or go run of a cmd/grove package.
func groveArgs(segment string) ([]string, bool) {
	f := strings.Fields(strings.TrimLeft(strings.TrimSpace(segment), "({"))
	if len(f) == 0 {
		return nil, false
	}
	if path.Base(f[0]) == "grove" {
		return f[1:], true
	}
	if len(f) > 2 && f[0] == "go" && f[1] == "run" {
		for i, a := range f[2:] {
			if a == "cmd/grove" || strings.HasSuffix(a, "/cmd/grove") {
				return f[3+i:], true
			}
		}
	}
	return nil, false
}

// Total is what a set of attempts, such as one work's, cost together. An
// attempt counts in full toward each work its selection includes.
type Total struct {
	Attempts   int
	USD        float64
	Turns      int
	TurnsLow   bool          // an attempt recorded only its last result event's turns: Turns is a lower bound
	Duration   time.Duration // start to finish, finished attempts only
	Unfinished int           // no result.json: in none of the sums
	Unpriced   int           // finished without a result event: not in USD or turns
}

// Sum totals views.
func Sum(views []View) Total {
	t := Total{Attempts: len(views)}
	for _, v := range views {
		r := v.Result
		if r == nil {
			t.Unfinished++
			continue
		}
		t.Duration += max(r.Finished.Sub(v.Launch.Started), 0)
		if f := r.Events.Result; f != nil {
			t.USD += f.CostUSD
			if r.Events.Turns == 0 && f.Turns != 0 {
				t.Turns, t.TurnsLow = t.Turns+f.Turns, true
			} else {
				t.Turns += r.Events.Turns
			}
		} else {
			t.Unpriced++
		}
	}
	return t
}

func (t Total) String() string {
	low := ""
	if t.TurnsLow {
		low = "≥"
	}
	noun := "attempts"
	if t.Attempts == 1 {
		noun = "attempt"
	}
	s := fmt.Sprintf("Total: %d %s, $%.2f, %s%d turns, %dm", t.Attempts, noun, t.USD, low, t.Turns, int(t.Duration.Round(time.Minute).Minutes()))
	var out []string
	if t.Unfinished != 0 {
		out = append(out, fmt.Sprintf("%d unfinished in none of these", t.Unfinished))
	}
	if t.Unpriced != 0 {
		out = append(out, fmt.Sprintf("%d without a result event not in the cost or turns", t.Unpriced))
	}
	if len(out) != 0 {
		s += "; " + strings.Join(out, ", ")
	}
	return s
}
