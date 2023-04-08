package blob

import "io"

type Blob struct {
}

func New(writer io.ReadWriter) *Blob {

	return &Blob{}
}
