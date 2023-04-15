// Package workspace /
// Workspace class is responsible for the files in the working tree - the files edited.
package workspace

import (
	"os"
	"path"
	"path/filepath"
	"sort"
)

type Workspace struct {
	dir string
}

// paths to ignore when listing files in the working tree
var ignore = []string{".", "..", ".envi", ".git", ".idea"}

func New(dir string) *Workspace {
	return &Workspace{
		dir: dir,
	}
}

func (ws *Workspace) ListFiles() ([]string, error) {
	return ws.listFiles(ws.dir)
}

func (ws *Workspace) listFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var paths []string
	for _, entry := range entries {
		if has(ignore, entry.Name()) {
			continue
		}

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
