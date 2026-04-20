# Changelog

All notable changes to gixy are documented here.

---

## [Unreleased]

### Added

#### Commit generation (`gixy commit generate`)
- **Interactive action prompt** — replaced the simple `[y/N]` confirmation with a full action menu: `[y]es / [e]dit / [r]egenerate / [N]o`
- **Regenerate** — press `r` to call the AI again with the same diff and get a fresh suggestion without restarting
- **Edit in `$EDITOR`** — press `e` to open the suggested message in your `$EDITOR` (falls back to `vi`) for manual tweaks; loops back to the action prompt after saving
- **`--context` flag** — pass extra free-text to guide the AI (e.g. `gixy commit generate --context "fixes login bug reported by QA"`); context is appended to the prompt
- **Optional description** — after confirming, prompts for a free-text description appended as the commit message body (press Enter to skip)
- **Optional issue link** — after confirming, prompts for an issue URL appended as `Issue linked: <url>` in the commit footer (press Enter to skip)
- **Git repo guard** — shows a clean error if run outside a git repository instead of a confusing git error
- **Empty diff guard** — shows a clear error if no changes are staged, prompting the user to `git add` first

#### Provider management (`gixy provider`)
- **`gixy provider add`** — interactive wizard to add an AI provider; shows provider menu first (Gemini / OpenAI / Anthropic / Ollama), then provider-specific fields with sensible defaults, then asks for a config name
- **`gixy provider list`** — lists all configured providers with the active one marked `*`; prompts to switch by number, or press Enter to keep current
- **`gixy provider remove <name>`** — removes a configured provider by name (alias: `delete`); warns if the removed provider was active and others remain

#### New AI providers
- **OpenAI** — uses the chat completions API (`gpt-4o` default); requires an API key
- **Anthropic** — uses the Messages API (`claude-3-5-sonnet-20241022` default); requires an API key
- **Ollama** — uses a local Ollama instance (`llama3.2` default, `http://localhost:11434`); no API key needed

#### Profile management (`gixy profile`)
- **`gixy profile show`** — shows the active git identity for the current repo; matches it against saved gixy profiles and highlights the profile name if found
- **`gixy profile edit <name>`** — edit an existing profile's name and email interactively; press Enter on any field to keep the current value

#### General
- **Version** — `gixy --version` / `gixy -v` now prints the version (`0.1.0`)

### Changed
- `gixy commit config add` now accepts `openai`, `anthropic`, and `ollama` in addition to `gemini` for the `--provider` flag
- Config schema: `CommitGenEntry` gained an optional `host` field used by the Ollama provider
- `commitgen.GenerateCommitMessage` and `commitgen.BuildPrompt` now accept an `extraContext` string parameter (empty string = same behaviour as before)
