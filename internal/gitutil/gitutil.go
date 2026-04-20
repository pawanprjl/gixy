package gitutil

import (
"fmt"
"os/exec"
)

// EnsureGitRepo returns a user-friendly error if the current directory is not
// inside a git repository.
func EnsureGitRepo() error {
cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
if err := cmd.Run(); err != nil {
return fmt.Errorf("not inside a git repository")
}
return nil
}
