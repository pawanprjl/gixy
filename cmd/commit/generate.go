package commit

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pawanprjl/gixy/internal/colors"
	"github.com/pawanprjl/gixy/internal/commitgen"
	"github.com/pawanprjl/gixy/internal/config"
	"github.com/urfave/cli/v3"
)

var GenerateCommand = cli.Command{
	Name:  "generate",
	Usage: "Generate a commit message from staged changes using AI",
	Action: func(ctx context.Context, _ *cli.Command) error {
		return runGenerate(ctx)
	},
}

func runGenerate(ctx context.Context) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return cli.Exit(fmt.Errorf("load config: %w", err), 1)
	}
	if cfg.CommitGen == nil || cfg.CommitGen.Active == "" {
		return cli.Exit(colors.Red("commit generation is not configured; run `gixy commit config add` first"), 1)
	}

	entry, exists := cfg.CommitGen.Providers[cfg.CommitGen.Active]
	if !exists {
		return cli.Exit(colors.Red(fmt.Sprintf("active provider %q not found; run `gixy commit config use` to set a valid provider", cfg.CommitGen.Active)), 1)
	}

	diff, err := commitgen.GetStagedDiff(ctx)
	if err != nil {
		return cli.Exit(err, 1)
	}

	fmt.Println(colors.Cyan("Generating commit message..."))

	provider, err := commitgen.NewProvider(commitgen.ProviderConfig{
		Provider: entry.Provider,
		Model:    entry.Model,
		APIKey:   entry.APIKey,
	})
	if err != nil {
		return cli.Exit(fmt.Errorf("init provider: %w", err), 1)
	}

	message, err := commitgen.GenerateCommitMessage(ctx, diff, provider)
	if err != nil {
		return cli.Exit(fmt.Errorf("generate message: %w", err), 1)
	}

	fmt.Println()
	fmt.Println(colors.Yellow("Suggested commit message:"))
	fmt.Println(message)
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print(colors.Cyan("Use this message? [y/N]: "))
	answer, err := reader.ReadString('\n')
	if err != nil {
		return cli.Exit(fmt.Errorf("read answer: %w", err), 1)
	}
	answer = strings.TrimSpace(strings.ToLower(answer))

	if answer != "y" {
		fmt.Println("Aborted.")
		return nil
	}

	cmd := exec.CommandContext(ctx, "git", "commit", "-m", message)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return cli.Exit(fmt.Errorf("git commit: %w", err), 1)
	}
	return nil
}
