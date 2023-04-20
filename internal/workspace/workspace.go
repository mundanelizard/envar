// Package workspace /
// Workspace class is responsible for the files in the working tree - the files edited.
package workspace

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

type Workspace struct {
	dir   string
	globs []string
}

func New(dir string) *Workspace {
	return &Workspace{
		dir:   dir,
		globs: []string{"**/*.env"},
	}
}

func (ws *Workspace) ListFiles() ([]string, error) {
	err := ws.loadMatch()

	if err != nil {
		return nil, err
	}

	return ws.listFiles(ws.dir)
}

func (ws *Workspace) loadMatch() error {
	file, err := os.ReadFile(path.Join(".envmatch"))
	if err != nil {
		return err
	}

	ws.globs = append(ws.globs, strings.Split(string(file), "\n")...)
	return nil
}

func (ws *Workspace) match(path string) bool {
	for _, glob := range ws.globs {
		ok, _ := filepath.Match(glob, path)
		if ok {
			return true
		}
	}

	return false
}

func (ws *Workspace) listFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var paths []string
	for _, entry := range entries {
		// skipping files that do not match the glob
		if !ws.match(entry.Name()) {
			continue
		}

		fmt.Println("Committing", entry.Name())

		newDir := path.Join(dir, entry.Name())

		if !entry.IsDir() {
			// extracting the path relative to the workspace directory
			p, _ := filepath.Rel(ws.dir, newDir)
			paths = append(paths, p)
			continue
		}

		subPaths, err := ws.listFiles(newDir)
		if err != nil {
			return nil, err
		}

		paths = append(paths, subPaths...)
	}

	// sorting an array of strings
	sort.Strings(paths)

	return paths, nil
}

func (ws *Workspace) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (ws *Workspace) Stat(name string) (os.FileInfo, error) {
	stat, err := os.Stat(path.Join(ws.dir, name))
	if err != nil {
		return nil, err
	}

	return stat, err
}

func has(list []string, key string) bool {
	for _, value := range list {
		if value == key {
			return true
		}
	}

	return false
}
