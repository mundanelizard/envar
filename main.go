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
	app.AddCommand(command.Push())
	app.AddCommand(command.Pull())
	app.AddCommand(command.Revoke())
	app.AddCommand(command.Share())
	app.AddCommand(command.Clone())
	app.AddCommand(command.Login())
	app.AddCommand(command.Signup())

	app.Execute(os.Args)
}
