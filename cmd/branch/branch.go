package branch

import (
"bufio"
"context"
"fmt"
"os"
"os/exec"
"strconv"
"strings"

"github.com/pawanprjl/gixy/internal/colors"
"github.com/pawanprjl/gixy/internal/gitutil"
"github.com/urfave/cli/v3"
)

var BranchCommand = cli.Command{
Name:  "branch",
Usage: "Interactively list and switch recent git branches",
Action: func(ctx context.Context, _ *cli.Command) error {
return runBranch(ctx)
},
}

func runBranch(ctx context.Context) error {
if err := gitutil.EnsureGitRepo(); err != nil {
return cli.Exit(colors.Red(err.Error()), 1)
}

out, err := exec.CommandContext(ctx, "git", "branch", "--sort=-committerdate").Output()
if err != nil {
return cli.Exit(fmt.Errorf("list branches: %w", err), 1)
}

lines := strings.Split(strings.TrimSpace(string(out)), "\n")
if len(lines) == 0 || (len(lines) == 1 && lines[0] == "") {
fmt.Println(colors.Yellow("No local branches found."))
return nil
}

type branchEntry struct {
name    string
current bool
}

branches := make([]branchEntry, 0, len(lines))
for _, line := range lines {
if line == "" {
continue
}
current := strings.HasPrefix(line, "*")
name := strings.TrimSpace(strings.TrimPrefix(line, "*"))
branches = append(branches, branchEntry{name: name, current: current})
}

fmt.Println()
for i, b := range branches {
marker := "  "
nameStr := b.name
if b.current {
marker = colors.Green("*")
nameStr = colors.Green(b.name)
}
fmt.Printf("  %s %d) %s\n", marker, i+1, nameStr)
}
fmt.Println()

reader := bufio.NewReader(os.Stdin)
fmt.Printf("%s (Enter to cancel, or number to switch): ", colors.Cyan("Select branch"))
line, err := reader.ReadString('\n')
if err != nil {
return cli.Exit(fmt.Errorf("read input: %w", err), 1)
}
line = strings.TrimSpace(line)

if line == "" {
fmt.Println("No branch selected.")
return nil
}

idx, err := strconv.Atoi(line)
if err != nil || idx < 1 || idx > len(branches) {
return cli.Exit(colors.Red(fmt.Sprintf("invalid choice %q; enter a number between 1 and %d", line, len(branches))), 1)
}

chosen := branches[idx-1]
if chosen.current {
fmt.Println(colors.Cyan(fmt.Sprintf("Already on %q.", chosen.name)))
return nil
}

cmd := exec.CommandContext(ctx, "git", "checkout", chosen.name)
cmd.Stdout = os.Stdout
cmd.Stderr = os.Stderr
return cmd.Run()
}
