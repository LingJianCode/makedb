package makedb

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Hint struct {
	Timestamp int64
	Ksz       int64
	ValueSz   int64
	ValuePos  int64
	Key       string
}

func NewHint() Hint {
	return Hint{}
}

func (h *Hint) Marshal() []byte {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(h)
	if err != nil {
		log.Println(err)
		return nil
	}
	return buf.Bytes()
}

func (h *Hint) Unmarshal(i []byte) {
	dec := gob.NewDecoder(bytes.NewBuffer(i))
	err := dec.Decode(&h)
	if err != nil {
		log.Fatal(err)
	}
}
