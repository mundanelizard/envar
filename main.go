package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func NewApp(name string) *Command {
	return &Command{
		Name: name,
		Action: func(command *ParsedValues, args []string) {
			fmt.Println(args)
			fmt.Println(command)
		},
		root: true,
	}
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

	values := map[string]interface{}{}
	set := flag.NewFlagSet(cmd.Name, flag.ExitOnError)

	var value *string

	for _, f := range cmd.Flags {
		value = f.Attach(set).(*string)
		values[f.GetID()] = value
	}

	args = args[1:]
	err := set.Parse(args)

	if err != nil {
		log.Fatalln(err)
	}

	// TODO => validate the response
	//for _, f := range cmd.Flags {
	//	value := values[f.GetID()]
	//	set.PrintDefaults()
	//	os.Exit(1)
	//	// update this to work properly - i guess right now
	//}

	cmd.Action(&values, args)
	return

}

type Command struct {
	Name     string
	Usage    string
	Commands []*Command
	Flags    []Flag
	Action   func(command *ParsedValues, args []string)
	root     bool
}

type ParsedValues = map[string]interface{}

type Flag interface {
	Attach(set *flag.FlagSet) any
	GetID() string
}

type IntFlag struct {
	Name     string
	Usage    string
	Required bool
	Value    int
}

func (f *IntFlag) Attach(set *flag.FlagSet) any {
	return set.Int(f.Name, f.Value, f.Usage)
}

func (f *IntFlag) GetID() string {
	return f.Name
}

type StringFlag struct {
	Name     string
	Usage    string
	Required bool
	Value    string
}

func (f *StringFlag) Attach(set *flag.FlagSet) any {
	return set.String(f.Name, f.Value, f.Usage)
}

func (f *StringFlag) GetID() string {
	return f.Name
}

type BoolFlag struct {
	Name  string
	Usage string
	Value bool
}

func (f *BoolFlag) Attach(set *flag.FlagSet) any {
	return set.Bool(f.Name, f.Value, f.Usage)
}

func (f *BoolFlag) BoolFlag() string {
	return f.Name
}

// build a cli utility

func main() {
	var app = NewApp("jit")

	app.AddCommand(&Command{
		Name: "commit",
		Flags: []Flag{
			&StringFlag{
				Name:     "message",
				Usage:    "message commit message",
				Required: true,
				Value:    "",
			},
		},
		Action: handleCommit,
	})

	app.Execute(os.Args)
}

func handleCommit(commands *ParsedValues, args []string) {
	fmt.Println("Handling commit")
	fmt.Println(commands)
	fmt.Println(args)
	for key, val := range *commands {
		fmt.Print(key)
		fmt.Print(" - ")
		fmt.Print(*((val).(*string)))
		fmt.Println()
	}
}
