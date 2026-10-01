package project

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const LegacyReviewClosing = "Open findings: none"

// ReviewSummary describes the reviewer's terminal report, never approval.
// The body retains the findings and evidence; the consumer binds the record
// to its candidate and applies the owner's delegation.
type ReviewSummary struct {
	Complete            bool
	Blockers, Followups int
}

var reviewSummary = regexp.MustCompile(`^Review v1: (complete|incomplete); blockers=([0-9]+); follow-ups=([0-9]+)$`)

// ReadReviewSummary accepts the final nonempty line only. Legacy nonzero
// reports have no classification and cannot be reinterpreted as follow-ups.
func ReadReviewSummary(source []byte) (ReviewSummary, error) {
	lines := strings.Split(strings.TrimSpace(string(source)), "\n")
	last := strings.TrimSpace(lines[len(lines)-1])
	if last == LegacyReviewClosing {
		for _, line := range lines[:len(lines)-1] {
			if strings.HasPrefix(strings.TrimSpace(line), "Review v") {
				return ReviewSummary{}, fmt.Errorf("a versioned review cannot fall back to a legacy closing line")
			}
		}
		return ReviewSummary{Complete: true}, nil
	}
	m := reviewSummary.FindStringSubmatch(last)
	if m == nil {
		return ReviewSummary{}, fmt.Errorf("does not end with %q or a valid Review v1 summary", LegacyReviewClosing)
	}
	blockers, bErr := strconv.Atoi(m[2])
	followups, fErr := strconv.Atoi(m[3])
	if bErr != nil || fErr != nil {
		return ReviewSummary{}, fmt.Errorf("review summary counts exceed the supported integer range")
	}
	return ReviewSummary{Complete: m[1] == "complete", Blockers: blockers, Followups: followups}, nil
}
