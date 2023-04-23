package command

import (
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/mundanelizard/envi/internal/database"
	"github.com/mundanelizard/envi/internal/entry"
	"github.com/mundanelizard/envi/internal/refs"
	"github.com/mundanelizard/envi/internal/workspace"
	"github.com/mundanelizard/envi/pkg/cli"
)

func Commit() *cli.Command {
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
	db := database.New(path.Join(wd, ".envi", "objects"))
	rs := refs.New(path.Join(wd, ".envi", "refs"))

	entries, err := getWorkspaceEntries(db)
	if err != nil {
		logger.Fatal(err)
		return
	}

	t := database.BuildTree(entries)

	t.Traverse(func(t *database.Tree) {
		err = db.Store(t)
		if err != nil {
			logger.Fatal(err)
			return
		}
	})

	// retrieving the previous commit id
	pid, err := rs.Read()
	if err != nil {
		logger.Fatal(err)
		return
	}

	if len(pid) == 0 && len(entries) == 0 {
		fmt.Println("Working on clean repository, nothing to commit.")
		return
	}

	stale, err := detectEmptyCommit(db, pid, t.Id())
	if err != nil {
		logger.Fatal(err)
		return
	}

	if stale {
		fmt.Println("Working on a clean tree, nothing to commit")
		return
	}

	// retreiving the user data
	user, err := srv.RetrieveUser()
	if err != nil {
		logger.Fatal(err)
		return
	}

	aut := database.NewAuthor(user.Username, time.Now())
	message, _ := values.GetString("message")

	com := database.NewCommit(pid, t.Id(), aut, message)
	err = db.Store(com)
	if err != nil {
		logger.Fatal(err)
		return
	}

	// updating the ref head to contain the current commit
	err = rs.Update(com.Id())
	if err != nil {
		logger.Fatal(err)
		return
	}

	fmt.Printf("[%s] %s\n", com.Id(), strings.Split(message, "\n")[0])
}

func detectEmptyCommit(db *database.Db, pid, currTreeId string) (bool, error) {
	if len(pid) == 0 {
		return false, nil
	}

	obj, err := db.Read(pid)
	if err != nil {
		return false, err
	}

	parent, err := database.NewCommitFromByteArray(pid, obj)
	if err != nil {
		return false, err
	}

	if parent.TreeId() == currTreeId {
		return true, nil
	}

	return false, nil
}

func getWorkspaceEntries(db *database.Db) ([]database.Enterable, error) {
	ws := workspace.New(wd)

	paths, err := ws.ListFiles()
	if err != nil {
		logger.Fatal(err)
	}

	entries := make([]database.Enterable, 0, len(paths))

	for _, p := range paths {
		var data []byte
		data, err = ws.ReadFile(p)
		if err != nil {
			return nil, err
		}

		b := database.NewBlob(data)

		err = db.Store(b)
		if err != nil {
			return nil, err
		}

		stat, err := ws.Stat(p)
		if err != nil {
			return nil, err
		}

		e := entry.New(p, b.Id(), stat)

		entries = append(entries, e)
	}

	return entries, nil
}
