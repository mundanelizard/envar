package workspace

import (
	"os"
)

type workspace struct {
	wd string
}

var ignore = []string{".", "..", ".envi", ".git"}

func New(wd string) *workspace {
	return &workspace{
		wd: wd,
	}
}

func (ws *workspace) ListFiles() ([]os.DirEntry, error) {
	dirEntries, err := os.ReadDir(ws.wd)

	if err != nil {
		return nil, err
	}

	filteredNames := make([]os.DirEntry, 0)

	for _, de := range dirEntries {
		if has(ignore, de.Name()) {
			continue
		}

		filteredNames = append(filteredNames, de)
	}

	return filteredNames, nil
}

func has(list []string, key string) bool {
	for _, value := range list {
		if value == key {
			return true
		}
	}

	return false
}
