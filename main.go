package main

import (
	"flag"
	"makedb/global"
	"makedb/initialize"
	"makedb/server"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "./makedb.yml", "config file")
	initialize.InitViper(configFile)
	global.MAKEDB_LOG = initialize.Zap()
	defer global.MAKEDB_LOG.Sync()
	server.StartServer()
}
