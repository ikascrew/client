package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ikascrew/client"
	"github.com/ikascrew/client/config"
	"github.com/ikascrew/client/tool"

	"golang.org/x/xerrors"
)

var joycon bool
var powermate bool

func init() {
	flag.BoolVar(&joycon, "joy", false, "Use JoyCon")
	flag.BoolVar(&powermate, "pm", false, "Use Powermate")
}

func main() {

	flag.Parse()

	err := run()
	if err != nil {
		fmt.Printf("ikascrew client error: %+v\n", err)
		os.Exit(1)
	}

	fmt.Println("bye!")
}

func run() error {

	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		return xerrors.Errorf("ikascrew client command required[start,create]")
	}

	var err error
	command := args[0]

	switch command {
	case "start":

		opts := make([]config.Option, 0)
		if joycon {
			opts = append(opts, config.Controller(config.ControllerTypeJoyCon))
		}
		if powermate {
			opts = append(opts, config.UsePowermate())
		}
		err = client.Start(opts...)
	case "create":
		id := args[1]
		err = tool.CreateProject(id)
	default:
		err = fmt.Errorf("ikascrew client sub command not found %s[start,create]", command)
	}

	if err != nil {
		return xerrors.Errorf("ikascrew client command[%s]: %w", command, err)
	}

	return nil
}
