package subcmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

// A Command specifies a sub-command for a program's command-line interface.
type Command struct {
	Name        string
	Description string
	Do          func(args []string)
	SubCommands []Command
}

// Run parses os.Args and dispatches to the correct subcommand given by cmds.
// It produces an error message listing the commands with their descriptions if
// a nonexistent subcommand is provided, or if the command is "help", "-h",
// "-help", or "--help". This error message may be customized by altering Usage.
//
// Run panics if any command is named "help", "-h", "-help", or "--help",
// or if any two commands have the same name.
func Run(cmds []Command) {
	run(cmds, 1, os.Args)
}

// checkingIndex is the index of os.Args to dispatch to the correct subcommand given
// by cmds.
func run(cmds []Command, checkingIndex int, args []string) {
	byName := make(map[string]func([]string))
	subCmdByName := make(map[string][]Command)
	for _, cmd := range cmds {
		if _, ok := helpWords[cmd.Name]; ok {
			panicf("subcmd: cannot name a command %q", cmd.Name)
		}
		if _, ok := byName[cmd.Name]; ok {
			panicf("subcmd: duplicate command %q given to Run", cmd.Name)
		}
		if (cmd.Do != nil) == (cmd.SubCommands != nil) {
			panicf("subcmd: need to assign either Do or SubCommands in command %q", cmd.Name)
		}
		if cmd.Do != nil {
			byName[cmd.Name] = cmd.Do
		} else {
			subCmdByName[cmd.Name] = cmd.SubCommands
		}
	}
	if len(args) < 2 {
		usageExit(cmds, checkingIndex, 1)
	}
	if _, ok := helpWords[args[1]]; ok {
		usageExit(cmds, checkingIndex, 0)
	}
	subCmds, ok := subCmdByName[args[1]]
	if ok {
		run(subCmds, checkingIndex+1, args[1:])
		return
	}
	do, ok := byName[args[1]]
	if ok {
		do(args[2:])
		return
	}
	usageExit(cmds, checkingIndex, 1)
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

func usageExit(cmds []Command, checkingIndex, status int) {
	Usage(cmds, checkingIndex)
	os.Exit(status)
}

// Usage prints a help message listing the possible commands.
var Usage = func(cmds []Command, checkingIndex int) {
	parsedArgs := strings.Join(os.Args[:checkingIndex], " ")
	fmt.Fprintf(os.Stderr, "Usage:\n\n  %s COMMAND\n\nPossible commands are:\n\n", parsedArgs)
	PrintDefaults(cmds)
	fmt.Fprintf(
		os.Stderr,
		"\nRun '%s COMMAND -h' to see more information about a command.\n",
		parsedArgs,
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
