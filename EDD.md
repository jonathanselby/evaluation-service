# Evaluation Service — Engineering Design Document

**Version:** 0.1

## 1. Purpose

Evaluation Service is a small Go command-line component implementing one reusable job.

The core design rule is:

> **Measure the useful behaviour you claim. A regression is a build failure, not a slide.**

The V0.1 engineering question is: Can a deterministic evaluation gate cover useful behaviours such as refusal, citation and leakage without model-as-judge complexity?

## 2. Non-goals

The component:

- does not invoke models;
- does not generate evaluation sets;
- does not claim statistical significance;
- does not replace human usefulness research;

These exclusions are part of the design. V0.1 is not a smaller version of a future platform; it is intended to remain independently understandable and reusable.

## 3. Component model

```text
CALLER / INPUT
      |
      v
+-------------------------+
| cmd/evalgate            |
+------------+------------+
             | typed contract
             v
+-------------------------+
| internal/evaluation     |
+------------+------------+
             |
             v
     STRUCTURED RESULT
```

| Component | Responsibility |
|---|---|
| CLI (`cmd/evalgate`) | parse configuration/input, invoke the core, print structured output and map exit codes |
| Core (`internal/evaluation`) | own the component's one reusable contract |
| `examples/bafflewick` | executable fictional inputs or a test-backed acceptance example |

## 4. Core model

The core contract is represented by real Go types rather than a second schema invented only for documentation.

```go
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
```

Fields not shown above remain defined in the package source; `generated/EDD.md` is the implementation-derived evidence record.

## 5. Core behaviour

1. index candidate outputs by case ID.
2. evaluate required and forbidden terms case-insensitively.
3. evaluate refusal and citation requirements.
4. calculate a weighted pass ratio.
5. apply the configured minimum threshold.
6. when supplied, compare the score with a baseline and fail a regression beyond the allowance.

## 6. Contract and abstraction

The core package is the stable in-process boundary. The CLI is deliberately thin and contains no duplicate policy logic. Where the component crosses an external edge, that edge is isolated behind the smallest practical function or interface so unit tests can exercise behaviour without a live service.

The component does not import any other Jonno's Garage repository. Composition happens in the consuming application; reuse does not become runtime coupling.

## 7. Credential boundary

No credentials or network calls exist.

## 8. Runtime flow

```text
INPUT
  |
  v
PARSE + VALIDATE
  |
  +-- malformed -----------------> EXECUTION ERROR
  |
  v
CORE CONTRACT
  |
  +-- governed refusal/failure --> STRUCTURED RESULT + EXIT 1
  |
  v
SUCCESS -------------------------> STRUCTURED RESULT + EXIT 0
```

There is no silent failover. A component may deliberately degrade only where that behaviour is part of its explicit contract.

## 9. Configuration and CLI contract

Primary invocation:

```text
`evalgate -suite <suite.json> -results <candidate-results.json>`
```

Every executable also supports:

```text
-version
```

Exit codes:

```text
0  success
1  governed / deterministic contract failure
2  invalid CLI usage
3  malformed input, file I/O or execution failure
```

This is a one-shot command-line component: governed contract refusal is visible to the caller, while malformed input, file or execution failure is kept distinct from that business result.

## 10. Logging and trace

Normal results are structured JSON on standard output. Diagnostics are written to standard error. The component emits no hidden telemetry.

Security-relevant or external operations expose outcome information in their explicit result/audit contract rather than relying on unstructured logs.

## 11. Failure handling

| Failure | Behaviour |
|---|---|
| invalid threshold or suite | execution/input error |
| missing candidate case | case fails visibly |
| score below threshold | governed failure, exit 1 |
| regression exceeds allowance | governed failure, exit 1 |

Failure is visible. V0.1 does not implement retries unless a retry is part of the component's one defined job.

## 12. State and persistence

No persistent runtime state. Suites and candidate results are immutable inputs to one evaluation run.

## 13. Test strategy

| ID | Acceptance behaviour | Executable evidence |
|---|---|---|
| T-01 | passing suite enforces required terms, citations and refusals | `TestGate` |
| T-02 | forbidden leakage fails the case/gate | `TestLeakFails` |
| T-03 | regression beyond the configured allowance fails | `TestRegressionFails` |
| T-04 | identical inputs are deterministic; missing required terms/citations/refusals fail | `TestDeterministicAndContractFailures` |
| T-05 | documented Bafflewick set passes through the CLI | Bafflewick CLI acceptance in CI |
| T-06 | version and `linux/arm64` build contracts hold | CI `Version contract` and `Linux arm64 build` |

CI also runs `go vet ./...` and `go test -race ./...`. V0.1 tests deterministic contract scoring; they do not claim statistical validity, human usefulness or model-quality research.

## 14. Security boundary

This is an evaluation primitive, not a security control. Forbidden-term checks can detect configured leakage examples but do not prove the absence of sensitive data in arbitrary outputs.

Host hardening, filesystem permissions, transport security, secret storage and deployment identity remain deployment responsibilities unless this component explicitly owns them.

## 15. Operational behaviour

- one small Go executable;
- Go standard library only in V0.1;
- deterministic core where the job permits it;
- structured machine-readable output;
- stable exit codes;
- no background daemon unless the protocol itself requires a foreground loop;
- no hidden database, cache or second source of truth.

The repository's minimum-hardware proof is a `linux/arm64` build. Raspberry Pi runtime acceptance belongs to the end-to-end application unless a component-specific resource claim requires separate measurement.

## 16. Known limitations

- string checks are deterministic but deliberately shallow;
- no fuzzy matching, semantic grading or LLM-as-judge;
- scores are only meaningful for the quality of the supplied case set;

These are explicit boundaries, not implied future features.

## 17. V0.1 design conclusion

V0.1 makes evaluation executable and build-breaking while keeping scoring mechanics simple enough to understand without a model.
