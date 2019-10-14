package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cespare/subcmd"
)

var cmds = []subcmd.Command{
	{
		Name:        "foo",
		Description: "perform foo tasks",
		Do:          foo,
	},
	{
		Name:        "xyz",
		Description: "do some other thing",
		Do:          xyz,
	},
}

var fooCmds = []subcmd.Command{
	{
		Name:        "bar",
		Description: "something about bar",
		Do:          foobar,
	},
	{
		Name:        "baz",
		Description: "something about baz",
		Do:          foobaz,
	},
}

func foo(args []string) {
	cmdName := fmt.Sprintf("%s foo", os.Args[0])
	r := subcmd.New(cmdName, fooCmds, flag.ExitOnError)
	r.Run(args)
}

func foobar(args []string) {
	fs := flag.NewFlagSet("foo bar", flag.ExitOnError)
	a := fs.Bool("a", false, "Set option a")
	fs.Parse(args)
	fmt.Println("a:", *a)
}

func foobaz(args []string) {
	fs := flag.NewFlagSet("foo baz", flag.ExitOnError)
	a := fs.Bool("a", false, "Set option a")
	fs.Parse(args)
	fmt.Println("a:", *a)
}

func xyz(args []string) {
	fs := flag.NewFlagSet("xyz", flag.ExitOnError)
	n := fs.Int("n", 10, "Number of blah")
	fs.Parse(args)
	fmt.Println("n:", *n)
}

func main() {
	subcmd.Run(cmds)
}
