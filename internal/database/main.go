package database

import "github.com/mundanelizard/envi/internal/blob"

type Db struct {
}

func (db *Db) Store(blob *blob.Blob) error {
	return nil
}

func New(dir string) *Db {

	return &Db{}
}
