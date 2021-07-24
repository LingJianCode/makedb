package makedb

import (
	"bufio"
	"fmt"
	"strings"
	"testing"
)

func TestResp(t *testing.T) {
	testCmdArray := []string{
		"-error\r\n",
		"+OK\r\n",
		":100\r\n",
		"$3\r\nGET\r\n",
		"*3\r\n$3\r\nSET\r\n$4\r\nling\r\n$4\r\njian\r\n",
		"*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n",
		"*3\r\n$3\r\nset\r\n$5\r\nlevel\r\n$4\r\ngood\r\n",
	}
	for _, cmd := range testCmdArray {
		r := strings.NewReader(cmd)
		reader := bufio.NewReader(r)
		a, _ := Resp(reader)
		fmt.Println(a)
	}
}
