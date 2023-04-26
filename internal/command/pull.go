package command

import (
	"fmt"
	"os"
	"path"

	"github.com/mundanelizard/envi/internal/command/helpers"
	"github.com/mundanelizard/envi/pkg/cli"
)

func Pull() *cli.Command {
	secret := &cli.StringFlag{
		Value: "",
		Flag: cli.Flag{
			Name:     "secret",
			Usage:    "Repository secret to use when encrypting the codebase - ie `envi pull -secret='SECRET'`",
			Required: true,
		},
	}

	return &cli.Command{
		Name:   "pull",
		Action: handlePull,
		Flags:  []cli.Flagger{secret},
	}
}

func handlePull(values *cli.ActionArgs, _ []string) {
	secret, _ := values.GetString("secret")

	repo, err := os.ReadFile(path.Join(wd, ".envi", "remote"))
	if err != nil {
		logger.Fatal(err)
		return
	}

	// download the latest file from the server
	encDir, err := srv.PullRepo(string(repo))
	if err != nil {
		logger.Fatal(err)
		return
	}

	// decrypt the file with the secret in the pull
	comDir, err := helpers.DecryptCompressedEnvironment(encDir, secret)
	if err != nil {
		logger.Fatal(err)
		return
	}

	dest := path.Join(wd, ".envi")

	// delete current directory
	err = os.RemoveAll(dest)
	if err != nil {
		logger.Fatal(err)
		return
	}

	// decompressing file to the current destination
	err = helpers.DecompressEnvironment(comDir, dest)
	if err != nil {
		logger.Fatal(err)
		return
	}

	// replace the current directory with the current file
	err = populateEnvironment(dest)
	if err != nil {
		logger.Fatal(err)
		return
	}

	err = os.Remove(encDir)
	if err != nil {
		logger.Fatal(err)
		return
	}

	err = os.Remove(comDir)
	if err != nil {
		logger.Fatal(err)
		return
	}

	fmt.Println("Successfully pulled latest branch from remote.")
}
