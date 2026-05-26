---
name: local-first
description: Discover and use man pages, --help, apropos, info pages, system logs, and kernel params before reaching for external documentation. Prefer over websearch or ctx7 for installed CLI tools, system utilities, and infrastructure commands. Local sources are version-matched and cost zero context tokens. When local sources are exhausted, fall back to skill("find-docs") → ctx7, then websearch last.
allowed-tools: bash, man, info, pacman, apropos, whatis
metadata:
  tags: "local man-pages help skills offline system infrastructure hardware discovery"
---

# Local-First Information Sources

Discover and use local documentation before external fetches.
Local sources are authoritative for installed tools and cost zero context tokens.

## Discovery Procedure

For any CLI tool or system utility, run discovery BEFORE reaching for
external docs. Each step takes <100ms. Stop at the first successful result.

### Step 1: Does the tool exist?
```bash
command -v <tool>        # exit 0 = installed, no output = not found
```

### Step 2: Does it have a man page?
```bash
man -w <tool>            # exit 0 = found, prints path
                         # exit 1 = no man page → try step 3
man -f <tool>            # one-line description (whatis)
```

### Man page sections

If a topic has pages in multiple sections, specify the section:
`man <section> <name>`

| Section | Covers | Example |
|---------|--------|---------|
| 1 | Executable programs / shell commands | `man 1 bash` |
| 5 | File formats and conventions | `man 5 crontab` |
| 7 | Miscellaneous (macros, conventions) | `man 7 regex` |
| 8 | System administration commands | `man 8 mount` |

### Step 3: Does it have --help or -h?
```bash
<tool> --help 2>&1 | head -20    # GNU-style help
<tool> -h 2>&1 | head -20        # Short help
<tool> help 2>&1 | head -20      # Subcommand help (gh, git, pacman)
<tool> <sub> --help 2>&1 | head -20  # Subcommand-specific
```

### Step 4: Don't know the tool name? Search by keyword.
```bash
apropos "<keyword>"              # Search all man page descriptions
man -k "<keyword>"               # Same as apropos
apropos -s <section> "<keyword>" # Filter by section (1,5,7,8)
# Common: apropos -s 8 for admin commands, -s 5 for config formats
```

### Step 5: Is it a GNU project? Check info pages.
```bash
info -w <topic>                  # Locate info page (exit 0 = found)
info <topic>                     # Read the info page directly
info -f <file>                   # Read a specific info file
```

### Step 6: What package owns this tool? (Arch)
```bash
pacman -Qo $(which <tool>)      # Which package provides this binary
pacman -Ql <pkg> | grep man      # Does the package ship man pages
```

## Decision Tree

Is the question about...
├─ A CLI tool installed on hinterland?
│  └─ Run discovery procedure above → use strongest local source found
│
├─ Infrastructure/hardware troubleshooting?
│  └─ Check dmesg, journalctl, sysctl, modinfo FIRST (see below)
│
├─ A library, framework, SDK, or cloud API?
│  └─ Local sources unlikely → use skill("find-docs") → ctx7 CLI
│
├─ A GitHub CLI pattern?
│  └─ skill("gh-cli-patterns") OR gh --help (discovery: `man -w gh`)
│
├─ Systemd services or journal?
│  └─ `man -w systemctl` is definitive → use man for directives
│
├─ A project workflow or convention?
│  └─ Check AGENTS.md, README.md, .apm/instructions/
│
└─ Local discovery exhausted?
   └─ Fall back to skill("find-docs") → ctx7, then websearch

## System-State Discovery (Infrastructure)

For hardware, kernel, and system-level questions, local sources
are THE authoritative answer — they match your exact kernel, drivers,
and firmware.

```bash
# Kernel parameters (matching your running kernel)
sysctl -a | grep <param>

# Kernel module info (matching your loaded modules)
modinfo <module>       # Parameters, version, dependencies
lsmod | grep <module>  # Is it loaded?

# Hardware detection (matches your actual hardware)
dmesg | grep -i <component>    # What the kernel saw at boot
lspci -k | grep -i <device>    # Which driver is handling it
lsusb | grep -i <device>       # USB device tree
lsblk -f                        # Block devices + filesystems

# System logging (your actual error messages)
journalctl -xe | grep -i <service>
journalctl -u <unit> --since "1 hour ago"

# Installed package inspection (Arch)
pacman -Qi <pkg>        # Version, description, deps, install date
pacman -Ql <pkg>        # All files owned by a package
pkgfile <filename>      # Which package provides a file
```

## When to use local-first vs find-docs

### Use local-first when...
- The question is about an installed tool's flags, behavior, or syntax
- Hardware troubleshooting (your exact drives, kernel, drivers)
- System administration (systemd units, mounts, networking)
- Package management (pacman/yay operations, dependencies)
- Any CLI tool where `man -w <tool>` exits 0

### Use skill("find-docs") → ctx7 when...
- Library, framework, SDK, or cloud service API (React, APM, Nomad)
- Cloud provider APIs (DigitalOcean Terraform, Supabase, AWS)
- Tools NOT installed on this system
- Config file formats that man section 5 doesn't cover

### Concrete examples

| Wrong path | Right path |
|---|---|
| ctx7 for smartctl flags | `man -w smartctl` finds man8 instantly |
| Websearch "how to mount ext4" | `man -w mount` + `man -w fdisk` exist |
| ctx7 for systemd directives | `man -w systemd.service` is comprehensive |
| Websearch "pacman remove orphan" | `pacman --help` + `man pacman` |
| Websearch "nvme drive temperature" | `man -w nvme` + `man -w smartctl` |
| ctx7 for React hooks | Correct — library, not CLI tool |
| ctx7 for APM manifest schema | Correct — SDK/library reference |
