package provider

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/pawanprjl/gixy/internal/colors"
	"github.com/pawanprjl/gixy/internal/config"
	"github.com/urfave/cli/v3"
)

type providerDef struct {
	name         string
	label        string
	defaultModel string
	needsAPIKey  bool
	defaultHost  string
}

var providerDefs = []providerDef{
	{name: "gemini", label: "Gemini", defaultModel: "gemini-2.0-flash", needsAPIKey: true},
	{name: "openai", label: "OpenAI", defaultModel: "gpt-4o", needsAPIKey: true},
	{name: "anthropic", label: "Anthropic", defaultModel: "claude-3-5-sonnet-20241022", needsAPIKey: true},
	{name: "ollama", label: "Ollama", defaultModel: "llama3.2", defaultHost: "http://localhost:11434"},
}

var AddCommand = cli.Command{
	Name:  "add",
	Usage: "Interactively add an AI provider for commit generation",
	Action: func(_ context.Context, _ *cli.Command) error {
		return runAdd()
	},
}

func runAdd() error {
	reader := bufio.NewReader(os.Stdin)

	names := make([]string, len(providerDefs))
	for i, p := range providerDefs {
		names[i] = p.name
	}

	fmt.Println()
	fmt.Printf("Available providers: %s\n", strings.Join(names, ", "))
	fmt.Println()

	providerName, err := prompt(reader, "Provider", providerDefs[0].name)
	if err != nil {
		return err
	}

	var def *providerDef
	for i := range providerDefs {
		if providerDefs[i].name == providerName {
			def = &providerDefs[i]
			break
		}
	}
	if def == nil {
		return cli.Exit(colors.Red(fmt.Sprintf("invalid provider %q; available: %s", providerName, strings.Join(names, ", "))), 1)
	}

	entry, err := setupProvider(reader, def)
	if err != nil {
		return err
	}

	name, err := prompt(reader, "Name for this provider config (e.g. my-gemini)", "")
	if err != nil {
		return err
	}
	if name == "" {
		return cli.Exit(colors.Red("name cannot be empty"), 1)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return cli.Exit(fmt.Errorf("load config: %w", err), 1)
	}

	if cfg.CommitGen == nil {
		cfg.CommitGen = &config.CommitGenConfig{
			Providers: make(map[string]config.CommitGenEntry),
		}
	}

	if _, exists := cfg.CommitGen.Providers[name]; exists {
		return cli.Exit(colors.Red(fmt.Sprintf("provider %q already exists; remove it first with `gixy commit config remove %s`", name, name)), 1)
	}

	cfg.CommitGen.Providers[name] = entry
	if cfg.CommitGen.Active == "" {
		cfg.CommitGen.Active = name
	}

	if err := config.SaveConfig(cfg); err != nil {
		return cli.Exit(fmt.Errorf("save config: %w", err), 1)
	}

	fmt.Println()
	fmt.Println(colors.Green(fmt.Sprintf("Provider %q added (%s / %s).", name, entry.Provider, entry.Model)))
	if cfg.CommitGen.Active == name {
		fmt.Println(colors.Cyan(fmt.Sprintf("Active provider set to %q.", name)))
	}
	return nil
}

func setupProvider(r *bufio.Reader, def *providerDef) (config.CommitGenEntry, error) {
	model, err := prompt(r, "Model", def.defaultModel)
	if err != nil {
		return config.CommitGenEntry{}, err
	}
	if model == "" {
		return config.CommitGenEntry{}, cli.Exit(colors.Red("model cannot be empty"), 1)
	}

	entry := config.CommitGenEntry{Provider: def.name, Model: model}

	if def.needsAPIKey {
		apiKey, err := prompt(r, "API key", "")
		if err != nil {
			return config.CommitGenEntry{}, err
		}
		if apiKey == "" {
			return config.CommitGenEntry{}, cli.Exit(colors.Red("API key cannot be empty"), 1)
		}
		entry.APIKey = apiKey
	} else {
		host, err := prompt(r, "Host", def.defaultHost)
		if err != nil {
			return config.CommitGenEntry{}, err
		}
		entry.Host = host
	}

	return entry, nil
}

// prompt prints "<label> [<default>]: " and reads a line.
// If the user presses Enter with no input, the default value is returned.
func prompt(r *bufio.Reader, label, defaultVal string) (string, error) {
	if defaultVal != "" {
		fmt.Printf("%s [%s]: ", colors.Cyan(label), defaultVal)
	} else {
		fmt.Printf("%s: ", colors.Cyan(label))
	}
	line, err := r.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("read input: %w", err)
	}
	line = strings.TrimSpace(line)
	if line == "" {
		return defaultVal, nil
	}
	return line, nil
}
