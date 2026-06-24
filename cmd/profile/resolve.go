package profile

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pawanprjl/gixy/internal/config"
	"github.com/pawanprjl/gixy/internal/sshutil"
	"github.com/urfave/cli/v3"
)

var ResolveCommand = cli.Command{
	Name:   "resolve",
	Hidden: true,
	Usage:  "Print shell env for the current directory's profile (used by shell hooks)",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "shell",
			Value: "posix",
			Usage: "Output format: posix or fish",
		},
	},
	Action: resolveProfile,
}

// resolveProfile prints shell-evalable env that injects the cwd's profile identity
// and SSH key into a single git invocation, without mutating any global state.
// It fails soft (prints nothing, exits 0) on any error since it runs on every
// directory change.
func resolveProfile(_ context.Context, cmd *cli.Command) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil
	}
	cwd = filepath.Clean(cwd)

	profileName := longestPrefixMatch(cwd, cfg.PathMappings)
	if profileName == "" {
		profileName = cfg.DefaultProfile
	}
	if profileName == "" {
		return nil
	}

	p, exists := cfg.Profiles[profileName]
	if !exists {
		return nil
	}

	shell := cmd.String("shell")
	var b strings.Builder

	// Identity — skip if the repo sets an explicit local user.email, so git's
	// normal precedence (local config wins) is preserved.
	if !hasLocalConfig("user.email") {
		writeEnv(&b, shell, "GIT_AUTHOR_NAME", p.Name)
		writeEnv(&b, shell, "GIT_AUTHOR_EMAIL", p.Email)
		writeEnv(&b, shell, "GIT_COMMITTER_NAME", p.Name)
		writeEnv(&b, shell, "GIT_COMMITTER_EMAIL", p.Email)
	}

	// SSH key — skip if the repo sets an explicit local core.sshCommand.
	if !hasLocalConfig("core.sshCommand") {
		if keyDir, err := sshutil.KeyDir(profileName); err == nil {
			keyPath := filepath.Join(keyDir, "id_ed25519")
			if _, err := os.Stat(keyPath); err == nil {
				writeEnv(&b, shell, "GIT_SSH_COMMAND",
					fmt.Sprintf(`ssh -i "%s" -o IdentitiesOnly=yes`, keyPath))
			}
		}
	}

	fmt.Print(b.String())
	return nil
}

// hasLocalConfig reports whether the current repo has a non-empty local value for
// the given git config key. Outside a repo (or when unset) it returns false. git's
// stderr is discarded (Stderr stays nil → /dev/null), so nothing leaks to the hook.
func hasLocalConfig(key string) bool {
	out, err := exec.Command("git", "config", "--local", "--get", key).Output()
	return err == nil && strings.TrimSpace(string(out)) != ""
}

// writeEnv appends one shell-safe export/set line for the given variable.
func writeEnv(b *strings.Builder, shell, name, value string) {
	quoted := shellSingleQuote(value)
	if shell == "fish" {
		fmt.Fprintf(b, "set -lx %s %s\n", name, quoted)
		return
	}
	fmt.Fprintf(b, "export %s=%s\n", name, quoted)
}

// shellSingleQuote wraps s in single quotes, escaping embedded single quotes so the
// result is a single safe token in both POSIX shells and fish.
func shellSingleQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}
