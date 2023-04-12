// Package blob /
// Blob represents a file contents
package blob

const (
	TYPE = "blob"
)

type Blob struct {
	data  string
	bytes []byte
	id    string
}

func New(data []byte) *Blob {
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
	return TYPE
}

/*
type Tree struct {
    Oid     string
    Entries []Entry
}

type Entry struct {
    Name string
    Oid  string
}

func (t *Tree) Type() string {
    return "tree"
}

func (t *Tree) String() string {
    var buf bytes.Buffer
    for _, entry := range t.Entries {
        mode := "100644"
        fmt.Fprintf(&buf, "%s %s\x00", mode, entry.Name)
        buf.Write(hexDecode(entry.Oid))
    }
    return buf.String()
}

func hexDecode(hex string) []byte {
    data, err := hex.DecodeString(hex)
    if err != nil {
        panic(err)
    }
    return data
}


*/
