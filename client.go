package makedb

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

type Client struct {
	conn net.Conn
}

func NewClient(host string, port string) Client {
	conn, _ := net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
	return Client{conn: conn}
}

func (c *Client) Set(key string, value string) []string {
	var msg string
	count := 3
	msg += "$3\r\n"
	msg += "SET" + "\r\n"
	msg += fmt.Sprintf("$%d\r\n", len(key))
	msg += key + "\r\n"
	msg += fmt.Sprintf("$%d\r\n", len(value))
	msg += value + "\r\n"
	cmd := fmt.Sprintf("*%d\r\n%s", count, msg)
	//log.Println(cmd)
	io.WriteString(c.conn, cmd)
	reader := bufio.NewReader(c.conn)
	ress, err := Resp(reader)
	if err != nil {
		log.Println(err)
	}
	return ress
}

func (c *Client) Get(key string) []string {
	var msg string
	count := 2
	msg += "$3\r\n"
	msg += "GET" + "\r\n"
	msg += fmt.Sprintf("$%d\r\n", len(key))
	msg += key + "\r\n"
	cmd := fmt.Sprintf("*%d\r\n%s", count, msg)
	//log.Println(cmd)
	io.WriteString(c.conn, cmd)
	reader := bufio.NewReader(c.conn)
	ress, err := Resp(reader)
	if err != nil {
		log.Println(err)
	}
	return ress
}

func (c *Client) Del(key string) []string {
	var msg string
	count := 2
	msg += "$3\r\n"
	msg += "DEL" + "\r\n"
	msg += fmt.Sprintf("$%d\r\n", len(key))
	msg += key + "\r\n"
	cmd := fmt.Sprintf("*%d\r\n%s", count, msg)
	//log.Println(cmd)
	io.WriteString(c.conn, cmd)
	reader := bufio.NewReader(c.conn)
	ress, err := Resp(reader)
	if err != nil {
		log.Println(err)
	}
	return ress
}
