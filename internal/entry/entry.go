package entry

import "os"

const (
	ModeRegular    = "100644"
	ModeExecutable = "100755"
)

type Entry struct {
	name string
	id   string
	stat os.FileInfo
}

func New(name, id string, stat os.FileInfo) *Entry {
	return &Entry{
		name: name,
		id:   id,
		stat: stat,
	}
}

func (e *Entry) Name() string {
	return e.name
}

func (e *Entry) Id() string {
	return e.id
}

func (e *Entry) Mode() string {
	if e.stat.Mode().IsRegular() {
		return ModeRegular
	}

	return ModeExecutable
}
