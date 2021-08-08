package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"makedb"
	"os"
	"strings"
	"testing"
)

func TestGet(t *testing.T) {
	var op, key, value string
	fmt.Scanln(&op, &key, &value)
	fmt.Println(op, key, value)
}

func TestStoreRead(t *testing.T) {
	fd, _ := os.Open("../data/store.002")
	defer fd.Close()
	entry := makedb.NewEntry()
	stream, err := ioutil.ReadAll(fd)
	if err != nil {
		fmt.Println(err)
	}
	dec := gob.NewDecoder(bytes.NewBuffer(stream))
	err1 := dec.Decode(&entry)
	if err1 != nil {
		if err1 == io.EOF {
			log.Println(err1)
		} else {
			log.Fatal(err1)
		}
	}
	fmt.Println(entry)

}

func TestStoreScan(t *testing.T) {
	fd, _ := os.Open("../data/store.002")
	defer fd.Close()
	r := bufio.NewReader(fd)
	for {
		entry := makedb.NewEntry()
		slice, err := r.ReadBytes('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		dec := gob.NewDecoder(bytes.NewBuffer(slice[:len(slice)-1]))
		err1 := dec.Decode(&entry)
		if err1 != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(entry.Value)
	}
}

func TestBufioRead(T *testing.T) {
	r := strings.NewReader("hello\nworld\n")
	reader := bufio.NewReader(r)
	for {
		str, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(str[:len(str)-1]))
	}
}

func TestListDir(T *testing.T) {
	fileInfos, err := ioutil.ReadDir("../data")
	if err != nil {
		fmt.Println(err)
	}
	for _, fi := range fileInfos {
		//fmt.Println(fi.Name())
		filename := fi.Name()
		if strings.HasPrefix(filename, "store1") {
			fmt.Println(filename)
		}
	}
}
