package command

import (
	"fmt"
	"github.com/mundanelizard/envi/internal/database"
	"github.com/mundanelizard/envi/internal/lockfile"
	"github.com/mundanelizard/envi/internal/refs"
	"os"
	"path"
	"path/filepath"
	"strings"

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

	fmt.Println("Successfully pulled latest branch from remote.")
}

type tree struct {
	id   string
	path string
}

func populateEnvironment(envDir string) error {
	db := database.New(path.Join(envDir, "objects"))
	rs := refs.New(path.Join(envDir, "refs"))

	// retrieving the previous commit id
	pid, err := rs.Read()
	if err != nil {
		return err
	}

	if len(pid) == 0 {
		return nil
	}

	data, err := db.Read(pid)
	if err != nil {
		return err
	}

	commit, err := database.NewCommitFromByteArray(pid, data)
	initialTree := tree{
		id:   commit.TreeId(),
		path: filepath.Dir(envDir),
	}
	trees := []tree{initialTree}

	for len(trees) != 0 {
		t := trees[0]
		trees = trees[1:]

		// read current tree
		data, err := db.Read(t.id)
		if err != nil {
			return err
		}

		err = os.MkdirAll(t.path, 0655)
		if err != nil && !os.IsExist(err) {
			return nil
		}

		entries, err := extractEntriesFromTree(data)

		for _, entry := range entries {
			p := path.Join(t.path, entry.name)

			if entry.t == "tree" {
				t = tree{
					id:   entry.id,
					path: p,
				}
				trees = append(trees, t)
				continue
			}

			data, err := db.Read(entry.id)
			if err != nil {
				return err
			}

			content := string(data)
			content = strings.Join(strings.Split(content, "\x00")[1:], "")

			err = lockfile.WriteWithLock(p, []byte(content))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type treeEntry struct {
	id   string
	t    string
	name string
}

func extractEntriesFromTree(raw []byte) ([]treeEntry, error) {
	data := string(raw)

	chunks := strings.Split(data, "\x00")

	// clean data
	data = strings.Join(chunks[1:], "")

	chunks = strings.Split(data, "\n")

	var entries []treeEntry

	for _, chunk := range chunks {
		if len(chunk) == 0 {
			continue
		}

		columns := strings.Split(chunk, " ")

		t := "blob"
		name := columns[1]
		id := columns[2]

		if columns[0] != "100644" {
			t = "tree"
		}

		entry := treeEntry{
			t:    t,
			id:   id,
			name: name,
		}

		entries = append(entries, entry)
	}

	return entries, nil
}
