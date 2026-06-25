package initcmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pawanprjl/gixy/internal/colors"
	"github.com/urfave/cli/v3"
)

var InitCommand = cli.Command{
	Name:      "init",
	Usage:     "Print shell integration hook for auto-activating profiles per git command",
	ArgsUsage: "[-]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "shell",
			Usage: "Shell type: zsh, bash, or fish (auto-detected if omitted)",
		},
	},
	Action: printInit,
}

func printInit(_ context.Context, cmd *cli.Command) error {
	shell := cmd.String("shell")
	rawMode := cmd.Args().First() == "-"

	// Auto-detect shell if not explicitly provided
	if shell == "" {
		shell = detectShell()
	}

	// Raw hook output: `gixy init -` or `gixy init --shell <shell>`
	if rawMode || cmd.IsSet("shell") {
		switch shell {
		case "zsh":
			fmt.Print(zshHook)
		case "bash":
			fmt.Print(bashHook)
		case "fish":
			fmt.Print(fishHook)
		default:
			return cli.Exit(colors.Red(fmt.Sprintf("unsupported shell %q — supported: zsh, bash, fish", shell)), 1)
		}
		return nil
	}

	// No args: show human-friendly setup instructions
	return printSetupInstructions(shell)
}

func printSetupInstructions(shell string) error {
	var configFile string

	switch shell {
	case "zsh":
		configFile = "~/.zshrc"
	case "bash":
		configFile = "~/.bashrc"
	case "fish":
		configFile = "~/.config/fish/config.fish"
	default:
		fmt.Println("To enable gixy shell integration, add the following to your shell config:")
		fmt.Println()
		fmt.Println("  " + colors.Bold(`eval "$(gixy init -)"`) + "   " + colors.Dim("# zsh / bash"))
		fmt.Println("  " + colors.Bold(`gixy init - | source`) + "    " + colors.Dim("# fish"))
		return nil
	}

	evalLine := `eval "$(gixy init -)"`
	if shell == "fish" {
		evalLine = `gixy init - | source`
	}

	fmt.Printf("To enable gixy shell integration, add the following to %s:\n", colors.Bold(configFile))
	fmt.Println()
	fmt.Printf("  %s\n", colors.Bold(evalLine))
	fmt.Println()
	fmt.Printf("Then restart your shell or run: %s\n", colors.Bold("source "+configFile))
	return nil
}

func detectShell() string {
	shellEnv := os.Getenv("SHELL")
	if shellEnv == "" {
		return ""
	}
	return filepath.Base(shellEnv)
}

// The hooks sync the cwd's repo to its mapped profile (writing the identity +
// SSH key into the repo's local git config) before running git. gixy is only
// spawned when $PWD changes since the last git call; once a repo is synced, git
// itself applies the profile for every tool, so no further gixy involvement is
// needed.
const zshHook = `# gixy shell integration
function git() {
  if [ "$PWD" != "$__gixy_pwd" ]; then
    command gixy profile sync
    __gixy_pwd="$PWD"
  fi
  command git "$@"
}
`

const bashHook = `# gixy shell integration
function git() {
  if [ "$PWD" != "$__gixy_pwd" ]; then
    command gixy profile sync
    __gixy_pwd="$PWD"
  fi
  command git "$@"
}
`

const fishHook = `# gixy shell integration
function git
  if test "$PWD" != "$__gixy_pwd"
    command gixy profile sync
    set -g __gixy_pwd "$PWD"
  end
  command git $argv
end
`
