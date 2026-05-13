# Stack-Specific Guidance Rules

This reference expands on the technology-specific rules scaffolded into `GEMINI.md`

## C# / .NET
- **Target Framework:** .NET 10, C# 14
- **Code Style:** Prefer `sealed` modifiers on classes not designed for inheritance
- **Framework:** Microsoft Agent Framework v1.3.0 for agentic workloads
- **Build:** `dotnet build` with `--configuration Release`
- **Test:** `dotnet test` with `--logger trx`
- **Lint:** `dotnet format` with `.editorconfig`

## Python
- **Package Manager:** Use `uv` for all dependency management
- **Virtual Environment:** Activate venv before any operation: `source .venv/bin/activate` or `uv run`
- **Install:** `uv pip install -r requirements.txt` or `uv sync`
- **Lint:** `ruff check .` and `ruff format --check .`
- **Test:** `pytest` with `--cov` for coverage
- **Type Check:** `mypy .` (if configured)

## JavaScript / Node.js
- **Package Manager:** npm with isolation: `npm install --cache "$srcdir/npm-cache"`
- **Build:** `npm run build` (if defined) or `npx tsc` for TypeScript
- **Test:** `npm test` (if defined) or `npx jest`
- **Lint:** `npx eslint .` (if configured)
- **Distribution:** Strip npm-internal metadata from `package.json` in dist:
  Remove `devDependencies`, `scripts`, and `jest`/`eslint` config keys before publishing

## Arch Linux / AUR
- **Package Standard:** Follow `PKGBUILD(5)` man page specifications
- **Tool:** Use `pkgctl` for package management: `pkgctl build`, `pkgctl release`
- **Integrity:** Verify checksums with `updpkgsums` before building
- **Lint:** `namcap PKGBUILD` and `namcap *.pkg.tar.zst`
- **Commit:** Use AUR-compliant `.SRCINFO` generation via `makepkg --printsrcinfo > .SRCINFO`

## Infrastructure as Code (IaC)
- **Pattern:** Composition pattern — small, reusable modules over monolithic configs
- **Decision Tables:** Use markdown decision tables in `DECISIONS.md` to document tradeoffs
- **Validation:** `terraform validate` before `terraform plan`
- **Lint:** `tflint` or `terraform fmt -check -recursive`
- **Secrets:** Never hardcode — use variables with `sensitive = true` or external vault
