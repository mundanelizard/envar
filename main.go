package main

import (
	"fmt"
	"github.com/mundanelizard/envi/internal/database"
	"github.com/mundanelizard/envi/internal/entry"
	"github.com/mundanelizard/envi/internal/index"
	"github.com/mundanelizard/envi/internal/refs"
	"github.com/mundanelizard/envi/internal/workspace"
	"github.com/mundanelizard/envi/pkg/cli"
	log "github.com/mundanelizard/envi/pkg/logger"
	"os"
	"path"
	"sort"
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

	app.AddCommand(NewInitCommand())
	app.AddCommand(NewAddCommand())
	app.AddCommand(NewCommitCommand())

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

func NewAddCommand() *cli.Command {
	return &cli.Command{
		Name:   "add",
		Flags:  []cli.Flagger{},
		Action: handleAdd,
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

func handleAdd(values *cli.ActionArgs, args []string) {
	ws := workspace.New(wd)
	db := database.New(od)
	ix := index.New(path.Join(ed, "index"))

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

	pid, err := rs.Read()
	if err != nil {
		logger.Fatal(err)
		return
	}

	fmt.Println("Parent Id: ", pid)

	err = detectEmptyCommit(db, pid, t.Id())
	if err != nil {
		logger.Fatal(err)
		return
	}

	com := database.NewCommit(pid, t.Id(), aut, message)
	err = db.Store(com)
	if err != nil {
		logger.Fatal(err)
		return
	}

	err = rs.Update(com.Id())
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

func detectEmptyCommit(db *database.Db, pid, currTreeId string) error {
	if len(pid) == 0 {
		return nil
	}

	obj, err := db.Read(pid)
	if err != nil {
		return err
	}

	parent, err := database.NewCommitFromByteArray(pid, obj)
	if err != nil {
		return err
	}

	if parent.TreeId() == currTreeId {
		return err
	}

	return nil
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
