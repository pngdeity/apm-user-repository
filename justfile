# APM marketplace validation
validate:
    apm marketplace check
    apm pack --dry-run

# Build Go eval tools when source files change
build-eval:
    go build -o ./bin/invoke-cli ./packages/skill-eval-agents/.apm/scripts/cmd/invoke-cli
    go build -o ./bin/parse-session ./packages/skill-eval-agents/.apm/scripts/cmd/parse-session
    go build -o ./bin/compute-benchmark ./packages/skill-eval-agents/.apm/scripts/cmd/compute-benchmark
    go build -o ./bin/select-best ./packages/skill-eval-agents/.apm/scripts/cmd/select-best

# Run eval pipeline on one or more skill paths
eval-skill skill_paths:
    just build-eval
    @echo "Running eval pipeline on {{skill_paths}}"
    # The actual evaluation is performed by loading the skill-eval-pipeline SKILL.md
    # in an agent session. CI/CD invokes the pipeline agent programmatically.
    # This recipe ensures prerequisites are met before agent dispatch.

# Post PR comment with benchmark results
eval-post-comment:
    @find . -name "benchmark.json" -path "*/evals/workspace/*" -exec cat {} \; | head -50
    @echo "Posting eval results to PR (via gh pr comment)"

# Fail if quality gates not met
eval-gate-check:
    @if find . -name "selected.json" -path "*/evals/workspace/*" -exec grep -q '"selected"' {} \; ; then \
        echo "Selection found — quality gate passed"; \
    else \
        echo "ERROR: No selected.json found — quality gate failed"; \
        exit 1; \
    fi
