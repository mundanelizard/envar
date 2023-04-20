// Package database /
// this is responsible for managing the files in .envi/objects.
package database

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
)

type Db struct {
	dir string
}

type Storable interface {
	Id() string
	SetId(id string)
	String() string
	Type() string
}

func New(dir string) *Db {
	return &Db{
		dir: dir,
	}
}

func (db *Db) Read(id string) ([]byte, error) {
	file, err := os.Open(path.Join(db.dir, id[:2], id[2:]))

	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return decompress(data)
}

func (db *Db) Store(blob Storable) error {
	data := blob.String()
	size := len(data)
	content := []byte(fmt.Sprintf("%s %d\x00%s", blob.Type(), size, data))
	digest := hash(content)
	blob.SetId(digest)
	return db.write(blob, content)
}

func (db *Db) write(blob Storable, content []byte) error {
	id := blob.Id()
	if len(id) == 0 {
		return errors.New("invalid Blob.Id: requires a blob id to store blob")
	}

	dirName := path.Join(db.dir, id[0:2])
	objPath := path.Join(dirName, id[2:])

	exist, err := db.exists(objPath)
	switch {
	case err != nil:
		return err
	case exist:
		return nil
	}

	buf, err := compress(content)
	if err != nil {
		return err
	}

	return db.writeFile(objPath, buf)
}

func (db *Db) writeFile(path string, buf bytes.Buffer) error {
	dir := filepath.Dir(path)
	exist, err := db.exists(dir)
	switch {
	case err != nil:
		return err
	case !exist:
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	tempPath := path + ".temp"

	file, err := os.OpenFile(tempPath, os.O_CREATE|os.O_RDWR|os.O_EXCL, 0755)
	if err != nil {
		return err
	}

	_, err = file.Write(buf.Bytes())
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	err = os.Rename(tempPath, path)
	if err != nil {
		return err
	}

	return nil
}

func (db *Db) exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func hash(bv []byte) string {
	h := sha1.New()
	h.Write(bv)
	return hex.EncodeToString(h.Sum(nil))
}

func compress(content []byte) (bytes.Buffer, error) {
	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	_, err := w.Write(content)
	if err != nil {
		return bytes.Buffer{}, err
	}
	err = w.Close()
	if err != nil {
		return bytes.Buffer{}, err
	}

	return buf, nil
}

func decompress(data []byte) ([]byte, error) {
	reader, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return io.ReadAll(reader)
}
