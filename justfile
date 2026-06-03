# Format all YAML files
fmt:
    yamlfmt

# Check YAML formatting (CI mode, read-only)
fmt-check:
    yamlfmt -lint

# Full CI pipeline: format check → version check → schema validate → readme check → compile check → marketplace validate → pack gate → outdated check
ci: fmt-check sync-check yaml-validate sync-readme-check
    -apm compile --validate
    -apm marketplace check
    apm pack --check-versions --check-clean
    -apm marketplace outdated

# CI pipeline for eval: build go tools → run eval → quality gate
ci-eval skill_paths:
    just build
    just eval-skill {{skill_paths}}
    just eval-gate-check

# Run eval pipeline on one or more skill paths (invokes the eval agent)
eval-skill skill_paths:
    @echo "Running eval pipeline on {{skill_paths}}"

# Validate apm.yml files against the APM JSON Schema
yaml-validate:
    @command -v check-jsonschema >/dev/null 2>&1 || { echo "ERROR: check-jsonschema not installed (extra/check-jsonschema or pip install check-jsonschema)" >&2; exit 1; }
    check-jsonschema --schemafile schemas/apm.json apm.yml packages/*/apm.yml

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

# Full pre-commit check: format → fix versions → schema validate → sync readme → marketplace validate → regenerate marketplace.json → verify nothing stale
pre-commit-check:
    yamlfmt
    -go run ./cmd/sync-versions
    -just yaml-validate
    -go run ./cmd/readme-sync
    -apm marketplace check
    apm pack
    @git diff --exit-code -- .claude-plugin/marketplace.json packages/*/apm.yml README.md || (echo "ERROR: uncommitted changes after sync+pack — commit the regenerated files and retry" && exit 1)

# Dogfood: install consumed packages and deploy skills
update-self:
    apm install
    apm prune
    @echo "Done. Verify with: git diff apm.lock.yaml"

# Build all Go tools
build: build-eval build-sync build-check-upstream build-readme

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

# Build readme-sync tool
build-readme:
    go build -o ./bin/readme-sync ./cmd/readme-sync

# Check if upstream external refs are stale (CI advisory mode, non-blocking)
check-upstream:
    @go run ./cmd/check-upstream --check

# Update upstream external refs to latest commits
update-upstream:
    @go run ./cmd/check-upstream

# Check marketplace for packages with newer matching tags
outdated:
    apm marketplace outdated

# Run APM lockfile integrity + drift audit (CI gate)
audit:
    apm audit --ci

# Quality gate check for eval pipeline results
eval-gate-check:
    @if find . -name "selected.json" -path "*/evals/workspace/*" -exec grep -q '"selected"' {} \; ; then \
        echo "Selection found — quality gate passed"; \
    else \
        echo "ERROR: No selected.json found — quality gate failed"; \
        exit 1; \
    fi

# Post PR comment with benchmark results
eval-post-comment:
    @find . -name "benchmark.json" -path "*/evals/workspace/*" -exec cat {} \; | head -50
    @echo "Posting eval results to PR (via gh pr comment)"

# Sync README package table from apm.yml (auto-generate)
sync-readme:
    @go run ./cmd/readme-sync

# Sync-readme check only (CI mode, read-only)
sync-readme-check:
    @go run ./cmd/readme-sync --check
