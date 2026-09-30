package project

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"regexp"
	"slices"
	"strings"
	"time"
	"unicode/utf8"

	"go.yaml.in/yaml/v3"
)

// Record retains the original Markdown so inspection never reserializes it.
type Record struct {
	ID, Type, Title, Status string
	Kind, Size              string
	Priority                *int
	Created, Updated        *time.Time
	DependsOn, Members      []string
	Blocks, RelatesTo       []string
	Work                    []string // plan and review: the work they belong to
	Examined                string   // review: the Git commit it examined
	Candidate               string   // work: the commit offered for judgment; required in review and accepted
	Approved                string   // work: the candidate accepted; always equal to Candidate
	ApprovedBy              string   // work: "owner" or "policy sha256:…", the acceptance's authority
	ApprovedContext         string   // work: AcceptanceContext when accepted, which a later edit can make stale
	Formerly                string   // the ID or document path convert replaced
	Path                    string
	Source                  []byte
}

type Diagnostic struct {
	Path, Field, Message string
	Line                 int
}

func (d Diagnostic) String() string {
	path := d.Path
	if d.Line > 0 {
		path += fmt.Sprintf(":%d", d.Line)
	}
	if d.Field != "" {
		return path + ": " + d.Field + ": " + d.Message
	}
	return path + ": " + d.Message
}

type metadata struct {
	path   string
	offset int
	fields map[string]*yaml.Node
	errors []Diagnostic
}

func (m *metadata) problem(field, message string) {
	line := 0
	if n := m.fields[field]; n != nil {
		line = n.Line + m.offset
	}
	m.errors = append(m.errors, Diagnostic{Path: m.path, Field: field, Message: message, Line: line})
}

func parseMapping(path string, source []byte, offset int) *metadata {
	m := &metadata{path: path, offset: offset, fields: make(map[string]*yaml.Node)}
	if !utf8.Valid(source) {
		m.problem("", "invalid UTF-8")
		return m
	}
	decoder := yaml.NewDecoder(bytes.NewReader(source))
	var doc yaml.Node
	if err := decoder.Decode(&doc); err != nil {
		m.problem("", "expected a YAML mapping: "+err.Error())
		return m
	}
	var extra yaml.Node
	if err := decoder.Decode(&extra); err != io.EOF {
		if err != nil {
			m.problem("", "invalid YAML document: "+err.Error())
		} else {
			m.problem("", "expected exactly one YAML document")
		}
	}
	if len(doc.Content) != 1 || doc.Content[0].Kind != yaml.MappingNode || doc.Content[0].Tag != "!!map" {
		m.problem("", "expected a YAML mapping")
		return m
	}
	m.addFields(doc.Content[0].Content)
	return m
}

// addFields takes a mapping node's alternating keys and values.
func (m *metadata) addFields(nodes []*yaml.Node) {
	for i := 0; i < len(nodes); i += 2 {
		key, value := nodes[i], nodes[i+1]
		if key.Kind != yaml.ScalarNode || key.Tag != "!!str" {
			m.problem("", "mapping keys must be strings; merge keys are not supported")
			continue
		}
		if _, exists := m.fields[key.Value]; exists {
			m.problem(key.Value, "duplicate YAML key")
			continue
		}
		m.fields[key.Value] = value
	}
}

func (m *metadata) stringField(key string, required bool) string {
	n, ok := m.fields[key]
	if !ok {
		if required {
			m.problem(key, "required field is missing")
		}
		return ""
	}
	if n.Kind == yaml.AliasNode {
		m.problem(key, "YAML aliases are not supported")
		return ""
	}
	if n.Kind != yaml.ScalarNode || n.Tag != "!!str" || strings.TrimSpace(n.Value) == "" {
		m.problem(key, "expected a nonempty string")
		return ""
	}
	return n.Value
}

func (m *metadata) integerField(key string, required bool) (int, bool) {
	n, ok := m.fields[key]
	if !ok {
		if required {
			m.problem(key, "required field is missing")
		}
		return 0, false
	}
	var value int
	if n.Kind != yaml.ScalarNode || n.Tag != "!!int" || n.Decode(&value) != nil {
		m.problem(key, "expected an integer")
		return 0, false
	}
	return value, true
}

func (m *metadata) listField(key string) []string {
	n, ok := m.fields[key]
	if !ok {
		return nil
	}
	if n.Kind != yaml.SequenceNode || n.Tag != "!!seq" {
		m.problem(key, "expected a list of record IDs")
		return nil
	}
	result := make([]string, 0, len(n.Content))
	seen := map[string]bool{}
	for _, item := range n.Content {
		if item.Kind != yaml.ScalarNode || item.Tag != "!!str" || strings.TrimSpace(item.Value) == "" {
			m.problem(key, "each target must be a nonempty string; aliases are not supported")
			continue
		}
		if seen[item.Value] {
			m.problem(key, "duplicate target "+item.Value)
		}
		seen[item.Value] = true
		result = append(result, item.Value)
	}
	return result
}

// TypeInfo is one record type. Statuses[0] is what new writes; a type without
// statuses has no lifecycle. Fields are the type's own, beyond Envelope.
type TypeInfo struct {
	Name     string
	Statuses []string
	Fields   []string
}

// Types is the whole record vocabulary, in ID display order. These tables,
// Envelope, Kinds and Sizes are what validation checks, what its messages
// list, and what grove guide model must print (model_test.go).
// Work's done is schema 3's completion claim, which migration keeps and
// nothing newly writes: schema 4 derives Done from acceptance and delivery.
var Types = []TypeInfo{
	{"work", []string{"proposed", "active", "review", "accepted", "abandoned", "done"}, []string{"kind", "size", "priority", "members", "depends_on", "candidate", "approved", "approved_by", "approved_context"}},
	{"question", []string{"open", "resolved"}, []string{"blocks"}},
	{"decision", []string{"proposed", "accepted", "rejected", "superseded"}, nil},
	{"term", []string{"proposed", "settled"}, nil},
	{"plan", []string{"current", "superseded"}, []string{"work"}},
	{"review", []string{"current", "superseded"}, []string{"work", "examined"}},
	{"page", nil, nil},
}

// schema3Work is work as schema 3 had it, which only migration reads.
var schema3Work = TypeInfo{"work", []string{"proposed", "active", "review", "done", "abandoned"}, []string{"kind", "size", "priority", "members", "depends_on", "candidate", "approved"}}

// Envelope is the fields every record may carry; a type without statuses
// carries no status.
var Envelope = []string{"id", "type", "title", "status", "relates_to", "created", "updated", "formerly"}

// Kinds and Sizes are the values of work's kind and size.
var Kinds = []string{"feature", "fix", "refactor", "investigation", "tooling", "release"}
var Sizes = []string{"small", "medium", "large"}

// Keys is every frontmatter field a record of this type may carry.
func (t TypeInfo) Keys() []string {
	keys := slices.Clone(Envelope)
	if len(t.Statuses) == 0 {
		keys = slices.DeleteFunc(keys, func(key string) bool { return key == "status" })
	}
	return append(keys, t.Fields...)
}

// TypeNames is every type's name, in Types order.
func TypeNames() []string {
	names := make([]string, len(Types))
	for i, t := range Types {
		names[i] = t.Name
	}
	return names
}

// Choices renders values for a message: "a, b or c".
func Choices(values []string) string {
	if len(values) < 2 {
		return strings.Join(values, "")
	}
	return strings.Join(values[:len(values)-1], ", ") + " or " + values[len(values)-1]
}

// NeutralPrefix starts every ID Grove issues, whatever the type.
const NeutralPrefix = "G"

// Type returns the named type, or nil.
func Type(name string) *TypeInfo {
	for i := range Types {
		if Types[i].Name == name {
			return &Types[i]
		}
	}
	return nil
}

var datePattern = regexp.MustCompile("^[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}Z$")

// IDForm is the shape of a record ID, unanchored: the UTC creation date
// YYMMDD and five lowercase Crockford base32 characters (G-260926-2da4n).
const IDForm = NeutralPrefix + "-[0-9]{6}-[0-9a-hjkmnp-tv-z]{5}"

// IDPattern matches a record ID. An ID is an identity only, so a reclassified
// or converted record keeps its own whatever its type.
var IDPattern = regexp.MustCompile("^" + IDForm + "$")
var commitPattern = regexp.MustCompile("^[0-9a-f]{7,40}$")
var digestPattern = regexp.MustCompile("^sha256:[0-9a-f]{64}$")

// Schema is the record schema this binary reads and writes; Schema3 is the
// one before it, which only migration reads.
const (
	Schema  = 4
	Schema3 = 3
)

func (m *metadata) dateField(key string) *time.Time {
	if _, ok := m.fields[key]; !ok {
		return nil
	}
	value := m.stringField(key, false)
	if value == "" {
		return nil
	}
	date, err := time.Parse(time.RFC3339, value)
	if !datePattern.MatchString(value) || err != nil {
		m.problem(key, "expected a quoted UTC timestamp YYYY-MM-DDTHH:MM:SSZ")
		return nil
	}
	return &date
}

// Revision identifies exact file content: "sha256:" plus the lowercase hex
// SHA-256 of every byte, including any BOM and line endings. Timestamps are an
// authoring convention; callers detecting stale input must compare content.
func Revision(source []byte) string {
	sum := sha256.Sum256(source)
	return "sha256:" + hex.EncodeToString(sum[:])
}

// ParseRecord validates one record's source. Writers use it so the candidate
// they produce is judged by the same rules the reader applies.
func ParseRecord(path string, source []byte) (*Record, []Diagnostic) {
	return parseRecord(path, source, Schema)
}

func parseRecord(path string, source []byte, schema int) (*Record, []Diagnostic) {
	r := &Record{Path: path, Source: source}
	fail := func(message string) (*Record, []Diagnostic) {
		return r, []Diagnostic{{Path: path, Field: "frontmatter", Message: message}}
	}
	if !utf8.Valid(source) {
		return fail("invalid UTF-8")
	}
	text := strings.TrimPrefix(string(source), "\ufeff")
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	if len(lines) == 0 || strings.TrimRight(lines[0], " \t") != "---" {
		return fail("expected an opening --- line")
	}
	end := 1
	for end < len(lines) && strings.TrimRight(lines[end], " \t") != "---" {
		end++
	}
	if end == len(lines) {
		return fail("missing closing --- line")
	}
	m := parseMapping(path, []byte(strings.Join(lines[1:end], "\n")), 1)
	r.ID = m.stringField("id", true)
	r.Type = m.stringField("type", true)
	r.Title = m.stringField("title", true)
	t := Type(r.Type)
	if schema == Schema3 && t != nil && t.Name == "work" {
		t = &schema3Work
	}
	allowed := Envelope
	if t != nil {
		allowed = t.Keys()
	}
	// The type field alone classifies a record, so an unknown or missing one
	// is an error rather than a generic page.
	if r.Type != "" && t == nil {
		m.problem("type", "unknown record type; expected "+Choices(TypeNames()))
	}
	if !IDPattern.MatchString(r.ID) {
		m.problem("id", "expected a canonical ID, e.g. G-260925-7k2qm")
	}
	r.Formerly = m.stringField("formerly", false)
	if t == nil || len(t.Statuses) != 0 { // a page has no lifecycle, so status is an unknown field on it
		if r.Status = m.stringField("status", true); t != nil && r.Status != "" && !slices.Contains(t.Statuses, r.Status) {
			m.problem("status", "expected "+Choices(t.Statuses)+" for "+r.Type)
		}
	}
	if r.Type == "work" {
		r.Kind, r.Size = m.stringField("kind", false), m.stringField("size", false)
		if r.Kind != "" && !slices.Contains(Kinds, r.Kind) {
			m.problem("kind", "expected "+Choices(Kinds))
		}
		if r.Size != "" && !slices.Contains(Sizes, r.Size) {
			m.problem("size", "expected "+Choices(Sizes))
		}
		if priority, ok := m.integerField("priority", false); ok {
			r.Priority = &priority
			if priority < 1 || priority > 5 {
				m.problem("priority", "expected 1 (highest) through 5 (lowest)")
			}
		}
		r.Members, r.DependsOn = m.listField("members"), m.listField("depends_on")
		// A candidate is one commit, so what a review examined and what was
		// delivered can be compared to it. Its absence on a done record means
		// the record predates the review lifecycle and claims only branch-local
		// completion; the reader never rewrites that.
		if r.Candidate = m.stringField("candidate", false); r.Candidate != "" && !commitPattern.MatchString(r.Candidate) {
			m.problem("candidate", "expected a quoted Git commit of 7 to 40 lowercase hex digits")
		} else if r.Candidate == "" && (r.Status == "review" || r.Status == "accepted") {
			m.problem("candidate", "required while status is "+r.Status+": set candidate=COMMIT, the commit offered for judgment, in the same update")
		}
		// Approval is of one commit (G-260921-btyck): the field must name the
		// candidate, so a changed candidate cannot inherit it. In schema 4 it
		// is the acceptance itself, with its authority and the context it
		// accepted (G-260930-2qa4a); schema 3 kept it beside review.
		approvedIn := []string{"accepted", "done"}
		if schema == Schema3 {
			approvedIn = []string{"review", "done"}
		}
		if r.Approved = m.stringField("approved", false); r.Approved != "" {
			switch {
			case !commitPattern.MatchString(r.Approved):
				m.problem("approved", "expected a quoted Git commit of 7 to 40 lowercase hex digits")
			case r.Approved != r.Candidate:
				if r.Candidate == "" {
					m.problem("approved", "approval is of the candidate, and there is none: set candidate to the approved commit, or unset approved")
				} else {
					m.problem("approved", "approval is of one commit and must name the candidate "+r.Candidate+"; a changed candidate needs its own approval")
				}
			case !slices.Contains(approvedIn, r.Status) && slices.Contains(t.Statuses, r.Status): // an invalid status has its own refusal
				m.problem("approved", "approval holds only while status is "+Choices(approvedIn)+", not "+r.Status+": unset approved, or set status "+approvedIn[0])
			}
		}
		if schema == Schema {
			r.ApprovedBy, r.ApprovedContext = m.stringField("approved_by", false), m.stringField("approved_context", false)
			if r.ApprovedBy != "" && r.ApprovedBy != "owner" && !(strings.HasPrefix(r.ApprovedBy, "policy ") && digestPattern.MatchString(strings.TrimPrefix(r.ApprovedBy, "policy "))) {
				m.problem("approved_by", "expected owner, or policy sha256:HEX naming the grove.yaml revision a sweep acted under")
			}
			if r.ApprovedContext != "" && !digestPattern.MatchString(r.ApprovedContext) {
				m.problem("approved_context", "expected sha256: and 64 lowercase hex digits, the acceptance context")
			}
			if r.Status == "accepted" {
				for _, f := range []struct{ name, value string }{{"approved", r.Approved}, {"approved_by", r.ApprovedBy}, {"approved_context", r.ApprovedContext}} {
					if f.value == "" {
						m.problem(f.name, "required while status is accepted: grove approve writes the acceptance")
					}
				}
			} else if slices.Contains(t.Statuses, r.Status) {
				for _, f := range []struct{ name, value string }{{"approved_by", r.ApprovedBy}, {"approved_context", r.ApprovedContext}} {
					if f.value != "" {
						m.problem(f.name, "belongs to an acceptance and holds only while status is accepted, not "+r.Status)
					}
				}
			}
		}
	}
	if r.Type == "question" {
		r.Blocks = m.listField("blocks")
	}
	if r.Type == "plan" || r.Type == "review" {
		r.Work = m.listField("work")
	}
	if r.Type == "review" {
		if r.Examined = m.stringField("examined", false); r.Examined != "" && !commitPattern.MatchString(r.Examined) {
			m.problem("examined", "expected a quoted Git commit of 7 to 40 lowercase hex digits")
		}
	}
	for key := range m.fields {
		if !slices.Contains(allowed, key) {
			if t == nil {
				m.problem(key, "unknown field")
			} else {
				m.problem(key, "not a field of "+r.Type+" records, which take "+Choices(allowed))
			}
		}
	}
	r.RelatesTo = m.listField("relates_to")
	r.Created, r.Updated = m.dateField("created"), m.dateField("updated")
	if r.Created != nil && r.Updated != nil && r.Updated.Before(*r.Created) {
		m.problem("updated", "must not precede created")
	}
	return r, m.errors
}
