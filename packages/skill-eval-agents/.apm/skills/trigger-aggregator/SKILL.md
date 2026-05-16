---
name: trigger-aggregator
description: Aggregate trigger evaluation results across CLIs, compute trigger rates, and optimize skill descriptions using train/validation splits.
allowed-tools: bash
metadata:
  tags: "trigger aggregation evaluation optimization description-tuning split-validation"
compatibility: Works with trigger results JSON from the trigger-evaluator. Uses deterministic splits with a fixed seed for reproducibility.
---

# Trigger Aggregator Skill

Aggregates per-CLI trigger results, computes trigger rates, and optimizes skill descriptions using a train/validation split methodology. This agent performs LLM reasoning to revise the `description` field based on empirical trigger data.

## Workflow

### Step 1: Load Trigger Results

Read trigger results from `evals/workspace/trigger-results/<cli>.json` for all tested CLIs. Each file contains per-query trigger counts from the trigger-evaluator step.

### Step 2: Compute Trigger Rates

For each query, across all CLIs, compute the trigger rate:

```
trigger_rate = trigger_count / total_runs
```

Where `trigger_count` is the sum of successful triggers across all CLIs for that query, and `total_runs` is the total number of invocations (typically 3 per CLI × number of CLIs).

### Step 3: Train/Validation Split

Apply a 60/40 train/validation split to the queries. Shuffle queries deterministically using a fixed seed (use seed `42` or `SKILL_EVAL_SEED` from environment):

1. Compute a hash-based sort key for each query `id` using `hash(id + seed)`.
2. Sort by this hash to produce a stable but shuffled ordering.
3. Assign the first 60% to the train set, the remaining 40% to the validation set.

This ensures the same queries land in the same splits across runs.

### Step 4: Identify Failures on Train Set

A query is a "failure" if:
- `should_trigger: true` but trigger rate < 0.5 (skill not activating when it should)
- `should_trigger: false` but trigger rate > 0.5 (skill activating when it shouldn't)

Categorize failures as:
- **Under-triggering**: should_trigger=true but rate < 0.5
- **Over-triggering**: should_trigger=false but rate > 0.5

### Step 5: Revise Description (on Train Failures)

If failures exist on the train set, revise the skill's `description` field in the YAML frontmatter. This agent modifies **only** the `description` field — not the full SKILL.md body.

Revision principles (from agentskills.io best practices):

- **Generalize from failures**: Don't copy failed query keywords directly into the description. Instead, identify the *category* of task that failed and broaden the description to cover it.
- **Don't overfit**: If "write a React component" under-triggers, adding the literal phrase "React component" will fix that one query but break others. Instead, generalize to "frontend component development" or "UI implementation patterns."
- **Cover both directions**: If the skill over-triggers on certain query categories, add exclusion clauses (e.g., "Not for simple code formatting or linting tasks").
- **Stay under 1024 characters**: The `description` field has a hard cap per the agentskills.io spec.

### Step 6: Iterate on Train Set

Re-run trigger evaluation on the train set only with the revised description. Continue iterating (up to 5 iterations) until:
- All train set queries pass (no failures), OR
- 5 iterations are exhausted with no further improvement

If stuck after 5 iterations, try a structurally different description framing — for example, shift from task-oriented ("Use when writing ADRs") to capability-oriented ("Use for structured technical decision documentation") or vice versa.

### Step 7: Evaluate on Validation Set

Once the train set passes (or iterations are exhausted), evaluate the final revised description on the held-out validation set. Do NOT further revise based on validation results — this is the held-out test of generalization.

### Step 8: Write Output

Write `evals/workspace/trigger-aggregation.json`:

```json
{
  "train": {
    "queries": ["q1", "q3", "q5"],
    "failures_initial": 3,
    "failures_final": 0,
    "iterations": 2
  },
  "validation": {
    "queries": ["q2", "q4"],
    "failures": 0,
    "validation_pass_rate": 1.0
  },
  "original_description": "original description text",
  "optimized_description": "revised description text",
  "rates": {
    "q1": {"cli_opencode": 1.0, "cli_gemini": 0.66},
    "q2": {"cli_opencode": 0.0, "cli_gemini": 0.0}
  }
}
```

## Gotchas

- **Avoid overfitting**: Copying failed query keywords verbatim into the description is the most common mistake. The validation set exists precisely to catch this. If train passes but validation fails, the description is overfit.
- **Stuck after 5 iterations**: If the failure rate plateaus after 5 iterations, the issue is likely in the skill body (instructions), not the description. Flag this for the revision-synthesizer stage.
- **Seed determinism is critical**: Must use a fixed seed for the train/validation split. If the seed changes between runs, results become non-reproducible and the validation metric is meaningless.
- **Description length tradeoffs**: Adding exclusion clauses ("Not for...") consumes characters that could be used for inclusion clauses. Prioritize what the skill DOES, not just what it avoids.
- **Cross-CLI variability**: A query may trigger reliably on opencode but not gemini. If the variance between CLIs is high (>0.5 difference in rates), investigate whether the CLI's routing algorithm differs rather than blaming the description.
- **Edge case: zero total_runs**: If all CLI invocations for a query timed out, the trigger rate is undefined. Exclude these queries from train/validation splits and note them in output.

## Verification

Verify this skill produces correct output:

1. Create synthetic trigger results with known failure patterns: 3 queries in train set with under-triggering, 2 queries in validation set fully passing.
2. Run the trigger-aggregator. Confirm it detects the 3 failures, produces a revised description, and iterates until train passes.
3. Confirm the validation set results are computed but do NOT influence further description changes.
4. Verify the fixed seed produces identical train/validation splits across 3 consecutive runs with the same input data.
5. Confirm that a description exceeding 1024 characters after revision is rejected and the agent retries with a shorter description.
