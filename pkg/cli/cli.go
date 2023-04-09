package cli

import "fmt"

type ActionFunc func(command *ActionArgs, args []string)

func New(name string) *Command {
	// default handler
	fn := func(command *ActionArgs, args []string) {
		fmt.Printf("%s\n\n", name)
		fmt.Printf("Arguments: %v\n", args)
		fmt.Printf("Command: %v\n", command)
	}
	return NewWithAction(name, fn)
}

func NewWithAction(name string, action ActionFunc) *Command {
	return &Command{
		Action: action,
		Name:   name,
		root:   true,
	}
}

func (cmd *Command) SetAction(action ActionFunc) {
	cmd.Action = action
}
