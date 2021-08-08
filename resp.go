package makedb

import (
	"bufio"
	"errors"
	"log"
	"strconv"
)

func Resp(reader *bufio.Reader) ([]string, error) {
	flag, err := reader.ReadByte()
	var res []string
	if err != nil {
		return nil, err
	}
	switch flag {
	case '-':
		c := handleError(reader)
		res = append(res, c)
	case '+':
		c := handleSimpleString(reader)
		res = append(res, c)
	case ':':
		c := handleInt(reader)
		res = append(res, c)
	case '$':
		c, _ := handleString(reader)
		res = append(res, c)
	case '*':
		res = handleArray(reader)
	default:
		return nil, errors.New("-unkown handle")
	}
	return res, nil
}

func handleError(reader *bufio.Reader) string {
	//利用\n来做标识，以后需要对\r做验证
	msg, _ := reader.ReadBytes('\n')
	//fmt.Println(string(msg[:len(msg)-2]))
	return string(msg[:len(msg)-2])
}

func handleSimpleString(reader *bufio.Reader) string {
	msg, _ := reader.ReadBytes('\n')
	//fmt.Println(string(msg[:len(msg)-2]))
	return string(msg[:len(msg)-2])
}

func handleInt(reader *bufio.Reader) string {
	msg, _ := reader.ReadBytes('\n')
	return string(msg[:len(msg)-2])
}

func handleString(reader *bufio.Reader) (string, error) {
	msg, _ := reader.ReadBytes('\n')
	l, _ := strconv.Atoi(string(msg[:len(msg)-2]))
	//fmt.Println(reflect.TypeOf(len))
	res := make([]byte, l+2)
	_, err := reader.Read(res)
	if err != nil {
		log.Println(err)
		return "", err
	}
	//fmt.Println(string(res[:len(res)-2]))
	return string(res[:len(res)-2]), nil
}

func handleArray(reader *bufio.Reader) []string {
	msg, _ := reader.ReadBytes('\n')
	l, _ := strconv.Atoi(string(msg[:len(msg)-2]))
	var res []string
	for i := 0; i < l; i++ {
		flag, _ := reader.ReadByte()
		switch flag {
		case '-':
			handleError(reader)
		case '+':
			handleSimpleString(reader)
		case ':':
			handleInt(reader)
		case '$':
			x, _ := handleString(reader)
			res = append(res, x)
		case '*':
			//handleArray(reader)
			log.Fatal("not support")
		}
	}
	//fmt.Println(res)
	return res
}
