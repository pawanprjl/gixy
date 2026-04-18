# agents.md — AI Agent Instructions for gixy

This file provides context for AI coding agents (e.g. GitHub Copilot) working in this repository.

---

## Project Overview

`gixy` is a git proxy CLI tool written in Go. It wraps `git` with additional commands (like profile management) and transparently forwards anything it doesn't recognize directly to `git`. The goal is to progressively add git-enhancing capabilities without breaking any existing git workflows.

---

## Tech Stack

| Component       | Choice                                         |
| --------------- | ---------------------------------------------- |
| Language        | Go (1.21+)                                     |
| CLI framework   | [urfave/cli v2](https://github.com/urfave/cli) |
| Git passthrough | `os/exec` — exec the system `git` binary       |
| Config storage  | JSON file at `$XDG_CONFIG_HOME/gixy/config`    |
| Config path     | Resolved via `os.UserConfigDir()`              |

---

## Intended Project Structure

```
gixy/
├── main.go                    # Entrypoint: registers commands, sets up passthrough
├── cmd/
│   └── <feature>/
│       └── <subcommand>.go    # One file per subcommand
├── internal/
│   ├── config/
│   │   └── store.go           # Read/write ~/.config/gixy/config (JSON)
│   └── git/
│       └── passthrough.go     # exec.Command("git", args...) wrapper
├── go.mod
└── go.sum
```

Each feature gets its own subdirectory under `cmd/`. Internal packages are shared utilities — config I/O and git execution. `main.go` only wires things together.

---

## Key Architectural Patterns

### 1. Command routing

`urfave/cli` registers known commands. Any args that don't match a registered command are caught by `CommandNotFound` and forwarded to `git` via `os/exec`.

```go
// Pseudocode for passthrough default action
app.ExCommandNotFound = func(ctx *cli.Context, command string) {
    git.Passthrough(append([]string{command}, ctx.Args().Slice()...))
}
```

### 2. Config store

All persistent data lives in a single JSON file at `$XDG_CONFIG_HOME/gixy/config` (defaults to `~/.config/gixy/config`). The store in `internal/config/store.go` is the shared read/write layer used by all commands:

- `LoadConfig()` — reads and parses the config file
- `SaveConfig(cfg)` — marshals and writes back

The config path is resolved via `os.UserConfigDir()`, which handles XDG on Linux and `~/Library/Application Support` on macOS automatically.

### 3. Git passthrough

When gixy doesn't own a command, it delegates to the system `git` binary with `os.Stdin`, `os.Stdout`, and `os.Stderr` attached directly — so interactive git commands (e.g. `git rebase -i`) work exactly as expected.

---

## How to Add a New Command

1. Create a new directory `cmd/<feature>/` for the feature group
2. Add one `.go` file per subcommand (e.g. `add.go`, `list.go`, `use.go`)
3. In each file, define a `cli.Command` struct with `Name`, `Usage`, and `Action`
4. Register the command (or subcommand group) in `main.go`
5. If the command reads/writes persistent data, go through `internal/config/store.go` — do not open the config file directly
6. If the command needs to invoke git, use `internal/git/passthrough.go`
7. Keep `main.go` thin — it only wires commands, it does not contain business logic

---

## Conventions

- **Error handling** — wrap errors with `fmt.Errorf("context: %w", err)`; surface to the user via `cli.Exit(err, 1)`
- **No global state** — pass config/store as function arguments, not package-level vars
- **No silent failures** — always return or log errors; never swallow them
- **urfave/cli v2** — use `*cli.Context` action signatures; prefer `cli.Command` structs over the older fluent API
- **Config directory creation** — `os.MkdirAll` the config dir on first write; don't assume it exists
- **Passthrough fidelity** — when forwarding to git, attach `os.Stdin`, `os.Stdout`, `os.Stderr` to the exec'd process so interactive commands (e.g. `git rebase -i`) work correctly

---

## Out of Scope (for now)

- Writing directly to `~/.gitconfig` (global git config)
- SSH key management
- GPG signing key management
- GUI or TUI interfaces
