package index

import (
	"bytes"
	"crypto/sha1"
	"errors"
	"hash"
	"io"
)

var (
	ErrEndOfFile = errors.New("unexpected end-of-file while reading index")
	ErrInvalid   = errors.New("checksum does not match value stored on disk")
)

type Checksum struct {
	file   io.ReadSeeker
	digest hash.Hash
}

const (
	ChecksumSize = 20
)

func NewChecksum(file io.ReadSeeker) *Checksum {
	return &Checksum{
		file:   file,
		digest: sha1.New(),
	}
}

func (cs *Checksum) Read(size int64) ([]byte, error) {
	data := make([]byte, size)
	n, err := cs.file.Read(data)
	if err != nil && err != io.EOF {
		return nil, err
	}
	if int64(n) != size {
		return nil, ErrEndOfFile
	}
	cs.digest.Write(data)
	return data, nil
}

func (cs *Checksum) verifyChecksum() error {
	sum := make([]byte, ChecksumSize)
	_, err := cs.file.Read(sum)
	if err != nil {
		return err
	}
	if !bytes.Equal(sum, cs.digest.Sum(nil)) {
		return ErrInvalid
	}
	return nil
}
