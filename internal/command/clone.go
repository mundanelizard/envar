package command

import (
	"errors"
	"os"
	"path"

	"github.com/mundanelizard/envi/internal/command/helpers"
	"github.com/mundanelizard/envi/pkg/cli"
)

func Clone() *cli.Command {
	return &cli.Command{
		Name:   "pull",
		Action: handleClone,
	}
}

func handleClone(values *cli.ActionArgs, args []string) {
	// check if user is authenticated
	_, err := srv.RetrieveUser()
	if err != nil {
		logger.Fatal(err)
		return
	}

	if len(args) != 2 {
		logger.Fatal(errors.New("expected args of length 1"))
		return
	}

	secret, _ := values.GetString("secret")
	repo := args[0]

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
	dest := path.Join(wd, repoName)

	err = helpers.DecompressEnvironment(comDir, dest)
	if err != nil {
		logger.Fatal(err)
		return
	}

	// replace the current directory with the current file
	err = populateEnvironment()
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

func populateEnvironment() error {
	return nil
}
