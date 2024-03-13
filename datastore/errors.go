package datastore

import "errors"

var (
	ErrKeyNotExist      = errors.New("key does not exist in database.")
	ErrKeyCheckSumWrong = errors.New("entry's checksum is wrong, maybe entry be modified illegally.")
)
