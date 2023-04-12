package tree

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/mundanelizard/envi/internal/entry"
)

const (
	TYPE = "tree"
	MODE = "100644"
)

type Tree struct {
	id      string
	entries []*entry.Entry
}

func New(entries []*entry.Entry) *Tree {
	return &Tree{
		entries: entries,
	}
}

// implementing methods for database.Storable

func (t *Tree) String() string {
	var buf bytes.Buffer

	for _, e := range t.entries {
		_, err := fmt.Fprintf(&buf, "%s %s\x00", MODE, e.Name())
		if err != nil {
			panic(err)
		}
		buf.Write(hexDecode(e.Id()))
	}

	return buf.String()
}

func (t *Tree) SetId(id string) {
	t.id = id
}

func (t *Tree) Id() string {
	return t.id
}

func (t *Tree) Type() string {
	return TYPE
}

func hexDecode(h string) []byte {
	data, err := hex.DecodeString(h)
	if err != nil {
		panic(err)
	}
	return data
}
