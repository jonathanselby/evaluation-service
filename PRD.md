# Evaluation Service — PRD

**Version:** 0.1

## Problem

AI-enabled systems can accumulate claims about quality without a repeatable gate that can fail a build. A useful first evaluation layer must be simple enough to run in continuous integration, deterministic enough to diagnose, and explicit about leakage, refusal and citation behaviour.

## Product objective

Provide a deliberately small Go command-line application that can score a fixed evaluation set against recorded system outputs and fail a regression gate deterministically.

> **Measure the useful behaviour you claim. A regression is a build failure, not a slide.**

The product test for V0.1 is:

> Can evaluation remain small, deterministic and CI-friendly while still testing citations, refusals and forbidden leakage?

## Users

Engineers who need a repeatable gate for retrieval, model or workflow outputs.

## In scope

- versioned JSON evaluation sets;
- deterministic checks for required terms, forbidden terms, refusal and citations;
- weighted case scoring;
- minimum score threshold;
- optional baseline score and maximum regression allowance;
- machine-readable result and stable exit codes;
- Go standard library only;

## Out of scope

- LLM-as-judge scoring;
- model invocation;
- test-set generation;
- statistical-significance claims;
- human-usefulness research;
- observability-platform functions;

These exclusions are deliberate. A future application need should widen this component only when the same reusable job has genuinely changed; otherwise the capability belongs in the consuming application or another component.

## Acceptance criteria

- the same inputs always produce the same score;
- missing required terms fail the relevant case;
- forbidden terms fail the relevant case;
- required refusal and citation behaviour are checked;
- a score below threshold exits as governed failure;
- a regression beyond the configured allowance exits as governed failure;
- the repository builds and tests for `linux/arm64`;

## Definition of done for V0.1

V0.1 is complete when all acceptance criteria are represented by automated tests or executable acceptance checks, CI runs `go vet ./...` and `go test -race ./...`, the repository cross-builds for `linux/arm64`, the production boundary and known limitations are explicit, `generated/EDD.md` is recovered from the implementation, and `generated/COMPARISON.md` shows whether the intended one-job boundary survived implementation.
