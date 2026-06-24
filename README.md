# gixy

A git identity & SSH profile manager — switch author identity and SSH keys per GitHub account, automatically.

`gixy` handles gixy-specific commands (like profile management). For everything else, use `git` directly — there is no proxy layer.

---

## Features

- **Profile management** — define named author profiles (name + email), switch them globally in one command
- **Auto-activation** — map folder paths to profiles; the correct profile is applied automatically for every `git` command, based on the directory you run it in
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

When you run `git` inside `~/projects/work/some-repo`, gixy resolves the matching profile and injects its identity and SSH key **into that single command only** — without touching `~/.gitconfig` or the `~/.ssh/id_ed25519` symlink. The most specific matching path wins, so `~/projects/work/client-acme` can have its own mapping that overrides `~/projects/work`. If no mapping matches and no default is set, git runs unchanged.

Because each `git` invocation resolves independently, multiple terminals in different projects never interfere with each other — and multi-account GitHub works on plain `git@github.com` remotes (gixy sets `GIT_SSH_COMMAND` with `IdentitiesOnly=yes`, so the right key is offered without `~/.ssh/config` host aliases). gixy is only invoked when you change directory, so repeated git commands in the same repo add no overhead.

**Respecting per-repo overrides:** if a repo has an explicit local `user.email` (or `core.sshCommand`), gixy leaves it alone so git's normal precedence applies.

#### Known limitations

- `git config user.name` reflects your **global baseline** (set by `gixy profile global`), not the per-command injected profile — though commits are stamped with the correct identity. Use `gixy profile global <name>` to set the baseline that non-shell tools see.
- Tools that invoke git outside your interactive shell (IDEs, GUIs, CI, `/usr/bin/git`) bypass the wrapper and use that global baseline.

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
