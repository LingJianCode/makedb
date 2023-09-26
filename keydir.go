package makedb

type KeydirElement struct {
	DataFile *DataFile
	// ValuePos is entry's offset in file.
	ValuePos int64
	// ValueSize is entry's length in disk.
	ValueSize uint32
	Tstamp    uint64
}

func NewKeydirElement(df *DataFile, vp int64, vs uint32, ts uint64) *KeydirElement {
	return &KeydirElement{
		DataFile:  df,
		ValuePos:  vp,
		ValueSize: vs,
		Tstamp:    ts,
	}
}
