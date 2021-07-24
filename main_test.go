package main

import (
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {
	var op, key, value string
	fmt.Scanln(&op, &key, &value)
	fmt.Println(op, key, value)
}
