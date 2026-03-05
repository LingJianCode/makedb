# makedb

## The purpose of this project
To learn database system.

## What is makedb?
makedb is a persistent k-v database that based on bitcask. Bitcask's paper is [here](https://riak.com/assets/bitcask-intro.pdf).

## Task
- [x] Get/Put k/v with string.
- [x] Support http protocol.
- [ ] Write ahead log.
- [ ] Rotate active file.
- [ ] Support redis protocol.
- [ ] Get/Put k/v with list.

## Usage

### Start server
```bash
go run main.go
```

Or use custom config file:
```bash
go run main.go -config ./makedb.yml
```

### Put key
```bash
curl -X PUT http://127.0.0.1:8080/key/value
```
Response:
```json
{"status":"success"}
```

### Get key
```bash
curl -X GET http://127.0.0.1:8080/key
```
Response:
```json
{"status":"success","key":"key","value":"value"}
```

## Configuration
Edit `makedb.yml` to configure:
- `http_port`: Server port (default: 8080)
- `data_path`: Data storage directory (default: ./data)
- `log_level`: Log level (debug/info/warn/error, default: info)
- `log_format`: Log format (console/json, default: console)