# agents.md — AI Agent Instructions for gixy

This file provides context for AI coding agents (e.g. GitHub Copilot) working in this repository.

---

## Project Overview

`gixy` is a purpose-built CLI companion for git, written in Go. It provides git-enhancing commands (like profile management) alongside the user's existing git workflow. gixy handles gixy-specific tasks; users continue using `git` directly for git operations. There is no proxy layer.

---

## Tech Stack

| Component       | Choice                                         |
| --------------- | ---------------------------------------------- |
| Language        | Go (1.21+)                                     |
| CLI framework   | [urfave/cli v3](https://github.com/urfave/cli) |
| Config storage  | JSON file at `~/.config/gixy/config`           |
| Config path     | Resolved via `os.UserHomeDir()` + `.config`    |

---

## Intended Project Structure

```
gixy/
├── main.go                    # Entrypoint: registers commands
├── cmd/
│   └── <feature>/
│       └── <subcommand>.go    # One file per subcommand
├── internal/
│   └── config/
│       └── store.go           # Read/write ~/.config/gixy/config (JSON)
├── go.mod
└── go.sum
```

Each feature gets its own subdirectory under `cmd/`. Internal packages are shared utilities — config I/O. `main.go` only wires things together.

---

## Key Architectural Patterns

### 1. Config store

All persistent data lives in a single JSON file at `~/.config/gixy/config` on all platforms. The store in `internal/config/store.go` is the shared read/write layer used by all commands:

- `LoadConfig()` — reads and parses the config file
- `SaveConfig(cfg)` — marshals and writes back

The config path is resolved via `os.UserHomeDir()` joined with `.config/gixy/config`, ensuring a consistent `~/.config` location on all platforms including macOS.

---

## How to Add a New Command

1. Create a new directory `cmd/<feature>/` for the feature group
2. Add one `.go` file per subcommand (e.g. `add.go`, `list.go`, `use.go`)
3. In each file, define a `cli.Command` struct with `Name`, `Usage`, and `Action`
4. Register the command (or subcommand group) in `main.go`
5. If the command reads/writes persistent data, go through `internal/config/store.go` — do not open the config file directly
6. Keep `main.go` thin — it only wires commands, it does not contain business logic

---

## Conventions

- **Scope discipline** — implement only what is explicitly listed in the agreed plan. Do not add extra commands, flags, helpers, or features beyond what was drafted, even if they seem useful or are mentioned elsewhere (e.g. README roadmap)
- **Error handling** — wrap errors with `fmt.Errorf("context: %w", err)`; surface to the user via `cli.Exit(err, 1)`
- **No global state** — pass config/store as function arguments, not package-level vars
- **No silent failures** — always return or log errors; never swallow them
- **urfave/cli v3** — use `context.Context, *cli.Command` action signatures (v3 changed the first arg from `*cli.Context` to `context.Context`); prefer `cli.Command` structs
- **Config directory creation** — `os.MkdirAll` the config dir on first write; don't assume it exists
- **File creation (VS Code)** — VS Code may auto-insert a `package <name>` line when creating a new `.go` file; when writing file contents via a tool, always verify the final file has exactly one `package` declaration at the top to avoid `syntax error: non-declaration statement outside function body` compile errors

---

## Out of Scope (for now)

- Writing directly to `~/.gitconfig` (global git config)
- SSH key management
- GPG signing key management
- GUI or TUI interfaces
