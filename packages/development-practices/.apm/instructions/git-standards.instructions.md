---
description: Git and version control standards — upstream-first workflow, signed commits, safety protocols
applyTo: "**"
---

## Git & Version Control
- **Upstream First:** Maintain a strict "upstream-first" workflow for all forks.
- **Safety:** Always confirm with the user before committing to the default branch (e.g., `main` or `master`).
- **Commits:** Mandatory commit signing is required. If `git commit -S` fails, stop and provide the exact command for manual signing.
- **Workflow Tools:** 
    - Use `git refresh` to sync `main` with `upstream/master`.
    - Observe the `post-checkout` hook warnings regarding stale branches.

- **Exclusion Rule:** NEVER commit the `.gemini/` or `.agents/` directories to version control.
