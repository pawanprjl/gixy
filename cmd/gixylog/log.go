package gixylog

import (
"context"
"fmt"
"os/exec"
"strings"

"github.com/pawanprjl/gixy/internal/colors"
"github.com/pawanprjl/gixy/internal/gitutil"
"github.com/urfave/cli/v3"
)

var LogCommand = cli.Command{
Name:  "log",
Usage: "Show a pretty color-coded git log",
Flags: []cli.Flag{
&cli.IntFlag{
Name:  "count",
Usage: "Number of commits to show",
Value: 20,
},
},
Action: func(ctx context.Context, cmd *cli.Command) error {
return runLog(ctx, cmd.Int("count"))
},
}

var typeColors = map[string]func(string) string{
"feat":     colors.Green,
"fix":      colors.Red,
"docs":     colors.Cyan,
"chore":    colors.Yellow,
"test":     colors.Yellow,
"refactor": colors.Cyan,
"style":    colors.Yellow,
"perf":     colors.Green,
"ci":       colors.Yellow,
"build":    colors.Yellow,
"revert":   colors.Red,
}

func colorSubject(subject string) string {
for typ, colorFn := range typeColors {
if strings.HasPrefix(subject, typ+":") || strings.HasPrefix(subject, typ+"(") {
return colorFn(subject)
}
}
return subject
}

func runLog(ctx context.Context, count int) error {
if err := gitutil.EnsureGitRepo(); err != nil {
return cli.Exit(colors.Red(err.Error()), 1)
}

out, err := exec.CommandContext(ctx, "git", "log",
fmt.Sprintf("-n%d", count),
"--format=%C(dim)%h%Creset %s",
"--color=never",
).Output()
if err != nil {
return cli.Exit(fmt.Errorf("git log: %w", err), 1)
}

lines := strings.Split(strings.TrimSpace(string(out)), "\n")
if len(lines) == 0 || (len(lines) == 1 && lines[0] == "") {
fmt.Println(colors.Yellow("No commits found."))
return nil
}

fmt.Println()
for _, line := range lines {
if line == "" {
continue
}
parts := strings.SplitN(line, " ", 2)
if len(parts) != 2 {
fmt.Println(line)
continue
}
hash := "\033[2m" + parts[0] + "\033[0m"
subject := colorSubject(parts[1])
fmt.Printf("  %s  %s\n", hash, subject)
}
fmt.Println()
return nil
}
