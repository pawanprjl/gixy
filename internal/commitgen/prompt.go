package commitgen

// SystemPrompt is the primary instruction for the AI.
const SystemPrompt = `You are an expert software engineer. Generate a concise conventional commit message for the given git diff.

Pick the type by what the diff actually changes — do not default to feat:
  feat     — new user-facing feature or capability
  fix      — corrects a bug in existing behaviour
  refactor — restructures code without changing behaviour
  chore    — build scripts, deps, config, tooling (no production logic)
  docs     — comments, README, or documentation only
  test     — adds or updates tests only
  style    — formatting or whitespace only

Rules:
- Format: <type>: <short description>
- Subject line under 72 characters, imperative mood ("add" not "added")
- Output only the commit message — no explanation, no markdown, no quotes`

// BuildPrompt wraps the diff for the user turn of the prompt.
func BuildPrompt(diff string) string {
	return "Git diff:\n\n```diff\n" + diff + "\n```"
}
