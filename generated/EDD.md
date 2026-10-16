# Evaluation Service — Generated Engineering Design Document

**Evidence commit:** `36953ce98c75136b1ffcd58c1bf477a99fd4af17`  
**Method:** code-derived review using the `codebase-to-edd` evidence discipline  
**Generated:** 2026-10-17

This document describes the implementation at the evidence commit. It does not use `PRD.md` or the hand-written `EDD.md` as proof of runtime behaviour.

## 1. Executive system description

**KNOWN** — `evaluation-service` is a Go command-line repository with one executable entry point at `cmd/evalgate/main.go` and one material core package at `internal/evaluation`. The executable is a thin parse/invoke/output edge; material behaviour is implemented in the core package.

**KNOWN** — `Evaluate` indexes candidate results by case ID, applies deterministic term/refusal/citation checks, calculates weighted score, then applies threshold and regression rules.

## 2. Scope and evidence basis

Evidence inspected:

- `cmd/evalgate/main.go` — process entry and exit-code behaviour;
- `internal/evaluation/*.go` — core types and runtime behaviour;
- `internal/evaluation/*_test.go` — executable boundary evidence;
- `.github/workflows/test.yml` — continuous-integration contract;
- `examples/bafflewick/` — acceptance inputs or instructions;
- `go.mod` — dependency boundary;
- `VERSION` and `CHANGELOG.md` — release metadata.

Source snapshot contains 18 tracked files outside `generated/` at this evidence commit.

## 3. Repository map

```text
cmd/evalgate/          process entry / CLI edge
internal/evaluation/         core contract and tests
examples/bafflewick/    fictional acceptance evidence
.github/workflows/      vet, race tests, build and acceptance
architecture.txt        hand-maintained plain-text view
PRD.md / EDD.md         intended product and design documents (not evidence for this generated EDD)
```

## 4. Component model

| Runtime part | Evidence | Responsibility |
|---|---|---|
| executable | `cmd/evalgate/main.go` | parse flags/files, invoke core, encode result, map process exit status |
| core | `internal/evaluation` | own the reusable component behaviour |
| acceptance | `examples/bafflewick`, package tests | demonstrate a representative boundary without CRAB |

**KNOWN** — there is no import of another Jonno's Garage component. The repository can compile independently.

## 5. Runtime flows

**KNOWN** — `Evaluate` indexes candidate results by case ID, applies deterministic term/refusal/citation checks, calculates weighted score, then applies threshold and regression rules.

Representative flow:

```text
INPUT
  |
  v
cmd/evalgate
  |
  v
internal/evaluation
  |
  +-- governed failure --> structured result / process status 1 where applicable
  |
  +-- execution error --> stderr / process status 3 where applicable
  |
  v
structured success
```

## 6. State and persistence

**KNOWN** — No persistent state; suite and candidate JSON are per-run inputs.

## 7. Interfaces

**KNOWN** — CLI JSON files -> `evaluation.Evaluate`; no network or filesystem writes.

## 8. External dependencies

**KNOWN** — No external Go module dependencies: `go.mod` contains only the module path and Go version.

The implementation uses the Go standard library only at this commit.

## 9. Configuration

**KNOWN** — configuration is explicit through CLI flags, JSON input files and, only where documented by code, environment variables. The executable supports `-version`; CI tests the exact `0.1.0` version string.

## 10. Failure and recovery behaviour

**KNOWN** — deterministic denials/failures are represented separately from malformed input or execution errors in the CLI design. There is no generic retry loop or hidden failover layer in the repository.

**KNOWN** — the test suite includes: `TestGate`, `TestLeakFails`, `TestRegressionFails`, `TestDeterministicAndContractFailures`.

## 11. Security boundaries

**KNOWN** — Forbidden-term checks are explicit test rules, not a general data-loss-prevention control.

**UNKNOWN** — deployment filesystem permissions, host hardening and surrounding transport controls are not established by this repository unless implemented directly by the component.

## 12. Operational behaviour

**KNOWN** — `.github/workflows/test.yml` runs `go vet ./...`, `go test -race ./...`, builds the executable, checks the version contract, performs a `linux/arm64` cross-build and runs a Bafflewick acceptance check.

**KNOWN** — no background telemetry framework, cache, service mesh or hidden data store is present in the codebase.

## 13. Known / inferred / unknown

| Classification | Finding |
|---|---|
| KNOWN | one executable and one material core package implement the component |
| KNOWN | Go standard library only |
| KNOWN | CI exercises race tests and ARM64 build |
| KNOWN | failure behaviour is visible rather than silently retried |
| INFERRED | the package boundary is intended for reuse because the core is independent of the CLI and of CRAB |
| UNKNOWN | production throughput, latency and memory behaviour on a Raspberry Pi; V0.1 has a build proof, not a runtime benchmark |
| UNKNOWN | behaviour of deployment-owned systems not represented by a test adapter or repository code |

## 14. Open questions

No open question prevents V0.1 from demonstrating its declared one-job contract. The material unknowns are production integration and runtime-resource measurements, which belong to the application or a later component version unless they change this component's reusable job.

## 15. Architecture diagram

```text
CALLER
  |
  v
+-----------------------+
| cmd/evalgate         |
+-----------+-----------+
            |
            v
+-----------------------+
| internal/evaluation      |
+-----------+-----------+
            |
            v
       RESULT / EFFECT
```

See `architecture.txt` for the component-specific flow.

## 16. Evidence trace

The following properties are recovered from executable source/tests rather than the hand-written design:

| Property | Code evidence | Test / acceptance evidence |
|---|---|---|
| deterministic scoring | `Evaluate` is pure over suite + candidates | `TestDeterministicAndContractFailures` |
| leakage gate | forbidden terms fail the case | `TestLeakFails` |
| citation/refusal contracts | case flags are checked explicitly | `TestGate`, `TestDeterministicAndContractFailures` |
| regression gate | baseline minus score compared with allowance | `TestRegressionFails` |

This evidence trace does not establish deployment controls or business rationale that are absent from code. Those remain design/deployment concerns rather than inferred implementation facts.
