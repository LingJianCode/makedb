package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"testing"
	"time"
)

func TestSendCommand(T *testing.T) {
	conn, err := net.Dial("tcp", "localhost:4444")
	defer conn.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("connected...")
	go mustCopy(os.Stdout, conn)
	conn.Write([]byte("*3\r\n$3\r\nset\r\n$4\r\nling\r\n$4\r\njian\r\n"))
	time.Sleep(1 * time.Second)
	conn.Write([]byte("*2\r\n$3\r\nget\r\n$4\r\nling\r\n"))
	time.Sleep(1 * time.Second)
}
