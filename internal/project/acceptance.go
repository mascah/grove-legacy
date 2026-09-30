package project

import (
	"strings"
)

// Appended are the openings of the paragraphs Grove appends to a work
// record's body: verdicts, feedback, reopenings, delivery notes and the
// schema 4 migration. They record judgment and history, not what was judged.
var Appended = []string{"Verdict on candidate ", "Feedback on candidate ", "Reopened with ", "Integrated under ", "Delivered ", "Migrated to schema 4"}

// AcceptanceContext identifies what an acceptance accepted (G-260930-2qa4a):
// the SHA-256, as Revision gives it, of the title, a newline and the body,
// its line endings normalized, less the ## Next section, which carries
// progress, and less each paragraph Grove appends. Frontmatter is outside
// it, so the acceptance fields never feed their own digest. Any other edit,
// to outcome, scope, acceptance or evidence, changes it, so an acceptance
// errs toward no longer applying rather than toward a false Done.
func AcceptanceContext(r *Record) string {
	text := strings.ReplaceAll(strings.TrimPrefix(string(r.Source), "\ufeff"), "\r\n", "\n")
	lines := strings.Split(text, "\n")
	// The body follows the closing --- line, which ParseRecord found.
	body := lines[len(lines):]
	for i := 1; i < len(lines); i++ {
		if strings.TrimRight(lines[i], " \t") == "---" {
			body = lines[i+1:]
			break
		}
	}
	var kept, paragraph []string
	inNext := false
	flush := func() {
		if len(paragraph) != 0 && !startsAny(paragraph[0], Appended) {
			kept = append(kept, paragraph...)
			kept = append(kept, "")
		}
		paragraph = nil
	}
	for _, line := range body {
		if strings.HasPrefix(line, "## ") {
			flush()
			inNext = strings.TrimSpace(line) == "## Next"
		}
		switch {
		case inNext:
		case strings.TrimSpace(line) == "":
			flush()
		default:
			paragraph = append(paragraph, strings.TrimRight(line, " \t"))
		}
	}
	flush()
	return Revision([]byte(r.Title + "\n" + strings.Join(kept, "\n")))
}

func startsAny(s string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}
