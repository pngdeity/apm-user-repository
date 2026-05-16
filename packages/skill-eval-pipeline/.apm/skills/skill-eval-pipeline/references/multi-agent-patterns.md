# Multi-Agent Patterns

The architectural patterns applied in the skill-eval-pipeline, derived from Google ADK multi-agent design principles.

## Applied Patterns

### 1. Coordinator/Dispatcher

**What:** A single orchestrator agent dispatches specialized sub-agents and monitors their progress via shared state. The orchestrator does not perform domain work — it only decides *when* and *in what order* to dispatch.

**Where:** The `skill-eval-pipeline` SKILL.md orchestrator dispatches 7 sub-agents (`struct-validator`, `trigger-evaluator`, `trigger-aggregator`, `quality-evaluator`, `output-grader`, `revision-synthesizer`, `candidate-selector`).

**Why this pattern:** Evaluation pipelines have sequential dependencies (can't grade before quality eval runs) but benefit from specialized agents per domain (structural validation is deterministic; grading requires LLM judgment). A single agent doing all 7 jobs would be too large and error-prone.

**Trade-offs:**
- **Pro:** Clean separation of concerns. Each sub-agent has a focused persona and limited tool set.
- **Pro:** Sub-agents can be developed, tested, and versioned independently.
- **Con:** Orchestrator is a single point of failure. If the orchestrator crashes mid-pipeline, state recovery depends on `state.json` integrity.

### 2. Sequential Pipeline

**What:** Agents execute in a fixed order where each agent's output is the next agent's input. The pipeline has explicit gating between stages.

**Where:** The overall pipeline flow: `struct → trigger → aggregate → quality → grade → revise → select`.

**Why this pattern:** Skill evaluation is inherently sequential. You can't measure quality delta without first running quality eval. You can't select a revision before generating revisions. Gating at each stage prevents wasted compute on downstream stages that would fail anyway.

**Trade-offs:**
- **Pro:** Predictable execution order. Easy to debug — each stage has one well-defined input and output.
- **Pro:** Early abort saves compute (structural failure at Stage 1 prevents all downstream work).
- **Con:** Total pipeline time is the sum of all stage times. No stage-level parallelism between different stages.

### 3. Parallel Fan-Out

**What:** A single stage dispatches multiple independent work units that run concurrently, then aggregates results when all complete.

**Where:**
- **Stage 2 (Trigger Evaluator):** Runs opencode and gemini trigger evaluations internally as parallel operations. Both CLIs are independent — opencode results don't affect gemini results.
- **Stage 6 (Revision Synthesizer):** Generates 3 revision variants with different strategies internally. Each variant is independent — Revision A doesn't affect Revision B.

**Why this pattern:** Trigger evaluation across CLIs is embarrassingly parallel — each CLI invocation is independent. Revision generation with different strategies is similarly independent. Parallelizing these stages roughly halves their wall-clock time.

**Trade-offs:**
- **Pro:** Significant wall-clock reduction for independent work units.
- **Pro:** Failure isolation — if opencode eval fails but gemini succeeds, we still get partial results.
- **Con:** Resource contention — parallel CLI invocations may share API rate limits.
- **Con:** More complex error handling — must decide whether to continue if one fan-out branch fails.

### 4. Generator-Critic

**What:** One agent (Generator) produces outputs and another agent (Critic) evaluates them. The Critic provides structured feedback that feeds back to the Generator for improvement.

**Where:** The `quality-evaluator` (Generator) produces with-skill and without-skill outputs. The `output-grader` (Critic) evaluates those outputs against assertions and produces structured grading with evidence. The grading signals then feed into the `revision-synthesizer` for improvement.

**Why this pattern:** Separating generation from evaluation prevents self-assessment bias. An agent that produces an output is poorly positioned to judge its own quality. The LLM-judge pattern (output-grader) provides objective, evidence-backed assessment.

**Trade-offs:**
- **Pro:** Eliminates self-assessment bias. The grader has no stake in the generator's output.
- **Pro:** The grader can identify assertion quality issues (always-pass, always-fail) that the generator would miss.
- **Con:** Two agents instead of one means double the agent overhead and coordination.
- **Con:** The grader's LLM judgment may have its own biases (position bias, verbosity bias).

### 5. Iterative Refinement

**What:** An initial version is evaluated, improvement signals are collected, and revised versions are generated. This creates a feedback loop: evaluate → identify weaknesses → revise → re-evaluate.

**Where:** The pipeline generates the original skill's eval signals (Stages 1-5), then synthesizes 3 improved variants (Stage 6), then re-evaluates all candidates to select the best (Stage 7).

**Why this pattern:** Single-pass skill creation rarely produces optimal results. Trigger descriptions need real-world activation data to tune. Procedure quality needs concrete assertion failures to identify gaps. The revision loop turns eval signals into actionable improvements.

**Trade-offs:**
- **Pro:** Data-driven improvement. Every revision is linked to specific eval signal evidence.
- **Pro:** The refinement is bounded (3 variants) to prevent infinite loops.
- **Con:** Re-evaluating all candidates in Stage 7 multiplies compute cost (4× quality eval).
- **Con:** Overfitting risk — revisions may optimize for the specific test suite rather than general quality.

### 6. Shared Session State

**What:** Agents communicate exclusively through a shared file tree. There is no direct messaging, no RPC, no event bus. Each agent reads its inputs from well-known file paths and writes its outputs to well-known file paths.

**Where:** All agents read from and write to `evals/workspace/` following the state protocol in [state-protocol.md](state-protocol.md).

**Why this pattern:** File-based state is the simplest cross-agent communication mechanism. It requires no infrastructure (no message broker, no database, no service discovery). It is trivially debuggable (inspect files on disk). It supports resumption after crashes (state is persistent). It works identically in local development and CI/CD.

**Trade-offs:**
- **Pro:** Zero infrastructure dependencies. Works everywhere filesystems work.
- **Pro:** Crash-resilient. State survives agent and orchestrator restarts.
- **Pro:** Debuggable. Inspect any intermediate output by reading the file.
- **Con:** No real-time notifications. An agent must poll or wait for the previous agent to complete.
- **Con:** Race conditions if two agents write to the same file simultaneously (mitigated by sequential pipeline design).
- **Con:** File format versioning — schema changes require all agents to be updated.

## Pattern Interaction Diagram

```
                     ┌────────────────────────────────────────┐
                     │          COORDINATOR/DISPATCHER         │
                     │         (skill-eval-pipeline)           │
                     └────────────────┬───────────────────────┘
                                      │
            ┌─────────────────────────┼─────────────────────────┐
            │                         │                         │
            ▼                         ▼                         ▼
   ┌─────────────────┐     ┌──────────────────┐     ┌──────────────────┐
   │ SEQUENTIAL      │     │ PARALLEL         │     │ SHARED SESSION   │
   │ PIPELINE        │     │ FAN-OUT          │     │ STATE            │
   │                 │     │                  │     │                  │
   │ Stage 1 ────────┤     │ Stage 2:         │     │ evals/workspace/ │
   │   │             │     │  opencode ∥ gemini│     │   state.json     │
   │   ▼             │     │                  │     │   trigger-results│
   │ Stage 2 ────────┤     │ Stage 6:         │     │   quality-results│
   │   │             │     │  RevA ∥ RevB ∥   │     │   revisions/     │
   │   ▼             │     │  RevC            │     │   selected.json  │
   │ Stage 3 ────────┤     │                  │     │                  │
   │   │             │     └──────────────────┘     └──────────────────┘
   │   ▼             │
   │ Stage 4 ────────┤     ┌──────────────────┐
   │   │             │     │ GENERATOR-CRITIC │
   │   ▼             │     │                  │
   │ Stage 5 ────────┤     │ quality-evaluator│
   │   │             │     │   (generates)    │
   │   ▼             │     │       │          │
   │ Stage 6 ────────┤     │       ▼          │
   │   │             │     │ output-grader    │
   │   ▼             │     │   (critiques)    │
   │ Stage 7 ────────┤     │       │          │
   │                 │     │       ▼          │
   └─────────────────┘     │ revision-        │
            │              │ synthesizer      │
            │              │   (refines)      │
            ▼              └──────────────────┘
   ┌─────────────────┐
   │ ITERATIVE       │
   │ REFINEMENT      │
   │                 │
   │ original → eval │
   │   → signals     │
   │   → 3 revisions │
   │   → re-eval     │
   │   → select best │
   └─────────────────┘
```

## Why Not Other Patterns

### Why not Peer-to-Peer (P2P)?
All agents are specialists dispatched by a single orchestrator. No agent needs to communicate directly with another agent — they communicate through state files. P2P would add complexity without benefit.

### Why not Hierarchical Teams?
The pipeline is flat — all agents report to the orchestrator. There are no "manager" agents that dispatch their own sub-agents. A hierarchical structure would make sense if, for example, the `trigger-evaluator` dispatched its own sub-agent for each CLI, but we chose internal parallel execution instead.

### Why not Blackboard?
Blackboard systems allow agents to opportunistically contribute to a shared knowledge base. Our pipeline is strictly sequential with explicit gating — agents must complete in order. A blackboard would allow agents to operate out of order, which breaks our gating model.

### Why not Voting/Mixture-of-Agents?
Voting requires multiple agents producing the same output type and a mechanism to aggregate votes. Our agents produce fundamentally different outputs (validation reports, benchmark scores, SKILL.md files). Voting doesn't apply.
