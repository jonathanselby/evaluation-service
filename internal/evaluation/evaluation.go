package evaluation

import (
	"errors"
	"fmt"
	"strings"
)

type Case struct {
	ID               string   `json:"id"`
	Weight           float64  `json:"weight,omitempty"`
	RequiredTerms    []string `json:"required_terms,omitempty"`
	ForbiddenTerms   []string `json:"forbidden_terms,omitempty"`
	RequireRefusal   bool     `json:"require_refusal,omitempty"`
	RequireCitations bool     `json:"require_citations,omitempty"`
}
type Suite struct {
	Name          string   `json:"name"`
	Threshold     float64  `json:"threshold"`
	Baseline      *float64 `json:"baseline,omitempty"`
	MaxRegression float64  `json:"max_regression,omitempty"`
	Cases         []Case   `json:"cases"`
}
type Candidate struct {
	ID        string   `json:"id"`
	Answer    string   `json:"answer,omitempty"`
	Refused   bool     `json:"refused,omitempty"`
	Citations []string `json:"citations,omitempty"`
}
type CaseResult struct {
	ID      string   `json:"id"`
	Passed  bool     `json:"passed"`
	Reasons []string `json:"reasons,omitempty"`
	Weight  float64  `json:"weight"`
}
type Result struct {
	Score      float64      `json:"score"`
	Threshold  float64      `json:"threshold"`
	Passed     bool         `json:"passed"`
	Regression *float64     `json:"regression,omitempty"`
	Cases      []CaseResult `json:"cases"`
}

func Evaluate(s Suite, cands []Candidate) (Result, error) {
	if s.Threshold < 0 || s.Threshold > 1 {
		return Result{}, errors.New("threshold must be between 0 and 1")
	}
	if s.MaxRegression < 0 {
		return Result{}, errors.New("max_regression cannot be negative")
	}
	if s.Baseline != nil && (*s.Baseline < 0 || *s.Baseline > 1) {
		return Result{}, errors.New("baseline must be between 0 and 1")
	}

	caseIDs := map[string]bool{}
	for _, tc := range s.Cases {
		if strings.TrimSpace(tc.ID) == "" {
			return Result{}, errors.New("case id required")
		}
		if caseIDs[tc.ID] {
			return Result{}, fmt.Errorf("duplicate case id %q", tc.ID)
		}
		caseIDs[tc.ID] = true
		if tc.Weight < 0 {
			return Result{}, fmt.Errorf("case %q weight cannot be negative", tc.ID)
		}
	}

	m := map[string]Candidate{}
	for _, c := range cands {
		if strings.TrimSpace(c.ID) == "" {
			return Result{}, errors.New("candidate id required")
		}
		if _, exists := m[c.ID]; exists {
			return Result{}, fmt.Errorf("duplicate candidate id %q", c.ID)
		}
		m[c.ID] = c
	}

	var got, total float64
	var cr []CaseResult
	for _, tc := range s.Cases {
		w := tc.Weight
		if w == 0 {
			w = 1
		}
		total += w
		c, ok := m[tc.ID]
		r := CaseResult{ID: tc.ID, Passed: true, Weight: w}
		if !ok {
			r.Passed = false
			r.Reasons = append(r.Reasons, "missing candidate result")
		} else {
			low := strings.ToLower(c.Answer)
			for _, x := range tc.RequiredTerms {
				if !strings.Contains(low, strings.ToLower(x)) {
					r.Passed = false
					r.Reasons = append(r.Reasons, fmt.Sprintf("missing required term %q", x))
				}
			}
			for _, x := range tc.ForbiddenTerms {
				if strings.Contains(low, strings.ToLower(x)) {
					r.Passed = false
					r.Reasons = append(r.Reasons, fmt.Sprintf("contains forbidden term %q", x))
				}
			}
			if tc.RequireRefusal && !c.Refused {
				r.Passed = false
				r.Reasons = append(r.Reasons, "refusal required")
			}
			if tc.RequireCitations && len(c.Citations) == 0 {
				r.Passed = false
				r.Reasons = append(r.Reasons, "citation required")
			}
		}
		if r.Passed {
			got += w
		}
		cr = append(cr, r)
	}
	score := 1.0
	if total > 0 {
		score = got / total
	}
	res := Result{Score: score, Threshold: s.Threshold, Passed: score >= s.Threshold, Cases: cr}
	if s.Baseline != nil {
		reg := *s.Baseline - score
		res.Regression = &reg
		if reg > s.MaxRegression {
			res.Passed = false
		}
	}
	return res, nil
}
