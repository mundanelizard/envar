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
			Usage:    "Share user is required - ie 'envi share -user='username'",
			Required: true,
		},
	}

	return &cli.Command{
		Name:   "revoke",
		Action: handleRevoke,
		Flags:  []cli.Flagger{user},
	}
}

func handleRevoke(values *cli.ActionArgs, args []string) {
	repo, err := os.ReadFile(path.Join(wd, ".envi", "remote"))
	if err != nil {
		logger.Fatal(err)
		return
	}

	user, _ := values.GetString("user")

	err = srv.RevokeRepo(string(repo), user)

	if err != nil {
		logger.Fatal(err)
	}

	fmt.Println("Revoked repository access from", user)
}
