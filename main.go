package main

import (
	"flag"
	"makedb/initialize"
	"makedb/server"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "./makedb.yml", "config file")
	initialize.InitViper(configFile)
	server.StartServer()
}
