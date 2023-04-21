// Package workspace /
// Workspace class is responsible for the files in the working tree - the files edited.
package workspace

import (
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Workspace struct {
	dir   string
	globs []string
}

func New(dir string) *Workspace {
	return &Workspace{
		dir:   dir,
		globs: []string{},
	}
}

func (ws *Workspace) ListFiles() ([]string, error) {
	err := ws.loadMatch()
	if err != nil {
		return nil, err
	}

	paths, err := ws.listFiles(ws.dir)
	if err != nil {
		return nil, err
	}

	var filteredPaths []string
	for _, path := range paths {
		// skipping files that do not match the glob
		if !ws.match(path) {
			continue
		}

		filteredPaths = append(filteredPaths, path)
	}

	return filteredPaths, nil
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

func (ws *Workspace) loadMatch() error {
	file, err := os.ReadFile(path.Join(".envmatch"))
	if err != nil {
		return err
	}

	globs := strings.Split(string(file), "\n")

	for _, glob := range globs {
		trimmedGlob := strings.Trim(glob, "<>")
		if len(trimmedGlob) == 0 {
			continue
		}

		ws.globs = append(ws.globs, trimmedGlob)
	}

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
	// todo => swap this with fs.WalkDir :)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var paths []string
	for _, entry := range entries {
		fullPath := path.Join(dir, entry.Name())

		if !entry.IsDir() {
			// extracting the path relative to the workspace directory
			p, _ := filepath.Rel(ws.dir, fullPath)
			paths = append(paths, p)
			continue
		}

		subPaths, err := ws.listFiles(fullPath)
		if err != nil {
			return nil, err
		}

		paths = append(paths, subPaths...)
	}

	return paths, nil
}
