package entry

type Entry struct {
	name string
	id   string
}

func New(name, id string) *Entry {
	return &Entry{
		name: name,
		id:   id,
	}
}

func (e *Entry) Name() string {
	return e.name
}

func (e *Entry) Id() string {
	return e.id
}
