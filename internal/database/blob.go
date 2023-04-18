package database

type Blob struct {
	data  string
	bytes []byte
	id    string
}

func NewBlob(data []byte) *Blob {
	str := string(data)
	return &Blob{
		data:  str,
		bytes: nil,
	}
}

// implementing methods for database.Storable

func (blob *Blob) Size() int {
	return len(blob.data)
}

func (blob *Blob) String() string {
	return blob.data
}

func (blob *Blob) SetId(id string) {
	blob.id = id
}

func (blob *Blob) Id() string {
	return blob.id
}

func (blob *Blob) Type() string {
	return "blob"
}
