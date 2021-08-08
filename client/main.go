package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"makedb"
	"net"
	"os"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	conn, err := net.Dial("tcp", "localhost:4444")
	defer conn.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("connected...")
	go printRes(os.Stdout, conn)
	sendCommand(conn)
}

func printRes(dst io.Writer, src io.Reader) {
	reader := bufio.NewReader(src)
	for {
		ress, err := makedb.Resp(reader)
		if err != nil {
			log.Println(err)
			//避免由于服务端关闭服务导致客户端goroutine无法结束
			break
		}
		for _, res := range ress {
			io.WriteString(dst, res)
		}
	}
}

func sendCommand(conn net.Conn) {
	defer conn.Close()
	var op, key, value, msg string

	var count int
	for true {
		count = 0
		op, key, value, msg = "", "", "", ""
		//fmt.Print(">")
		fmt.Scanf("%s %s %s", &op, &key, &value)
		if op == "exit" {
			break
		}
		if op != "" {
			count++
			msg += fmt.Sprintf("$%d\r\n", len(op))
			msg += op + "\r\n"
		}
		if key != "" {
			count++
			msg += fmt.Sprintf("$%d\r\n", len(key))
			msg += key + "\r\n"
		}
		if value != "" {
			count++
			msg += fmt.Sprintf("$%d\r\n", len(value))
			msg += value + "\r\n"
		}
		cmd := fmt.Sprintf("*%d\r\n%s", count, msg)
		//fmt.Println(cmd)
		io.WriteString(conn, cmd)
	}
}
