package command

import (
	"fmt"
	"os"
	"path"

	"github.com/mundanelizard/envi/pkg/cli"
)

func Share() *cli.Command {
	user := &cli.StringFlag{
		Value: "",
		Flag: cli.Flag{
			Name:     "user",
			Usage:    "Share user is required - ie 'envi share -user='username'",
			Required: true,
		},
	}

	role := &cli.StringFlag{
		Value: "",
		Flag: cli.Flag{
			Name:     "role",
			Usage:    "Share role is required - ie 'envi share -role='R'",
			Required: true,
		},
	}

	return &cli.Command{
		Name:   "share",
		Action: handleShare,
		Flags:  []cli.Flagger{user, role},
	}
}

func handleShare(values *cli.ActionArgs, args []string) {
	repo, err := os.ReadFile(path.Join(wd, ".envi", "remote"))
	if err != nil {
		logger.Fatal(err)
		return
	}

	user, _ := values.GetString("user")
	role, _ := values.GetString("role")

	err = srv.ShareRepo(string(repo), user, role)

	if err != nil {
		logger.Fatal(err)
	}

	fmt.Println("Shared repository with", user, "with", role, "access")
}
