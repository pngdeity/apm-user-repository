---
name: writing-adrs
description: Guidance for creating and maintaining high-quality Architecture Decision Records (ADR) using the Nygard/MADR hybrid format. This skill ensures that technical decisions are documented with clear rationale, empirical trade-off analysis, and long-term auditability.
---

# Skill: Writing Architecture Decision Records (ADR)

Guidance for creating and maintaining high-quality Architecture Decision Records (ADR) using the Nygard/MADR hybrid format. This skill ensures that technical decisions are documented with clear rationale, empirical trade-off analysis, and long-term auditability.

## 1. File Naming Convention
All ADR files MUST follow a strict naming pattern to ensure chronological sorting and readability:
*   **Format:** `NNNN-slugified-title.md`
*   **Index:** A four-digit, zero-padded integer representing the decision sequence (e.g., `0001`, `0042`).
*   **Slug:** A lower-case, hyphen-separated version of the decision title.

## 2. Mandatory Document Structure
Every ADR must contain the following sections in the specified order:

### Metadata Block
Located at the top of the file, immediately following the H1 title.
*   **Status:** Current state of the decision (Proposed, Accepted, Superseded, Deprecated).
*   **Date:** The date the decision was drafted or last updated.
*   **Author:** The primary author or maintainer of the record.
*   **Parent Documents:** (Optional) Links to preceding ADRs that this decision extends or modifies.

### 1. Context and Problem Statement
A narrative description of the technical challenge or friction point. It should explain *why* a decision is required now and describe the limitations of the current approach without proposing the solution yet.

### 2. Decision Drivers
A bulleted list of 3–5 high-level constraints or goals influencing the choice (e.g., Resilience, Latency, Cost, Team Velocity). These serve as the criteria against which alternatives are measured.

### 3. Architecture Decision
The core of the document. Use an assertive "We will..." statement.
*   **Sub-components:** Use sub-headers (3.1, 3.2) if the decision involves multiple parts.
*   **Rationale:** For every tool, pattern, or constraint introduced, provide a brief "Rationale" explaining why it was chosen over others.

### 4. Alternatives Considered
A list of other approaches that were evaluated. Every alternative MUST include a **Rejection Reason** to prevent the same debate from being reopened without new information.

### 5. Consequences
A binary assessment of the decision's impact.
*   **Positive:** The direct benefits, "wins," and improvements to the system or developer experience.
*   **Negative / Trade-offs:** The "tax," increased complexity, or technical debt introduced by this choice. Every decision has a cost; this section must not be empty.

## 3. Writing Style and Tone
*   **Assertive Voice:** Use "We will..." or "We standardized on..." rather than passive or non-committal language.
*   **Empirical Focus:** Base rationales on technical trade-offs, benchmarks, or architectural principles rather than personal preference.
*   **Surgical Scope:** Each ADR should focus on one logical decision. If multiple unrelated decisions are being made, split them into separate records.

## 4. Content Guidelines & Anti-Patterns
To ensure ADRs remain true architectural artifacts and not retroactive implementation logs, they MUST adhere to the following constraints:

*   **The "Implementation Bleed" Anti-Pattern:** Do not mention specific class names, line numbers, or non-core framework details. Describe the *conceptual interface* and the *pattern*. The ADR must be authored as if the code does not exist yet.
*   **The "Infinite Resources" Anti-Pattern:** Architecture is the management of physical constraints. Never assume infinite capacity (e.g., infinite context windows, RAM, or API limits). ADRs must explicitly define physical boundaries and dictate the boundary strategy (e.g., chunking, rejection, or degradation) for when those limits are breached.
*   **Explicit Hierarchy of Constraints:** Broad, absolute statements (e.g., "Schema integrity must be prioritized regardless of trade-offs") are dangerous. When multiple drivers exist, the ADR must define an explicit Hierarchy of Constraints (e.g., 1. System Stability, 2. Schema Integrity, 3. Execution Speed) to guide future developers when resolving conflicts.
*   **Interface Requirements over Concrete Classes:** When dictating a structural change, justify the "Why" of the interface boundaries. Explain how those boundaries adhere to best practices (e.g., ISP, SRP), explicitly leaving concrete implementation details out of the record.

## 5. ADR Template
```markdown
# Architecture Decision Record {NNNN}: {Title}

**Status:** {Status}  
**Date:** {Current Date}  
**Author:** {Name/Role}  
**Parent Documents:** {Links if applicable}

## 1. Context and Problem Statement
{Describe the problem, the context, and the motivation for this change.}

## 2. Decision Drivers
* {Driver 1: Description}
* {Driver 2: Description}
* {Driver 3: Description}

## 3. Architecture Decision
We will {describe the chosen solution}.

### 3.1. {Component A}
{Details of component A}
* **Rationale:** {Why this was chosen}

### 3.2. {Component B}
{Details of component B}
* **Rationale:** {Why this was chosen}

## 4. Alternatives Considered
* **{Alternative 1}:** {Short description}
    * **Rejection Reason:** {Specific reason why this was not selected}
* **{Alternative 2}:** {Short description}
    * **Rejection Reason:** {Specific reason why this was not selected}

## 5. Consequences

### Positive
* {Benefit 1}
* {Benefit 2}

### Negative / Trade-offs
* {Drawback/Tax 1}
* {Drawback/Tax 2}
```

## 6. Verification Checklist
Before finalizing an ADR, ensure:
1. [ ] The filename uses 4-digit padding (e.g., `0005-title.md`).
2. [ ] The status is clearly defined.
3. [ ] Every alternative considered has a clear "Rejection Reason."
4. [ ] There is at least one "Negative / Trade-off" listed.
5. [ ] The record is added to the central index (e.g., `README.md` in the ADR directory).
