# gixy

A git identity & SSH profile manager — switch author identity and SSH keys per GitHub account, automatically.

`gixy` handles gixy-specific commands (like profile management). For everything else, use `git` directly — there is no proxy layer.

---

## Features

- **Profile management** — define named author profiles (name + email), switch them globally in one command
- **Auto-activation** — map folder paths to profiles; the correct profile activates automatically when you `cd` into a mapped directory
- **SSH key management** — each profile gets its own ed25519 keypair; `profile use` symlinks the active keys to `~/.ssh/id_ed25519`
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

# List all saved profiles (* marks the active profile)
gixy profile list

# Activate a profile globally: sets git user.name/email in ~/.gitconfig
# and symlinks ~/.ssh/id_ed25519{,.pub} to the profile's keypair
gixy profile use <name>

# Show SSH key paths, fingerprint, and public key for a profile
gixy profile keys <name>

# Edit an existing profile (empty input keeps the current value)
gixy profile edit <name>

# Delete a profile
gixy profile delete <name>
```

### Auto-activation (pyenv-style)

Map folder paths to profiles so the right identity activates automatically when you `cd` into a directory.

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
gixy profile default personal

# Clear the fallback profile
gixy profile default --clear
```

When you `cd` into `~/projects/work/some-repo`, gixy automatically runs `profile use work` in the background — switching your global git identity and SSH keys. The most specific matching path wins, so `~/projects/work/client-acme` can have its own mapping that overrides `~/projects/work`.

If no mapping matches and no default is set, gixy does nothing.

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
- [ ] `gixy profile status` — show currently active profile and which mapping triggered it

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
