package commitgen

import (
	"context"
	"fmt"

	"github.com/pawanprjl/gixy/internal/commitgen/provider/gemini"
)

// Provider is the interface that all AI provider implementations must satisfy.
// It is intentionally minimal so the package can be extracted as a standalone library.
type Provider interface {
	GenerateCommitMessage(ctx context.Context, diff string) (string, error)
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
