package subcmd

import (
	"fmt"
	"os"
	"text/tabwriter"
)

// A Command specifies a sub-command for a program's command-line interface.
type Command struct {
	Name        string
	Description string
	Do          func(args []string)
}

// Run parses os.Args and dispatches to the correct subcommand given by cmds.
// It produces an error message listing the commands with their descriptions if
// a nonexistent subcommand is provided, or if the command is "help", "-h",
// "-help", or "--help". This error message may be customized by altering Usage.
//
// Run panics if any command is named "help", "-h", "-help", or "--help",
// or if any two commands have the same name.
func Run(cmds []Command) {
	byName := make(map[string]func([]string))
	for _, cmd := range cmds {
		if _, ok := helpWords[cmd.Name]; ok {
			panicf("subcmd: cannot name a command %q", cmd.Name)
		}
		if _, ok := byName[cmd.Name]; ok {
			panicf("subcmd: duplicate command %q given to Run", cmd.Name)
		}
		byName[cmd.Name] = cmd.Do
	}
	if len(os.Args) < 2 {
		usageExit(cmds, 1)
	}
	if _, ok := helpWords[os.Args[1]]; ok {
		usageExit(cmds, 0)
	}
	do, ok := byName[os.Args[1]]
	if !ok {
		usageExit(cmds, 1)
	}
	do(os.Args[2:])
}

func panicf(format string, args ...interface{}) {
	panic(fmt.Sprintf(format, args...))
}

var helpWords = map[string]struct{}{
	"help":   {},
	"-h":     {},
	"-help":  {},
	"--help": {},
}

func usageExit(cmds []Command, status int) {
	Usage(cmds)
	os.Exit(status)
}

// Usage prints a help message listing the possible commands.
var Usage = func(cmds []Command) {
	fmt.Fprintf(os.Stderr, "Usage:\n\n  %s COMMAND\n\nPossible commands are:\n\n", os.Args[0])
	PrintDefaults(cmds)
	fmt.Fprintf(
		os.Stderr,
		"\nRun '%s COMMAND -h' to see more information about a command.\n",
		os.Args[0],
	)
}

// PrintDefaults formats a list of commands. For each command, the output is
//   Name    Description
func PrintDefaults(cmds []Command) {
	tw := tabwriter.NewWriter(os.Stderr, 0, 0, 4, ' ', 0)
	for _, cmd := range cmds {
		fmt.Fprintf(tw, "  %s\t%s\n", cmd.Name, cmd.Description)
	}
	tw.Flush()
}
