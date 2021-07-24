package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)
import "makedb"

var db = makedb.NewDbms()

func main() {
	//打印日志在文件中的行数
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	listener, err := net.Listen("tcp", "localhost:4444")
	fmt.Println("start...")
	defer listener.Close()
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for true {
		handle, err := makedb.Resp(reader)
		if err != nil {
			if err == io.EOF {
				log.Println(err)
				break
			}
			log.Println(err)
			//避免由于客户端关闭导致服务端goroutine无法结束
			break
		}
		//避免在客户端发送数据前从连接取数据，在客户端发送数据前读取数据会导致读取的数据为nil
		if handle == nil {
			continue
		}
		fmt.Println(handle)
		switch strings.ToUpper(handle[0]) {
		case "GET":
			key := handle[1]
			val, err := db.Get(key)
			if err != nil {
				log.Println(err)
				break
			}
			if val == "" {
				io.WriteString(conn, "-not exists\r\n")
				break
			}
			//conn.Write([]byte(val))
			//fmt.Fprint(conn , val+"\r\n")
			io.WriteString(conn, val+"\r\n")
		case "SET":
			key := handle[1]
			value := handle[2]
			err := db.Set(key, value)
			if err != nil {
				io.WriteString(conn, "-set "+key+" failed\r\n")
			}
			io.WriteString(conn, "+OK\r\n")
		case "DEL":
			key := handle[1]
			err := db.Del(key)
			if err != nil {
				fmt.Println("del ", key, "failed")
				io.WriteString(conn, "-del "+key+" failed\r\n")
			}
			io.WriteString(conn, "+OK\r\n")
		case "EXIT":
			fmt.Println("exit...")
			break
		default:
			fmt.Println("plase input set/get/del xxx or exit.")
		}
	}
}
