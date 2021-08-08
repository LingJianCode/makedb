package makedb

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"time"
)

type Entry struct {
	Crc       uint32
	Timestamp int64
	Ksz       int64
	ValueSz   int64
	Key       string
	Value     string
	flag      int
}

func NewEntryAndCalc(key string, value string) Entry {
	data := Entry{
		Timestamp: time.Now().Unix(),
		Key:       key,
		Value:     value,
		flag:      0,
	}
	data.calc()
	return data
}

func NewEntry() Entry {
	return Entry{}
}

func (e *Entry) calcCrc() {
	i := crc32.NewIEEE()
	str := fmt.Sprintf("%d%d%d%s%s", e.Timestamp, e.Ksz, e.ValueSz, e.Key, e.Value)
	io.WriteString(i, str)
	e.Crc = i.Sum32()
	//fmt.Println(i.Sum32())
}

func (e *Entry) calcKsz() {
	e.Ksz = int64(len(e.Key))
	//fmt.Println(d.Ksz)
}

func (e *Entry) calcValueSz() {
	e.ValueSz = int64(len(e.Value))
	//fmt.Println(d.ValueSz)
}

func (e *Entry) calc() {
	e.calcKsz()
	e.calcValueSz()
	e.calcCrc()
}

//采用gob序列化
func (e *Entry) Marshal() []byte {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(e)
	if err != nil {
		log.Println(err)
		return nil
	}
	return buf.Bytes()
}

//采用gob反序列化
func (e *Entry) Unmarshal(i []byte) {
	dec := gob.NewDecoder(bytes.NewBuffer(i))
	err := dec.Decode(&e)
	if err != nil {
		log.Fatal(err)
	}
}
