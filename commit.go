package main

import (
	"fmt"
	"github.com/mundanelizard/envi/internal/database"
	"github.com/mundanelizard/envi/internal/entry"
	"github.com/mundanelizard/envi/internal/refs"
	"github.com/mundanelizard/envi/internal/workspace"
	"github.com/mundanelizard/envi/pkg/cli"
	"os"
	"strings"
	"time"
)

func newCommitCommand() *cli.Command {
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

func handleCommit(values *cli.ActionArgs, args []string) {
	ws := workspace.New(wd)
	db := database.New(od)
	rs := refs.New(ed)

	paths, err := ws.ListFiles()
	if err != nil {
		logger.Fatal(err)
	}

	entries := make([]database.Enterable, 0, len(paths))

	for _, p := range paths {
		var data []byte
		data, err = ws.ReadFile(p)
		if err != nil {
			logger.Fatal(err)
			return
		}

		b := database.NewBlob(data)

		err = db.Store(b)
		if err != nil {
			logger.Fatal(err)
			return
		}

		stat, err := ws.Stat(p)
		if err != nil {
			logger.Fatal(err)
			return
		}

		e := entry.New(p, b.Id(), stat)
		entries = append(entries, e)
	}

	t := database.BuildTree(entries)

	t.Traverse(func(t *database.Tree) {
		err = db.Store(t)
		if err != nil {
			logger.Fatal(err)
			return
		}
	})

	aut := database.NewAuthor(os.Getenv("ENVI_AUTHOR_NAME"), os.Getenv("ENVI_AUTHOR_EMAIL"), time.Now())
	message, _ := values.GetString("message")

	pid, err := rs.ReadHead()
	if err != nil {
		logger.Fatal(err)
	}

	// todo => check if tree id is the latest tree id and skip the creation
	currId, err := rs.ReadHead()
	commit, err := database.ReadCommit(currId)

	if commit.TreeId() == t.Id() {
		fmt.Println("Nothing to commit, working tree clean")
		return
	}

	com := database.NewCommit(pid, t.Id(), aut, message)
	err = db.Store(com)
	if err != nil {
		logger.Fatal(err)
		return
	}

	err = rs.UpdateHead(com.Id())
	if err != nil {
		logger.Fatal(err)
		return
	}

	err = rs.UpdateHistory(com.Id())
	if err != nil {
		logger.Fatal(err)
		return
	}

	meta := ""
	if len(pid) == 0 {
		meta = "(root-commit)"
	}

	fmt.Printf("[%s %s] %s\n", meta, com.Id(), strings.Split(message, "\n")[0])
}
