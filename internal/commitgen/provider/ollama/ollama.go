package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const defaultHost = "http://localhost:11434"

// Provider implements commitgen.Provider using a local Ollama instance.
type Provider struct {
	host         string
	model        string
	systemPrompt string
}

// New creates a new Ollama Provider. Host defaults to http://localhost:11434 if empty.
func New(model, host, systemPrompt string) (*Provider, error) {
	if model == "" {
		return nil, fmt.Errorf("model must not be empty")
	}
	if host == "" {
		host = defaultHost
	}
	return &Provider{model: model, host: strings.TrimRight(host, "/"), systemPrompt: systemPrompt}, nil
}

type chatRequest struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
	Error string `json:"error,omitempty"`
}

// Generate sends the prompt to Ollama and returns the generated text.
func (p *Provider) Generate(ctx context.Context, prompt string) (string, error) {
	body, err := json.Marshal(chatRequest{
		Model: p.model,
		Messages: []message{
			{Role: "system", Content: p.systemPrompt},
			{Role: "user", Content: prompt},
		},
		Stream: false,
	})
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.host+"/api/chat", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("ollama request: %w", err)
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
		if result.Error != "" {
			errMsg = result.Error
		}
		return "", fmt.Errorf("ollama API error %d: %s", resp.StatusCode, errMsg)
	}

	return strings.TrimSpace(result.Message.Content), nil
}
