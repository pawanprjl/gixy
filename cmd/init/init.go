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
	Name:  "init",
	Usage: "Print shell integration hook for auto-activating profiles on cd",
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
	if shell == "" {
		shell = detectShell()
	}

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

func detectShell() string {
	shellEnv := os.Getenv("SHELL")
	if shellEnv == "" {
		return ""
	}
	return filepath.Base(shellEnv)
}

const zshHook = `# gixy shell integration
autoload -Uz add-zsh-hook
function _gixy_auto_activate() {
  gixy profile activate --silent
}
add-zsh-hook chpwd _gixy_auto_activate
_gixy_auto_activate
`

const bashHook = `# gixy shell integration
function _gixy_auto_activate() {
  gixy profile activate --silent
}
function cd() {
  builtin cd "$@" && _gixy_auto_activate
}
_gixy_auto_activate
`

const fishHook = `# gixy shell integration
function _gixy_auto_activate
  gixy profile activate --silent
end
function cd
  builtin cd $argv; and _gixy_auto_activate
end
_gixy_auto_activate
`
