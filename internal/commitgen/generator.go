package commitgen

import (
	"context"
	"fmt"

	"github.com/pawanprjl/gixy/internal/commitgen/provider/gemini"
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
}

// NewProvider returns a Provider implementation for the given ProviderConfig.
// Currently supports "gemini". Returns an error for unknown providers.
func NewProvider(cfg ProviderConfig) (Provider, error) {
	switch cfg.Provider {
	case "gemini":
		return gemini.New(cfg.Model, cfg.APIKey)
	default:
		return nil, fmt.Errorf("unsupported provider %q; supported: gemini", cfg.Provider)
	}
}

// GenerateCommitMessage builds the prompt from the staged diff and calls the provider.
// This is the single entry point for commit message generation.
func GenerateCommitMessage(ctx context.Context, diff string, p Provider) (string, error) {
	prompt := BuildPrompt(diff)
	msg, err := p.Generate(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("generate commit message: %w", err)
	}
	return msg, nil
}
