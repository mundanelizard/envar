package command

import (
	"os"

	log "github.com/mundanelizard/envi/pkg/logger"
)

var logger = log.New(os.Stdout, log.Info)
var wd string

func loadWd() {
	var err error
	wd, err = os.Getwd()

	if err != nil {
		logger.Fatal(err)
		return
	}
}

func init() {
	loadWd()
}