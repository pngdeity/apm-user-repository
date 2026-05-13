---
name: config-modernizer
description: Audit, update, and modernize software configuration files (dotfiles) with a focus on technical accuracy and system compatibility. Use when reviewing terminal emulator configs, shell rc files, editor settings, or any program configuration.
allowed-tools: bash, grep, read_file
metadata:
  tags: "configuration dotfiles modernization audit shell terminal"
compatibility: Requires bash shell and manual page access. Designed for Linux environments.
---
# Skill: Configuration Modernizer

A specialized workflow for auditing, updating, and modernizing software configuration files (dotfiles) with a focus on technical accuracy and system compatibility.

## Instructions

### 1. Discovery & Context Gathering
*   **Identify the Program:** Use `apropos <program>` and `man -f <program>` to locate binary paths and manual sections.
*   **Locate Configs:** Read the `FILES` section of the program's primary man page. This is the authoritative source for finding where the program expects its configuration (e.g., `~/.config/`, `/etc/`, `~/.<program>rc`).
*   **Check the Version:** Run `<program> --version` and check the local package manager (e.g., `yay -Q <program>`) to ensure you are auditing for the correct release. Features and syntax change significantly between versions.

### 2. Environmental Audit
*   **Environment Variables:** Check `$TERM`, `$SHELL`, `$XDG_SESSION_TYPE`, and `$DISPLAY`/`$WAYLAND_DISPLAY`. Many configurations (especially for CLI tools and terminal emulators) behave differently based on these values.
*   **Sandbox Detection:** Identify if the program is running inside a multiplexer (Check `$TMUX`) or a container, as this often requires "passthrough" or "escape" configurations.

### 3. The Review (The Four C's)
Perform a line-by-line audit of all identified configuration files using these criteria:
*   **Correctness:** Identify typos, invalid options, or deprecated syntax. Use shell-native check commands if available (e.g., `bash -n`, `python -m py_compile`).
*   **Completeness:** Compare the current files against the program's default templates (often found in `/usr/share/doc/` or the official GitHub repository). Identify missing standard features.
*   **Currentness:** Ensure the configuration leverages modern protocols (e.g., GPU acceleration, Wayland support, Kitty graphics protocol) instead of legacy fallbacks.
*   **Best Practices:** Audit for performance (e.g., optimized search commands, efficient redraw rates) and readability (e.g., logical grouping, helpful comments).

### 4. Strategy & Implementation
*   **Document Findings:** Present a concise summary of issues found before proposing changes.
*   **Surgical Edits:** Use the `replace` tool for targeted updates. Do not overwrite entire files unless they are beyond repair or the user requests a factory reset.
*   **Verification:** After editing, always run a syntax check or dry-run of the program to ensure the new configuration is valid.

## Mandatory Safety Rule
*   **User Confirmation:** You MUST present your findings and proposed plan to the user. Do not modify configuration files until a Directive is issued.
