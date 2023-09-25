package makedb

import (
	"fmt"
	"testing"
)

func TestInit(t *testing.T) {
	path := "./data"
	key1 := []byte("ling")
	value1 := []byte("jian123")
	key2 := []byte("cheng")
	value2 := []byte("du")
	ds, err := Init(path)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	ds.put(key1, value1)
	ds.put(key2, value2)
	fmt.Println("---")
	v, _ := ds.get(key2)
	fmt.Println(string(v))
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
