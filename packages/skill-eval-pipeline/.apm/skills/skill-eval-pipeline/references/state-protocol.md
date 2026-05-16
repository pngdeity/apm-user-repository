# State Protocol

The shared state file tree used by all sub-agents for communication. No direct messaging — agents read input from and write output to this directory.

## Directory Structure

```
evals/workspace/
├── state.json                    # Pipeline state and progress tracking
├── struct-validation.json        # Output of struct-validator (Stage 1)
├── trigger-results/
│   ├── opencode.json             # Output of trigger-evaluator, opencode CLI
│   └── gemini.json               # Output of trigger-evaluator, gemini CLI
├── trigger-aggregation.json      # Output of trigger-aggregator (Stage 3)
├── quality-results/
│   └── iteration-1/
│       ├── metadata.json         # Iteration metadata
│       ├── eval-<id>/
│       │   ├── with_skill/
│       │   │   ├── outputs/
│       │   │   │   └── output.txt
│       │   │   ├── grading.json   # Written by output-grader
│       │   │   └── timing.json
│       │   └── without_skill/
│       │       ├── outputs/
│       │       │   └── output.txt
│       │       ├── grading.json   # Written by output-grader
│       │       └── timing.json
│       ├── benchmark.json         # Aggregate benchmark written by output-grader
│       └── selected.json          # Written by candidate-selector (Stage 7)
├── revisions/
│   ├── revision-A.md
│   ├── revision-B.md
│   ├── revision-C.md
│   └── rationale.md              # Why each change was made
└── selected-SKILL.md             # Winning SKILL.md copied by candidate-selector
```

## File Schemas

### `state.json`

The single source of truth for pipeline state. Every agent reads this on startup and updates it on completion.

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

| Field | Type | Description |
|-------|------|-------------|
| `stage` | string | Current pipeline stage: `init`, `struct`, `trigger`, `trigger_complete`, `trigger_aggregated`, `quality`, `quality_complete`, `grading`, `grading_complete`, `revisions`, `revisions_synthesized`, `selection`, `complete`, `aborted` |
| `status` | string | `running`, `complete`, `aborted` |
| `skill_path` | string | Absolute or repo-relative path to the SKILL.md under evaluation |
| `revision_strategy` | string | Strategy for revision-synthesizer: `trigger_optimization`, `quality_deepening`, `concision`, `balanced` |
| `struct_validation` | string | `pass`, `fail`, or `null` (not yet run) |
| `trigger_eval_opencode_complete` | boolean | Whether opencode trigger eval is done |
| `trigger_eval_gemini_complete` | boolean | Whether gemini trigger eval is done |
| `optimized_description_available` | boolean | Whether trigger-aggregator produced an optimized description |
| `quality_iteration` | integer | Current iteration number for quality evaluation |
| `grading_complete` | boolean | Whether output-grader has run |
| `revision_count` | integer | Number of revisions synthesized |
| `selected_variant` | string | Winner from candidate-selector: `original`, `revision-A`, `revision-B`, `revision-C` |
| `improvement_over_original` | float | Composite score delta over original |
| `errors` | array | Accumulated error messages across stages |
| `warnings` | array | Accumulated warning messages across stages |
| `started_at` | string | ISO 8601 timestamp of pipeline start |
| `updated_at` | string | ISO 8601 timestamp of last state update |

### `struct-validation.json`

Written by `struct-validator`. Described fully in the struct-validator agent persona. Key field: `status` must be `pass` for the pipeline to proceed.

### `trigger-results/opencode.json` and `trigger-results/gemini.json`

Written by `trigger-evaluator`. Each contains per-query activation data with 3-run statistical sampling.

### `trigger-aggregation.json`

Written by `trigger-aggregator`. Contains train/val split, baseline rates, and optionally an optimized description string.

### `quality-results/iteration-<n>/eval-<id>/timing.json`

Written by `quality-evaluator`:

```json
{
  "eval_id": "eval-1",
  "run": "with_skill",
  "wall_seconds": 4.8,
  "cpu_seconds": 4.5,
  "exit_code": 0
}
```

### `quality-results/iteration-<n>/eval-<id>/grading.json`

Written by `output-grader`. Contains per-assertion verdicts and evidence citations.

### `quality-results/iteration-<n>/benchmark.json`

Written by `output-grader` (aggregate across all eval cases). Contains the critical `skill_delta` value used for gating. After candidate-selector aggregation, a combined `benchmark.json` may also be written at the workspace root for cross-candidate comparison.

### `revisions/revision-<variant>.md`

Written by `revision-synthesizer`. Each is a complete, standalone SKILL.md file.

### `selected.json`

Written by `candidate-selector` via Go tool:

```json
{
  "selected": "revision-B",
  "selected_path": "evals/workspace/revisions/revision-B.md",
  "scores": {
    "original": {"composite": 0.72, "quality_delta": 0.15, "activation_rate": 0.70, "timing_ratio": 1.05},
    "revision-A": {"composite": 0.78, "quality_delta": 0.18, "activation_rate": 0.82, "timing_ratio": 1.08},
    "revision-B": {"composite": 0.85, "quality_delta": 0.32, "activation_rate": 0.80, "timing_ratio": 1.02},
    "revision-C": {"composite": 0.74, "quality_delta": 0.20, "activation_rate": 0.72, "timing_ratio": 0.95}
  },
  "rankings": [
    {"candidate": "revision-B", "composite_score": 0.85, "pass_rate": 0.78, "trigger_rate": 0.88, "token_efficiency_delta": -0.05},
    {"candidate": "revision-A", "composite_score": 0.78, "pass_rate": 0.75, "trigger_rate": 0.82, "token_efficiency_delta": -0.12},
    {"candidate": "revision-C", "composite_score": 0.74, "pass_rate": 0.73, "trigger_rate": 0.80, "token_efficiency_delta": 0.02},
    {"candidate": "original", "composite_score": 0.72, "pass_rate": 0.72, "trigger_rate": 0.85, "token_efficiency_delta": 0.0}
  ],
  "improvement_over_original": 0.13,
  "recommendation": "apply_revision_B"
}
```

## State Transitions

```
init → struct → trigger → trigger_aggregated → quality → grading → revisions → selection → complete
  │        │                    │                   │                   │
  └────────┼────────────────────┼───────────────────┼───────────────────┤
           ▼                    ▼                   ▼                   ▼
         aborted              aborted             aborted             aborted
```

### Resumption Rules

When the orchestrator starts and finds `state.json` with `status: "running"`:

1. Read `stage`.
2. Verify that all files expected prior to that stage exist. If any are missing, revert to the last complete stage.
3. Dispatch from the next incomplete stage.
4. If `status` is `complete`, report the result and exit.
5. If `status` is `aborted`, report the abort reason and exit.

### Idempotency Guarantee

Each sub-agent should check for existing output before running:
- If `struct-validation.json` exists → skip struct validation.
- If `trigger-results/opencode.json` exists → skip opencode trigger eval (but run gemini if missing).
- If `revisions/revision-A.md` exists → skip revision synthesis.
- If `selected-SKILL.md` exists and `status` is `complete` → pipeline already finished.

The one exception is `candidate-selector`, which must re-evaluate all candidates to ensure fair comparison — it does not reuse prior eval results.
