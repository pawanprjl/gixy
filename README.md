# gixy

A CLI companion for git that adds workflow enhancements on top of your existing git setup.

`gixy` handles gixy-specific commands (like profile management). For everything else, use `git` directly — there is no proxy layer.

---

## Features

- **Profile management** — define named author profiles (name + email), switch them globally in one command
- **SSH key management** — each profile gets its own ed25519 keypair; `profile use` symlinks the active keys to `~/.ssh/id_ed25519`
- **AI commit messages** — generate conventional commit messages from staged changes using an AI provider of your choice
- **Multiple AI providers** — supports Gemini, OpenAI, Anthropic, and Ollama (local)
- **Interactive commit flow** — accept, edit in `$EDITOR`, or regenerate a suggestion before committing
- **Config stored at** `~/.config/gixy/config`
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

### Commit message generation

Generate a conventional commit message from your staged changes.

**1. Add a provider (one-time setup)**

```sh
gixy provider add
```

This runs an interactive wizard that walks you through selecting a provider, entering a model name, and providing an API key (or host for Ollama). The first provider added is automatically set as active.

**2. Stage your changes and generate**

```sh
git add <files>
gixy commit generate
```

gixy reads the staged diff, calls the active provider, and displays a suggested commit message with an interactive action menu:

```
Suggested commit message:
feat(auth): add OAuth2 login support

[y] accept  [e] edit  [r] regenerate  [N] abort
```

| Key | Action                                             |
| --- | -------------------------------------------------- |
| `y` | Accept the message and run `git commit`            |
| `e` | Open the message in `$EDITOR` (falls back to `vi`) |
| `r` | Call the AI again for a new suggestion             |
| `N` | Abort without committing                           |

**Optional flags**

| Flag               | Description                                              |
| ------------------ | -------------------------------------------------------- |
| `--context <text>` | Extra context to guide the AI (e.g. ticket title, scope) |
| `--issue <url>`    | Appends the URL as an issue footer in the commit message |

### Provider management

```sh
# Interactive wizard to add a new provider
gixy provider add

# List all configured providers (* marks the active one)
gixy provider list

# Switch the active provider
gixy provider use <name>

# Remove a provider (alias: delete)
gixy provider remove <name>
```

### Supported AI providers

| Provider    | Default model                | Auth                            |
| ----------- | ---------------------------- | ------------------------------- |
| `gemini`    | `gemini-2.0-flash`           | API key                         |
| `openai`    | `gpt-4o`                     | API key                         |
| `anthropic` | `claude-3-5-sonnet-20241022` | API key                         |
| `ollama`    | `llama3.2`                   | None (local); configurable host |

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
  "commit_gen": {
    "active": "personal-gemini",
    "providers": {
      "personal-gemini": {
        "provider": "gemini",
        "model": "gemini-2.0-flash",
        "api_key": "AIza..."
      },
      "local-llama": {
        "provider": "ollama",
        "model": "llama3.2",
        "host": "http://localhost:11434"
      }
    }
  }
}
```

> **Note:** API keys are stored in plaintext in this file. gixy sets file permissions to `0600` on write.

SSH keypairs are stored separately under `~/.ssh/gixy/<profile-name>/id_ed25519{,.pub}`.

---

## Roadmap

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
- [urfave/cli v3](https://github.com/urfave/cli) — CLI framework

---

## License

MIT — see [LICENSE](LICENSE)
