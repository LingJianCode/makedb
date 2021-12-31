package main

import (
	"bufio"
	"fmt"
	"github.com/robfig/cron"
	"io"
	"log"
	"net"
	"os"
	"strings"
)
import "makedb"

var keydir = make(makedb.KeyDir)

func main() {
	//打印日志在文件中的行数
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	var fd *os.File
	fd = makedb.InitKeyDir(&keydir, makedb.StoragePath)
	defer fd.Close()
	if fd == nil {
		fd, _ = os.OpenFile(makedb.ActiveFilePath, os.O_APPEND|os.O_CREATE, 0644)
	}

	//定时merge
	c := cron.New()
	spec := "0 */10 * * * ?"
	c.AddFunc(spec, func() {
		merge(fd)
	})
	c.Start()

	listener, err := net.Listen("tcp", "0.0.0.0:4444")
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
		go handleConn(conn, fd)
	}
}

func handleConn(conn net.Conn, fd *os.File) {
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
			val := get(key)
			if val == "" {
				io.WriteString(conn, "-not exists\r\n")
				break
			}
			//conn.Write([]byte(val))
			//fmt.Fprint(conn , val+"\r\n")
			//io.WriteString(conn, val+"\r\n")
			//fmt.Println(fmt.Sprintf("$%d%s\r\n", len(val), val))
			io.WriteString(conn, fmt.Sprintf("$%d\r\n%s\r\n", len(val), val))
		case "SET":
			key := handle[1]
			value := handle[2]
			err := set(&fd, key, value)
			if err != nil {
				io.WriteString(conn, "-set "+key+" failed\r\n")
			} else {
				io.WriteString(conn, "+OK\r\n")
			}
		case "DEL":
			key := handle[1]
			del(key)
			io.WriteString(conn, "+OK\r\n")
		case "EXIT":
			fmt.Println("exit...")
			break
		case "COMMAND":
			io.WriteString(conn, "+OK\r\n")
			break
		default:
			fmt.Println("plase input set/get/del xxx or exit.")
			io.WriteString(conn, "+OK\r\n")
			break
		}
	}
}
