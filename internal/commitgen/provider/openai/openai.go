package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Provider implements commitgen.Provider using the OpenAI chat completions API.
type Provider struct {
	apiKey       string
	model        string
	systemPrompt string
}

// New creates a new OpenAI Provider.
func New(model, apiKey, systemPrompt string) (*Provider, error) {
	if model == "" {
		return nil, fmt.Errorf("model must not be empty")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("api key must not be empty")
	}
	return &Provider{model: model, apiKey: apiKey, systemPrompt: systemPrompt}, nil
}

type chatRequest struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// Generate sends the prompt to OpenAI and returns the generated text.
func (p *Provider) Generate(ctx context.Context, prompt string) (string, error) {
	body, err := json.Marshal(chatRequest{
		Model: p.model,
		Messages: []message{
			{Role: "system", Content: p.systemPrompt},
			{Role: "user", Content: prompt},
		},
	})
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.openai.com/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("openai request: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	var result chatResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return "", fmt.Errorf("parse response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		errMsg := strings.TrimSpace(string(data))
		if result.Error != nil {
			errMsg = result.Error.Message
		}
		return "", fmt.Errorf("openai API error %d: %s", resp.StatusCode, errMsg)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("openai returned no choices")
	}
	return strings.TrimSpace(result.Choices[0].Message.Content), nil
}
