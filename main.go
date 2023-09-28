package main

import (
	"flag"
	"makedb/conf"
	"makedb/server"

	"github.com/spf13/viper"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "./makedb.yml", "config file")
	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		panic("read config file error.")
	}
	var m conf.Makedb
	err = viper.UnmarshalKey("makedb", &m)
	if err != nil {
		panic("UnmarshalKey config file error.")
	}
	server.StartServer(&m)
}
