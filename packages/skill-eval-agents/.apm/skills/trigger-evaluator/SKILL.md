---
name: trigger-evaluator
description: Evaluate whether a skill's description triggers correctly on test queries across agent CLIs. Use when testing skill activation accuracy against opencode and gemini CLI targets.
allowed-tools: go, bash
metadata:
  tags: "trigger evaluation activation accuracy cli opencode gemini skill-routing"
compatibility: Requires Go toolchain and the invoke-cli binary. Targets opencode and gemini CLIs for skill activation testing.
---

# Trigger Evaluator Skill

Evaluates how reliably a skill's `description` field causes agent CLIs to activate the skill for relevant queries. Runs each test query multiple times to account for model nondeterminism and records per-run trigger state.

## Workflow

### Step 1: Load Test Queries

Read `<skill-dir>/evals/evals.json`. This file contains an array of test queries with the schema:

```json
[
  {"id": "q1", "query": "write an ADR for auth", "should_trigger": true},
  {"id": "q2", "query": "add a comment to the login function", "should_trigger": false}
]
```

`should_trigger: true` means the skill description SHOULD cause the CLI to activate this skill. `should_trigger: false` means the skill should NOT activate.

### Step 2: Run Trigger Tests

For each query, invoke the target CLI 3 times per CLI target (opencode and gemini):

```bash
go run ./cmd/invoke-cli \
  --cli <opencode|gemini> \
  --prompt "<query>" \
  --workspace <workspace> \
  --skill <skill-path>
```

Use a **120-second timeout** per invocation. The `invoke-cli` tool outputs JSON with a `skill_activated` boolean field.

### Step 3: Record Results

For each CLI target, produce a result file at `evals/workspace/trigger-results/<cli>.json`:

```json
{
  "cli": "opencode",
  "results": [
    {
      "id": "q1",
      "query": "write an ADR for auth",
      "should_trigger": true,
      "runs": [true, false, true],
      "trigger_count": 2,
      "total_runs": 3
    },
    {
      "id": "q2",
      "query": "add a comment to the login function",
      "should_trigger": false,
      "runs": [false, false, false],
      "trigger_count": 0,
      "total_runs": 3
    }
  ]
}
```

### Step 4: Handle Failures

If a CLI invocation times out (exceeds 120s), record that run as `null` and reduce `total_runs` accordingly. If all 3 runs for a query time out, flag the query as `status: "unavailable"` and exclude it from downstream aggregation with a note in the output.

If `invoke-cli` crashes or returns malformed JSON, retry once. If the retry also fails, mark all remaining queries for that CLI as `status: "cli_error"` and write a partial results file.

## Gotchas

- **opencode may not invoke skills for single-step tasks**: opencode's router can skip skill activation for queries it considers trivial or single-action. A query like "list files" (should_trigger: false for most skills) may silently bypass the skill system — this is expected behavior, not a test failure.
- **gemini stream-json mode buffers differently**: The gemini CLI emits skill activation events in stream-json mode, which can buffer lines. Ensure `invoke-cli` reads the full stream before parsing, or activation events may be missed.
- **Model nondeterminism requires multiple runs**: A single run can produce a false positive or false negative due to model variability. Three runs is the minimum; for high-stakes evaluations, consider 5 runs and use majority vote.
- **Prompt length affects activation**: Very long prompts (over 2000 tokens) can dilute the skill description's influence on routing. Keep test queries concise and focused.
- **Description caching**: Some CLI implementations cache skill descriptions at startup. If you modify the description between test runs, restart the CLI or clear its cache.
- **Skill directory must be discoverable**: The `--skill` path must be registered or resolvable by the CLI's skill discovery mechanism. If the skill is not found, `skill_activated` will always be false regardless of query relevance.

## Verification

Verify this skill produces correct output:

1. Create a fixture skill with a description that mentions "kubernetes deployment".
2. Create an `evals.json` with two queries: one about "deploy to k8s" (should_trigger: true) and one about "sort an array" (should_trigger: false).
3. Run the trigger evaluator against both CLI targets.
4. Confirm `evals/workspace/trigger-results/opencode.json` and `gemini.json` exist.
5. Confirm that `trigger_count` is a number between 0 and 3 (or fewer if timeouts occurred) for each query.
6. Confirm that timed-out runs produce `null` entries in the `runs` array, not false.
