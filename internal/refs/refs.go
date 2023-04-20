package refs

import (
	"io"
	"os"
	"path"

	"github.com/mundanelizard/envi/internal/lockfile"
)

type Refs struct {
	dir      string
	headPath string
}

func New(dir string) *Refs {
	return &Refs{
		dir:      path.Join(dir, "refs"),
		headPath: path.Join(dir, "refs", "HEAD"),
	}
}

func (r *Refs) Update(cid string) error {
	lock := lockfile.New(r.headPath)

	err := lock.Hold()
	if err != nil {
		return err
	}

	err = lock.Write([]byte(cid + "\n"))
	if err != nil {
		return err
	}

	err = lock.Commit()
	if err != nil {
		return err
	}

	return os.WriteFile(r.headPath, []byte(cid), 0755)
}

func (r *Refs) Read() (string, error) {
	file, err := os.Open(r.headPath)
	defer file.Close()

	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		} else {
			return "", err
		}
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
