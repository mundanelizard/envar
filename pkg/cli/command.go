package cli

import (
	"flag"
	"log"
)

type Command struct {
	Name     string
	Usage    string
	Commands []*Command
	Flags    []Flagger
	Action   func(command *ActionArgs, args []string)
	root     bool
}

func (cmd *Command) AddCommand(command *Command) {
	cmd.Commands = append(cmd.Commands, command)
}

func (cmd *Command) Execute(args []string) {
	if !cmd.root && args[0] != cmd.Name {
		log.Fatalln("Invalid mapping " + args[0] + " to " + cmd.Name)
	}

	for _, c := range cmd.Commands {
		if len(args) == 1 {
			break
		}

		if args[1] == c.Name {
			c.Execute(args[1:])
			return
		}
	}

	values := NewActionArgs()
	set := flag.NewFlagSet(cmd.Name, flag.ExitOnError)

	for _, f := range cmd.Flags {
		values.Set(f.GetID(), f.Attach(set))
	}

	args = args[1:]
	err := set.Parse(args)

	if err != nil {
		log.Fatalln(err)
	}

	// TODO => validate the response

	for _, f := range cmd.Flags {
		ptr := values.Get(f.GetID())

		err = f.Validate(ptr)

		if err != nil {
			log.Fatalln(err)
		}
	}

	cmd.Action(values, args)

	return

}
