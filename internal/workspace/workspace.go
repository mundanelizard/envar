// Package workspace /
// Workspace class is responsible for the files in the working tree - the files edited.
package workspace

import (
	"os"
)

type Workspace struct {
	wd string
}

// paths to ignore when listing files in the working tree
var ignore = []string{".", "..", ".envi", ".git"}

func New(wd string) *Workspace {
	return &Workspace{
		wd: wd,
	}
}

func (ws *Workspace) ListFiles() ([]string, error) {
	dirEntries, err := os.ReadDir(ws.wd)

	if err != nil {
		return nil, err
	}

	paths := make([]string, 0)

	for _, de := range dirEntries {
		// todo => optimise this by using a map for the check.
		if de.IsDir() || has(ignore, de.Name()) {
			continue
		}

		paths = append(paths, de.Name())
	}

	return paths, nil
}

func (ws *Workspace) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func has(list []string, key string) bool {
	for _, value := range list {
		if value == key {
			return true
		}
	}

	return false
}
