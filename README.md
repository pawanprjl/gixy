# gixy

A git proxy CLI that adds superpowers to your everyday git workflow.

`gixy` wraps git transparently — any command that isn't a gixy-native feature is passed directly to git as-is. This means you can alias `git` to `gixy` and never notice the difference, while gaining extra capabilities on top.

---

## Features

- **Profile management** — define named author profiles (name + email) and switch between them per-repo in one command
- **Transparent git passthrough** — any unrecognized command is forwarded verbatim to `git`
- **XDG-compliant config** — profiles stored in `~/.config/gixy/config` (respects `$XDG_CONFIG_HOME`)
- **Extensible** — designed to grow with more git-enhancing commands over time

---

## Installation

```sh
go install github.com/pawanprjl/gixy@latest
```

Requires Go 1.21+. Make sure `$GOPATH/bin` (or `$GOBIN`) is in your `PATH`.

---

## Usage

### Profile management

```sh
# Add a new profile
gixy profile add <profile-name>
# Prompts for: Name, Email

# List all saved profiles
gixy profile list

# Apply a profile to the current repository
gixy profile use <profile-name>
# Writes user.name and user.email to the repo's local .git/config
```

### Git passthrough

Any command not matched by gixy is forwarded directly to `git`:

```sh
gixy commit -m "feat: add something"
gixy push origin main
gixy log --oneline
gixy status
```

This means you can safely use `gixy` as a drop-in replacement for `git`.

---

## Profile storage

Profiles are stored in:

```
$XDG_CONFIG_HOME/gixy/config
```

Which defaults to `~/.config/gixy/config` when `$XDG_CONFIG_HOME` is not set.

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
  }
}
```

---

## Roadmap

- [ ] SSH key path per profile
- [ ] GPG signing key per profile
- [ ] `profile show` — display the active profile for the current repo
- [ ] Shell completions

---

## Development

```sh
git clone https://github.com/pawanprjl/gixy.git
cd gixy
go run . <command>
```

**Built with:**

- [Go](https://golang.org/)
- [urfave/cli](https://github.com/urfave/cli) — CLI framework

---

## License

MIT
