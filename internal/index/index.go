package index

import (
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"github.com/mundanelizard/envi/internal/lockfile"
	"hash"
	"os"
)

const (
	RegularMode    = 0100644
	ExecutableMode = 0100755
	MaxPathSize    = 0xfff
)

type Index struct {
	entries  map[string]*Entry
	lockfile *lockfile.Lockfile
	digest   hash.Hash
}

func New(path string) *Index {
	return &Index{
		entries:  map[string]*Entry{},
		lockfile: lockfile.New(path),
	}
}

func (i *Index) Add(path, id string, stat os.FileInfo) {
	e := NewEntry(path, id, stat)
	i.entries[path] = e
}

func (i *Index) WriteUpdates() error {
	err := i.lockfile.Hold()
	if err != nil {
		return err
	}

	i.beginWrite()
	header := []byte("DIRC")
	header = append(header, make([]byte, 8)...)
	binary.BigEndian.PutUint32(header[4:], 2)
	binary.BigEndian.PutUint32(header[8:], uint32(len(i.entries)))
	err = i.write(header)
	if err != nil {
		return err
	}

	for _, entry := range i.entries {
		err = i.write(entry.Bytes())
		if err != nil {
			return err
		}
	}

	err = i.finishWrite()
	if err != nil {
		return err
	}

	return nil
}

func (i *Index) beginWrite() {
	// Create a new SHA1 hash object
	i.digest = sha1.New()
}

func (i *Index) write(data []byte) error {
	// Write the data to the lock file
	err := i.lockfile.Write(data)
	if err != nil {
		return err
	}

	// Update the SHA1 hash with the data
	i.digest.Write(data)

	return nil
}

func (i *Index) finishWrite() error {
	// Write the final SHA1 digest to the lock file
	err := i.lockfile.Write(i.digest.Sum(nil))
	if err != nil {
		return err
	}

	// Commit the changes to the lock file
	err = i.lockfile.Commit()
	if err != nil {
		return err
	}

	return nil
}

type Entry struct {
	ctime     int64
	ctimeNSec int
	mtime     int64
	mtimeNSec int
	dev       int64
	ino       int
	mode      int
	uid       int
	gid       int
	size      int64
	oid       string
	flags     int
	path      string
}

func NewEntry(pathname string, oid string, stat os.FileInfo) *Entry {
	path := pathname
	mode := RegularMode
	if (stat.Mode() & os.ModeType) == 0x8000 {
		mode = ExecutableMode
	}

	flags := len(path)
	if flags > MaxPathSize {
		flags = MaxPathSize
	}

	// todo => swap out for system time in the future if you comeback
	time := stat.ModTime()

	return &Entry{
		ctime:     time.Unix(),
		ctimeNSec: time.Nanosecond(),
		mtime:     time.Unix(),
		mtimeNSec: time.Nanosecond(),
		dev:       time.Unix(),
		ino:       int(time.Unix()),
		mode:      mode,
		uid:       0,
		gid:       0,
		size:      stat.Size(),
		oid:       oid,
		flags:     flags,
		path:      path,
	}
}

const (
	EntryBlock = 8
)

func (e *Entry) Bytes() []byte {
	data := make([]byte, 40)

	// Packing 10 32bit unsigned big-endian numbers into the first bytes
	binary.BigEndian.PutUint32(data[0:], uint32(e.ctime))
	binary.BigEndian.PutUint32(data[4:], uint32(e.ctimeNSec))
	binary.BigEndian.PutUint32(data[8:], uint32(e.mtime))
	binary.BigEndian.PutUint32(data[12:], uint32(e.mtimeNSec))
	binary.BigEndian.PutUint32(data[16:], uint32(e.dev))
	binary.BigEndian.PutUint32(data[20:], uint32(e.ino))
	binary.BigEndian.PutUint32(data[24:], uint32(e.mode))
	binary.BigEndian.PutUint32(data[28:], uint32(e.uid))
	binary.BigEndian.PutUint32(data[32:], uint32(e.gid))
	binary.BigEndian.PutUint32(data[36:], uint32(e.size))

	oid, err := hex.DecodeString(e.oid)
	if err != nil {
		panic(err)
	}

	data = append(data, oid...)
	data = append(data, make([]byte, 2)...)

	binary.BigEndian.PutUint16(data[60:], uint16(e.flags))

	data = append(data, []byte(e.path)...)
	data = append(data, []byte("\x00")...)

	// padding with null bytes
	padding := EntryBlock - (len(data) % EntryBlock)
	data = append(data, make([]byte, padding)...)

	return data
}
