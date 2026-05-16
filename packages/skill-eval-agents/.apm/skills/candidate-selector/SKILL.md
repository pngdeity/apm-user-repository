---
name: candidate-selector
description: Select the best skill variant from multiple candidates based on composite evaluation scores. Use as the final stage of the skill eval pipeline to pick the measurably best SKILL.md.
allowed-tools: go, bash
metadata:
  tags: "selection candidate ranking composite-score evaluation pipeline-final"
compatibility: Requires Go toolchain with compute-benchmark and select-best binaries. Consumes output from quality-evaluator, output-grader, and trigger-evaluator stages.
---

# Candidate Selector Skill

The final pipeline stage. Re-evaluates all skill candidates (original + up to 3 revisions) under identical conditions and selects the best one using a composite score that balances pass rate, trigger reliability, and token efficiency.

## Candidates

Candidates to evaluate:

| Candidate | Source |
|-----------|--------|
| Original | The source skill at `<skill-path>` |
| Minimalist | `evals/workspace/revisions/revision-A.md` |
| Additive | `evals/workspace/revisions/revision-B.md` |
| Restructured | `evals/workspace/revisions/revision-C.md` |

If a revision file does not exist (e.g., revision-synthesizer failed to produce it), skip that candidate.

## Workflow

### Step 1: Validate All Candidates

Before running evals, validate each candidate:

```bash
skills-ref validate <candidate-path>
```

Exclude any candidate that fails validation. Note the exclusion in the final report.

### Step 2: Run Quality Evaluator on Each Candidate

For each valid candidate, run the full quality evaluator workflow:

1. Use the candidate's SKILL.md as the skill source.
2. Run ALL test cases from `<skill-dir>/evals/evals.json` through both with-skill and without-skill configurations.
3. Capture outputs to `evals/workspace/candidate-evals/<candidate-name>/`.

This re-evaluates the original too — even though it already has eval results from the initial pipeline run. Re-running ensures a fair comparison; model nondeterminism means prior results may differ from current results on the same candidate.

All evaluations must use:
- The **same queries** from evals.json
- The **same CLI** targets (opencode, gemini)
- The **same model** version
- The **same workspace** setup (clean workspace for each run)

### Step 3: Run Output Grader on Each Candidate

For each candidate's eval outputs, run the output-grader to produce `grading.json` and `grading-summary.json`.

### Step 4: Compute Benchmark Scores

Run the benchmark aggregator:

```bash
./bin/compute-benchmark --workspace evals/workspace/quality-results/iteration-<N>
```

This reads all timing.json and grading.json files across candidates and produces `evals/workspace/quality-results/iteration-<N>/benchmark.json` with per-candidate metrics:

```json
{
  "candidates": {
    "original": {
      "pass_rate": 0.72,
      "trigger_rate": 0.85,
      "token_efficiency_delta": 0.0,
      "composite_score": null
    },
    "revision-A": {
      "pass_rate": 0.75,
      "trigger_rate": 0.82,
      "token_efficiency_delta": -0.12,
      "composite_score": null
    }
  }
}
```

- **pass_rate**: From grading.json summary (with-skill pass rate averaged across eval cases).
- **trigger_rate**: From trigger results (proportion of should_trigger=true queries that triggered, aggregated across CLIs).
- **token_efficiency_delta**: Improvement in token usage vs. no-skill baseline. Negative means the skill uses MORE tokens than no-skill (which is acceptable if pass_rate gain justifies it). Formula: `1.0 - (with_skill_tokens / without_skill_tokens)`. Positive means skill reduces tokens.
- **composite_score**: Computed by select-best stage.

### Step 5: Select Best Candidate

Run the selector:

```bash
./bin/select-best --workspace evals/workspace/quality-results/iteration-<N>
```

This computes the composite score for each candidate:

```
composite_score = (pass_rate * 0.5) + (trigger_rate * 0.3) + (token_efficiency_delta * 0.2)
```

The weights prioritize output quality (50%), then reliable triggering (30%), then efficiency (20%).

The selector produces `evals/workspace/selected.json`:

```json
{
  "selected": "revision-B",
  "selected_path": "evals/workspace/revisions/revision-B.md",
  "composite_score": 0.691,
  "rankings": [
    {"candidate": "revision-B", "composite_score": 0.691, "pass_rate": 0.78, "trigger_rate": 0.88, "token_efficiency_delta": -0.05},
    {"candidate": "revision-A", "composite_score": 0.682, "pass_rate": 0.75, "trigger_rate": 0.82, "token_efficiency_delta": -0.12},
    {"candidate": "revision-C", "composite_score": 0.670, "pass_rate": 0.73, "trigger_rate": 0.80, "token_efficiency_delta": 0.02},
    {"candidate": "original", "composite_score": 0.655, "pass_rate": 0.72, "trigger_rate": 0.85, "token_efficiency_delta": 0.0}
  ],
  "threshold_checks": {
    "pass_rate_minimum": 0.5,
    "trigger_rate_minimum": 0.5,
    "all_candidates_above_threshold": true
  }
}
```

### Step 6: Check Minimum Thresholds

Minimum thresholds:
- **pass_rate ≥ 0.5**
- **trigger_rate ≥ 0.5**

If NO candidate meets both thresholds, the selection FAILS. Write a failure report to `evals/workspace/selected.json`:

```json
{
  "selected": null,
  "selected_path": null,
  "composite_score": null,
  "rankings": [...],
  "threshold_checks": {
    "pass_rate_minimum": 0.5,
    "trigger_rate_minimum": 0.5,
    "all_candidates_above_threshold": false
  },
  "failure_reasons": {
    "revision-A": "pass_rate 0.42 below threshold 0.5",
    "revision-B": "trigger_rate 0.45 below threshold 0.5",
    "revision-C": "pass_rate 0.38 below threshold 0.5, trigger_rate 0.44 below threshold 0.5",
    "original": "pass_rate 0.41 below threshold 0.5"
  }
}
```

### Step 7: Write Final Report

If a candidate is selected:
1. Copy the winning SKILL.md to `evals/workspace/selected-SKILL.md`.
2. Write `selected.json` with full rankings (as shown in Step 5).
3. Output a summary: "Selected candidate `<name>` with composite score `<score>`. Pass rate: `<pass_rate>`, Trigger rate: `<trigger_rate>`."

If no candidate meets thresholds:
1. Write `selected.json` with failure reasons.
2. Do NOT write `selected-SKILL.md`.
3. Output the failure summary with specific reasons for each candidate.

### Step 8: Report Recommendations

If no candidate was selected, provide actionable recommendations:
- If all candidates fail pass_rate: "Consider restructuring the skill body to address the specific assertion failures listed in grading.json."
- If all candidates fail trigger_rate: "Consider revising the description field using the trigger-aggregator's optimized output."
- If mixed failures: "Analyze per-candidate failure reasons. Consider running additional revision-synthesizer iterations targeting the specific failure modes."

## Gotchas

- **Re-evaluate the original too**: The original's prior results are stale. Model nondeterminism means the same skill evaluated twice can produce different pass rates. Re-running ensures the comparison is apples-to-apples.
- **Identical evaluation conditions are mandatory**: Same queries, same CLI, same model, same prompts. If any variable differs between candidates, the ranking is invalid. Validate that all evals.json test cases match exactly.
- **Token efficiency delta can be negative**: A skill that produces better outputs may use MORE tokens than no-skill. This is acceptable and expected for additive revisions. The composite score balances this against pass rate gains. A negative delta of -0.3 with a pass_rate gain of +0.25 is typically a net positive.
- **Candidate file must exist and be valid**: If a revision file is missing or fails validation, skip it. Don't try to patch or fix it — that's the revision-synthesizer's job. Report the skip in the rankings as `status: "skipped"` with a reason.
- **Don't trust prior eval results**: The initial quality evaluator run was potentially on a different model version, CLI version, or had different workspace state. Always re-run for selection.
- **Thresholds are minimums, not goals**: A candidate that barely meets pass_rate=0.5 is probably not production-ready. The thresholds exist to filter out broken candidates, not to signify quality. Use the composite score for actual ranking.

## Verification

Verify this skill produces correct output:

1. Create 4 fixture candidates (original, A, B, C) with known pass rates: 0.6, 0.7, 0.55, 0.8 respectively.
2. Run candidate-selector end-to-end.
3. Confirm the selector picks revision-B (pass_rate=0.8, assuming trigger_rate and token_efficiency are equal) as the winner.
4. Confirm `selected-SKILL.md` contains the winning candidate's content.
5. Confirm `selected.json` has full rankings sorted by composite_score descending.
6. Modify all candidates to have pass_rate < 0.5. Re-run. Confirm selection FAILS and failure_reasons are written per candidate.
7. Confirm the original candidate was re-evaluated (not using cached results) by checking that timing.json has timestamps after the selector run started.
