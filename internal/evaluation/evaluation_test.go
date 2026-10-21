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
	if err != nil { t.Fatal(err) }
	two, err := Evaluate(s, c)
	if err != nil { t.Fatal(err) }
	if !reflect.DeepEqual(one, two) { t.Fatalf("same inputs produced different results: %#v %#v", one, two) }
	missing, _ := Evaluate(s, []Candidate{{ID: "a", Answer: "wrong"}, {ID: "b", Refused: false}})
	if missing.Passed { t.Fatal("missing required term/citation/refusal accepted") }
}
func TestDuplicateIDsRejected(t *testing.T) {
	s := Suite{Cases: []Case{{ID: "a"}, {ID: "a"}}}
	if _, err := Evaluate(s, nil); err == nil { t.Fatal("duplicate case ids accepted") }
	s = Suite{Cases: []Case{{ID: "a"}}}
	if _, err := Evaluate(s, []Candidate{{ID: "a"}, {ID: "a"}}); err == nil { t.Fatal("duplicate candidate ids accepted") }
}
func TestInvalidNumericPolicyRejected(t *testing.T) {
	badBaseline := 1.1
	for _, s := range []Suite{
		{Cases: []Case{{ID: "a", Weight: -1}}},
		{MaxRegression: -0.1, Cases: []Case{{ID: "a"}}},
		{Baseline: &badBaseline, Cases: []Case{{ID: "a"}}},
	} {
		if _, err := Evaluate(s, []Candidate{{ID: "a"}}); err == nil { t.Fatalf("invalid suite accepted: %+v", s) }
	}
}
func TestZeroWeightKeepsDefaultWeight(t *testing.T) {
	r, err := Evaluate(Suite{Cases: []Case{{ID: "a", Weight: 0}}}, []Candidate{{ID: "a"}})
	if err != nil || len(r.Cases) != 1 || r.Cases[0].Weight != 1 { t.Fatalf("zero/default weight semantics changed: %+v %v", r, err) }
}
