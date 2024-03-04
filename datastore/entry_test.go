package datastore

import (
	"fmt"
	"hash/crc32"
	"os"
	"testing"
)

const (
	KEY   = "ling"
	VALUE = "jian123"
	FILE  = "./data/makedb.active"
)

func TestEncode(t *testing.T) {
	// fmt.Println(crc32.ChecksumIEEE([]byte("lingjian")))
	// tt := time.Now().Unix()
	// fmt.Println(tt, uint64(tt), reflect.TypeOf(tt))
	// buf := make([]byte, EntryHeader)
	// binary.BigEndian.PutUint64(buf, 1)
	// fmt.Println(buf)
	// binary.LittleEndian.PutUint64(buf, 1)
	// fmt.Println(buf)

	e := NewEntry([]byte(KEY), []byte(VALUE))
	b, _ := e.Encode()
	f, err := os.OpenFile(FILE, os.O_RDWR|os.O_CREATE, 0644)
	defer f.Close()
	if err != nil {
		t.Fail()
	}
	f.Write(b)
}

func TestDecode(t *testing.T) {
	f, _ := os.OpenFile(FILE, os.O_RDWR|os.O_CREATE, 0644)
	defer f.Close()
	buf := make([]byte, EntryHeader)
	n, err := f.Read(buf)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	e, _ := Decode(buf)
	e.Key = make([]byte, e.KeySize)
	e.Value = make([]byte, e.ValueSize)
	n, _ = f.ReadAt(e.Key, EntryHeader)
	f.ReadAt(e.Value, int64(EntryHeader+n))
	fmt.Println(string(e.Key))
	fmt.Println(string(e.Value))
	data := make([]byte, EntryHeader+e.KeySize+e.ValueSize)
	copy(data[0:EntryHeader], buf[:])
	copy(data[EntryHeader:EntryHeader+e.KeySize], e.Key)
	copy(data[EntryHeader+e.KeySize:], e.Value)
	if string(e.Key) != KEY || string(e.Value) != VALUE || e.Crc != crc32.ChecksumIEEE(data[4:]) {
		fmt.Println(e.Crc, crc32.ChecksumIEEE(data[4:]))
		t.Fail()
	}
}
