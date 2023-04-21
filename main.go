package main

import (
	"os"

	"github.com/mundanelizard/envi/internal/command"
	"github.com/mundanelizard/envi/pkg/cli"
)

func main() {
	var app = cli.New("jit")

	app.AddCommand(command.Init())
	app.AddCommand(command.Add())
	app.AddCommand(command.Commit())

	app.Execute(os.Args)
}
