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

// resolveProfile prints shell env injecting the cwd's profile per git invocation; fails soft.
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

	_, profileName := resolveProfileName(cwd, cfg)
	if profileName == "" {
		return nil
	}

	p, exists := cfg.Profiles[profileName]
	if !exists {
		return nil
	}

	shell := cmd.String("shell")
	var b strings.Builder

	// Skip identity if the repo sets a local user.email (let local config win).
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

// resolveProfileName returns the profile for cwd (longest mapping prefix, else default); matchedPath is empty when defaulted.
func resolveProfileName(cwd string, cfg *config.Config) (matchedPath, name string) {
	path, profile := longestPrefixMatch(cwd, cfg.PathMappings)
	if profile == "" {
		return "", cfg.DefaultProfile
	}
	return path, profile
}

// longestPrefixMatch returns the most-specific mapping prefix of cwd as (path, profile), or ("", "").
func longestPrefixMatch(cwd string, mappings map[string]string) (string, string) {
	bestLen := -1
	bestPath := ""
	bestProfile := ""

	for mappedPath, profileName := range mappings {
		clean := filepath.Clean(mappedPath)
		// cwd == mapping or inside it; trailing slash stops /work matching /workspace.
		if cwd == clean || strings.HasPrefix(cwd+string(filepath.Separator), clean+string(filepath.Separator)) {
			if len(clean) > bestLen {
				bestLen = len(clean)
				bestPath = clean
				bestProfile = profileName
			}
		}
	}

	return bestPath, bestProfile
}

// hasLocalConfig reports whether the repo has a non-empty local value for key (false outside a repo).
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

// shellSingleQuote wraps s in single quotes, escaping embedded quotes, as one safe token.
func shellSingleQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}
