# technical-documentation

Architecture Decision Record (ADR) and Technical Design Document (Micro-Leaf) authoring workflows.

## Primitives

### Skills

- **Technical Documentation** — Entry point skill that delegates to ADR and Design Leaf workflows based on task context
- **Writing ADRs** (reference) — Nygard/MADR hybrid format with template, verification checklist, and anti-patterns
- **Writing Design Files** (reference) — Micro-Leaf philosophy with interface contract templates and failure mode analysis

## Documentation Philosophy

1. ADRs first — define the Why and What at the system level
2. Leaves second — define the How for each atomic unit
3. Code third — implement only after interface and failure modes are approved
