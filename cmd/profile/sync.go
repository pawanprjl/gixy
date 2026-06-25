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

// gixyProfileKey is the local-config marker recording which profile gixy pinned
// into a repo. Its presence means gixy owns the repo's identity/SSH settings.
const gixyProfileKey = "gixy.profile"

var SyncCommand = cli.Command{
	Name:   "sync",
	Hidden: true,
	Usage:  "Sync the current repo's local git config to its mapped profile (used by shell hook)",
	Action: syncProfile,
}

// syncProfile pins the cwd's repo to its mapped profile via local git config.
// It is fail-soft: any problem returns nil so the shell hook never blocks git.
func syncProfile(_ context.Context, _ *cli.Command) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil
	}
	cwd = filepath.Clean(cwd)

	if !insideGitRepo() {
		return nil
	}

	_, profileName := resolveProfileName(cwd, cfg)
	if profileName == "" {
		return nil
	}
	p, ok := cfg.Profiles[profileName]
	if !ok {
		return nil
	}

	switch marker := localConfigGet(gixyProfileKey); {
	case marker == profileName:
		// Already managed and current — nothing to do.
		return nil
	case marker == "" && hasLocalConfig("user.email"):
		// Unmanaged repo with a hand-set identity — respect it, never clobber.
		return nil
	default:
		// Fresh repo, or the mapped profile changed — (re)pin it.
		_ = writeRepoProfile(profileName, p)
		return nil
	}
}

// writeRepoProfile writes the profile's identity, SSH key, and the gixy marker
// into the current repo's local git config.
func writeRepoProfile(profileName string, p config.Profile) error {
	if err := setLocalConfig("user.name", p.Name); err != nil {
		return err
	}
	if err := setLocalConfig("user.email", p.Email); err != nil {
		return err
	}

	// Clear any sshCommand we set for a previous profile, then set the new
	// profile's key — so re-pinning to a keyless profile doesn't leave a stale one.
	unsetLocalConfig("core.sshCommand")
	if keyDir, err := sshutil.KeyDir(profileName); err == nil {
		keyPath := filepath.Join(keyDir, "id_ed25519")
		if _, err := os.Stat(keyPath); err == nil {
			if err := setLocalConfig("core.sshCommand",
				fmt.Sprintf(`ssh -i "%s" -o IdentitiesOnly=yes`, keyPath)); err != nil {
				return err
			}
		}
	}

	return setLocalConfig(gixyProfileKey, profileName)
}

// insideGitRepo reports whether the cwd is inside a git work tree.
func insideGitRepo() bool {
	out, err := exec.Command("git", "rev-parse", "--is-inside-work-tree").Output()
	return err == nil && strings.TrimSpace(string(out)) == "true"
}

// hasLocalConfig reports whether the repo has a non-empty local value for key.
func hasLocalConfig(key string) bool {
	return localConfigGet(key) != ""
}

// localConfigGet returns the repo-local value for key, or "" if unset/outside a repo.
func localConfigGet(key string) string {
	out, err := exec.Command("git", "config", "--local", "--get", key).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// setLocalConfig sets a repo-local git config value.
func setLocalConfig(key, value string) error {
	return exec.Command("git", "config", "--local", key, value).Run()
}

// unsetLocalConfig removes a repo-local git config key; a missing key is not an error.
func unsetLocalConfig(key string) {
	_ = exec.Command("git", "config", "--local", "--unset", key).Run()
}
