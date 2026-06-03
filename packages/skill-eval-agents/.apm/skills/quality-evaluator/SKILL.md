---
name: quality-evaluator
description: Evaluate skill output quality by running test cases with and without the skill and comparing results. Use within the eval pipeline to measure a skill's value-add.
allowed-tools: go, bash
metadata:
  tags: "quality evaluation comparison with-skill without-skill value-add baseline"
compatibility: Requires Go toolchain and invoke-cli. Runs opencode and gemini in controlled sessions for quality measurement.
---

# Quality Evaluator Skill

Runs test cases through agent CLIs both with and without the target skill
loaded, capturing outputs and timing data. This produces the raw material for
downstream grading and comparison.

## Workflow

### Step 1: Load Test Cases

Read quality eval test cases from `<skill-dir>/evals/evals.json`:

```json
[
  {
    "id": "eval-1",
    "prompt": "Create a chart showing monthly revenue from the CSV file",
    "expected_output": "A PNG bar chart with labeled axes",
    "files": ["data.csv"],
    "assertions": [
      "Output includes a PNG image file",
      "X-axis has month labels",
      "Y-axis has revenue labels"
    ]
  }
]
```

`files` (optional): input files to stage in the workspace before running the
prompt.

### Step 2: Prepare Workspace per Iteration

Create a clean workspace for each iteration:

```
evals/workspace/quality-results/iteration-<N>/
├── eval-<id>/
│   ├── with_skill/
│   │   ├── workspace/    # fresh workspace for this run
│   │   └── outputs/      # captured outputs
│   └── without_skill/
│       ├── workspace/
│       └── outputs/
```

Each run gets a clean workspace — no leftover files, no cached agent state, no
prior context. This is mandatory for valid comparison.

### Step 3: Stage Input Files

For test cases that specify `files`, copy the referenced files into the
workspace directory before invoking the CLI. Files are relative to
`<skill-dir>/evals/`.

### Step 4: Run With Skill

Invoke the CLI with the skill directory available in the agent's scope:

```bash
go run ./cmd/invoke-cli \
  --cli <opencode|gemini> \
  --prompt "<prompt>" \
  --workspace evals/workspace/quality-results/iteration-<N>/eval-<id>/with_skill/workspace \
  --skill <skill-path>
```

### Step 5: Run Without Skill

Invoke the CLI without providing a skill path — omit the `--skill` flag so the
agent does not discover or load the target skill:

```bash
go run ./cmd/invoke-cli \
  --cli <opencode|gemini> \
  --prompt "<prompt>" \
  --workspace evals/workspace/quality-results/iteration-<N>/eval-<id>/without_skill/workspace
```

### Step 6: Capture Outputs

For each run, save:

- All files produced in the workspace to `<run>/outputs/`
- The agent's text response to `<run>/outputs/response.txt`
- Timing data to `<run>/outputs/timing.json`:
  ```json
  {
    "token_count": { "input": 500, "output": 1200 },
    "duration_ms": 8500
  }
  ```

### Step 7: Handle Errors

If the agent produces an error, capture the error output rather than aborting:

- Save error text to `<run>/outputs/error.txt`
- Still capture any partial outputs
- Note the error in timing.json with `"error": true`

## Gotchas

- **Clean context per run is mandatory**: Reusing sessions or workspaces across
  runs will contaminate results. Always use `--new-session` or equivalent flags.
  The gemini CLI in particular may hold context across invocations unless
  explicitly reset.
- **gemini needs `--new-session` flag**: Without `--new-session`, gemini carries
  forward prior conversation context, which can make the with-skill and
  without-skill runs share state and invalidate the comparison.
- **File output capture must be exhaustive**: Some agents produce files in
  unexpected locations (temp directories, global caches). Ensure `invoke-cli`
  captures the full workspace after the run, not just the initial directory.
- **Stale workspace cleanup**: If a previous iteration left files, they may be
  picked up as inputs by the next run. Always `rm -rf` the workspace directory
  before creating it fresh.
- **Timing variance**: Token counts and durations can vary +/-20% between runs
  on the same prompt due to model nondeterminism. Treat small timing differences
  as noise; use token efficiency delta in candidate-selector, not absolute
  timing.
- **Long-running prompts**: If a test case takes over 5 minutes, consider it for
  exclusion or increase the CLI timeout. A single hung evaluation should not
  block the entire pipeline.

## Verification

Verify this skill produces correct output:

1. Create a fixture skill and an `evals.json` with one test case: "List all
   Python files in the workspace."
2. Stage a workspace with 3 `.py` files and 2 `.txt` files.
3. Run quality evaluator for iteration 1.
4. Confirm `with_skill/outputs/response.txt` and
   `without_skill/outputs/response.txt` both exist.
5. Confirm `timing.json` is present in both output directories.
6. Confirm the `without_skill` workspace does not contain any skill-loaded
   artifacts.
7. Re-run and confirm iteration-2 was created in a clean workspace (no files
   from iteration-1).
