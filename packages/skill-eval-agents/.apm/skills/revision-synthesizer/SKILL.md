---
name: revision-synthesizer
description: Synthesize improved SKILL.md variants from eval signals. Use after grading to generate revised skill versions incorporating feedback from assertion failures, human review, and execution transcripts.
allowed-tools: skills-ref, bash
metadata:
  tags: "revision synthesis optimization variants improvement feedback skills-ref"
compatibility: Requires skills-ref for self-validation. Designed to run in parallel with different strategy parameters for variant generation.
---

# Revision Synthesizer Skill

Generates improved SKILL.md variants by synthesizing signals from grading failures, trigger rates, execution transcripts, and (optionally) human feedback. This agent is designed to be invoked **3 times in parallel** with different strategy parameters to produce diverse candidates for the candidate-selector.

## Variant Strategies

This skill must be invoked with a `--variant` parameter set to one of:

### Variant A — Minimalist
Cut instructions, simplify, remove redundant steps. Reduce token count while maintaining or improving pass rate. The hypothesis: shorter skills reduce distraction and let the model focus on what matters.

Tactics:
- Remove verbose explanations and restatements of obvious facts
- Merge repeated verification steps into a single check
- Eliminate "optional" or "nice-to-have" sections that rarely affect output quality
- Replace multi-step sequences with single consolidated instructions where the model handles decomposition

### Variant B — Additive
Add examples, edge case handling, clearer step descriptions, more detailed verification. The hypothesis: richer instructions produce more consistent output across diverse inputs.

Tactics:
- Add concrete examples for each instruction step (input → expected output)
- Add explicit edge case handling (e.g., "If the file is empty, return early with a warning")
- Expand verification steps with pass/fail criteria
- Add error recovery paths for common failure modes found in grading results

### Variant C — Restructured
Reorganize instruction flow, change framing, use a different procedural structure. The hypothesis: the current structure causes the model to miss context or follow steps out of order.

Tactics:
- Reorder the workflow: put high-impact steps first, verification steps immediately after their relevant instructions
- Change framing: shift from imperative checklist to decision-tree or state-machine structure
- Group related steps under thematic sub-headings
- Move reference material to end (the model reads sequentially, so what it reads last has higher recency bias)

## Workflow

### Step 1: Gather All Signals

Read the following inputs:

1. **Current SKILL.md**: The full source skill file.
2. **grading.json files**: All pass/fail results from the output-grader, both with-skill and without-skill configurations.
3. **Execution transcripts** (if available): Full agent response texts from quality-evaluator outputs at `evals/workspace/iteration-*/eval-*/{with_skill,without_skill}/outputs/response.txt`.
4. **Timing data**: `timing.json` from each eval run.
5. **Human feedback** (if available): Any review notes or annotations from a human evaluator.

### Step 2: Analyze Failure Patterns

From grading results, identify:

- Which assertions fail consistently (across multiple eval cases)?
- Do failures cluster around specific instruction types (e.g., output format instructions, verification steps)?
- Is the skill helping? Compare with-skill vs without-skill pass rates. If without-skill is higher on some assertions, the skill may be actively harmful for those cases.
- Do certain instructions produce over-constrained output that the assertion checker disagrees with?

### Step 3: Generate Revision

Apply the variant strategy's tactics to produce a revised SKILL.md. Follow agentskills.io best practices:

- **Generalize from feedback, don't narrow-patch**: If "the chart lacks axis labels" fails, don't add "always add axis labels to charts." Instead, add a general instruction like "Verify all visual outputs include proper labels and units."
- **Keep skill lean**: Every instruction should have a demonstrable impact on at least one eval assertion. If an instruction never affects grading, consider removing it.
- **Explain the why**: Prefer reasoning-based instructions ("Check artifact completeness because partial outputs appear correct but lack required metadata") over arbitrary rules ("Always run step 3 before step 4").
- **Bundle repeated work**: If multiple steps require reading the same file, have them read it once and reference the result.

### Step 4: Write Output

Write the complete revised SKILL.md to:

```
evals/workspace/revisions/revision-<variant>.md
```

This must be a **complete valid SKILL.md**, not a diff or a partial fragment. Include:
- Full YAML frontmatter with `name` and `description`
- Complete markdown body with all sections

The `name` field must match the original skill name (not the variant name). The `description` field may be updated if the trigger-aggregator or failure analysis suggests changes.

### Step 5: Self-Validate

Before writing the final file, validate the revision:

```bash
skills-ref validate <output-path>
```

If validation fails, fix the issues and re-validate. Do not write an invalid SKILL.md. If the revision cannot be made valid after 3 attempts, write the best-effort version and annotate with a `<!-- VALIDATION_FAILED: <reason> -->` comment at the top of the file.

### Step 6: Document Changes

Add a `## Revision Notes` section at the bottom of the SKILL.md (or as a separate `revision-<variant>-notes.md` file) documenting:
- Which failures drove which changes
- Why the variant strategy was chosen for each change
- Expected impact on pass rate (qualitative prediction)
- Token count delta vs. original

## Gotchas

- **Explain WHY each change was made**: The revision notes are as important as the revision itself. Downstream humans need to understand the rationale. "Removed step 4 because" is insufficient — explain what failure pattern step 4 was causing.
- **Avoid adding instructions the model already knows**: If the model consistently fails to "use kebab-case for filenames," adding "remember to use kebab-case" won't help. Instead, restructure the instruction to include a verification step: "After creating files, check that all filenames are in kebab-case."
- **Removing instructions can be more effective than adding**: If pass rates plateau despite additive changes, the skill may be too long, causing the model to skip instructions. The minimalist variant exists specifically to test this hypothesis.
- **Output is a complete valid SKILL.md, not a diff**: The candidate-selector reads the revision file directly as a skill source. Partial files or diff patches will break the pipeline.
- **Validation must pass**: An invalid SKILL.md will cause the candidate-selector to skip this variant with no partial credit. Validate before writing.
- **Don't optimize the description alone**: If the trigger-aggregator already optimized the description, don't undo that work. If the description needs further changes, justify them in revision notes and cross-reference the trigger-aggregator's optimized version.

## Verification

Verify this skill produces correct output:

1. Create a fixture SKILL.md with known issues: verbose redundant steps, missing edge case handling, poor instruction ordering.
2. Create grading.json showing 3 assertion failures: one related to output format, one to missing verification, one to edge case handling.
3. Invoke revision-synthesizer with `--variant B` (Additive).
4. Confirm the output at `evals/workspace/revisions/revision-B.md` is a complete valid SKILL.md (runs `skills-ref validate` clean).
5. Confirm the output includes added examples, edge case handling, and more detailed verification steps.
6. Invoke revision-synthesizer with `--variant A` (Minimalist).
7. Confirm the minimalist revision is shorter (fewer lines or lower token count) than the original while still passing validation.
8. Confirm each revision includes `## Revision Notes` documenting the rationale for changes.
