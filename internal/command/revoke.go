package command

import (
	"fmt"
	"os"
	"path"

	"github.com/mundanelizard/envi/pkg/cli"
)

func Revoke() *cli.Command {
	user := &cli.StringFlag{
		Value: "",
		Flag: cli.Flag{
			Name:     "user",
			Usage:    "Share user is required - ie 'envi revoke -user='username'",
			Required: true,
		},
	}

	role := &cli.StringFlag{
		Value: "",
		Flag: cli.Flag{
			Name:     "role",
			Usage:    "Share role is required - ie 'envi revoke -role='R'",
			Required: true,
		},
	}

	return &cli.Command{
		Name:   "revoke",
		Action: handleRevoke,
		Flags:  []cli.Flagger{user, role},
	}
}

func handleRevoke(values *cli.ActionArgs, args []string) {
	repo, err := os.ReadFile(path.Join(wd, ".envi", "remote"))
	if err != nil {
		logger.Fatal(err)
		return
	}

	user, _ := values.GetString("user")
	role, _ := values.GetString("role")

	err = srv.RevokeRepo(string(repo), user, role)

	if err != nil {
		logger.Fatal(err)
	}

	fmt.Println("Revoked repository access from", user)
}
