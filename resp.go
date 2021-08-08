package makedb

import (
	"bufio"
	"errors"
	"log"
	"strconv"
)

func Resp(reader *bufio.Reader) ([]string, error) {
	flag, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	var res []string
	switch flag {
	case '-':
		c := analysisError(reader)
		res = append(res, c)
	case '+':
		c := analysisSimpleString(reader)
		res = append(res, c)
	case ':':
		c := analysisInt(reader)
		res = append(res, c)
	case '$':
		c, _ := analysisString(reader)
		res = append(res, c)
	case '*':
		res = analysisArray(reader)
	default:
		return nil, errors.New("-unkown handle")
	}
	return res, nil
}

func analysisError(reader *bufio.Reader) string {
	//利用\n来做标识，以后需要对\r做验证
	msg, _ := reader.ReadBytes('\n')
	//fmt.Println(string(msg[:len(msg)-2]))
	return string(msg[:len(msg)-2])
}

func analysisSimpleString(reader *bufio.Reader) string {
	msg, _ := reader.ReadBytes('\n')
	//fmt.Println(string(msg[:len(msg)-2]))
	return string(msg[:len(msg)-2])
}

func analysisInt(reader *bufio.Reader) string {
	msg, _ := reader.ReadBytes('\n')
	return string(msg[:len(msg)-2])
}

func analysisString(reader *bufio.Reader) (string, error) {
	msg, _ := reader.ReadBytes('\n')
	//len-2,2 = len("\r\n")，去掉\r\n
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

func analysisArray(reader *bufio.Reader) []string {
	msg, _ := reader.ReadBytes('\n')
	l, _ := strconv.Atoi(string(msg[:len(msg)-2]))
	var res []string
	for i := 0; i < l; i++ {
		flag, _ := reader.ReadByte()
		switch flag {
		case '-':
			analysisError(reader)
		case '+':
			analysisSimpleString(reader)
		case ':':
			analysisInt(reader)
		case '$':
			x, _ := analysisString(reader)
			res = append(res, x)
		case '*':
			//handleArray(reader)
			log.Fatal("not support")
		}
	}
	//fmt.Println(res)
	return res
}

func RespSendError(e string) {

}

func RespSendSimpleString(ss string) {

}

func RespSendInt(i int) {

}

func RespSendString(s string) {

}

func RespSendArray(sa []string) {

}
