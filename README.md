# Evaluation Service

[![test](https://github.com/jonathanselby/evaluation-service/actions/workflows/test.yml/badge.svg)](https://github.com/jonathanselby/evaluation-service/actions/workflows/test.yml)

**Current version:** 0.1.0

**Measure the useful behaviour you claim. A regression is a build failure, not a slide.**

`evaluation-service` is a deliberately small Go command-line component for scoring recorded system outputs against a fixed evaluation set and failing a regression gate deterministically. It is independently usable outside CRAB.

## Why it exists

> Can evaluation remain small, deterministic and CI-friendly while still testing citations, refusals and forbidden leakage?

V0.1 answers that question with code, deterministic tests and a Bafflewick acceptance case, not architecture prose alone.

## What this does

- versioned JSON evaluation sets;
- deterministic checks for required terms, forbidden terms, refusal and citations;
- weighted case scoring;
- minimum score threshold;
- optional baseline score and maximum regression allowance;
- machine-readable result and stable exit codes;
- Go standard library only;

## What this does not do

- LLM-as-judge scoring;
- model invocation;
- test-set generation;
- statistical-significance claims;
- human-usefulness research;
- observability-platform functions;

## Install

Requires Go 1.26+.

```bash
go install github.com/jonathanselby/evaluation-service/cmd/evalgate@main
```

Check the executable:

```bash
evalgate -version
```

Expected version: `0.1.0`.

## Run it

Run the fictional Bafflewick evaluation set:

```bash
evalgate \
  -suite examples/bafflewick/suite.json \
  -results examples/bafflewick/results.json
```

The gate passes only when the known-answer case carries a citation, the unknown case refuses and the forbidden privileged phrase is absent. The same command exits `1` when the threshold or regression contract fails.

## Bafflewick acceptance

The full acceptance notes and fixtures are under [`examples/bafflewick/`](examples/bafflewick/). Continuous integration runs the example plus the material boundary tests in addition to the race-tested package suite.

## Stable exit codes

| Code | Meaning |
|---:|---|
| `0` | success |
| `1` | governed or deterministic contract failure |
| `2` | invalid CLI usage |
| `3` | malformed input, file I/O or execution failure |

## Design cross-check

This repository keeps intended design and implementation evidence separate. [`EDD.md`](EDD.md) states the intended component boundary. [`generated/EDD.md`](generated/EDD.md) is recovered from code, tests, configuration and operational files without using the hand-written PRD/EDD as runtime evidence. [`generated/COMPARISON.md`](generated/COMPARISON.md) compares the two and records mismatches rather than silently reconciling them.

## Engineering artefacts

- `PRD.md` — product boundary and testable acceptance criteria;
- `EDD.md` — hand-written design and production boundary;
- `architecture.txt` — plain-text runtime view;
- `internal/` — the one reusable Go contract;
- `cmd/` — thin executable edge;
- `examples/bafflewick/` — fictional acceptance case;
- `.github/workflows/test.yml` — vet, race tests, version contract, Bafflewick acceptance and `linux/arm64` cross-build;
- `generated/EDD.md` — code-derived design evidence pinned to a source commit;
- `generated/COMPARISON.md` — intended design versus recovered implementation;
- `SKILL.md` — compact adoption instructions;
- `RELEASE.md`, `CHANGELOG.md` and `VERSION` — release discipline.

## Repository shape

```text
evaluation-service/
├── .github/workflows/test.yml
├── PRD.md
├── EDD.md
├── README.md
├── SKILL.md
├── RELEASE.md
├── VERSION
├── CHANGELOG.md
├── LICENSE
├── architecture.txt
├── go.mod
├── cmd/evalgate/main.go
├── internal/evaluation/
├── examples/bafflewick/
└── generated/
    ├── EDD.md
    └── COMPARISON.md
```

## Design rule

> **Measure the useful behaviour you claim. A regression is a build failure, not a slide.**

## Licence

MIT. See [`LICENSE`](LICENSE).
