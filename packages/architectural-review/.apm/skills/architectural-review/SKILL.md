---
name: architectural-review
description: Workflow for writing and auditing Architecture Decision Records (ADR) and Technical Design Documents (TDD). Use when creating new architectural decisions or reviewing functional designs for correctness and logic. Not for modifying existing skills, performing code reviews, or implementing code changes.
allowed-tools: grep
metadata:
  tags: "architecture adr tdd review design documentation"
compatibility: Generic — no environment restrictions. Works in any project with markdown documentation.
---

# Architectural Review Skill

This skill enforces high-fidelity architectural and technical documentation standards. It provides a structured workflow for drafting and auditing ADRs and TDDs.

## Workflow: Drafting Documentation

When tasked with writing an ADR or TDD, follow the schemas defined in [references/schemas.md](references/schemas.md).

1. **Information Gathering:** Identify all "Decision Drivers" and "Alternatives Considered."
2. **Drafting:** Use imperative mood and scannable tables/bullets.
3. **Verification:** Ensure mandatory sections (like "Alternatives Considered" for ADRs) are present.

## Workflow: The 3-Phase Audit

When tasked with reviewing an existing document, perform this exhaustive audit:

### Phase 1: Structural Audit
- **Checklist:** Verify all required sections from [references/schemas.md](references/schemas.md) are present.
- **Grading:** Mark as `[Fail]` if the "Alternatives Considered" section is missing or if the "Interface Contract" lacks strict types.
  - **If structural audit fails:** Stop and report the missing sections. Do not proceed to Phase 2 until the document structure is corrected.

### Phase 2: Logic & Soundness Audit
- **Blast Radius:** Evaluate the impact of the decision on other system components.
- **Complexity:** Identify potential "Gold Plating" or over-engineering.
- **Security:** Check for "Fail Securely" violations and PII leakage.
- **Pure Function Test (TDD Only):** Verify if the module is truly atomic and side-effect-free.

### Phase 3: Consistency Report
- **Lineage:** Ensure the document correctly references parent ADRs or related TDDs.
- **Drift:** Identify if the proposed design contradicts existing mandates in `GEMINI.md`.

## Final Output Format

Deliver your review document-by-document. For each document, provide:
1. A `[Pass/Fail]` grade for Structure and Logic.
2. Detailed feedback and actionable fixes.
3. The Phase 3 Consistency Report.

## Verification
Verify this skill produces correct output:
1. Create a fixture ADR document with a missing "Alternatives Considered" section.
2. Run the 3-Phase Audit against it.
3. Confirm the Structural Audit returns `[Fail]` and correctly identifies the missing section.
4. Create a fixture TDD with all required sections present. Confirm it passes Phase 2 logic audit.
