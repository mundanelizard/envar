package lockfile

import (
	"errors"
	"os"
)

var (
	ErrMissingParent = errors.New("lockfile: missing parent")
	ErrNoPermission  = errors.New("lockfile: no permission")
	ErrStaleLock     = errors.New("lockfile: stale lock")
	ErrLockDenied    = errors.New("lockfile: lock denied")
)

type Lockfile struct {
	filePath string
	lockPath string
	lock     *os.File
}

func New(filePath string) *Lockfile {
	return &Lockfile{
		filePath: filePath,
		lockPath: filePath + ".lock",
	}
}

func (l *Lockfile) Hold() error {
	if l.lock != nil {
		return nil
	}

	f, err := os.OpenFile(l.lockPath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0755)

	if err != nil {
		switch {
		case os.IsExist(err):
			return ErrLockDenied
		case os.IsNotExist(err):
			return ErrMissingParent
		case os.IsPermission(err):
			return ErrNoPermission
		}

		return err
	}

	l.lock = f

	return nil
}

func (l *Lockfile) Write(data []byte) error {
	if l.lock == nil {
		return ErrStaleLock
	}

	_, err := l.lock.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (l *Lockfile) Commit() error {
	if l.lock == nil {
		return ErrStaleLock
	}

	err := l.lock.Close()
	if err != nil {
		return err
	}

	err = os.Rename(l.lockPath, l.filePath)
	if err != nil {
		return err
	}

	return nil
}
