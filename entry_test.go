package study

import (
	"encoding/binary"
	"fmt"
	"os"
	"testing"
)

func TestEncode(t *testing.T) {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, 1)
	fmt.Println(buf)
	binary.LittleEndian.PutUint64(buf, 1)
	fmt.Println(buf)

	key := "ling"
	value := "jianttt"
	e := NewEntry([]byte(key), []byte(value))
	b, _ := e.Encode()
	fmt.Println(b)
	f, _ := os.OpenFile("./file.db", os.O_RDWR|os.O_CREATE, 0644)
	defer f.Close()

	f.Write(b)

}

func TestDecode(t *testing.T) {

	f, _ := os.OpenFile("./file.db", os.O_RDWR|os.O_CREATE, 0644)
	defer f.Close()
	buf := make([]byte, 8)
	n, err := f.Read(buf)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	key := make([]byte, binary.BigEndian.Uint32(buf[0:4]))
	value := make([]byte, binary.BigEndian.Uint32(buf[4:8]))
	n, _ = f.ReadAt(key, 8)
	f.ReadAt(value, int64(8+n))
	fmt.Println(string(key))
	fmt.Println(string(value))
}
