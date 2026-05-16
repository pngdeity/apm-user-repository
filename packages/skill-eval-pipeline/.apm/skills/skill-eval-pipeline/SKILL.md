---
name: skill-eval-pipeline
description: Orchestrate end-to-end AI agent skill evaluation — structural validation, trigger testing across opencode and gemini CLI, output quality measurement with graded assertions, and automated revision synthesis. Use when a SKILL.md needs measurable quality verification or when implementing CI/CD quality gates for agent context files.
allowed-tools: read_file, write_file, run_shell_command, glob, list_directory, task
metadata:
  tags: "eval testing quality skill-creation ci-cd pipeline"
compatibility: Requires opencode CLI, gemini CLI, go 1.24+, skills-ref. Designed for CI/CD automation via GitHub Actions.
---

# Skill Eval Pipeline

Coordinator/dispatcher orchestrator for the multi-agent skill evaluation pipeline. This agent does NOT perform evaluation work itself — it dispatches sub-agents in order, monitors shared state via the workspace file tree, and enforces gating decisions between stages.

## 1. Architecture

The pipeline uses a **Coordinator/Dispatcher** pattern: a single orchestrator agent dispatches 7 specialized sub-agents in a Sequential Pipeline with Parallel Fan-Out at key stages. Agents communicate exclusively through the shared workspace file tree — there is no direct messaging between sub-agents. See [state-protocol.md](references/state-protocol.md) for the complete workspace schema and [multi-agent-patterns.md](references/multi-agent-patterns.md) for the architectural patterns applied.

### Pipeline Diagram

```
                         ┌─────────────────────────┐
                         │     skill-eval-pipeline  │
                         │      (Orchestrator)      │
                         └────────────┬────────────┘
                                      │
         ┌────────────────────────────┼────────────────────────────┐
         │                            │                            │
         ▼                            ▼                            ▼
┌─────────────────┐        ┌─────────────────┐        ┌─────────────────────┐
│ struct-validator │        │     Stage Gate   │        │    Shared State      │
│   (Agent 1)      │───────▶│  pass → continue │───────▶│ evals/workspace/     │
│                  │        │  fail → ABORT    │        │   state.json         │
└─────────────────┘        └─────────────────┘        └──────────┬──────────┘
                                                                  │
                    ┌─────────────────────────────────────────────┘
                    │
                    ▼
     ┌──────────────────────────────┐
     │      trigger-evaluator       │
     │         (Agent 2)            │
     │  ┌──────────┬──────────┐     │
     │  │ opencode │  gemini  │     │  ◄── Parallel Fan-Out
     │  │ 3x runs  │ 3x runs  │     │
     │  └──────────┴──────────┘     │
     └──────────────┬───────────────┘
                    │
                    ▼
     ┌──────────────────────────────┐
     │    trigger-aggregator        │
     │        (Agent 3)             │
     │  60/40 split → desc optimize │
     └──────────────┬───────────────┘
                    │
         ┌─────────┼─────────┐
         │  Gate   │         │
         │  rates  │         │
         │  < 0%   │         │
         │  WARN   │         │
         └─────────┘         │
                    │        │
                    ▼        ▼
     ┌──────────────────────────────┐
     │    quality-evaluator         │
     │       (Agent 4)              │
     │  ┌────────┬────────┐         │
     │  │without │  with  │         │
     │  │ skill  │ skill  │         │
     │  └────────┴────────┘         │
     └──────────────┬───────────────┘
                    │
                    ▼
     ┌──────────────────────────────┐
     │     output-grader            │  ◄── Generator-Critic:
     │       (Agent 5)              │      quality-evaluator
     │  LLM judge → grading.json    │      generates, grader
     │  + benchmark.json            │      critiques
     └──────────────┬───────────────┘
                    │
         ┌─────────┼─────────┐
         │  Gate   │         │
         │ delta   │         │
         │ ≤ 0     │         │
         │ ABORT   │         │
         └─────────┘         │
                    │        │
                    ▼        ▼
     ┌──────────────────────────────┐
     │   revision-synthesizer       │
     │        (Agent 6)             │
     │  ┌──────┬──────┬──────┐      │  ◄── Parallel Fan-Out
     │  │ RevA │ RevB │ RevC │      │      3 variants
     │  └──────┴──────┴──────┘      │
     │  + skills-ref self-validate  │
     └──────────────┬───────────────┘
                    │
                    ▼
     ┌──────────────────────────────┐
     │   candidate-selector         │
     │       (Agent 7)              │
     │  re-eval all → composite     │
     │  score → select best         │
     │  → selected-SKILL.md         │
     └──────────────┬───────────────┘
                    │
                    ▼
              ┌──────────┐
              │  OUTPUT  │
              │  winning │
              │ SKILL.md │
              └──────────┘
```

### Pipeline Stages (Sequential)

| Stage | Agent | Depends On | Parallel? |
|-------|-------|-----------|-----------|
| 1. Structural Validation | `struct-validator` | — | No |
| 2. Trigger Evaluation | `trigger-evaluator` | Stage 1 pass | Yes (opencode ∥ gemini) |
| 3. Trigger Aggregation | `trigger-aggregator` | Stage 2 complete | No |
| 4. Quality Evaluation | `quality-evaluator` | Stage 2+3 complete | No |
| 5. Output Grading | `output-grader` | Stage 4 complete | No |
| 6. Revision Synthesis | `revision-synthesizer` | All prior stages | Yes (3 variants ∥) |
| 7. Candidate Selection | `candidate-selector` | Stage 6 complete | No |

## 2. Reading and Resuming State

All state is persisted to `evals/workspace/state.json`. The orchestrator reads this file at startup to determine:

```json
{
  "stage": "trigger",
  "status": "running",
  "skill_path": "packages/my-skill/skills/my-skill/SKILL.md",
  "revision_strategy": "balanced",
  "struct_validation": "pass",
  "trigger_eval_opencode_complete": true,
  "trigger_eval_gemini_complete": true,
  "optimized_description_available": false,
  "quality_iteration": 1,
  "grading_complete": false,
  "revision_count": 0,
  "selected_variant": null,
  "improvement_over_original": null,
  "errors": [],
  "warnings": [],
  "started_at": "2026-05-16T12:00:00Z",
  "updated_at": "2026-05-16T12:02:00Z"
}
```

To **resume** an interrupted pipeline:
1. Read `evals/workspace/state.json`.
2. If `status` is `running`, identify `stage` and dispatch from the next incomplete stage.
3. If `status` is `aborted`, report the abort reason and exit — do not resume.
4. If `status` is `complete`, report the selected variant and exit.

## 3. Gating Logic

### Stage 1 Gate: Structural Validation
- **Pass:** `struct-validation.json` has `"status": "pass"` → proceed to Stage 2.
- **Fail:** `"status": "fail"` → **ABORT**. Set `state.json` status to `aborted`. A structurally invalid skill cannot be meaningfully evaluated. Report errors to the user.

### Stage 2+3 Gate: Trigger Activation
- **Normal:** Activation rate > 0 on at least one CLI → proceed to Stage 4.
- **Warning:** Activation rate == 0 on all CLIs → record a warning in state but **CONTINUE**. The skill may still add quality value even if never auto-triggered. Flag for manual review.
- The pipeline never aborts on trigger failure alone.

### Stage 4+5 Gate: Quality Delta
- **Positive delta:** `skill_delta > 0` in `benchmark.json` → proceed to Stage 6.
- **Zero or negative delta:** `skill_delta <= 0` → **ABORT**. The skill adds no measurable quality and revision is unlikely to help. Flag as `"recommendation": "skill_adds_no_value"`.
- **Rationale:** If the skill doesn't improve outputs, no amount of description tuning or procedure refinement will create value. The skill may have fundamental design issues requiring human intervention.

### Stage 6+7 Gate: Revision Selection
- **Improvement:** At least one revision has a higher composite score than the original → select it.
- **No improvement:** All revisions score equal or lower → select the original. Record as `"recommendation": "keep_original"`.

## 4. Dispatching Sub-Agents

Sub-agents are defined in `.apm/agents/` and dispatched via the orchestrator's `task` tool:

```
dispatch: task(
  subagent_name: "struct-validator",
  description: "Run structural validation",
  prompt: "Validate the skill at <skill_path>. Read from evals/workspace/state.json for the skill_path. Write results to evals/workspace/struct-validation.json."
)
```

Each dispatch must:
1. Include the `skill_path` in the prompt (read from `state.json`).
2. Specify the expected output file path so the sub-agent knows where to write.
3. Include any stage-specific parameters (e.g., `revision_strategy` for the revision synthesizer).
4. Wait for the sub-agent to complete before dispatching the next stage (except for parallel fan-out at Stage 2 and Stage 6).

### Parallel Dispatch Pattern (Stage 2: Trigger Evaluator)

The trigger evaluator runs opencode and gemini internally as parallel fan-out. Dispatch the single `trigger-evaluator` agent — it handles both CLIs internally. This avoids the orchestrator needing to manage parallel sub-agent dispatch.

### Parallel Dispatch Pattern (Stage 6: Revision Synthesizer)

Dispatch the single `revision-synthesizer` agent — it generates all 3 variants internally. The orchestrator does not need to dispatch 3 parallel agents.

## 5. Error Handling Strategy

| Failure Mode | Action |
|-------------|--------|
| Structural validation fails | **ABORT** immediately. Report errors. No further stages run. |
| CLI not available (binary missing) | **WARN** in trigger results. Continue with available CLI. |
| All trigger activations == 0 | **WARN**. Continue to quality evaluation. Flag for manual review. |
| Quality delta <= 0 | **ABORT**. Skill adds no measurable value. Revision cannot fix this. |
| Quality evaluator timeout (>180s per run) | Record timeout, mark test case as failed, continue to next case. |
| Revision fails self-validation | Retry up to 2 times. If all retries fail, skip that variant. |
| Candidate selector finds no improvement | Select original. Pipeline completes normally with `keep_original`. |
| File not found (missing eval output) | **ABORT** current stage. Report which file is missing. |

## 6. Checklist

- [ ] **Stage 0: Init** — Create `evals/workspace/` directory. Write initial `state.json` with `skill_path`, `stage: "init"`, `status: "running"`.
- [ ] **Stage 0: Verify prerequisites** — Check `skills-ref` is available. Check at least one CLI (opencode, gemini) is on PATH. Check `evals/evals.json` exists.
- [ ] **Stage 1: Structural validation** — Dispatch `struct-validator`. Gate: pass → continue, fail → abort.
- [ ] **Stage 2: Trigger evaluation** — Dispatch `trigger-evaluator`. Runs opencode and gemini in parallel internally.
- [ ] **Stage 3: Trigger aggregation** — Dispatch `trigger-aggregator`. Gate: any activation > 0 → continue, all zero → warn but continue.
- [ ] **Stage 4: Quality evaluation** — Dispatch `quality-evaluator`. Runs all test cases with-skill and without-skill.
- [ ] **Stage 5: Output grading** — Dispatch `output-grader`. Gate: delta > 0 → continue, delta ≤ 0 → abort.
- [ ] **Stage 6: Revision synthesis** — Dispatch `revision-synthesizer`. Generates 3 variants with self-validation.
- [ ] **Stage 7: Candidate selection** — Dispatch `candidate-selector`. Re-evaluates all candidates, selects best.
- [ ] **Stage 8: Finalize** — Verify `selected-SKILL.md` exists. Write final state as `"status": "complete"`. Report summary.

## 7. Gotchas

### Agent Timeouts
- The `trigger-evaluator` has a 30-minute timeout. If the eval suite has many queries (e.g., 20 queries × 2 CLIs × 3 runs = 120 CLI invocations), this may be insufficient. Monitor progress and increase timeout if needed.
- The `quality-evaluator` has a 45-minute timeout. Long test cases can consume this quickly.

### CLI Auth Requirements
- Both `opencode` and `gemini` require valid API keys. The orchestrator should verify `$OPENCODE_API_KEY` or `$GEMINI_API_KEY` is set before dispatching sub-agents. If missing, abort with a clear message.
- CLI authentication tokens may expire during long pipeline runs. If a sub-agent reports authentication errors, re-check environment variables.

### Workspace Cleanup
- The workspace at `evals/workspace/` can grow large (multiple quality eval iterations with full output files). Before starting a new pipeline run, clean the workspace:
  ```bash
  rm -rf evals/workspace/*
  ```
- Do NOT clean workspace when resuming an interrupted run.

### Deterministic Splits
- The trigger aggregator's 60/40 train/val split must be deterministic (keyed on query ID hash). If the split changes between runs, the validation measurement is contaminated and comparison across revisions is invalid.

### Idempotency
- Once a stage completes and writes its output file, re-dispatching the same stage should detect the existing output and skip. Each sub-agent should check for existing output before running (except the candidate-selector, which must re-evaluate).

### Compute Cost
- A full pipeline run with 5 test cases, 10 trigger queries, and 4 candidates (original + 3 revisions) requires approximately 80-120 CLI invocations. At ~5 seconds per invocation, expect 7-10 minutes of compute. With API rate limiting, this may extend to 15-20 minutes.

### Go Tool Dependencies
- The `compute-benchmark` and `select-best` Go tools must be compiled and available at `./bin/compute-benchmark` and `./bin/select-best` (repo root). If not present, the candidate-selector will fail. Verify these exist before starting the pipeline.
