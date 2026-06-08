package commitgen

// SystemPrompt is the primary instruction for the AI.
const SystemPrompt = `You are an expert software engineer. Generate a conventional commit message for the given git diff.

Pick the type by what the diff actually changes — do not default to feat:
  feat     — new user-facing feature or capability
  fix      — corrects a bug in existing behaviour
  refactor — restructures code without changing behaviour
  chore    — build scripts, deps, config, tooling (no production logic)
  docs     — comments, README, or documentation only
  test     — adds or updates tests only
  style    — formatting or whitespace only

Format:
  <type>: <short subject>

  <optional body>

Rules:
- Subject line: under 72 characters, imperative mood ("add" not "added")
- Body: include only when the diff is non-trivial and the subject alone is insufficient to explain the intent or rationale. Write 1-3 sentences explaining *why* the change was made, not what changed. Separate from subject with a blank line.
- Omit the body for simple changes (single-line fixes, formatting, dependency bumps, trivial chores).
- Output only the commit message — no explanation, no markdown, no quotes`

// BuildPrompt wraps the diff for the user turn of the prompt.
// extraContext is optional free-text appended to guide the AI.
func BuildPrompt(diff DiffResult, extraContext string) string {
	var p string
	if diff.Truncated {
		p = "The staged diff is large and was truncated. The stat summary below shows all changed files; the diff that follows is the first portion of the actual changes:\n\n" + diff.Content + "\n\n[diff truncated]"
	} else {
		p = "Git diff:\n\n```diff\n" + diff.Content + "\n```"
	}
	if extraContext != "" {
		p += "\n\nAdditional context: " + extraContext
	}
	return p
}
