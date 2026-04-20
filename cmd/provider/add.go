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

var AddCommand = cli.Command{
	Name:  "add",
	Usage: "Interactively add an AI provider for commit generation",
	Action: func(_ context.Context, _ *cli.Command) error {
		return runAdd()
	},
}

func runAdd() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println()
	fmt.Println("Select a provider:")
	fmt.Println("  1) Gemini")
	fmt.Println("  2) OpenAI")
	fmt.Println("  3) Anthropic")
	fmt.Println("  4) Ollama")
	fmt.Println()

	choiceStr, err := prompt(reader, "Choice", "1")
	if err != nil {
		return err
	}

	var entry config.CommitGenEntry

	switch choiceStr {
	case "1", "":
		entry, err = setupGemini(reader)
	case "2":
		entry, err = setupOpenAI(reader)
	case "3":
		entry, err = setupAnthropic(reader)
	case "4":
		entry, err = setupOllama(reader)
	default:
		return cli.Exit(colors.Red(fmt.Sprintf("invalid choice %q; enter 1, 2, 3, or 4", choiceStr)), 1)
	}
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

func setupGemini(r *bufio.Reader) (config.CommitGenEntry, error) {
	model, err := prompt(r, "Model", "gemini-2.0-flash")
	if err != nil {
		return config.CommitGenEntry{}, err
	}
	apiKey, err := prompt(r, "API key", "")
	if err != nil {
		return config.CommitGenEntry{}, err
	}
	if apiKey == "" {
		return config.CommitGenEntry{}, cli.Exit(colors.Red("API key cannot be empty"), 1)
	}
	return config.CommitGenEntry{Provider: "gemini", Model: model, APIKey: apiKey}, nil
}

func setupOpenAI(r *bufio.Reader) (config.CommitGenEntry, error) {
	model, err := prompt(r, "Model", "gpt-4o")
	if err != nil {
		return config.CommitGenEntry{}, err
	}
	apiKey, err := prompt(r, "API key", "")
	if err != nil {
		return config.CommitGenEntry{}, err
	}
	if apiKey == "" {
		return config.CommitGenEntry{}, cli.Exit(colors.Red("API key cannot be empty"), 1)
	}
	return config.CommitGenEntry{Provider: "openai", Model: model, APIKey: apiKey}, nil
}

func setupAnthropic(r *bufio.Reader) (config.CommitGenEntry, error) {
	model, err := prompt(r, "Model", "claude-3-5-sonnet-20241022")
	if err != nil {
		return config.CommitGenEntry{}, err
	}
	apiKey, err := prompt(r, "API key", "")
	if err != nil {
		return config.CommitGenEntry{}, err
	}
	if apiKey == "" {
		return config.CommitGenEntry{}, cli.Exit(colors.Red("API key cannot be empty"), 1)
	}
	return config.CommitGenEntry{Provider: "anthropic", Model: model, APIKey: apiKey}, nil
}

func setupOllama(r *bufio.Reader) (config.CommitGenEntry, error) {
	model, err := prompt(r, "Model", "llama3.2")
	if err != nil {
		return config.CommitGenEntry{}, err
	}
	if model == "" {
		return config.CommitGenEntry{}, cli.Exit(colors.Red("model cannot be empty"), 1)
	}
	host, err := prompt(r, "Host", "http://localhost:11434")
	if err != nil {
		return config.CommitGenEntry{}, err
	}
	return config.CommitGenEntry{Provider: "ollama", Model: model, Host: host}, nil
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
