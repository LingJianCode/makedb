package main

import (
	"flag"
	"fmt"
	"makedb/global"
	"makedb/initialize"
	"makedb/server"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "./makedb.yml", "config file")
	initialize.InitViper(configFile)
	global.MAKEDB_LOG = initialize.Zap()
	defer global.MAKEDB_LOG.Sync()

	s := server.StartServer()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	global.MAKEDB_LOG.Info("shutting down server...")
	if err := s.Ds.Close(); err != nil {
		global.MAKEDB_LOG.Error(fmt.Sprintf("error closing datastore: %v", err))
	}
	global.MAKEDB_LOG.Info("server stopped")
}
