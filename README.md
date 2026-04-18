# gixy

A CLI companion for git that adds workflow enhancements on top of your existing git setup.

`gixy` handles gixy-specific commands (like profile management). For everything else, use `git` directly — there is no proxy layer.

---

## Features

- **Profile management** — define named author profiles (name + email) and switch between them per-repo in one command
- **AI commit messages** — generate conventional commit messages from staged changes using AI (currently supports Gemini)
- **XDG-compliant config** — config stored in `~/.config/gixy/config` (respects `$XDG_CONFIG_HOME`)
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

# Delete a profile
gixy profile delete <profile-name>
```

### Commit message generation

Generate a commit message from your staged changes using an AI provider.

**1. Add a provider (one-time setup)**

```sh
gixy commit config add <name> --provider gemini --model gemini-2.0-flash --api-key <key>
```

| Flag         | Description                                |
| ------------ | ------------------------------------------ |
| `--provider` | AI provider. Currently supported: `gemini` |
| `--model`    | Model name (e.g. `gemini-2.0-flash`)       |
| `--api-key`  | API key for the provider                   |

The first provider added is automatically set as active.

**2. Stage your changes and generate**

```sh
git add <files>
gixy commit generate
```

gixy reads the staged diff, calls the active provider, and displays a suggested commit message. You'll be prompted to confirm before the commit is made:

```
Suggested commit message:
feat(auth): add OAuth2 login support

Use this message? [y/N]:
```

If you confirm with `y`, gixy runs `git commit -m <message>` for you.

**Managing providers**

```sh
# List all configured providers (* marks the active one)
gixy commit config list

# Switch the active provider
gixy commit config use <name>

# Remove a provider (alias: delete)
gixy commit config remove <name>
```

## Config storage

All gixy data is stored in a single JSON file:

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
  },
  "commit_gen": {
    "active": "personal-gemini",
    "providers": {
      "personal-gemini": {
        "provider": "gemini",
        "model": "gemini-2.0-flash",
        "api_key": "AIza..."
      }
    }
  }
}
```

> **Note:** API keys are stored in plaintext in this file. Make sure the file permissions are restrictive (gixy sets `0600` on write).

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
