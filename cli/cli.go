package cli

import "fmt"

func New(name string) *Command {
	// default handler
	fn := func(command *ActionArgs, args []string) {
		fmt.Printf("%s\n\n", name)
		fmt.Printf("Arguments: %v\n", args)
		fmt.Printf("Command: %v\n", command)
	}
	return NewWithAction(name, fn)
}

func NewWithAction(name string, action func(command *ActionArgs, args []string)) *Command {
	return &Command{
		Name:   name,
		Action: action,
		root:   true,
	}
}
