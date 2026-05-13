---
name: technical-documentation
description: Comprehensive guidance for writing architectural and technical documentation. Use when the user asks to create or update Architecture Decision Records (ADRs) or Technical Design Documents (Micro-Leaves).
---

# Technical Documentation (ADRs & Design Leaves)

This project maintains a strict hierarchy of technical documentation to ensure architectural integrity and implementation reliability.

## 1. Architectural Decisions (ADRs)
Use this workflow when making high-level technology choices, defining system patterns, or establishing project-wide constraints.
*   **Workflow & Template:** Refer to [references/adrs.md](./references/adrs.md)
*   **Key Focus:** Rationale, alternatives considered, and long-term consequences.

## 2. Technical Implementation (Micro-Leaves)
Use this workflow when designing atomic units of logic, defining interface contracts, or mapping out failure modes before writing code.
*   **Workflow & Template:** Refer to [references/design.md](./references/design.md)
*   **Key Focus:** Interface contracts, blast radius isolation, and failure-mode resilience.

## 3. Documentation Philosophy
1.  **ADRs first:** Define the "Why" and the "What" at the system level.
2.  **Leaves second:** Define the "How" for each atomic unit of that decision.
3.  **Code third:** Implement only after the interface and failure modes are approved.
