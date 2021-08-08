package main

import (
	"fmt"
	"log"
	"makedb"
	"os"
	"time"
)

var keydir = make(makedb.KeyDir)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	var fd *os.File
	fd = makedb.InitKeyDir(&keydir, makedb.StoragePath)
	defer fd.Close()
	if fd == nil {
		fd, _ = os.OpenFile(makedb.ActiveFilePath, os.O_APPEND|os.O_CREATE, 0644)
	}

	fmt.Println(fd)
	fmt.Println(keydir["ling"].FileId)

	fmt.Println("connected...")
Loop:
	for true {
		var op, key, value string
		fmt.Print(">")
		fmt.Scanln(&op, &key, &value)
		switch op {
		case "get":
			get(key)
		case "set":
			//fmt.Println(fd, key, value)
			set(fd, key, value)
		case "del":
			del(key)
		case "exit":
			fmt.Println("exit...")
			fd.Close()
			break Loop
		case "merge":
			merge(fd)
		case "rotation":
			fmt.Println(fd)
			rotation(&fd)
			fmt.Println(fd)
		default:
			fmt.Println("plase input set/get/del xxx or exit.")
		}
	}
}

func merge(fd *os.File) {
	mergeFileName := fmt.Sprintf(makedb.MergeFilePreFix+".%d", time.Now().Unix())
	mergeFd, _ := os.OpenFile(makedb.StoragePath+"/"+mergeFileName, os.O_APPEND|os.O_CREATE, 0644)
	hintFileName := fmt.Sprintf("%s_"+mergeFileName, makedb.HintFilePreFIx)
	hintFd, _ := os.OpenFile(makedb.StoragePath+"/"+hintFileName, os.O_APPEND|os.O_CREATE, 0644)

	removeFdMap := make(map[string]*os.File)
	var pos int64 = 0
	for k, v := range keydir {
		if v.FileId != fd {
			data := make([]byte, v.ValueSz)
			_, err := v.FileId.ReadAt(data, v.ValuePos)
			if err != nil {
				log.Fatal(err)
			}
			//写入merge文件
			n, _ := mergeFd.Write(append(data, makedb.Separator))
			mergeFd.Sync()
			//更新keydir
			keydir[k] = makedb.ValueIndex{
				FileId:    mergeFd,
				ValueSz:   v.ValueSz,
				ValuePos:  pos,
				Timestamp: v.Timestamp,
			}

			//写入hint文件
			hint := makedb.Hint{
				Timestamp: time.Now().Unix(),
				Ksz:       int64(len(k)),
				ValueSz:   v.ValueSz,
				ValuePos:  pos,
				Key:       k,
			}
			_, err = hintFd.Write(append(hint.Marshal(), makedb.Separator))
			if err != nil {
				log.Fatal(err)
			}
			hintFd.Sync()
			//更新offset
			pos += int64(n)
			removeFdMap[v.FileId.Name()] = v.FileId

		}
	}

	//清理old file
	for k, v := range removeFdMap {
		fmt.Println(k)
		v.Close()
		err := os.Remove(v.Name())
		if err != nil {
			log.Println(err)
		}
	}
}

func rotation(fd **os.File) {
	(*fd).Close()
	newFileName := fmt.Sprintf(makedb.StoragePath+"/"+makedb.StorageFilePreFix+".%d", time.Now().Unix())
	err := os.Rename(makedb.ActiveFilePath, newFileName)
	if err != nil {
		log.Println(err)
	}
	oldFd, _ := os.Open(newFileName)
	for k, v := range keydir {
		//fmt.Println("rotation", v.FileId, fd)
		if v.FileId == *fd {
			keydir[k] = makedb.ValueIndex{
				FileId:    oldFd,
				ValueSz:   v.ValueSz,
				ValuePos:  v.ValuePos,
				Timestamp: v.Timestamp,
			}
		}
	}
	*fd, _ = os.OpenFile(makedb.ActiveFilePath, os.O_APPEND|os.O_CREATE, 0644)
}

func get(key string) makedb.Entry {
	val, ok := keydir[key]
	if !ok {
		fmt.Println("not exist")
		return makedb.Entry{}
	}
	data := make([]byte, val.ValueSz)
	_, err := val.FileId.ReadAt(data, val.ValuePos)
	if err != nil {
		log.Println(err)
	}
	e := makedb.NewEntry()
	e.Unmarshal(data)
	fmt.Println(e.Value)
	return e
}

func set(fd *os.File, key string, value string) {
	entry := makedb.NewEntryAndCalc(key, value)
	fi, _ := fd.Stat()
	//写入文件是多加一个分隔符\n，为了扫描文件时进行读取数据，然后进行反序列化
	n, err := fd.Write(append(entry.Marshal(), makedb.Separator))
	if err != nil {
		log.Fatal(err)
	}
	//实时刷盘
	fd.Sync()
	keydir[key] = makedb.ValueIndex{
		FileId: fd,
		//去掉Separator长度
		ValueSz:   int64(n - 1),
		ValuePos:  fi.Size(),
		Timestamp: time.Now().Unix(),
	}
	fmt.Println("OK")
}

func del(key string) {
	delete(keydir, key)
}

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}
