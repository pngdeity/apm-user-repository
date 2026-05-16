---
name: output-grader
description: Grade skill evaluation outputs against assertions using LLM judgment. Use after quality evaluator completes to produce pass/fail grading with concrete evidence.
allowed-tools: bash
metadata:
  tags: "grading evaluation assertions pass-fail evidence llm-judge quality-assessment"
compatibility: Reads quality evaluator outputs. Uses LLM reasoning for evidence-based grading. Designed for post-evaluation pass/fail determination.
---

# Output Grader Skill

Acts as an LLM judge to grade skill outputs against predefined assertions. Produces structured grading reports with concrete evidence for each pass or fail decision.

## Workflow

### Step 1: Load Assertions and Outputs

For each test case defined in `<skill-dir>/evals/evals.json`:

1. Read the `assertions` array.
2. Read the actual outputs from `evals/workspace/quality-results/iteration-<N>/eval-<id>/with_skill/outputs/` and `without_skill/outputs/`.
3. Read the `response.txt` and any output files for each run.

### Step 2: Evaluate Each Assertion

For each assertion, perform LLM-based evaluation against both with-skill and without-skill outputs:

**Grading criteria:**

- **PASS**: The assertion is demonstrably true. Must cite concrete evidence — a specific file name, a quoted string from the output, a measurable property (file size, count, format).
- **FAIL**: The assertion is false or cannot be verified from the available outputs.

**Evidence requirements:**
- Evidence must quote or reference specific output content (e.g., "Found `chart.png` (45KB) in `with_skill/outputs/`", or "Response line 12: 'X-axis: Q1, Q2, Q3, Q4'").
- "The output looks correct" is not evidence. Be specific.
- If the output is ambiguous, FAIL with evidence explaining the ambiguity.

### Step 3: Review Assertion Quality

Before finalizing grades, evaluate the assertions themselves for quality issues:

- **Always-pass assertion**: The behavior is present even without the skill. Flag with `assertion_quality: "always-pass"`. These assertions don't measure skill value and should be removed from future eval sets.
- **Always-fail assertion**: The behavior fails in both with-skill and without-skill configurations. Flag with `assertion_quality: "always-fail"`. These indicate either broken assertions or genuinely impossible tasks.
- **Unverifiable assertion**: Cannot be checked from outputs (e.g., "The code is secure" without a security scanner). Flag with `assertion_quality: "unverifiable"`.

### Step 4: Produce Grading Files

For each eval case, write per-configuration grading.json files that match the state-protocol schema:

`evals/workspace/quality-results/iteration-<N>/eval-<id>/with_skill/grading.json`:

```json
{
  "eval_id": "eval-1",
  "configuration": "with_skill",
  "assertion_results": [
    {
      "text": "Output includes a PNG image file",
      "passed": true,
      "evidence": "Found chart.png (45KB) in with_skill/outputs/",
      "assertion_quality": "valid"
    },
    {
      "text": "X-axis has month labels",
      "passed": true,
      "evidence": "Chart shows labels: Jan, Feb, Mar on x-axis (confirmed via chart metadata)",
      "assertion_quality": "valid"
    },
    {
      "text": "The chart uses optimal color scheme",
      "passed": false,
      "evidence": "Cannot objectively determine 'optimal' from output alone — subjective criterion",
      "assertion_quality": "unverifiable"
    }
  ],
  "summary": {"passed": 2, "failed": 1, "total": 3, "pass_rate": 0.66}
}
```

`evals/workspace/quality-results/iteration-<N>/eval-<id>/without_skill/grading.json`:

```json
{
  "eval_id": "eval-1",
  "configuration": "without_skill",
  "assertion_results": [
    {
      "text": "Output includes a PNG image file",
      "passed": false,
      "evidence": "No PNG files found in without_skill/outputs/. Only response.txt present.",
      "assertion_quality": "valid"
    },
    {
      "text": "X-axis has month labels",
      "passed": false,
      "evidence": "No chart generated",
      "assertion_quality": "valid"
    },
    {
      "text": "The chart uses optimal color scheme",
      "passed": false,
      "evidence": "No chart generated",
      "assertion_quality": "unverifiable"
    }
  ],
  "summary": {"passed": 0, "failed": 3, "total": 3, "pass_rate": 0.0},
  "assertion_quality_flags": ["unverifiable"]
}
```

`assertion_quality` is recorded per-assertion in both files. The `assertion_quality_flags` summary array collects all flagged assertion qualities in the without_skill file.

### Step 5: Summarize per Iteration

Optionally produce an aggregate summary across all eval cases in the iteration at `evals/workspace/quality-results/iteration-<N>/grading-summary.json`.

## Gotchas

- **Assertions that always pass in both configs**: If an assertion passes regardless of skill presence, it provides no signal. These should be removed from the eval set — they inflate pass rates without measuring skill value.
- **Assertions that always fail in both configs**: If both with-skill and without-skill fail, the task may be impossible given the tools available, or the assertion is poorly specified. Do not penalize the skill for these — flag and exclude from pass rate.
- **Don't give the benefit of the doubt**: If evidence is missing or incomplete, the assertion FAILS. "Probably correct" is not a grade. Require explicit, quotable evidence.
- **Subjective assertions are untestable**: Assertions like "the output is good" or "the code is clean" cannot be objectively graded. Flag as `unverifiable`. Rewrite to be specific: "The output includes a docstring for every public function" or "The code passes `ruff check` with zero errors."
- **With-skill and without-skill grading order matters**: Grade both configurations independently. Don't let one configuration's result bias the other. Explicitly read each set of outputs before making judgments.
- **File encoding issues**: Some output files may be binary (images, PDFs). For these, grade based on file existence, size, and format metadata rather than content inspection.

## Verification

Verify this skill produces correct output:

1. Create fixture outputs: with-skill has `chart.png` and a response mentioning axis labels; without-skill has only a text response.
2. Create assertions: "Output includes a PNG file" (should pass with-skill, fail without-skill), "Response mentions revenue" (should pass both), "Chart is interactive" (should fail both).
3. Run output-grader. Confirm grading.json has correct pass/fail per assertion per configuration.
4. Confirm "Response mentions revenue" is flagged as `assertion_quality: "always-pass"`.
5. Confirm "Chart is interactive" is flagged as `assertion_quality: "always-fail"`.
6. Confirm all PASS entries include specific evidence (not "looks good" or equivalent).
