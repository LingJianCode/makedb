package study

import "encoding/binary"

type Entry struct {
	KeySize   uint32
	ValueSize uint32
	Key       []byte
	Value     []byte
}

func NewEntry(key, value []byte) *Entry {
	return &Entry{
		KeySize:   uint32(len(key)),
		ValueSize: uint32(len(value)),
		Key:       key,
		Value:     value,
	}
}

func (e *Entry) getSize() int64 {
	return int64(e.KeySize + e.ValueSize + 8)
}

func (e *Entry) Encode() ([]byte, error) {
	buf := make([]byte, e.getSize())
	binary.BigEndian.PutUint32(buf[0:4], e.KeySize)
	binary.BigEndian.PutUint32(buf[4:8], e.ValueSize)
	copy(buf[8:8+e.KeySize], e.Key)
	copy(buf[8+e.KeySize:], e.Value)
	return buf, nil
}

func Decode(buf []byte) (*Entry, error) {
	e := Entry{}
	e.KeySize = binary.BigEndian.Uint32(buf[0:4])
	e.ValueSize = binary.BigEndian.Uint32(buf[4:8])
	return &e, nil
}
