package command

import (
	"os"
	"path"

	"github.com/mundanelizard/envi/internal/command/helpers"
	"github.com/mundanelizard/envi/pkg/cli"
)

func Clone() *cli.Command {
	secret := &cli.StringFlag{
		Value: "",
		Flag: cli.Flag{
			Name:     "secret",
			Usage:    "Repository secret to use when encrypting the codebase - ie `envi pull -secret='SECRET'`",
			Required: true,
		},
	}

	repo := &cli.StringFlag{
		Value: "",
		Flag: cli.Flag{
			Name:     "repo",
			Usage:    "Repository secret to use when encrypting the codebase - ie `envi pull -repo='repo-path'`",
			Required: true,
		},
	}

	return &cli.Command{
		Name:   "clone",
		Action: handleClone,
		Flags:  []cli.Flagger{secret, repo},
	}
}

func handleClone(values *cli.ActionArgs, _ []string) {
	secret, _ := values.GetString("secret")
	repo, _ := values.GetString("repo")

	// download the latest file from the server
	encDir, err := srv.PullRepo(repo)
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

	repoName := path.Base(repo)
	dest := path.Join(wd, repoName, ".envi")

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
}

func populateEnvironment(_ string) error {
	return nil
}
