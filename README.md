# gixy

A CLI companion for git that adds workflow enhancements on top of your existing git setup.

`gixy` handles gixy-specific commands (like profile management). For everything else, use `git` directly — there is no proxy layer.

---

## Features

- **Profile management** — define named author profiles (name + email) and switch between them per-repo in one command
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
