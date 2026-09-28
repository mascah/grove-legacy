package grove

import (
	"io/fs"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/mascah/grove/internal/project"
)

var code = regexp.MustCompile("`([^`]+)`")

// quoted is every `code` span in text, in order.
func quoted(text string) []string {
	var out []string
	for _, m := range code.FindAllStringSubmatch(text, -1) {
		out = append(out, m[1])
	}
	return out
}

// TestRecordModelIsTheContract keeps grove guide model a reference an agent
// can afford to print, and its enumerations the ones validation checks: each
// type's statuses and fields, the envelope, work's kinds and sizes, and the
// configuration keys, so the document and the code cannot drift apart.
func TestRecordModelIsTheContract(t *testing.T) {
	source, err := fs.ReadFile(Guides, GuideFiles["model"])
	if err != nil {
		t.Fatal(err)
	}
	model := string(source)
	if len(model) > 12*1024 {
		t.Errorf("the record model is %d bytes; it must stay within 12 KB", len(model))
	}
	// rows returns the table rows in section whose first cell is `name`.
	rows := func(section, name string) [][]string {
		var found [][]string
		for _, line := range strings.Split(section, "\n") {
			if cells := strings.Split(line, " | "); strings.HasPrefix(line, "| `"+name+"` |") {
				found = append(found, cells)
			}
		}
		return found
	}
	between := func(from, to string) string {
		i := strings.Index(model, from)
		j := strings.Index(model[i+1:], to)
		if i < 0 || j < 0 {
			t.Fatalf("the record model lacks %q or %q after it", from, to)
		}
		return model[i : i+1+j]
	}

	types := between("\n| Type | Statuses", "\n\n")
	for _, ty := range project.Types {
		r := rows(types, ty.Name)
		if len(r) != 1 || len(r[0]) != 3 {
			t.Errorf("the type table must have one row for %s: %q", ty.Name, r)
			continue
		}
		if got := quoted(r[0][1]); !slices.Equal(got, ty.Statuses) {
			t.Errorf("%s statuses: model %q, code %q", ty.Name, got, ty.Statuses)
		}
		if got := quoted(r[0][2]); !slices.Equal(got, ty.Fields) {
			t.Errorf("%s fields: model %q, code %q", ty.Name, got, ty.Fields)
		}
	}
	if got := quoted(between("Every record may carry", "then its type's own")); !slices.Equal(got, project.Envelope) {
		t.Errorf("envelope: model %q, code %q", got, project.Envelope)
	}

	fields := between("\n| Field | Form", "\n\n")
	all := slices.Clone(project.Envelope)
	for _, ty := range project.Types {
		for _, f := range ty.Fields {
			if !slices.Contains(all, f) {
				all = append(all, f)
			}
		}
	}
	for _, f := range all {
		n := 0
		for _, line := range strings.Split(fields, "\n") {
			if strings.HasPrefix(line, "| ") && slices.Contains(quoted(strings.Split(line, " | ")[0]), f) {
				n++
			}
		}
		if n != 1 {
			t.Errorf("the field table must describe %s once, not %d times", f, n)
		}
	}
	for name, values := range map[string][]string{"kind": project.Kinds, "size": project.Sizes} {
		if r := rows(fields, name); len(r) != 1 || !slices.Equal(quoted(r[0][1]), values) {
			t.Errorf("%s values: model %q, code %q", name, r, values)
		}
	}

	config := between("\n| Key | Value", "\n\n")
	var keys []string
	for _, line := range strings.Split(config, "\n")[3:] {
		keys = append(keys, quoted(strings.Split(line, " | ")[0])...)
	}
	if !slices.Equal(keys, project.ConfigKeys) {
		t.Errorf("configuration keys: model %q, code %q", keys, project.ConfigKeys)
	}
	for name, want := range map[string][]string{"run": project.RunKeys, "policy": project.PolicyKeys} {
		r := rows(config, name)
		if len(r) != 1 {
			t.Errorf("the configuration table must have one row for %s", name)
			continue
		}
		// The row ends with the keys, after the commands it names.
		if got := quoted(r[0][1]); !slices.Equal(got[max(0, len(got)-len(want)):], want) {
			t.Errorf("%s keys: model %q, code %q", name, got, want)
		}
	}
	policy := between("Under `policy:`", "\n\n")
	for _, key := range slices.Concat(project.PolicyKeys, project.PolicyResolveKeys, project.PolicyApproveKeys) {
		if !strings.Contains(policy, "`"+key) {
			t.Errorf("the policy paragraph does not name %s", key)
		}
	}
}
