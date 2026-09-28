// Package grove embeds the shared workflow guides, the review guide and the
// record model they cite, so a built binary carries the workflow of its own
// commit: the files stay the one editable owner, and grove guide prints them
// wherever it runs. It also owns the entrypoint templates init writes and the
// build identity that names all of them.
package grove

import (
	"bytes"
	"embed"
	"io/fs"
)

// Guides holds docs/work-execution.md, docs/work-shaping.md,
// docs/work-review.md and the record model they cite, docs/record-model.md,
// which names no record so it reads the same in any project.
//
//go:embed docs/work-execution.md docs/work-shaping.md docs/work-review.md docs/record-model.md
var Guides embed.FS

// GuideFiles maps each guide name grove guide takes to its file in Guides.
var GuideFiles = map[string]string{"work": "docs/work-execution.md", "shape": "docs/work-shaping.md", "review": "docs/work-review.md", "model": "docs/record-model.md"}

// WorkParts names the parts of the work guide after its head, in file order,
// each by the heading it starts at; a part runs to the next one's heading.
// grove guide work prints the head, everything before the first part, and
// --part NAME one part alone.
var WorkParts = []struct{ Name, Heading string }{
	{"prepare", "## 4. Prepare"},
	{"implement", "## 5. Implement through evidence"},
	{"review", "## 6. Review"},
	{"checkpoint", "## 7. Checkpoint and resume"},
	{"handoff", "## 8. Hand off into Review and return"},
	{"judge", "## Judging and integrating a candidate"},
	{"invocation", "## Invocation"},
}

// WorkGuide returns the work guide's head for "", the named part, or the whole
// file for "all"; the name was validated against WorkParts.
func WorkGuide(part string) []byte {
	source, err := fs.ReadFile(Guides, GuideFiles["work"])
	if err != nil {
		panic(err)
	}
	if part == "all" {
		return source
	}
	cuts := []int{}
	for _, p := range WorkParts {
		i := bytes.Index(source, []byte("\n"+p.Heading+"\n"))
		if i < 0 {
			panic("the work guide lacks " + p.Heading)
		}
		cuts = append(cuts, i+1)
	}
	cuts = append(cuts, len(source))
	if part == "" {
		return source[:cuts[0]]
	}
	for i, p := range WorkParts {
		if p.Name == part {
			return source[cuts[i]:cuts[i+1]]
		}
	}
	panic("unknown work guide part " + part)
}
