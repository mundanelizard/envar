package main

import (
	"fmt"
	"github.com/mundanelizard/envi/cli"
	"log"
	"os"
	"path"
)

func main() {
	var app = cli.New("jit")

	app.AddCommand(NewCommitCommand())
	app.AddCommand(NewInitCommand())

	app.Execute(os.Args)
}

func NewCommitCommand() *cli.Command {
	message := &cli.StringFlag{
		Value: "",
		Flag: cli.Flag{
			Name:     "message",
			Usage:    "commit message",
			Required: true,
		},
	}

	return &cli.Command{
		Name: "commit",
		Flags: []cli.Flagger{
			message,
		},
		Action: handleCommit,
	}
}

func NewInitCommand() *cli.Command {
	message := &cli.StringFlag{
		Value: "",
		Flag: cli.Flag{
			Name:  "message",
			Usage: "init message",
		},
	}
	return &cli.Command{
		Name: "init",
		Flags: []cli.Flagger{
			message,
		},
		Action: handleInit,
	}
}

func handleCommit(values *cli.ActionArgs, args []string) {
	fmt.Println("Handling commit")
	fmt.Println(values)
	fmt.Println("args")
}

func handleInit(_ *cli.ActionArgs, args []string) {
	wd, err := os.Getwd()

	if err != nil {
		log.Fatalln(err)
	}

	if len(args) == 1 {
		wd = path.Join(wd, args[0])
	}

	wd = path.Join(wd, ".envi")
	dirs := []string{"objects", "refs"}

	for _, dir := range dirs {
		err = os.MkdirAll(path.Join(wd, dir), 0755)
		if err != nil {
			log.Fatalln(err)
		}
	}

	fmt.Printf("Initialised empty envi directory in %s\n", wd)
}