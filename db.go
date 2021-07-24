package makedb

import "fmt"

type Dbms struct {
	db map[string]string
}

func (d Dbms) Get(key string) (string, error) {
	value := d.db[key]
	if value == "" {
		return "", fmt.Errorf("value is nil")
	}
	return value, nil
}

func (d Dbms) Set(key string, value string) error {
	d.db[key] = value
	return nil
}

func (d Dbms) Del(key string) error {
	delete(d.db, key)
	return nil
}

func NewDbms() Dbms {
	return Dbms{make(map[string]string)}
}
