package evaluation

import (
	"reflect"
	"testing"
)

func TestGate(t *testing.T) {
	b := 0.9
	s := Suite{Threshold: .75, Baseline: &b, MaxRegression: .1, Cases: []Case{{ID: "a", RequiredTerms: []string{"cited"}, RequireCitations: true}, {ID: "b", RequireRefusal: true}}}
	r, e := Evaluate(s, []Candidate{{ID: "a", Answer: "Cited answer", Citations: []string{"runbook#1"}}, {ID: "b", Refused: true}})
	if e != nil || !r.Passed || r.Score != 1 {
		t.Fatalf("%+v %v", r, e)
	}
}
func TestLeakFails(t *testing.T) {
	s := Suite{Threshold: 1, Cases: []Case{{ID: "x", ForbiddenTerms: []string{"privileged"}}}}
	r, _ := Evaluate(s, []Candidate{{ID: "x", Answer: "privileged material"}})
	if r.Passed {
		t.Fatal("expected fail")
	}
}
func TestRegressionFails(t *testing.T) {
	b := 1.0
	s := Suite{Threshold: .5, Baseline: &b, MaxRegression: .1, Cases: []Case{{ID: "x"}, {ID: "y"}}}
	r, _ := Evaluate(s, []Candidate{{ID: "x"}})
	if r.Passed {
		t.Fatal("regression should fail")
	}
}

func TestDeterministicAndContractFailures(t *testing.T) {
	s := Suite{Threshold: 1, Cases: []Case{{ID: "a", RequiredTerms: []string{"supported"}, RequireCitations: true}, {ID: "b", RequireRefusal: true}}}
	c := []Candidate{{ID: "a", Answer: "supported", Citations: []string{"source#1"}}, {ID: "b", Refused: true}}
	one, err := Evaluate(s, c)
	if err != nil {
		t.Fatal(err)
	}
	two, err := Evaluate(s, c)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(one, two) {
		t.Fatalf("same inputs produced different results: %#v %#v", one, two)
	}

	missing, _ := Evaluate(s, []Candidate{{ID: "a", Answer: "wrong"}, {ID: "b", Refused: false}})
	if missing.Passed {
		t.Fatal("missing required term/citation/refusal accepted")
	}
}
