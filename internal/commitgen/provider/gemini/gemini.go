package gemini

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/genai"
)

const systemPrompt = `You are an expert software engineer. Generate a concise, conventional commit message for the following git diff.

Rules:
- Use the conventional commits format: <type>: <description>
- Types: feat, fix, docs, style, refactor, test, chore
- Keep the subject line under 72 characters
- Use the imperative mood ("add" not "added")
- Output only the commit message — no explanation, no markdown, no quotes`

// Provider implements commitgen.Provider using the Google Gemini API.
type Provider struct {
	client *genai.Client
	model  string
}

// New creates a new Gemini Provider. Returns an error if the API key or model is empty.
func New(model, apiKey string) (*Provider, error) {
	if model == "" {
		return nil, fmt.Errorf("model must not be empty")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("api key must not be empty")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("create gemini client: %w", err)
	}

	return &Provider{client: client, model: model}, nil
}

// GenerateCommitMessage sends the diff to Gemini and returns a commit message.
func (p *Provider) GenerateCommitMessage(ctx context.Context, diff string) (string, error) {
	prompt := systemPrompt + "\n\n```diff\n" + diff + "\n```"

	result, err := p.client.Models.GenerateContent(ctx, p.model, genai.Text(prompt), nil)
	if err != nil {
		return "", fmt.Errorf("gemini generate content: %w", err)
	}

	if len(result.Candidates) == 0 {
		return "", fmt.Errorf("gemini returned no candidates")
	}
	candidate := result.Candidates[0]
	if candidate.Content == nil || len(candidate.Content.Parts) == 0 {
		return "", fmt.Errorf("gemini returned empty content")
	}

	var sb strings.Builder
	for _, part := range candidate.Content.Parts {
		if part != nil {
			sb.WriteString(part.Text)
		}
	}
	msg := strings.TrimSpace(sb.String())
	if msg == "" {
		return "", fmt.Errorf("gemini returned blank message")
	}
	return msg, nil
}
