package datastore

import (
	"hash/crc32"
	"os"
	"path/filepath"
	"sync"
)

type DataFile struct {
	File       *os.File
	TailOffset int64
	// maybe we should use sync.RWMutex
	Mu *sync.Mutex
}

func NewDataFile(fd *os.File, off int64) *DataFile {
	return &DataFile{
		File:       fd,
		TailOffset: off,
		Mu:         &sync.Mutex{},
	}
}

func OpenActiveFile(path string) (*DataFile, error) {
	absp, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	activeFilePath := filepath.Join(absp, ACTIVE_FILE_NAME)
	fd, err := os.OpenFile(activeFilePath, os.O_RDWR|os.O_CREATE, PERM)
	if err != nil {
		return nil, err
	}
	fi, err := fd.Stat()
	if err != nil {
		return nil, err
	}
	return NewDataFile(fd, fi.Size()), nil
}

func (df *DataFile) WriteAt(e *Entry) (int, error) {
	df.Mu.Lock()
	defer df.Mu.Unlock()
	buf, err := e.Encode()
	if err != nil {
		return 0, err
	}
	n, err := df.File.WriteAt(buf, df.TailOffset)
	if err != nil {
		return 0, err
	}
	df.File.Sync()
	df.TailOffset += int64(n)
	return n, nil
}

func (df *DataFile) ReadAt(off int64, len uint32) (*Entry, error) {
	df.Mu.Lock()
	defer df.Mu.Unlock()
	buf := make([]byte, len)
	_, err := df.File.ReadAt(buf, off)
	if err != nil {
		return nil, err
	}
	e, err := Decode(buf)
	if err != nil {
		return nil, err
	}
	if e.Crc != crc32.ChecksumIEEE(buf[4:]) {
		return nil, ErrKeyCheckSumWrong
	}
	e.Key = make([]byte, e.KeySize)
	e.Value = make([]byte, e.ValueSize)
	copy(e.Key, buf[EntryHeader:EntryHeader+e.KeySize])
	copy(e.Value, buf[EntryHeader+e.KeySize:])
	return e, nil
}

// entry's length = EntryHeader + entry's KeySize + entry's ValueSize
// This function is used when instance starting.
func (df *DataFile) Read(off int64) (*Entry, error) {
	df.Mu.Lock()
	defer df.Mu.Unlock()
	buf := make([]byte, EntryHeader)
	_, err := df.File.ReadAt(buf, off)
	if err != nil {
		return nil, err
	}
	e, err := Decode(buf)
	if err != nil {
		return nil, err
	}
	e.Key = make([]byte, e.KeySize)
	e.Value = make([]byte, e.ValueSize)
	n, err := df.File.ReadAt(e.Key, off+EntryHeader)
	if err != nil {
		return nil, err
	}
	df.File.ReadAt(e.Value, off+EntryHeader+int64(n))
	if err != nil {
		return nil, err
	}
	data := make([]byte, EntryHeader+e.KeySize+e.ValueSize)
	copy(data[0:EntryHeader], buf[:])
	copy(data[EntryHeader:EntryHeader+e.KeySize], e.Key)
	copy(data[EntryHeader+e.KeySize:], e.Value)
	if e.Crc != crc32.ChecksumIEEE(data[4:]) {
		return nil, ErrKeyCheckSumWrong
	}
	return e, nil
}
