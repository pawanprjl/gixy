package commitgen

import (
	"context"
	"fmt"

	"github.com/pawanprjl/gixy/internal/commitgen/provider/anthropic"
	"github.com/pawanprjl/gixy/internal/commitgen/provider/gemini"
	"github.com/pawanprjl/gixy/internal/commitgen/provider/ollama"
	"github.com/pawanprjl/gixy/internal/commitgen/provider/openai"
	"github.com/pawanprjl/gixy/internal/config"
)

// Provider is the interface that all AI provider implementations must satisfy.
// Implementations are thin adapters: they accept a fully-formed prompt and return raw text.
type Provider interface {
	Generate(ctx context.Context, prompt string) (string, error)
}

// ProviderConfig holds the configuration required to construct a Provider.
type ProviderConfig struct {
	Provider string
	Model    string
	APIKey   string
	Host     string // used by ollama
}

// NewProvider returns a Provider implementation for the given ProviderConfig.
// Currently supports "gemini", "openai", "ollama", and "anthropic".
func NewProvider(cfg ProviderConfig) (Provider, error) {
	switch cfg.Provider {
	case "gemini":
		return gemini.New(cfg.Model, cfg.APIKey, SystemPrompt)
	case "openai":
		return openai.New(cfg.Model, cfg.APIKey, SystemPrompt)
	case "ollama":
		return ollama.New(cfg.Model, cfg.Host, SystemPrompt)
	case "anthropic":
		return anthropic.New(cfg.Model, cfg.APIKey, SystemPrompt)
	default:
		return nil, fmt.Errorf("unsupported provider %q; supported: gemini, openai, ollama, anthropic", cfg.Provider)
	}
}

// NewProviderFromEntry constructs a Provider directly from a config.CommitGenEntry,
// avoiding the need for callers to manually map fields into ProviderConfig.
func NewProviderFromEntry(entry config.CommitGenEntry) (Provider, error) {
	return NewProvider(ProviderConfig{
		Provider: entry.Provider,
		Model:    entry.Model,
		APIKey:   entry.APIKey,
		Host:     entry.Host,
	})
}

// GenerateCommitMessage builds the prompt from the staged diff and calls the provider.
// extraContext is optional free-text to guide the AI (can be empty).
func GenerateCommitMessage(ctx context.Context, diff DiffResult, extraContext string, p Provider) (string, error) {
	prompt := BuildPrompt(diff.Content, extraContext, diff.IsStat)
	msg, err := p.Generate(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("generate commit message: %w", err)
	}
	return msg, nil
}
