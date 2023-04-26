package command

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/mundanelizard/envi/internal/lockfile"
	"github.com/mundanelizard/envi/pkg/cli"
)

func Init() *cli.Command {
	message := &cli.StringFlag{
		Value: "",
		Flag: cli.Flag{
			Name:  "message",
			Usage: "init message",
		},
	}

	return &cli.Command{
		Name: "init",
		Flags: []cli.Flagger{
			message,
		},
		Action: handleInit,
	}
}

func handleInit(_ *cli.ActionArgs, args []string) {
	cwd := wd
	if len(args) == 1 {
		cwd = path.Join(cwd, args[0])
	}

	enviDir := path.Join(cwd, ".envi")

	// checking if an envi repository has been initialised
	stat, err := os.Stat(enviDir)
	if err == nil && stat.IsDir() {
		fmt.Println("Reinitialised current repository")
		return
	}

	err = createEnviSubdirectories(enviDir)
	if err != nil {
		logger.Fatal(err)
		return
	}

	// creating a repository on the server and saving the remote
	repoName := filepath.Base(cwd)
	endpoint, err := srv.CreateNewRepo(repoName)
	if err != nil {
		_, _ = logger.Error(err)
		_ = os.RemoveAll(enviDir)
		return
	}

	err = lockfile.WriteWithLock(path.Join(enviDir, "remote"), []byte(endpoint))
	if err != nil {
		logger.Fatal(err)
		return
	}

	// creating a envi file and populating it to only accept .env files
	err = lockfile.WriteWithLock(path.Join(cwd, ".envmatch"), []byte("**/*.env\n.env"))
	if err != nil {
		logger.Fatal(err)
		return
	}

	// updating/creating a .gitignore file to ignore .env, .envfiles, and .envi
	err = lockfile.AppendWithLock(path.Join(cwd, ".gitignore"), []byte("\n.env\n.envfile\n.envi"))
	if err != nil {
		logger.Fatal(err)
		return
	}

	fmt.Printf("Initialised empty envi directory in %s\n", wd)
}

func createEnviSubdirectories(enviDir string) error {
	dirs := []string{"objects", "refs"}

	for _, dir := range dirs {
		err := os.MkdirAll(path.Join(enviDir, dir), 0755)
		if err != nil {
			return err
		}
	}

	return nil
}
