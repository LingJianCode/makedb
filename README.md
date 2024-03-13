# makedb

## The purpose of this project
To learn database system.

## What is makedb?
makedb is a persisent k-v database that based on bitcask.bitcask's paper is [here](https://riak.com/assets/bitcask-intro.pdf). 

## Task
- [x] Get/Put k/v with string.
- [x] Support http protocol.
- [ ] Write ahead log.
- [ ] Rotate active file.
- [ ] Support redis protocol.
- [ ] Get/Put k/v with list.

## usage
- start server
```bash
go run main.go
```
- put key
```bash
curl -X PUT http://127.0.0.1:8080/key/value
```
- get key
```bash
curl -X GET http://127.0.0.1:8080/key
```