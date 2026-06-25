# gixy

A git identity & SSH profile manager — switch author identity and SSH keys per GitHub account, automatically.

`gixy` handles gixy-specific commands (like profile management). For everything else, use `git` directly — there is no proxy layer.

---

## Features

- **Profile management** — define named author profiles (name + email), switch them globally in one command
- **Auto-activation** — map folder paths to profiles; gixy writes the right profile into each repo's local git config, so the correct identity & SSH key are used by every tool (your shell, Claude Code, IDEs, CI) — not just your interactive shell
- **SSH key management** — each profile gets its own ed25519 keypair; `profile global` symlinks the baseline keys to `~/.ssh/id_ed25519`
- **Config stored at** `~/.config/gixy/config`
- **Extensible** — designed to grow with more git-enhancing commands over time

---

## Installation

```sh
go install github.com/pawanprjl/gixy@latest
```

Requires Go 1.26.2+. Make sure `$GOPATH/bin` (or `$GOBIN`) is in your `PATH`.

---

## Usage

### Profile management

```sh
# Add a new profile (prompts for name and email, then generates an SSH keypair)
gixy profile add <name>

# List all saved profiles (* marks the profile that applies to the current directory)
gixy profile list

# Show which profile applies to the current directory and why (mapping + global baseline)
gixy profile status

# Show one profile's identity, SSH key, and the folders mapped to it
gixy profile show <name>

# Set the GLOBAL baseline identity + SSH key: writes git user.name/email to ~/.gitconfig
# and symlinks ~/.ssh/id_ed25519{,.pub}. This is what non-shell git (IDEs, GUIs, CI) and
# plain `ssh` see — per-directory auto-activation works independently of it.
gixy profile global <name>

# Show SSH key paths, fingerprint, and public key for a profile
gixy profile keys <name>

# Edit an existing profile (empty input keeps the current value)
gixy profile edit <name>

# Delete a profile
gixy profile delete <name>

# Pin the CURRENT repo to a profile by writing its identity + SSH key into the
# repo's local .git/config (so every git tool — Claude Code, IDEs, CI — uses it).
# Defaults to the profile mapped to this directory; pass a name to override.
gixy profile pin [name]

# Remove gixy's identity/SSH settings from the current repo's local config
gixy profile unpin
```

### Auto-activation

Map folder paths to profiles so the right identity and SSH key are used automatically for every `git` command, based on the directory you run it in.

**1. Set up shell integration (one-time)**

Add this to your `~/.zshrc`, `~/.bashrc`, or `~/.config/fish/config.fish`:

```sh
eval "$(gixy init)"
```

gixy auto-detects your shell. To specify explicitly: `gixy init --shell zsh|bash|fish`.

**2. Map folders to profiles**

```sh
# Map a directory (and all its subdirectories) to a profile
gixy profile map add work ~/projects/work
gixy profile map add personal ~/projects/personal

# List all mappings
gixy profile map list

# Remove a mapping
gixy profile map remove ~/projects/work

# Set a fallback profile for unmapped directories
gixy profile map default personal

# Clear the fallback profile
gixy profile map default --clear
```

When you first run `git` inside `~/projects/work/some-repo`, gixy resolves the matching profile and **writes its identity + SSH key into that repo's local `.git/config`** (`user.name`, `user.email`, `core.sshCommand`, plus a `gixy.profile` marker). From then on git itself applies the profile for **every** caller — your shell, Claude Code, IDEs, GUIs, CI — with no further gixy involvement. The most specific matching path wins, so `~/projects/work/client-acme` can override `~/projects/work`. If no mapping matches and no default is set, git runs unchanged.

Because the profile lives in the repo's own config, multiple terminals never interfere with each other, and multi-account GitHub works on plain `git@github.com` remotes (`core.sshCommand` uses `IdentitiesOnly=yes`, so the right key is offered without `~/.ssh/config` host aliases). The shell hook only runs when you change directory, and once a repo is synced it does no work at all.

**Respecting per-repo overrides:** if a repo has a hand-set local `user.email` and no `gixy.profile` marker, gixy never touches it.

**Pinning manually:** auto-sync only fires in a shell that has the hook. To pin a repo without it (or to refresh one after editing a profile's email/key), run `gixy profile pin` inside the repo. Revert with `gixy profile unpin`.

#### Known behavior

- gixy writes to each repo's **local** `.git/config` — untracked, reversible (`gixy profile unpin`), and inspectable via `git config --local --list`.
- Editing a profile's identity/key does not retroactively rewrite already-pinned repos with the same profile name; run `gixy profile pin` in those repos to refresh them. (Auto-sync re-pins only when a folder's mapped *profile* changes.)
- `gixy profile global <name>` still sets a global baseline (`~/.gitconfig` + `~/.ssh/id_ed25519` symlink) for repos that aren't mapped/pinned.

---

## Config storage

All gixy data is stored in a single JSON file:

```
~/.config/gixy/config
```

Example config:

```json
{
  "profiles": {
    "personal": {
      "name": "Jane Doe",
      "email": "jane@personal.dev"
    },
    "work": {
      "name": "Jane Doe",
      "email": "jane.doe@company.com"
    }
  },
  "path_mappings": {
    "/home/jane/projects/work": "work",
    "/home/jane/projects/personal": "personal"
  },
  "default_profile": "personal"
}
```

gixy sets file permissions to `0600` on write.

SSH keypairs are stored separately under `~/.ssh/gixy/<profile-name>/id_ed25519{,.pub}`.

---

## Roadmap

- [ ] Shell completions
- [x] `gixy profile status` — show currently active profile and which mapping triggered it

---

## Development

```sh
git clone https://github.com/pawanprjl/gixy.git
cd gixy
go run . <command>
```

**Built with:**

- [Go](https://golang.org/)
- [urfave/cli v3](https://github.com/urfave/cli) — CLI framework

---

## License

MIT — see [LICENSE](LICENSE)
