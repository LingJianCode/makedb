package makedb

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"hash/crc32"
	"io"

	//"github.com/vmihailenco/msgpack"
	"log"
	"os"
	"testing"
)

func TestCalcCrc(t *testing.T) {
	a := NewEntryAndCalc("ling", "jian")
	fmt.Println(a)
}

//
//func TestStore(T *testing.T)  {
//	//使用msgpack序列化
//	a := NewData("ling", "jian")
//	fmt.Println(a)
//	fd, _ := os.Create("./data/store")
//	data, _ := msgpack.Marshal(a)
//	n, _ := fd.Write(data)
//	fmt.Println(n)
//	fd.Close()
//}
//
//func TestReadFromFile(T *testing.T){
//	//使用msgpack反序列化
//	fd, _ := os.Open("./data/store")
//	a := make([]byte, 100)
//	var b Data
//	//{1268759064 1627305646 4 4 ling jian 0}
//	fd.Read(a)
//	msgpack.Unmarshal(a, &b)
//	fmt.Println(b)
//}

func TestGobWriteFile(T *testing.T) {
	//使用gob进行序列化
	a := NewEntryAndCalc("ling", "jian")
	fmt.Println(a)
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(a)
	data := buf.Bytes()
	if err != nil {
		log.Fatal(err)
		return
	}
	fd, _ := os.Create("./data/store.001")
	fi, _ := fd.Stat()
	fmt.Println(fi)
	n, _ := fd.Write(data)
	fmt.Println(n)
	fd.Close()
}

func TestGobReadFile(T *testing.T) {
	fd, _ := os.Open("./data/store.001")
	a := make([]byte, 111)
	var b Entry
	//{3883620618 1627306442 4 4 ling jian 0}
	fd.Read(a)
	dec := gob.NewDecoder(bytes.NewBuffer(a))
	err := dec.Decode(&b)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(b)
}

func TestData_Marshal(t *testing.T) {
	//109
	a := NewEntryAndCalc("go", "lang")
	fd, _ := os.OpenFile("./data/store.001", os.O_APPEND, 0644)
	fi, _ := fd.Stat()
	fmt.Println("write before:", fi.Size())
	n, _ := fd.Write(a.Marshal())

	fmt.Println(n)
	//fi2, _ := fd.Stat()
	fmt.Println("write After:", fi.Size())
	fd.Close()

	crc := crc32.ChecksumIEEE([]byte("e.Meta.Value"))
	fmt.Println(crc)

	i := crc32.NewIEEE()
	io.WriteString(i, "e.Meta.Value")
	fmt.Println(i.Sum32())
}

func TestData_Unmarshal(t *testing.T) {

	a := NewEntry()
	fd, _ := os.Open("./data/store.001")

	b := make([]byte, 111)
	fd.Read(b)
	a.Unmarshal(b)
	fmt.Println(a)

	c := make([]byte, 109)
	fd.ReadAt(c, 111)
	a.Unmarshal(c)
	fd.Close()
	fmt.Println(a)
}
