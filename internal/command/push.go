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
	return &cli.Command{
		Name:   "push",
		Action: handlePush,
	}
}

func handlePush(values *cli.ActionArgs, args []string) {
	db := database.New(path.Join(wd, ".envi", "objects"))
	rs := refs.New(path.Join(wd, ".envi", "refs"))

	_, err := srv.RetrieveUser()
	if err != nil {
		logger.Fatal(err)
		return
	}

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

	err = srv.PushRepo(string(repo), commit.TreeId(), encDir, secret)
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

	fmt.Println("Environment pushed! Save encryption key", secret)
}
