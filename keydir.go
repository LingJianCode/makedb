package makedb

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const StorageFilePreFix = "store"
const Separator = '\n'
const ActiveFileName = "store.gob"
const StoragePath = "./data"
const ActiveFilePath = StoragePath + "/" + ActiveFileName
const MergeFilePreFix = "merge"
const HintFilePreFIx = "hint"

//HashRecord[key]value
type KeyDir map[string]ValueIndex

type ValueIndex struct {
	FileId    *os.File
	ValueSz   int64
	ValuePos  int64
	Timestamp int64
}

func InitKeyDir(keydir *KeyDir, fileDir string) *os.File {
	fileInfos, err := ioutil.ReadDir(fileDir)
	if err != nil {
		fmt.Println(err)
	}

	var activeFileId *os.File
	for _, fi := range fileInfos {
		filename := fi.Name()
		fmt.Println(filename)
		if strings.HasPrefix(filename, StorageFilePreFix) {
			var fd *os.File
			var fdErr error
			if filename == ActiveFileName {
				fd, fdErr = os.OpenFile(fileDir+"/"+filename, os.O_APPEND|os.O_CREATE, 0644)
				activeFileId = fd
			} else {
				fd, fdErr = os.Open(fileDir + "/" + filename)
			}
			if fdErr != nil {
				log.Fatal(fdErr)
			}
			r := bufio.NewReader(fd)
			var pos int64 = 0
			for {
				entry := NewEntry()
				dataByte, dataByteErr := r.ReadBytes(Separator)
				if dataByteErr != nil {
					if dataByteErr == io.EOF {
						break
					}
					log.Fatal(dataByteErr)
				}
				dataByteLen := len(dataByte)
				recodeSz := int64(dataByteLen)
				//去掉Separator长度
				valueSz := dataByteLen - 1
				entry.Unmarshal(dataByte[:valueSz])
				val, ok := (*keydir)[entry.Key]
				var isNewest bool = false
				if !ok {
					isNewest = true
				} else {
					if val.Timestamp < entry.Timestamp {
						isNewest = true
					}
				}
				fmt.Println(entry)
				if isNewest {
					(*keydir)[entry.Key] = ValueIndex{
						FileId:    fd,
						ValueSz:   int64(valueSz),
						ValuePos:  pos,
						Timestamp: entry.Timestamp,
					}
				}
				pos += recodeSz
			}
		} else if strings.HasPrefix(filename, MergeFilePreFix) {
			fd, fdErr := os.Open(fileDir + "/" + filename)
			if fdErr != nil {
				log.Fatal(fdErr)
			}

			hintFd, hintFdErr := os.Open(fileDir + "/" + HintFilePreFIx + "_" + filename)
			if hintFdErr != nil {
				log.Fatal(hintFdErr)
			}
			r := bufio.NewReader(hintFd)
			for {
				hint := NewHint()
				dataByte, dataByteErr := r.ReadBytes(Separator)
				if dataByteErr != nil {
					if dataByteErr == io.EOF {
						break
					}
					log.Fatal(dataByteErr)
				}
				hint.Unmarshal(dataByte[:len(dataByte)-1])
				fmt.Println(hint)
				(*keydir)[hint.Key] = ValueIndex{
					FileId:    fd,
					ValueSz:   hint.ValueSz,
					ValuePos:  hint.ValuePos,
					Timestamp: hint.Timestamp,
				}
			}
		}
	}
	return activeFileId
}

//func (k *KeyDir) Marshal() []byte  {
//	buf := new(bytes.Buffer)
//	enc := gob.NewEncoder(buf)
//	err := enc.Encode(k)
//	if err != nil{
//		log.Println(err)
//		return nil
//	}
//	return buf.Bytes()
//}
//
//func (k *KeyDir) Unmarshal(i []byte) {
//	dec := gob.NewDecoder(bytes.NewBuffer(i))
//	err := dec.Decode(&k)
//	if err != nil {
//		if err == io.EOF {
//			log.Println(err)
//		}else {
//			log.Fatal(err)
//		}
//	}
//}
