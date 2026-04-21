package anthropic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const apiVersion = "2023-06-01"

// Provider implements commitgen.Provider using the Anthropic Messages API.
type Provider struct {
	apiKey       string
	model        string
	systemPrompt string
}

// New creates a new Anthropic Provider.
func New(model, apiKey, systemPrompt string) (*Provider, error) {
	if model == "" {
		return nil, fmt.Errorf("model must not be empty")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("api key must not be empty")
	}
	return &Provider{model: model, apiKey: apiKey, systemPrompt: systemPrompt}, nil
}

type messagesRequest struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	System    string    `json:"system"`
	Messages  []message `json:"messages"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type messagesResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// Generate sends the prompt to Anthropic and returns the generated text.
func (p *Provider) Generate(ctx context.Context, prompt string) (string, error) {
	body, err := json.Marshal(messagesRequest{
		Model:     p.model,
		MaxTokens: 1024,
		System:    p.systemPrompt,
		Messages: []message{
			{Role: "user", Content: prompt},
		},
	})
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.anthropic.com/v1/messages", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", p.apiKey)
	req.Header.Set("anthropic-version", apiVersion)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("anthropic request: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	var result messagesResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return "", fmt.Errorf("parse response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		errMsg := strings.TrimSpace(string(data))
		if result.Error != nil {
			errMsg = result.Error.Message
		}
		return "", fmt.Errorf("anthropic API error %d: %s", resp.StatusCode, errMsg)
	}

	if len(result.Content) == 0 {
		return "", fmt.Errorf("anthropic returned empty content")
	}
	return strings.TrimSpace(result.Content[0].Text), nil
}
