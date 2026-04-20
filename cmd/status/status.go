package status

import (
"context"
"fmt"
"os/exec"
"strings"

"github.com/pawanprjl/gixy/internal/colors"
"github.com/pawanprjl/gixy/internal/config"
"github.com/pawanprjl/gixy/internal/gitutil"
"github.com/urfave/cli/v3"
)

var StatusCommand = cli.Command{
Name:  "status",
Usage: "Show a compact summary of the working tree and active profile",
Action: func(ctx context.Context, _ *cli.Command) error {
return runStatus(ctx)
},
}

func runStatus(ctx context.Context) error {
if err := gitutil.EnsureGitRepo(); err != nil {
return cli.Exit(colors.Red(err.Error()), 1)
}

// Current branch
branchOut, err := exec.CommandContext(ctx, "git", "branch", "--show-current").Output()
if err != nil {
return cli.Exit(fmt.Errorf("get branch: %w", err), 1)
}
branch := strings.TrimSpace(string(branchOut))

// File counts from porcelain output
porcelainOut, err := exec.CommandContext(ctx, "git", "status", "--porcelain").Output()
if err != nil {
return cli.Exit(fmt.Errorf("git status: %w", err), 1)
}

var staged, unstaged, untracked int
for _, line := range strings.Split(string(porcelainOut), "\n") {
if len(line) < 2 {
continue
}
x, y := line[0], line[1]
if x == '?' && y == '?' {
untracked++
continue
}
if x != ' ' && x != '!' {
staged++
}
if y != ' ' && y != '!' {
unstaged++
}
}

// Active profile identity
nameOut, _ := exec.CommandContext(ctx, "git", "config", "user.name").Output()
emailOut, _ := exec.CommandContext(ctx, "git", "config", "user.email").Output()
gitName := strings.TrimSpace(string(nameOut))
gitEmail := strings.TrimSpace(string(emailOut))

// Match against saved profiles
cfg, _ := config.LoadConfig()
matchedProfile := ""
if cfg != nil {
for profileName, p := range cfg.Profiles {
if p.Name == gitName && p.Email == gitEmail {
matchedProfile = profileName
break
}
}
}

fmt.Println()
fmt.Printf("  %s  %s\n", colors.Cyan("branch"), colors.Green(branch))
fmt.Println()

stagedStr := fmt.Sprintf("%d", staged)
if staged > 0 {
stagedStr = colors.Green(stagedStr)
}
unstStr := fmt.Sprintf("%d", unstaged)
if unstaged > 0 {
unstStr = colors.Yellow(unstStr)
}
untrStr := fmt.Sprintf("%d", untracked)
if untracked > 0 {
untrStr = colors.Red(untrStr)
}
fmt.Printf("  staged %-4s  unstaged %-4s  untracked %s\n", stagedStr, unstStr, untrStr)

if gitName != "" {
fmt.Println()
if matchedProfile != "" {
fmt.Printf("  %s  %s (%s <%s>)\n", colors.Cyan("profile"), colors.Green(matchedProfile), gitName, gitEmail)
} else {
fmt.Printf("  %s  %s <%s>\n", colors.Cyan("profile"), gitName, gitEmail)
}
}
fmt.Println()
return nil
}
