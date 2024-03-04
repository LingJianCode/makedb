package datastore

import (
	"io"
	"os"
	"path/filepath"
	"sync"
)

const (
	ACTIVE_FILE_NAME  = "makedb.active"
	OLD_FILE_PREFIX   = "old."
	MERGE_FILE_PREFIX = "merge."
	HINT_FILE_PREFIX  = "hint."
	PERM              = 0664
)

type DataStore struct {
	Keydir     map[string]*KeydirElement
	ActiveFile *DataFile
	FileList   []*os.File
	Mu         sync.Mutex
}

func NewDataStore(path string) (*DataStore, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}
	df, err := OpenActiveFile(path)
	if err != nil {
		return nil, err
	}
	ds := &DataStore{Keydir: map[string]*KeydirElement{}, ActiveFile: df, FileList: []*os.File{}}

	ds.UpdateKeydirFromFile(ds.ActiveFile)

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	files, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}
	f.Close()
	if err != nil {
		return nil, err
	}
	absp, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.Name() == ACTIVE_FILE_NAME {
			continue
		}
		fd, err := os.OpenFile(filepath.Join(absp, file.Name()), os.O_RDONLY, PERM)
		if err != nil {
			return nil, err
		}
		ds.FileList = append(ds.FileList, fd)
		err = ds.UpdateKeydirFromFile(NewDataFile(fd, 0))
		if err != nil {
			return nil, err
		}
	}
	return ds, nil
}

func (ds *DataStore) UpdateKeydirFromFile(df *DataFile) error {
	var off int64 = 0
	for {
		e, err := df.Read(off)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		len := int64(EntryHeader + e.KeySize + e.ValueSize)
		ds.Keydir[string(e.Key)] = NewKeydirElement(df, off, uint32(len), e.Tstamp)
		off += len
	}
	return nil
}

func (ds *DataStore) Get(key []byte) ([]byte, error) {
	ds.Mu.Lock()
	defer ds.Mu.Unlock()
	ke, ok := ds.Keydir[string(key)]
	if !ok {
		return nil, ErrKeyNotExist
	}
	e, err := ke.DataFile.ReadAt(ke.ValuePos, ke.ValueSize)
	if err != nil {
		return nil, err
	}
	return e.Value, nil
}

func (ds *DataStore) Put(key []byte, value []byte) error {
	ds.Mu.Lock()
	defer ds.Mu.Unlock()
	e := NewEntry(key, value)
	// record off before write
	off := ds.ActiveFile.TailOffset
	n, err := ds.ActiveFile.WriteAt(e)
	if err != nil {
		return err
	}
	ds.Keydir[string(key)] = NewKeydirElement(ds.ActiveFile, int64(off), uint32(n), e.Tstamp)
	return nil
}

func (ds *DataStore) Close() {
	ds.Mu.Lock()
	defer ds.Mu.Unlock()
	ds.ActiveFile.File.Close()
	for _, fd := range ds.FileList {
		fd.Close()
	}
}

func (ds *DataStore) RotateActiveFile() error {
	ds.Mu.Lock()
	defer ds.Mu.Unlock()
	return nil
}
