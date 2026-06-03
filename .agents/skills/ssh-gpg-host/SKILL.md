---
name: ssh-gpg-host
description: Host-specific constraints for SSH authentication, GPG operations, and git commit signing on Arch Linux. Use when the user mentions SSH, GPG, signing, agent, keys, git commit signing, authentication, pinentry, or ssh-agent.
---

# SSH & GPG Host Configuration (hinterland)

## Description

Host-specific constraints for SSH authentication, GPG operations, and git commit signing
on the Arch Linux system "hinterland." Activates when the user mentions SSH, GPG, signing,
agent, keys, git commit signing, authentication, pinentry, or ssh-agent.

## Architecture

This host uses a **unified gpg-agent** for both SSH and GPG. The standard `ssh-agent` is
intentionally masked. `SSH_AUTH_SOCK` is set in three layers:

| Layer | File | Scope |
|-------|------|-------|
| systemd | `~/.config/environment.d/60-ssh-auth.conf` | All systemd user services |
| Shell | `~/.config/bash/bashrc` lines 17–20 | Interactive shells (login + non-login) |
| X session | `~/.config/X11/xinitrc` lines 44–48 | X session children + dbus propagation |

The bash_profile at `~/.config/bash/bash_profile` sources `~/.bashrc` (shim) → `~/.config/bash/bashrc`
at line 85. It does **not** contain its own SSH exports (they were moved to bashrc).

Startx finds xinitrc via `XINITRC` env var set in bash_profile line 27.

## Constraints

- **DO NOT** enable or unmask `ssh-agent.socket` or `ssh-agent.service`. They are
  intentionally masked. Re-enabling poisons systemd user env with
  `SSH_AUTH_SOCK=/run/user/1000/ssh-agent.socket`.

- **DO NOT** add `SSH_AUTH_SOCK`, `GPG_TTY`, or `gpg-connect-agent` exports to
  `~/.config/bash/bash_profile`. They belong in `~/.config/bash/bashrc`.

- **DO NOT** add `GPG_AGENT_INFO` anywhere. GPG 1.x relic, unused by GPG 2.x.

- **DO NOT** mask `gnome-keyring-daemon.socket`. It runs without SSH component
  (`--components=pkcs11,secrets`), handling secret-storage only.

- **DO NOT** modify the `sshcontrol` file at `~/.gnupg/sshcontrol` without
  understanding keygrip registration. This file is deprecated upstream in favor
  of the `Use-for-ssh` key attribute. See `man gpg-agent`.

## Git Commit Signing

Enabled for default profile (`~/.config/git/pngdeity-github.gitconfig`):
`gpgsign = true`, `signingkey = 0D762ADE0E00C8DF`. Work repos under `~/itp/`
use `nsomers2-github.gitconfig` (no signing).

Git config at `~/.config/git/config` uses `pushInsteadOf` to rewrite HTTPS URLs
to SSH for push operations. SSH must be functional for `git push`.

Pinentry for commit signing depends on `GPG_TTY` and `updatestartuptty` set in bashrc.

## Key Files

| File | Role |
|------|------|
| `~/.config/environment.d/60-ssh-auth.conf` | Systemd user env (requires daemon-reload) |
| `~/.config/X11/xinitrc` | X session init |
| `~/.config/bash/bashrc` | Interactive shell config (SSH exports at lines 17–20) |
| `~/.config/bash/bash_profile` | Login shell config (no SSH exports) |
| `~/.gnupg/gpg-agent.conf` | Agent options, cache TTLs |
| `~/.gnupg/sshcontrol` | Registered SSH keygrips (deprecated) |
| `~/.config/git/config` | Git profiles, URL rewrites |
| `~/.config/git/pngdeity-github.gitconfig` | Default identity + signing |
| `~/.config/git/nsomers2-github.gitconfig` | Work identity (no signing) |

## Troubleshooting

If SSH stops working, check:

```
systemctl --user is-active gpg-agent.service      # must be: active
ls $(gpgconf --list-dirs agent-ssh-socket)         # must exist
echo $SSH_AUTH_SOCK                                 # must match gpgconf output
ssh-add -l                                          # must list ECDSA key
systemctl --user status ssh-agent.socket            # must be: masked, inactive
```

If `openssh` package update re-enables vendor presets:

```
systemctl --user mask ssh-agent.socket ssh-agent.service
```

Reload environment.d changes:

```
systemctl --user daemon-reload
```

## Reference

Full documentation: `~/downloads/ssh-gpg-configuration.md`
Man pages: `gpg-agent(1)`, `gpgconf(1)`, `environment.d(5)`, `ssh-add(1)`, `systemctl(1)`
