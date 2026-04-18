package gemini

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/genai"
)

// Provider implements commitgen.Provider using the Google Gemini API.
type Provider struct {
	client       *genai.Client
	model        string
	systemPrompt string
}

// New creates a new Gemini Provider. Returns an error if the API key or model is empty.
func New(model, apiKey, systemPrompt string) (*Provider, error) {
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

	return &Provider{client: client, model: model, systemPrompt: systemPrompt}, nil
}

// Generate sends the prompt to Gemini and returns the generated text.
func (p *Provider) Generate(ctx context.Context, prompt string) (string, error) {
	config := &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{{Text: p.systemPrompt}},
		},
	}

	result, err := p.client.Models.GenerateContent(ctx, p.model, genai.Text(prompt), config)
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
