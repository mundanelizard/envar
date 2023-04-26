package command

import (
	"fmt"
	"os"
	"path"

	"github.com/mundanelizard/envi/internal/command/helpers"
	"github.com/mundanelizard/envi/internal/database"
	"github.com/mundanelizard/envi/internal/refs"
	"github.com/mundanelizard/envi/pkg/cli"
)

func Push() *cli.Command {
	secret := &cli.StringFlag{
		Value: "",
		Flag: cli.Flag{
			Name:     "secret",
			Usage:    "Repository secret to use when encrypting the codebase - ie `envi push -secret='SECRET'`",
			Required: true,
		},
	}

	return &cli.Command{
		Name:   "push",
		Action: handlePush,
		Flags:  []cli.Flagger{secret},
	}
}

func handlePush(values *cli.ActionArgs, args []string) {
	db := database.New(path.Join(wd, ".envi", "objects"))
	rs := refs.New(path.Join(wd, ".envi", "refs"))

	repo, err := os.ReadFile(path.Join(wd, ".envi", "remote"))
	if err != nil {
		logger.Fatal(err)
		return
	}

	pid, err := rs.Read()
	if err != nil {
		logger.Fatal(err)
		return
	}

	obj, err := db.Read(pid)
	if err != nil {
		logger.Fatal(err)
		return
	}

	commit, err := database.NewCommitFromByteArray(pid, obj)
	if err != nil {
		logger.Fatal(err)
		return
	}

	secret, _ := values.GetString("secret")

	comDir, encDir, err := helpers.CompressAndEncryptRepo(wd, string(repo), secret)
	if err != nil {
		logger.Fatal(err)
		return
	}

	liveRepo, err := srv.RetrieveRepo(string(repo))
	if err != nil {
		logger.Fatal(err)
	}

	oldCommitTreeId := ""

	if len(liveRepo.CommitId) != 0 {
		oldCommit, err := database.NewCommitFromByteArray(liveRepo.CommitId, obj)
		if err != nil {
			logger.Fatal(err)
			return
		}
		oldCommitTreeId = oldCommit.TreeId()
	}

	err = srv.PushRepo(string(repo), oldCommitTreeId, commit.TreeId(), commit.Id(), encDir, secret)
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

	fmt.Println("Successfully pushed", commit.Id(), "to remote server.")
}
