---
name: conventional-commits
description: Enforce the Conventional Commits specification for git commit messages. Use when writing commit messages, setting up commit linting, or configuring automated versioning and changelog generation.
allowed-tools: git
metadata:
  tags: "git commits conventional-commits versioning changelog"
compatibility: Generic — requires git. No environment restrictions.
---

# Conventional Commits

## Description
Enforces the [Conventional Commits](https://www.conventionalcommits.org/) specification for all git commit messages within this repository. This standard provides a lightweight convention on top of commit messages, ensuring a clear and machine-readable commit history which facilitates automated versioning and changelog generation.

## Instructions

### 1. Commit Message Format
All commit messages must follow this structural pattern:
```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### 2. Mandatory Types
You must use one of the following types to categorize your changes:
- **feat**: A new feature.
- **fix**: A bug fix.
- **docs**: Documentation only changes.
- **style**: Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc).
- **refactor**: A code change that neither fixes a bug nor adds a feature.
- **perf**: A code change that improves performance.
- **test**: Adding missing tests or correcting existing tests.
- **build**: Changes that affect the build system or external dependencies (example scopes: gulp, broccoli, npm).
- **ci**: Changes to our CI configuration files and scripts (example scopes: Travis, Circle, BrowserStack, SauceLabs).
- **chore**: Other changes that don't modify src or test files.
- **revert**: Reverts a previous commit.

### 3. Guidelines
- **Subject Line**: The subject line must be a succinct description of the change. Use the imperative, present tense: "change" not "changed" nor "changes". Do not capitalize the first letter and do not end with a period.
- **Scope (Optional)**: A scope may be provided to a commit's type, to provide additional contextual information and is contained within parenthesis, e.g., `feat(parser): add ability to parse arrays`.
- **Breaking Changes**: Breaking changes must be indicated by an `!` after the type/scope, or by including `BREAKING CHANGE:` at the beginning of the footer section. A breaking change can be part of a commit of any type.
- **Body (Optional)**: Use the body to explain the "what" and "why" of the change, as opposed to the "how".
- **Footer (Optional)**: Use the footer to reference GitHub issues that this commit closes or to provide other metadata.

### 4. Verification
Before finalizing a commit, verify that it adheres to this specification. If using an automated tool or agent to commit, ensure it is configured to follow these rules strictly.
