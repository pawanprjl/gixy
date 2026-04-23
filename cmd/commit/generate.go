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
	"github.com/pawanprjl/gixy/internal/gitutil"
	"github.com/urfave/cli/v3"
)

var GenerateCommand = cli.Command{
	Name:  "generate",
	Usage: "Generate a commit message from staged changes using AI",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "context",
			Usage: "Extra context to guide the AI (e.g. 'fixes login bug reported by QA')",
		},
		&cli.StringFlag{
			Name:  "issue",
			Usage: "Issue URL or identifier to append to the commit message (e.g. 'https://github.com/org/repo/issues/42')",
		},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		return runGenerate(ctx, cmd.String("context"), cmd.String("issue"))
	},
}

func runGenerate(ctx context.Context, extraContext, issueLink string) error {
	if err := gitutil.EnsureGitRepo(); err != nil {
		return cli.Exit(colors.Red(err.Error()), 1)
	}

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

	if diff.IsStat {
		fmt.Println(colors.Yellow("Diff is large; sending file summary to AI instead of full diff."))
	}

	provider, err := commitgen.NewProviderFromEntry(entry)
	if err != nil {
		return cli.Exit(fmt.Errorf("init provider: %w", err), 1)
	}

	if extraContext != "" {
		fmt.Println(colors.Cyan("Context: ") + extraContext)
	}
	fmt.Println(colors.Cyan("Generating commit message..."))
	message, err := commitgen.GenerateCommitMessage(ctx, diff, extraContext, provider)
	if err != nil {
		return cli.Exit(fmt.Errorf("generate message: %w", err), 1)
	}

	reader := bufio.NewReader(os.Stdin)

	message, err = negotiateMessage(ctx, diff, extraContext, provider, message, reader)
	if err != nil {
		return err
	}
	if message == "" {
		return nil
	}

	finalMessage := message
	if issueLink != "" {
		finalMessage += "\n\nIssue linked: " + issueLink
	}

	cmd := exec.CommandContext(ctx, "git", "commit", "-m", finalMessage)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return cli.Exit(fmt.Errorf("git commit: %w", err), 1)
	}
	return nil
}

// negotiateMessage presents the generated message to the user and loops until
// they accept, edit, or regenerate it. Returns an empty string if the user aborts.
func negotiateMessage(ctx context.Context, diff commitgen.DiffResult, extraContext string, provider commitgen.Provider, message string, reader *bufio.Reader) (string, error) {
	for {
		fmt.Println()
		fmt.Println(colors.Yellow("Suggested commit message:"))
		fmt.Println(message)
		fmt.Println()
		fmt.Print(colors.Cyan("Action? [y]es / [e]dit / [r]egenerate / [N]o: "))

		answer, err := reader.ReadString('\n')
		if err != nil {
			return "", cli.Exit(fmt.Errorf("read answer: %w", err), 1)
		}
		answer = strings.TrimSpace(strings.ToLower(answer))

		switch answer {
		case "r":
			fmt.Println(colors.Cyan("Regenerating..."))
			message, err = commitgen.GenerateCommitMessage(ctx, diff, extraContext, provider)
			if err != nil {
				return "", cli.Exit(fmt.Errorf("generate message: %w", err), 1)
			}
		case "e":
			edited, err := openInEditor(message)
			if err != nil {
				return "", cli.Exit(fmt.Errorf("edit message: %w", err), 1)
			}
			message = edited
		case "y":
			return message, nil
		default:
			fmt.Println("Aborted.")
			return "", nil
		}
	}
}

func openInEditor(content string) (string, error) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}

	tmp, err := os.CreateTemp("", "gixy-commit-*.txt")
	if err != nil {
		return "", fmt.Errorf("create temp file: %w", err)
	}
	defer os.Remove(tmp.Name())

	if _, err := tmp.WriteString(content); err != nil {
		tmp.Close()
		return "", fmt.Errorf("write temp file: %w", err)
	}
	tmp.Close()

	cmd := exec.Command(editor, tmp.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("editor: %w", err)
	}

	data, err := os.ReadFile(tmp.Name())
	if err != nil {
		return "", fmt.Errorf("read temp file: %w", err)
	}
	return strings.TrimSpace(string(data)), nil
}
