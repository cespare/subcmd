// Package subcmd works with the flags package to implement sub-commands in the
// manner of git and similar tools.
package subcmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"
)

// A Command specifies a sub-command for a program's command-line interface.
type Command struct {
	Name        string              // the command's one-word name
	Description string              // a short description of the command
	Do          func(args []string) // command implementation
}

// A Runner runs sub-commands. To change Usage or ErrorHandling, alter these
// after creating a runner with New but before calling Runner.Run.
type Runner struct {
	cmds          []Command
	errorHandling flag.ErrorHandling

	// Usage prints the runner's usage.
	// If Usage is nil, the package-level Usage is called instead.
	Usage func()
}

// New creates a Runner with the given name and command list. The error-handling
// behavior of Run is controlled by errorHandling and has the same semantics as
// for flag.FlagSet.
//
// New panics if any command is named "help", "-h", "-help", or "--help",
// or if any two commands have the same name.
func New(name string, cmds []Command, errorHandling flag.ErrorHandling) *Runner {
	names := make(map[string]struct{})
	for _, cmd := range cmds {
		if _, ok := helpWords[cmd.Name]; ok {
			panicf("subcmd: cannot name a command %q", cmd.Name)
		}
		if _, ok := names[cmd.Name]; ok {
			panicf("subcmd: duplicate command %q given to Run", cmd.Name)
		}
		names[cmd.Name] = struct{}{}
	}
	return &Runner{
		cmds:          cmds,
		errorHandling: errorHandling,
		Usage:         func() { defaultUsage(name, cmds) },
	}
}

// ErrHelp is the error returned if the first argument is "help", "-h", "-help",
// or "--help".
var ErrHelp = errors.New("subcmd: help requested")

// Run parses args and dispatches to the correct subcommand.
// It produces an error message listing the commands with their descriptions if
// a nonexistent subcommand is provided, or if the command is "help", "-h",
// "-help", or "--help". This error message may be customized by altering
// r.Usage.
func (r *Runner) Run(args []string) error {
	if len(args) < 1 {
		return r.errorExit(args, errors.New("subcmd: no sub-command provided"))
	}
	if _, ok := helpWords[args[0]]; ok {
		return r.errorExit(args, ErrHelp)
	}
	for _, cmd := range r.cmds {
		if cmd.Name == args[0] {
			cmd.Do(args[1:])
			return nil
		}
	}
	err := fmt.Errorf("subcmd: no such command %q", args[0])
	return r.errorExit(args, err)
}

func (r *Runner) errorExit(args []string, err error) error {
	switch r.errorHandling {
	case flag.ContinueOnError:
		return err
	case flag.PanicOnError:
		panic(err)
	case flag.ExitOnError:
		r.Usage()
		if err == ErrHelp {
			os.Exit(0)
		}
		os.Exit(2)
	default:
		panicf("subcmd: bad ErrorHandling value %d", r.errorHandling)
	}
	panic("unreached")
}

// Run parses os.Args and dispatches to the correct subcommand given by cmds.
// It produces an error message listing the commands with their descriptions if
// a nonexistent subcommand is provided, or if the command is "help", "-h",
// "-help", or "--help". This error message may be customized by altering Usage.
// If the command provided isn't one of those in cmds, Run calls os.Exit(2)
// after printing the error message.
//
// Run panics if any command is named "help", "-h", "-help", or "--help",
// or if any two commands have the same name.
func Run(cmds []Command) {
	r := New(os.Args[0], cmds, flag.ExitOnError)
	r.Usage = func() { Usage(cmds) }
	r.Run(os.Args[1:])
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

// Usage prints a help message listing the possible commands.
// The function is a variable that may be changed to point at a custom function.
var Usage = func(cmds []Command) {
	defaultUsage(os.Args[0], cmds)
}

func defaultUsage(name string, cmds []Command) {
	fmt.Fprintf(os.Stderr, "Usage:\n\n  %s COMMAND\n\nPossible commands are:\n\n", name)
	PrintDefaults(cmds)
	fmt.Fprintf(os.Stderr, "\nRun '%s COMMAND -h' to see more information about a command.\n", name)
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
