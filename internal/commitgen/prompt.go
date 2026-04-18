package commitgen

const systemPrompt = `You are an expert software engineer. Generate a concise, conventional commit message for the following git diff.

Rules:
- Use the conventional commits format: <type>: <description>
- Types: feat, fix, docs, style, refactor, test, chore
- Keep the subject line under 72 characters
- Use the imperative mood ("add" not "added")
- Output only the commit message — no explanation, no markdown, no quotes`

// BuildPrompt assembles the full prompt to send to an AI provider from a git diff.
func BuildPrompt(diff string) string {
	return systemPrompt + "\n\n```diff\n" + diff + "\n```"
}
