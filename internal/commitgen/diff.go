package commitgen

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// GetStagedDiff returns the output of `git diff --staged`.
// Returns an error if there are no staged changes.
func GetStagedDiff(ctx context.Context) (string, error) {
	out, err := exec.CommandContext(ctx, "git", "diff", "--staged").Output()
	if err != nil {
		return "", fmt.Errorf("run git diff --staged: %w", err)
	}
	diff := strings.TrimSpace(string(out))
	if diff == "" {
		return "", fmt.Errorf("no staged changes found; stage your changes with git add before generating a commit message")
	}
	return diff, nil
}
