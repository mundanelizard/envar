package main

import (
	"fmt"
	"github.com/mundanelizard/envi/pkg/cli"
	"os"
	"path"
)

func newInitCommand() *cli.Command {
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

func handleInit(_ *cli.ActionArgs, args []string) {
	newEnvDir := ed
	if len(args) == 1 {
		newEnvDir = path.Join(newEnvDir, args[0])
	}

	dirs := []string{"objects", "refs"}

	for _, dir := range dirs {
		err := os.MkdirAll(path.Join(newEnvDir, dir), 0755)
		if err != nil {
			logger.Fatal(err)
		}
	}

	fmt.Printf("Initialised empty envi directory in %s", wd)
}
