package commitgen

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

const maxDiffBytes = 12_000

// DiffResult holds the diff content and whether it was truncated to stat-only.
type DiffResult struct {
	Content string
	IsStat  bool
}

// GetStagedDiff returns the staged diff, falling back to --stat if the full
// diff exceeds maxDiffBytes to avoid hitting AI token limits.
func GetStagedDiff(ctx context.Context) (DiffResult, error) {
	out, err := exec.CommandContext(ctx, "git", "diff", "--staged").Output()
	if err != nil {
		return DiffResult{}, fmt.Errorf("run git diff --staged: %w", err)
	}
	diff := strings.TrimSpace(string(out))
	if diff == "" {
		return DiffResult{}, fmt.Errorf("no staged changes found; stage your changes with git add before generating a commit message")
	}

	if len(diff) <= maxDiffBytes {
		return DiffResult{Content: diff, IsStat: false}, nil
	}

	// Diff is too large — fall back to --stat.
	statOut, err := exec.CommandContext(ctx, "git", "diff", "--staged", "--stat").Output()
	if err != nil {
		return DiffResult{}, fmt.Errorf("run git diff --staged --stat: %w", err)
	}
	return DiffResult{Content: strings.TrimSpace(string(statOut)), IsStat: true}, nil
}
