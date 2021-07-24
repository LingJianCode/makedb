package main

import "fmt"

func main() {
	db := Dbms{make(map[string]string)}
	fmt.Println("connected...")
Loop:
	for true {
		var op, key, value string
		fmt.Print(">")
		fmt.Scanln(&op, &key, &value)
		switch op {
		case "get":
			val, err := db.get(key)
			if err != nil {
				fmt.Print(err)
			}
			fmt.Println(val)
		case "set":
			err := db.set(key, value)
			if err != nil {
				fmt.Println("set ", key, "failed")
			}
			fmt.Println("ok")
		case "del":
			err := db.del(key)
			if err != nil {
				fmt.Println("del ", key, "failed")
			}
			fmt.Println("ok")
		case "exit":
			fmt.Println("exit...")
			break Loop
		default:
			fmt.Println("plase input set/get/del xxx or exit.")
		}
	}
}

type Dbms struct {
	db map[string]string
}

func (d Dbms) get(key string) (string, error) {
	value := d.db[key]
	if value == "" {
		return "", fmt.Errorf("value is nil")
	}
	return value, nil
}

func (d Dbms) set(key string, value string) error {
	d.db[key] = value
	return nil
}

func (d Dbms) del(key string) error {
	delete(d.db, key)
	return nil
}
