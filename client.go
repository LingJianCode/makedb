package makedb

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

type Client struct {
	op    string
	key   string
	value string
}

func (c *Client) sendCmd(conn net.Conn) []string {
	defer conn.Close()
	var msg string
	var count int
	count = 0
	if c.op == "exit" {
		os.Exit(0)
	}
	if c.op != "" {
		count++
		msg += fmt.Sprintf("$%d\r\n", len(c.op))
		msg += c.op + "\r\n"
	}
	if c.key != "" {
		count++
		msg += fmt.Sprintf("$%d\r\n", len(c.key))
		msg += c.key + "\r\n"
	}
	if c.value != "" {
		count++
		msg += fmt.Sprintf("$%d\r\n", len(c.value))
		msg += c.value + "\r\n"
	}
	cmd := fmt.Sprintf("*%d\r\n%s", count, msg)
	//log.Println(cmd)
	io.WriteString(conn, cmd)
	reader := bufio.NewReader(conn)
	ress, err := Resp(reader)
	if err != nil {
		log.Println(err)
	}
	return ress
}
