package makedb

import (
	"fmt"
	"testing"
)

var (
	PATH     = "./data"
	TEST_SET = map[string]string{
		"ling":  "jian123",
		"Cheng": "Du123",
		"Si":    "Chuan123",
		"good":  "bad123",
		"true":  "false123",
	}
)

func TestInit(t *testing.T) {
	ds, err := Init(PATH)
	defer ds.Close()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	// _, err := os.Stat(path)
	// if os.IsNotExist(err) {
	// 	err = os.MkdirAll(path, os.ModePerm)
	// 	if err != nil {
	// 		t.Fail()
	// 	}
	// }
	// absp, _ := filepath.Abs(path)
	// fmt.Println(filepath.Join(absp, ACTIVE_FILE_NAME))
	// fmt.Println("--")
	// f, _ := os.Open(path)
	// files, err := f.Readdir(-1)
	// f.Close()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// for _, file := range files {
	// 	fmt.Println(filepath.Join(absp, file.Name()))
	// }
}

func TestPut(t *testing.T) {
	ds, err := Init(PATH)
	defer ds.Close()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	for k, v := range TEST_SET {
		fmt.Println(k, v)
		err = ds.Put([]byte(k), []byte(v))
		if err != nil {
			fmt.Println(err)
			t.Fail()
		}
	}
}

func TestGET(t *testing.T) {
	ds, err := Init(PATH)
	defer ds.Close()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	for k, v := range TEST_SET {
		tv, err := ds.Get([]byte(k))
		fmt.Println(k, string(tv))
		if err != nil {
			fmt.Println(err)
			t.Fail()
		}
		if string(tv) != v {
			t.Fail()
		}
	}
}
