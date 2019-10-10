package main

import (
	"flag"
	"fmt"

	"github.com/cespare/subcmd"
)

var cmds = []subcmd.Command{
	{
		Name:        "foobar",
		Description: "whiffle through the tulgey wood",
		Do:          foo,
	},
	{
		Name:        "baz",
		Description: "seek the manxome foe",
		SubCommands: []subcmd.Command{
			{
				Name:        "bazinga",
				Description: "ble bu li ke",
				Do:          bazinga,
			},
		},
	},
}

func foo(args []string) {
	fs := flag.NewFlagSet("foobar", flag.ExitOnError)
	a := fs.Bool("a", false, "Set option a")
	fs.Parse(args)
	fmt.Println("a:", *a)
}

func bazinga(args []string) {
	fs := flag.NewFlagSet("bazinga", flag.ExitOnError)
	n := fs.Int("n", 10, "Number of blah")
	fs.Parse(args)
	fmt.Println("n:", *n)
}

func main() {
	subcmd.Run(cmds)
}
