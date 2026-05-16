---
name: struct-validator
description: Validate SKILL.md files for structural compliance against the agentskills.io specification. Use when a new or modified SKILL.md needs frontmatter validation, directory schema verification, or format compliance checking.
allowed-tools: skills-ref
metadata:
  tags: "validation structure compliance yaml frontmatter skill-format"
compatibility: Requires skills-ref CLI. Works against any SKILL.md file following the agentskills.io specification.
---

# Struct Validator Skill

Validates SKILL.md files against the agentskills.io structural specification. This agent is deterministic — no LLM reasoning is needed for the validation pass. Only execution of the validation tool and structured reporting of results.

## Workflow

### Step 1: Run Validation

Execute the `skills-ref validate` command on the target skill directory:

```bash
skills-ref validate <path-to-skill-directory>
```

Capture both stdout and stderr. The command will return a non-zero exit code on failure.

### Step 2: Parse Output

Parse the CLI output for errors. The `skills-ref validate` output follows this pattern:
- Exit code 0: validation passed
- Exit code non-zero: one or more structural errors found, each reported as a line beginning with `ERROR:` or `WARN:`

### Step 3: Write Results

If the command exits clean (exit code 0), write the pass result:

```json
{"status": "pass", "errors": []}
```

If errors are present, format them as a structured JSON report. Each error entry includes the error message, severity (`error` or `warn`), and the file path:

```json
{
  "status": "fail",
  "exit_code": 2,
  "errors": [
    {
      "severity": "error",
      "file": "SKILL.md",
      "message": "name field 'StructValidator' does not match directory name 'struct-validator'"
    },
    {
      "severity": "error",
      "file": "SKILL.md",
      "message": "Missing YAML document start delimiter '---' at line 1"
    }
  ]
}
```

Write the result to `evals/workspace/struct-validation.json`.

### Step 4: Report to Pipeline

Return the result object for downstream pipeline steps. If `status` is `fail`, stop the eval pipeline here — no further stages can proceed on an invalid SKILL.md.

## Gotchas

- **Name mismatch with directory**: The `name` field in YAML frontmatter must exactly match the parent directory name (kebab-case). Even character-case differences trigger a validation failure.
- **Uppercase in name field**: The `name` field must be kebab-case. Any uppercase letters will fail validation even if the directory name matches structurally.
- **Missing YAML delimiters**: The frontmatter block must start and end with `---` on its own line. A common mistake is placing content before the opening `---` or omitting the closing `---`.
- **Trailing whitespace in description**: Some validators reject `description` fields with trailing whitespace. Trim all frontmatter values.
- **Description over 1024 characters**: The `description` field is capped at 1024 characters per spec. The validator will flag longer descriptions even if they appear structurally valid.
- **Stale validation cache**: The `skills-ref` tool may cache prior validation results. Run with `--no-cache` if you get unexpected pass results on known-bad files.
- **UTF-8 BOM**: Files saved with a UTF-8 byte order mark can cause the YAML parser to fail. Ensure files are plain UTF-8 without BOM.

## Verification

Verify this skill produces correct output:

1. Create a fixture skill directory at `/tmp/test-skill/skills/fixture-validator/` with an intentionally invalid `SKILL.md` (missing closing `---`, name field mismatch, description over 1024 chars).
2. Run `skills-ref validate /tmp/test-skill/skills/fixture-validator/` and confirm non-zero exit.
3. Execute the struct-validator workflow against the fixture. Confirm `struct-validation.json` is written with `"status":"fail"` and at least 3 error entries.
4. Fix the fixture to be valid and re-run. Confirm `struct-validation.json` reads `{"status":"pass","errors":[]}`.
5. Confirm that a description field with exactly 1025 characters is flagged; one with exactly 1024 characters passes.
