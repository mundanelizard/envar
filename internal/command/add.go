package command


import (
	"path"
	"sort"

	"github.com/mundanelizard/envi/internal/database"
	"github.com/mundanelizard/envi/internal/index"
	"github.com/mundanelizard/envi/internal/workspace"
	"github.com/mundanelizard/envi/pkg/cli"
)


func Add() *cli.Command {
	return &cli.Command{
		Name:   "add",
		Flags:  []cli.Flagger{},
		Action: handleAdd,
	}
}

func handleAdd(values *cli.ActionArgs, args []string) {
	ws := workspace.New(wd)
	db := database.New(path.Join(wd, ".envi", "objects"))
	ix := index.New(path.Join(wd, ".envi", "index"))

	err := ix.Load()

	if err != nil {
		logger.Fatal(err)
		return
	}

	var paths []string

	if len(args) != 0 {
		sort.Strings(args)
		paths = args
	} else {
		files, err := ws.ListFiles()
		if err != nil {
			logger.Fatal(err)
			return
		}

		paths = files
	}

	for _, p := range paths {
		data, err := ws.ReadFile(p)
		if err != nil {
			logger.Fatal(err)
			return
		}

		stat, err := ws.Stat(p)
		if err != nil {
			logger.Fatal(err)
			return
		}

		blob := database.NewBlob(data)
		err = db.Store(blob)
		if err != nil {
			logger.Fatal(err)
			return
		}
		ix.Add(p, blob.Id(), stat)
	}

	err = ix.WriteUpdates()
	if err != nil {
		logger.Fatal(err)
		return
	}
}

