# Stack Compatibility Map

Maps detected build system files to `compatibility_tokens` for skill pre-filtering.
Used by `init-project-guidance` Step 3 (Dynamic Skill Linking).

## Detection → Token Mapping

| Detected File | compatibility_tokens |
|---|---|
| `*.csproj`, `*.sln` | `["dotnet", "generic"]` |
| `package.json` | `["node", "npm", "javascript", "generic"]` |
| `requirements.txt`, `pyproject.toml`, `uv.lock` | `["python", "generic"]` |
| `Makefile`, `GNUmakefile` | `["make", "generic"]` |
| `PKGBUILD` | `["arch", "generic"]` |
| `main.tf`, `*.tf` | `["terraform", "iac", "generic"]` |
| `Dockerfile` | `["docker", "generic"]` |
| `Cargo.toml` | `["rust", "generic"]` |
| `go.mod` | `["go", "generic"]` |
| No build system detected | `["generic"]` |

## Filter Logic

1. Read `skill-index.json` from the scaffolding repo root.
2. For each skill, check if `skill.compatibility_tokens ∩ stack.tokens` is non-empty.
3. If intersection is non-empty → include the skill.
4. If intersection is empty → exclude with reason logged.
5. `custom-skill-creator` is always included (meta-skill, universally needed).
6. The `compatibility` prose field on each skill is human-readable documentation only — matching uses `compatibility_tokens`.

## Token Semantics

- **`generic`**: No environment restrictions. Always present in both skill and stack tokens. Guarantees Generic skills pass all filters.
- **`node`**: Requires Node.js runtime.
- **`make`**: Requires Make build tool.
- **`dotnet`**: Requires .NET SDK.
- **`python`**: Requires Python runtime.
- **Stack-specific tokens** (`terraform`, `docker`, `rust`, `go`, `arch`, `npm`, `javascript`): Indicate tool or ecosystem specificity.

## Token Field in skill-index.json

```json
{
  "skills": {
    "architectural-review": {
      "compatibility_tokens": ["generic"]
    },
    "ci-cd-pipeline": {
      "compatibility_tokens": ["make", "generic"]
    }
  }
}
```
