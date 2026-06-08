package commitgen

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

const maxDiffBytes = 50_000

// DiffResult holds the diff content and metadata about how it was retrieved.
type DiffResult struct {
	Content   string
	Truncated bool // true when the diff was cut to fit within maxDiffBytes
}

// GetStagedDiff returns the staged diff. If the diff exceeds maxDiffBytes it is
// truncated at a line boundary and prefixed with a --stat summary so the AI
// still sees the full scope of changes alongside as much detail as fits.
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
		return DiffResult{Content: diff}, nil
	}

	// Diff exceeds the limit — truncate at a line boundary and prepend the stat
	// summary so the AI sees the full scope of changes plus as much detail as fits.
	statOut, err := exec.CommandContext(ctx, "git", "diff", "--staged", "--stat").Output()
	if err != nil {
		// Stat failed too; return what we have, truncated.
		return DiffResult{Content: truncateAtLine(diff, maxDiffBytes), Truncated: true}, nil
	}
	stat := strings.TrimSpace(string(statOut))
	truncated := truncateAtLine(diff, maxDiffBytes)
	return DiffResult{Content: stat + "\n\n" + truncated, Truncated: true}, nil
}

// truncateAtLine cuts s to at most maxBytes, trimming back to the last newline
// so the result never ends mid-line.
func truncateAtLine(s string, maxBytes int) string {
	if len(s) <= maxBytes {
		return s
	}
	cut := s[:maxBytes]
	if idx := strings.LastIndexByte(cut, '\n'); idx != -1 {
		cut = cut[:idx]
	}
	return cut
}
