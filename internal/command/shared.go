package command

import (
	"os"

	"github.com/mundanelizard/envi/internal/server"
	log "github.com/mundanelizard/envi/pkg/logger"
)

var logger = log.New(os.Stdout, log.Info)
var wd string
var srv = server.New("https://localhost:9000/")

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