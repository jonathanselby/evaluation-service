package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/jonathanselby/evaluation-service/internal/evaluation"
	"os"
)

const version = "0.1.0"

func main() {
	var suite, out string
	var v bool
	flag.StringVar(&suite, "suite", "", "evaluation suite JSON")
	flag.StringVar(&out, "results", "", "candidate results JSON")
	flag.BoolVar(&v, "version", false, "print version")
	flag.Parse()
	if v {
		fmt.Println("evalgate " + version)
		return
	}
	if suite == "" || out == "" {
		flag.Usage()
		os.Exit(2)
	}
	var s evaluation.Suite
	var c []evaluation.Candidate
	if e := read(suite, &s); e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(3)
	}
	if e := read(out, &c); e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(3)
	}
	r, e := evaluation.Evaluate(s, c)
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(3)
	}
	json.NewEncoder(os.Stdout).Encode(r)
	if !r.Passed {
		os.Exit(1)
	}
}
func read(p string, v any) error {
	b, e := os.ReadFile(p)
	if e != nil {
		return e
	}
	return json.Unmarshal(b, v)
}
