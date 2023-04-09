package main

import (
	"fmt"
	"github.com/mundanelizard/envi/internal/blob"
	"github.com/mundanelizard/envi/internal/database"
	"github.com/mundanelizard/envi/internal/workspace"
	"github.com/mundanelizard/envi/pkg/cli"
	log "github.com/mundanelizard/envi/pkg/logger"
	"os"
	"path"
)

var logger = log.New(os.Stdout, log.Info)

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
	baseDir, err := os.Getwd()

	if err != nil {
		logger.Fatal(err)
		return
	}

	enviDir := path.Join(baseDir, ".envi")
	dbDir := path.Join(enviDir, "objects")

	ws := workspace.New(enviDir)
	db := database.New(dbDir)

	files, err := ws.ListFiles()

	for _, file := range files {
		data, err := ws.ReadFile(file)
		if err != nil {
			logger.Fatal(err)
			return
		}

		b := blob.New(data)

		err = db.Store(b)

		if err != nil {
			logger.Fatal(err)
			return
		}
	}

	fmt.Println("Handling commit")
	fmt.Println(values)
	fmt.Println("args")
}

func handleInit(_ *cli.ActionArgs, args []string) {
	wd, err := os.Getwd()

	if err != nil {
		logger.Fatal(err)
		return
	}

	if len(args) == 1 {
		wd = path.Join(wd, args[0])
	}

	wd = path.Join(wd, ".envi")
	dirs := []string{"objects", "refs"}

	for _, dir := range dirs {
		err = os.MkdirAll(path.Join(wd, dir), 0755)
		if err != nil {
			logger.Fatal(err)
		}
	}

	fmt.Printf("Initialised empty envi directory in %s\n", wd)
}
