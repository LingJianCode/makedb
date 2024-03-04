package datastore

import (
	"encoding/binary"
	"hash/crc32"
	"time"
)

const EntryHeader = 20

type Entry struct {
	Crc       uint32
	Tstamp    uint64
	KeySize   uint32
	ValueSize uint32
	Key       []byte
	Value     []byte
}

func NewEntry(key, value []byte) *Entry {
	return &Entry{
		Tstamp:    uint64(time.Now().Unix()),
		KeySize:   uint32(len(key)),
		ValueSize: uint32(len(value)),
		Key:       key,
		Value:     value,
	}
}

func (e *Entry) getSize() int64 {
	return int64(EntryHeader + e.KeySize + e.ValueSize)
}

func (e *Entry) Encode() ([]byte, error) {
	buf := make([]byte, e.getSize())
	binary.LittleEndian.PutUint64(buf[4:12], e.Tstamp)
	binary.LittleEndian.PutUint32(buf[12:16], e.KeySize)
	binary.LittleEndian.PutUint32(buf[16:20], e.ValueSize)
	copy(buf[EntryHeader:EntryHeader+e.KeySize], e.Key)
	copy(buf[EntryHeader+e.KeySize:], e.Value)
	binary.LittleEndian.PutUint32(buf[0:4], crc32.ChecksumIEEE(buf[4:]))
	return buf, nil
}

func Decode(buf []byte) (*Entry, error) {
	e := Entry{}
	e.Crc = binary.LittleEndian.Uint32(buf[0:4])
	e.Tstamp = binary.LittleEndian.Uint64(buf[4:12])
	e.KeySize = binary.LittleEndian.Uint32(buf[12:16])
	e.ValueSize = binary.LittleEndian.Uint32(buf[16:20])
	return &e, nil
}
