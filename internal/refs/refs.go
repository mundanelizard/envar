package refs

import (
	"github.com/mundanelizard/envi/internal/lockfile"
	"os"
	"path"
)

type Refs struct {
	dir      string
	headPath string
}

func New(dir string) *Refs {
	return &Refs{
		dir:      dir,
		headPath: path.Join(dir, "HEAD"),
	}
}

func (r *Refs) UpdateHead(cid string) error {
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

func (r *Refs) ReadHead() (string, error) {
	file, err := os.Open(r.headPath)
	defer file.Close()

	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		} else {
			return "", err
		}
	}

	buf := make([]byte, 0)
	_, err = file.Read(buf)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func (r *Refs) UpdateHistory(commitId string) error {

	return nil
}