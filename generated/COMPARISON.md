# Evaluation Service — Hand-Written vs Generated EDD

**Implementation evidence:** `36953ce98c75136b1ffcd58c1bf477a99fd4af17`  
**Generated:** 2026-10-17

## Purpose

This comparison tests whether the V0.1 implementation of `evaluation-service` communicates the architecture intended by the hand-written `EDD.md`. The generated EDD is derived from source, tests, configuration, examples and CI without using the PRD/hand-written EDD as runtime proof.

## Result

**Pass.** The material V0.1 runtime architecture is recoverable from code and the intended one-job boundary is preserved. Excluded capabilities have not silently appeared in the component.

## Comparison

| Design area | Hand-written intent | Recovered implementation | Assessment |
|---|---|---|---|
| one-job boundary | score recorded outputs and fail a regression gate | one CLI + deterministic evaluation package; no model invocation | MATCH |
| determinism | same inputs produce same result | no clock/random/network state in `Evaluate` | MATCH |
| required terms | required evidence strings checked case-insensitively | explicit term loop | MATCH |
| forbidden terms | forbidden leakage fails case | explicit forbidden-term loop | MATCH |
| refusal contract | cases may require refusal | `RequireRefusal` checked | MATCH |
| citation contract | cases may require citations | `RequireCitations` checked | MATCH |
| weighted score | case weights contribute to score | default weight 1; weighted numerator/denominator | MATCH |
| threshold/regression | gate fails below threshold or excess regression | both rules in returned `Passed` | MATCH |
| operational contract | vet, race, CLI example, ARM64 | all encoded in CI | MATCH |

## What code communicates well

1. **The scoring rule is inspectable and contains no hidden judge/model call.**
2. **Leakage, citation and refusal are ordinary deterministic test properties rather than AI-specific rhetoric.**
3. **Baseline regression is separated from absolute threshold.**

## What implementation cannot establish by itself

- whether a particular evaluation set represents user usefulness;
- statistical confidence or sampling adequacy;
- how candidate outputs are produced upstream.

These unknowns are not treated as implementation defects unless the PRD/EDD says V0.1 owns them. The generated EDD marks uncertainty rather than inventing deployment or product intent.

## V0.1 architectural finding

V0.1 is an evaluation gate, not an evaluation science platform; the code preserves that small boundary.
