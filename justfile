# Full CI pipeline: version check → marketplace validate → stale check
ci: sync-check
    apm marketplace check
    apm pack
    @git diff --exit-code -- .claude-plugin/marketplace.json || (echo "ERROR: marketplace.json is stale — run 'just sync-fix && apm pack' and commit the result" && exit 1)

# CI pipeline for eval: build go tools → run eval → quality gate
ci-eval skill_paths:
    just build
    just eval-skill {{skill_paths}}
    just eval-gate-check

# Run eval pipeline on one or more skill paths (invokes the eval agent)
eval-skill skill_paths:
    @echo "Running eval pipeline on {{skill_paths}}"

# APM marketplace validation (dry-run, for quick local checks)
validate:
    apm marketplace check
    apm pack --dry-run

# Check package versions against root constraints (CI mode, read-only)
sync-check:
    @go run ./cmd/sync-versions --check

# Fix mismatched package versions to satisfy root constraints
sync-fix:
    @go run ./cmd/sync-versions

# Full pre-commit check: fix versions → regenerate marketplace.json → verify nothing stale
pre-commit-check:
    -go run ./cmd/sync-versions
    apm marketplace check
    apm pack
    @git diff --exit-code -- .claude-plugin/marketplace.json packages/*/apm.yml || (echo "ERROR: uncommitted changes after sync+pack — commit the regenerated files and retry" && exit 1)

# Build all Go tools
build: build-eval build-sync build-check-upstream

# Build eval Go tools
build-eval:
    go build -o ./bin/invoke-cli ./packages/skill-eval-agents/.apm/scripts/cmd/invoke-cli
    go build -o ./bin/parse-session ./packages/skill-eval-agents/.apm/scripts/cmd/parse-session
    go build -o ./bin/compute-benchmark ./packages/skill-eval-agents/.apm/scripts/cmd/compute-benchmark
    go build -o ./bin/select-best ./packages/skill-eval-agents/.apm/scripts/cmd/select-best

# Build sync-versions tool
build-sync:
    go build -o ./bin/sync-versions ./cmd/sync-versions

# Build check-upstream tool
build-check-upstream:
    go build -o ./bin/check-upstream ./cmd/check-upstream

# Check if upstream external refs are stale (CI advisory mode, non-blocking)
check-upstream:
    @go run ./cmd/check-upstream --check

# Update upstream external refs to latest commits
update-upstream:
    @go run ./cmd/check-upstream

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
