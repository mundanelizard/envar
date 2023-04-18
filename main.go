package main

import (
	"github.com/mundanelizard/envi/pkg/cli"
	log "github.com/mundanelizard/envi/pkg/logger"
	"os"
	"path"
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

	app.AddCommand(newInitCommand())
	app.AddCommand(newAddCommand())
	app.AddCommand(newCommitCommand())

	app.Execute(os.Args)
}
