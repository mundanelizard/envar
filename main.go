package main

import (
	"fmt"
	"github.com/mundanelizard/envi/internal/author"
	"github.com/mundanelizard/envi/internal/blob"
	"github.com/mundanelizard/envi/internal/commit"
	"github.com/mundanelizard/envi/internal/database"
	"github.com/mundanelizard/envi/internal/entry"
	"github.com/mundanelizard/envi/internal/tree"
	"github.com/mundanelizard/envi/internal/workspace"
	"github.com/mundanelizard/envi/pkg/cli"
	log "github.com/mundanelizard/envi/pkg/logger"
	"os"
	"path"
	"strings"
	"time"
)

var logger = log.New(os.Stdout, log.Info)
var wd string
var ed string
var od string

func loadWd() {
	var err error
	wd, err = os.Getwd()

	if err != nil {
		logger.Fatal(err)
		return
	}
}

func loadDd(enviDir string) {
	od = path.Join(enviDir, "objects")
}

func loadEd(baseDir string) {
	ed = path.Join(baseDir, ".envi")
}

func init() {
	loadWd()
	loadEd(wd)
	loadDd(ed)
}

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
			Usage:    "Commit message is required - ie 'envi commit -m 'commit message'",
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
	ws := workspace.New(wd)
	db := database.New(od)

	files, err := ws.ListFiles()
	if err != nil {
		logger.Fatal(err)
	}

	entries := make([]*entry.Entry, 0, len(files))

	for _, path := range files {
		data, err := ws.ReadFile(path)
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

		e := entry.New(path, b.Id())
		entries = append(entries, e)
	}

	t := tree.New(entries)
	err = db.Store(t)
	if err != nil {
		logger.Fatal(err)
	}

	aut := author.New(os.Getenv("ENVI_AUTHOR_NAME"), os.Getenv("ENVI_AUTHOR_EMAIL"), time.Now())
	message, _ := values.GetString("message")

	com := commit.New(t.Id(), aut, message)
	db.Store(com)

	os.WriteFile(path.Join(ed, "HEAD"), []byte(com.Id()), 0755)

	fmt.Printf("[(root-commit) %s] %s\n", com.Id(), strings.Split(message, "\n")[0])
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

	logger.Info(fmt.Sprintf("Initialised empty envi directory in %s", wd))
}
