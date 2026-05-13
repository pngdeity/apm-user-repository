---
name: writing-design-files
description: Guidance for creating atomic Technical Design Documents (TDDs) using the "Micro-Leaf" philosophy. This skill ensures that complex systems are decomposed into verifiable, isolated units of logic ("Leaves") before implementation begins.
---

# Skill: Writing Technical Design Documents (Micro-Leaves)

Guidance for creating atomic Technical Design Documents (TDDs) using the "Micro-Leaf" philosophy. This skill ensures that complex systems are decomposed into verifiable, isolated units of logic ("Leaves") before implementation begins.

## 1. The Micro-Leaf Philosophy
*   **Blast Radius Isolation:** Atomic units ensure that failures are localized and easily diagnosable.
*   **Interface-First Design:** Define the boundaries (Inputs/Outputs) before the internal logic.
*   **Pure Function Preference:** Separate data transformation from I/O boundaries whenever possible to maximize testability.
*   **Proving the Physics:** Use design to identify where isolated "Tracer Bullet" scripts are needed to verify external API or library behavior in a vacuum.

## 2. Mandatory Document Structure
Each Design Document (or section within a document) must contain these six sections:

### 1. Metadata
*   **Leaf ID & Title:** A hierarchical identifier (e.g., `Leaf 1.1: Name`) to map the leaf to its functional branch.
*   **Parent ADRs:** Links to the Architectural Decision Records that mandate this specific implementation.
*   **Status:** The current state of the design (Draft, Reviewed, Approved, Superseded).

### 2. Objective & Non-Goals
*   **Objective:** A concise statement of what this leaf *must* accomplish.
*   **Non-Goals:** Explicitly define what this leaf *will not* do. This is critical for preventing scope creep and ensuring isolation.

### 3. Interface Contract
Define the "Black Box" boundaries:
*   **Inputs:** All data required to execute the logic (parameters, environment variables, state).
*   **Outputs/Side Effects:** What the leaf returns and what changes it makes to the external system (files written, network calls made, database updates).

### 4. Execution Logic & State Machine
A step-by-step description of the implementation flow. Use numbered lists or pseudocode to describe the internal transformations and state transitions.

### 5. Failure Modes & Resilience
Identify what can go wrong and how the system responds. Categories to check:
*   **I/O Errors:** (e.g., File system full, network timeout).
*   **Constraint Violations:** (e.g., Input data doesn't match the schema).
*   **Dependency Failures:** (e.g., External API returns 429 or 500).
*   **Security Boundaries:** (e.g., Potential path traversal or injection points).
*   **Mitigation:** For every failure, specify the action (e.g., Raise specialized Exception, Retry with Backoff, Fail Open/Closed).

### 6. Testing Strategy
Instructions for verifying the leaf in isolation:
*   **Mocking Requirements:** What external systems must be simulated.
*   **Validation:** The specific assertions required to prove success.
*   **Edge Cases:** The specific "pathological" inputs that must be tested.

## 3. Writing Style and Tone
*   **Surgical Precision:** Keep the scope of a single leaf to a single logical task.
*   **Implementation Neutral:** Focus on logic and contracts rather than specific syntax, unless a language-specific feature is the core of the design.
*   **Defensive Posture:** Write the Failure Modes section as if the leaf is expected to fail.

## 4. Technical Design Template
```markdown
### Technical Design Document: Leaf {X.Y} ({Title})

**1. Metadata**
* **Leaf ID & Title:** Leaf {X.Y}: {Title}
* **Parent ADRs:** {ADR-XXXX}
* **Status:** {Status}

**2. Objective & Non-Goals**
* **Objective:** {Concise goal}
* **Non-Goals:** {What is explicitly out of scope}

**3. Interface Contract**
* **Inputs:** 
    * `{param_1}`: {type} ({description})
* **Outputs/Side Effects:**
    * {What is returned or changed}

**4. Execution Logic & State Machine**
1. {Step 1}
2. {Step 2}
3. {Step 3}

**5. Failure Modes & Resilience**
* **{Failure Type}:** {Scenario} -> **Mitigation:** {Specific Action}
* **{Failure Type}:** {Scenario} -> **Mitigation:** {Specific Action}

**6. Testing Strategy**
* **Isolation:** {How to test without the rest of the system}
* **Verification:** {Success criteria}
```

## 5. Verification Checklist
1. [ ] Is the leaf truly atomic? (Can it be explained in one sentence?)
2. [ ] Are the Non-Goals clearly defined?
3. [ ] Does the Interface Contract specify both data and side effects?
4. [ ] Does the Failure Modes section cover at least three distinct scenarios?
5. [ ] Can this leaf be tested without running the entire application?
